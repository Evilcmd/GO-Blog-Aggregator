-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id) VALUES($1, $2, $3, $4) RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM FEEDS;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: UpdateLastFetched :one
UPDATE feeds
SET last_fetched_at=NOW(), updated_at=NOW()
WHERE id = $1
RETURNING *;