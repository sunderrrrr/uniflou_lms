import type { Deadline } from "@/features/plan/api";
import { Button } from "@/shared/ui/Button/Button";
import { Card } from "@/shared/ui/Card/Card";
import styles from "./DeadlineList.module.css";

interface DeadlineListProps {
  items: Deadline[];
  onRemove?: (id: string) => void;
}

export function DeadlineList({ items, onRemove }: DeadlineListProps) {
  if (!items.length) {
    return <p className={styles.empty}>Дедлайнов пока нет</p>;
  }
  return (
    <ul className={styles.list}>
      {items.map((d, i) => (
        <li key={d.id}>
          <Card index={i}>
            <strong>{d.title}</strong>
            <span className={styles.date}>
              {new Date(d.date).toLocaleDateString("ru-RU")}
            </span>
            {onRemove ? (
              <Button variant="ghost" size="sm" onClick={() => onRemove(d.id)}>
                Удалить
              </Button>
            ) : null}
          </Card>
        </li>
      ))}
    </ul>
  );
}
