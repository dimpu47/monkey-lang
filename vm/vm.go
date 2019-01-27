package vm

// package vm implements our virtual machine which is passed a sequence of
// bytecode instructions to execute which has been parsed and compiled by
// the lexer/parser and compiler in previous steps

import (
	"fmt"

	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/compiler"
	"github.com/prologic/monkey-lang/object"
)

const StackSize = 2048

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER && rightType == object.INTEGER {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s",
		leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.Add:
		result = leftValue + rightValue
	case code.Sub:
		result = leftValue - rightValue
	case code.Mul:
		result = leftValue * rightValue
	case code.Div:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER || right.Type() == object.INTEGER {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.Equal:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.NotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.Equal:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.NotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.GreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {

		case code.LoadConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.LoadTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.LoadFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.Add, code.Sub, code.Mul, code.Div:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.Equal, code.NotEqual, code.GreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.Bang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		case code.Minus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.Pop:
			vm.pop()

		}
	}

	return nil
}
