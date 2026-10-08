import { apiFetch } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';
import type { AuthUser } from './store';

interface AuthResponse {
  tokens: { access_token: string; refresh_token: string; expires_in: number };
  user: AuthUser;
}

export async function loginRequest(email: string, password: string) {
  return apiFetch<AuthResponse>(endpoints.auth.login, {
    method: 'POST',
    body: JSON.stringify({ email, password }),
    skipAuth: true,
  });
}

export async function registerRequest(email: string, password: string, role: AuthUser['role']) {
  return apiFetch<AuthResponse>(endpoints.auth.register, {
    method: 'POST',
    body: JSON.stringify({ email, password, role }),
    skipAuth: true,
  });
}

export async function resetPasswordRequest(email: string) {
  return apiFetch<{ message: string }>(endpoints.auth.reset, {
    method: 'POST',
    body: JSON.stringify({ email }),
    skipAuth: true,
  });
}

export async function resetPasswordConfirm(token: string, newPassword: string) {
  return apiFetch<{ message: string }>(endpoints.auth.resetConfirm, {
    method: 'POST',
    body: JSON.stringify({ token, new_password: newPassword }),
    skipAuth: true,
  });
}

export async function logoutRequest() {
  return apiFetch<void>(endpoints.auth.logout, { method: 'POST' });
}
