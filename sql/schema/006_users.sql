-- +goose Up
CREATE TYPE user_type AS ENUM ('customer', 'agent', 'manager');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pswd TEXT NOT NULL,
    user_role user_type NOT NULL,
    role_id UUID UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;