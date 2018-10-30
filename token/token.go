package token

// Package token implements types and constants to support tokenizing
// the input source before passing the stream of tokens on to the parser.

const (
	// ILLEGAL represents an illegal token
	ILLEGAL = "ILLEGAL"

	// EOF end of file
	EOF = "EOF"

	//
	// Identifiers + literals
	//

	// IDENT an identifier, e.g: add, foobar, x, y, ...
	IDENT = "IDENT"
	// INT an integer, e.g: 1234
	INT = "INT"

	//
	// Operators
	//

	// ASSIGN the assignment operator
	ASSIGN = "="
	// PLUS the addition operator
	PLUS = "+"

	//
	// Delimiters
	//

	// COMMA a comma
	COMMA = ","
	// SEMICOLON a semi-colon
	SEMICOLON = ";"

	// LPAREN a left paranthesis
	LPAREN = "("
	// RPAREN a right parenthesis
	RPAREN = ")"
	// LBRACE a left brace
	LBRACE = "{"
	// RBRACE a right brace
	RBRACE = "}"

	//
	// Keywords
	//

	// FUNCTION the `fn` keyword (function)
	FUNCTION = "FUNCTION"
	// LET the `let` keyword (let)
	LET = "LET"
)

// Type represents the type of a token
type Type string

// Token holds a single token type and its literal value
type Token struct {
	Type    Type
	Literal string
}
