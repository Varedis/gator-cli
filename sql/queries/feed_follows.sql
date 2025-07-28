-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, user_id, feed_id)
  VALUES ($1, $2, $3)
  RETURNING *
)
SELECT
  inserted_feed_follow.*,
  users.name AS user_name,
  feeds.name AS feed_name
FROM
  inserted_feed_follow
INNER JOIN users
  ON users.id = feed_follows.user_id
INNER JOIN feeds
  ON feeds.id = feed_follows.feed_id;
