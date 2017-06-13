package filter

//go:generate goyacc -v "" -o parser.go expr.y
//go:generate gofmt -w parser.go
//go:generate nex -e -o lexer.go expr.nex
//go:generate gofmt -w lexer.go

import (
	"errors"
	"strings"
)

func Postgresql(query string) (string, error) {

	lexer := NewLexer(strings.NewReader(query))

	parseStatus := yyParse(lexer)

	if parseStatus == 0 {
		return lexer.parseResult.(string), nil
	}

	return "", errors.New(lexer.parseResult.(string))
}
