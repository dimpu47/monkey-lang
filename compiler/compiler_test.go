package compiler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
)

type Instructions string

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(
	expected []code.Instructions,
	actual code.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot =%q",
			concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot =%q",
				i, concatted, actual)
		}
	}

	return nil
}

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {

		case []code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T",
					i, actual[i])
			}

			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s",
					i, err)
			}

		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s",
					i, err)
			}

		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		}
	}

	return nil
}

func testConstants2(t *testing.T, expected []interface{}, actual []object.Object) {
	assert := assert.New(t)

	assert.Equal(len(expected), len(actual))

	for i, constant := range expected {
		switch constant := constant.(type) {

		case Instructions:
			assert.Equal(
				string(constant),
				actual[i].(*object.CompiledFunction).Instructions.String(),
			)
		case string:
			assert.Equal(constant, actual[i].(*object.String).Value)
		case int:
			assert.Equal(int64(constant), actual[i].(*object.Integer).Value)
		}
	}
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

type compilerTestCase2 struct {
	input        string
	constants    []interface{}
	instructions string
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Log(tt.input)
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Log(tt.input)
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Log(tt.input)
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func runCompilerTests2(t *testing.T, tests []compilerTestCase2) {
	t.Helper()

	assert := assert.New(t)

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		assert.NoError(err)

		bytecode := compiler.Bytecode()
		assert.Equal(tt.instructions, bytecode.Instructions.String())

		testConstants2(t, tt.constants, bytecode.Constants)
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Add),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.Pop),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Sub),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Mul),
				code.Make(code.Pop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Div),
				code.Make(code.Pop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.Minus),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.Pop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadFalse),
				code.Make(code.Pop),
			},
		},
		{
			input:             "null",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadNull),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.GreaterThan),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.GreaterThan),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 >= 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.GreaterThanEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 <= 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.GreaterThanEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Equal),
				code.Make(code.Pop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.NotEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.LoadFalse),
				code.Make(code.Equal),
				code.Make(code.Pop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.LoadFalse),
				code.Make(code.NotEqual),
				code.Make(code.Pop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadTrue),
				code.Make(code.Bang),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
            if (true) { 10 }; 3333;
            `,
			constants:    []interface{}{10, 3333},
			instructions: "0000 LoadTrue\n0001 JumpIfFalse 10\n0004 LoadConstant 0\n0007 Jump 11\n0010 LoadNull\n0011 Pop\n0012 LoadConstant 1\n0015 Pop\n",
		},
		{
			input: `
            if (true) { 10 } else { 20 }; 3333;
            `,
			constants:    []interface{}{10, 20, 3333},
			instructions: "0000 LoadTrue\n0001 JumpIfFalse 10\n0004 LoadConstant 0\n0007 Jump 13\n0010 LoadConstant 1\n0013 Pop\n0014 LoadConstant 2\n0017 Pop\n",
		},
		{
			input: `
			x := 0; if (true) { x = 1; }; if (false) { x = 2; }
            `,
			constants:    []interface{}{0, 1, 2},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadTrue\n0008 JumpIfFalse 20\n0011 LoadConstant 1\n0014 AssignGlobal 0\n0017 Jump 21\n0020 LoadNull\n0021 Pop\n0022 LoadFalse\n0023 JumpIfFalse 35\n0026 LoadConstant 2\n0029 AssignGlobal 0\n0032 Jump 36\n0035 LoadNull\n0036 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestIteration(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			while (true) { 10 };
            `,
			expectedConstants: []interface{}{10},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.LoadTrue),
				// 0001
				code.Make(code.JumpIfFalse, 11),
				// 0004
				code.Make(code.LoadConstant, 0),
				// 0007
				code.Make(code.Pop),
				// 0008
				code.Make(code.Jump, 0),
				// 0011
				code.Make(code.LoadNull),
				// 0012
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalBindExpressions(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
			one := 1;
			two := 2;
            `,
			constants:    []interface{}{1, 2},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadConstant 1\n0010 BindGlobal 1\n0013 Pop\n",
		},
		{
			input: `
			one := 1;
            one;
            `,
			constants:    []interface{}{1},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadGlobal 0\n0010 Pop\n",
		},
		{
			input: `
			one := 1;
			two := one;
            two;
            `,
			constants:    []interface{}{1},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadGlobal 0\n0010 BindGlobal 1\n0013 Pop\n0014 LoadGlobal 1\n0017 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []interface{}{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []interface{}{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Add),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeArray, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.LoadConstant, 2),
				code.Make(code.MakeArray, 3),
				code.Make(code.Pop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.Add),
				code.Make(code.LoadConstant, 2),
				code.Make(code.LoadConstant, 3),
				code.Make(code.Sub),
				code.Make(code.LoadConstant, 4),
				code.Make(code.LoadConstant, 5),
				code.Make(code.Mul),
				code.Make(code.MakeArray, 3),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeHash, 0),
				code.Make(code.Pop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.LoadConstant, 2),
				code.Make(code.LoadConstant, 3),
				code.Make(code.LoadConstant, 4),
				code.Make(code.LoadConstant, 5),
				code.Make(code.MakeHash, 6),
				code.Make(code.Pop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.LoadConstant, 2),
				code.Make(code.Add),
				code.Make(code.LoadConstant, 3),
				code.Make(code.LoadConstant, 4),
				code.Make(code.LoadConstant, 5),
				code.Make(code.Mul),
				code.Make(code.MakeHash, 4),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestSelectorExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `{"foo": 1}.foo`,
			expectedConstants: []interface{}{"foo", 1, "foo"},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.MakeHash, 2),
				code.Make(code.LoadConstant, 2),
				code.Make(code.GetItem),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.LoadConstant, 2),
				code.Make(code.MakeArray, 3),
				code.Make(code.LoadConstant, 3),
				code.Make(code.LoadConstant, 4),
				code.Make(code.Add),
				code.Make(code.GetItem),
				code.Make(code.Pop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []interface{}{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.MakeHash, 2),
				code.Make(code.LoadConstant, 2),
				code.Make(code.LoadConstant, 3),
				code.Make(code.Sub),
				code.Make(code.GetItem),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.LoadConstant, 0),
					code.Make(code.LoadConstant, 1),
					code.Make(code.Add),
					code.Make(code.Return),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeClosure, 2, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `fn() { 1; return 2 }`,
			expectedConstants: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.LoadConstant, 0),
					code.Make(code.Pop),
					code.Make(code.LoadConstant, 1),
					code.Make(code.Return),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeClosure, 2, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionsWithoutReturn(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadNull),
					code.Make(code.Return),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeClosure, 0, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.Mul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.Sub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.Sub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.Sub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d",
			compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}
	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.Add)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.Add {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.Add)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != code.Mul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
			previous.Opcode, code.Mul)
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `fn() { return 24 }();`,
			constants: []interface{}{
				24,
				Instructions("0000 LoadConstant 0\n0003 Return\n"),
			},
			instructions: "0000 MakeClosure 1 0\n0004 Call 0\n0006 Pop\n",
		},
		{
			input: `
			noArg := fn() { return 24 };
            noArg();
            `,
			constants: []interface{}{
				24,
				Instructions("0000 SetSelf 0\n0002 LoadConstant 0\n0005 Return\n"),
			},
			instructions: "0000 LoadGlobal 0\n0003 MakeClosure 1 1\n0007 BindGlobal 0\n0010 Pop\n0011 LoadGlobal 0\n0014 Call 0\n0016 Pop\n",
		},
		{
			input: `
			oneArg := fn(a) { return a };
            oneArg(24);
            `,
			constants: []interface{}{
				Instructions("0000 SetSelf 0\n0002 LoadLocal 0\n0004 Return\n"),
				24,
			},
			instructions: "0000 LoadGlobal 0\n0003 MakeClosure 0 1\n0007 BindGlobal 0\n0010 Pop\n0011 LoadGlobal 0\n0014 LoadConstant 1\n0017 Call 1\n0019 Pop\n",
		},
		{
			input: `
			manyArg := fn(a, b, c) { a; b; return c };
            manyArg(24, 25, 26);
            `,
			constants: []interface{}{
				Instructions("0000 SetSelf 0\n0002 LoadLocal 0\n0004 Pop\n0005 LoadLocal 1\n0007 Pop\n0008 LoadLocal 2\n0010 Return\n"),
				24,
				25,
				26,
			},
			instructions: "0000 LoadGlobal 0\n0003 MakeClosure 0 1\n0007 BindGlobal 0\n0010 Pop\n0011 LoadGlobal 0\n0014 LoadConstant 1\n0017 LoadConstant 2\n0020 LoadConstant 3\n0023 Call 3\n0025 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
			x := 1
			x = 2
			`,
			constants:    []interface{}{1, 2},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadConstant 1\n0010 AssignGlobal 0\n0013 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestAssignmentStatementScopes(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
			num := 0;
            fn() { num  = 55; }
            `,
			constants: []interface{}{
				0,
				55,
				Instructions("0000 LoadConstant 1\n0003 AssignGlobal 0\n0006 LoadNull\n0007 Return\n"),
			},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 MakeClosure 2 0\n0011 Pop\n",
		},
		{
			input: `
			fn() { num := 0; num  = 55; }
            `,
			constants: []interface{}{
				0,
				55,
				Instructions("0000 LoadConstant 0\n0003 BindLocal 0\n0005 Pop\n0006 LoadConstant 1\n0009 AssignLocal 0\n0011 LoadNull\n0012 Return\n"),
			},
			instructions: "0000 MakeClosure 2 0\n0004 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestBindExpressionScopes(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
			num := 55;
            fn() { return num }
            `,
			constants: []interface{}{
				55,
				Instructions("0000 LoadGlobal 0\n0003 Return\n"),
			},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 MakeClosure 1 0\n0011 Pop\n",
		},
		{
			input: `
            fn() {
				num := 55;
                return num
            }
            `,
			constants: []interface{}{
				55,
				Instructions("0000 LoadConstant 0\n0003 BindLocal 0\n0005 Pop\n0006 LoadLocal 0\n0008 Return\n"),
			},
			instructions: "0000 MakeClosure 1 0\n0004 Pop\n",
		},
		{
			input: `
            fn() {
				a := 55;
				b := 77;
                return a + b
            }
            `,
			constants: []interface{}{
				55,
				77,
				Instructions("0000 LoadConstant 0\n0003 BindLocal 0\n0005 Pop\n0006 LoadConstant 1\n0009 BindLocal 1\n0011 Pop\n0012 LoadLocal 0\n0014 LoadLocal 1\n0016 Add\n0017 Return\n"),
			},
			instructions: "0000 MakeClosure 2 0\n0004 Pop\n",
		},
		{
			input: `
			a := 0;
			a := a + 1;
			`,
			constants: []interface{}{
				0,
				1,
			},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 LoadGlobal 0\n0010 LoadConstant 1\n0013 Add\n0014 BindGlobal 0\n0017 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            len([]);
            push([], 1);
            `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadBuiltin, 7),
				code.Make(code.MakeArray, 0),
				code.Make(code.Call, 1),
				code.Make(code.Pop),
				code.Make(code.LoadBuiltin, 10),
				code.Make(code.MakeArray, 0),
				code.Make(code.LoadConstant, 0),
				code.Make(code.Call, 2),
				code.Make(code.Pop),
			},
		},
		{
			input: `fn() { return len([]) }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.LoadBuiltin, 7),
					code.Make(code.MakeArray, 0),
					code.Make(code.Call, 1),
					code.Make(code.Return),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.MakeClosure, 0, 0),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase2{
		{
			input: `
            fn(a) {
                return fn(b) {
                    return a + b
                }
            }
            `,
			constants: []interface{}{
				Instructions("0000 LoadFree 0\n0002 LoadLocal 0\n0004 Add\n0005 Return\n"),
				Instructions("0000 LoadLocal 0\n0002 MakeClosure 0 1\n0006 Return\n"),
			},
			instructions: "0000 MakeClosure 1 0\n0004 Pop\n",
		},
		{
			input: `
            fn(a) {
                return fn(b) {
                    return fn(c) {
                        return a + b + c
                    }
                }
            };
            `,
			constants: []interface{}{
				Instructions("0000 LoadFree 0\n0002 LoadFree 1\n0004 Add\n0005 LoadLocal 0\n0007 Add\n0008 Return\n"),
				Instructions("0000 LoadFree 0\n0002 LoadLocal 0\n0004 MakeClosure 0 2\n0008 Return\n"),
				Instructions("0000 LoadLocal 0\n0002 MakeClosure 1 1\n0006 Return\n"),
			},
			instructions: "0000 MakeClosure 2 0\n0004 Pop\n",
		},
		{
			input: `
			global := 55;

            fn() {
				a := 66;

                return fn() {
					b := 77;

                    return fn() {
						c := 88;

                        return global + a + b + c;
                    }
                }
            }
            `,
			constants: []interface{}{
				55,
				66,
				77,
				88,
				Instructions("0000 LoadConstant 3\n0003 BindLocal 0\n0005 Pop\n0006 LoadGlobal 0\n0009 LoadFree 0\n0011 Add\n0012 LoadFree 1\n0014 Add\n0015 LoadLocal 0\n0017 Add\n0018 Return\n"),
				Instructions("0000 LoadConstant 2\n0003 BindLocal 0\n0005 Pop\n0006 LoadFree 0\n0008 LoadLocal 0\n0010 MakeClosure 4 2\n0014 Return\n"),
				Instructions("0000 LoadConstant 1\n0003 BindLocal 0\n0005 Pop\n0006 LoadLocal 0\n0008 MakeClosure 5 1\n0012 Return\n"),
			},
			instructions: "0000 LoadConstant 0\n0003 BindGlobal 0\n0006 Pop\n0007 MakeClosure 6 0\n0011 Pop\n",
		},
	}

	runCompilerTests2(t, tests)
}
