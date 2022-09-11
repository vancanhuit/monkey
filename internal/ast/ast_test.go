package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vancanhuit/monkey/internal/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.Identifier,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	require.Equal(t, program.String(), "let myVar = anotherVar;")
}
