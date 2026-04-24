package dto

type CreateDemoEventsRequest struct {
	Count int `json:"count"`
}

type CreateDemoEventsResponse struct {
	Created  int      `json:"created"`
	EventIDs []string `json:"event_ids"`
}
