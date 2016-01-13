-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE jobs
ADD COLUMN user_id text;

UPDATE jobs
SET user_id = 'unknown';

ALTER TABLE jobs 
ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE jobs
DROP COLUMN user_id;