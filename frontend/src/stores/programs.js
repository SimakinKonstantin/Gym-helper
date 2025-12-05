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
        const programs = await getPrograms();
        this.items = Array.isArray(programs) ? programs : [];
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
      // Получаем информацию о программе из списка, если она там есть
      const programFromList = this.items.find((p) => p.id === id);
      if (programFromList) {
        this.current = programFromList;
      } else {
        // Если программы нет в списке, пытаемся получить отдельно
        this.current = (await getProgram(id)) ?? null;
      }
      // Получаем расписание тренировок по дням
      const schedule = await getProgramTrainings(id);
      this.currentTrainings = Array.isArray(schedule) ? schedule : [];
    },
    async addTraining({ programId, day, trainingId }) {
      await attachTrainingToProgram({
        program_id: programId,
        day,
        training_id: trainingId ?? null,
      });
      await this.fetchProgram(programId);
    },
    async removeTraining({ programId, day, trainingId }) {
      await detachTrainingFromProgram({
        program_id: programId,
        day,
        training_id: trainingId ?? null,
      });
      await this.fetchProgram(programId);
    },
  },
});

