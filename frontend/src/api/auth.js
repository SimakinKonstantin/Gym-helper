import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api';

const authClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10_000,
});

export const signIn = async ({ login, password }) => {
  const { data } = await authClient.post('/sign-in', { login, password });
  return data;
};

export const signUp = async ({ login, password }) => {
  await authClient.post('/sign-up', { login, password });
};

