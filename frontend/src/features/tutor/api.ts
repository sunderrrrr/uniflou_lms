import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";

export interface TutorStudent {
  id: string;
  email: string;
  goal_percent: number;
}

export interface TutorRequest {
  id: string;
  student_id: string;
  student_email: string;
  requester_email: string;
  message: string;
  created_at: string;
}

export async function fetchTutorStudents() {
  return apiFetch<TutorStudent[]>(endpoints.tutorStudents);
}

export async function fetchTutorStudent(id: string) {
  return apiFetch<{
    student: { id: string; email: string };
    profile: {
      university_name: string;
      specialty_name: string;
      pass_score: number;
    };
    progress: number;
    weak: unknown[];
    comments: { body: string; needs_call: boolean; created_at: string }[];
  }>(endpoints.tutorStudent(id));
}

export function setTutorStudentTarget(
  id: string,
  input: { university: string; specialty: string; pass_score: number },
) {
  return apiFetch(endpoints.tutorTarget(id), {
    method: "PUT",
    body: JSON.stringify(input),
  });
}

export async function postTutorComment(
  id: string,
  body: string,
  needsCall: boolean,
) {
  return apiFetch(endpoints.tutorComment(id), {
    method: "POST",
    body: JSON.stringify({ body, needs_call: needsCall }),
  });
}

export function requestTutor(studentId?: string, message?: string) {
  return apiFetch<void>(endpoints.tutorRequest, {
    method: "POST",
    body: JSON.stringify({ student_id: studentId, message }),
  });
}

export function fetchTutorRequests() {
  return apiFetch<TutorRequest[]>(endpoints.tutorRequests);
}

export function markTutorRequestRead(id: string) {
  return apiFetch<void>(endpoints.tutorRequestRead(id), { method: "PATCH" });
}
