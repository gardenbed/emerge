package parser

import (
	"io"

	"github.com/moorara/algo/lexer"

	ebnflexer "github.com/gardenbed/emerge/internal/ebnf/lexer"
)

// Parser is a parser (a.k.a. syntax analyzer) for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
type Parser struct {
	L lexer.Lexer
}

// New creates a new parser (a.k.a. syntax analyzer) for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
func New(filename string, src io.Reader) (*Parser, error) {
	L, err := ebnflexer.New(filename, src)
	if err != nil {
		return nil, err
	}

	return &Parser{
		L: L,
	}, nil
}
