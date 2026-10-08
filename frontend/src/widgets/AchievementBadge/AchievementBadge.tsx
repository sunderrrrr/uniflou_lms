import styles from './AchievementBadge.module.css';

interface AchievementBadgeProps {
  code: string;
  label: string;
}

export function AchievementBadge({ code, label }: AchievementBadgeProps) {
  return (
    <div className={styles.badge} title={code}>
      <span className={styles.icon}>★</span>
      <span>{label}</span>
    </div>
  );
}
