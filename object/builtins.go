package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"unicode/utf8"
)

var Builtins = map[string]*Builtin{
	"len":   &Builtin{Name: "len", Fn: Len},
	"input": &Builtin{Name: "input", Fn: Input},
	"print": &Builtin{Name: "print", Fn: Print},
	"first": &Builtin{Name: "first", Fn: First},
	"last":  &Builtin{Name: "last", Fn: Last},
	"rest":  &Builtin{Name: "rest", Fn: Rest},
	"push":  &Builtin{Name: "push", Fn: Push},
	"pop":   &Builtin{Name: "pop", Fn: Pop},
	"exit":  &Builtin{Name: "exit", Fn: Exit},
}

var BuiltinsIndex []*Builtin

func init() {
	var keys []string
	for k := range Builtins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		BuiltinsIndex = append(BuiltinsIndex, Builtins[k])
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func Len(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	case *String:
		return &Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s",
			args[0].Type())
	}
}

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

func Print(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return nil
}

func First(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `first` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return nil
}

func Last(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `last` must be ARRAY, got %s",
			args[0].Type())
	}

	arr := args[0].(*Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}

func Rest(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `rest` must be ARRAY, got %s",
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

func Push(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `push` must be ARRAY, got %s",
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

func Pop(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != ARRAY {
		return newError("argument to `pop` must be ARRAY, got %s",
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

func Exit(args ...Object) Object {
	if len(args) == 1 {
		if args[0].Type() != INTEGER {
			return newError("argument to `exit` must be INTEGER, got %s",
				args[0].Type())
		}
		os.Exit(int(args[0].(*Integer).Value))
	} else {
		os.Exit(0)
	}
	return nil
}
