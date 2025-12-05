INSERT INTO statistics (user_name, training_id, start_time, finish_time, result_values, status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;