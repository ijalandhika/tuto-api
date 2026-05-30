-- name: CreateParent :one 
INSERT INTO parents (
    email,
    password_hash, 
    display_name, 
    marketing_opt_in
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetParentByEmail :one 
SELECT * FROM parents 
WHERE email = $1 
LIMIT 1;

-- name: GetParentByID :one 
SELECT * FROM parents 
WHERE id = $1
LIMIT 1;

-- name: CreateSession :one 
INSERT INTO auth_sessions (
    actor_type, 
    actor_id, 
    token_hash,
    device_id,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT * FROM auth_sessions
WHERE token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > now()
LIMIT 1;

-- name: RevokeSession :exec
UPDATE auth_sessions
  SET revoked_at = now()
WHERE token_hash = $1;

-- name: RevokeAllSessionsByActor :exec
UPDATE auth_sessions
    SET revoked_at = now()
WHERE actor_type = $1
    AND actor_id = $2
    AND revoked_at IS NULL;