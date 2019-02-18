package object

// TypeOf ...
func TypeOf(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	return &String{Value: string(args[0].Type())}
}
