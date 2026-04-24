package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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
	// 1. claim job
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

	// 2. external call (no tx)
	startedAt := time.Now()

	var execErr error

	switch channel {
	case "email":
		execErr = dispatcher.SendEmailMock(payload)
	default:
		execErr = fmt.Errorf("unsupported channel: %s", channel)
	}

	finishedAt := time.Now()
	attemptNo := attempts + 1

	// 3. db tx (atomic update job + record attempt)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction for job %s: %v", jobID, err)
		return
	}
	defer tx.Rollback(ctx)

	if execErr != nil {
		isDeadLetter := attemptNo >= maxAttempts
		// ------- failure handling-------
		if isDeadLetter {
			// dead letter
			_, err = tx.Exec(ctx, `
				UPDATE jobs
				SET status = 'dead_lettered',
					attempts = $2,
					last_error = $3,
					updated_at = NOW()
				WHERE id = $1
			`, jobID, attemptNo, execErr.Error())
			if err != nil {
				log.Printf("failed to dead-letter job %s: %v", jobID, err)
				return
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO dead_letters (
					id, job_id, reason, payload_snapshot, created_at
				) 
				SELECT $1, id, $2, payload, NOW()
				FROM jobs WHERE id = $3
				ON CONFLICT (job_id) DO NOTHING
			`, uuid.New(), execErr.Error(), jobID)
			if err != nil {
				log.Printf("failed to insert dead letter for job %s: %v", jobID,
					err)
				return
			}
		} else {
			// retry
			backoff := computeBackoff(attemptNo)

			_, err = tx.Exec(ctx, `
				UPDATE jobs
				SET status = 'pending',
					attempts = $2,
					last_error = $3,
					available_at = NOW() + $4 * INTERVAL '1 second',
					updated_at = NOW()
				WHERE id = $1
			`, jobID, attemptNo, execErr.Error(), int(backoff.Seconds()))
			if err != nil {
				log.Printf("failed to reschedule job %s: %v", jobID, err)
				return
			}
		}

		// insert record attempt
		_, err = tx.Exec(ctx, `
			INSERT INTO delivery_attempts (
				id, job_id, attempt_no, provider, status,
				error_message, started_at, finished_at 
			) VALUES ($1, $2, $3, $4, 'failed', $5, $6, $7)
		`,
			uuid.New(),
			jobID,
			attemptNo,
			channel,
			execErr.Error(),
			startedAt,
			finishedAt,
		)
		if err != nil {
			log.Printf("failed to insert job attempt for job %s: %v", jobID,
				err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			log.Printf("failed to commit transaction for job %s: %v", jobID, err)
			return
		}

		if isDeadLetter {
			log.Printf("job %s dead-lettered after %d attempts: %v", jobID, attemptNo, execErr)
		} else {
			log.Printf("job %s failed on attempt %d, will retry after backoff: %v", jobID, attemptNo, execErr)
		}
		return
	}

	// ------- success handling -------
	_, err = tx.Exec(ctx, `
		UPDATE jobs
		SET status = 'succeeded', 
			attempts = $2, 
			last_error = NULL,
			updated_at = NOW()
		WHERE id = $1
	`, jobID, attemptNo)
	if err != nil {
		log.Printf("failed to mark job %s as succeeded: %v", jobID, err)
		return
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO delivery_attempts (
			id, job_id, attempt_no, provider, status,
			started_at, finished_at 
		) VALUES ($1, $2, $3, $4, 'succeeded', $5, $6)
	`,
		uuid.New(),
		jobID,
		attemptNo,
		channel,
		startedAt,
		finishedAt,
	)
	if err != nil {
		log.Printf("failed to insert job attempt for job %s: %v", jobID, err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction for job %s: %v", jobID, err)
		return
	}

	log.Printf("job %s succeeded on attempt %d", jobID, attemptNo)
}

func computeBackoff(attempts int) time.Duration {
	switch attempts {
	case 1:
		return 1 * time.Second
	case 2:
		return 2 * time.Second
	case 3:
		return 3 * time.Second
	default:
		return 4 * time.Second
	}
}
