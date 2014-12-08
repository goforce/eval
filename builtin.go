package eval

import (
	"fmt"
	"math/big"
	"strings"
	"time"
	"unicode"
)

const (
	ISO8601 string = "2006-01-02T15:04:05.999Z0700"
)

func builtin(name string, args []interface{}, context Context) (val interface{}, err error) {
	switch name {
	//salesforce text functions: CASESAFEID, GETSESSIONID, HYPERLINK, IMAGE, ISPICKVAL not implemented as too specific
	// TEXT should be avoided, use FORMAT instead
	case "BEGINS":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		s2 := MustBeString(args, 1)
		return strings.Index(s1, s2) == 0, nil
	case "CONTAINS":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		s2 := MustBeString(args, 1)
		return strings.Index(s1, s2) != -1, nil
	case "FIND":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		s2 := MustBeString(args, 1)
		return new(big.Rat).SetInt64(int64(strings.Index(s1, s2) + 1)), nil
	case "INCLUDES":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		s2 := MustBeString(args, 1)
		return strings.Index(";"+s1+";", ";"+s2+";") != -1, nil
	case "LEFT":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		n1 := GetNumberAsInt(args, 1)
		return substr(s1, 0, n1), nil
	case "LEN":
		NumOfParams(args, 1)
		s1 := MustBeString(args, 0)
		return new(big.Rat).SetInt64(int64(len([]rune(s1)))), nil
	case "LOWER":
		NumOfParams(args, 1)
		s1 := MustBeString(args, 0)
		return strings.ToLower(s1), nil
	case "LPAD":
		n1 := GetNumberAsInt(args, 1)
		s1 := MustBeString(args, 0)
		s2 := " "
		if len(args) > 2 {
			NumOfParams(args, 3)
			s2 = MustBeString(args, 2)
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
		NumOfParams(args, 3)
		s1 := MustBeString(args, 0)
		n1 := GetNumberAsInt(args, 1)
		if n1 <= 0 {
			n1 = 1
		}
		n2 := GetNumberAsInt(args, 2)
		if n2 <= 0 {
			return "", nil
		}
		return substr(s1, n1-1, n2), nil
	case "RIGHT":
		NumOfParams(args, 2)
		s1 := MustBeString(args, 0)
		n1 := GetNumberAsInt(args, 1)
		if n1 <= 0 {
			return "", nil
		}
		return substr(s1, -n1, -1), nil
	case "RPAD":
		n1 := GetNumberAsInt(args, 1)
		s1 := MustBeString(args, 0)
		s2 := " "
		if len(args) > 2 {
			NumOfParams(args, 3)
			s2 = MustBeString(args, 2)
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
		NumOfParams(args, 3)
		s1 := MustBeString(args, 0)
		s2 := MustBeString(args, 1)
		s3 := MustBeString(args, 2)
		return strings.Replace(s1, s2, s3, -1), nil
	case "TEXT":
		NumOfParams(args, 1)
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
		panic(fmt.Sprint("unsupported type:", v1))
	case "TRIM":
		NumOfParams(args, 1)
		s1 := MustBeString(args, 0)
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
		NumOfParams(args, 1)
		s1 := MustBeString(args, 0)
		return strings.ToUpper(s1), nil
	case "ISBLANK":
		NumOfParams(args, 1)
		if args[0] == nil {
			return true, nil
		} else if s, ok := args[0].(string); ok && s == "" {
			return true, nil
		}
		return false, nil
	case "VALUE":
		NumOfParams(args, 1)
		s1 := MustBeString(args, 0)
		f, ok := new(big.Rat).SetString(s1)
		if ok {
			return f, nil
		}
		panic(fmt.Sprint("not a number", s1))
	// salesforce date functions. DATEVALUE accepts additional format parameter
	case "DATE":
		NumOfParams(args, 3)
		n1 := GetNumberAsInt(args, 0)
		n2 := GetNumberAsInt(args, 1)
		n3 := GetNumberAsInt(args, 2)
		return time.Date(n1, time.Month(n2), n3, 0, 0, 0, 0, context.cast().localTimeZone), nil
	case "DATEVALUE":
		s1 := MustBeString(args, 0)
		s2 := "2006-01-02"
		if len(args) > 1 {
			NumOfParams(args, 2)
			s2 = MustBeString(args, 1)
		}
		return context.ParseDate(s2, s1)
	case "DAY":
		NumOfParams(args, 1)
		d1 := MustBeDate(args, 0)
		return new(big.Rat).SetInt64(int64(d1.Day())), nil
	case "MONTH":
		NumOfParams(args, 1)
		d1 := MustBeDate(args, 0)
		return new(big.Rat).SetInt64(int64(d1.Month())), nil
	case "NOW":
		NumOfParams(args, 0)
		return time.Now().In(context.cast().localTimeZone), nil
	case "TODAY":
		NumOfParams(args, 0)
		return time.Now().In(context.cast().localTimeZone).Truncate(time.Hour * 24), nil
	case "YEAR":
		NumOfParams(args, 1)
		d1 := MustBeDate(args, 0)
		return new(big.Rat).SetInt64(int64(d1.Year())), nil
	// salesforce numerical functions: TBC
	case "ABS":
		NumOfParams(args, 1)
		n1 := MustBeNumber(args, 0)
		return n1.Abs(n1), nil
	// salesforce logical functions: AND, NOT, OR should not be used, use logical operators
	case "CASE":
		MinNumOfParams(args, 3)
		i := 1
		for ; i < len(args)-1; i += 2 {
			switch args[0].(type) {
			case nil:
				if args[i] == nil {
					return args[i+1], nil
				}
			case *big.Rat:
				n := MustBeNumber(args, i)
				if args[0].(*big.Rat).Cmp(n) == 0 {
					return args[i+1], nil
				}
			case string:
				s := MustBeString(args, i)
				if args[0].(string) == s {
					return args[i+1], nil
				}
			case bool:
				b1 := MustBeBool(args, i)
				if args[0].(bool) == b1 {
					return args[i+1], nil
				}
			default:
				panic(fmt.Sprint("unsupported type:", args[i]))
			}
		}
		if i < len(args) {
			return args[i], nil
		}
		panic("missing default value")
	case "IF":
		NumOfParams(args, 3)
		cond := MustBeBool(args, 0)
		if cond {
			return args[1], nil
		}
		return args[2], nil
	case "NULLVALUE":
		NumOfParams(args, 2)
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
				s := MustBeString(args, i+1)
				if s != "" {
					a = append(a, s)
				}
			}
		}
		s1 := MustBeString(args, 0)
		return strings.Join(a, s1), nil
	}
	panic("unknown function")
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
