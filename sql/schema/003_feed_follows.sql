-- +goose Up
CREATE TABLE feed_follows (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (feed_id)
  REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;