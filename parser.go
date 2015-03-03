package eval

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type parser struct {
	scanner scanner
	tok     token
	lit     string
}

// ParseString parses string and returns expression or error
func ParseString(src string) (expr Expr, err error) {
	return Parse([]byte(src))
}

func Parse(src []byte) (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	p := parser{}
	p.scanner = newScanner(src)
	p.next()
	expr = p.parseExpr(nil)
	p.next()
	if p.tok != EOE {
		return nil, errors.New(p.scanner.errmsg("no matching opening parenthesis"))
	}
	return
}

// ParseIdent is shortcut function to get value from context
func ParseIdent(name string) Expr {
	return &ident{name: name}
}

func (p *parser) next() {
	p.tok, p.lit = p.scanner.scan()
}

func (p *parser) parseExpr(x Expr) Expr {
	if x == nil {
		x = p.parseUnaryExpr()
	}
	if p.tok == RPAREN || p.tok == EOE || p.tok == COMMA {
		return x
	}
	op := p.tok
	p.next()
	y := p.parseUnaryExpr()
	//y := p.parseOperand()
	if p.tok == RPAREN || p.tok == EOE || p.tok == COMMA {
		return &binary{x: x, op: op, y: y}
	}
	if prec[op] >= prec[p.tok] {
		return p.parseExpr(&binary{x: x, op: op, y: y})
	} else {
		return &binary{x: x, op: op, y: p.parseExpr(y)}
	}
}

func (p *parser) parseParenExpr() Expr {
	p.next() // consume opening parenthesis
	x := p.parseExpr(nil)
	if p.tok != RPAREN {
		panic(p.scanner.errmsg("no closing parenthesis"))
	}
	p.next() // consume closing parenthesis
	return x
}

func (p *parser) parseUnaryExpr() Expr {
	switch p.tok {
	case ADD, SUB, NOT:
		op := p.tok
		p.next()
		x := p.parseOperand()
		return &unary{op: op, x: x}
	}
	return p.parseOperand()
}

func (p *parser) parseOperand() Expr {
	switch p.tok {
	case IDENT:
		x := &ident{name: p.lit}
		p.next()
		switch strings.ToUpper(x.name) {
		case "TRUE":
			return &literal{value: true}
		case "FALSE":
			return &literal{value: false}
		case "NULL":
			return &literal{value: nil}
		}
		if p.tok == LPAREN {
			return p.parseCall(x)
		}
		return x
	case STRING:
		x := &literal{value: p.lit}
		p.next()
		return x
	case NUMBER:
		v, ok := new(big.Rat).SetString(p.lit)
		if !ok {
			panic(p.scanner.errmsg("not a number: " + p.lit))
		}
		x := &literal{value: v}
		p.next()
		return x
	case LPAREN:
		return p.parseParenExpr()
	case RPAREN:
		return nil
	}
	panic(p.scanner.errmsg("operand expected"))
}

func (p *parser) parseCall(ident *ident) Expr {
	p.next() // consume LPAREN
	var args []Expr = make([]Expr, 0)
	for {
		arg := p.parseExpr(nil)
		if arg != nil {
			args = append(args, arg)
		}
		if p.tok != COMMA && p.tok != RPAREN {
			panic(p.scanner.errmsg("comma or closing parenthesis expected"))
		}
		if p.tok == COMMA {
			p.next()
		} else if p.tok == RPAREN || p.tok == EOE {
			break
		}
	}
	p.next() // consume RPAREN
	ident.name = strings.ToUpper(ident.name)
	return &call{ident: ident, args: args}
}
