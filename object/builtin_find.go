package object

import (
	"strings"
)

// Find ...
func Find(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	if haystack, ok := args[0].(*String); ok {
		if needle, ok := args[1].(*String); ok {
			index := strings.Index(haystack.Value, needle.Value)
			return &Integer{Value: int64(index)}
		} else {
			return newError("expected arg #2 to be `str` got got=%T", args[1])
		}
	} else if haystack, ok := args[0].(*Array); ok {
		needle := args[1]
		index := -1
		for i, el := range haystack.Elements {
			if cmp, ok := el.(Comparable); ok && cmp.Equal(needle) {
				index = i
				break
			}
		}
		return &Integer{Value: int64(index)}
	} else {
		return newError("expected arg #1 to be `str` or `array` got got=%T", args[0])
	}
}
