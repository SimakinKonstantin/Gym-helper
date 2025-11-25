<template>
  <section class="card auth-card stack">
    <div class="stack-sm">
      <h2>{{ isSignUp ? 'Регистрация' : 'Вход' }}</h2>
      <p class="muted">
        Используйте логин и пароль для {{ isSignUp ? 'создания нового' : 'входа в' }}
        аккаунта. Успешная авторизация сохранит JWT в LocalStorage и будет автоматически
        подставлять заголовки <code>X-User-Login</code> и <code>Authorization</code>.
      </p>
    </div>

    <div class="nav-links">
      <button
        class="ghost"
        :class="{ active: !isSignUp }"
        type="button"
        @click="mode = 'sign-in'"
      >
        Войти
      </button>
      <button
        class="ghost"
        :class="{ active: isSignUp }"
        type="button"
        @click="mode = 'sign-up'"
      >
        Регистрация
      </button>
    </div>

    <form class="stack" @submit.prevent="handleSubmit">
      <label class="stack-sm">
        Логин
        <input
          v-model.trim="form.login"
          placeholder="sportlover"
          required
          minlength="3"
        />
      </label>

      <label class="stack-sm">
        Пароль
        <input
          v-model="form.password"
          type="password"
          placeholder="Минимум 6 символов"
          required
          minlength="6"
        />
      </label>

      <p v-if="auth.error" class="muted" style="color: #dc2626">
        {{ auth.error }}
      </p>

      <button type="submit" :disabled="auth.loading">
        {{ isSignUp ? 'Создать аккаунт' : 'Войти' }}
      </button>
    </form>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();

const mode = ref('sign-in');
const form = reactive({
  login: '',
  password: '',
});

const isSignUp = computed(() => mode.value === 'sign-up');

const handleSubmit = async () => {
  if (!form.login || !form.password) return;
  try {
    if (isSignUp.value) {
      await auth.signUp(form);
    } else {
      await auth.signIn(form);
    }
    const redirect = route.query.redirect ?? '/programs';
    router.push(redirect);
  } catch {
    // ошибки отображаются из стора
  }
};
</script>
