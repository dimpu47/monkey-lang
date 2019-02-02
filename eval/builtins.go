package eval

import (
	"github.com/prologic/monkey-lang/object"
)

var builtins = map[string]*object.Builtin{
	"len":   object.LookupBuiltin("len"),
	"input": object.LookupBuiltin("input"),
	"print": object.LookupBuiltin("print"),
	"first": object.LookupBuiltin("first"),
	"last":  object.LookupBuiltin("last"),
	"rest":  object.LookupBuiltin("rest"),
	"push":  object.LookupBuiltin("push"),
	"pop":   object.LookupBuiltin("pop"),
	"exit":  object.LookupBuiltin("exit"),
}
