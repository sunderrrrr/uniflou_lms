import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface WeeklyPoint {
  week: string;
  score: number;
}

export interface DashboardData {
  goal_percent: number;
  today_task_ids: string[];
  weekly_scores: WeeklyPoint[];
  next_deadline?: { title: string; date: string } | null;
}

export interface ProgressRow {
  id: string;
  topic_id: string;
  subject_id: string;
  score: number;
}

export interface MockExamResult {
  id: string;
  subject_id: string;
  title: string;
  total_tasks: number;
  score: number;
  created_at: string;
}

export async function fetchDashboard() {
  return apiFetch<DashboardData>(endpoints.meProgress);
}

export async function fetchProgress() {
  return apiFetch<{ goal_percent: number; items: ProgressRow[] }>(
    endpoints.progress,
  );
}

export async function fetchWeakTopics() {
  return apiFetch<ProgressRow[]>(endpoints.progressWeak);
}

export function saveMockResult(input: {
  subject_id: string;
  title: string;
  total_tasks: number;
  score: number;
}) {
  return apiFetch<MockExamResult>(endpoints.mockResults, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function fetchMockResults() {
  return apiFetch<MockExamResult[]>(endpoints.mockResults);
}
