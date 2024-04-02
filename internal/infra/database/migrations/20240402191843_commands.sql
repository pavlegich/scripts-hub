-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS commands (
    id serial PRIMARY KEY,
    name varchar(32) UNIQUE NOT NULL,
    script text
);

-- create indexes
CREATE INDEX IF NOT EXISTS command_name_idx ON commands (name);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP INDEX command_name_idx;
DROP TABLE users;
