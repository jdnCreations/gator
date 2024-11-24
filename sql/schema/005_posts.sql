-- +goose Up
CREATE TABLE posts (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  title TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID not null,
  FOREIGN KEY (feed_id)
  REFERENCES feeds(id)
);

-- +goose Down
DROP table posts;