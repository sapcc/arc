
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tags
(
  agent_id text REFERENCES agents(agent_id) ON DELETE CASCADE,
	project text NOT NULL,
	"value" text NOT NULL,
  created_at timestamp without time zone NOT NULL,
	PRIMARY KEY(agent_id, value)
);
CREATE INDEX index_tags_on_agent_id on tags (agent_id);
CREATE INDEX index_tags_on_value on tags ("value");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX index_tags_on_agent_id;
DROP INDEX index_tags_on_value;
DROP TABLE tags;
