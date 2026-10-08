import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { analyzePlan } from "@/features/ai/api";
import {
  addDeadline,
  addUniversityToPlan,
  fetchPlan,
  fetchUniversities,
  removeDeadline,
  removePlanUniversity,
} from "@/features/plan/api";
import { AIAnalysisBlock } from "@/widgets/AIAnalysisBlock/AIAnalysisBlock";
import { DeadlineList } from "@/widgets/DeadlineList/DeadlineList";
import {
  UniversityCatalog,
  UniversityList,
} from "@/widgets/UniversityList/UniversityList";
import { Skeleton } from "@/shared/ui/Skeleton/Skeleton";
import styles from "./StudentPages.module.css";

export function PlanPage() {
  const qc = useQueryClient();
  const [deadlineTitle, setDeadlineTitle] = useState("");
  const [deadlineDate, setDeadlineDate] = useState("");
  const [deadlineType, setDeadlineType] = useState("exam");
  const [deadlineError, setDeadlineError] = useState("");
  const planQ = useQuery({ queryKey: ["plan"], queryFn: fetchPlan });
  const catalogQ = useQuery({
    queryKey: ["universities"],
    queryFn: fetchUniversities,
  });

  const addMut = useMutation({
    mutationFn: ({ id, priority }: { id: string; priority: string }) =>
      addUniversityToPlan(id, priority),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plan"] }),
  });

  const removeMut = useMutation({
    mutationFn: (planId: string) => removePlanUniversity(planId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plan"] }),
  });

  const addDeadlineMut = useMutation({
    mutationFn: () =>
      addDeadline({
        title: deadlineTitle.trim(),
        date: deadlineDate,
        type: deadlineType,
      }),
    onSuccess: async () => {
      setDeadlineTitle("");
      setDeadlineDate("");
      setDeadlineError("");
      await qc.invalidateQueries({ queryKey: ["plan"] });
    },
    onError: (error) =>
      setDeadlineError(
        error instanceof Error ? error.message : "Не удалось сохранить дедлайн",
      ),
  });

  const removeDeadlineMut = useMutation({
    mutationFn: removeDeadline,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["plan"] }),
  });

  return (
    <div>
      <div className={styles.hero}>
        <div>
          <p className={styles.accent}>куда идём</p>
          <h1>Мой план</h1>
        </div>
      </div>

      <h2 className={styles.sectionTitle}>Вузы</h2>
      {planQ.isLoading ? (
        <Skeleton height={120} />
      ) : (
        <UniversityList
          plans={planQ.data?.plans ?? []}
          onRemove={(id) => removeMut.mutate(id)}
        />
      )}

      <h2 className={styles.sectionTitle}>Добавить из каталога</h2>
      {catalogQ.isLoading ? (
        <Skeleton height={120} />
      ) : (
        <UniversityCatalog
          catalog={catalogQ.data ?? []}
          onAdd={(id, priority) => addMut.mutate({ id, priority })}
        />
      )}

      <h2 className={styles.sectionTitle}>Дедлайны</h2>
      <DeadlineList
        items={planQ.data?.deadlines ?? []}
        onRemove={(id) => removeDeadlineMut.mutate(id)}
      />
      <form
        className={styles.deadlineForm}
        onSubmit={(event) => {
          event.preventDefault();
          addDeadlineMut.mutate();
        }}
      >
        <label>
          Название
          <input
            required
            maxLength={255}
            value={deadlineTitle}
            onChange={(event) => setDeadlineTitle(event.target.value)}
          />
        </label>
        <label>
          Дата
          <input
            required
            type="date"
            value={deadlineDate}
            onChange={(event) => setDeadlineDate(event.target.value)}
          />
        </label>
        <label>
          Тип
          <select
            value={deadlineType}
            onChange={(event) => setDeadlineType(event.target.value)}
          >
            <option value="exam">Экзамен</option>
            <option value="application">Подача документов</option>
            <option value="other">Другое</option>
          </select>
        </label>
        <button type="submit" disabled={addDeadlineMut.isPending}>
          {addDeadlineMut.isPending ? "Сохраняем…" : "Добавить дедлайн"}
        </button>
      </form>
      {deadlineError ? <p role="alert">{deadlineError}</p> : null}

      <AIAnalysisBlock title="Анализ плана" onAnalyze={analyzePlan} />
    </div>
  );
}
