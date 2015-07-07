-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE agents
(
  agent_id text PRIMARY KEY,
	project text,
	organization text,
  facts jsonb,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE agents;