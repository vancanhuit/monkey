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
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	testCases := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Identifier, "five"},
		{token.Assign, "="},
		{token.Integer, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "ten"},
		{token.Assign, "="},
		{token.Integer, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LeftParen, "("},
		{token.Identifier, "x"},
		{token.Comma, ","},
		{token.Identifier, "y"},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.Identifier, "x"},
		{token.Plus, "+"},
		{token.Identifier, "y"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Identifier, "result"},
		{token.Assign, "="},
		{token.Identifier, "add"},
		{token.LeftParen, "("},
		{token.Identifier, "five"},
		{token.Comma, ","},
		{token.Identifier, "ten"},
		{token.RightParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Integer, "5"},
		{token.Semicolon, ";"},
		{token.Integer, "5"},
		{token.LessThan, "<"},
		{token.Integer, "10"},
		{token.GreaterThan, ">"},
		{token.Integer, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LeftParen, "("},
		{token.Integer, "5"},
		{token.LessThan, "<"},
		{token.Integer, "10"},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Else, "else"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Integer, "10"},
		{token.Equal, "=="},
		{token.Integer, "10"},
		{token.Semicolon, ";"},
		{token.Integer, "10"},
		{token.NotEqual, "!="},
		{token.Integer, "9"},
		{token.Semicolon, ";"},
		{token.String, "foobar"},
		{token.String, "foo bar"},
		{token.LeftBracket, "["},
		{token.Integer, "1"},
		{token.Comma, ","},
		{token.Integer, "2"},
		{token.RightBracket, "]"},
		{token.Semicolon, ";"},
		{token.LeftBrace, "{"},
		{token.String, "foo"},
		{token.Colon, ":"},
		{token.String, "bar"},
		{token.RightBrace, "}"},
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
