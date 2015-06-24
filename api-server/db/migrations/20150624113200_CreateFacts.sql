
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE facts
(
  agent_id text PRIMARY KEY,
  facts jsonb,
  createdat timestamp NOT NULL,
  updatedat timestamp NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE facts;
