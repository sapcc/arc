
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE agents ALTER COLUMN tags SET DEFAULT '{}'::jsonb;
UPDATE agents SET tags='{}'::jsonb WHERE tags IS NULL;
ALTER TABLE agents ALTER COLUMN tags SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration s rolled back
ALTER TABLE agents ALTER COLUMN tags DROP NOT NULL;
ALTER TABLE agents ALTER COLUMN tags DROP DEFAULT;
