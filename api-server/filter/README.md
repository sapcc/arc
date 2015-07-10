Arc API filter parser
====================
This package contains a parser to transform the filter syntax exposed by the API to a filter expression that can by used by the underlying fact storage system.

As we currently store the facts in a postgresql cloumn of type `jsonb` the parser currently produces a WHERE clause that uses [postgresql specifc json operators](http://www.postgresql.org/docs/9.4/static/functions-json.html).

Example:

`column1 = "1" OR column2 != 2`

is transformed to


`( facts->>'column1' = '1' OR (facts->>'column2')::numeric <> 2 )`


Implementation details
======================
Both the lexer and parser are auto generated. The generated output `parser.go` and `lexer.go` are checked into version control, so unless you want to change the parser you don't need to bother with the following.

The lexer is generated using [nex](https://github.com/blynn/nex) from the file `filter.nex`. It can be installed with `go get https://github.com/blynn/nex`


The parser is generated using go's yacc implmentation `go tool yacc` from the file `filter.y`. It comes with a standard go installation.

The `filter.go` file contains `//go:generate` stanzas to generate both the lexer and the parser for you. To modify the parser you just need to edit `filter.y` or `filter.nex` and run `go generate` in the package directory.
