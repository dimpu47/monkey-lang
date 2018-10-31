package main

// Package main implements the main process which invokes the interpreter's
// REPL and waits for user input before lexing, parsing nad evaulating.

import (
	"fmt"
	"os"
	"os/user"

	"git.mills.io/prologic/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
