package eval

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
	"unicode"
)

const (
	ISO8601 string = "2006-01-02T15:04:05.999Z0700"
)

func builtin(name string, args []interface{}, context *context) (val interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	switch name {
	//salesforce text functions: CASESAFEID, GETSESSIONID, HYPERLINK, IMAGE, ISPICKVAL not implemented as too specific
	// TEXT should be avoided, use FORMAT instead
	case "BEGINS":
		NumOfParams(args, 2, "BEGINS")
		s1 := MustBeString(args, 0, "BEGINS")
		s2 := MustBeString(args, 1, "BEGINS")
		return strings.Index(s1, s2) == 0, nil
	case "CONTAINS":
		NumOfParams(args, 2, "CONTAINS")
		s1 := MustBeString(args, 0, "CONTAINS")
		s2 := MustBeString(args, 1, "CONTAINS")
		return strings.Index(s1, s2) != -1, nil
	case "FIND":
		NumOfParams(args, 2, "FIND")
		s1 := MustBeString(args, 0, "FIND")
		s2 := MustBeString(args, 1, "FIND")
		return new(big.Rat).SetInt64(int64(strings.Index(s1, s2))), nil
	case "INCLUDES":
		NumOfParams(args, 2, "INCLUDES")
		s1 := MustBeString(args, 0, "INCLUDES")
		s2 := MustBeString(args, 1, "INCLUDES")
		return strings.Index(";"+s1+";", ";"+s2+";") != -1, nil
	case "LEFT":
		NumOfParams(args, 2, "LEFT")
		s1 := MustBeString(args, 0, "LEFT")
		n1 := GetNumberAsInt(args, 1, "LEFT")
		return substr(s1, 0, n1), nil
	case "LEN":
		NumOfParams(args, 1, "LEN")
		s1 := MustBeString(args, 0, "LEN")
		return new(big.Rat).SetInt64(int64(len([]rune(s1)))), nil
	case "LOWER":
		NumOfParams(args, 1, "LOWER")
		s1 := MustBeString(args, 0, "LOWER")
		return strings.ToLower(s1), nil
	case "LPAD":
		n1 := GetNumberAsInt(args, 1, "LPAD")
		s1 := MustBeString(args, 0, "LPAD")
		s2 := " "
		if len(args) > 2 {
			NumOfParams(args, 3, "LPAD")
			s2 = MustBeString(args, 2, "LPAD")
		}
		s1r := []rune(s1)
		s2r := []rune(s2)
		if n1 == len(s1r) {
			return s1, nil
		} else if n1 < len(s1r) {
			return string(s1r[:n1]), nil
		}
		n1 = n1 - len(s1r)
		var sr []rune
		for i := 0; len(sr) < n1; i++ {
			sr = append(sr, s2r[i%len(s2r)])
		}
		return string(append(sr, s1r...)), nil
	case "MID":
		NumOfParams(args, 3, "MID")
		s1 := MustBeString(args, 0, "MID")
		n1 := GetNumberAsInt(args, 1, "MID")
		if n1 <= 0 {
			n1 = 1
		}
		n2 := GetNumberAsInt(args, 2, "MID")
		if n2 <= 0 {
			return "", nil
		}
		return substr(s1, n1-1, n2), nil
	case "RIGHT":
		NumOfParams(args, 2, "RIGHT")
		s1 := MustBeString(args, 0, "RIGHT")
		n1 := GetNumberAsInt(args, 1, "RIGHT")
		if n1 <= 0 {
			return "", nil
		}
		return substr(s1, -n1, -1), nil
	case "RPAD":
		n1 := GetNumberAsInt(args, 1, "LPAD")
		s1 := MustBeString(args, 0, "LPAD")
		s2 := " "
		if len(args) > 2 {
			NumOfParams(args, 3, "LPAD")
			s2 = MustBeString(args, 2, "LPAD")
		}
		s1r := []rune(s1)
		s2r := []rune(s2)
		if n1 == len(s1r) {
			return s1, nil
		} else if n1 < len(s1r) {
			return string(s1r[:n1]), nil
		}
		for i := 0; len(s1r) < n1; i++ {
			s1r = append(s1r, s2r[i%len(s2r)])
		}
		return string(s1r), nil
	case "SUBSTITUTE":
		NumOfParams(args, 3, "SUBSTITUTE")
		s1 := MustBeString(args, 0, "SUBSTITUTE")
		s2 := MustBeString(args, 1, "SUBSTITUTE")
		s3 := MustBeString(args, 2, "SUBSTITUTE")
		return strings.Replace(s1, s2, s3, -1), nil
	case "TEXT":
		NumOfParams(args, 1, "TEXT")
		v1 := args[0]
		switch v1.(type) {
		case string:
			return v1.(string), nil
		case big.Rat:
			return v1.(*big.Rat).String(), nil
		case bool:
			if v1.(bool) {
				return "true", nil
			} else {
				return "false", nil
			}
		case time.Time:
			return v1.(time.Time).Format(ISO8601), nil
		}
		return nil, errors.New(fmt.Sprint("function TEXT: unknown parameter type:", v1))
	case "TRIM":
		NumOfParams(args, 1, "TRIM")
		s1 := MustBeString(args, 0, "TRIM")
		s1 = strings.TrimSpace(s1)
		s1r := []rune(s1)
		sr := make([]rune, 0, len(s1r))
		var first bool = false
		for _, r := range s1r {
			sp := unicode.IsSpace(r)
			if first && sp {
				sr = append(sr, ' ')
				first = false
			} else if !sp {
				sr = append(sr, r)
				first = true
			}
		}
		return string(sr), nil
	case "UPPER":
		NumOfParams(args, 1, "UPPER")
		s1 := MustBeString(args, 0, "UPPER")
		return strings.ToUpper(s1), nil
	case "VALUE":
		NumOfParams(args, 1, "VALUE")
		s1 := MustBeString(args, 0, "VALUE")
		f, ok := new(big.Rat).SetString(s1)
		if ok {
			return f, nil
		}
		return nil, errors.New(fmt.Sprint("function VALUE: not a number:", s1))
	// salesforce date functions. DATEVALUE accepts additional format parameter
	case "DATE":
		NumOfParams(args, 3, "DATE")
		n1 := GetNumberAsInt(args, 0, "DATE")
		n2 := GetNumberAsInt(args, 1, "DATE")
		n3 := GetNumberAsInt(args, 2, "DATE")
		return time.Date(n1, time.Month(n2), n3, 0, 0, 0, 0, context.localTimeZone), nil
	case "DATEVALUE":
		s1 := MustBeString(args, 0, "DATEVALUE")
		s2 := "2006-01-02"
		if len(args) > 1 {
			NumOfParams(args, 2, "DATEVALUE")
			s2 = MustBeString(args, 1, "DATEVALUE")
		}
		return context.ParseDate(s2, s1)
	case "DAY":
		NumOfParams(args, 1, "DAY")
		d1 := MustBeDate(args, 0, "DAY")
		return new(big.Rat).SetInt64(int64(d1.Day())), nil
	case "MONTH":
		NumOfParams(args, 1, "MONTH")
		d1 := MustBeDate(args, 0, "MONTH")
		return new(big.Rat).SetInt64(int64(d1.Month())), nil
	case "NOW":
		NumOfParams(args, 0, "NOW")
		return time.Now().In(context.localTimeZone), nil
	case "TODAY":
		NumOfParams(args, 0, "TODAY")
		return time.Now().In(context.localTimeZone).Truncate(time.Hour * 24), nil
	case "YEAR":
		NumOfParams(args, 1, "YEAR")
		d1 := MustBeDate(args, 0, "YEAR")
		return new(big.Rat).SetInt64(int64(d1.Year())), nil
		// salesforce logical functions: AND, NOT, OR should not be used, use logical operators
	case "CASE":
		if len(args) < 3 {
			return nil, errors.New(fmt.Sprint("function CASE: expected at least 3 parameters, actual: ", len(args)))
		}
		i := 1
		for ; i < len(args)-1; i += 2 {
			switch args[0].(type) {
			case nil:
				if args[i] == nil {
					return args[i+1], nil
				}
			case *big.Rat:
				n := MustBeNumber(args, i, "CASE")
				if args[0].(*big.Rat).Cmp(n) == 0 {
					return args[i+1], nil
				}
			case string:
				s := MustBeString(args, i, "CASE")
				if args[0].(string) == s {
					return args[i+1], nil
				}
			case bool:
				v, ok := args[i].(bool)
				if !ok {
					return nil, errors.New(fmt.Sprint("function DECODE: parameter", i, "not a boolean:", args[i]))
				}
				if args[0].(bool) == v {
					return args[i+1], nil
				}
			default:
				return nil, errors.New(fmt.Sprint("function DECODE: illegal type in parameter", i, ":", args[i]))
			}
		}
		if i < len(args) {
			return args[i], nil
		}
		return nil, errors.New("function CASE: no value found")
	case "IF":
		if len(args) != 3 {
			return nil, errors.New(fmt.Sprint("function IF: expected 3 parameters, actual:", len(args)))
		}
		cond, ok := args[0].(bool)
		if !ok {
			return nil, errors.New(fmt.Sprint("function IF: expected boolean as first parameter, actual:", args[0]))
		}
		if cond {
			return args[1], nil
		}
		return args[2], nil
	case "NULLVALUE":
		if len(args) != 2 {
			return nil, errors.New(fmt.Sprint("function NULLVALUE: expected 2 parameters, actual:", len(args)))
		}
		if args[0] == nil {
			return args[1], nil
		}
		return args[0], nil
	// additional convenience functions for text
	// join(delimiter, strings...) joins non empty strings listed as arguments using delimiter (empty strings are skipped)
	case "JOIN":
		var a []string
		for i, v := range args[1:] {
			if v != nil {
				s := MustBeString(args, i, "JOIN")
				if s != "" {
					a = append(a, s)
				}
			}
		}
		s1 := MustBeString(args, 0, "JOIN")
		return strings.Join(a, s1), nil
	}
	return nil, errors.New(fmt.Sprint("unknown function: ", name))
}

func substr(s string, b int, l int) string {
	r := []rune(s)
	if len(r) == 0 || l == 0 {
		return ""
	}
	if b < 0 {
		if -b >= len(r) {
			b = 0
		} else {
			b = len(r) + b
		}
	}
	if l < 0 {
		l = len(r)
	}
	e := b + l
	if e > len(r) {
		e = len(r)
	}
	return string(r[b:e])
}