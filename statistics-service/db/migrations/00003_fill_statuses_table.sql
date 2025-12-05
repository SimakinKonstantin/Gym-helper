-- +goose Up
-- +goose StatementBegin
INSERT INTO statuses (name) VALUES
                         ('processing'),
                         ('done');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM statuses WHERE name = 'processing' or name = 'done';
-- +goose StatementEnd
