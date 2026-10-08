import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  fetchChildProgress,
  fetchChildReports,
  fetchChildren,
} from "@/features/parent/api";
import { requestTutor } from "@/features/tutor/api";
import { apiFetch } from "@/shared/api/client";
import { endpoints } from "@/shared/api/endpoints";
import { ProgressRing } from "@/widgets/ProgressRing/ProgressRing";
import { DeadlineList } from "@/widgets/DeadlineList/DeadlineList";
import { Card } from "@/shared/ui/Card/Card";
import { Button } from "@/shared/ui/Button/Button";
import { Input } from "@/shared/ui/Input/Input";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import type { Deadline } from "@/features/plan/api";
import styles from "./ParentDashboard.module.css";

export function ParentDashboardPage() {
  const [selected, setSelected] = useState<string | null>(null);
  const [studentEmail, setStudentEmail] = useState("");
  const qc = useQueryClient();

  const childrenQ = useQuery({
    queryKey: ["children"],
    queryFn: fetchChildren,
  });
  const childId = selected ?? childrenQ.data?.[0]?.id ?? null;

  const progressQ = useQuery({
    queryKey: ["child-progress", childId],
    queryFn: () => fetchChildProgress(childId!),
    enabled: Boolean(childId),
  });

  const reportsQ = useQuery({
    queryKey: ["child-reports", childId],
    queryFn: () => fetchChildReports(childId!),
    enabled: Boolean(childId),
  });

  const linkMut = useMutation({
    mutationFn: () =>
      apiFetch(endpoints.linkStudent, {
        method: "POST",
        body: JSON.stringify({ student_email: studentEmail }),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["children"] }),
  });

  const tutorRequestMut = useMutation({
    mutationFn: () =>
      childId
        ? requestTutor(
            childId,
            "Родитель просит связаться по поводу подготовки.",
          )
        : Promise.reject(new Error("Сначала привяжите ученика")),
  });

  const deadlines = (reportsQ.data?.deadlines ?? []) as Deadline[];

  return (
    <div>
      <h1>Дашборд родителя</h1>
      <p className={styles.lead}>
        Спокойный обзор прогресса без лишних деталей
      </p>

      <Card index={0} className={styles.linkCard}>
        <h3>Привязать ученика</h3>
        <div className={styles.linkRow}>
          <Input
            label="Email ученика"
            value={studentEmail}
            onChange={(e) => setStudentEmail(e.target.value)}
          />
          <Button
            onClick={() => linkMut.mutate()}
            disabled={!studentEmail || linkMut.isPending}
          >
            Привязать
          </Button>
        </div>
      </Card>

      {childrenQ.isLoading ? (
        <Skeleton height={80} />
      ) : (
        <div className={styles.tabs}>
          {(childrenQ.data ?? []).map((c) => (
            <Button
              key={c.id}
              variant={childId === c.id ? "primary" : "ghost"}
              size="sm"
              onClick={() => setSelected(c.id)}
            >
              {c.email}
            </Button>
          ))}
        </div>
      )}

      {!childId ? (
        <p>Добавьте ученика по email, чтобы видеть прогресс.</p>
      ) : (
        <div className={styles.grid}>
          <Card index={1}>
            {progressQ.isLoading ? (
              <Skeleton height={160} />
            ) : (
              <>
                <ProgressRing percent={progressQ.data?.goal_percent ?? 0} />
                <p className={styles.summary}>{progressQ.data?.summary}</p>
              </>
            )}
          </Card>
          <Card index={2}>
            <h3>Отчёт за неделю</h3>
            <ul>
              {(reportsQ.data?.weekly_report ?? []).map((line) => (
                <li key={line}>{line}</li>
              ))}
            </ul>
            <Button
              variant="ghost"
              size="sm"
              style={{ marginTop: 12 }}
              onClick={() => tutorRequestMut.mutate()}
              disabled={tutorRequestMut.isPending || tutorRequestMut.isSuccess}
            >
              {tutorRequestMut.isPending
                ? "Отправляем…"
                : tutorRequestMut.isSuccess
                  ? "Запрос отправлен"
                  : "Написать наставнику"}
            </Button>
            {tutorRequestMut.isError ? (
              <p role="alert">Не удалось отправить запрос наставнику.</p>
            ) : null}
          </Card>
          <Card index={3}>
            <h3>Дедлайны</h3>
            <DeadlineList items={deadlines} />
            <h3 style={{ marginTop: 16 }}>Заметки наставника</h3>
            <ul>
              {(
                reportsQ.data?.tutor_notes as { body: string }[] | undefined
              )?.map((n, i) => <li key={i}>{n.body}</li>) ?? null}
            </ul>
          </Card>
        </div>
      )}
    </div>
  );
}
