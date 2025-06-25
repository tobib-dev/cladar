-- +goose Up
UPDATE managers
SET dept_id = (
    SELECT id FROM departments LIMIT 1
)
WHERE dept_id IS NULL;

ALTER TABLE managers ALTER COLUMN dept_id SET NOT NULL;

-- +goose Down
ALTER TABLE managers ALTER COLUMN dept_id DROP NOT NULL;
