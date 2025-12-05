<template>
  <div class="app-root">
    <header class="app-header">
      <div class="brand">
        <h1>Система управления тренировками</h1>
        <p>Создавайте и запускайте тренировки</p>
      </div>
      <nav class="nav-links">
        <RouterLink to="/app/programs" v-if="isAuthenticated">Программы</RouterLink>
        <RouterLink to="/app/trainings" v-if="isAuthenticated">Тренировки</RouterLink>
        <RouterLink to="/app/training-history" v-if="isAuthenticated">История тренировок</RouterLink>
        <RouterLink to="/auth" v-if="!isAuthenticated">Вход</RouterLink>
        <button v-else class="ghost" @click="handleLogout">
          Выйти ({{ auth.login }})
        </button>
      </nav>
    </header>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from './stores/auth';

const router = useRouter();
const auth = useAuthStore();
const isAuthenticated = computed(() => auth.isAuthenticated);

const handleLogout = () => {
  auth.logout();
  router.push('/auth');
};
</script>
