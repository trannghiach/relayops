package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"relayops/internal/broker/message"
)

type Consumer struct {
	db *pgxpool.Pool
	nc *nats.Conn
}

func NewConsumer(db *pgxpool.Pool, nc *nats.Conn) *Consumer {
	return &Consumer{
		db: db,
		nc: nc,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	_, err := c.nc.Subscribe("events.>", func(msg *nats.Msg) {
		var evt message.EventCreatedMessage
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			return
		}

		jobs, err := PlanJobs(evt.EventID, evt.EventType)
		if err != nil {
			log.Printf("failed to plan jobs for event %s: %v", evt.EventID, err)
			c.markEventFailed(ctx, evt.EventID)
			return
		}

		tx, err := c.db.Begin(ctx)
		if err != nil {
			log.Printf("failed to begin tx for event %s: %v", evt.EventID, err)
			c.markEventFailed(ctx, evt.EventID)
			return
		}
		defer tx.Rollback(ctx)

		for _, job := range jobs {
			_, err := tx.Exec(ctx, `
				INSERT INTO jobs (
					id, event_id, job_type, channel, status, payload, 
					available_at, attempts, max_attempts, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			`,
				job.ID,
				job.EventID,
				job.JobType,
				job.Channel,
				job.Status,
				job.Payload,
				job.AvailableAt,
				job.Attempts,
				job.MaxAttempts,
				job.CreatedAt,
				job.UpdatedAt,
			)
			if err != nil {
				log.Printf("failed to insert job %s: %v", job.ID, err)
				c.markEventFailed(ctx, evt.EventID)
				return
			}
		}

		_, err = tx.Exec(ctx, `
				UPDATE events
				SET status = 'planned'
				WHERE id = $1
			`, evt.EventID)
		if err != nil {
			log.Printf("failed to update event %s status: %v", evt.EventID, err)
			c.markEventFailed(ctx, evt.EventID)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			log.Printf("failed to commit tx for event %s: %v", evt.EventID, err)
			c.markEventFailed(ctx, evt.EventID)
			return
		}

		log.Printf("planned event successfully: event_id=%s jobs=%d", evt.EventID, len(jobs))
	})

	return err
}

func (c *Consumer) markEventFailed(ctx context.Context, eventID string) {
	_, err := c.db.Exec(ctx, `
		UPDATE events
		SET status = 'failed'
		WHERE id = $1
	`, eventID)
	if err != nil {
		log.Printf("failed to mark event %s as failed: %v", eventID, err)
	}
}
