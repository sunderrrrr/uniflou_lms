import styles from './Toast.module.css';

interface ToastProps {
  message: string;
  tone?: 'info' | 'error' | 'success';
}

export function Toast({ message, tone = 'info' }: ToastProps) {
  return <div className={`${styles.toast} ${styles[tone]}`}>{message}</div>;
}
