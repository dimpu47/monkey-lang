package compiler

import (
	"fmt"
	"testing"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/code"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
)

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

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
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
	tests := []compilerTestCase{
		{
			input: `
            if (true) { 10 }; 3333;
            `,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.LoadTrue),
				// 0001
				code.Make(code.JumpIfFalse, 10),
				// 0004
				code.Make(code.LoadConstant, 0),
				// 0007
				code.Make(code.Jump, 11),
				// 0010
				code.Make(code.LoadNull),
				// 0011
				code.Make(code.Pop),
				// 0012
				code.Make(code.LoadConstant, 1),
				// 0015
				code.Make(code.Pop),
			},
		},
		{
			input: `
            if (true) { 10 } else { 20 }; 3333;
            `,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.LoadTrue),
				// 0001
				code.Make(code.JumpIfFalse, 10),
				// 0004
				code.Make(code.LoadConstant, 0),
				// 0007
				code.Make(code.Jump, 13),
				// 0010
				code.Make(code.LoadConstant, 1),
				// 0013
				code.Make(code.Pop),
				// 0014
				code.Make(code.LoadConstant, 2),
				// 0017
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            let one = 1;
            let two = 2;
            `,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.BindGlobal, 0),
				code.Make(code.LoadConstant, 1),
				code.Make(code.BindGlobal, 1),
			},
		},
		{
			input: `
            let one = 1;
            one;
            `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.BindGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.Pop),
			},
		},
		{
			input: `
            let one = 1;
            let two = one;
            two;
            `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.LoadConstant, 0),
				code.Make(code.BindGlobal, 0),
				code.Make(code.LoadGlobal, 0),
				code.Make(code.BindGlobal, 1),
				code.Make(code.LoadGlobal, 1),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
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
				code.Make(code.Index),
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
				code.Make(code.Index),
				code.Make(code.Pop),
			},
		},
	}

	runCompilerTests(t, tests)
}
