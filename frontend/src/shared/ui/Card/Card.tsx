import type { ReactNode } from 'react';
import styles from './Card.module.css';

interface CardProps {
  children: ReactNode;
  index?: number;
  className?: string;
}

export function Card({ children, index = 0, className = '' }: CardProps) {
  const odd = index % 2 === 1;
  return (
    <div className={`${styles.card} ${odd ? styles.odd : styles.even} ${className}`}>{children}</div>
  );
}
