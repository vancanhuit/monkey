package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vancanhuit/monkey/internal/ast"
	"github.com/vancanhuit/monkey/internal/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)

	testCases := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tc := range testCases {
		stmt := program.Statements[i]
		require.Equal(t, "let", stmt.TokenLiteral())
		letStmt, ok := stmt.(*ast.LetStatement)
		require.True(t, ok)
		require.Equal(t, letStmt.Name.Value, tc.expectedIdentifier)
		require.Equal(t, letStmt.Name.TokenLiteral(), tc.expectedIdentifier)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		require.True(t, ok)
		require.Equal(t, "return", returnStmt.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	identifier, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, "foobar", identifier.Value)
	require.Equal(t, "foobar", identifier.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.True(t, ok)
	require.Equal(t, int64(5), literal.Value)
	require.Equal(t, "5", literal.TokenLiteral())
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok)
		require.Equal(t, tc.operator, expr.Operator)

		integer, ok := expr.Right.(*ast.IntegerLiteral)
		require.True(t, ok)
		require.Equal(t, integer.Value, tc.integerValue)

		require.Equal(t, fmt.Sprintf("%d", integer.Value), integer.TokenLiteral())
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		require.True(t, ok)

		integer, ok := expr.Left.(*ast.IntegerLiteral)
		require.True(t, ok)
		require.Equal(t, integer.Value, tc.leftValue)

		require.Equal(t, tc.operator, expr.Operator)

		integer, ok = expr.Right.(*ast.IntegerLiteral)
		require.True(t, ok)
		require.Equal(t, integer.Value, tc.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)

		actual := program.String()
		require.Equal(t, actual, tc.expected)
	}
}
