-- name: CreateClaim :one
INSERT INTO claims (id, customer_id, agent_id, claim_type, created_at, updated_at, current_status, award)
VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW(), 'pending', $4)
RETURNING *;
