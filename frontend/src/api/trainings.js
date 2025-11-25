import http from './http';

export const getTrainings = async () => {
  const { data } = await http.get('/trainings');
  return data;
};

export const createTraining = async (payload) => {
  const { data } = await http.post('/trainings', payload);
  return data;
};

export const updateTraining = async (id, payload) => {
  await http.patch(`/trainings/${id}`, payload);
};

export const removeTraining = async (id) => {
  await http.delete(`/trainings/${id}`);
};

export const getTraining = async (id) => {
  const { data } = await http.get(`/trainings/${id}`);
  return data;
};

