package eval

import (
	"github.com/prologic/monkey-lang/object"
)

var builtins = map[string]*object.Builtin{
	"len":   object.LookupBuiltin("len"),
	"puts":  object.LookupBuiltin("puts"),
	"first": object.LookupBuiltin("first"),
	"last":  object.LookupBuiltin("last"),
	"rest":  object.LookupBuiltin("rest"),
	"push":  object.LookupBuiltin("push"),
}
