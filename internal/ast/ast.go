package ast

import (
	"bytes"

	"github.com/vancanhuit/monkey/internal/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (expr *Identifier) expressionNode() {}
func (expr *Identifier) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *Identifier) String() string {
	return expr.Value
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (stmt *LetStatement) statementNode() {}
func (stmt *LetStatement) TokenLiteral() string {
	return stmt.Token.Literal
}
func (stmt *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(stmt.TokenLiteral() + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" = ")
	if stmt.Value != nil {
		out.WriteString(stmt.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (stmt *ReturnStatement) statementNode() {}
func (stmt *ReturnStatement) TokenLiteral() string {
	return stmt.Token.Literal
}
func (stmt *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(stmt.TokenLiteral() + " ")
	if stmt.Value != nil {
		out.WriteString(stmt.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (stmt *ExpressionStatement) statementNode() {}
func (stmt *ExpressionStatement) TokenLiteral() string {
	return stmt.Token.Literal
}
func (stmt *ExpressionStatement) String() string {
	if stmt.Expression != nil {
		return stmt.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (stmt *IntegerLiteral) expressionNode() {}
func (stmt *IntegerLiteral) TokenLiteral() string {
	return stmt.Token.Literal
}
func (stmt *IntegerLiteral) String() string {
	return stmt.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (expr *PrefixExpression) expressionNode() {}
func (expr *PrefixExpression) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Operator)
	out.WriteString(expr.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (expr *InfixExpression) expressionNode() {}
func (expr *InfixExpression) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Left.String())
	out.WriteString(" " + expr.Operator + " ")
	out.WriteString(expr.Right.String())
	out.WriteString(")")
	return out.String()
}
