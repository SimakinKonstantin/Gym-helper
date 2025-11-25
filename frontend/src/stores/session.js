import { defineStore } from 'pinia';

const SESSION_KEY = 'workout-app:session';

export const useSessionStore = defineStore('session', {
  state: () => ({
    activeTraining: null,
    completedSets: {},
  }),
  getters: {
    isActive: (state) => Boolean(state.activeTraining),
    completionRate: (state) => {
      if (!state.activeTraining) return 0;
      const totalSets = state.activeTraining.exercises.reduce(
        (acc, exercise) => acc + exercise.sets.length,
        0,
      );
      const completed = Object.keys(state.completedSets).length;
      return totalSets === 0 ? 0 : Math.round((completed / totalSets) * 100);
    },
  },
  actions: {
    hydrate() {
      const raw = localStorage.getItem(SESSION_KEY);
      if (!raw) return;
      try {
        const parsed = JSON.parse(raw);
        this.activeTraining = parsed.activeTraining;
        this.completedSets = parsed.completedSets ?? {};
      } catch {
        this.clear();
      }
    },
    persist() {
      localStorage.setItem(
        SESSION_KEY,
        JSON.stringify({
          activeTraining: this.activeTraining,
          completedSets: this.completedSets,
        }),
      );
    },
    startSession(training) {
      this.activeTraining = training;
      this.completedSets = {};
      this.persist();
    },
    toggleSetCompletion(exerciseIndex, setIndex) {
      const key = `${exerciseIndex}-${setIndex}`;
      if (this.completedSets[key]) {
        delete this.completedSets[key];
      } else {
        this.completedSets[key] = true;
      }
      this.persist();
    },
    clear() {
      this.activeTraining = null;
      this.completedSets = {};
      localStorage.removeItem(SESSION_KEY);
    },
  },
});

