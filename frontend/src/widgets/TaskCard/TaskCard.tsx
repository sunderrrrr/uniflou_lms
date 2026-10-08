import { Link } from 'react-router-dom';
import { Card } from '@/shared/ui/Card/Card';
import { Button } from '@/shared/ui/Button/Button';
import styles from './TaskCard.module.css';

interface TaskCardProps {
  title: string;
  subtitle?: string;
  taskId?: string;
  index?: number;
}

export function TaskCard({ title, subtitle, taskId, index = 0 }: TaskCardProps) {
  return (
    <Card index={index} className={styles.card}>
      <h3>{title}</h3>
      {subtitle ? <p className={styles.sub}>{subtitle}</p> : null}
      {taskId ? (
        <Link to={`/student/trainer/task/${taskId}`}>
          <Button size="sm">Открыть</Button>
        </Link>
      ) : null}
    </Card>
  );
}
