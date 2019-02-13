package object

import (
	"strconv"
)

// Int ...
func Int(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *Boolean:
		if arg.Value {
			return &Integer{Value: 1}
		}
		return &Integer{Value: 0}
	case *Integer:
		return arg
	case *String:
		n, err := strconv.ParseInt(arg.Value, 10, 64)
		if err != nil {
			return newError("could not parse string to int: %s", err)
		}
		return &Integer{Value: n}
	default:
		return newError("argument to `int` not supported, got %s",
			args[0].Type())
	}
}
