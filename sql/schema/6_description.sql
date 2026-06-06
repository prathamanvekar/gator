-- +goose Up
ALTER TABLE posts 
ALTER COLUMN description SET NOT NULL;

-- +goose Down
ALTER TABLE posts 
ALTER COLUMN description DROP NOT NULL;