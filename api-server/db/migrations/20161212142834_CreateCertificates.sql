
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE certificates(
  fingerprint text PRIMARY KEY,
  common_name text NOT NULL,
  country text,
  locality text,
  organization text,
  organizational_unit text,
  not_before TIMESTAMP WITHOUT TIME ZONE not NULL,
  not_after TIMESTAMP WITHOUT TIME ZONE not NULL,
  pem TEXT NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE certificates;