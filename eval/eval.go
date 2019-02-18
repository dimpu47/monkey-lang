package eval

// Package eval implements the evaluator -- a tree-walker implemtnation that
// recursively walks the parsed AST (abstract syntax tree) and evaluates
// the nodes according to their semantic meaning

import (
	"fmt"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/object"
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Eval evaluates the node and returns an object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.Return{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return fromNativeBoolean(node.Value)
	case *ast.Null:
		return NULL

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.BindExpression:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			if immutable, ok := value.(object.Immutable); ok {
				env.Set(ident.Value, immutable.Clone())
			} else {
				env.Set(ident.Value, value)
			}

			return NULL
		}
		return newError("expected identifier on left got=%T", node.Left)

	case *ast.AssignmentExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			env.Set(ident.Value, value)
		} else if ie, ok := node.Left.(*ast.IndexExpression); ok {
			obj := Eval(ie.Left, env)
			if isError(obj) {
				return obj
			}

			if array, ok := obj.(*object.Array); ok {
				index := Eval(ie.Index, env)
				if isError(index) {
					return index
				}
				if idx, ok := index.(*object.Integer); ok {
					array.Elements[idx.Value] = value
				} else {
					return newError("cannot index array with %#v", index)
				}
			} else if hash, ok := obj.(*object.Hash); ok {
				key := Eval(ie.Index, env)
				if isError(key) {
					return key
				}
				if hashKey, ok := key.(object.Hashable); ok {
					hashed := hashKey.HashKey()
					hash.Pairs[hashed] = object.HashPair{Key: key, Value: value}
				} else {
					return newError("cannot index hash with %T", key)
				}
			} else {
				return newError("object type %T does not support item assignment", obj)
			}
		} else {
			return newError("expected identifier or index expression got=%T", left)
		}

		return NULL

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.Return); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		if right.Type() == object.BOOLEAN {
			return evalBooleanPrefixOperatorExpression(operator, right)
		}
		return evalIntegerPrefixOperatorExpression(operator, right)
	case "~", "-":
		return evalIntegerPrefixOperatorExpression(operator, right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBooleanPrefixOperatorExpression(operator string, right object.Object) object.Object {
	if right.Type() != object.BOOLEAN {
		return newError("unknown operator: %s%s", operator, right.Type())
	}

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

func evalIntegerPrefixOperatorExpression(operator string, right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: %s%s", operator, right.Type())
	}

	value := right.(*object.Integer).Value
	switch operator {
	case "!":
		return FALSE
	case "~":
		return &object.Integer{Value: ^value}
	case "-":
		return &object.Integer{Value: -value}
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return fromNativeBoolean(left == right)
	case operator == "!=":
		return fromNativeBoolean(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return fromNativeBoolean(left == right)
	case "!=":
		return fromNativeBoolean(left != right)
	case "&&":
		return &object.Boolean{Value: leftVal && rightVal}
	case "||":
		return &object.Boolean{Value: leftVal || rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case "<=":
		return fromNativeBoolean(leftVal <= rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case ">=":
		return fromNativeBoolean(leftVal >= rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "<":
		return fromNativeBoolean(leftVal < rightVal)
	case ">":
		return fromNativeBoolean(leftVal > rightVal)
	case "==":
		return fromNativeBoolean(leftVal == rightVal)
	case "!=":
		return fromNativeBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			result = Eval(we.Consequence, env)
		} else {
			break
		}
	}

	if result != nil {
		return result
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		env := extendFunctionEnv(fn, args)
		return unwrapReturnValue(Eval(fn.Body, env))

	case *object.Builtin:
		if result := fn.Fn(args...); result != nil {
			return result
		}
		return NULL

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := fn.Env.Clone()

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexAssignmentExpression(left, index, value object.Object) object.Object {
	switch {
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.STRING && index.Type() == object.INTEGER:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return &object.String{Value: ""}
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}
