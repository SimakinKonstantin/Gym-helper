<template>
  <section class="grid">
    <article class="card stack">
      <header>
        <h2>{{ editingId ? 'Редактировать тренировку' : 'Новая тренировка' }}</h2>
      </header>
      <form class="stack" @submit.prevent="handleSubmit">
        <label class="stack-sm">
          Название
          <input v-model.trim="form.name" placeholder="Ноги + плечи" required />
        </label>
        <div class="stack">
          <div
            v-for="(exercise, exerciseIndex) in form.exercises"
            :key="exerciseIndex"
            class="card stack-sm"
          >
            <div class="nav-links">
              <strong>Упражнение {{ exerciseIndex + 1 }}</strong>
              <button
                v-if="form.exercises.length > 1"
                class="ghost"
                type="button"
                @click="removeExercise(exerciseIndex)"
              >
                Удалить
              </button>
            </div>
            <input
              v-model="exercise.name"
              placeholder="Присед"
              required
            />
            <div class="stack-sm">
              <strong>Подходы</strong>
              <div
                v-for="(set, setIndex) in exercise.sets"
                :key="setIndex"
                class="nav-links set-row"
              >
                <span class="set-row__title">Подход {{ setIndex + 1 }}</span>
                <label class="set-field">
                  <span class="set-field__label">Повторения</span>
                  <input
                    type="number"
                    min="1"
                    v-model.number="set.reps"
                    placeholder="Сколько повторов"
                  />
                </label>
                <label class="set-field">
                  <span class="set-field__label">Вес, кг</span>
                  <input
                    type="number"
                    min="0"
                    step="2.5"
                    v-model.number="set.weight"
                    placeholder="Вес снаряда"
                  />
                </label>
                <label class="set-field">
                  <span class="set-field__label">Сженных калорий за повтор</span>
                  <input
                    type="number"
                    min="0"
                    v-model.number="set.calPerSet"
                    placeholder="Энергозатраты"
                  />
                </label>
                <button
                  v-if="exercise.sets.length > 1"
                  class="ghost"
                  type="button"
                  @click="removeSet(exerciseIndex, setIndex)"
                >
                  ×
                </button>
              </div>
              <button type="button" @click="addSet(exerciseIndex)">
                + Подход
              </button>
            </div>
          </div>
          <button type="button" @click="addExercise">+ Упражнение</button>
        </div>

        <div class="nav-links">
          <button type="submit">
            {{ editingId ? 'Сохранить изменения' : 'Создать тренировку' }}
          </button>
          <button v-if="editingId" class="ghost" type="button" @click="resetForm">
            Отмена
          </button>
        </div>
      </form>
    </article>

    <article class="card stack">
      <header class="stack-sm">
        <h3>Все тренировки</h3>
        <p class="muted">Выберите тренировку, чтобы запустить или отредактировать.</p>
      </header>
      <div class="stack">
        <div v-for="training in trainings" :key="training.id" class="card stack">
          <div class="stack-sm">
            <h4>{{ training.name }}</h4>
            <p class="muted">
              {{ training.exercises?.length ?? 0 }} упражнений
            </p>
          </div>
          <div class="nav-links">
            <RouterLink :to="`/app/trainings/${training.id}/start`">Запуск</RouterLink>
            <button class="ghost" type="button" @click="edit(training)">Редактировать</button>
            <button class="ghost" type="button" @click="remove(training.id)">Удалить</button>
          </div>
        </div>
        <p v-if="!trainings.length" class="muted">Пока нет тренировок.</p>
      </div>
    </article>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useTrainingsStore } from '@/stores/trainings';

const route = useRoute();
const trainingsStore = useTrainingsStore();
const { items: trainings } = storeToRefs(trainingsStore);
const editingId = ref(null);

const emptyExercise = () => ({
  name: '',
  sets: [{ reps: 10, weight: 20, calPerSet: 0 }],
});

const form = reactive({
  name: '',
  exercises: [emptyExercise()],
});

const resetForm = () => {
  form.name = '';
  form.exercises = [emptyExercise()];
  editingId.value = null;
};

const openEditForm = (trainingId) => {
  const training = trainings.value.find((t) => t.id === Number(trainingId));
  if (training) {
    edit(training);
  }
};

onMounted(async () => {
  await trainingsStore.fetchTrainings();
  // Проверяем query параметр для автоматического открытия формы редактирования
  const editId = route.query.edit;
  if (editId) {
    openEditForm(editId);
  }
});

// Отслеживаем изменения query параметра
watch(
  () => route.query.edit,
  (editId) => {
    if (editId) {
      openEditForm(editId);
    } else {
      resetForm();
    }
  },
);

const handleSubmit = async () => {
  const payload = {
    name: form.name,
    exercises: form.exercises.map((exercise) => ({
      name: exercise.name,
      sets: exercise.sets.map((set) => ({
        reps: Number(set.reps),
        weight: Number(set.weight),
        cal_per_set: Number(set.calPerSet ?? set.cal_per_set ?? 0),
      })),
    })),
  };

  if (editingId.value) {
    await trainingsStore.update(editingId.value, payload);
  } else {
    await trainingsStore.create(payload);
  }
  resetForm();
};

const addExercise = () => {
  form.exercises.push(emptyExercise());
};

const removeExercise = (index) => {
  form.exercises.splice(index, 1);
};

const addSet = (exerciseIndex) => {
  form.exercises[exerciseIndex].sets.push({ reps: 10, weight: 20, calPerSet: 0 });
};

const removeSet = (exerciseIndex, setIndex) => {
  form.exercises[exerciseIndex].sets.splice(setIndex, 1);
};

const edit = (training) => {
  editingId.value = training.id;
  form.name = training.name;
  form.exercises = training.exercises.map((exercise) => ({
    name: exercise.name,
    sets: exercise.sets.map((set) => ({
      reps: set.reps,
      weight: set.weight,
      calPerSet: set.cal_per_set ?? set.calPerSet ?? 0,
    })),
  }));
};

const remove = async (id) => {
  if (!confirm('Удалить тренировку?')) return;
  if (editingId.value === id) resetForm();
  await trainingsStore.remove(id);
};
</script>

<style scoped>
.set-row {
  gap: 1rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

.set-row__title {
  font-weight: 600;
  min-width: 80px;
}

.set-field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  min-width: 120px;
}

.set-field__label {
  font-size: 0.75rem;
  text-transform: uppercase;
  color: #6b7280;
  letter-spacing: 0.02em;
}
</style>

