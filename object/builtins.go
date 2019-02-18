package object

import (
	"fmt"
	"sort"
)

// Builtins ...
var Builtins = map[string]*Builtin{
	"len":    &Builtin{Name: "len", Fn: Len},
	"input":  &Builtin{Name: "input", Fn: Input},
	"print":  &Builtin{Name: "print", Fn: Print},
	"first":  &Builtin{Name: "first", Fn: First},
	"last":   &Builtin{Name: "last", Fn: Last},
	"rest":   &Builtin{Name: "rest", Fn: Rest},
	"push":   &Builtin{Name: "push", Fn: Push},
	"pop":    &Builtin{Name: "pop", Fn: Pop},
	"exit":   &Builtin{Name: "exit", Fn: Exit},
	"assert": &Builtin{Name: "assert", Fn: Assert},
	"bool":   &Builtin{Name: "bool", Fn: Bool},
	"int":    &Builtin{Name: "int", Fn: Int},
	"str":    &Builtin{Name: "str", Fn: Str},
	"typeof": &Builtin{Name: "typeof", Fn: TypeOf},
	"args":   &Builtin{Name: "args", Fn: Args},
}

// BuiltinsIndex ...
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
