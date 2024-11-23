-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users where $1 = name;

-- name: DeleteAll :exec
TRUNCATE feed_follows, feeds, users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one
select * from users where $1 = id;