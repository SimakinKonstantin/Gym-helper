<template>
  <section class="stack" v-if="trainingDetails">
    <article class="card stack">
      <RouterLink to="/app/training-history" class="muted">← История тренировок</RouterLink>
      <div>
        <h2>{{ trainingDetails.training_name || `Тренировка #${trainingDetails.training_id}` }}</h2>
        <div class="training-info">
          <div class="info-item">
            <span class="muted">Время начала:</span>
            <strong>{{ formatDateTime(trainingDetails.start_time) }}</strong>
          </div>
          <div class="info-item">
            <span class="muted">Время конца:</span>
            <strong>{{ formatDateTime(trainingDetails.finish_time) }}</strong>
          </div>
          <div class="info-item">
            <span class="muted">Продолжительность:</span>
            <strong>{{ formatDuration(trainingDetails.start_time, trainingDetails.finish_time) }}</strong>
          </div>
          <div class="info-item" v-if="trainingDetails.kcal">
            <span class="muted">Сожжено ккал:</span>
            <strong>{{ trainingDetails.kcal }}</strong>
          </div>
        </div>
      </div>
    </article>

    <article class="card stack" v-if="trainingDetails.result_values && trainingDetails.result_values.length > 0">
      <h3>Результаты упражнений</h3>
      <div
        v-for="(exercise, exerciseIndex) in trainingDetails.result_values"
        :key="exerciseIndex"
        class="exercise-result-card"
      >
        <h4>{{ exercise.name }}</h4>
        <div class="sets-table-wrapper">
          <table class="sets-table">
            <thead>
              <tr>
                <th>Подход</th>
                <th>Вес, кг</th>
                <th>Планировалось повторов</th>
                <th>Выполнено повторов</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(set, setIndex) in exercise.sets"
                :key="setIndex"
              >
                <td>{{ setIndex + 1 }}</td>
                <td>{{ set.weight || '—' }}</td>
                <td>{{ set.original_reps || '—' }}</td>
                <td
                  :class="getRepsClass(set.original_reps, set.real_reps)"
                >
                  {{ set.real_reps || '—' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </article>

    <article class="card stack" v-if="trainingDetails.comment">
      <h3>Комментарий к тренировке</h3>
      <p class="comment-text">{{ trainingDetails.comment }}</p>
    </article>

    <article class="card" v-if="String(trainingDetails.status || '').toLowerCase() === 'processing'">
      <p class="status-processing">⏳ Производится расчет результатов...</p>
    </article>
  </section>
  <section v-else class="stack">
    <article class="card">
      <p class="muted">Загрузка...</p>
    </article>
  </section>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { getTrainingDetails } from '@/api/statistics';
import { getTraining } from '@/api/trainings';

const props = defineProps({
  id: {
    type: [String, Number],
    required: true,
  },
});

const route = useRoute();
const trainingDetails = ref(null);
let pollingInterval = null;

const formatDateTime = (dateString) => {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date);
};

const formatDuration = (startTime, finishTime) => {
  if (!startTime || !finishTime) return '—';
  const start = new Date(startTime);
  const finish = new Date(finishTime);
  const diffMs = finish - start;
  const hours = Math.floor(diffMs / (1000 * 60 * 60));
  const minutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((diffMs % (1000 * 60)) / 1000);
  
  if (hours > 0) {
    return `${hours}ч ${minutes}м ${seconds}с`;
  } else if (minutes > 0) {
    return `${minutes}м ${seconds}с`;
  } else {
    return `${seconds}с`;
  }
};

const getRepsClass = (originalReps, realReps) => {
  if (originalReps === null || originalReps === undefined || realReps === null || realReps === undefined) {
    return '';
  }
  if (realReps > originalReps) {
    return 'reps-more';
  } else if (realReps < originalReps) {
    return 'reps-less';
  }
  return '';
};

const fetchDetails = async () => {
  try {
    // Используем props.id, если доступен, иначе route.params.id
    const id = props.id || route.params.id;
    if (!id) {
      console.error('ID тренировки не найден:', { 
        propsId: props.id, 
        routeParamsId: route.params.id,
        allParams: route.params,
        fullRoute: route
      });
      return;
    }
    console.log('Загружаем детали тренировки с ID:', id);
    const data = await getTrainingDetails(String(id));
    
    // Если нет названия тренировки, получаем его отдельно
    if (!data.training_name && data.training_id) {
      try {
        const training = await getTraining(data.training_id);
        data.training_name = training.name;
      } catch (err) {
        console.warn('Не удалось получить название тренировки:', err);
      }
    }
    
    trainingDetails.value = data;
    
    // Если статус processing, продолжаем polling
    const status = String(data.status || '').toLowerCase();
    if (status === 'processing') {
      if (!pollingInterval) {
        pollingInterval = setInterval(() => {
          fetchDetails();
        }, 3000);
      }
    } else {
      // Если статус изменился на done, останавливаем polling
      stopPolling();
    }
  } catch (error) {
    console.error('Ошибка загрузки деталей тренировки:', error);
  }
};

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval);
    pollingInterval = null;
  }
};

onMounted(() => {
  fetchDetails();
});

// Отслеживаем изменения параметра id
watch(
  () => props.id || route.params.id,
  (newId) => {
    if (newId) {
      stopPolling();
      fetchDetails();
    }
  },
);

onUnmounted(() => {
  stopPolling();
});
</script>

<style scoped>
.training-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.exercise-result-card {
  margin-bottom: 24px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 12px;
}

.exercise-result-card:last-child {
  margin-bottom: 0;
}

.exercise-result-card h4 {
  margin: 0 0 12px 0;
  font-size: 1.1rem;
}

.sets-table-wrapper {
  overflow-x: auto;
}

.sets-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.sets-table th,
.sets-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.sets-table th {
  background: #f1f5f9;
  font-weight: 600;
  color: #475569;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.sets-table td {
  color: #0f172a;
}

.sets-table tbody tr:last-child td {
  border-bottom: none;
}

.reps-more {
  background-color: #10b981;
  color: white;
  font-weight: 600;
  text-align: center;
}

.reps-less {
  background-color: #f59e0b;
  color: white;
  font-weight: 600;
  text-align: center;
}

.comment-text {
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.status-processing {
  color: #f59e0b;
  font-weight: 600;
  text-align: center;
  padding: 16px;
}
</style>

