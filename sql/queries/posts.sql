-- name: CreatePost :one
INSERT INTO posts (id, title, url, description, published_at, feed_id)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
  posts.*,
  feeds.name as feed_name
FROM
  posts
INNER JOIN feed_follows
  ON feed_follows.feed_id = posts.feed_id
INNER JOIN feeds
  ON feeds.id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;

-- name: FindPostsForUser :many
SELECT
  posts.*,
  feeds.name as feed_name
FROM
  posts
INNER JOIN feed_follows
  ON feed_follows.feed_id = posts.feed_id
INNER JOIN feeds
  ON feeds.id = posts.feed_id
WHERE feed_follows.user_id = $1
  AND (
    posts.title LIKE '%' || $2 || '%' OR
    posts.description LIKE '%' || $2 || '%'
  )
ORDER BY posts.published_at DESC
LIMIT $3;

