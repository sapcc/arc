-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE agents
ADD COLUMN updated_with text,
ADD COLUMN updated_by text;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE agents
DROP COLUMN updated_with,
DROP COLUMN updated_by;