import type { InputHTMLAttributes } from 'react';
import styles from './Input.module.css';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export function Input({ label, error, id, className = '', ...rest }: InputProps) {
  const inputId = id ?? rest.name;
  return (
    <label className={`${styles.wrap} ${className}`} htmlFor={inputId}>
      {label ? <span className={styles.label}>{label}</span> : null}
      <input id={inputId} className={styles.input} {...rest} />
      {error ? <span className={styles.error}>{error}</span> : null}
    </label>
  );
}
