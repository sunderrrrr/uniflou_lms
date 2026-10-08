import { useMutation } from '@tanstack/react-query';
import { Button } from '@/shared/ui/Button/Button';
import { Card } from '@/shared/ui/Card/Card';
import styles from './AIAnalysisBlock.module.css';

interface AIAnalysisBlockProps {
  title?: string;
  onAnalyze: () => Promise<{ response: string }>;
}

export function AIAnalysisBlock({ title = 'ИИ-анализ', onAnalyze }: AIAnalysisBlockProps) {
  const mutation = useMutation({ mutationFn: onAnalyze });

  return (
    <Card index={3} className={styles.block}>
      <div className={styles.head}>
        <h3>{title}</h3>
        <Button
          variant="primary"
          size="sm"
          onClick={() => mutation.mutate()}
          disabled={mutation.isPending}
        >
          {mutation.isPending ? 'Анализ…' : 'Проанализировать через ИИ'}
        </Button>
      </div>
      {mutation.data ? (
        <p className={styles.text}>{mutation.data.response}</p>
      ) : (
        <p className={styles.hint}>Получите персональные рекомендации на основе ваших данных.</p>
      )}
      {mutation.isError ? (
        <p className={styles.error}>{mutation.error instanceof Error ? mutation.error.message : 'Ошибка'}</p>
      ) : null}
    </Card>
  );
}
