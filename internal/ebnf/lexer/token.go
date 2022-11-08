package lexer

import "fmt"

// Tag is the type for token tags (names).
type Tag int

const (
	// QUO is the token tag for `"`.
	QUO Tag = 256 + iota
	// SOL is the token tag for `\`.
	SOL
	// DEF is the token tag for `=`.
	DEF
	// ALT is the token tag for `|`.
	ALT
	// LPAREN is the token tag for `(`.
	LPAREN
	// RPAREN is the token tag for `)`.
	RPAREN
	// LBRACK is the token tag for `[`.
	LBRACK
	// RBRACK is the token tag for `]`.
	RBRACK
	// LBRACE is the token tag for `{`.
	LBRACE
	// RBRACE is the token tag for `}`.
	RBRACE
	// LLBRACE is the token tag for `{{`.
	LLBRACE
	// RRBRACE is the token tag for `}}`.
	RRBRACE
	// GRAMMER is the token tag for `grammar`.
	GRAMMER
	// IDENT is the token tag for `[a-z][0-9a-z_]*`.
	IDENT
	// TOKEN is the token tag for `[A-Z][0-9A-Z_]*`.
	TOKEN
	// STRING is the token tag for `\"([\x21\x23-\x5B\x5D-\x7E]\|\\[\x21-\x7E]?)*\"`.
	STRING
	// REGEX is the token tag for `\/([\x21-\x2E\x30-\x5B\x5D-\x7E]\|\\[\x21-\x7E]?)*\/`.
	REGEX
)

// String implements the fmt.Stringer interface.
func (t Tag) String() string {
	switch t {
	case QUO:
		return "QUO"
	case SOL:
		return "SOL"
	case DEF:
		return "DEF"
	case ALT:
		return "ALT"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACK:
		return "LBRACK"
	case RBRACK:
		return "RBRACK"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LLBRACE:
		return "LLBRACE"
	case RRBRACE:
		return "RRBRACE"
	case GRAMMER:
		return "GRAMMER"
	case IDENT:
		return "IDENT"
	case TOKEN:
		return "TOKEN"
	case STRING:
		return "STRING"
	case REGEX:
		return "REGEX"
	default:
		return fmt.Sprintf("Tag(%d)", t)
	}
}

// Token is the type for a tuple consisting of a lexeme and an abstract symbol representing the kind of lexeme unit.
// It also contains some metadata about the lexeme, such as the starting position of the lexeme.
type Token struct {
	Tag    Tag
	Lexeme string
	Pos    int
}

// String implements the fmt.Stringer interface.
func (t Token) String() string {
	return fmt.Sprintf("%s<%s,%d>", t.Tag, t.Lexeme, t.Pos)
}
