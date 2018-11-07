package ast

// Packge ast implement the Abstract Syntax Tree that represents the parsed
// source code before being passed on to the interpreter for evaluation.

import (
	"git.mills.io/prologic/monkey/token"
)

// Node defines an interface for all nodes in the AST.
type Node interface {
	TokenLiteral() string
}

// Statement defines the interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression defines the interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node. All programs consist of a slice of Statement(s)
type Program struct {
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement the `let` statement represents the AST node that binds an
// expression to an identifier
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier is a node that holds the literal value of an identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ReturnStatement the `return` statement that represents the AST node that
// holds a return value to the outter stack in the call stack.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
