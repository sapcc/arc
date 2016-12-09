
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tokens(
  id TEXT PRIMARY KEY,
  profile TEXT,
  subject JSON,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE tokes;
