import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface AdminUser {
  id: string;
  email: string;
  role: "student" | "parent" | "tutor" | "admin";
  email_verified: boolean;
  created_at: string;
}

export interface AdminTask {
  id: string;
  subject_id: string;
  topic_id: string;
  task_number: number;
  type: "single_choice" | "short_answer";
  question: string;
  options?: string[];
  correct_answer: string;
  explanation: string;
  source: string;
}

export interface AdminSubject {
  id: string;
  name: string;
  code: string;
  topics: { id: string; name: string }[];
}

interface Page<T> {
  items: T[];
  total: number;
  page: number;
}

export interface AdminTaskInput {
  subject_id: string;
  topic_id: string;
  task_number: number;
  type: AdminTask["type"];
  question: string;
  options?: string[];
  correct_answer: string;
  explanation: string;
  source: string;
}

export const fetchAdminUsers = () =>
  apiFetch<Page<AdminUser>>(`${endpoints.adminUsers}?limit=100`);
export const fetchAdminSubjects = () =>
  apiFetch<AdminSubject[]>(endpoints.adminSubjects);
export const fetchAdminTasks = () =>
  apiFetch<Page<AdminTask>>(`${endpoints.adminTasks}?limit=100`);

export function createTutor(email: string, password: string) {
  return apiFetch<{ id: string; email: string; role: string }>(
    endpoints.adminTutors,
    {
      method: "POST",
      body: JSON.stringify({ email, password }),
    },
  );
}

export function updateUserRole(id: string, role: AdminUser["role"]) {
  return apiFetch<void>(endpoints.adminUserRole(id), {
    method: "PATCH",
    body: JSON.stringify({ role }),
  });
}

export function deleteUser(id: string) {
  return apiFetch<void>(endpoints.adminUser(id), { method: "DELETE" });
}

export function saveAdminTask(input: AdminTaskInput, id?: string) {
  return apiFetch<AdminTask>(
    id ? endpoints.adminTask(id) : endpoints.adminTasks,
    {
      method: id ? "PUT" : "POST",
      body: JSON.stringify(input),
    },
  );
}

export function deleteAdminTask(id: string) {
  return apiFetch<void>(endpoints.adminTask(id), {
    method: "DELETE",
  });
}

export function createAdminSubject(input: { name: string; code: string }) {
  return apiFetch<AdminSubject>(endpoints.adminSubjects, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function createAdminTopic(input: { subject_id: string; name: string }) {
  return apiFetch<{ id: string; name: string }>(endpoints.adminTopics, {
    method: "POST",
    body: JSON.stringify(input),
  });
}
