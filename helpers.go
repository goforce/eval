// Helper functions to facilitate creation of custom Functions. All functions do panic on error.

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

func NumOfParams(args []interface{}, expected int) {
	if expected != len(args) {
		panic(fmt.Sprint("expected ", expected, " parameters, actual ", len(args)))
	}
}

func MustBeString(args []interface{}, index int) string {
	val, ok := args[index].(string)
	if !ok {
		panic(fmt.Sprint("parameter ", index, " not a string ", args[index]))
	}
	return val
}

func MustBeNumber(args []interface{}, index int) *big.Rat {
	val, ok := args[index].(*big.Rat)
	if !ok {
		panic(fmt.Sprint("parameter ", index, " not a number ", args[index]))
	}
	return val
}

func MustBeDate(args []interface{}, index int) time.Time {
	val, ok := args[index].(time.Time)
	if !ok {
		panic(fmt.Sprint("parameter ", index, " not a date ", args[index]))
	}
	return val
}

func MustBeBool(args []interface{}, index int) bool {
	val, ok := args[index].(bool)
	if !ok {
		panic(fmt.Sprint("parameter ", index, " not a boolean ", args[index]))
	}
	return val
}

func MustBeNumberAsInt(args []interface{}, index int) int {
	f, _ := MustBeNumber(args, index).Float64()
	return int(f)
}

func GetNumber(args []interface{}, index int) *big.Rat {
	switch args[index].(type) {
	case *big.Rat:
		return args[index].(*big.Rat)
	case string:
		r, ok := new(big.Rat).SetString(args[index].(string))
		if ok {
			return r
		}
	}
	panic(fmt.Sprint("parameter ", index, " not a number ", args[index]))
}

func GetNumberAsInt(args []interface{}, index int) int {
	f, _ := GetNumber(args, index).Float64()
	return int(f)
}

func MinNumOfParams(args []interface{}, expected int) {
	if expected > len(args) {
		panic(fmt.Sprint("should have at least ", expected, " parameters, actual ", len(args)))
	}
}
