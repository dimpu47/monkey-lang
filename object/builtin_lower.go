package object

import (
	"strings"
)

// Lower ...
func Lower(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if str, ok := args[0].(*String); ok {
		return &String{Value: strings.ToLower(str.Value)}
	}
	return newError("expected `str` argument to `lower` got=%T", args[0])
}
