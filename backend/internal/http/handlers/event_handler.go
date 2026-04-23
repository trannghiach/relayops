package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"relayops/internal/broker/message"
	"relayops/internal/http/dto"
)

type EventHandler struct {
	db *pgxpool.Pool
	nc *nats.Conn
}

func NewEventHandler(db *pgxpool.Pool, nc *nats.Conn) *EventHandler {
	return &EventHandler{db: db, nc: nc}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	id := uuid.New()

	payloadBytes, _ := json.Marshal(req.Payload)
	metadataBytes, _ := json.Marshal(req.Metadata)

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO events (id, event_type, source, tenant_id, aggregate_type, aggregate_id, payload, metadata, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'received')
	`, id, req.EventType, req.Source, req.TenantID, req.AggregateType, req.AggregateID, payloadBytes, metadataBytes)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// publish message
	msg := message.EventCreatedMessage{
		EventID:   id.String(),
		EventType: req.EventType,
	}

	msgBytes, _ := json.Marshal(msg)
	h.nc.Publish("events."+req.EventType, msgBytes)

	resp := dto.CreateEventResponse{
		ID:        id.String(),
		Status:    "received",
		CreatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
