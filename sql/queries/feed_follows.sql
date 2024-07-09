-- name: CreateFeedFollow :one
INSERT INTO feed_follows(id, created_at, updated_at, user_id, feeds_id) VALUES($1, $2, $3, $4, $5) RETURNING *;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows where id = $1 RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows;