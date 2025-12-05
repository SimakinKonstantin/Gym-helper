import { defineStore } from 'pinia';
import {
  getTrainings,
  createTraining,
  updateTraining,
  removeTraining,
  getTraining,
} from '@/api/trainings';
import { useAuthStore } from './auth';

export const useTrainingsStore = defineStore('trainings', {
  state: () => ({
    items: [],
    current: null,
    loading: false,
    error: null,
  }),
  actions: {
    async fetchTrainings() {
      this.loading = true;
      try {
        const trainings = await getTrainings();
        this.items = Array.isArray(trainings) ? trainings : [];
      } catch (error) {
        this.error = error;
      } finally {
        this.loading = false;
      }
    },
    async fetchTraining(id) {
      this.current = (await getTraining(id)) ?? null;
    },
    async create(payload) {
      const auth = useAuthStore();
      await createTraining({
        ...payload,
        user_login: auth.login,
      });
      await this.fetchTrainings();
    },
    async update(id, payload) {
      const auth = useAuthStore();
      await updateTraining(id, {
        ...payload,
        user_login: auth.login,
      });
      await this.fetchTrainings();
    },
    async remove(id) {
      await removeTraining(id);
      this.items = this.items.filter((training) => training.id !== id);
    },
  },
});

