import type {
  ApiResponse,
  DeadLetter,
  DeliveryAttempt,
  EventItem,
  JobDetail,
  JobItem,
  ListResponse,
  MetricsSummary,
} from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL!;

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed: ${res.status}`);
  }

  return res.json() as Promise<T>;
}

export const api = {
  getMetrics: () =>
    request<ApiResponse<MetricsSummary>>("/metrics/summary"),

  getEvents: () =>
    request<ListResponse<EventItem>>("/events"),

  getJobs: () =>
    request<ListResponse<JobItem>>("/jobs"),

  getJob: (id: string) =>
    request<ApiResponse<JobDetail>>(`/jobs/${id}`),

  getJobAttempts: (id: string) =>
    request<ListResponse<DeliveryAttempt>>(`/jobs/${id}/attempts`),

  getDeadLetters: () =>
    request<ListResponse<DeadLetter>>("/dead-letters"),

  retryJob: (id: string) =>
    request<ApiResponse<{ job_id: string; action: string }>>(`/jobs/${id}/retry`, {
      method: "POST",
    }),

  replayJob: (id: string) =>
    request<ApiResponse<{ job_id: string; action: string }>>(`/jobs/${id}/replay`, {
      method: "POST",
    }),
};