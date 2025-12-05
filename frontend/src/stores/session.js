import { defineStore } from 'pinia';

const SESSION_KEY = 'workout-app:session';

const toNumber = (value, fallback = 0) => {
  const result = Number(value);
  return Number.isFinite(result) ? result : fallback;
};

const pickCalories = (set) => {
  if (typeof set?.cal_per_set === 'number') return set.cal_per_set;
  if (typeof set?.calPerSet === 'number') return set.calPerSet;
  if (typeof set?.calories_per_set === 'number')
    return set.calories_per_set;
  return 0;
};

const createSetId = () =>
  `set-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

const normalizeSet = (set, exerciseFallback = null) => {
  const originalWeight = toNumber(
    set.originalWeight ?? set.weight ?? set.plannedWeight,
    0,
  );
  const originalReps = toNumber(set.originalReps ?? set.reps, 0);
  
  // Сначала пытаемся получить cal_per_set из set, если нет - из exercise (для обратной совместимости)
  let calPerSet = pickCalories(set);
  if (calPerSet === 0 && exerciseFallback) {
    if (typeof exerciseFallback?.cal_per_set === 'number') {
      calPerSet = exerciseFallback.cal_per_set;
    } else if (typeof exerciseFallback?.calPerSet === 'number') {
      calPerSet = exerciseFallback.calPerSet;
    }
  }

  return {
    uid: set.uid ?? createSetId(),
    originalWeight,
    weight: toNumber(set.weight ?? set.plannedWeight, originalWeight),
    originalReps,
    reps: toNumber(set.reps ?? set.completedReps, originalReps),
    completed: Boolean(set.completed),
    cal_per_set: calPerSet,
  };
};

const normalizeExercise = (exercise) => ({
  name: exercise.name,
  sets: (exercise.sets ?? []).map((set) => normalizeSet(set, exercise)),
});

const normalizeTraining = (training) => {
  if (!training) return null;
  return {
    id: training.id,
    name: training.name,
    exercises: (training.exercises ?? []).map((exercise) =>
      normalizeExercise(exercise),
    ),
  };
};

const cloneExercise = (exercise) =>
  JSON.parse(JSON.stringify(exercise));

export const useSessionStore = defineStore('session', {
  state: () => ({
    activeTraining: null,
    originalTraining: null,
    startedAt: null,
    finishedAt: null,
  }),
  getters: {
    isActive: (state) => Boolean(state.activeTraining),
    completionRate: (state) => {
      if (!state.activeTraining) return 0;
      const { exercises } = state.activeTraining;
      const totalSets = exercises.reduce(
        (acc, exercise) => acc + exercise.sets.length,
        0,
      );
      if (totalSets === 0) return 0;
      const completed = exercises.reduce(
        (acc, exercise) =>
          acc + exercise.sets.filter((set) => set.completed).length,
        0,
      );
      return Math.round((completed / totalSets) * 100);
    },
    hasExerciseDeviation: (state) => (exerciseIndex) => {
      if (!state.activeTraining || !state.originalTraining) return false;
      const active = state.activeTraining.exercises?.[exerciseIndex];
      const original = state.originalTraining.exercises?.[exerciseIndex];
      if (!active || !original) return false;
      if (active.removed) return true;
      if (active.sets.length !== original.sets.length) return true;
      return active.sets.some((set, index) => {
        const originalSet = original.sets[index];
        if (!originalSet) return true;
        return (
          toNumber(set.weight) !== toNumber(originalSet.weight) ||
          toNumber(set.reps) !== toNumber(originalSet.reps)
        );
      });
    },
  },
  actions: {
    hydrate() {
      const raw = localStorage.getItem(SESSION_KEY);
      if (!raw) return;
      try {
        const parsed = JSON.parse(raw);
        this.activeTraining = normalizeTraining(parsed.activeTraining);
        this.originalTraining = normalizeTraining(parsed.originalTraining);
        this.startedAt = parsed.startedAt ?? null;
        this.finishedAt = parsed.finishedAt ?? null;
      } catch {
        this.clear();
      }
    },
    persist() {
      localStorage.setItem(
        SESSION_KEY,
        JSON.stringify({
          activeTraining: this.activeTraining,
          originalTraining: this.originalTraining,
          startedAt: this.startedAt,
          finishedAt: this.finishedAt,
        }),
      );
    },
    startSession(training) {
      const normalized = normalizeTraining(training);
      if (!normalized) {
        this.clear();
        return;
      }
      this.activeTraining = JSON.parse(JSON.stringify(normalized));
      this.originalTraining = JSON.parse(JSON.stringify(normalized));
      this.startedAt = new Date().toISOString();
      this.finishedAt = null;
      this.persist();
    },
    updateSet(exerciseIndex, setIndex, payload) {
      if (!this.activeTraining) return;
      const exercise = this.activeTraining.exercises?.[exerciseIndex];
      if (!exercise) return;
      const set = exercise.sets?.[setIndex];
      if (!set) return;
      if (payload.weight !== undefined) {
        set.weight = toNumber(payload.weight, set.weight);
      }
      if (payload.reps !== undefined) {
        set.reps = Math.max(0, toNumber(payload.reps, set.reps));
      }
      this.persist();
    },
    toggleSetCompletion(exerciseIndex, setIndex) {
      if (!this.activeTraining) return;
      const exercise = this.activeTraining.exercises?.[exerciseIndex];
      if (!exercise) return;
      const set = exercise.sets?.[setIndex];
      if (!set) return;
      set.completed = !set.completed;
      this.persist();
    },
    addSet(exerciseIndex) {
      if (!this.activeTraining) return;
      const exercise = this.activeTraining.exercises?.[exerciseIndex];
      if (!exercise) return;
      const sample = exercise.sets[exercise.sets.length - 1] ?? {
        weight: 0,
        reps: 0,
        originalReps: 0,
        cal_per_set: 0,
      };
      const newSet = normalizeSet({
        weight: sample.weight,
        reps: sample.reps,
        originalReps: sample.originalReps ?? sample.reps,
        cal_per_set: sample.cal_per_set ?? 0,
      });
      exercise.sets.push(newSet);
      this.persist();
    },
    removeSet(exerciseIndex, setIndex) {
      if (!this.activeTraining) return;
      const exercise = this.activeTraining.exercises?.[exerciseIndex];
      if (!exercise) return;
      if (exercise.sets.length <= 1) {
        exercise.removed = true;
      } else {
        exercise.sets.splice(setIndex, 1);
      }
      this.persist();
    },
    resetExercise(exerciseIndex) {
      if (!this.activeTraining || !this.originalTraining) return;
      const originalExercise =
        this.originalTraining.exercises?.[exerciseIndex];
      if (!originalExercise) return;
      this.activeTraining.exercises[exerciseIndex] = cloneExercise(
        originalExercise,
      );
      this.persist();
    },
    completeExercise(exerciseIndex) {
      if (!this.activeTraining) return;
      const exercise = this.activeTraining.exercises?.[exerciseIndex];
      if (!exercise) return;
      exercise.sets = exercise.sets.map((set) => ({
        ...set,
        completed: true,
      }));
      this.persist();
    },
    buildResultPayload({ userLogin, finishTime }) {
      if (!this.activeTraining || !this.startedAt) return null;
      const finish = finishTime ?? new Date().toISOString();
      return {
        training_id: this.activeTraining.id,
        user_login: userLogin,
        start_time: this.startedAt,
        finish_time: finish,
        result_values: this.activeTraining.exercises
          .filter((exercise) => !exercise.removed)
          .map((exercise) => ({
            name: exercise.name,
            sets: exercise.sets.map((set) => ({
              weight: toNumber(set.weight, 0),
              original_reps: toNumber(set.originalReps, 0),
              real_reps: toNumber(set.reps, toNumber(set.originalReps, 0)),
              cal_per_set: toNumber(set.cal_per_set, 0),
            })),
          })),
      };
    },
    clear() {
      this.activeTraining = null;
      this.originalTraining = null;
      this.startedAt = null;
      this.finishedAt = null;
      localStorage.removeItem(SESSION_KEY);
    },
  },
});

