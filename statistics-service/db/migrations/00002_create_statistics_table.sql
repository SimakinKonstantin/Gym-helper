-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS statistics (
    id SERIAL PRIMARY KEY,
    training_id INT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    finish_time TIMESTAMP NOT NULL,
    result_values JSON NOT NULL,
    status INTEGER NOT NULL REFERENCES statuses(id),
    kcal SMALLINT,
    comment TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS statistics;
-- +goose StatementEnd
