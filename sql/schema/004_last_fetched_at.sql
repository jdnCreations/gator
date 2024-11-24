-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;

-- +goose Down
DROP TABLE feed_follows, feeds;