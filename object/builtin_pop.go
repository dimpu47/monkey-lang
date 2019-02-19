package object

// Pop ...
func Pop(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `pop` must be array, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	if length == 0 {
		return newError("cannot pop from an empty array")
	}

	element := arr.Elements[length-1]
	arr.Elements = arr.Elements[:length-1]

	return element
}
