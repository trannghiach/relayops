package domain

import "time"

type Job struct {
	ID          string
	EventID     string
	JobType     string
	Channel     string
	Status      string
	Recipient   *string
	TargetURL   *string
	Payload     []byte
	DedupeKey   *string
	ScheduledAt *time.Time
	AvailableAt time.Time
	Attempts    int
	MaxAttempts int
	LastError   *string
	ClaimedAt   *time.Time
	WorkerID    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
