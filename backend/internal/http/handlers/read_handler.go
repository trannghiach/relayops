package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"relayops/internal/http/dto"
)

type ReadHandler struct {
	db *pgxpool.Pool
}

func NewReadHandler(db *pgxpool.Pool) *ReadHandler {
	return &ReadHandler{db: db}
}

func parseLimitOffset(r *http.Request) (int, int) {
	limit := 20
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	return limit, offset
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *ReadHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)

	rows, err := h.db.Query(r.Context(), `
		SELECT id, event_type, source, status, created_at
		FROM events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	items := []dto.EventListItem{}

	for rows.Next() {
		var item dto.EventListItem
		if err := rows.Scan(&item.ID, &item.EventType, &item.Source, &item.Status, &item.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, dto.ListResponse[dto.EventListItem]{
		Data:   items,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *ReadHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item dto.EventDetail

	err := h.db.QueryRow(r.Context(), `
		SELECT id, event_type, source, tenant_id, aggregate_type, aggregate_id, 
		payload, metadata, status, created_at
		FROM events
		WHERE id = $1
	`, id).Scan(
		&item.ID,
		&item.EventType,
		&item.Source,
		&item.TenantID,
		&item.AggregateType,
		&item.AggregateID,
		&item.Payload,
		&item.Metadata,
		&item.Status,
		&item.CreatedAt,
	)

	if err != nil {
		writeError(w, http.StatusNotFound, "event not found")
		return
	}

	writeJSON(w, http.StatusOK, dto.Response[dto.EventDetail]{
		Data: item,
	})
}

func (h *ReadHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)

	rows, err := h.db.Query(r.Context(), `
		SELECT id, event_id, job_type, channel, status, 
				attempts, max_attempts, available_at, created_at
		FROM jobs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	items := []dto.JobListItem{}

	for rows.Next() {
		var item dto.JobListItem
		if err := rows.Scan(
			&item.ID,
			&item.EventID,
			&item.JobType,
			&item.Channel,
			&item.Status,
			&item.Attempts,
			&item.MaxAttempts,
			&item.AvailableAt,
			&item.CreatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, dto.ListResponse[dto.JobListItem]{
		Data:   items,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *ReadHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var item dto.JobDetail

	err := h.db.QueryRow(r.Context(), `
		SELECT id, event_id, job_type, channel, status, payload,
				attempts, max_attempts, last_error, available_at, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`, id).Scan(
		&item.ID,
		&item.EventID,
		&item.JobType,
		&item.Channel,
		&item.Status,
		&item.Payload,
		&item.Attempts,
		&item.MaxAttempts,
		&item.LastError,
		&item.AvailableAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	writeJSON(w, http.StatusOK, dto.Response[dto.JobDetail]{
		Data: item,
	})
}

func (h *ReadHandler) ListJobAttempts(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	limit, offset := parseLimitOffset(r)

	rows, err := h.db.Query(r.Context(), `
		SELECT id, job_id, attempt_no, provider, status, 
				response_code, response_body, error_message, 
				started_at, finished_at
		FROM delivery_attempts
		WHERE job_id = $1
		ORDER BY attempt_no ASC
		LIMIT $2 OFFSET $3
	`, jobID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	items := []dto.DeliveryAttemptItem{}

	for rows.Next() {
		var item dto.DeliveryAttemptItem
		if err := rows.Scan(
			&item.ID,
			&item.JobID,
			&item.AttemptNo,
			&item.Provider,
			&item.Status,
			&item.ReponseCode,
			&item.ResponseBody,
			&item.ErrorMessage,
			&item.StartedAt,
			&item.FinishedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, dto.ListResponse[dto.DeliveryAttemptItem]{
		Data:   items,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *ReadHandler) ListDeadLetters(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)

	rows, err := h.db.Query(r.Context(), `
		SELECT dl.id, dl.job_id, dl.reason, dl.payload_snapshot, dl.created_at
		FROM dead_letters dl
		JOIN jobs j ON j.id = dl.job_id
		WHERE j.status = 'dead_lettered'
		ORDER BY dl.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	items := []dto.DeadLetterItem{}

	for rows.Next() {
		var item dto.DeadLetterItem
		if err := rows.Scan(
			&item.ID,
			&item.JobID,
			&item.Reason,
			&item.PayloadSnapshot,
			&item.CreatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, dto.ListResponse[dto.DeadLetterItem]{
		Data:   items,
		Limit:  limit,
		Offset: offset,
	})
}

// temp
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, dto.ErrorResponse{
		Error: message,
	})
}
