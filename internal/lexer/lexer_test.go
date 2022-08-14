package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vancanhuit/monkey/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`

	testCases := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.GREATER_THAN, ">"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LEFT_PAREN, "("},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.INTEGER, "10"},
		{token.EQUAL, "=="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for _, tc := range testCases {
		tok := l.NextToken()

		//t.Logf("test case [%d]\n", i)
		require.Equal(t, tok.Type, tc.expectedType)
		require.Equal(t, tok.Literal, tc.expectedLiteral)
	}
}
