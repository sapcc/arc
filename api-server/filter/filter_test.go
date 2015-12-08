// +build !integration

package filter

import (
	"fmt"
	"testing"
)

func TestStringEqualFilter(t *testing.T) {
	filter, err := Postgresql(`column = "query"`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}

	expected := fmt.Sprintf(`%s->>'column' = 'query'`, tagsColumn)

	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

	filter, err = Postgresql(`"query" = column`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}
}

func TestStringNotEqualFilter(t *testing.T) {
	filter, err := Postgresql(`column != "query"`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}

	expected := fmt.Sprintf(`%s->>'column' <> 'query'`, tagsColumn)

	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

	filter, err = Postgresql(`"query" != column`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}
}

func TestNumberEqualFilter(t *testing.T) {
	filter, err := Postgresql(`column = 3`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}

	expected := fmt.Sprintf(`(%s->>'column')::numeric = 3`, tagsColumn)

	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

	filter, err = Postgresql(`3 = column`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}
}

func TestNumberNotEqualFilter(t *testing.T) {
	filter, err := Postgresql(`column != 3`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}

	expected := fmt.Sprintf(`(%s->>'column')::numeric <> 3`, tagsColumn)

	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

	filter, err = Postgresql(`3 != column`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}
}

func TestEmptyFilter(t *testing.T) {
	_, err := Postgresql(``)
	if err == nil {
		t.Error("Expected a parse error")
	}
}

func TestCompoundFilter(t *testing.T) {
	filter, err := Postgresql(`@fact1 = "1" OR tag2 != "2" AND NOT (tag3 = "3" AND @fact4 = "4")`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	expected := "( facts->>'fact1' = '1' OR ( tags->>'tag2' <> '2' AND NOT ( ( tags->>'tag3' = '3' AND facts->>'fact4' = '4' ) ) ) )"
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

}
