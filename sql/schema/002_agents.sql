-- +goose Up
CREATE TABLE agents (
	id UUID PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	dept TEXT NOT NULL
);

-- +goose Down
DROP TABLE agents;