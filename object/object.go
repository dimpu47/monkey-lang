package object

// Package object implements the object system (or value system) of Monkey
// used to both represent values as the evaluator encounters and constructs
// them as well as how the user interacts with values.

import (
	"bytes"
	"fmt"
	"strings"

	"git.mills.io/prologic/monkey/ast"
)

const (
	// INTEGER is the Integer object type
	INTEGER = "INTEGER"

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
)

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
