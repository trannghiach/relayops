package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"relayops/internal/http/dto"
)

type ControlHandler struct {
	db *pgxpool.Pool
}

func NewControlHandler(db *pgxpool.Pool) *ControlHandler {
	return &ControlHandler{db: db}
}

func (h *ControlHandler) RetryJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	cmdTag, err := h.db.Exec(r.Context(), `
		UPDATE jobs
		SET status = 'pending', 
			available_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status IN ('failed', 'pending')
	`, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "job not eligible for retry", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, dto.JobActionResponse{
		JobID:  jobID,
		Action: "retried",
	})
}

func (h *ControlHandler) ReplayJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	cmdTag, err := h.db.Exec(r.Context(), `
		UPDATE jobs
		SET status = 'pending', 
			last_error = NULL,
			available_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status = 'dead_lettered'
	`, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "job not in dead letter state for replay", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, dto.JobActionResponse{
		JobID:  jobID,
		Action: "replayed",
	})
}
