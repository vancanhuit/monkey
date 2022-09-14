package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	Illegal = "ILLEGAL"

	Identifier = "IDENTIFIER"
	Integer    = "INTEGER"
	String     = "STRING"

	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"

	LessThan    = "<"
	GreaterThan = ">"

	Comma     = ","
	Semicolon = ";"
	Colon     = ":"

	LeftParen    = "("
	RightParen   = ")"
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"

	Function = "FUNCTION"
	Let      = "LET"
	Return   = "RETURN"
	If       = "IF"
	Else     = "ELSE"
	True     = "TRUE"
	False    = "FALSE"

	Equal    = "=="
	NotEqual = "!="
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"if":     If,
	"else":   Else,
	"return": Return,
	"true":   True,
	"false":  False,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return Identifier
}
