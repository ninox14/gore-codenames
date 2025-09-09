-- name: CreateUser :one
INSERT INTO users (id, name)
VALUES ($1, $2)
RETURNING id, created_at, name;

-- name: GetUserByID :one
SELECT id, created_at, name
FROM users
WHERE id = $1;

-- name: UpdateUserName :exec
UPDATE users
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, created_at, name
FROM users
ORDER BY created_at DESC;