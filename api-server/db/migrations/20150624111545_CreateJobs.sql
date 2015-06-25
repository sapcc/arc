
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE jobs
(
  id text PRIMARY KEY,
  version integer NOT NULL,
  sender text NOT NULL,
  "to" text NOT NULL,
  timeout integer NOT NULL,
  agent text NOT NULL,
  action text NOT NULL,
  payload text NOT NULL,
  status integer NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE jobs

