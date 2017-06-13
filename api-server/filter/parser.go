//line expr.y:2
package filter

import __yyfmt__ "fmt"

//line expr.y:3
import (
	"fmt"
	"strings"
)

var factsColumn = "facts"
var tagsColumn = "tags"

var likeSyntax = strings.NewReplacer("*", "%", "+", "_")

//line expr.y:17
type yySymType struct {
	yys int
	str string
	num int
}

const FIELD = 57346
const STRING = 57347
const NUMBER = 57348
const NEQ = 57349
const NLIKE = 57350
const AND = 57351
const OR = 57352
const NOT = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"FIELD",
	"STRING",
	"NUMBER",
	"'='",
	"'('",
	"')'",
	"'^'",
	"NEQ",
	"NLIKE",
	"AND",
	"OR",
	"NOT",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line expr.y:129
func stringKey(field string) string {

	column := tagsColumn
	if field[0] == '@' {
		field = field[1:]
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

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 40

var yyAct = [...]int{

	6, 7, 8, 25, 4, 10, 9, 10, 9, 17,
	10, 3, 19, 18, 20, 13, 31, 2, 15, 14,
	16, 11, 12, 5, 30, 21, 37, 23, 24, 22,
	28, 29, 26, 27, 36, 35, 34, 33, 32, 1,
}
var yyPact = [...]int{

	-4, -1000, -8, -4, -4, -1000, 8, 2, 18, -4,
	-4, -1000, -6, 27, 25, 19, 11, 34, 33, 32,
	31, 30, 22, -3, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 39, 17, 23,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3,
}
var yyR2 = [...]int{

	0, 1, 2, 3, 3, 3, 1, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 15, 8, -3, 4, 5, 6, 14,
	13, -2, -2, 7, 11, 10, 12, 7, 11, 10,
	12, 7, 11, -2, -2, 9, 5, 6, 5, 6,
	5, 5, 4, 4, 4, 4, 4, 4,
}
var yyDef = [...]int{

	0, -2, 1, 0, 0, 6, 0, 0, 0, 0,
	0, 2, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 4, 5, 3, 7, 15, 9, 17,
	11, 13, 8, 10, 12, 14, 16, 18,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	8, 9, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 7, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 10,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 11, 12, 13, 14, 15,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line expr.y:39
		{
			yylex.(*Lexer).parseResult = yyDollar[1].str
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line expr.y:45
		{
			yyVAL.str = fmt.Sprintf(`NOT %s`, yyDollar[2].str)
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:50
		{
			yyVAL.str = fmt.Sprintf(`( %s )`, yyDollar[2].str)
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:55
		{
			yyVAL.str = fmt.Sprintf(`( %s OR %s )`, yyDollar[1].str, yyDollar[3].str)
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:60
		{
			yyVAL.str = fmt.Sprintf(`( %s AND %s )`, yyDollar[1].str, yyDollar[3].str)
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line expr.y:65
		{
			yyVAL.str = yyDollar[1].str
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:71
		{
			yyVAL.str = fmt.Sprintf("%s = '%s'", stringKey(yyDollar[1].str), yyDollar[3].str)
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:76
		{
			yyVAL.str = fmt.Sprintf("%s = '%s'", stringKey(yyDollar[3].str), yyDollar[1].str)
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:81
		{
			yyVAL.str = fmt.Sprintf("%s <> '%s'", stringKey(yyDollar[1].str), yyDollar[3].str)
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:86
		{
			yyVAL.str = fmt.Sprintf("%s <> '%s'", stringKey(yyDollar[3].str), yyDollar[1].str)
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:91
		{
			yyVAL.str = fmt.Sprintf("%s LIKE '%s'", stringKey(yyDollar[1].str), likeSyntax.Replace(yyDollar[3].str))
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:96
		{
			yyVAL.str = fmt.Sprintf("%s LIKE '%s'", stringKey(yyDollar[3].str), likeSyntax.Replace(yyDollar[1].str))
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:101
		{
			yyVAL.str = fmt.Sprintf("%s NOT LIKE '%s'", stringKey(yyDollar[1].str), likeSyntax.Replace(yyDollar[3].str))
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:106
		{
			yyVAL.str = fmt.Sprintf("%s NOT LIKE '%s'", stringKey(yyDollar[3].str), likeSyntax.Replace(yyDollar[1].str))
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:111
		{
			yyVAL.str = fmt.Sprintf("%s = %d", numKey(yyDollar[1].str), yyDollar[3].num)
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:116
		{
			yyVAL.str = fmt.Sprintf("%s = %d", numKey(yyDollar[3].str), yyDollar[1].num)
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:121
		{
			yyVAL.str = fmt.Sprintf("%s <> %d", numKey(yyDollar[1].str), yyDollar[3].num)
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line expr.y:126
		{
			yyVAL.str = fmt.Sprintf("%s <> %d", numKey(yyDollar[3].str), yyDollar[1].num)
		}
	}
	goto yystack /* stack new state and value */
}
