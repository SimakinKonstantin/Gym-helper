import { createRouter, createWebHashHistory } from 'vue-router';
import AuthView from '@/views/AuthView.vue';
import ProgramsView from '@/views/ProgramsView.vue';
import ProgramDetailsView from '@/views/ProgramDetailsView.vue';
import TrainingsView from '@/views/TrainingsView.vue';
import TrainingSessionView from '@/views/TrainingSessionView.vue';
import TrainingHistoryView from '@/views/TrainingHistoryView.vue';
import TrainingHistoryDetailsView from '@/views/TrainingHistoryDetailsView.vue';
import { useAuthStore } from '@/stores/auth';

const routes = [
  {
    path: '/auth',
    name: 'auth',
    component: AuthView,
    meta: { public: true },
  },
  {
    path: '/',
    redirect: '/app/programs',
  },
  {
    path: '/app/programs',
    name: 'programs',
    component: ProgramsView,
  },
  {
    path: '/app/programs/:id',
    name: 'program-details',
    component: ProgramDetailsView,
    props: true,
  },
  {
    path: '/app/trainings',
    name: 'trainings',
    component: TrainingsView,
  },
  {
    path: '/app/trainings/:id/start',
    name: 'training-session',
    component: TrainingSessionView,
    props: true,
  },
  {
    path: '/app/training-history',
    name: 'training-history',
    component: TrainingHistoryView,
  },
  {
    path: '/app/training-history/:id',
    name: 'training-history-details',
    component: TrainingHistoryDetailsView,
    props: true,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // Если есть сохраненная позиция (например, при нажатии кнопки "назад")
    if (savedPosition) {
      return savedPosition;
    }
    // Иначе прокручиваем наверх
    // Используем небольшую задержку для асинхронных компонентов
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({ top: 0, behavior: 'smooth' });
      }, 100);
    });
  },
});

router.beforeEach((to, _from, next) => {
  if (to.meta.public) {
    next();
    return;
  }
  const auth = useAuthStore();
  if (!auth.isAuthenticated) {
    next({ name: 'auth', query: { redirect: to.fullPath } });
    return;
  }
  next();
});

export default router;

