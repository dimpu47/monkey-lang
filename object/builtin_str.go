package object

import (
	"fmt"
)

// Str ...
func Str(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	arg, ok := args[0].(fmt.Stringer)
	if !ok {
		return newError("argument to `str` not supported, got %s",
			args[0].Type())
	}

	return &String{Value: arg.String()}
}
