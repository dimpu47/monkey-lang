package repl

// Package repl implements the Read-Eval-Print-Loop or interactive console
// by lexing, parsing and evaluating the input in the interpreter

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/prologic/monkey-lang/eval"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/object"
	"github.com/prologic/monkey-lang/parser"
)

// PROMPT is the REPL prompt displayed for each input
const PROMPT = ">> "

// MonkeyFace is the REPL's face of shock and horror when you encounter a
// parser error :D
const MonkeyFace = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

// Exec parses and executes the program given by f and returns the resulting
// environment, any errors are printed to stderr
func Exec(f io.Reader) (env *object.Environment) {
	env = object.NewEnvironment()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %s", err)
		return
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	eval.Eval(program, env)
	return
}

// Start starts the REPL in a continious loop
func Start(in io.Reader, out io.Writer, env *object.Environment) {
	scanner := bufio.NewScanner(in)

	if env == nil {
		env = object.NewEnvironment()
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		obj := eval.Eval(program, env)
		if obj != nil {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MonkeyFace)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
