
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE registries
(
  registry_id text PRIMARY KEY,
	agent_id text
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE registries;