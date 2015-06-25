
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE logs
(
  job_id text REFERENCES jobs(id) ON DELETE CASCADE,
  content text NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  PRIMARY KEY(job_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE logs;

