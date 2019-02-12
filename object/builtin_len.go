package object

import (
	"unicode/utf8"
)

// Len ...
func Len(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	case *String:
		return &Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}
