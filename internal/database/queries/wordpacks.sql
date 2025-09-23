-- name: CreateWordpack :one
INSERT INTO wordpacks (
    name, description, created_by, is_default, words
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetWordpack :one
SELECT * FROM wordpacks
WHERE id = $1 LIMIT 1;

-- name: GetWordpackByName :one
SELECT * FROM wordpacks
WHERE name = $1 LIMIT 1;

-- name: ListWordpacks :many
SELECT * FROM wordpacks
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: ListWordpacksByUser :many
SELECT * FROM wordpacks
WHERE created_by = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: ListDefaultWordpacks :many
SELECT * FROM wordpacks
WHERE is_default = true
ORDER BY name;

-- name: UpdateWordpack :one
UPDATE wordpacks
SET
    name = $2,
    description = $3,
    is_default = $4,
    words = $5
WHERE id = $1
RETURNING *;

-- name: DeleteWordpack :exec
DELETE FROM wordpacks
WHERE id = $1;

-- name: SearchWordpacks :many
SELECT * FROM wordpacks
WHERE name ILIKE $1 OR description ILIKE $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: CountWordpacks :one
SELECT COUNT(*) FROM wordpacks;

-- name: CountWordpacksByUser :one
SELECT COUNT(*) FROM wordpacks
WHERE created_by = $1;

-- name: GetWordpacksByIDs :many
SELECT * FROM wordpacks
WHERE id = ANY($1::int[]);

-- name: CheckWordpackExists :one
SELECT EXISTS(SELECT 1 FROM wordpacks WHERE id = $1);

-- name: CheckWordpackNameExists :one
SELECT EXISTS(SELECT 1 FROM wordpacks WHERE name = $1 AND id != $2);

-- name: UpdateWordpackWords :one
UPDATE wordpacks
SET words = $2
WHERE id = $1
RETURNING *;

-- name: ToggleWordpackDefault :one
UPDATE wordpacks
SET is_default = $2
WHERE id = $1
RETURNING *;