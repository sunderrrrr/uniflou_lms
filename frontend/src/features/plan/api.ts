import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface University {
  id: string;
  name: string;
  city?: string;
  target_scores?: Record<string, number>;
}

export interface StudentPlanRow {
  id: string;
  priority: string;
  university: University;
}

export interface Deadline {
  id: string;
  title: string;
  date: string;
  type?: string;
}

export async function fetchPlan() {
  return apiFetch<{ plans: StudentPlanRow[]; deadlines: Deadline[] }>(
    endpoints.plan,
  );
}

export async function fetchUniversities() {
  return apiFetch<University[]>(endpoints.planUniversities);
}

export async function addUniversityToPlan(
  universityId: string,
  priority: string,
) {
  return apiFetch(endpoints.planUniversities, {
    method: "POST",
    body: JSON.stringify({ university_id: universityId, priority }),
  });
}

export async function removePlanUniversity(planId: string) {
  return apiFetch<void>(`${endpoints.planUniversities}/${planId}`, {
    method: "DELETE",
  });
}

export async function fetchDeadlines() {
  return apiFetch<Deadline[]>(endpoints.planDeadlines);
}

export function addDeadline(input: {
  title: string;
  date: string;
  type: string;
}) {
  return apiFetch<Deadline>(endpoints.planDeadlines, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function removeDeadline(id: string) {
  return apiFetch<void>(`${endpoints.planDeadlines}/${id}`, {
    method: "DELETE",
  });
}
