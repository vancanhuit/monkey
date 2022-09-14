package ast

import (
	"bytes"
	"strings"

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

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (stmt *BlockStatement) statementNode() {}
func (stmt *BlockStatement) TokenLiteral() string {
	return stmt.Token.Literal
}
func (stmt *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range stmt.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (expr *IfExpression) expressionNode() {}
func (expr *IfExpression) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(expr.Condition.String())
	out.WriteString(" ")
	out.WriteString(expr.Consequence.String())

	if expr.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(expr.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (expr *FunctionLiteral) expressionNode() {}
func (expr *FunctionLiteral) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range expr.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(expr.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(expr.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (expr *CallExpression) expressionNode() {}
func (expr *CallExpression) TokenLiteral() string {
	return expr.Token.Literal
}
func (expr *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range expr.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(expr.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (e *StringLiteral) expressionNode() {}
func (e *StringLiteral) TokenLiteral() string {
	return e.Token.Literal
}
func (e *StringLiteral) String() string {
	return e.Token.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (o *ArrayLiteral) expressionNode() {}
func (o *ArrayLiteral) TokenLiteral() string {
	return o.Token.Literal
}
func (o *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range o.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (o *IndexExpression) expressionNode() {}
func (o *IndexExpression) TokenLiteral() string {
	return o.Token.Literal
}
func (o *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(o.Left.String())
	out.WriteString("[")
	out.WriteString(o.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (o *HashLiteral) expressionNode() {}
func (o *HashLiteral) TokenLiteral() string {
	return o.Token.Literal
}
func (o *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range o.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
