package handlers

import (
	"net/http"
	"relayops/internal/http/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsHandler struct {
	db *pgxpool.Pool
}

func NewMetricsHandler(db *pgxpool.Pool) *MetricsHandler {
	return &MetricsHandler{db: db}
}

func (h *MetricsHandler) GetMetricsSummary(w http.ResponseWriter, r *http.Request) {
	var res dto.MetricsSummary

	err := h.db.QueryRow(r.Context(), `
		SELECT 
			(SELECT COUNT(*) FROM events) AS total_events,
			(SELECT COUNT(*) FROM jobs WHERE status = 'pending') AS pending_jobs,
			(SELECT COUNT(*) FROM jobs WHERE status = 'succeeded') AS succeeded_jobs,
			(SELECT COUNT(*) FROM jobs WHERE status = 'dead_lettered') AS dead_lettered_jobs,
			(SELECT COUNT(*) FROM delivery_attempts) AS total_attempts
	`).Scan(
		&res.TotalEvents,
		&res.PendingJobs,
		&res.SucceededJobs,
		&res.DeadLetteredJobs,
		&res.TotalAttempts,
	)
	if err != nil {
		http.Error(w, "failed to query metrics", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, dto.Response[dto.MetricsSummary]{
		Data: res,
	})
}
