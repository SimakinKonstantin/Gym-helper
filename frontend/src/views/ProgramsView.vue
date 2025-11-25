<template>
  <section class="stack">
    <article class="card stack">
      <div>
        <h2>Программы</h2>
        <p class="muted">
          Создавайте программы, затем добавляйте в них тренировки по дням.
        </p>
      </div>
      <form class="stack-sm" @submit.prevent="handleCreate">
        <label class="stack-sm">
          Название программы
          <input
            v-model.trim="programName"
            placeholder="Подготовка к полумарафону"
            required
          />
        </label>
        <button type="submit" :disabled="!programName">Создать</button>
      </form>
    </article>

    <article class="card">
      <header class="stack-sm">
        <h3>Список программ</h3>
        <span class="muted" v-if="programs.length === 0">
          Пока пусто — добавьте первую программу.
        </span>
      </header>
      <div class="grid">
        <div v-for="program in programs" :key="program.id" class="card stack">
          <h4>{{ program.name }}</h4>
          <p class="muted">ID: {{ program.id }}</p>
          <div class="nav-links">
            <RouterLink :to="`/programs/${program.id}`">Открыть</RouterLink>
            <button class="ghost" @click="remove(program.id)">Удалить</button>
          </div>
        </div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useProgramsStore } from '@/stores/programs';

const programName = ref('');
const programsStore = useProgramsStore();
const { items: programs } = storeToRefs(programsStore);

onMounted(() => {
  programsStore.fetchPrograms();
});

const handleCreate = async () => {
  if (!programName.value) return;
  await programsStore.createProgram(programName.value);
  programName.value = '';
};

const remove = async (id) => {
  if (!confirm('Удалить программу?')) return;
  await programsStore.deleteProgram(id);
};
</script>

