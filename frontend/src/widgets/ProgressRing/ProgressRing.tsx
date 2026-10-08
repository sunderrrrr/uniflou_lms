import { useEffect, useState } from 'react';
import styles from './ProgressRing.module.css';

interface ProgressRingProps {
  percent: number;
  size?: 'lg' | 'sm';
  label?: string;
}

export function ProgressRing({ percent, size = 'lg', label = 'до цели' }: ProgressRingProps) {
  const [value, setValue] = useState(0);
  const dim = size === 'lg' ? 180 : 80;
  const stroke = 8;
  const r = (dim - stroke) / 2;
  const c = 2 * Math.PI * r;
  const offset = c - (value / 100) * c;

  useEffect(() => {
    const t = setTimeout(() => setValue(Math.min(100, Math.max(0, percent))), 50);
    return () => clearTimeout(t);
  }, [percent]);

  return (
    <div className={styles.wrap} style={{ width: dim, height: dim }}>
      <svg width={dim} height={dim} viewBox={`0 0 ${dim} ${dim}`}>
        <circle cx={dim / 2} cy={dim / 2} r={r} stroke="var(--line)" strokeWidth={stroke} fill="none" />
        <circle
          cx={dim / 2}
          cy={dim / 2}
          r={r}
          stroke="var(--accent)"
          strokeWidth={stroke}
          fill="none"
          strokeLinecap="round"
          strokeDasharray={c}
          strokeDashoffset={offset}
          transform={`rotate(-90 ${dim / 2} ${dim / 2})`}
          className={styles.ring}
        />
      </svg>
      <div className={styles.center}>
        <strong>{Math.round(value)}%</strong>
        <span>{label}</span>
      </div>
    </div>
  );
}
