import { useEffect, useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useLocation, useNavigate } from "react-router-dom";
import { analyzeTrainerSession } from "@/features/ai/api";
import type { Task, AttemptResult } from "@/features/tasks/api";
import { submitAttempt } from "@/features/tasks/api";
import { saveMockResult } from "@/features/progress/api";
import { AIAnalysisBlock } from "@/widgets/AIAnalysisBlock/AIAnalysisBlock";
import { Button } from "@/shared/ui/Button/Button";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import styles from "./TrainerSessionPage.module.css";

interface SessionState {
  tasks: Task[];
  mode: "weak" | "topic" | "mock";
  subject: string;
  subjectId: string;
  taskNumber: string;
  quantity: number;
  durationSeconds?: number;
}

interface RecordedAnswer {
  result?: AttemptResult;
  timeSpent: number;
  skipped: boolean;
}

function getOptions(value: Task["options"]): string[] {
  if (Array.isArray(value)) return value.map(String);
  if (typeof value === "string") {
    try {
      const parsed: unknown = JSON.parse(value);
      return Array.isArray(parsed) ? parsed.map(String) : [];
    } catch {
      return [];
    }
  }
  return [];
}

export function TrainerSessionPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const session = location.state as SessionState | null;
  const tasks = session?.tasks ?? [];
  const [index, setIndex] = useState(0);
  const [answer, setAnswer] = useState("");
  const [questionStartedAt, setQuestionStartedAt] = useState(Date.now());
  const [answers, setAnswers] = useState<RecordedAnswer[]>([]);
  const [remainingTime, setRemainingTime] = useState(
    session?.durationSeconds ?? 0,
  );
  const [analysisError, setAnalysisError] = useState("");
  const mockResultSaved = useRef(false);
  const current = tasks[index];
  const done =
    tasks.length > 0 &&
    (index >= tasks.length ||
      (session?.mode === "mock" && remainingTime === 0));
  const correctCount = answers.filter((item) => item.result?.is_correct).length;
  const answeredCount = answers.filter((item) => item.result).length;
  const totalTime = answers.reduce((sum, item) => sum + item.timeSpent, 0);
  const options = current ? getOptions(current.options) : [];

  useEffect(() => {
    if (!session || tasks.length === 0)
      navigate("/student/trainer", { replace: true });
  }, [navigate, session, tasks.length]);

  useEffect(() => {
    if (session?.mode !== "mock" || done || remainingTime === 0) return;
    const interval = window.setInterval(
      () => setRemainingTime((value) => Math.max(value - 1, 0)),
      1000,
    );
    return () => window.clearInterval(interval);
  }, [done, remainingTime, session?.mode]);

  const attemptMutation = useMutation({
    mutationFn: async () => {
      const timeSpent = Math.floor((Date.now() - questionStartedAt) / 1000);
      const result = await submitAttempt(current.id, answer.trim(), timeSpent);
      return { result, timeSpent };
    },
    onSuccess: ({ result, timeSpent }) => {
      setAnswers((previous) => [
        ...previous,
        { result, timeSpent, skipped: false },
      ]);
      setAnswer("");
      setQuestionStartedAt(Date.now());
      setIndex((previous) => previous + 1);
    },
  });

  const mockResultMutation = useMutation({ mutationFn: saveMockResult });

  useEffect(() => {
    if (
      !done ||
      session?.mode !== "mock" ||
      !session.subjectId ||
      mockResultSaved.current
    )
      return;
    mockResultSaved.current = true;
    mockResultMutation.mutate({
      subject_id: session.subjectId,
      title: `Пробный вариант: ${session.subject}`,
      total_tasks: tasks.length,
      score: Math.round((correctCount / tasks.length) * 100),
    });
  }, [correctCount, done, mockResultMutation, session, tasks.length]);

  if (!session || (!current && !done)) return null;

  function skipCurrent() {
    setAnswers((previous) => [
      ...previous,
      {
        timeSpent: Math.floor((Date.now() - questionStartedAt) / 1000),
        skipped: true,
      },
    ]);
    setAnswer("");
    setQuestionStartedAt(Date.now());
    setIndex((previous) => previous + 1);
  }

  return (
    <div className={styles.page}>
      <Button
        variant="ghost"
        size="sm"
        onClick={() => navigate("/student/trainer")}
      >
        ← К выбору заданий
      </Button>
      <div className={styles.sessionHead}>
        <div>
          <p>
            {session.subject} <span>/</span> {session.taskNumber}
          </p>
          <h1>{done ? "Результат тренировки" : "Тренировка"}</h1>
        </div>
        {!done ? (
          <div className={styles.sessionCounters}>
            {session.mode === "mock" ? (
              <strong className={styles.timer} aria-label="Осталось времени">
                {Math.floor(remainingTime / 3600)}:
                {String(Math.floor((remainingTime % 3600) / 60)).padStart(
                  2,
                  "0",
                )}
                :{String(remainingTime % 60).padStart(2, "0")}
              </strong>
            ) : null}
            <strong>
              {index + 1} <span>/</span> {tasks.length}
            </strong>
          </div>
        ) : null}
      </div>
      {!done && current ? (
        <Card index={0}>
          <p className={styles.question}>{current.question}</p>
          {current.type === "single_choice" ? (
            <div className={styles.options}>
              {options.map((option) => (
                <label className={styles.option} key={option}>
                  <input
                    type="radio"
                    name="answer"
                    value={option}
                    checked={answer === option}
                    onChange={() => setAnswer(option)}
                  />
                  {option}
                </label>
              ))}
            </div>
          ) : (
            <Input
              label="Ваш ответ"
              value={answer}
              onChange={(event) => setAnswer(event.target.value)}
              autoFocus
            />
          )}
          {attemptMutation.isError ? (
            <p role="alert">Ответ не удалось сохранить. Повторите отправку.</p>
          ) : null}
          <div className={styles.actions}>
            <Button
              onClick={() => attemptMutation.mutate()}
              disabled={!answer.trim() || attemptMutation.isPending}
            >
              {attemptMutation.isPending ? "Проверяем…" : "Ответить"}
            </Button>
            <Button
              variant="ghost"
              onClick={skipCurrent}
              disabled={attemptMutation.isPending}
            >
              Пропустить
            </Button>
          </div>
        </Card>
      ) : null}

      {done ? (
        <>
          <section className={styles.stats} aria-label="Итоги тренировки">
            <div>
              <span>Верных ответов</span>
              <strong>
                {correctCount} из {answeredCount}
              </strong>
            </div>
            <div>
              <span>Точность</span>
              <strong>
                {answeredCount
                  ? Math.round((correctCount / answeredCount) * 100)
                  : 0}
                %
              </strong>
            </div>
            <div>
              <span>Время</span>
              <strong>
                {Math.floor(totalTime / 60)}:
                {String(totalTime % 60).padStart(2, "0")}
              </strong>
            </div>
          </section>
          {session.mode === "mock" && mockResultMutation.isPending ? (
            <p role="status">Сохраняем результат пробника…</p>
          ) : null}
          {session.mode === "mock" && mockResultMutation.isError ? (
            <p role="alert">
              Результат не сохранился. Итоги ответов остаются на экране.
            </p>
          ) : null}
          <div className={styles.review}>
            {answers.map(({ result, skipped }, answerIndex) => (
              <Card
                key={result?.attempt.id ?? `skip-${answerIndex}`}
                index={answerIndex}
              >
                <h3
                  className={
                    skipped
                      ? styles.skipped
                      : result?.is_correct
                        ? styles.correct
                        : styles.incorrect
                  }
                >
                  {answerIndex + 1}.{" "}
                  {skipped
                    ? "Пропущено"
                    : result?.is_correct
                      ? "Верно"
                      : "Нужно повторить"}
                </h3>
                {result ? <p>{result.explanation}</p> : null}
              </Card>
            ))}
          </div>
          {answers.some((item) => item.result) ? (
            <AIAnalysisBlock
              title="Анализ всей тренировки"
              onAnalyze={async () => {
                setAnalysisError("");
                try {
                  return await analyzeTrainerSession({
                    attempt_ids: answers.flatMap((item) =>
                      item.result ? [item.result.attempt.id] : [],
                    ),
                    subject: session.subject,
                    task_number: String(session.taskNumber),
                  });
                } catch (error) {
                  setAnalysisError(
                    error instanceof Error
                      ? error.message
                      : "Не удалось выполнить AI-анализ",
                  );
                  throw error;
                }
              }}
            />
          ) : (
            <p>Для ИИ-анализа нужны сохранённые ответы.</p>
          )}
          {analysisError ? <p role="alert">{analysisError}</p> : null}
          <div className={styles.actions}>
            <Button onClick={() => navigate("/student/trainer")}>
              Новая тренировка
            </Button>
          </div>
        </>
      ) : null}
    </div>
  );
}
