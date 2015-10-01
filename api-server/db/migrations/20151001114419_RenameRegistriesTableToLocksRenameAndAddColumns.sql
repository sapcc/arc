
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE registries
RENAME TO locks;
ALTER TABLE locks
RENAME COLUMN registry_id to lock_id;
ALTER TABLE locks
ADD COLUMN created_at timestamp without time zone NOT NULL DEFAULT(NOW()),
DROP CONSTRAINT registries_pkey,
ADD CONSTRAINT locks_pkey  PRIMARY KEY (lock_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE locks
RENAME TO registries;
ALTER TABLE registries
RENAME COLUMN lock_id to registry_id;
ALTER TABLE registries
DROP COLUMN created_at,
DROP CONSTRAINT locks_pkey,
ADD CONSTRAINT registries_pkey  PRIMARY KEY (registry_id);