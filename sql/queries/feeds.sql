-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListFeeds :many
SELECT
  feeds.name,
  feeds.url,
  users.name AS user
FROM
  feeds
INNER JOIN users
  ON feeds.user_id = users.id; 

-- name: GetFeedByURL :one
SELECT id
FROM
  feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
  updated_at = NOW(),
  last_updated_at = NOW()
WHERE id = $1;