
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX index_agents_on_project on agents (project);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX index_agents_on_project;
