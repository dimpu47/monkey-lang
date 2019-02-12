package object

import (
	"os"
)

// Exit ...
func Exit(args ...Object) Object {
	if len(args) == 1 {
		if args[0].Type() != INTEGER {
			return newError("argument to `exit` must be INTEGER, got %s",
				args[0].Type())
		}
		os.Exit(int(args[0].(*Integer).Value))
	} else {
		os.Exit(0)
	}
	return nil
}
