package object

import (
	"io/ioutil"
)

// Write ...
func Write(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}

	arg, ok := args[0].(*String)
	if !ok {
		return newError("argument #1 to `write` expected to be `str` got=%T", args[0].Type())
	}
	filename := arg.Value

	arg, ok = args[1].(*String)
	if !ok {
		return newError("argument #2 to `write` expected to be `str` got=%T", args[1].Type())
	}
	data := []byte(arg.Value)

	err := ioutil.WriteFile(filename, data, 0755)
	if err != nil {
		return newError("error writing file: %s", err)
	}

	return &Null{}
}
