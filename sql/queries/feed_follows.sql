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
  ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds
  ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFollowFeedsForUser :many
SELECT
  feeds.name AS feed_name,
  users.name AS user_name
FROM
  feed_follows
INNER JOIN users
  ON users.id = feed_follows.user_id
INNER JOIN feeds
  ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;
