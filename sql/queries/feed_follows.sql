-- name: CreateFeedFollow :one
WITH inserted_feed_follow as (INSERT INTO feed_follows (
  user_id,
  feed_id
) VALUES (
  $1,
  $2
) RETURNING *)

SELECT inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
  FROM inserted_feed_follow
  INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
  INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as feed_name 
FROM feed_follows
JOIN feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec 
DELETE FROM feed_follows
WHERE user_id = (
  select id from users
  where users.name = $1
) AND 
feed_id = (
  select id from feeds
  where url = $2
);