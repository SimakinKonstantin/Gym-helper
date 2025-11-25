import http from './http';

export const getPrograms = async () => {
  const { data } = await http.get('/programs');
  return data;
};

export const createProgram = async (payload) => {
  await http.post('/programs', payload);
};

export const removeProgram = async (id) => {
  await http.delete(`/programs/${id}`);
};

export const getProgram = async (id) => {
  const { data } = await http.get(`/programs/${id}`);
  return data;
};

export const getProgramTrainings = async (programId) => {
  const { data } = await http.get(`/programs/${programId}/trainings`);
  return data;
};

export const attachTrainingToProgram = async (payload) => {
  await http.post('/programs/trainings', payload);
};

export const detachTrainingFromProgram = async (payload) => {
  await http.delete('/programs/trainings', { data: payload });
};

