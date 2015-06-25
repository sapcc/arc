
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE log_parts
(
  job_id text REFERENCES jobs(id) ON DELETE CASCADE,
  "number" integer NOT NULL,
  content text,
  final boolean NOT NULL DEFAULT FALSE,
  created_at timestamp without time zone NOT NULL,
  PRIMARY KEY(job_id, "number")
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE log_parts;

