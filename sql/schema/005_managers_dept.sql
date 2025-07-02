-- +goose Up
ALTER TABLE managers
ADD COLUMN dept_id UUID
REFERENCES departments(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE managers
DROP COLUMN dept_id;