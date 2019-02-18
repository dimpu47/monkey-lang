package token

// Package token implements types and constants to support tokenizing
// the input source before passing the stream of tokens on to the parser.

const (
	// ILLEGAL represents an illegal token
	ILLEGAL = "ILLEGAL"

	// EOF end of file
	EOF = "EOF"

	// COMMENT a line comment, e.g: # this is a comment
	COMMENT = "COMMENT"

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

	// BIND the bind operator
	BIND = ":="
	// ASSIGN the assignment operator
	ASSIGN = "="
	// PLUS the addition operator
	PLUS = "+"
	// MINUS the substraction operator
	MINUS = "-"
	// MULTIPLY the multiplication operator
	MULTIPLY = "*"
	// DIVIDE the division operator
	DIVIDE = "/"
	// MODULO the modulo operator
	MODULO = "%"

	//
	// Bitwise / Logical operators
	//

	// Bitwise AND
	BitwiseAND = "&"
	// Bitwise OR
	BitwiseOR = "|"
	// Bitwise XOR
	BitwiseXOR = "^"
	// Bitwise NOT
	BitwiseNOT = "~"

	//
	// Logical operators
	//

	// NOT the not operator
	NOT = "!"
	// AND the and operator
	AND = "&&"
	// OR the or operator
	OR = "||"

	//
	// Comparision operators
	//

	// LT the less than comparision operator
	LT = "<"
	// LTE  the less than or equal comparision operator
	LTE = "<="

	// GT the greater than comparision operator
	GT = ">"
	// GTE the grather than or equal comparision operator
	GTE = ">="

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
	// COLON a comon
	COLON = ":"
	// DOT a dot
	DOT = "."

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
	// TRUE the `true` keyword (true)
	TRUE = "TRUE"
	// FALSE the `false` keyword (false)
	FALSE = "FALSE"
	// NULL the `null` keyword (null)
	NULL = "NULL"
	// IF the `if` keyword (if)
	IF = "IF"
	// ELSE the `else` keyword (else)
	ELSE = "ELSE"
	// RETURN the `return` keyword (return)
	RETURN = "RETURN"
	// WHILE the `while` keyword (while)
	WHILE = "WHILE"
)

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
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
