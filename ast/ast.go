package ast

// Packge ast implement the Abstract Syntax Tree that represents the parsed
// source code before being passed on to the interpreter for evaluation.

import (
	"bytes"

	"git.mills.io/prologic/monkey/token"
)

// Node defines an interface for all nodes in the AST.
type Node interface {
	TokenLiteral() string
	String() string
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

// String returns a stringified version of the AST for debugging
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// String returns a stringified version of the `let` node
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier is a node that holds the literal value of an identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns a stringified version of the identifier node
func (i *Identifier) String() string { return i.Value }

// ReturnStatement the `return` statement that represents the AST node that
// holds a return value to the outter stack in the call stack.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a stringified version of the `return` node
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression node
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a stringified version of the expression node
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents a literal integer node
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the expression node
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression represents a prefix expression node
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a stringified version of the expression node
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
