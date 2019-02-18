package object

// Args ...
func Args(args ...Object) Object {
	elements := make([]Object, len(Arguments))
	for i, arg := range Arguments {
		elements[i] = &String{Value: arg}
	}
	return &Array{Elements: elements}
}
