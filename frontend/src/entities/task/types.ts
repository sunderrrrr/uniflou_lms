export type TaskType = 'single_choice' | 'short_answer';

export interface TaskEntity {
  id: string;
  type: TaskType;
  question: string;
}
