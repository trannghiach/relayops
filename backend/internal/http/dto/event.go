package dto

import "time"

type CreateEventRequest struct {
	EventType     string                 `json:"event_type"`
	Source        string                 `json:"source"`
	TenantID      string                 `json:"tenant_id"`
	AggregateType string                 `json:"aggregate_type"`
	AggregateID   string                 `json:"aggregate_id"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type CreateEventResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
