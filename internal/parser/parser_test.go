package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vancanhuit/monkey/internal/ast"
	"github.com/vancanhuit/monkey/internal/lexer"
)

func testIdentifier(t *testing.T, expr ast.Expression, value string) {
	t.Helper()

	identifier, ok := expr.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, identifier.Value, value)
	require.Equal(t, identifier.TokenLiteral(), value)
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) {
	t.Helper()

	i, ok := expr.(*ast.IntegerLiteral)
	require.True(t, ok)
	require.Equal(t, i.Value, value)
	require.Equal(t, i.TokenLiteral(), fmt.Sprintf("%d", i.Value))
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) {
	b, ok := expr.(*ast.Boolean)
	require.True(t, ok)
	require.Equal(t, b.Value, value)
	require.Equal(t, b.TokenLiteral(), fmt.Sprintf("%t", b.Value))
}

func testLiteralExpression(
	t *testing.T, expr ast.Expression, expected interface{}) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, expr, int64(v))
	case int64:
		testIntegerLiteral(t, expr, v)
	case string:
		testIdentifier(t, expr, v)
	case bool:
		testBooleanLiteral(t, expr, v)
	}
}

func testInfixExpression(
	t *testing.T, expr ast.Expression,
	left interface{}, operator string, right interface{}) {

	infixExpr, ok := expr.(*ast.InfixExpression)
	require.True(t, ok)

	testLiteralExpression(t, infixExpr.Left, left)
	require.Equal(t, infixExpr.Operator, operator)
	testLiteralExpression(t, infixExpr.Right, right)
}

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()

		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		require.Equal(t, "let", stmt.TokenLiteral())
		letStmt, ok := stmt.(*ast.LetStatement)
		require.True(t, ok)
		testIdentifier(t, letStmt.Name, tc.expectedIdentifier)
		testLiteralExpression(t, letStmt.Value, tc.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt := program.Statements[0]

		require.Equal(t, "return", stmt.TokenLiteral())
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		require.True(t, ok)
		testLiteralExpression(t, returnStmt.Value, tc.expectedValue)
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

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestBooleanExpression(t *testing.T) {
	testCases := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
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

		testBooleanLiteral(t, stmt.Expression, tc.expectedBoolean)
	}
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
	testIntegerLiteral(t, stmt.Expression, int64(5))
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
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

		testLiteralExpression(t, expr.Right, tc.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		testInfixExpression(
			t, stmt.Expression, tc.leftValue, tc.operator, tc.rightValue)
	}
}

func TestIfExpressionParsing(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	expr, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok)

	testInfixExpression(t, expr.Condition, "x", "<", "y")

	require.Len(t, expr.Consequence.Statements, 1)

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testIdentifier(t, consequence.Expression, "x")
}

func TestIfElseExpressionParsing(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	expr, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok)

	testInfixExpression(t, expr.Condition, "x", "<", "y")

	require.Len(t, expr.Consequence.Statements, 1)

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testIdentifier(t, consequence.Expression, "x")

	require.Len(t, expr.Alternative.Statements, 1)
	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	testCases := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		require.Equal(t, len(function.Parameters), len(tc.expectedParams))

		for i, identifier := range tc.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok)
	require.Len(t, function.Parameters, 2)

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	require.Len(t, function.Body.Statements, 1)

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	expr, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t, ok)

	testIdentifier(t, expr.Function, "add")

	require.Len(t, expr.Arguments, 3)

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
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
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok)
	require.Equal(t, "hello world", literal.Value)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	require.True(t, ok)
	require.Len(t, array.Elements, 3)

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	require.True(t, ok)

	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 3)

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		require.True(t, ok)
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 0)
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	require.Equal(t, len(expected), len(hash.Pairs))

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.Boolean)
		require.True(t, ok)

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	require.Equal(t, len(expected), len(hash.Pairs))

	for key, value := range hash.Pairs {
		integer, ok := key.(*ast.IntegerLiteral)
		require.True(t, ok)

		expectedValue := expected[integer.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok)

	require.Len(t, hash.Pairs, 3)

	testCases := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		require.True(t, ok)

		testFunc, ok := testCases[literal.String()]
		require.True(t, ok)

		testFunc(value)
	}
}
