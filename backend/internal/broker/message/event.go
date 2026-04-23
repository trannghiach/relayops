package message

type EventCreatedMessage struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
}
