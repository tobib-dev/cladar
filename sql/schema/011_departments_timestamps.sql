-- +goose Up
ALTER TABLE departments ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();
ALTER TABLE departments ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE departments
DROP COLUMN created_at,
DROP COLUMN updated_at;
