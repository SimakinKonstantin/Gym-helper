-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS programs_trainings (
    training_id INTEGER REFERENCES trainings(id) ON DELETE CASCADE,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    day SMALLINT NOT NULL,
    PRIMARY KEY (training_id, program_id, day)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS programs_trainings;
-- +goose StatementEnd
