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
	// STRING a string, e.g: "1234"
	STRING = "STRING"

	//
	// Operators
	//

	// ASSIGN the assignment operator
	ASSIGN = "="
	// PLUS the addition operator
	PLUS = "+"
	// MINUS the substraction operator
	MINUS = "-"
	// BANG the factorial operator
	BANG = "!"
	// ASTERISK the multiplication operator
	ASTERISK = "*"
	// SLASH the division operator
	SLASH = "/"

	// LT the less than comparision operator
	LT = "<"
	// GT the greater than comparision operator
	GT = ">"

	// EQ the equality operator
	EQ = "=="
	// NEQ the inequality operator
	NEQ = "!="

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
	// LBRACKET a left bracket
	LBRACKET = "["
	// RBRACKET a right bracket
	RBRACKET = "]"

	//
	// Keywords
	//

	// FUNCTION the `fn` keyword (function)
	FUNCTION = "FUNCTION"
	// LET the `let` keyword (let)
	LET = "LET"
	// TRUE the `true` keyword (true)
	TRUE = "TRUE"
	// FALSE the `false` keyword (false)
	FALSE = "FALSE"
	// IF the `if` keyword (if)
	IF = "IF"
	// ELSE the `else` keyword (else)
	ELSE = "ELSE"
	// RETURN the `return` keyword (return)
	RETURN = "RETURN"
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Type represents the type of a token
type Type string

// Token holds a single token type and its literal value
type Token struct {
	Type    Type
	Literal string
}

// LookupIdent looks up the identifier in ident and returns the appropriate
// token type depending on whether the identifier is user-defined or a keyword
func LookupIdent(ident string) Type {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}
