package object

// Rest ...
func Rest(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `rest` must be array, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &Array{Elements: newElements}
	}

	return nil
}
