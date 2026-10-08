import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import {
  fetchTrainerCatalog,
  fetchTrainerTopics,
  fetchNextBatch,
  fetchMockSession,
} from "@/features/tasks/api";
import { Button } from "@/shared/ui/Button/Button";
import styles from "./TrainerPage.module.css";

export function TrainerPage() {
  const navigate = useNavigate();
  const [mode, setMode] = useState<"weak" | "topic" | "mock">("weak");
  const catalog = useQuery({
    queryKey: ["trainer-catalog"],
    queryFn: fetchTrainerCatalog,
  });
  const topicsQuery = useQuery({
    queryKey: ["trainer-topics"],
    queryFn: fetchTrainerTopics,
    enabled: mode === "topic",
  });
  const subjects = catalog.data ?? [];
  const [subjectId, setSubjectId] = useState("");
  const [topicId, setTopicId] = useState("");
  const [quantity, setQuantity] = useState(5);
  const selectedSubject = subjects.find((subject) => subject.id === subjectId);
  const selectedTopic = topicsQuery.data?.find((topic) => topic.id === topicId);

  const start = useMutation({
    mutationFn: () =>
      mode === "mock"
        ? fetchMockSession(subjectId)
        : fetchNextBatch(
            mode,
            quantity,
            mode === "topic" ? topicId : undefined,
          ),
    onSuccess: (tasks) =>
      navigate("/student/trainer/session", {
        state: {
          tasks,
          mode,
          subject:
            mode === "mock"
              ? (selectedSubject?.name ?? "")
              : (selectedTopic?.subject_name ?? ""),
          subjectId: mode === "mock" ? subjectId : "",
          taskNumber:
            mode === "topic"
              ? (selectedTopic?.name ?? "Тема")
              : mode === "weak"
                ? "Персональный подбор"
                : "Пробник",
          quantity,
          durationSeconds: mode === "mock" ? 235 * 60 : undefined,
        },
      }),
  });

  const canStart =
    mode === "weak" ||
    (mode === "topic" ? Boolean(topicId) : Boolean(subjectId));

  return (
    <div className={styles.page}>
      <h1>Тренажёр</h1>
      <p className={styles.lead}>Выберите формат тренировки</p>

      <div className={styles.controls}>
        <div
          className={styles.modes}
          role="group"
          aria-label="Режим тренировки"
        >
          <button
            type="button"
            aria-pressed={mode === "weak"}
            className={mode === "weak" ? styles.modeActive : ""}
            onClick={() => setMode("weak")}
          >
            Персональный подбор
          </button>
          <button
            type="button"
            aria-pressed={mode === "topic"}
            className={mode === "topic" ? styles.modeActive : ""}
            onClick={() => setMode("topic")}
          >
            По теме
          </button>
          <button
            type="button"
            aria-pressed={mode === "mock"}
            className={mode === "mock" ? styles.modeActive : ""}
            onClick={() => setMode("mock")}
          >
            Пробник
          </button>
        </div>

        {mode === "topic" ? (
          <label className={styles.field}>
            <span>Тема</span>
            <select
              value={topicId}
              onChange={(event) => setTopicId(event.target.value)}
            >
              <option value="">Выберите тему</option>
              {topicsQuery.data?.map((topic) => (
                <option value={topic.id} key={topic.id}>
                  {topic.subject_name} · {topic.name} ({topic.count})
                </option>
              ))}
            </select>
          </label>
        ) : null}

        {mode === "mock" ? (
          <label className={styles.field}>
            <span>Предмет</span>
            <select
              value={subjectId}
              onChange={(event) => setSubjectId(event.target.value)}
            >
              <option value="">Выберите предмет</option>
              {subjects.map((subject) => (
                <option value={subject.id} key={subject.id}>
                  {subject.name}
                </option>
              ))}
            </select>
          </label>
        ) : null}

        {mode !== "mock" ? (
          <label className={styles.field}>
            <span>Количество заданий</span>
            <select
              value={quantity}
              onChange={(event) => setQuantity(Number(event.target.value))}
            >
              {[3, 5, 10].map((value) => (
                <option key={value} value={value}>
                  {value}
                </option>
              ))}
            </select>
          </label>
        ) : null}

        <div className={styles.startRow}>
          <Button
            onClick={() => start.mutate()}
            disabled={!canStart || start.isPending}
          >
            {start.isPending
              ? "Подготавливаем…"
              : mode === "mock"
                ? "Начать пробник"
                : "Начать тренировку"}
          </Button>
          {mode === "mock" ? (
            <span>Ограничение времени: 3 ч 55 мин</span>
          ) : null}
        </div>
      </div>

      {catalog.isLoading ? <p role="status">Загружаем каталог…</p> : null}
      {catalog.isError || topicsQuery.isError || start.isError ? (
        <p role="alert">Не удалось загрузить тренировку. Попробуйте ещё раз.</p>
      ) : null}
      {!catalog.isLoading && !catalog.isError && subjects.length === 0 ? (
        <p>В каталоге пока нет заданий.</p>
      ) : null}
    </div>
  );
}
