import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export type UserRole = 'student' | 'parent' | 'tutor' | 'admin';

export interface AuthUser {
  id: string;
  email: string;
  role: UserRole;
}

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthUser | null;
  setSession: (tokens: { access_token: string; refresh_token: string }, user: AuthUser) => void;
  setAccessToken: (token: string) => void;
  clear: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      user: null,
      setSession: (tokens, user) =>
        set({
          accessToken: tokens.access_token,
          refreshToken: tokens.refresh_token,
          user,
        }),
      setAccessToken: (token) => set({ accessToken: token }),
      clear: () => set({ accessToken: null, refreshToken: null, user: null }),
    }),
    { name: 'uniflou-auth' },
  ),
);
