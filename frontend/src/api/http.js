import axios from 'axios';
import { storeToRefs } from 'pinia';
import { useAuthStore } from '@/stores/auth';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api';

const http = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10_000,
});

http.interceptors.request.use((config) => {
  const authStore = useAuthStore();
  const { login, token } = storeToRefs(authStore);
  if (login.value) {
    config.headers['X-User-Login'] = login.value;
    config.headers['User-Login-Id'] = login.value;
    config.headers.login = login.value;
  }
  if (token.value) {
    config.headers.Authorization = `Bearer ${token.value}`;
  }
  return config;
});

export default http;

