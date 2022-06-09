package ebnf

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
