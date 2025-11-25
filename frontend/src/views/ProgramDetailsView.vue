<template>
  <section class="stack" v-if="program">
    <article class="card stack">
      <div>
        <RouterLink to="/programs" class="muted">← Назад к программам</RouterLink>
        <h2>{{ program.name }}</h2>
        <p class="muted">ID: {{ program.id }}</p>
      </div>
      <form class="grid" @submit.prevent="handleAttach">
        <label class="stack-sm">
          День
          <select v-model.number="form.day" required>
            <option disabled value="">Выберите день</option>
            <option v-for="day in 7" :key="day" :value="day">
              День {{ day }}
            </option>
          </select>
        </label>
        <label class="stack-sm">
          Тренировка
          <select v-model.number="form.trainingId" required>
            <option disabled value="">Выберите тренировку</option>
            <option v-for="training in trainings" :key="training.id" :value="training.id">
              {{ training.name }}
            </option>
          </select>
        </label>
        <div class="stack-sm">
          <span class="muted">&nbsp;</span>
          <button type="submit" :disabled="!form.trainingId || !form.day">
            Добавить в программу
          </button>
        </div>
      </form>
    </article>

    <article class="card stack">
      <header>
        <h3>Тренировки программы</h3>
      </header>
      <div class="stack" v-if="programTrainings.length">
        <div
          v-for="training in programTrainings"
          :key="`${training.id}-${training.day ?? training.name}`"
          class="card stack"
        >
          <div class="stack-sm">
            <h4>{{ training.name }}</h4>
            <p class="muted">
              День {{ training.day ?? '—' }} · {{ training.exercises?.length ?? 0 }}
              упражнений
            </p>
          </div>
          <div class="nav-links">
            <RouterLink :to="`/trainings/${training.id}/start`">
              Запустить
            </RouterLink>
            <button class="ghost" @click="detach(training)">Удалить</button>
          </div>
        </div>
      </div>
      <p v-else class="muted">
        В этой программе пока нет тренировок. Добавьте хотя бы одну.
      </p>
    </article>
  </section>
</template>

<script setup>
import { onMounted, reactive, watch } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useProgramsStore } from '@/stores/programs';
import { useTrainingsStore } from '@/stores/trainings';

const route = useRoute();
const programsStore = useProgramsStore();
const trainingsStore = useTrainingsStore();
const { current: program, currentTrainings: programTrainings } =
  storeToRefs(programsStore);
const { items: trainings } = storeToRefs(trainingsStore);

const form = reactive({
  day: '',
  trainingId: '',
});

const load = async () => {
  const id = Number(route.params.id);
  if (Number.isNaN(id)) return;
  await Promise.allSettled([
    programsStore.fetchProgram(id),
    trainingsStore.fetchTrainings(),
  ]);
};

onMounted(load);

watch(
  () => route.params.id,
  () => load(),
);

const handleAttach = async () => {
  const id = Number(route.params.id);
  await programsStore.addTraining({
    programId: id,
    day: Number(form.day),
    trainingId: Number(form.trainingId),
  });
  form.day = '';
  form.trainingId = '';
};

const detach = async (training) => {
  const id = Number(route.params.id);
  const day = training.day ?? Number(prompt('Укажите день тренировки'));
  if (!day) return;
  await programsStore.removeTraining({
    programId: id,
    day,
    trainingId: training.id,
  });
};
</script>

