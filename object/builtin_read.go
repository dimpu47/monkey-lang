package object

import (
	"io/ioutil"
)

// Read ...
func Read(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	arg, ok := args[0].(*String)
	if !ok {
		return newError("argument to `read` expected to be `str` got=%T", args[0].Type())
	}

	filename := arg.Value
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return newError("error reading file: %s", err)
	}

	return &String{Value: string(data)}
}
