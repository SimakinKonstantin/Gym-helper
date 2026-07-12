<template>
  <section class="stack" v-if="training">
    <article class="card stack">
      <RouterLink to="/app/trainings" class="muted">← Все тренировки</RouterLink>
      <div class="session-header">
        <div>
          <h2>{{ training.name }}</h2>
          <p class="muted">
            {{ training.exercises?.length ?? 0 }} упражнений ·
            {{ session.completionRate }}% выполнено
          </p>
        </div>
        <div v-if="isActiveForCurrent" class="timer-box">
          <span class="muted">Длительность</span>
          <strong>{{ formattedDuration }}</strong>
        </div>
      </div>
      <div class="nav-links">
        <button
          v-if="!isActiveForCurrent"
          :disabled="hasForeignSession"
          type="button"
          @click="start"
        >
          Запустить тренировку
        </button>
        <template v-else>
          <button
            :disabled="isFinishing"
            type="button"
            @click="handleFinish"
          >
            {{ isFinishing ? 'Сохраняем…' : 'Завершить тренировку' }}
          </button>
          <button
            :disabled="isFinishing"
            type="button"
            class="ghost"
            @click="handleCancel"
          >
            Отменить тренировку
          </button>
        </template>
      </div>
      <p v-if="hasForeignSession" class="notice">
        Уже запущена тренировка «{{ session.activeTraining?.name }}». Завершите
        её, чтобы начать новую.
      </p>
    </article>

    <article class="card stack" v-if="isActiveForCurrent">
      <div
        v-for="(exercise, exerciseIndex) in activeExercises"
        :key="`${exercise.name}-${exerciseIndex}`"
        class="card stack exercise-card"
      >
        <header class="exercise-header">
          <div>
            <strong>{{ exercise.name }}</strong>
            <p class="muted">
              {{ exercise.sets.length }} подходов
            </p>
          </div>
          <div
            class="exercise-flags"
            v-if="isExerciseModified(exerciseIndex)"
          >
            <span class="badge badge--warning"
              >Тренировка не соответствует программе</span
            >
            <button
              class="ghost sm"
              type="button"
              @click="resetExercise(exerciseIndex)"
            >
              Вернуть как было
            </button>
          </div>
        </header>
        <div class="stack" v-if="!exercise.removed">
          <div
            v-for="(set, setIndex) in exercise.sets"
            :key="set.uid ?? `${exerciseIndex}-${setIndex}`"
            class="set-row"
          >
            <label class="input-field">
              <span class="muted">Вес, кг</span>
              <input
                type="number"
                min="0"
                step="0.5"
                :value="set.weight"
                @input="handleNumberInput(exerciseIndex, setIndex, 'weight', $event)"
              />
            </label>
            <label class="input-field">
              <span class="muted">Повторы</span>
              <input
                type="number"
                min="0"
                step="1"
                :value="set.reps"
                @input="handleNumberInput(exerciseIndex, setIndex, 'reps', $event)"
              />
            </label>
            <label class="checkbox-inline">
              <input
                type="checkbox"
                :checked="set.completed"
                @change="toggleSet(exerciseIndex, setIndex)"
              />
              Выполнено
            </label>
            <button
              type="button"
              class="ghost sm"
              @click="removeSet(exerciseIndex, setIndex)"
            >
              Удалить подход
            </button>
          </div>
        </div>
        <div v-else class="removed-exercise">
          <p class="muted">Упражнение удалено из этой тренировки.</p>
          <button
            type="button"
            class="ghost sm"
            @click="resetExercise(exerciseIndex)"
          >
            Вернуть упражнение
          </button>
        </div>
        <div class="nav-links wrap">
          <button
            type="button"
            class="ghost"
            @click="addSet(exerciseIndex)"
            :disabled="exercise.removed"
          >
            + Добавить подход
          </button>
          <button
            type="button"
            class="ghost"
            @click="completeExercise(exerciseIndex)"
            :disabled="exercise.removed || exercise.sets.every((set) => set.completed)"
          >
            Отметить упражнение выполненным
          </button>
        </div>
      </div>
      <button
        type="button"
        class="finish-button"
        :disabled="isFinishing"
        @click="handleFinish"
      >
        {{ isFinishing ? 'Сохраняем…' : 'Завершить тренировку' }}
      </button>
    </article>
    <article class="card stack" v-else>
      <p class="muted">
        Запустите тренировку, чтобы редактировать упражнения и отслеживать
        прогресс.
      </p>
    </article>
  </section>
  <p v-else class="muted">Загружаем тренировку…</p>
</template>

<script setup>
import {
  computed,
  onBeforeUnmount,
  onMounted,
  ref,
  watch,
} from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useTrainingsStore } from '@/stores/trainings';
import { useSessionStore } from '@/stores/session';
import { processTrainingStatistics } from '@/api/statistics';
import { useAuthStore } from '@/stores/auth';

const route = useRoute();
const trainingsStore = useTrainingsStore();
const session = useSessionStore();
session.hydrate();

const authStore = useAuthStore();
const { login } = storeToRefs(authStore);

const { current: training } = storeToRefs(trainingsStore);
const trainingId = computed(() => Number(route.params.id));

