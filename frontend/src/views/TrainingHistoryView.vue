<template>
  <section class="stack">
    <article class="card stack">
      <header>
        <h2>История тренировок</h2>
        <p class="muted">Просмотр всех выполненных тренировок</p>
      </header>
    </article>

    <article class="card">
      <div class="stack">
        <div
          v-for="training in trainingHistory"
          :key="training.id"
          class="card stack"
          :class="{ 'card--processing': String(training.status || '').toLowerCase() === 'processing' }"
        >
          <div class="stack-sm">
            <div class="nav-links" style="justify-content: space-between; align-items: flex-start;">
              <div>
                <h4 v-if="training.training_name">{{ training.training_name }}</h4>
                <h4 v-else>Тренировка #{{ training.training_id }}</h4>
                <p class="muted">
                  Начало: {{ formatDate(training.start_time) }}<br />
                  Конец: {{ training.finish_time ? formatDate(training.finish_time) : '—' }}
                </p>
                <p v-if="String(training.status || '').toLowerCase() === 'processing'" class="status-processing">
                  ⏳ Производится расчет...
                </p>
              </div>
              <div>
                <RouterLink
                  v-if="String(training.status || '').toLowerCase() === 'done' && training.id"
                  :to="`/app/training-history/${training.id}`"
                >
                  Открыть
                </RouterLink>
                <span v-else-if="String(training.status || '').toLowerCase() === 'processing'" class="muted">Недоступно</span>
                <RouterLink
                  v-else-if="training.id"
                  :to="`/app/training-history/${training.id}`"
                >
                  Открыть
                </RouterLink>
                <span v-else class="muted">ID не найден</span>
              </div>
            </div>
          </div>
        </div>
        <p v-if="!trainingHistory.length && !loading" class="muted">
          Пока нет завершенных тренировок.
        </p>
        <p v-if="loading" class="muted">Загрузка...</p>
      </div>
    </article>
  </section>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { getTrainingHistory } from '@/api/statistics';
import { getTraining } from '@/api/trainings';

const trainingHistory = ref([]);
const loading = ref(false);
let pollingInterval = null;

const formatDate = (dateString) => {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
};

const fetchHistory = async (isPolling = false) => {
  try {
    if (!isPolling) {
      loading.value = true;
    }
    const data = await getTrainingHistory();
    const history = Array.isArray(data) ? data : [data];
    
    // Получаем названия тренировок для тех, у кого их нет
    const historyWithNames = await Promise.all(
      history.map(async (item) => {
        if (!item.training_name && item.training_id) {
          try {
            const training = await getTraining(item.training_id);
            return { ...item, training_name: training.name };
          } catch (err) {
            console.warn('Не удалось получить название тренировки:', err);
            return item;
          }
        }
        return item;
      })
    );
    
    // Сортируем по времени окончания (сначала самые новые)
    trainingHistory.value = historyWithNames.sort((a, b) => {
      const timeA = a.finish_time ? new Date(a.finish_time).getTime() : 0;
      const timeB = b.finish_time ? new Date(b.finish_time).getTime() : 0;
      // Если у тренировки нет finish_time, ставим её в конец
      if (!a.finish_time && !b.finish_time) return 0;
      if (!a.finish_time) return 1;
      if (!b.finish_time) return -1;
      // Сортируем по убыванию (новые сначала)
      return timeB - timeA;
    });
    
    // Логируем для отладки
    if (isPolling) {
      const processingCount = historyWithNames.filter(
        (t) => String(t.status || '').toLowerCase() === 'processing'
      ).length;
      if (processingCount === 0) {
        console.log('Все тренировки завершены, polling остановлен');
      }
    }
  } catch (error) {
    console.error('Ошибка загрузки истории тренировок:', error);
  } finally {
    if (!isPolling) {
      loading.value = false;
    }
  }
};

const startPolling = () => {
  // Останавливаем предыдущий polling, если он был
  stopPolling();
  
  // Проверяем, есть ли тренировки со статусом processing
  const hasProcessing = trainingHistory.value.some((t) => {
    const status = String(t.status || '').toLowerCase();
    return status === 'processing';
  });
  
  if (hasProcessing) {
    // Если есть processing, обновляем каждые 3 секунды
    pollingInterval = setInterval(() => {
      fetchHistory(true); // Передаем флаг, что это polling
    }, 3000);
  }
};

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval);
    pollingInterval = null;
  }
};

onMounted(async () => {
  await fetchHistory();
  startPolling();
});

onUnmounted(() => {
  stopPolling();
});

// Перезапускаем polling при изменении статусов
watch(
  () => trainingHistory.value,
  () => {
    // Небольшая задержка, чтобы убедиться, что данные обновились
    setTimeout(() => {
      stopPolling();
      startPolling();
    }, 100);
  },
  { deep: true },
);
</script>

<style scoped>
.card--processing {
  border-left: 4px solid #f59e0b;
}

.status-processing {
  color: #f59e0b;
  font-weight: 600;
  margin-top: 8px;
}
</style>

