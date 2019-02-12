package object

import (
	"fmt"
	"os"
)

// Assert ...
func Assert(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != BOOLEAN {
		return newError("argument #1 to `assert` must be BOOLEAN, got %s",
			args[0].Type())
	}
	if args[1].Type() != STRING {
		return newError("argument #2 to `assert` must be STRING, got %s",
			args[0].Type())
	}

	if !args[0].(*Boolean).Value {
		fmt.Printf("Assertion Error: %s", args[1].(*String).Value)
		os.Exit(1)
	}

	return nil
}
