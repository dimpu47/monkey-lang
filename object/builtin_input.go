package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Input ...
func Input(args ...Object) Object {
	if len(args) > 0 {
		obj, ok := args[0].(*String)
		if !ok {
			return newError(
				"argument to `input` not supported, got %s",
				args[0].Type(),
			)
		}
		fmt.Fprintf(os.Stdout, obj.Value)
	}

	buffer := bufio.NewReader(os.Stdin)

	line, _, err := buffer.ReadLine()
	if err != nil && err != io.EOF {
		return newError(fmt.Sprintf("error reading input from stdin: %s", err))
	}
	return &String{Value: string(line)}
}
