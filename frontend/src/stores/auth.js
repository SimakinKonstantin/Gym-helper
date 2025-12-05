import { defineStore } from 'pinia';
import { signIn, signUp } from '@/api/auth';

const STORAGE_KEY = 'workout-app:user';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    login: null,
    token: null,
    loading: false,
    error: null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    authHeader: (state) => {
      const headers = {};
      if (state.login) {
        headers['X-User-Login'] = state.login;
        headers.login = state.login;
      }
      if (state.token) {
        headers.Authorization = `Bearer ${state.token}`;
      }
      return headers;
    },
  },
  actions: {
    hydrate() {
      const persisted = localStorage.getItem(STORAGE_KEY);
      if (!persisted) return;
      try {
        const parsed = JSON.parse(persisted);
        this.login = parsed.login ?? null;
        this.token = parsed.token ?? null;
      } catch {
        this.clearStorage();
      }
    },
    persist() {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          login: this.login,
          token: this.token,
        }),
      );
    },
    setSession({ login, token }) {
      this.login = login;
      this.token = token;
      this.persist();
    },
    async signIn(payload) {
      this.loading = true;
      this.error = null;
      try {
        const { token } = await signIn(payload);
        this.setSession({ login: payload.login, token });
      } catch (err) {
        // Если статус 401, показываем специальное сообщение
        if (err?.response?.status === 401) {
          this.error = 'Неверный логин или пароль';
        } else {
          this.error = err?.response?.data ?? 'Не удалось войти';
        }
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async signUp(payload) {
      this.loading = true;
      this.error = null;
      try {
        await signUp(payload);
        await this.signIn(payload);
      } catch (err) {
        this.error = err?.response?.data ?? 'Не удалось зарегистрироваться';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    logout() {
      this.login = null;
      this.token = null;
      this.error = null;
      this.clearStorage();
    },
    clearStorage() {
      localStorage.removeItem(STORAGE_KEY);
    },
  },
});

