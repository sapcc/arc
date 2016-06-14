
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE jobs ADD COLUMN "user" jsonb DEFAULT '{}'::jsonb;
UPDATE jobs SET "user"=('{"id":"' || "user_id" || '"}')::jsonb;
ALTER TABLE jobs DROP COLUMN "user_id";

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE jobs ADD COLUMN "user_id" text;
UPDATE jobs SET "user_id"=("user"->>'id');
ALTER TABLE jobs DROP COLUMN "user";