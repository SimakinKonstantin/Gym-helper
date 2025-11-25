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
        this.items = await getTrainings();
      } catch (error) {
        this.error = error;
      } finally {
        this.loading = false;
      }
    },
    async fetchTraining(id) {
      this.current = await getTraining(id);
    },
    async create(payload) {
      const auth = useAuthStore();
      await createTraining({
        ...payload,
        userLogin: auth.login,
      });
      await this.fetchTrainings();
    },
    async update(id, payload) {
      await updateTraining(id, payload);
      await this.fetchTrainings();
    },
    async remove(id) {
      await removeTraining(id);
      this.items = this.items.filter((training) => training.id !== id);
    },
  },
});

