-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trainings (
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    exercises JSON NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trainings;
-- +goose StatementEnd
