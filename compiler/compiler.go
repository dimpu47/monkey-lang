package compiler

// package compiler turns an AST into a sequence of bytecode instructions
// suitable for execution by our virtual machine

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Scope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	Debug bool

	l         int
	constants []object.Object

	scopes     []Scope
	scopeIndex int

	symbolTable *SymbolTable
}

func New() *Compiler {
	mainScope := Scope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, builtin := range object.BuiltinsIndex {
		symbolTable.DefineBuiltin(i, builtin.Name)
	}

	return &Compiler{
		constants: []object.Object{},

		scopes:     []Scope{mainScope},
		scopeIndex: 0,

		symbolTable: symbolTable,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	c := New()
	c.symbolTable = symbolTable
	c.constants = constants
	return c
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.LoadGlobal, s.Index)
	case LocalScope:
		c.emit(code.LoadLocal, s.Index)
	case BuiltinScope:
		c.emit(code.LoadBuiltin, s.Index)
	case FreeScope:
		c.emit(code.LoadFree, s.Index)
	}
}

func (c *Compiler) enterScope() {
	scope := Scope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.Return))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.Return
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) Compile(node ast.Node) error {
	if c.Debug {
		log.Printf(
			"%sCompiling %T: %s\n",
			strings.Repeat(" ", c.l), node, node.String(),
		)
	}

	switch node := node.(type) {

	case *ast.Program:
		c.l++
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.BindExpression:
		var symbol Symbol

		if ident, ok := node.Left.(*ast.Identifier); ok {
			symbol, ok = c.symbolTable.Resolve(ident.Value)
			if !ok {
				symbol = c.symbolTable.Define(ident.Value)
			} else {
				// Local shadowing of previously defined "free" variable in a
				// function now begin rehound to a locally scopped variable.
				if symbol.Scope == FreeScope {
					symbol = c.symbolTable.Define(ident.Value)
				}
			}

			c.l++
			err := c.Compile(node.Value)
			c.l--
			if err != nil {
				return err
			}

			if symbol.Scope == GlobalScope {
				c.emit(code.BindGlobal, symbol.Index)
			} else {
				c.emit(code.BindLocal, symbol.Index)
			}
		} else {
			return fmt.Errorf("expected identifier got=%s", node.Left)
		}

	case *ast.AssignmentExpression:
		if ident, ok := node.Left.(*ast.Identifier); ok {
			symbol, ok := c.symbolTable.Resolve(ident.Value)
			if !ok {
				return fmt.Errorf("undefined variable %s", ident.Value)
			}

			c.l++
			err := c.Compile(node.Value)
			c.l--
			if err != nil {
				return err
			}

			if symbol.Scope == GlobalScope {
				c.emit(code.AssignGlobal, symbol.Index)
			} else {
				c.emit(code.AssignLocal, symbol.Index)
			}
		} else if ie, ok := node.Left.(*ast.IndexExpression); ok {
			c.l++
			err := c.Compile(ie.Left)
			c.l--
			if err != nil {
				return err
			}

			c.l++
			err = c.Compile(ie.Index)
			c.l--
			if err != nil {
				return err
			}

			c.l++
			err = c.Compile(node.Value)
			c.l--
			if err != nil {
				return err
			}

			c.emit(code.SetItem)
		} else {
			return fmt.Errorf("expected identifier or index expression got=%s", node.Left)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	case *ast.ExpressionStatement:
		c.l++
		err := c.Compile(node.Expression)
		c.l--
		if err != nil {
			return err
		}

		c.emit(code.Pop)

	case *ast.BlockStatement:
		c.l++
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
		c.l--

		if c.lastInstructionIs(code.Pop) {
			c.removeLastPop()
		} else {
			if !c.lastInstructionIs(code.Return) {
				c.emit(code.LoadNull)
			}
		}

	case *ast.IfExpression:
		c.l++
		err := c.Compile(node.Condition)
		c.l--
		if err != nil {
			return err
		}

		// Emit an `JumpIfFalse` with a bogus value
		jumpIfFalsePos := c.emit(code.JumpIfFalse, 0xFFFF)

		c.l++
		err = c.Compile(node.Consequence)
		c.l--
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.Pop) {
			c.removeLastPop()
		}

		// Emit an `Jump` with a bogus value
		jumpPos := c.emit(code.Jump, 0xFFFF)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpIfFalsePos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.LoadNull)
		} else {
			c.l++
			err := c.Compile(node.Alternative)
			c.l--
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.Pop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.WhileExpression:
		jumpConditionPos := len(c.currentInstructions())

		c.l++
		err := c.Compile(node.Condition)
		c.l--
		if err != nil {
			return err
		}

		// Emit an `JumpIfFalse` with a bogus value
		jumpIfFalsePos := c.emit(code.JumpIfFalse, 0xFFFF)

		c.l++
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		c.l--

		// Pop off the LoadNull(s) from ast.BlockStatement(s)
		c.emit(code.Pop)

		c.emit(code.Jump, jumpConditionPos)

		afterConsequencePos := c.emit(code.LoadNull)
		c.changeOperand(jumpIfFalsePos, afterConsequencePos)

	case *ast.PrefixExpression:
		c.l++
		err := c.Compile(node.Right)
		c.l--
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.Not)
		case "~":
			c.emit(code.BitwiseNOT)
		case "-":
			c.emit(code.Minus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			c.l++
			err := c.Compile(node.Right)
			c.l--
			if err != nil {
				return err
			}

			c.l++
			err = c.Compile(node.Left)
			c.l--
			if err != nil {
				return err
			}
			if node.Operator == "<=" {
				c.emit(code.GreaterThanEqual)
			} else {
				c.emit(code.GreaterThan)
			}
			return nil
		}

		c.l++
		err := c.Compile(node.Left)
		c.l--
		if err != nil {
			return err
		}

		c.l++
		err = c.Compile(node.Right)
		c.l--
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.Add)
		case "-":
			c.emit(code.Sub)
		case "*":
			c.emit(code.Mul)
		case "/":
			c.emit(code.Div)
		case "%":
			c.emit(code.Mod)
		case "|":
			c.emit(code.BitwiseOR)
		case "^":
			c.emit(code.BitwiseXOR)
		case "&":
			c.emit(code.BitwiseAND)
		case "||":
			c.emit(code.Or)
		case "&&":
			c.emit(code.And)
		case ">":
			c.emit(code.GreaterThan)
		case ">=":
			c.emit(code.GreaterThanEqual)
		case "==":
			c.emit(code.Equal)
		case "!=":
			c.emit(code.NotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IndexExpression:
		c.l++
		err := c.Compile(node.Left)
		c.l--
		if err != nil {
			return err
		}

		c.l++
		err = c.Compile(node.Index)
		c.l--
		if err != nil {
			return err
		}

		c.emit(code.GetItem)

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			c.l++
			err := c.Compile(k)
			c.l--
			if err != nil {
				return err
			}
			c.l++
			err = c.Compile(node.Pairs[k])
			c.l--
			if err != nil {
				return err
			}
		}

		c.emit(code.MakeHash, len(node.Pairs)*2)

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			c.l++
			err := c.Compile(el)
			c.l--
			if err != nil {
				return err
			}
		}

		c.emit(code.MakeArray, len(node.Elements))

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.LoadConstant, c.addConstant(str))

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.LoadConstant, c.addConstant(integer))

	case *ast.FunctionLiteral:
		c.enterScope()

		if node.Name != "" {
			symbol, ok := c.symbolTable.Resolve(node.Name)
			if !ok {
				return fmt.Errorf("undefined variable %s", node.Name)
			}

			// Redefine the symbol for the name assign to this closure as a
			// "free" variable so MakeClosure <idx> <nfree> has the correct
			// number of free variables (including self)
			symbol = c.symbolTable.DefineFree(symbol)
			c.emit(code.SetSelf, symbol.Index)
		}

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		c.l++
		err := c.Compile(node.Body)
		c.l--
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.Pop) {
			c.replaceLastPopWithReturn()
		}

		// If the function doesn't end with a return statement add one with a
		// `return null;` and also handle the edge-case of empty functions.
		if !c.lastInstructionIs(code.Return) {
			// empty function body (LoadNull from BlockStatement)
			if !c.lastInstructionIs(code.LoadNull) {
				c.emit(code.LoadNull)
			}
			c.emit(code.Return)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(code.MakeClosure, fnIndex, len(freeSymbols))

	case *ast.CallExpression:
		c.l++
		err := c.Compile(node.Function)
		c.l--
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			c.l++
			err := c.Compile(a)
			c.l--
			if err != nil {
				return err
			}
		}

		c.emit(code.Call, len(node.Arguments))

	case *ast.ReturnStatement:
		c.l++
		err := c.Compile(node.ReturnValue)
		c.l--
		if err != nil {
			return err
		}

		c.emit(code.Return)

	case *ast.Null:
		c.emit(code.LoadNull)

	case *ast.Boolean:
		if node.Value {
			c.emit(code.LoadTrue)
		} else {
			c.emit(code.LoadFalse)
		}
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
