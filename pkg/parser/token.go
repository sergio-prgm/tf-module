package parser

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	// Delimeters
	COMMA     = ","
	SEMICOLON = ";" // I don't think we need this

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// KEYWORDS
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	RESOURCE  = "RESOURCE"
	MODULE    = "MODULE"
	TERRAFORM = "TERRAFORM"
	PROVIDER  = "PROVIDER"
	IMPORT    = "IMPORT"
	VARIABLE  = "VARIABLE"
)

var Keywords = map[string]TokenType{
	"true":      TRUE,
	"false":     FALSE,
	"resource":  RESOURCE,
	"module":    MODULE,
	"terraform": TERRAFORM,
	"provider":  PROVIDER,
	"import":    IMPORT,
	"variable":  VARIABLE,
}
