import styles from './Skeleton.module.css';

interface SkeletonProps {
  height?: number;
  width?: string;
}

export function Skeleton({ height = 16, width = '100%' }: SkeletonProps) {
  return <div className={styles.skeleton} style={{ height, width }} aria-hidden />;
}
