-- +goose Up
ALTER TABLE managers ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE managers ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE managers
DROP COLUMN created_at,
DROP COLUMN updated_at;
