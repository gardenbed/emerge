package token

// Token is the type for all tokens in a grammar.
// Tokens are the leaves of abstract syntax tree.
type Token interface {
	Tag() Tag
	Pos() Pos
}

// Tag is the token tag.
type Tag int

const (
	INVALID Tag = iota

	CHAR    // [\x20-\x7E]
	NUM     // [0-9]+
	LETTERS // [A-Za-z]+

	START_OF_STRING // ^
	END_OF_STRING   // $
	ANY_CHAR        // .

	DIGIT          // "\d"
	NON_DIGIT      // "\D"
	WHITESPACE     // "\s"
	NON_WHITESPACE // "\S"
	WORD           // "\w"
	NON_WORD       // "\W"

	BLANK_CHARS  // "[:blank:]"
	SPACE_CHARS  // "[:space:]"
	DIGIT_CHARS  // "[:digit:]"
	XDIGIT_CHARS // "[:xdigit:]"
	UPPER_CHARS  // "[:upper:]"
	LOWER_CHARS  // "[:lower:]"
	ALPHA_CHARS  // "[:alpha:]"
	ALNUM_CHARS  // "[:alnum:]"
	WORD_CHARS   // "[:word:]"
	ASCII_CHARS  // "[:ascii:]"

	ZERO_OR_ONE  // ?
	ZERO_OR_MORE // *
	ONE_OR_MORE  // +
)

// Pos is the token position.
type Pos int
