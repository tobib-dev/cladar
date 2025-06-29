-- +goose Up
CREATE type STATUS AS ENUM ('declined', 'awarded', 'completed');

CREATE TABLE claims(
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL references customers(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL references agents(id) ON DELETE CASCADE,
    claim_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    current_status STATUS NOT NULL,
    award REAL
);

-- +goose Down
DROP TYPE STATUS;
DROP TABLE claims;
