package domain

import "time"

type Event struct {
	ID             string
	EventType      string
	Source         string
	TenantID       string
	AggregateType  string
	AggregateID    string
	Payload        []byte
	Metadata       []byte
	IdempotencyKey *string
	Status         string
	CreatedAt      time.Time
}
