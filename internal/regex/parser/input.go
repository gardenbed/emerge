package parser

import comb "github.com/gardenbed/emerge/internal/combinator"

// stringInput implements the combinator.Input interface for strings.
type stringInput struct {
	pos   int
	runes []rune
}

func newStringInput(s string) comb.Input {
	return &stringInput{
		pos:   0,
		runes: []rune(s),
	}
}

func (s *stringInput) Current() (rune, int) {
	return s.runes[0], s.pos
}

func (s *stringInput) Remaining() comb.Input {
	if len(s.runes) == 1 {
		return nil
	}

	return &stringInput{
		pos:   s.pos + 1,
		runes: s.runes[1:],
	}
}
