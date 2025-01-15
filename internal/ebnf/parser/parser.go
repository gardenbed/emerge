package parser

import (
	"io"

	"github.com/moorara/algo/lexer"

	ebnflexer "github.com/gardenbed/emerge/internal/ebnf/lexer"
)

// Lexer is an interface for the lexer.Lexer struct.
type Lexer interface {
	NextToken() (lexer.Token, error)
}

// Parser is a syntax analyzer for an extension of EBNF language.
type Parser struct {
	lex Lexer
}

// New creates a new syntax analyzer for an extension of EBNF language.
func New(filename string, src io.Reader) (*Parser, error) {
	lex, err := ebnflexer.New(filename, src)
	if err != nil {
		return nil, err
	}

	return &Parser{
		lex: lex,
	}, nil
}
