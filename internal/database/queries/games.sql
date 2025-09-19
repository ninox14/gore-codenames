-- name: CreateGame :one
INSERT INTO games (
    id,
    host_id,
    word_pack_id,
    game_state
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetGameByID :one
SELECT * FROM games
WHERE id = $1;

-- name: GetGamesByHost :many
SELECT * FROM games
WHERE host_id = $1
ORDER BY created_at DESC;

-- name: GetGamesByStatus :many
SELECT * FROM games
WHERE status = $1
ORDER BY created_at DESC;

-- name: UpdateGameStatus :one
UPDATE games
SET status = $2,
    started_at = CASE
        WHEN $2 = 'Started' AND started_at IS NULL THEN CURRENT_TIMESTAMP
        ELSE started_at
    END
WHERE id = $1
RETURNING *;

-- name: UpdateGameState :one
UPDATE games
SET game_state = $2
WHERE id = $1
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1;

-- name: GetRecentGames :many
SELECT * FROM games
ORDER BY created_at DESC
LIMIT $1;

-- name: GetGamesByWordPack :many
SELECT * FROM games
WHERE word_pack_id = $1
ORDER BY created_at DESC;

-- name: CountGamesByHost :one
SELECT COUNT(*) FROM games
WHERE host_id = $1;

-- name: GetGamesByHostAndStatus :many
SELECT * FROM games
WHERE host_id = $1
AND status = $2
ORDER BY created_at DESC;