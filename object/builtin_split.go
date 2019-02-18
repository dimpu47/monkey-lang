package object

import (
	"strings"
)

// Split ...
func Split(args ...Object) Object {
	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if obj, ok := args[0].(*String); ok {
		var sep string

		s := obj.Value

		if len(args) == 2 {
			if obj, ok := args[1].(*String); ok {
				sep = obj.Value
			} else {
				return newError("expected arg #2 to be `str` got=%T", args[1])
			}
		}

		tokens := strings.Split(s, sep)
		elements := make([]Object, len(tokens))
		for i, token := range tokens {
			elements[i] = &String{Value: token}
		}
		return &Array{Elements: elements}
	} else {
		return newError("expected arg #1 to be `str` got got=%T", args[0])
	}
}
