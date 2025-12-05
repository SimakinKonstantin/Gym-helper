<template>
  <section class="stack" v-if="program">
    <article class="card stack">
      <div>
        <RouterLink to="/app/programs" class="muted">← Назад к программам</RouterLink>
        <h2>{{ program.name }}</h2>
      </div>
      <form class="stack" @submit.prevent="handleAttach">
        <label class="stack-sm">
          День
          <div class="calendar">
            <button
              v-for="day in 30"
              :key="day"
              type="button"
              class="calendar-day"
              :class="{
                'calendar-day--selected': form.day === day,
                'calendar-day--has-training': hasTrainingOnDay(day),
                'calendar-day--rest-day': isRestDay(day),
              }"
              @click="form.day = day"
            >
              {{ day }}
            </button>
          </div>
        </label>
        <label class="stack-sm">
          Тренировка
          <select
            v-model.number="form.trainingId"
            :disabled="form.isRestDay || (form.day && isRestDay(Number(form.day)))"
            :required="!form.isRestDay"
          >
            <option disabled value="">Выберите тренировку</option>
            <option v-for="training in trainings" :key="training.id" :value="training.id">
              {{ training.name }}
            </option>
          </select>
        </label>
        <label class="stack-sm checkbox-inline">
          <span class="muted">Выходной день</span>
          <input 
            v-model="form.isRestDay" 
            type="checkbox"
            :disabled="form.day && (isRestDay(Number(form.day)) || hasTrainingOnDay(Number(form.day)))"
          />
        </label>
        <div v-if="form.day && isRestDay(Number(form.day))" class="stack-sm">
          <span class="muted" style="color: #ef4444;">
            ⚠️ В этот день уже установлен отдых. Удалите день отдыха перед добавлением тренировки.
          </span>
        </div>
        <div v-else-if="form.day && hasTrainingOnDay(Number(form.day)) && form.isRestDay" class="stack-sm">
          <span class="muted" style="color: #ef4444;">
            ⚠️ В этот день уже есть тренировка. Удалите тренировку перед добавлением дня отдыха.
          </span>
        </div>
        <div class="stack-sm">
          <span class="muted">&nbsp;</span>
          <button type="submit" :disabled="!form.day || !canAddToDay">
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
          v-for="daySchedule in programTrainings"
          :key="`day-${daySchedule.day}`"
          class="card stack"
        >
          <div class="stack-sm">
            <h4>День {{ daySchedule.day }}</h4>
            <p v-if="daySchedule.training" class="muted">
              {{ daySchedule.training.name }} · 
              {{ daySchedule.training.exercises?.length ?? 0 }} упражнений
            </p>
            <p v-else class="muted">Выходной день</p>
          </div>
          <div class="nav-links" v-if="daySchedule.training">
            <RouterLink :to="`/app/trainings/${daySchedule.training.id}/start`">
              Запустить
            </RouterLink>
            <RouterLink :to="`/app/trainings?edit=${daySchedule.training.id}`">
              Редактировать
            </RouterLink>
            <button class="ghost" @click="detach(daySchedule)">Удалить</button>
          </div>
          <div class="nav-links" v-else>
            <button class="ghost" @click="detach(daySchedule)">Удалить день отдыха</button>
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
import { computed, onMounted, reactive, watch } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useProgramsStore } from '@/stores/programs';
import { useTrainingsStore } from '@/stores/trainings';

const route = useRoute();
const programsStore = useProgramsStore();
const trainingsStore = useTrainingsStore();
const { current: program, currentTrainings: programTrainingsRaw } =
  storeToRefs(programsStore);
const { items: trainings } = storeToRefs(trainingsStore);

// Сортируем тренировки по дням
const programTrainings = computed(() => {
  return [...programTrainingsRaw.value].sort((a, b) => a.day - b.day);
});

const form = reactive({
  day: '',
  trainingId: '',
  isRestDay: false,
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

// Проверка, можно ли добавить элемент в выбранный день
const canAddToDay = computed(() => {
  if (!form.day) return true;
  const day = Number(form.day);
  const hasRest = isRestDay(day);
  const hasTraining = hasTrainingOnDay(day);
  
  // Если день уже содержит отдых, нельзя добавить ничего
  if (hasRest) return false;
  // Если день уже содержит тренировку и пытаются добавить отдых, нельзя
  if (hasTraining && form.isRestDay) return false;
  // Если день уже содержит тренировку и пытаются добавить еще одну тренировку, можно
  // Но если пытаются добавить отдых, нельзя
  return true;
});

const handleAttach = async () => {
  const id = Number(route.params.id);
  const day = Number(form.day);
  
  // Проверка: если в день уже есть отдых, нельзя добавить ничего
  if (isRestDay(day)) {
    alert('В этот день уже установлен отдых. Удалите день отдыха перед добавлением тренировки.');
    return;
  }
  
  // Проверка: если в день уже есть тренировка и пытаются добавить отдых, нельзя
  if (hasTrainingOnDay(day) && form.isRestDay) {
    alert('В этот день уже есть тренировка. Удалите тренировку перед добавлением дня отдыха.');
    return;
  }
  
  await programsStore.addTraining({
    programId: id,
    day: day,
    trainingId: form.isRestDay ? null : Number(form.trainingId),
  });
  form.day = '';
  form.trainingId = '';
  form.isRestDay = false;
};

const detach = async (daySchedule) => {
  const id = Number(route.params.id);
  const day = daySchedule.day;
  if (!day) return;
  await programsStore.removeTraining({
    programId: id,
    day,
    trainingId: daySchedule.training?.id ?? null,
  });
};

const hasTrainingOnDay = (day) => {
  return programTrainings.value.some(
    (schedule) => schedule.day === day && schedule.training !== null,
  );
};

const isRestDay = (day) => {
  return programTrainings.value.some(
    (schedule) => schedule.day === day && schedule.training === null,
  );
};
</script>

<style scoped>
.calendar {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 0.5rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
}

.calendar-day {
  aspect-ratio: 1;
  border: 2px solid #e5e7eb;
  background: white;
  border-radius: 0.375rem;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  color: #000000;
}

.calendar-day:hover {
  border-color: #3b82f6;
  background: #eff6ff;
  transform: scale(1.05);
}

.calendar-day--selected {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
  font-weight: 600;
}

.calendar-day--has-training {
  border-color: #10b981;
  background: #ecfdf5;
  color: #000000;
}

.calendar-day--has-training.calendar-day--selected {
  background: #10b981;
  color: white;
}

.calendar-day--rest-day {
  border-color: #f59e0b;
  background: #fffbeb;
  color: #000000;
}

.calendar-day--rest-day.calendar-day--selected {
  background: #f59e0b;
  color: white;
}
</style>

