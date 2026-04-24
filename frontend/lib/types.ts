export type ApiResponse<T> = {
  data: T;
};

export type ListResponse<T> = {
  data: T[];
  limit: number;
  offset: number;
};

export type MetricsSummary = {
  total_events: number;
  pending_jobs: number;
  succeeded_jobs: number;
  dead_lettered_jobs: number;
  total_attempts: number;
};

export type EventItem = {
  id: string;
  event_type: string;
  source: string;
  status: string;
  created_at: string;
};

export type JobItem = {
  id: string;
  event_id: string;
  job_type: string;
  channel: string;
  status: string;
  attempts: number;
  max_attempts: number;
  available_at: string;
  created_at: string;
};

export type JobDetail = JobItem & {
  payload: unknown;
  last_error?: string;
  updated_at: string;
};

export type DeliveryAttempt = {
  id: string;
  job_id: string;
  attempt_no: number;
  provider: string;
  status: string;
  response_code?: number;
  response_body?: string;
  error_message?: string;
  started_at: string;
  finished_at: string;
};

export type DeadLetter = {
  id: string;
  job_id: string;
  reason: string;
  payload_snapshot: unknown;
  created_at: string;
};