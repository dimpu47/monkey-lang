package eval

// Package eval implements the evaluator -- a tree-walker implemtnation that
// recursively walks the parsed AST (abstract syntax tree) and evaluates
// the nodes according to their semantic meaning

import (
	"git.mills.io/prologic/monkey/ast"
	"git.mills.io/prologic/monkey/object"
)

var (
	// TRUE is a cached Boolean object holding the `true` value
	TRUE = &object.Boolean{Value: true}

	// FALSE is a cached Boolean object holding the `false` value
	FALSE = &object.Boolean{Value: false}

	// NULL is a cached Null object
	NULL = &object.Null{}
)

func fromNativeBoolean(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Eval evaluates the node and returns an object
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return fromNativeBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}