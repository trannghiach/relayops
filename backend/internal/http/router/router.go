package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"relayops/internal/http/handlers"
)

func NewRouter(db *pgxpool.Pool, nc *nats.Conn) http.Handler {
	r := chi.NewRouter()

	h := handlers.NewEventHandler(db, nc)

	r.Get("/events", handlers.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/events", h.CreateEvent)
	})

	return r
}
