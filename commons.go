package eval

import ()

type Values func(string) (interface{}, bool)
type Functions func(string, []interface{}) (interface{}, error)
