package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"relayops/internal/broker/message"
	"relayops/internal/http/dto"
)

type DemoHandler struct {
	db            *pgxpool.Pool
	nc            *nats.Conn
	demoEnabled   bool
	demoAPIKey    string
	demoMaxEvents int
}

func NewDemoHandler(
	db *pgxpool.Pool,
	nc *nats.Conn,
	demoEnabled bool,
	demoAPIKey string,
	demoMaxEvents int,
) *DemoHandler {
	return &DemoHandler{
		db:            db,
		nc:            nc,
		demoEnabled:   demoEnabled,
		demoAPIKey:    demoAPIKey,
		demoMaxEvents: demoMaxEvents,
	}
}

func (h *DemoHandler) CreateDemoEvents(w http.ResponseWriter, r *http.Request) {
	if !h.demoEnabled {
		writeError(w, http.StatusForbidden, "demo mode is disabled")
		return
	}

	if h.demoAPIKey == "" || r.Header.Get("X-Demo-Key") != h.demoAPIKey {
		writeError(w, http.StatusForbidden, "invalid demo key")
		return
	}

	var req dto.CreateDemoEventsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Count <= 0 {
		req.Count = 10
	}
	if req.Count > h.demoMaxEvents {
		req.Count = h.demoMaxEvents
	}

	eventIDs := make([]string, 0, req.Count)

	for i := 0; i < req.Count; i++ {
		id := uuid.New()

		payloadBytes, _ := json.Marshal(map[string]any{
			"user_id": "demo-user",
			"email":   "demo@example.com",
		})

		metadataBytes, _ := json.Marshal(map[string]any{
			"demo": true,
		})

		_, err := h.db.Exec(r.Context(), `
			INSERT INTO events (
				id, event_type, source, payload, metadata, status
			)
			VALUES ($1,$2,$3,$4,$5,'received')
		`,
			id,
			"user.registered",
			"demo-generator",
			payloadBytes,
			metadataBytes,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		msg := message.EventCreatedMessage{
			EventID:   id.String(),
			EventType: "user.registered",
		}

		msgBytes, _ := json.Marshal(msg)

		if err := h.nc.Publish("events.user.registered", msgBytes); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		eventIDs = append(eventIDs, id.String())
	}

	writeJSON(w, http.StatusCreated, dto.Response[dto.CreateDemoEventsResponse]{
		Data: dto.CreateDemoEventsResponse{
			Created:  len(eventIDs),
			EventIDs: eventIDs,
		},
	})
}
