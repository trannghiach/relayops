package dto

import (
	"encoding/json"
	"time"
)

type EventListItem struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Source    string    `json:"source"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type EventDetail struct {
	ID            string          `json:"id"`
	EventType     string          `json:"event_type"`
	Source        string          `json:"source"`
	TenantID      *string         `json:"tenant_id,omitempty"`
	AggregateType *string         `json:"aggregate_type,omitempty"`
	AggregateID   *string         `json:"aggregate_id,omitempty"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

type JobListItem struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	JobType     string    `json:"job_type"`
	Channel     string    `json:"channel"`
	Status      string    `json:"status"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
	AvailableAt time.Time `json:"available_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type JobDetail struct {
	ID          string          `json:"id"`
	EventID     string          `json:"event_id"`
	JobType     string          `json:"job_type"`
	Channel     string          `json:"channel"`
	Status      string          `json:"status"`
	Payload     json.RawMessage `json:"payload"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	LastError   *string         `json:"last_error,omitempty"`
	AvailableAt time.Time       `json:"available_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type DeliveryAttemptItem struct {
	ID           string    `json:"id"`
	JobID        string    `json:"job_id"`
	AttemptNo    int       `json:"attempt_no"`
	Provider     string    `json:"provider"`
	Status       string    `json:"status"`
	ReponseCode  *int      `json:"response_code,omitempty"`
	ResponseBody *string   `json:"response_body,omitempty"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
}

type DeadLetterItem struct {
	ID              string          `json:"id"`
	JobID           string          `json:"job_id"`
	Reason          string          `json:"reason"`
	PayloadSnapshot json.RawMessage `json:"payload_snapshot"`
	CreatedAt       time.Time       `json:"created_at"`
}

type Response[T any] struct {
	Data T `json:"data"`
}

type ListResponse[T any] struct {
	Data   []T `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
