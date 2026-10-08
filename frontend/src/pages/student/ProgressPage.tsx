import { useQuery } from "@tanstack/react-query";
import {
  Bar,
  BarChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { analyzeProgress } from "@/features/ai/api";
import {
  fetchMockResults,
  fetchProgress,
  fetchWeakTopics,
} from "@/features/progress/api";
import { AIAnalysisBlock } from "@/widgets/AIAnalysisBlock/AIAnalysisBlock";
import { ProgressRing } from "@/widgets/ProgressRing/ProgressRing";
import { Card } from "@/shared/ui/Card/Card";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import styles from "./StudentPages.module.css";

export function ProgressPage() {
  const progressQ = useQuery({
    queryKey: ["progress"],
    queryFn: fetchProgress,
  });
  const weakQ = useQuery({ queryKey: ["weak"], queryFn: fetchWeakTopics });
  const mockResultsQ = useQuery({
    queryKey: ["mock-results"],
    queryFn: fetchMockResults,
  });

  const chartData =
    progressQ.data?.items.map((row, i) => ({
      name: `T${i + 1}`,
      score: row.score,
    })) ?? [];

  return (
    <div>
      <div className={styles.hero}>
        <div>
          <p className={styles.accent}>динамика</p>
          <h1>Прогресс</h1>
        </div>
        <ProgressRing percent={progressQ.data?.goal_percent ?? 0} size="sm" />
      </div>

      <div className={styles.grid2}>
        <Card index={0} className={styles.chartWrap}>
          <h3>По темам</h3>
          {progressQ.isLoading ? (
            <Skeleton height={200} />
          ) : chartData.length === 0 ? (
            <p>Пока нет данных о решённых заданиях.</p>
          ) : (
            <ResponsiveContainer width="100%" height={220}>
              <BarChart data={chartData}>
                <XAxis dataKey="name" tick={{ fontSize: 12 }} />
                <YAxis domain={[0, 100]} />
                <Tooltip />
                <Bar
                  dataKey="score"
                  fill="var(--accent)"
                  radius={[8, 8, 0, 0]}
                />
              </BarChart>
            </ResponsiveContainer>
          )}
        </Card>
        <Card index={1}>
          <h3>Топ слабых тем</h3>
          {weakQ.isLoading ? (
            <Skeleton height={100} />
          ) : (weakQ.data ?? []).length ? (
            <ul>
              {weakQ.data?.map((w) => (
                <li key={w.id}>
                  Тема {w.topic_id.slice(0, 8)}… — {w.score}%
                </li>
              ))}
            </ul>
          ) : (
            <p>Пока нет данных о слабых темах.</p>
          )}
        </Card>
      </div>

      <h2 className={styles.sectionTitle}>Результаты пробников</h2>
      {mockResultsQ.isLoading ? (
        <Skeleton height={72} />
      ) : (mockResultsQ.data ?? []).length ? (
        <div className={styles.mockResults}>
          {mockResultsQ.data?.map((result) => (
            <Card key={result.id}>
              <div>
                <strong>{result.title}</strong>
                <span>
                  {new Date(result.created_at).toLocaleDateString("ru-RU")}
                </span>
              </div>
              <strong>
                {result.score}% <small>· {result.total_tasks} заданий</small>
              </strong>
            </Card>
          ))}
        </div>
      ) : (
        <p>Пробники пока не пройдены.</p>
      )}

      <AIAnalysisBlock title="Глубокий анализ" onAnalyze={analyzeProgress} />
    </div>
  );
}
