package eval

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

func test_values(name string) (interface{}, bool) {
	switch name {
	case "string_a":
		return "a string", true
	case "string_b":
		return "b string", true
	case "number_1":
		return big.NewRat(1, 1), true
	default:
		return nil, false
	}
}

func test_functions(name string, args []interface{}) (interface{}, error) {
	return nil, NOFUNC{}
}

var test_context = NewContext().AddValues(test_values).AddFunctions(test_functions).SetTimeZone(time.FixedZone("MY", 0))

func TestParsing(t *testing.T) {
	{
		se := "2+2)  *2-1"
		e, err := ParseString(se)
		if err == nil {
			t.Error("parsing of missing opening parenthesis failed:", se)
		}
		if e != nil {
			t.Error("parsing of missing opening parenthesis failed, expression returned:", se)
		}
	}
}

func TestStringFunctions(t *testing.T) {
	s1 := " this \t is \n english "
	s1len := len([]rune(s1))
	s2 := " Вот ;есть; ;одна вещь "
	s2len := len([]rune(s2))

	mustErrorEvaluating(t, "begins()", "function begins: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "begins('',0,'')", "function begins: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "begins('')", "function begins: failed to check number of parameters, 1 parameter")
	mustErrorEvaluating(t, "begins('',true)", "function begins: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "begins('',0)", "function begins: failed to check type of parameters, number")
	mustResult(t, "begins('"+s1+"',' this')", true)
	mustResult(t, "begins('"+s1+"','this')", false)
	mustResult(t, "begins('"+s2+"','Вот ')", false)

	mustErrorEvaluating(t, "contains()", "function contains: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "contains('',0,'')", "function contains: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "contains('')", "function contains: failed to check number of parameters, 1 parameter")
	mustErrorEvaluating(t, "contains('',true)", "function contains: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "contains('',0)", "function contains: failed to check type of parameters, number")
	mustResult(t, "contains('"+s1+"',' this')", true)
	mustResult(t, "contains('"+s1+"','this')", true)
	mustResult(t, "contains('"+s2+"','Вот ')", true)
	mustResult(t, "contains('"+s2+"','Во т')", false)

	mustErrorEvaluating(t, "find()", "function find: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "find('',0,'')", "function find: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "find('')", "function find: failed to check number of parameters, 1 parameter")
	mustErrorEvaluating(t, "find('',true)", "function find: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "find('',0)", "function find: failed to check type of parameters, number")
	mustResult(t, "find('"+s1+"',' this')", big.NewRat(1, 1))
	mustResult(t, "find('"+s1+"','this')", big.NewRat(2, 1))
	mustResult(t, "find('"+s2+"','Вот ')", big.NewRat(2, 1))
	mustResult(t, "find('"+s2+"','Во т')", big.NewRat(0, 1))

	mustErrorEvaluating(t, "includes()", "function includes: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "includes('',0,'')", "function includes: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "includes('')", "function includes: failed to check number of parameters, 1 parameter")
	mustErrorEvaluating(t, "includes('',true)", "function includes: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "includes('',0)", "function includes: failed to check type of parameters, number")
	mustResult(t, "includes('"+s1+"',' this')", false)
	mustResult(t, "includes('"+s1+"','this')", false)
	mustResult(t, "includes('"+s2+"','Вот ')", false)
	mustResult(t, "includes('"+s2+"',' Вот ')", true)
	mustResult(t, "includes('"+s2+"',' ')", true)

	mustErrorEvaluating(t, "left()", "function left: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "left('',0,'')", "function left: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "left('')", "function left: failed to check number of parameters, 1 parameter")
	mustResult(t, "left('"+s1+"',3)", string([]rune(s1)[:3]))
	mustResult(t, "left('"+s2+"',3)", string([]rune(s2)[:3]))
	mustResult(t, "left('"+s2+"','3')", string([]rune(s2)[:3]))
	mustResult(t, "left('"+s2+"',0)", "")
	mustResult(t, "left('"+s2+"',"+strconv.Itoa(s2len)+")", s2)
	mustResult(t, "left('"+s2+"',"+strconv.Itoa(s2len+1)+")", s2)

	mustErrorEvaluating(t, "len()", "function len: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "len('','')", "function len: failed to check number of parameters, 2 parameters")
	mustErrorEvaluating(t, "len(0)", "function len: failed to check type of parameter, number parameter")
	mustResult(t, "len('"+s1+"')", big.NewRat(int64(s1len), 1))
	mustResult(t, "len('"+s2+"')", big.NewRat(int64(s2len), 1))

	mustErrorEvaluating(t, "lower()", "function lower: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "lower('','')", "function lower: failed to check number of parameters, 2 parameters")
	mustErrorEvaluating(t, "lower(true)", "function lower: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "lower(0)", "function lower: failed to check type of parameters, number")
	mustResult(t, "lower('"+s1+"')", strings.ToLower(s1))
	mustResult(t, "lower('"+s2+"')", strings.ToLower(s2))

	mustErrorEvaluating(t, "lpad()", "function lpad: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "lpad(true,0)", "function lpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "lpad(0,0)", "function lpad: failed to check type of parameters, number")
	mustErrorEvaluating(t, "lpad('','a')", "function lpad: failed to check type of parameters, not a number")
	mustErrorEvaluating(t, "lpad('',true)", "function lpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "lpad('',0,1)", "function lpad: failed to check type of parameters, number")
	mustErrorEvaluating(t, "lpad('',0,true)", "function lpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "lpad('',0,'','')", "function lpad: failed to check number of parameters, more than 3 parameters")
	mustResult(t, "lpad('"+s1+"',"+strconv.Itoa(s1len)+")", s1)
	mustResult(t, "lpad('"+s1+"',"+strconv.Itoa(s1len-1)+")", string([]rune(s1)[0:s1len-1]))
	mustResult(t, "lpad('"+s1+"',"+strconv.Itoa(s1len+1)+")", " "+s1)
	mustResult(t, "lpad('"+s2+"',"+strconv.Itoa(s2len)+")", s2)
	mustResult(t, "lpad('"+s2+"',"+strconv.Itoa(s2len-1)+")", string([]rune(s2)[0:s2len-1]))
	mustResult(t, "lpad('"+s2+"',"+strconv.Itoa(s2len+5)+",'Вот')", "ВотВо"+s2)

	mustErrorEvaluating(t, "mid()", "function mid: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "mid('',0)", "function mid: failed to check number of parameters, 2 parameters")
	mustErrorEvaluating(t, "mid('',0,0,0)", "function mid: failed to check number of parameters, 4 parameter")
	mustResult(t, "mid('"+s1+"',3,5)", string([]rune(s1)[2:7]))
	mustResult(t, "mid('"+s2+"',3,5)", string([]rune(s2)[2:7]))
	mustResult(t, "mid('"+s2+"','3','1')", string([]rune(s2)[2:3]))
	mustResult(t, "mid('"+s2+"',0,-1)", "")
	mustResult(t, "mid('"+s2+"',0,0)", "")
	mustResult(t, "mid('"+s2+"',0,"+strconv.Itoa(s2len)+")", s2)
	mustResult(t, "mid('"+s2+"',0,"+strconv.Itoa(s2len+1)+")", s2)

	mustErrorEvaluating(t, "right()", "function right: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "right('',0,'')", "function right: failed to check number of parameters, 3 parameters")
	mustErrorEvaluating(t, "right('')", "function right: failed to check number of parameters, 1 parameter")
	mustResult(t, "right('"+s1+"',3)", string([]rune(s1)[s1len-3:]))
	mustResult(t, "right('"+s2+"',3)", string([]rune(s2)[s2len-3:]))
	mustResult(t, "right('"+s2+"','3')", string([]rune(s2)[s2len-3:]))
	mustResult(t, "right('"+s2+"',0)", "")
	mustResult(t, "right('"+s2+"',"+strconv.Itoa(s2len)+")", s2)
	mustResult(t, "right('"+s2+"',"+strconv.Itoa(s2len+1)+")", s2)

	mustErrorEvaluating(t, "rpad()", "function rpad: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "rpad(true,0)", "function rpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "rpad(0,0)", "function rpad: failed to check type of parameters, number")
	mustErrorEvaluating(t, "rpad('','a')", "function rpad: failed to check type of parameters, not a number")
	mustErrorEvaluating(t, "rpad('',true)", "function rpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "rpad('',0,1)", "function rpad: failed to check type of parameters, number")
	mustErrorEvaluating(t, "rpad('',0,true)", "function rpad: failed to check type of parameters, boolean")
	mustErrorEvaluating(t, "rpad('',0,'','')", "function rpad: failed to check number of parameters, more than 3 parameters")
	mustResult(t, "rpad('"+s1+"',"+strconv.Itoa(s1len)+")", s1)
	mustResult(t, "rpad('"+s1+"',"+strconv.Itoa(s1len-1)+")", string([]rune(s1)[0:s1len-1]))
	mustResult(t, "rpad('"+s1+"',"+strconv.Itoa(s1len+1)+")", s1+" ")
	mustResult(t, "rpad('"+s2+"',"+strconv.Itoa(s2len)+")", s2)
	mustResult(t, "rpad('"+s2+"',"+strconv.Itoa(s2len-1)+")", string([]rune(s2)[0:s2len-1]))
	mustResult(t, "rpad('"+s2+"',"+strconv.Itoa(s2len+5)+",'Вот')", s2+"ВотВо")

	mustErrorEvaluating(t, "substitute()", "function substitute: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "substitute('','')", "function substitute: failed to check number of parameters, 2 parameters")
	mustErrorEvaluating(t, "substitute(0,0,0)", "function substitute: failed to check type of parameters, number parameter")
	mustResult(t, "substitute('replace me',' ','5')", "replace5me")
	mustResult(t, "substitute('Вот есть одна вещь','одна','вещь')", "Вот есть вещь вещь")
	mustResult(t, "replace('Вот есть одна вещь','Вот','вещь')", "вещь есть одна вещь")

	mustErrorEvaluating(t, "text()", "function text: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "text('','')", "function text: failed to check number of parameters, 2 parameters")
	mustResult(t, "text('"+s1+"')", s1)
	mustResult(t, "text(314)", "314")
	mustResult(t, "text(true)", "true")
	//	mustResult(t, "text(datetimevalue('2001-01-02T01:02:03Z'))", "2001-01-02T01:02:03Z")

	mustErrorEvaluating(t, "datevalue()", "function datevalue: failed to check number of parameters, no parameters")
	mustErrorEvaluating(t, "datevalue(0)", "function datevalue: failed to check type of parameters, number parameter")
	mustResult(t, "datevalue('2001-01-02')", time.Date(2001, 01, 02, 0, 0, 0, 0, test_context.(*context).localTimeZone))

}

func mustErrorEvaluating(t *testing.T, expression string, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = " should be:" + strings.Join(message, ", ")
	}
	e, err := ParseString(expression)
	if err != nil {
		t.Error("failed to parse:", expression, err)
		return
	}
	v, err := e.Eval(test_context)
	if err == nil {
		t.Error("no error returned on evaluate:", expression, msg, " instead value returned:", v)
		return
	}
}

func mustResult(t *testing.T, expression string, value interface{}) {
	e, err := ParseString(expression)
	if err != nil {
		t.Error("failed to parse:", expression, " error at parsing: ", err)
		return
	}
	v, err := e.Eval(test_context)
	if err != nil {
		t.Error("failed to evaluate:", expression, " error at evaluation: ", err)
		return
	}
	if r, ok := v.(*big.Rat); ok {
		if r.Cmp(value.(*big.Rat)) != 0 {
			t.Error("failed to evaluate:", expression, " expected:", value, " actual:", v)
		}
	} else if v != value {
		t.Error("failed to evaluate:", expression, " expected:", value, " actual:", v)
		return
	}
}
