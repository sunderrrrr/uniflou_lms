import { useMutation, useQuery } from "@tanstack/react-query";
import {
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { analyzeDashboard } from "@/features/ai/api";
import { fetchNextTask } from "@/features/tasks/api";
import { fetchDashboard } from "@/features/progress/api";
import { fetchDeadlines } from "@/features/plan/api";
import { requestTutor } from "@/features/tutor/api";
import { AIAnalysisBlock } from "@/widgets/AIAnalysisBlock/AIAnalysisBlock";
import { ProgressRing } from "@/widgets/ProgressRing/ProgressRing";
import { TaskCard } from "@/widgets/TaskCard/TaskCard";
import { Card } from "@/shared/ui/Card/Card";
import { Button } from "@/shared/ui/Button/Button";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import styles from "./StudentPages.module.css";

export function DashboardPage() {
  const requestMut = useMutation({
    mutationFn: () =>
      requestTutor(undefined, "Прошу связаться по поводу подготовки."),
  });
  const dashQ = useQuery({ queryKey: ["dashboard"], queryFn: fetchDashboard });
  const deadlinesQ = useQuery({
    queryKey: ["deadlines"],
    queryFn: fetchDeadlines,
  });
  const nextQ = useQuery({
    queryKey: ["next-task", "weak"],
    queryFn: () => fetchNextTask("weak"),
  });

  const goal = dashQ.data?.goal_percent ?? 0;
  const weekly = dashQ.data?.weekly_scores ?? [];
  const nextDeadline = deadlinesQ.data?.[0];

  return (
    <div>
      <div className={styles.hero}>
        <div>
          <p className={styles.accent}>сегодня к цели</p>
          <h1>Дашборд ученика</h1>
        </div>
        {dashQ.isLoading ? (
          <Skeleton height={180} width="180px" />
        ) : (
          <ProgressRing percent={goal} />
        )}
      </div>

      <div className={styles.grid}>
        <Card index={0}>
          <h3>Сегодня</h3>
          {nextQ.isLoading ? (
            <Skeleton height={60} />
          ) : nextQ.data ? (
            <TaskCard
              index={1}
              title="Задание из слабой темы"
              subtitle={nextQ.data.question.slice(0, 80)}
              taskId={nextQ.data.id}
            />
          ) : (
            <p>Пока нет доступных заданий.</p>
          )}
        </Card>
        <Card index={1} className={styles.chartWrap}>
          <h3>Баллы за 8 недель</h3>
          {dashQ.isLoading ? (
            <Skeleton height={180} />
          ) : weekly.length === 0 ? (
            <p>Пока нет данных о решённых заданиях.</p>
          ) : (
            <ResponsiveContainer width="100%" height={180}>
              <LineChart data={weekly}>
                <XAxis
                  dataKey="week"
                  tick={{ fontSize: 12, fill: "var(--ink-faint)" }}
                />
                <YAxis hide domain={[0, 100]} />
                <Tooltip />
                <Line
                  type="monotone"
                  dataKey="score"
                  stroke="var(--accent)"
                  strokeWidth={2.5}
                  dot={{ r: 4, fill: "var(--accent)" }}
                />
              </LineChart>
            </ResponsiveContainer>
          )}
        </Card>
        <Card index={2}>
          <h3>Ближайший дедлайн</h3>
          {deadlinesQ.isLoading ? (
            <Skeleton height={40} />
          ) : nextDeadline ? (
            <>
              <strong>{nextDeadline.title}</strong>
              <p>{new Date(nextDeadline.date).toLocaleDateString("ru-RU")}</p>
            </>
          ) : (
            <p className={styles.accent}>Пока нет дедлайнов</p>
          )}
          <Button
            variant="ghost"
            size="sm"
            style={{ marginTop: 12 }}
            onClick={() => requestMut.mutate()}
            disabled={requestMut.isPending || requestMut.isSuccess}
          >
            {requestMut.isPending
              ? "Отправляем…"
              : requestMut.isSuccess
                ? "Запрос отправлен"
                : "Позвать наставника"}
          </Button>
          {requestMut.isError ? (
            <p role="alert">Не удалось отправить запрос наставнику.</p>
          ) : null}
        </Card>
      </div>

      <AIAnalysisBlock onAnalyze={analyzeDashboard} />
    </div>
  );
}
