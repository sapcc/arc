
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION json_delete_keys("json" json, VARIADIC "keys_to_delete" TEXT[])
  RETURNS json
  LANGUAGE sql
  IMMUTABLE
  STRICT
AS $function$
SELECT COALESCE(
  (SELECT ('{' || string_agg(to_json("key") || ':' || "value", ',') || '}')
   FROM json_each("json")
   WHERE "key" <> ALL ("keys_to_delete")),
  '{}'
)::json
$function$
;
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP FUNCTION json_delete_keys("json" json, VARIADIC "keys_to_delete" TEXT[]);
