package worker

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"relayops/internal/domain"
)

func PlanJobs(eventID string, eventType string) ([]domain.Job, error) {
	now := time.Now()

	switch eventType {
	case "user.registered":
		payloadBytes, _ := json.Marshal(map[string]interface{}{
			"template": "welcome_email",
		})

		job := domain.Job{
			ID:          uuid.New().String(),
			EventID:     eventID,
			JobType:     "welcome_email",
			Channel:     "email",
			Status:      "pending",
			Payload:     payloadBytes,
			AvailableAt: now,
			Attempts:    0,
			MaxAttempts: 5,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		return []domain.Job{job}, nil
	default:
		return nil, errors.New("unsupported event type")
	}
}
