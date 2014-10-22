package eval

import (
	"errors"
	"fmt"
	"math/big"
	"time"
)

type Expr interface {
	Eval(*context) (interface{}, error)
	String() string
}

// use after calls to values or functions to ensure that no illegal types have been returned
func validate(v interface{}, name string) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	switch v.(type) {
	case string:
	case *big.Rat:
	case bool:
	case time.Time:
	default:
		if name == "" {
			return nil, errors.New("illegal value: '" + fmt.Sprint(v) + "'")
		} else {
			return nil, errors.New("illegal value: '" + fmt.Sprint(v) + "' in " + name)
		}
	}
	return v, nil
}

type ident struct {
	name string
}

func (e *ident) Eval(context *context) (interface{}, error) {
	for _, fn := range context.values {
		if v, ok := fn(e.name); ok {
			return validate(v, e.name)
		}
	}
	return nil, errors.New("unknown value: " + e.name)
}

func (e ident) String() string {
	return e.name
}

type literal struct {
	value interface{}
}

func (e *literal) Eval(context *context) (interface{}, error) {
	return e.value, nil
}

func (e literal) String() string {
	return fmt.Sprint("'", e.value, "'")
}

type call struct {
	ident *ident
	args  []Expr
}

func (e *call) Eval(context *context) (interface{}, error) {
	var list []interface{}
	for _, p := range e.args {
		v, err := p.Eval(context)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	for _, fn := range context.functions {
		if v, err := fn(e.ident.name, list); err == nil {
			return validate(v, e.ident.name)
		} else if _, ok := err.(NOFUNC); !ok {
			return nil, err
		}
	}
	v, err := builtin(e.ident.name, list, context)
	if err != nil {
		return nil, err
	}
	return validate(v, e.ident.name)
}

func (e call) String() string {
	return fmt.Sprint(e.ident, "(", e.args, ")")
}

type unary struct {
	op token
	x  Expr
}

func (e *unary) Eval(context *context) (interface{}, error) {
	v, err := e.x.Eval(context)
	if err != nil {
		return nil, err
	}
	switch e.op {
	case NOT:
		switch v.(type) {
		case bool:
			return !v.(bool), nil
		default:
			return nil, errors.New("not a boolean:" + fmt.Sprint(v))
		}
	case ADD:
		switch v.(type) {
		case *big.Rat:
			return v.(*big.Rat), nil
		default:
			return nil, errors.New("not a number:" + fmt.Sprint(v))
		}
	case SUB:
		switch v.(type) {
		case *big.Rat:
			return new(big.Rat).Neg(v.(*big.Rat)), nil
		default:
			return nil, errors.New("not a number:" + fmt.Sprint(v))
		}
	}
	return nil, errors.New("illegal unary operator" + string(e.op))
}

func (e unary) String() string {
	return fmt.Sprint(" ( ", e.op, e.x, " ) ")
}

type binary struct {
	x  Expr
	op token
	y  Expr
}

func (e *binary) Eval(context *context) (interface{}, error) {
	ix, err := e.x.Eval(context)
	if err != nil {
		return nil, err
	}
	iy, err := e.y.Eval(context)
	if err != nil {
		return nil, err
	}
	switch e.op {
	case ADD, LT, LTE, GT, GTE:
		r, ok, s := tryNumbers(ix, iy, e.op)
		if ok {
			return r, nil
		}
		r, ok, s = tryStrings(ix, iy, e.op)
		if ok {
			return r, nil
		}
		return nil, errors.New("not a string:" + fmt.Sprint(s))
	case MUL, DIV, SUB:
		r, ok, s := tryNumbers(ix, iy, e.op)
		if ok {
			return r, nil
		}
		return nil, errors.New("not a number:" + fmt.Sprint(s))
	case EQ, NEQ:
		r, ok, s := tryNumbers(ix, iy, e.op)
		if ok {
			return r, nil
		}
		r, ok, s = tryBools(ix, iy, e.op)
		if ok {
			return r, nil
		}
		r, ok, s = tryStrings(ix, iy, e.op)
		if ok {
			return r, nil
		}
		return nil, errors.New("not a string:" + fmt.Sprint(s))
	case AND, OR:
		r, ok, s := tryBools(ix, iy, e.op)
		if ok {
			return r, nil
		}
		return nil, errors.New("not a boolean:" + fmt.Sprint(s))
	}
	// TODO single equal sign not shown, error reporting should be fixed
	return nil, errors.New("illegal binary operator" + string(e.op))
}

func (e binary) String() string {
	return fmt.Sprint(" ( ", e.x, e.op, e.y, " ) ")
}

func tryNumbers(ix, iy interface{}, op token) (interface{}, bool, interface{}) {
	x, ok := ix.(*big.Rat)
	if !ok {
		return nil, false, ix
	}
	y, ok := iy.(*big.Rat)
	if !ok {
		return nil, false, iy
	}
	switch op {
	case ADD:
		return new(big.Rat).Add(x, y), true, nil
	case SUB:
		return new(big.Rat).Sub(x, y), true, nil
	case MUL:
		return new(big.Rat).Mul(x, y), true, nil
	case DIV:
		return new(big.Rat).Quo(x, y), true, nil
	case EQ:
		return x.Cmp(y) == 0, true, nil
	case NEQ:
		return x.Cmp(y) != 0, true, nil
	case LT, GTE:
		return x.Cmp(y) == -1, true, nil
	case GT, LTE:
		return x.Cmp(y) == 1, true, nil
	}
	return nil, false, nil
}

func tryStrings(ix, iy interface{}, op token) (interface{}, bool, interface{}) {
	x, ok := ix.(string)
	if !ok && ix != nil {
		return nil, false, ix
	}
	y, ok := iy.(string)
	if !ok && iy != nil {
		return nil, false, iy
	}
	switch op {
	case ADD:
		return x + y, true, nil
	case EQ:
		return x == y, true, nil
	case NEQ:
		return x != y, true, nil
	case LT:
		return x < y, true, nil
	case LTE:
		return x <= y, true, nil
	case GT:
		return x > y, true, nil
	case GTE:
		return x >= y, true, nil
	}
	return nil, false, nil
}

func tryBools(ix, iy interface{}, op token) (interface{}, bool, interface{}) {
	x, ok := ix.(bool)
	if !ok {
		return nil, false, ix
	}
	y, ok := iy.(bool)
	if !ok {
		return nil, false, iy
	}
	switch op {
	case AND:
		return x && y, true, nil
	case OR:
		return x || y, true, nil
	case EQ:
		return x == y, true, nil
	case NEQ:
		return x != y, true, nil
	}
	return nil, false, nil
}
