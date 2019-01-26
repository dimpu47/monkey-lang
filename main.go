package main

// Package main implements the main process which executes a program if
// a filename is supplied as an argument or invokes the interpreter's
// REPL and waits for user input before lexing, parsing nad evaulating.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/prologic/monkey-lang/repl"
)

var (
	interactive bool
	version     bool
	debug       bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug mode")

	flag.BoolVar(&interactive, "i", false, "enable interactive mode")
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

	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatalf("could not open source file %s: %s", flag.Arg(0), err)
		}
		env := repl.Exec(f)
		if interactive {
			repl.Start(os.Stdin, os.Stdout, env)
		}
	} else {
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout, nil)
	}
}
