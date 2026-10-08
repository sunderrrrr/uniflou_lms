import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface AIReport {
  id: string;
  response: string;
  context_type: string;
  created_at: string;
}

export async function analyzeDashboard() {
  return apiFetch<AIReport>(endpoints.aiDashboard, { method: "POST" });
}

export async function analyzePlan() {
  return apiFetch<AIReport>(endpoints.aiPlan, { method: "POST" });
}

export async function analyzeProgress() {
  return apiFetch<AIReport>(endpoints.aiProgress, { method: "POST" });
}

export async function analyzeTaskAttempt(attemptId: string) {
  return apiFetch<AIReport>(endpoints.aiTask(attemptId), { method: "POST" });
}

export async function analyzeTrainerSession(input: {
  attempt_ids: string[];
  subject: string;
  task_number: string;
}) {
  return apiFetch<AIReport>(endpoints.aiTrainerSession, {
    method: "POST",
    body: JSON.stringify(input),
  });
}
