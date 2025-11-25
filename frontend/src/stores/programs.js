import { defineStore } from 'pinia';
import {
  getPrograms,
  createProgram as apiCreateProgram,
  removeProgram as apiRemoveProgram,
  getProgram,
  getProgramTrainings,
  attachTrainingToProgram,
  detachTrainingFromProgram,
} from '@/api/programs';
import { useAuthStore } from './auth';

export const useProgramsStore = defineStore('programs', {
  state: () => ({
    items: [],
    current: null,
    currentTrainings: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchPrograms() {
      this.loading = true;
      try {
        this.items = await getPrograms();
      } catch (error) {
        this.error = error;
      } finally {
        this.loading = false;
      }
    },
    async createProgram(name) {
      const auth = useAuthStore();
      await apiCreateProgram({ name, user_login: auth.login });
      await this.fetchPrograms();
    },
    async deleteProgram(id) {
      await apiRemoveProgram(id);
      this.items = this.items.filter((program) => program.id !== id);
      if (this.current?.id === id) {
        this.current = null;
        this.currentTrainings = [];
      }
    },
    async fetchProgram(id) {
      this.current = await getProgram(id);
      this.currentTrainings = await getProgramTrainings(id);
    },
    async addTraining({ programId, day, trainingId }) {
      await attachTrainingToProgram({
        program_id: programId,
        day,
        training: trainingId,
      });
      await this.fetchProgram(programId);
    },
    async removeTraining({ programId, day, trainingId }) {
      await detachTrainingFromProgram({
        program_id: programId,
        day,
        training: trainingId,
      });
      await this.fetchProgram(programId);
    },
  },
});

