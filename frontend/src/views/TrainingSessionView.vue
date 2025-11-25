<template>
  <section class="stack" v-if="training">
    <article class="card stack">
      <RouterLink to="/trainings" class="muted">← Все тренировки</RouterLink>
      <h2>{{ training.name }}</h2>
      <p class="muted">
        {{ training.exercises?.length ?? 0 }} упражнений ·
        {{ session.completionRate }}% выполнено
      </p>
      <div class="nav-links">
        <button v-if="!isActiveForCurrent" @click="start">Запустить тренировку</button>
        <button v-else class="ghost" @click="session.clear">Завершить</button>
      </div>
    </article>

    <article class="card stack">
      <div
        v-for="(exercise, exerciseIndex) in training.exercises"
        :key="exerciseIndex"
        class="card stack-sm"
      >
        <strong>{{ exercise.name }}</strong>
        <div class="grid">
          <label
            v-for="(set, setIndex) in exercise.sets"
            :key="setIndex"
            class="card stack-sm"
          >
            <span>
              Подход {{ setIndex + 1 }} · {{ set.reps }} повторов ·
              {{ set.weight }} кг
            </span>
            <input
              type="checkbox"
              :checked="isCompleted(exerciseIndex, setIndex)"
              :disabled="!isActiveForCurrent"
              @change="toggle(exerciseIndex, setIndex)"
            />
          </label>
        </div>
      </div>
    </article>
  </section>
  <p v-else class="muted">Загружаем тренировку…</p>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useTrainingsStore } from '@/stores/trainings';
import { useSessionStore } from '@/stores/session';

const route = useRoute();
const trainingsStore = useTrainingsStore();
const session = useSessionStore();
session.hydrate();

const { current: training } = storeToRefs(trainingsStore);
const trainingId = computed(() => Number(route.params.id));

const load = async () => {
  await trainingsStore.fetchTraining(trainingId.value);
};

onMounted(load);

const isActiveForCurrent = computed(
  () => session.activeTraining?.id === training.value?.id && session.isActive,
);

const start = () => {
  if (!training.value) return;
  session.startSession(training.value);
};

const isCompleted = (exerciseIndex, setIndex) =>
  Boolean(session.completedSets[`${exerciseIndex}-${setIndex}`]);

const toggle = (exerciseIndex, setIndex) => {
  session.toggleSetCompletion(exerciseIndex, setIndex);
};
</script>

