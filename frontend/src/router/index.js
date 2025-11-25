import { createRouter, createWebHistory } from 'vue-router';
import AuthView from '@/views/AuthView.vue';
import ProgramsView from '@/views/ProgramsView.vue';
import ProgramDetailsView from '@/views/ProgramDetailsView.vue';
import TrainingsView from '@/views/TrainingsView.vue';
import TrainingSessionView from '@/views/TrainingSessionView.vue';
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
    redirect: '/programs',
  },
  {
    path: '/programs',
    name: 'programs',
    component: ProgramsView,
  },
  {
    path: '/programs/:id',
    name: 'program-details',
    component: ProgramDetailsView,
    props: true,
  },
  {
    path: '/trainings',
    name: 'trainings',
    component: TrainingsView,
  },
  {
    path: '/trainings/:id/start',
    name: 'training-session',
    component: TrainingSessionView,
    props: true,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
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

