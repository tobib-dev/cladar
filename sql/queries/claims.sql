-- name: CreateClaim :one
INSERT INTO claims (id, customer_id, agent_id, claim_type, created_at, updated_at, current_status, award)
VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW(), 'pending', $4)
RETURNING *;

-- name: GetAllClaims :many
SELECT * FROM claims;

-- name: GetClaimByID :one
SELECT * FROM claims
WHERE id = $1;

-- name: GetAllClaimsByCust :many
SELECT * FROM claims
WHERE customer_id = $1;

-- name: GetAllClaimsByAgent :many
SELECT * FROM claims
WHERE agent_id = $1;

-- name: ChangeAssignedAgent :one
UPDATE claims
SET agent_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ChangeClaimType :one
UPDATE claims
SET claim_type = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeclineClaim :one
UPDATE claims
SET current_status = 'declined',
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ApproveClaim :one
UPDATE claims
SET current_status = 'awarded',
    updated_at = NOW(),
    award = $2
WHERE id = $1
RETURNING *;

-- name: ChangeAwardAmount :one
UPDATE claims
SET award = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
