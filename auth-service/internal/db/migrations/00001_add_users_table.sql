-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(255) PRIMARY KEY,
    hash VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
