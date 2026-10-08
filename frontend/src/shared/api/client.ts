import { endpoints } from './endpoints';
import { parseApiError } from './errors';
import { useAuthStore } from '@/features/auth/store';

type RequestOptions = RequestInit & { skipAuth?: boolean };

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = useAuthStore.getState().refreshToken;
  if (!refreshToken) return null;
  const res = await fetch(endpoints.auth.refresh, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  if (!res.ok) {
    useAuthStore.getState().clear();
    return null;
  }
  const data = (await res.json()) as { access_token: string; refresh_token: string };
  useAuthStore.getState().setSession(
    { access_token: data.access_token, refresh_token: data.refresh_token },
    useAuthStore.getState().user!,
  );
  return data.access_token;
}

export async function apiFetch<T>(url: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers);
  if (!headers.has('Content-Type') && options.body) {
    headers.set('Content-Type', 'application/json');
  }
  if (!options.skipAuth) {
    const token = useAuthStore.getState().accessToken;
    if (token) headers.set('Authorization', `Bearer ${token}`);
  }

  let res = await fetch(url, { ...options, headers });

  if (res.status === 401 && !options.skipAuth) {
    const newToken = await refreshAccessToken();
    if (newToken) {
      headers.set('Authorization', `Bearer ${newToken}`);
      res = await fetch(url, { ...options, headers });
    }
  }

  if (!res.ok) {
    await parseApiError(res);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}
