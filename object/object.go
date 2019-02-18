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
	INTEGER = "int"

	// STRING is the String object type
	STRING = "str"

	// BOOLEAN is the Boolean object type
	BOOLEAN = "bool"

	// NULL is the Null object type
	NULL = "null"

	// RETURN is the Return object type
	RETURN = "return"

	// ERROR is the Error object type
	ERROR = "error"

	// FUNCTION is the Function object type
	FUNCTION = "fn"

	// COMPILED_FUNCTION is the CompiledFunction object type
	COMPILED_FUNCTION = "COMPILED_FUNCTION"

	// BUILTIN is the Builtin object type
	BUILTIN = "builtin"

	// CLOSURE is the Closure object type
	CLOSURE = "closure"

	// ARRAY is the Array object type
	ARRAY = "array"

	// HASH is the Hash object type
	HASH = "hash"
)

// Comparable is the interface for comparing two Object and their underlying
// values. It is the responsibility of the caller (left) to check for types.
// Returns `true` iif the types and values are identical, `false` otherwise.
type Comparable interface {
	Equal(other Object) bool
}

// Immutable is the interface for all immutable objects which must implement
// the Clone() method used by binding names to values.
type Immutable interface {
	Clone() Object
}

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
	String() string
	Inspect() string
}

// Integer is the integer type used to represent integer literals and holds
// an internal int64 value
type Integer struct {
	Value int64
}

func (i *Integer) Equal(other Object) bool {
	if obj, ok := other.(*Integer); ok {
		return i.Value == obj.Value
	}
	return false
}

func (i *Integer) String() string {
	return i.Inspect()
}

// Clone creates a new copy
func (i *Integer) Clone() Object {
	return &Integer{Value: i.Value}
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

func (s *String) Equal(other Object) bool {
	if obj, ok := other.(*String); ok {
		return s.Value == obj.Value
	}
	return false
}

func (s *String) String() string {
	return s.Value
}

// Clone creates a new copy
func (s *String) Clone() Object {
	return &String{Value: s.Value}
}

// Type returns the type of the object
func (s *String) Type() Type { return STRING }

// Inspect returns a stringified version of the object for debugging
func (s *String) Inspect() string { return fmt.Sprintf("%#v", s.Value) }

// Boolean is the boolean type and used to represent boolean literals and
// holds an interval bool value
type Boolean struct {
	Value bool
}

func (b *Boolean) Equal(other Object) bool {
	if obj, ok := other.(*Boolean); ok {
		return b.Value == obj.Value
	}
	return false
}

func (b *Boolean) String() string {
	return b.Inspect()
}

// Clone creates a new copy
func (b *Boolean) Clone() Object {
	return &Boolean{Value: b.Value}
}

// Type returns the type of the object
func (b *Boolean) Type() Type { return BOOLEAN }

// Inspect returns a stringified version of the object for debugging
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null is the null type and used to represent the absence of a value
type Null struct{}

func (n *Null) Equal(other Object) bool {
	_, ok := other.(*Null)
	return ok
}

func (n *Null) String() string {
	return n.Inspect()
}

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

func (rv *Return) String() string {
	return rv.Inspect()
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

func (e *Error) String() string {
	return e.Message
}

// Clone creates a new copy
func (e *Error) Clone() Object {
	return &Error{Message: e.Message}
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

func (cf *CompiledFunction) String() string {
	return cf.Inspect()
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

func (f *Function) String() string {
	return f.Inspect()
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
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) String() string {
	return b.Inspect()
}

// Type returns the type of the object
func (b *Builtin) Type() Type { return BUILTIN }

// Inspect returns a stringified version of the object for debugging
func (b *Builtin) Inspect() string {
	return fmt.Sprintf("<built-in function %s>", b.Name)
}

// Closure is the closure object type that holds a reference to a compiled
// functions and its free variables
type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) String() string {
	return c.Inspect()
}

// Type returns the type of the object
func (c *Closure) Type() Type { return CLOSURE }

// Inspect returns a stringified version of the object for debugging
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

// Array is the array literal type that holds a slice of Object(s)
type Array struct {
	Elements []Object
}

func (ao *Array) Equal(other Object) bool {
	if obj, ok := other.(*Array); ok {
		if len(ao.Elements) != len(obj.Elements) {
			return false
		}
		for i, el := range ao.Elements {
			cmp, ok := el.(Comparable)
			if !ok {
				return false
			}
			if !cmp.Equal(obj.Elements[i]) {
				return false
			}
		}

		return true
	}
	return false
}

func (ao *Array) String() string {
	return ao.Inspect()
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

func (h *Hash) Equal(other Object) bool {
	if obj, ok := other.(*Hash); ok {
		if len(h.Pairs) != len(obj.Pairs) {
			return false
		}
		for _, pair := range h.Pairs {
			left := pair.Value
			hashed := left.(Hashable)
			right, ok := obj.Pairs[hashed.HashKey()]
			if !ok {
				return false
			}
			cmp, ok := left.(Comparable)
			if !ok {
				return false
			}
			if !cmp.Equal(right.Value) {
				return false
			}
		}

		return true
	}
	return false
}

func (h *Hash) String() string {
	return h.Inspect()
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
