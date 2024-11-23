-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds where $1 = name;

-- name: GetFeedByUrl :one
SELECT * FROM feeds where $1 = url;

-- name: DeleteAllFeeds :exec
TRUNCATE feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;