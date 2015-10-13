-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE jobs
ADD COLUMN project text;

UPDATE jobs
SET project = agents.project
FROM agents
WHERE agents.agent_id = jobs.to;

ALTER TABLE jobs 
ALTER COLUMN project SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE jobs
DROP COLUMN project;

