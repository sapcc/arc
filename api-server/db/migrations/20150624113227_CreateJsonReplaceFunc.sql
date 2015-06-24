
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION json_replace(old_data json, new_data json)
RETURNS json
IMMUTABLE
LANGUAGE sql
AS $BODY$
  SELECT ('{'||string_agg(to_json(key)||':'||value, ',')||'}')::json
  FROM (
    SELECT * FROM json_each(old_data) WHERE key NOT IN (SELECT json_object_keys(new_data))
    UNION ALL
    SELECT * FROM json_each(new_data) WHERE json_typeof(value) <> 'null'
  ) t;
 $BODY$
;
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP FUNCTION json_replace(old_data json, new_data json);

