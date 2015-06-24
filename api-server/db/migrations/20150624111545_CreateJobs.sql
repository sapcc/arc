
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE jobs
(
  version integer NOT NULL,
  sender text NOT NULL,
  requestid text PRIMARY KEY,
  "to" text NOT NULL,
  timeout integer NOT NULL,
  agent text NOT NULL,
  action text NOT NULL,
  payload text NOT NULL,
  status integer NOT NULL,
  createdat timestamp NOT NULL,
  updatedat timestamp NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE jobs

