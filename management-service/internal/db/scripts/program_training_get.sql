SELECT t.id, t.name, t.exercises, pt.day AS day FROM programs_trainings AS pt
LEFT JOIN trainings AS t ON t.id = pt.training_id
WHERE pt.program_id = $1;