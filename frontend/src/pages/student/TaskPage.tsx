import { useEffect, useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { analyzeTaskAttempt } from '@/features/ai/api';
import { fetchTask, submitAttempt } from '@/features/tasks/api';
import { AIAnalysisBlock } from '@/widgets/AIAnalysisBlock/AIAnalysisBlock';
import { Button } from '@/shared/ui/Button/Button';
import { Card } from '@/shared/ui/Card/Card';
import { Input } from '@/shared/ui/Input/Input';
import { Skeleton } from '@/shared/ui/Skeleton/Skeleton';
import styles from './TaskPage.module.css';

function parseOptions(raw: unknown): string[] {
  if (Array.isArray(raw)) return raw.map(String);
  if (typeof raw === 'string') {
    try {
      const parsed = JSON.parse(raw) as unknown;
      if (Array.isArray(parsed)) return parsed.map(String);
    } catch {
      return [];
    }
  }
  return [];
}

export function TaskPage() {
  const { id = '' } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const exam = Boolean((location.state as { exam?: boolean } | null)?.exam);
  const [answer, setAnswer] = useState('');
  const [started] = useState(Date.now());
  const [result, setResult] = useState<Awaited<ReturnType<typeof submitAttempt>> | null>(null);
  const [timer, setTimer] = useState(exam ? 60 * 180 : 0);

  const taskQ = useQuery({ queryKey: ['task', id], queryFn: () => fetchTask(id), enabled: Boolean(id) });

  useEffect(() => {
    if (!exam || timer <= 0) return;
    const t = setInterval(() => setTimer((s) => Math.max(0, s - 1)), 1000);
    return () => clearInterval(t);
  }, [exam, timer]);

  const attemptMut = useMutation({
    mutationFn: () => submitAttempt(id, answer, Math.floor((Date.now() - started) / 1000)),
    onSuccess: (data) => setResult(data),
  });

  const options = parseOptions(taskQ.data?.options);

  return (
    <div className={styles.page}>
      {exam ? <p className={styles.timer}>Таймер: {Math.floor(timer / 60)}:{String(timer % 60).padStart(2, '0')}</p> : null}
      <Card index={0}>
        {taskQ.isLoading ? (
          <Skeleton height={120} />
        ) : taskQ.data ? (
          <>
            <h1>Задание</h1>
            <p className={styles.question}>{taskQ.data.question}</p>
            {taskQ.data.type === 'single_choice' ? (
              <div className={styles.options}>
                {options.map((opt) => (
                  <label key={opt} className={styles.option}>
                    <input
                      type="radio"
                      name="answer"
                      value={opt}
                      checked={answer === opt}
                      onChange={() => setAnswer(opt)}
                    />
                    {opt}
                  </label>
                ))}
              </div>
            ) : (
              <Input label="Ответ" value={answer} onChange={(e) => setAnswer(e.target.value)} />
            )}
            <div className={styles.actions}>
              <Button onClick={() => attemptMut.mutate()} disabled={!answer || attemptMut.isPending}>
                Проверить
              </Button>
              <Button variant="ghost" onClick={() => navigate('/student/trainer')}>
                Пропустить
              </Button>
            </div>
          </>
        ) : null}
      </Card>

      {result ? (
        <Card index={1} className={result.is_correct ? styles.ok : styles.bad}>
          <h3>{result.is_correct ? 'Верно!' : 'Неверно'}</h3>
          <p>{result.explanation}</p>
          <AIAnalysisBlock
            title="Разбор через ИИ"
            onAnalyze={() => analyzeTaskAttempt(result.attempt.id)}
          />
        </Card>
      ) : null}
    </div>
  );
}
