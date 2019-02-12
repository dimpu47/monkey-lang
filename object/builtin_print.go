package object

import (
	"fmt"
)

// Print ...
func Print(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.String())
	}

	return nil
}
