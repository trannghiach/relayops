package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"relayops/internal/dispatcher"
)

type Runner struct {
	db *pgxpool.Pool
}

func NewRunner(db *pgxpool.Pool) *Runner {
	return &Runner{db: db}
}

func (r *Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.runOnce(ctx)
		case <-ctx.Done():
			log.Println("runner shutting down")
			return
		}
	}
}

func (r *Runner) runOnce(ctx context.Context) {
	rows, err := r.db.Query(ctx, `
		SELECT id, channel, payload, attempts, max_attempts
		FROM jobs
		WHERE status = 'pending' 
		AND available_at <= NOW()
		LIMIT 10
	`)
	if err != nil {
		log.Printf("failed to query pending jobs: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var jobID string
		var channel string
		var payload []byte
		var attempts int
		var maxAttempts int

		if err := rows.Scan(&jobID, &channel, &payload, &attempts, &maxAttempts); err != nil {
			log.Printf("failed to scan job row: %v", err)
			continue
		}

		r.processJob(ctx, jobID, channel, payload, attempts, maxAttempts)
	}
}

func (r *Runner) processJob(
	ctx context.Context,
	jobID string,
	channel string,
	payload []byte,
	attempts int,
	maxAttempts int,
) {
	cmdTag, err := r.db.Exec(ctx, `
		UPDATE jobs
		SET status = 'processing', claimed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`, jobID)
	if err != nil {
		log.Printf("failed to claim job %s: %v", jobID, err)
		return
	}
	if cmdTag.RowsAffected() == 0 {
		return
	}

	log.Printf("processing job: id=%s channel=%s attempts=%d/%d", jobID, channel, attempts, maxAttempts)

	var execErr error
	switch channel {
	case "email":
		execErr = dispatcher.SendEmailMock(payload)
	default:
		execErr = fmt.Errorf("unsupported channel: %s", channel)
	}

	if execErr != nil {
		r.handleFailure(ctx, jobID, attempts, maxAttempts, execErr)
		return
	}

	_, err = r.db.Exec(ctx, `
		UPDATE jobs
		SET status = 'completed', updated_at = NOW()
		WHERE id = $1
	`, jobID)
	if err != nil {
		log.Printf("failed to mark job %s as completed: %v", jobID, err)
	}

	log.Printf("completed job: id=%s", jobID)
}

func (r *Runner) handleFailure(
	ctx context.Context,
	jobID string,
	attempts int,
	maxAttempts int,
	execErr error,
) {
	nextAttempts := attempts + 1

	if nextAttempts >= maxAttempts {
		_, err := r.db.Exec(ctx, `
			UPDATE jobs
			SET status = 'dead_lettered',
				attempts = $2,
				last_error = $3,
				updated_at = NOW()
				WHERE id = $1
		`, jobID, nextAttempts, execErr.Error())
		if err != nil {
			log.Printf("failed to dead-letter job %s: %v", jobID, err)
			return
		}

		log.Printf("job dead-lettered: id=%s attempts=%d error=%v", jobID, nextAttempts, execErr)
		return
	}

	backoff := computeBackoff(nextAttempts)

	_, err := r.db.Exec(ctx, `
		UPDATE jobs
		SET status = 'pending',
			attempts = $2,
			last_error = $3,
			available_at = NOW() + $4 * INTERVAL '1 second',
			updated_at = NOW()
		WHERE id = $1
	`, jobID, nextAttempts, execErr.Error(), int(backoff.Seconds()))
	if err != nil {
		log.Printf("failed to reschedule job %s: %v", jobID, err)
		return
	}

	log.Printf("job rescheduled: id=%s attempts=%d error=%v next_in=%s", jobID, nextAttempts, execErr, backoff)
}

func computeBackoff(attempts int) time.Duration {
	switch attempts {
	case 1:
		return 10 * time.Second
	case 2:
		return 30 * time.Second
	case 3:
		return 60 * time.Second
	default:
		return 5 * time.Minute
	}
}
