package filter

//go:generate go tool yacc -v "" -o parser.go expr.y
//go:generate nex -e -o lexer.go expr.nex

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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

func main() {

	input, _ := ioutil.ReadAll(os.Stdin)

	lex := NewLexer(bytes.NewReader(input))

	//var val yySymType
	//for t := lex.Lex(&val); t != 0; t = lex.Lex(&val) {
	//  fmt.Printf("type: %d, val: %s\n", t, val.str)
	//}

	status := yyParse(lex)
	if status == 0 {
		fmt.Println("parse result:", lex.parseResult)
	} else {
		fmt.Println("Error:", lex.parseResult)

		fmt.Print(string(input))
		fmt.Printf(strings.Repeat(" ", lex.Column()) + "^\n")

	}

}
