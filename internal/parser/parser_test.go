package parser

import (
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

	require.Len(t, p.errors, 0)
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
