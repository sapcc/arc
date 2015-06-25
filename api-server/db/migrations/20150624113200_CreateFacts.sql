
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE facts
(
  agent_id text PRIMARY KEY,
  facts jsonb,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE facts;
