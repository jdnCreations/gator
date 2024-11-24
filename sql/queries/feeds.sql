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

-- name: GetFeedById :one
select * from feeds where $1 = id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds where $1 = url;

-- name: DeleteAllFeeds :exec
TRUNCATE feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW() 
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * 
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
limit 1;