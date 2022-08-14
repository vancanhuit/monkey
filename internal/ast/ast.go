package ast

import "github.com/vancanhuit/monkey/internal/token"

type Node interface {
	TokenLiteral() string
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

type Identifier struct {
	Token token.Token
	Value string
}

func (expr *Identifier) expressionNode() {}
func (expr *Identifier) TokenLiteral() string {
	return expr.Token.Literal
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

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (stmt *ReturnStatement) statementNode() {}
func (stmt *ReturnStatement) TokenLiteral() string {
	return stmt.Token.Literal
}
