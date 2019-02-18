package object

import (
	"strings"
)

// Upper ...
func Upper(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if str, ok := args[0].(*String); ok {
		return &String{Value: strings.ToUpper(str.Value)}
	}
	return newError("expected `str` argument to `upper` got=%T", args[0])
}
