package object

// Exit ...
func Exit(args ...Object) Object {
	var status int
	if len(args) == 1 {
		if args[0].Type() != INTEGER {
			return newError("argument to `exit` must be INTEGER, got %s",
				args[0].Type())
		}
		status = int(args[0].(*Integer).Value)
	}

	ExitFunction(status)

	return nil
}
