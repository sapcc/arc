
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE token_cache
(
  key text PRIMARY KEY,
  value text NOT NULL,
  valid_until timestamp without time zone NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE token_cache;