-- +goose Up
ALTER TABLE messages ADD COLUMN image_url TEXT;

-- +goose Down
ALTER TABLE messages DROP COLUMN image_url;