const load = async () => {
  await trainingsStore.fetchTraining(trainingId.value);
};

onMounted(load);

const isActiveForCurrent = computed(
  () => session.activeTraining?.id === training.value?.id && session.isActive,
);

const hasForeignSession = computed(
  () =>
    session.isActive &&
    session.activeTraining?.id &&
    session.activeTraining.id !== training.value?.id,
);

const activeExercises = computed(() =>
  isActiveForCurrent.value ? session.activeTraining?.exercises ?? [] : [],
);

const elapsedSeconds = ref(0);
let timerId;

const syncElapsed = () => {
  if (!session.startedAt || !isActiveForCurrent.value) {
    elapsedSeconds.value = 0;
    return;
  }
  const start = new Date(session.startedAt).getTime();
  elapsedSeconds.value = Math.max(
    0,
    Math.floor((Date.now() - start) / 1000),
  );
};

watch(
  () => ({
    startedAt: session.startedAt,
    active: isActiveForCurrent.value,
  }),
  ({ startedAt, active }) => {
    if (timerId) {
      clearInterval(timerId);
      timerId = undefined;
    }
    if (startedAt && active) {
      syncElapsed();
      timerId = setInterval(syncElapsed, 1_000);
    } else {
      elapsedSeconds.value = 0;
    }
  },
  { immediate: true },
);

onBeforeUnmount(() => {
  if (timerId) {
    clearInterval(timerId);
  }
});

const formattedDuration = computed(() => {
  const seconds = elapsedSeconds.value;
  const hours = String(Math.floor(seconds / 3_600)).padStart(2, '0');
  const minutes = String(Math.floor((seconds % 3_600) / 60)).padStart(2, '0');
  const secs = String(seconds % 60).padStart(2, '0');
  return `${hours}:${minutes}:${secs}`;
});

const start = () => {
  if (!training.value || hasForeignSession.value) return;
  session.startSession(training.value);
};

const isExerciseModified = (exerciseIndex) =>
  session.hasExerciseDeviation?.(exerciseIndex) ?? false;

const handleNumberInput = (exerciseIndex, setIndex, field, event) => {
  const raw = event?.target?.value;
  const value = Number(raw);
  session.updateSet(exerciseIndex, setIndex, {
    [field]: Number.isFinite(value) ? value : 0,
  });
};

const toggleSet = (exerciseIndex, setIndex) => {
  session.toggleSetCompletion(exerciseIndex, setIndex);
};

const addSet = (exerciseIndex) => {
  session.addSet(exerciseIndex);
};

const removeSet = (exerciseIndex, setIndex) => {
  session.removeSet(exerciseIndex, setIndex);
};

const resetExercise = (exerciseIndex) => {
  session.resetExercise(exerciseIndex);
};

const completeExercise = (exerciseIndex) => {
  session.completeExercise(exerciseIndex);
};

const isFinishing = ref(false);

const handleCancel = () => {
  if (!isActiveForCurrent.value || isFinishing.value) return;
  if (
    !window.confirm(
      'Отменить тренировку? Все изменения будут потеряны и не будут сохранены.',
    )
  ) {
    return;
  }
  session.clear();
  elapsedSeconds.value = 0;
};

const handleFinish = async () => {
  if (!isActiveForCurrent.value || isFinishing.value) return;
  if (
    session.completionRate < 100 &&
    !window.confirm(
      'Есть невыполненные подходы. Завершить тренировку и сохранить результат?',
    )
  ) {
    return;
  }

  const payload = session.buildResultPayload({
    userLogin: login.value ?? '',
    finishTime: new Date().toISOString(),
  });
  if (!payload) return;

  try {
    isFinishing.value = true;
    await processTrainingStatistics(payload);
    window.alert('Тренировка успешно сохранена');
    session.clear();
    elapsedSeconds.value = 0;
    await trainingsStore.fetchTraining(trainingId.value);
  } catch (error) {
    console.error(error);
    window.alert('Не удалось сохранить тренировку. Попробуйте ещё раз.');
  } finally {
    isFinishing.value = false;
  }
};
</script>

<style scoped>
.session-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  flex-wrap: wrap;
}

.timer-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.25rem;
}

.notice {
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  background: #fff7ed;
  color: #c2410c;
  border: 1px solid #fdba74;
}

.exercise-card {
  gap: 1rem;
}

.exercise-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.exercise-flags {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.badge {
  padding: 0.25rem 0.5rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge--warning {
  background: #fef3c7;
  color: #92400e;
}

.set-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
  align-items: end;
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
}

.input-field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.875rem;
}

.input-field input {
  width: 100%;
}

.checkbox-inline {
  display: inline-flex;
  align-items: center;
  align-self: end;
  gap: 0.5rem;
  white-space: nowrap;
}

.checkbox-inline input[type='checkbox'] {
  width: 1rem;
  height: 1rem;
  margin: 0;
  padding: 0;
  flex: 0 0 auto;
}

.nav-links.wrap {
  flex-wrap: wrap;
  gap: 0.5rem;
}

.removed-exercise {
  padding: 0.75rem;
  border-radius: 0.5rem;
  border: 1px dashed #f97316;
  background: #fffbeb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
}

.finish-button {
  width: fit-content;
  margin-left: auto;
}
</style>
