package dto

type MetricsSummary struct {
	TotalEvents      int `json:"total_events"`
	PendingJobs      int `json:"pending_jobs"`
	SucceededJobs    int `json:"succeeded_jobs"`
	DeadLetteredJobs int `json:"dead_lettered_jobs"`
	TotalAttempts    int `json:"total_attempts"`
}
