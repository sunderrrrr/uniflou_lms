import type { StudentPlanRow, University } from "@/features/plan/api";
import { Card } from "@/shared/ui/Card/Card";
import { Button } from "@/shared/ui/Button/Button";
import styles from "./UniversityList.module.css";

const priorityLabel: Record<string, string> = {
  dream: "Мечта",
  realistic: "Реально",
  safety: "Страховка",
};

interface UniversityListProps {
  plans: StudentPlanRow[];
  onRemove?: (planId: string) => void;
}

export function UniversityList({ plans, onRemove }: UniversityListProps) {
  if (!plans.length) {
    return <p className={styles.empty}>Добавьте вузы в план</p>;
  }
  return (
    <div className={styles.grid}>
      {plans.map((p, i) => (
        <Card key={p.id} index={i} className={styles.item}>
          <span className={styles.badge}>
            {priorityLabel[p.priority] ?? p.priority}
          </span>
          <h3>{p.university?.name ?? "Вуз"}</h3>
          <p className={styles.city}>{p.university?.city}</p>
          {onRemove ? (
            <Button variant="ghost" size="sm" onClick={() => onRemove(p.id)}>
              Удалить
            </Button>
          ) : null}
        </Card>
      ))}
    </div>
  );
}

interface CatalogProps {
  catalog: University[];
  onAdd: (id: string, priority: string) => void;
}

export function UniversityCatalog({ catalog, onAdd }: CatalogProps) {
  if (!catalog.length) {
    return <p className={styles.empty}>Пока нет доступных вузов.</p>;
  }
  return (
    <div className={styles.grid}>
      {catalog.map((u, i) => (
        <Card key={u.id} index={i + 1} className={styles.item}>
          <h3>{u.name}</h3>
          <p className={styles.city}>{u.city}</p>
          <div className={styles.actions}>
            <Button size="sm" onClick={() => onAdd(u.id, "dream")}>
              Мечта
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => onAdd(u.id, "realistic")}
            >
              Реально
            </Button>
          </div>
        </Card>
      ))}
    </div>
  );
}
