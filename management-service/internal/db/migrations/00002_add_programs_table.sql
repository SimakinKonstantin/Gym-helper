-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS programs (
    id SERIAL PRIMARY KEY,
    user_login VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS programs;
-- +goose StatementEnd
