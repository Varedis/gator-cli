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
  last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT
  id,
  name,
  url
FROM
  feeds
ORDER BY last_fetched_at
NULLS FIRST
LIMIT 1;