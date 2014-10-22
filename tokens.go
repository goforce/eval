package eval

type token int

const (
	EOE token = iota
	BAD
	IDENT
	STRING
	NUMBER
	BOOL
	COMMA
	LPAREN
	RPAREN
	NOT
	ADD
	SUB
	MUL
	DIV
	EQ
	NEQ
	LT
	LTE
	GT
	GTE
	AND
	OR
)

// sequence is important!
var tokens = []struct {
	tok token
	seq []rune
}{
	{COMMA, []rune{','}},
	{LPAREN, []rune{'('}},
	{RPAREN, []rune{')'}},
	{ADD, []rune{'+'}},
	{SUB, []rune{'-'}},
	{MUL, []rune{'*'}},
	{DIV, []rune{'/'}},
	{EQ, []rune{'=', '='}},
	{NEQ, []rune{'!', '='}},
	{NOT, []rune{'!'}},
	{NEQ, []rune{'<', '>'}},
	{LTE, []rune{'<', '='}},
	{LT, []rune{'<'}},
	{GTE, []rune{'>', '='}},
	{GT, []rune{'>'}},
	{AND, []rune{'&', '&'}},
	{OR, []rune{'|', '|'}},
}

func repr(tok token) string {
	for _, v := range tokens {
		if tok == v.tok {
			return string(v.seq)
		}
	}
	return "unknown token"
}

var prec = map[token]int{
	MUL: 4,
	DIV: 4,
	ADD: 3,
	SUB: 3,
	EQ:  2,
	NEQ: 2,
	LT:  2,
	LTE: 2,
	GT:  2,
	GTE: 2,
	AND: 1,
	OR:  1,
}
