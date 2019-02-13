package object

// Bool ...
func Bool(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *Null:
		return &Boolean{Value: false}
	case *Boolean:
		return arg
	case *Integer:
		if arg.Value == 0 {
			return &Boolean{Value: false}
		}
		return &Boolean{Value: true}
	case *String:
		if len(arg.Value) > 0 {
			return &Boolean{Value: true}
		}
		return &Boolean{Value: false}
	case *Array:
		if len(arg.Elements) > 0 {
			return &Boolean{Value: true}
		}
		return &Boolean{Value: false}
	case *Hash:
		if len(arg.Pairs) > 0 {
			return &Boolean{Value: true}
		}
		return &Boolean{Value: false}

	default:
		return newError("argument to `bool` not supported, got %s",
			args[0].Type())
	}
}
