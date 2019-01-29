package object

// Package object implements the object system (or value system) of Monkey
// used to both represent values as the evaluator encounters and constructs
// them as well as how the user interacts with values.

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/code"
)

const (
	// INTEGER is the Integer object type
	INTEGER = "INTEGER"

	// STRING is the String object type
	STRING = "STRING"

	// BOOLEAN is the Boolean object type
	BOOLEAN = "BOOLEAN"

	// NULL is the Null object type
	NULL = "NULL"

	// RETURN is the Return object type
	RETURN = "RETURN"

	// ERROR is the Error object type
	ERROR = "ERROR"

	// FUNCTION is the Function object type
	FUNCTION = "FUNCTION"

	// COMPILED_FUNCTION is the CompiledFunction object type
	COMPILED_FUNCTION = "COMPILED_FUNCTION"

	// BUILTIN is the Builtin object type
	BUILTIN = "BUILTIN"

	// ARRAY is the Array object type
	ARRAY = "ARRAY"

	// HASH is the Hash object type
	HASH = "HASH"
)

// Hashable is the interface for all hashable objects which must implement
// the HashKey() method which reutrns a HashKey result.
type Hashable interface {
	HashKey() HashKey
}

// BuiltinFunction represents the builtin function type
type BuiltinFunction func(args ...Object) Object

// Type represents the type of an object
type Type string

// Object represents a value and implementations are expected to implement
// `Type()` and `Inspect()` functions
type Object interface {
	Type() Type
	Inspect() string
}

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value int64
}

// Type returns the type of the object
func (i *Integer) Type() Type { return INTEGER }

// Inspect returns a stringified version of the object for debugging
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// String is the string type used to represent string literals and holds
// an internal string value
type String struct {
	Value string
}

// Type returns the type of the object
func (s *String) Type() Type { return STRING }

// Inspect returns a stringified version of the object for debugging
func (s *String) Inspect() string { return s.Value }

// Boolean is the boolean type and used to represent boolean literals and
// holds an interval bool value
type Boolean struct {
	Value bool
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null is the null type and used to represent the absence of a value
type Null struct{}

// Type returns the type of the object
func (n *Null) Type() Type { return NULL }

// Inspect returns a stringified version of the object for debugging
func (n *Null) Inspect() string { return "null" }

// Return is the return type and used to hold the value of another object.
// This is used for `return` statements and this object is tracked through
// the evalulator and when encountered stops evaluation of the program,
// or body of a function.
type Return struct {
	Value Object
}

// Type returns the type of the object
func (rv *Return) Type() Type { return RETURN }

// Inspect returns a stringified version of the object for debugging
func (rv *Return) Inspect() string { return rv.Value.Inspect() }

// Error is the error type and used to hold a message denoting the details of
// error encountered. This object is trakced through the evaluator and when
// encountered stops evaulation of the program or body of a function.
type Error struct {
	Message string
}

// Type returns the type of the object
func (e *Error) Type() Type { return ERROR }

// Inspect returns a stringified version of the object for debugging
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// CompiledFunction is the compiled function type that holds the function's
// compiled body as bytecode instructions
type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

// Type returns the type of the object
func (cf *CompiledFunction) Type() Type { return COMPILED_FUNCTION }

// Inspect returns a stringified version of the object for debugging
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

// Function is the function type that holds the function's formal parameters,
// body and an environment to support closures.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns the type of the object
func (f *Function) Type() Type { return FUNCTION }

// Inspect returns a stringified version of the object for debugging
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Builtin  is the builtin object type that simply holds a reference to
// a BuiltinFunction type that takes zero or more objects as arguments
// and returns an object.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns a stringified version of the object for debugging
func (b *Builtin) Inspect() string { return "builtin function" }

// Array is the array literal type that holds a slice of Object(s)
type Array struct {
	Elements []Object
}

// Type returns the type of the object
func (ao *Array) Type() Type { return ARRAY }

// Inspect returns a stringified version of the object for debugging
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey represents a hash key object and holds the Type of Object
// hashed and its hash value in Value
type HashKey struct {
	Type  Type
	Value uint64
}

// HashKey returns a HashKey object
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns a HashKey object
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns a HashKey object
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashPair is an object that holds a key and value of type Object
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is a hash map and holds a map of HashKey to HashPair(s)
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the type of the object
func (h *Hash) Type() Type { return HASH }

// Inspect returns a stringified version of the object for debugging
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
