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

	expected := fmt.Sprintf(`%s->>'column' = 'query'`, dbCol)

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

	expected := fmt.Sprintf(`%s->>'column' <> 'query'`, dbCol)

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

	expected := fmt.Sprintf(`(%s->>'column')::numeric = 3`, dbCol)

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

	expected := fmt.Sprintf(`(%s->>'column')::numeric <> 3`, dbCol)

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
	filter, err := Postgresql(`column1 = "1" OR column2 != "2" AND NOT (column3 = "3" AND cloumn4 = "4")`)
	if err != nil {
		t.Error("Parsing query failed: ", err)
	}
	expected := "( facts->>'column1' = '1' OR ( facts->>'column2' <> '2' AND NOT ( ( facts->>'column3' = '3' AND facts->>'cloumn4' = '4' ) ) ) )"
	if filter != expected {
		t.Errorf("Unexpected filter result: `%s`", filter)
	}

}
