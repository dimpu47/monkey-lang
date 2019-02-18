package object

import (
	"strings"
)

// Join ...
func Join(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	if arr, ok := args[0].(*Array); ok {
		if sep, ok := args[1].(*String); ok {
			a := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				a[i] = el.String()
			}
			return &String{Value: strings.Join(a, sep.Value)}
		} else {
			return newError("expected arg #2 to be `str` got got=%T", args[1])
		}
	} else {
		return newError("expected arg #1 to be `array` got got=%T", args[0])
	}
}
