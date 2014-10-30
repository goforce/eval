// Helper functions to facilitate creation of custom Functions. All functions do panic on error.
// So panic should be translated to error
// 	defer func() {
//		if r := recover(); r != nil {
//			err = errors.New(fmt.Sprint(r))
//		}
//	}()
//

package eval

import (
	"fmt"
	"math/big"
	"time"
)

type NOFUNC struct{}

func (e NOFUNC) Error() string {
	return "function not defined"
}

func NumOfParams(args []interface{}, expected int, name string) {
	if expected != len(args) {
		panic(fmt.Sprint("function ", name, ": expected ", expected, " parameters, actual ", len(args)))
	}
}

func MustBeString(args []interface{}, index int, name string) string {
	val, ok := args[index].(string)
	if !ok {
		panic(fmt.Sprint("function ", name, ": parameter ", index, " not a string ", args[index]))
	}
	return val
}

func MustBeNumber(args []interface{}, index int, name string) *big.Rat {
	val, ok := args[index].(*big.Rat)
	if !ok {
		panic(fmt.Sprint("function ", name, ": parameter ", index, " not a number ", args[index]))
	}
	return val
}

func MustBeDate(args []interface{}, index int, name string) time.Time {
	val, ok := args[index].(time.Time)
	if !ok {
		panic(fmt.Sprint("function ", name, ": parameter ", index, " not a date ", args[index]))
	}
	return val
}

func MustBeBool(args []interface{}, index int, name string) bool {
	val, ok := args[index].(bool)
	if !ok {
		panic(fmt.Sprint("function ", name, ": parameter ", index, " not a boolean ", args[index]))
	}
	return val
}

func MustBeNumberAsInt(args []interface{}, index int, name string) int {
	f, _ := MustBeNumber(args, index, name).Float64()
	return int(f)
}

func GetNumber(args []interface{}, index int, name string) *big.Rat {
	HasParam(args, index, name)
	switch args[index].(type) {
	case *big.Rat:
		return args[index].(*big.Rat)
	case string:
		r, ok := new(big.Rat).SetString(args[index].(string))
		if ok {
			return r
		}
	}
	panic(fmt.Sprint("function ", name, ": parameter ", index, " not a number ", args[index]))
}

func GetNumberAsInt(args []interface{}, index int, name string) int {
	f, _ := GetNumber(args, index, name).Float64()
	return int(f)
}

func HasParam(args []interface{}, index int, name string) {
	if index >= len(args) {
		panic(fmt.Sprint("function ", name, ": should have at least ", index, " parameters, actual ", len(args)))
	}
}
