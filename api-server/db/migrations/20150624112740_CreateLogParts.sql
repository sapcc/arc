
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE log_parts
(
  job_id text REFERENCES jobs(requestid) ON DELETE CASCADE,
  number integer NOT NULL,
  content text,
  final boolean NOT NULL DEFAULT FALSE,
  createdat timestamp NOT NULL,
  CONSTRAINT log_parts_job_id_number_index UNIQUE (job_id, number)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE log_parts;

