
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE logs
(
  job_id text REFERENCES jobs(requestid) ON DELETE CASCADE,
  content text NOT NULL,
  createdat timestamp NOT NULL,
  updatedat timestamp NOT NULL,
  CONSTRAINT index_logs_on_job_id UNIQUE (job_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE logs;

