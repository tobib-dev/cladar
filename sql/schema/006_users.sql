-- +goose Up
CREATE TYPE user_type_new AS ENUM ('agent', 'manager');

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_type;

ALTER TYPE user_type_new RENAME TO user_type;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pswd TEXT NOT NULL,
    user_role user_type NOT NULL,
    role_id UUID UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TYPE user_type;
