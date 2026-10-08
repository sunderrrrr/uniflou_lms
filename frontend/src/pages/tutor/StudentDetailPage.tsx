import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useParams } from "react-router-dom";
import {
  fetchTutorStudent,
  postTutorComment,
  setTutorStudentTarget,
} from "@/features/tutor/api";
import { Card } from "@/shared/ui/Card/Card";
import { Button } from "@/shared/ui/Button/Button";
import { Input } from "@/shared/ui/Input/Input";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import styles from "./TutorPages.module.css";

export function StudentDetailPage() {
  const { id = "" } = useParams();
  const [body, setBody] = useState("");
  const [needsCall, setNeedsCall] = useState(false);
  const qc = useQueryClient();

  const q = useQuery({
    queryKey: ["tutor-student", id],
    queryFn: () => fetchTutorStudent(id),
    enabled: Boolean(id),
  });

  const commentMut = useMutation({
    mutationFn: () => postTutorComment(id, body, needsCall),
    onSuccess: () => {
      setBody("");
      qc.invalidateQueries({ queryKey: ["tutor-student", id] });
    },
  });

  const targetMut = useMutation({
    mutationFn: (input: {
      university: string;
      specialty: string;
      pass_score: number;
    }) => setTutorStudentTarget(id, input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tutor-student", id] }),
  });

  if (q.isLoading) return <Skeleton height={200} />;

  return (
    <div>
      <h1>{q.data?.student.email}</h1>
      <p className={styles.meta}>Прогресс к цели: {q.data?.progress}%</p>

      <div className={styles.grid}>
        <Card index={0}>
          <h3>Цель подготовки</h3>
          <form
            onSubmit={(event) => {
              event.preventDefault();
              const form = new FormData(event.currentTarget);
              targetMut.mutate({
                university: String(form.get("university") ?? ""),
                specialty: String(form.get("specialty") ?? ""),
                pass_score: Number(form.get("pass_score") ?? 0),
              });
            }}
          >
            <Input
              label="Вуз"
              name="university"
              defaultValue={q.data?.profile.university_name ?? ""}
            />
            <Input
              label="Направление"
              name="specialty"
              defaultValue={q.data?.profile.specialty_name ?? ""}
            />
            <Input
              label="Целевой балл"
              name="pass_score"
              type="number"
              min={0}
              max={100}
              defaultValue={q.data?.profile.pass_score || ""}
            />
            <Button type="submit" disabled={targetMut.isPending}>
              {targetMut.isPending ? "Сохраняем…" : "Сохранить цель"}
            </Button>
          </form>
          {targetMut.isError ? (
            <p role="alert">Не удалось сохранить цель.</p>
          ) : null}
        </Card>
        <Card index={0}>
          <h3>Слабые места</h3>
          <ul>
            {(
              q.data?.weak as { topic_id: string; score: number }[] | undefined
            )?.map((w) => (
              <li key={w.topic_id}>
                {w.topic_id.slice(0, 8)}… — {w.score}%
              </li>
            ))}
          </ul>
        </Card>
        <Card index={1}>
          <h3>Комментарий</h3>
          <Input
            label="Текст"
            value={body}
            onChange={(e) => setBody(e.target.value)}
          />
          <label className={styles.check}>
            <input
              type="checkbox"
              checked={needsCall}
              onChange={(e) => setNeedsCall(e.target.checked)}
            />
            Требуется созвон
          </label>
          <Button
            onClick={() => commentMut.mutate()}
            disabled={!body || commentMut.isPending}
          >
            Сохранить
          </Button>
        </Card>
        <Card index={2}>
          <h3>История комментариев</h3>
          <ul>
            {q.data?.comments.map((c) => (
              <li key={c.created_at}>
                {c.body} {c.needs_call ? "· созвон" : ""}
              </li>
            ))}
          </ul>
        </Card>
      </div>
    </div>
  );
}
