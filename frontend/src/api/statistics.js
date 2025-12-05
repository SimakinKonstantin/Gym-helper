import http from './http';

export const processTrainingStatistics = async (payload) => {
  await http.post('/statistics/process-training', payload);
};

export const getTrainingHistory = async () => {
  const response = await http.get('/statistics');
  return response.data;
};

export const getTrainingDetails = async (id) => {
  const response = await http.get(`/statistics/${id}`);
  return response.data;
};





