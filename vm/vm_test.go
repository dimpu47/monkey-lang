package vm

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prologic/monkey-lang/ast"
	"github.com/prologic/monkey-lang/compiler"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v) want=%d",
			actual, actual, expected)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Log(tt.input)
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Log(tt.input)
			t.Fatalf("vm error: %s", err)
		}
		if vm.sp != 0 {
			t.Log(tt.input)
			t.Fatal("vm stack pointer non-zero")
		}

		stackElem := vm.LastPopped()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {

	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}

		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash has wrong number of Pairs. want=%d, got=%d",
				len(expected), len(hash.Pairs))
			return
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}

			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d",
				len(expected), len(array.Elements))
			return
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}

	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}

	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}

	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)
			return
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q",
				expected.Message, errObj.Message)
		}

	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
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

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}

	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"null", nil},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(1 <= 2) == true", true},
		{"(1 <= 2) == false", false},
		{"(1 >= 2) == true", false},
		{"(1 >= 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
		{`"a" == "a"`, true},
		{`"a" < "b"`, true},
		{`"abc" == "abc"`, true},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
		{"if (true) { let a = 5; }", Null},
		{"if (true) { 10; let a = 5; }", Null},
		{"if (false) { 10 } else { let b = 5; }", Null},
		{"if (false) { 10 } else { 10; let b = 5; }", Null},
		{"if (true) { let a = 5; } else { 10 }", Null},
		{"let x = 0; if (true) { x = 1; }; if (false) { x = 2; }; x", 1},
		{"if (1 < 2) { 10 } else if (1 == 2) { 20 }", 10},
		{"if (1 > 2) { 10 } else if (1 == 2) { 20 } else { 30 }", 30},
	}

	runVmTests(t, tests)
}

func TestIterations(t *testing.T) {
	tests := []vmTestCase{
		{"while (false) { }", nil},
		{"let n = 0; while (n < 10) { n = n + 1 }; n", 10},
		{"let n = 10; while (n > 0) { n = n - 1 }; n", 0},
		{"let n = 0; while (n < 10) { n = n + 1 }", nil},
		{"let n = 10; while (n > 0) { n = n - 1 }", nil},
	}

	runVmTests(t, tests)
}

func TestAssignmentStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 0; one = 1", 1},
		{"let one = 0; one = 1; one", 1},
		{"let one = 0; one = 1; let two = 0; two = 2; one + two", 3},
		{"let one = 0; one = 1; let two = 0; two = one + one; one + two", 3},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
		{`"abc"[0]`, "a"},
		{`"abc"[1]`, "b"},
		{`"abc"[2]`, "c"},
		{`"abc"[3]`, ""},
		{`"abc"[-1]`, ""},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let fivePlusTen = fn() { return 5 + 10; };
        fivePlusTen();
        `,
			expected: 15,
		},
		{
			input: `
        let one = fn() { return 1; };
        let two = fn() { return 2; };
        one() + two()
        `,
			expected: 3,
		},
		{
			input: `
        let a = fn() { return 1 };
        let b = fn() { return a() + 1 };
        let c = fn() { return b() + 1 };
        c();
        `,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let earlyExit = fn() { return 99; 100; };
        earlyExit();
        `,
			expected: 99,
		},
		{
			input: `
        let earlyExit = fn() { return 99; return 100; };
        earlyExit();
        `,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let noReturn = fn() { };
        noReturn();
        `,
			expected: Null,
		},
		{
			input: `
        let noReturn = fn() { };
        let noReturnTwo = fn() { noReturn(); };
        noReturn();
        noReturnTwo();
        `,
			expected: Null,
		},
	}

	runVmTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let returnsOne = fn() { return 1; };
        let returnsOneReturner = fn() { return returnsOne; };
        returnsOneReturner()();
        `,
			expected: 1,
		},
		{
			input: `
        let returnsOneReturner = fn() {
            let returnsOne = fn() { return 1; };
            return returnsOne;
        };
        returnsOneReturner()();
        `,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; return one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; return one + two; };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; return one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; return three + four; };
			oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; return foobar; };
			let secondFoobar = fn() { let foobar = 100; return foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				return globalSeed - num;
			}
			let minusTwo = fn() {
				let num = 2;
				return globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let identity = fn(a) { return a; };
        identity(4);
        `,
			expected: 4,
		},
		{
			input: `
        let sum = fn(a, b) { return a + b; };
        sum(1, 2);
        `,
			expected: 3,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            return c;
        };
        sum(1, 2);
        `,
			expected: 3,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            return c;
        };
        sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            return c;
        };
        let outer = fn() {
            return sum(1, 2) + sum(3, 4);
        };
        outer();
        `,
			expected: 10,
		},
		{
			input: `
        let globalNum = 10;

        let sum = fn(a, b) {
            let c = a + b;
            return c + globalNum;
        };

        let outer = fn() {
            return sum(1, 2) + sum(3, 4) + globalNum;
        };
        outer() + globalNum;
        `,
			expected: 50,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() { return 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { return a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { return a + b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "wrong number of arguments. got=2, want=1",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`print("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{`first(1)`,
			&object.Error{
				Message: "argument to `first` must be ARRAY, got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{`last(1)`,
			&object.Error{
				Message: "argument to `last` must be ARRAY, got INTEGER",
			},
		},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`,
			&object.Error{
				Message: "argument to `push` must be ARRAY, got INTEGER",
			},
		},
		{`input()`, ""},
		{`pop([])`, &object.Error{
			Message: "cannot pop from an empty array",
		},
		},
		{`pop([1])`, 1},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let newClosure = fn(a) {
			return fn() { return a; };
		};
		let closure = newClosure(99);
		closure();
		`,
			expected: 99,
		},
		{
			input: `
		let newAdder = fn(a, b) {
			return fn(c) { return a + b + c };
		};
		let adder = newAdder(1, 2);
		adder(8);
		`,
			expected: 11,
		},
		{
			input: `
		let newAdder = fn(a, b) {
			let c = a + b;
			return fn(d) { return c + d };
		};
		let adder = newAdder(1, 2);
		adder(8);
		`,
			expected: 11,
		},
		{
			input: `
		let newAdderOuter = fn(a, b) {
			let c = a + b;
			return fn(d) {
				let e = d + c;
				return fn(f) { return e + f; };
			};
		};
		let newAdderInner = newAdderOuter(1, 2)
		let adder = newAdderInner(3);
		adder(8);
		`,
			expected: 14,
		},
		{
			input: `
		let a = 1;
		let newAdderOuter = fn(b) {
			return fn(c) {
				return fn(d) { return a + b + c + d };
			};
		};
		let newAdderInner = newAdderOuter(2)
		let adder = newAdderInner(3);
		adder(8);
		`,
			expected: 14,
		},
		{
			input: `
		let newClosure = fn(a, b) {
			let one = fn() { return a; };
			let two = fn() { return b; };
			return fn() { return one() + two(); };
		};
		let closure = newClosure(9, 90);
		closure();
		`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let fibonacci = fn(x) {
            if (x == 0) {
                return 0;
            } else {
                if (x == 1) {
                    return 1;
                } else {
                    return fibonacci(x - 1) + fibonacci(x - 2);
                }
            }
        };
        fibonacci(15);
        `,
			expected: 610,
		},
	}

	runVmTests(t, tests)
}

func TestTailCalls(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fact = fn(n, a) {
			  if (n == 0) {
				return a
			  }
			  return fact(n - 1, a * n)
			}

			fact(5, 1)
        	`,
			expected: 120,
		},

		// without tail recursion optimization this will cause a stack overflow
		{
			input: `
			let iter = fn(n, max) {
				if (n == max) {
					return n
				}
				return iter(n + 1, max)
			}
			iter(0, 9999)
			`,
			expected: 9999,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsInFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let double = fn(x) { return x * 2 };
			double(5);
			`,
			expected: 10,
		},
		{
			input: `
			let double        = fn(x) { return x * 2 };
			let double_double = fn(x) { return 2 * double(x); };
			double_double(5);
			`,
			expected: 20,
		},
		{
			input: `
			let double        = fn(x) { return x * 2 };
			let wrappedDouble = fn(x) { return double(x); };
			wrappedDouble(5);
			`,
			expected: 10,
		},
		{
			input: `
			let wrappedDouble = fn() {
				let double = fn(x) { return x * 2 };
				return double(5);
			};
			wrappedDouble();
			`,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}

func TestCallingRecursiveFunctionsInFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			// This works
			input: `
			let inner = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					return inner(x - 1);
				}
			};
			inner(1);
			`,
			expected: 0,
		},
		{
			// This also works
			input: `
			let inner = fn(x) {
				if (x == 0) {
					return 1;
				} else {
					return inner(x - 1);
				}
			};
			let wrapper = fn() {
				return inner(1);
			};
			wrapper();
			`,
			expected: 1,
		},
		{
			// This does _NOT_ work
			input: `
			let wrapper = fn() {
				let inner = fn(x) {
					if (x == 0) {
						return 2;
					} else {
						return inner(x - 1);
					}
				};
				return inner(1);
			};
			wrapper();
			`,
			expected: 2,
		},
	}

	runVmTests(t, tests)
}

func TestIntegration(t *testing.T) {
	matches, err := filepath.Glob("../testdata/*.monkey")
	if err != nil {
		t.Error(err)
	}

	for _, match := range matches {
		basename := path.Base(match)
		name := strings.TrimSuffix(basename, filepath.Ext(basename))

		t.Run(name, func(t *testing.T) {
			b, err := ioutil.ReadFile(match)
			if err != nil {
				t.Error(err)
			}

			input := string(b)
			program := parse(input)

			c := compiler.New()
			err = c.Compile(program)
			if err != nil {
				t.Log(input)
				t.Fatalf("compiler error: %s", err)
			}

			vm := New(c.Bytecode())

			err = vm.Run()
			if err != nil {
				t.Log(input)
				t.Fatalf("vm error: %s", err)
			}
			if vm.sp != 0 {
				t.Log(input)
				t.Fatal("vm stack pointer non-zero")
			}
		})
	}
}

func BenchmarkFibonacci(b *testing.B) {
	tests := map[string]string{
		"iterative": `
		let fib = fn(n) {
		   if (n < 3) {
			 return 1
		   }
		   let a = 1
		   let b = 1
		   let c = 0
		   let i = 0
		   while (i < n - 2) {
			 c = a + b
			 b = a
			 a = c
			 i = i + 1
		   }
		   return a
		}

		fib(35)
		`,
		"recursive": `
		let fib = fn(x) {
		  if (x == 0) {
			return 0
		  }
		  if (x == 1) {
			return 1
		  }
		  return fib(x-1) + fib(x-2)
		}

		fib(35)
		`,
		"tail-recursive": `
		let fib = fn(n, a, b) {
		  if (n == 0) {
			return a
		  }
		  if (n == 1) {
			return b
		  }
		  return fib(n - 1, b, a + b)
		}

		fib(35, 0, 1)
		`,
	}

	for name, input := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				program := parse(input)

				c := compiler.New()
				err := c.Compile(program)
				if err != nil {
					b.Log(input)
					b.Fatalf("compiler error: %s", err)
				}

				vm := New(c.Bytecode())

				err = vm.Run()
				if err != nil {
					b.Log(input)
					b.Fatalf("vm error: %s", err)
				}
				if vm.sp != 0 {
					b.Log(input)
					b.Fatal("vm stack pointer non-zero")
				}
			}
		})
	}
}
