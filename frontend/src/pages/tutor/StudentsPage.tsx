import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import {
  fetchTutorRequests,
  fetchTutorStudents,
  markTutorRequestRead,
} from "@/features/tutor/api";
import { Card } from "@/shared/ui/Card/Card";
import { ProgressRing } from "@/widgets/ProgressRing/ProgressRing";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import styles from "./TutorPages.module.css";

export function StudentsPage() {
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["tutor-students"],
    queryFn: fetchTutorStudents,
  });
  const requestsQ = useQuery({
    queryKey: ["tutor-requests"],
    queryFn: fetchTutorRequests,
  });
  const readMut = useMutation({
    mutationFn: markTutorRequestRead,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tutor-requests"] }),
  });

  return (
    <div>
      <h1>Ученики</h1>
      <section className={styles.requests}>
        <div className={styles.requestsHead}>
          <h2>Обращения</h2>
          <span>{requestsQ.data?.length ?? 0} новых</span>
        </div>
        {requestsQ.isLoading ? (
          <Skeleton height={72} />
        ) : requestsQ.isError ? (
          <p role="alert">Не удалось загрузить обращения.</p>
        ) : requestsQ.data?.length ? (
          requestsQ.data.map((request) => (
            <article className={styles.request} key={request.id}>
              <div>
                <Link to={`/tutor/students/${request.student_id}`}>
                  {request.student_email}
                </Link>
                <span>
                  {request.requester_email} ·{" "}
                  {new Date(request.created_at).toLocaleString("ru-RU")}
                </span>
                <p>{request.message}</p>
              </div>
              <button
                type="button"
                disabled={readMut.isPending}
                onClick={() => readMut.mutate(request.id)}
              >
                Отметить прочитанным
              </button>
            </article>
          ))
        ) : (
          <p className={styles.meta}>Новых обращений нет.</p>
        )}
      </section>
      <div className={styles.grid}>
        {q.isLoading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} height={140} />
          ))
        ) : (q.data ?? []).length === 0 ? (
          <p>Пока нет учеников.</p>
        ) : (
          q.data?.map((s, i) => (
            <Link key={s.id} to={`/tutor/students/${s.id}`}>
              <Card index={i} className={styles.card}>
                <h3>{s.email}</h3>
                <ProgressRing percent={s.goal_percent} size="sm" />
              </Card>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
