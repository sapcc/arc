//line expr.y:2
package filter

import __yyfmt__ "fmt"

//line expr.y:3
import (
	"fmt"
)

var dbCol = "facts"

//line expr.y:13
type yySymType struct {
	yys int
	str string
	num int
}

const FIELD = 57346
const STRING = 57347
const NUMBER = 57348
const NEQ = 57349
const AND = 57350
const OR = 57351
const NOT = 57352

var yyToknames = []string{
	"FIELD",
	"STRING",
	"NUMBER",
	"'='",
	"'('",
	"')'",
	"NEQ",
	"AND",
	"OR",
	"NOT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line expr.y:105
func stringKey(field string) string {
	return fmt.Sprintf("%s->>'%s'", dbCol, field)
}
func numKey(field string) string {
	return fmt.Sprintf("(%s->>'%s')::numeric", dbCol, field)
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 15
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 32

var yyAct = []int{

	6, 7, 8, 21, 4, 10, 9, 2, 10, 3,
	29, 11, 12, 10, 9, 28, 17, 19, 20, 18,
	15, 13, 27, 16, 14, 24, 25, 22, 23, 26,
	5, 1,
}
var yyPact = []int{

	-4, -1000, 2, -4, -4, -1000, 14, 13, 9, -4,
	-4, -1000, -6, 22, 20, 25, 18, 11, 6, -3,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 31, 7, 30,
}
var yyR1 = []int{

	0, 1, 2, 2, 2, 2, 2, 3, 3, 3,
	3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 2, 3, 3, 3, 1, 3, 3, 3,
	3, 3, 3, 3, 3,
}
var yyChk = []int{

	-1000, -1, -2, 13, 8, -3, 4, 5, 6, 12,
	11, -2, -2, 7, 10, 7, 10, 7, 10, -2,
	-2, 9, 5, 6, 5, 6, 4, 4, 4, 4,
}
var yyDef = []int{

	0, -2, 1, 0, 0, 6, 0, 0, 0, 0,
	0, 2, 0, 0, 0, 0, 0, 0, 0, 4,
	5, 3, 7, 11, 9, 13, 8, 10, 12, 14,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	8, 9, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 7,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 10, 11, 12, 13,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
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

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
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
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
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
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
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
			if yyn < 0 || yyn == yychar {
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
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
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
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
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
		//line expr.y:35
		{
			yylex.(*Lexer).parseResult = yyS[yypt-0].str
		}
	case 2:
		//line expr.y:41
		{
			yyVAL.str = fmt.Sprintf(`NOT %s`, yyS[yypt-0].str)
		}
	case 3:
		//line expr.y:46
		{
			yyVAL.str = fmt.Sprintf(`( %s )`, yyS[yypt-1].str)
		}
	case 4:
		//line expr.y:51
		{
			yyVAL.str = fmt.Sprintf(`( %s OR %s )`, yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 5:
		//line expr.y:56
		{
			yyVAL.str = fmt.Sprintf(`( %s AND %s )`, yyS[yypt-2].str, yyS[yypt-0].str)
		}
	case 6:
		//line expr.y:61
		{
			yyVAL.str = yyS[yypt-0].str
		}
	case 7:
		//line expr.y:67
		{
			yyVAL.str = fmt.Sprintf("%s = '%s'", stringKey(yyS[yypt-2].str), yyS[yypt-0].str)
		}
	case 8:
		//line expr.y:72
		{
			yyVAL.str = fmt.Sprintf("%s = '%s'", stringKey(yyS[yypt-0].str), yyS[yypt-2].str)
		}
	case 9:
		//line expr.y:77
		{
			yyVAL.str = fmt.Sprintf("%s <> '%s'", stringKey(yyS[yypt-2].str), yyS[yypt-0].str)
		}
	case 10:
		//line expr.y:82
		{
			yyVAL.str = fmt.Sprintf("%s <> '%s'", stringKey(yyS[yypt-0].str), yyS[yypt-2].str)
		}
	case 11:
		//line expr.y:87
		{
			yyVAL.str = fmt.Sprintf("%s = %d", numKey(yyS[yypt-2].str), yyS[yypt-0].num)
		}
	case 12:
		//line expr.y:92
		{
			yyVAL.str = fmt.Sprintf("%s = %d", numKey(yyS[yypt-0].str), yyS[yypt-2].num)
		}
	case 13:
		//line expr.y:97
		{
			yyVAL.str = fmt.Sprintf("%s <> %d", numKey(yyS[yypt-2].str), yyS[yypt-0].num)
		}
	case 14:
		//line expr.y:102
		{
			yyVAL.str = fmt.Sprintf("%s <> %d", numKey(yyS[yypt-0].str), yyS[yypt-2].num)
		}
	}
	goto yystack /* stack new state and value */
}
