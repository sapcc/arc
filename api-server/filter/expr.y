%{

package filter

import (
  "fmt"
  "strings"
)

var factsColumn = "facts"
var tagsColumn = "tags"

var likeSyntax = strings.NewReplacer("*", "%", "+", "_")

%}

%union {
  str string
  num int
}

%type <str> expr term comp

%token <str> FIELD
%token <str> STRING
%token <num> NUMBER

%token '=' '(' ')' '^'
%token NEQ NLIKE AND OR NOT

%left OR
%left AND
%left NOT

%%

expr:
   term
   {
      yylex.(*Lexer).parseResult = $1
   }

term:
  NOT term
  {
    $$ = fmt.Sprintf(`NOT %s`, $2)
  }
|
  '(' term ')'
  {
    $$ = fmt.Sprintf(`( %s )`, $2)
  }
|
  term OR term
  {
    $$ = fmt.Sprintf(`( %s OR %s )`, $1, $3)
  }
|
  term AND term
  {
    $$ = fmt.Sprintf(`( %s AND %s )`, $1, $3)
  }
|
  comp
  {
    $$ = $1
  }

comp:
  FIELD '=' STRING
   {
      $$ = fmt.Sprintf("%s = '%s'", stringKey($1), $3)
   }
|
  STRING '=' FIELD
   {
      $$ = fmt.Sprintf("%s = '%s'", stringKey($3), $1)
   }
|
  FIELD NEQ STRING
   {
      $$ = fmt.Sprintf("%s <> '%s'", stringKey($1), $3)
   }
|
  STRING NEQ FIELD
   {
      $$ = fmt.Sprintf("%s <> '%s'", stringKey($3), $1)
   }
|
  FIELD '^' STRING
   {
      $$ = fmt.Sprintf("%s LIKE '%s'", stringKey($1), likeSyntax.Replace($3))
   }
|
  STRING '^' FIELD
   {
      $$ = fmt.Sprintf("%s LIKE '%s'", stringKey($3), likeSyntax.Replace($1))
   }
|
  FIELD NLIKE STRING
   {
      $$ = fmt.Sprintf("%s NOT LIKE '%s'", stringKey($1), likeSyntax.Replace($3))
   }
|
  STRING NLIKE FIELD
   {
      $$ = fmt.Sprintf("%s NOT LIKE '%s'", stringKey($3), likeSyntax.Replace($1))
   }
|
  FIELD '=' NUMBER
  {
     $$ = fmt.Sprintf("%s = %d", numKey($1), $3)
  }
|
  NUMBER '=' FIELD
  {
     $$ = fmt.Sprintf("%s = %d", numKey($3), $1)
  }
|
  FIELD NEQ NUMBER
  {
    $$ = fmt.Sprintf("%s <> %d", numKey($1), $3)
  }
|
  NUMBER NEQ FIELD
  {
    $$ = fmt.Sprintf("%s <> %d", numKey($3), $1)
  }
%%

func stringKey(field string) string {

  column := tagsColumn
  if field[0] == '@' {
    field=field[1:]
    column = factsColumn
  }
  return fmt.Sprintf("%s->>'%s'", column, field)
}
func numKey(field string) string {
  column := tagsColumn
  if field[0] == '@' {
    field = field[1:]
    column = factsColumn
  }
  return fmt.Sprintf("(%s->>'%s')::numeric", column, field)
}





