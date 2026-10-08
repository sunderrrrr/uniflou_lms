import { apiFetch } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

export interface ChildUser {
  id: string;
  email: string;
}

export async function fetchChildren() {
  return apiFetch<ChildUser[]>(endpoints.parentChildren);
}

export async function fetchChildProgress(id: string) {
  return apiFetch<{
    goal_percent: number;
    summary: string;
    topics: unknown[];
  }>(endpoints.parentProgress(id));
}

export async function fetchChildReports(id: string) {
  return apiFetch<{
    weekly_report: string[];
    deadlines: unknown[];
    tutor_notes: unknown[];
  }>(endpoints.parentReports(id));
}
