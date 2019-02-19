package object

// Push ...
func Push(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `push` must be array, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)

	newElements := make([]Object, length+1, length+1)
	copy(newElements, arr.Elements)
	if immutable, ok := args[1].(Immutable); ok {
		newElements[length] = immutable.Clone()
	} else {
		newElements[length] = args[1]
	}

	return &Array{Elements: newElements}
}
