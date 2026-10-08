import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface Task {
  id: string;
  subject_id: string;
  topic_id: string;
  task_number: number;
  type: "single_choice" | "short_answer";
  question: string;
  options?: string[];
}

export interface TrainerSubject {
  id: string;
  name: string;
  code: string;
  task_numbers: { number: number; count: number }[];
}

export interface TrainerTopic {
  id: string;
  subject_id: string;
  subject_name: string;
  name: string;
  count: number;
}

export interface AttemptResult {
  attempt: { id: string; is_correct: boolean; answer: string };
  is_correct: boolean;
  explanation: string;
}

export async function fetchTrainerCatalog() {
  return apiFetch<TrainerSubject[]>(endpoints.trainerCatalog);
}

export async function fetchTrainerTopics() {
  return apiFetch<TrainerTopic[]>(endpoints.trainerTopics);
}

export async function fetchNextBatch(
  mode: string,
  limit: number,
  topicId?: string,
) {
  const query = new URLSearchParams({ mode, limit: String(limit) });
  if (topicId) query.set("topic_id", topicId);
  return apiFetch<Task[]>(`${endpoints.tasksNextBatch}?${query}`);
}

export async function fetchMockSession(subjectId: string) {
  const query = new URLSearchParams({ subject_id: subjectId });
  return apiFetch<Task[]>(`${endpoints.mockSession}?${query}`);
}

export async function fetchTrainingSession(
  subjectId: string,
  taskNumber: number,
  limit: number,
) {
  const query = new URLSearchParams({
    subject_id: subjectId,
    task_number: String(taskNumber),
    limit: String(limit),
  });
  return apiFetch<Task[]>(`${endpoints.trainerSession}?${query}`);
}

export async function fetchNextTask(mode: string, topicId?: string) {
  const q = new URLSearchParams({ mode });
  if (topicId) q.set("topic_id", topicId);
  return apiFetch<Task>(`${endpoints.tasksNext}?${q}`);
}

export async function fetchTask(id: string) {
  return apiFetch<Task>(endpoints.task(id));
}

export async function submitAttempt(
  taskId: string,
  answer: string,
  timeSpent: number,
) {
  return apiFetch<AttemptResult>(endpoints.taskAttempt(taskId), {
    method: "POST",
    body: JSON.stringify({ answer, time_spent: timeSpent }),
  });
}

export async function fetchTaskHistory(page = 1) {
  return apiFetch<{ items: unknown[]; total: number; page: number }>(
    `${endpoints.tasksHistory}?page=${page}&limit=20`,
  );
}
