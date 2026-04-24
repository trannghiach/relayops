package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"

	"relayops/internal/config"
	"relayops/internal/http/handlers"
)

func NewRouter(db *pgxpool.Pool, nc *nats.Conn, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Demo-Key"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	eventHandler := handlers.NewEventHandler(db, nc)
	readHandler := handlers.NewReadHandler(db)
	controlHandler := handlers.NewControlHandler(db)
	metricsHandler := handlers.NewMetricsHandler(db)
	demoHandler := handlers.NewDemoHandler(
		db,
		nc,
		cfg.DemoEnabled,
		cfg.DemoAPIKey,
		cfg.DemoMaxEvents,
	)

	r.Get("/health", handlers.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/events", eventHandler.CreateEvent)
		r.Post("/demo/events", demoHandler.CreateDemoEvents)

		r.Get("/events", readHandler.ListEvents)
		r.Get("/events/{id}", readHandler.GetEvent)

		r.Get("/jobs", readHandler.ListJobs)
		r.Get("/jobs/{id}", readHandler.GetJob)
		r.Get("/jobs/{id}/attempts", readHandler.ListJobAttempts)

		r.Get("/dead-letters", readHandler.ListDeadLetters)

		r.Post("/jobs/{id}/retry", controlHandler.RetryJob)
		r.Post("/jobs/{id}/replay", controlHandler.ReplayJob)

		r.Get("/metrics/summary", metricsHandler.GetMetricsSummary)
	})

	return r
}
