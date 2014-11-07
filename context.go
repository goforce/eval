package eval

import (
	"strings"
	"time"
)

type Context interface {
	AddFunctions(Functions) Context
	AddValues(Values) Context
	SetTimeZone(*time.Location) Context
	ParseDate(format, value string) (time.Time, error)
	cast() *context
}

type context struct {
	functions     []Functions
	values        []Values
	localTimeZone *time.Location
}

func NewContext() *context {
	return &context{functions: make([]Functions, 0, 5), values: make([]Values, 0, 5), localTimeZone: time.Now().Location()}
}

func (context *context) cast() *context {
	return context
}

func (context *context) AddFunctions(functions Functions) Context {
	if context == nil {
		context = NewContext()
	}
	if functions != nil {
		context.functions = append(context.functions, functions)
	}
	return context
}

func (context *context) AddValues(values Values) Context {
	if context == nil {
		context = NewContext()
	}
	if values != nil {
		context.values = append(context.values, values)
	}
	return context
}

func (context *context) SetTimeZone(location *time.Location) Context {
	if context == nil {
		context = NewContext()
	}
	if location != nil {
		context.localTimeZone = location
	}
	return context
}

func (context *context) ParseDate(format, value string) (time.Time, error) {
	layout := format
	for _, f := range mapping {
		layout = strings.Replace(layout, f.human, f.golang, 1)
	}
	if strings.Contains(layout, "Z") {
		return time.Parse(layout, value)
	} else {
		return time.ParseInLocation(layout, value, context.localTimeZone)
	}
}

type f struct {
	human  string
	golang string
}

var mapping = []f{
	f{"YYYY", "2006"},
	f{"YY", "06"},
	f{"MMMM", "January"},
	f{"MMM", "Jan"},
	f{"MM", "01"},
	f{"M", "1"},
	f{"DDDD", "Monday"},
	f{"DDD", "Mon"},
	f{"DD", "02"},
	f{"D", "2"},
	f{"hh", "15"},
	f{"h", "3"},
	f{"mm", "04"},
	f{"m", "4"},
	f{"ss", "05"},
	f{"s", "5"},
	f{"a", "PM"},
	f{"Z", "Z0700"},
}
