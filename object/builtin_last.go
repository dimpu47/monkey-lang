package object

// Last ...
func Last(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `last` must be array, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}
