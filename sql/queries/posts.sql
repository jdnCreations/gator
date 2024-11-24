-- name: CreatePost :one
INSERT INTO posts (title, url, description, published_at, feed_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
join feeds on posts.feed_id = feeds.id
join feed_follows on feeds.id = feed_follows.feed_id
where feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
