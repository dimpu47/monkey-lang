package main

// Package main implements the main process which executes a program if
// a filename is supplied as an argument or invokes the interpreter's
// REPL and waits for user input before lexing, parsing nad evaulating.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/prologic/monkey-lang/compiler"
	"github.com/prologic/monkey-lang/lexer"
	"github.com/prologic/monkey-lang/parser"
	"github.com/prologic/monkey-lang/repl"
)

var (
	engine      string
	interactive bool
	compile     bool
	version     bool
	debug       bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.BoolVar(&compile, "c", false, "compile input to bytecode")

	flag.BoolVar(&interactive, "i", false, "enable interactive mode")
	flag.StringVar(&engine, "e", "vm", "engine to use (eval or vm)")
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("%s %s", path.Base(os.Args[0]), FullVersion())
		os.Exit(0)
	}

	user, err := user.Current()
	if err != nil {
		log.Fatalf("could not determine current user: %s", err)
	}

	args := flag.Args()

	if compile {
		if len(args) < 1 {
			log.Fatal("no source file given to compile")
		}
		f, err := os.Open(args[0])
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		l := lexer.New(string(b))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			log.Fatal(p.Errors())
		}

		c := compiler.New()
		err = c.Compile(program)
		if err != nil {
			log.Fatal(err)
		}

		code := c.Bytecode()
		fmt.Printf("%s\n", code.Instructions)
	} else {
		opts := &repl.Options{
			Debug:       debug,
			Engine:      engine,
			Interactive: interactive,
		}
		repl := repl.New(user.Username, args, opts)
		repl.Run()
	}
}
