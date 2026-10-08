import { Link } from 'react-router-dom';
import { Button } from '@/shared/ui/Button/Button';

export function NotFoundPage() {
  return (
    <div style={{ minHeight: '100vh', display: 'grid', placeItems: 'center', padding: 24 }}>
      <div style={{ textAlign: 'center' }}>
        <img src="/logo.png" alt="" width={72} height={72} />
        <h1>404</h1>
        <p>Страница не найдена</p>
        <Link to="/">
          <Button>На главную</Button>
        </Link>
      </div>
    </div>
  );
}
