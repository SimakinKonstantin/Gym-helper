SELECT t.id, t.user_login, t.name, t.exercises FROM programs_trainings AS pt
JOIN trainings AS t ON t.id = pt.program_id
WHERE t.user_login = $1;
