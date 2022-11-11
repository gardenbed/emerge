package parser

import comb "github.com/gardenbed/emerge/internal/combinator"

// Input implements the input interface for strings.
type Input struct {
	pos   int
	runes []rune
}

func NewInput(s string) comb.Input {
	return &Input{
		pos:   0,
		runes: []rune(s),
	}
}

func (i *Input) Current() (rune, int) {
	return i.runes[0], i.pos
}

func (i *Input) Remaining() comb.Input {
	if len(i.runes) == 1 {
		return nil
	}

	return &Input{
		pos:   i.pos + 1,
		runes: i.runes[1:],
	}
}
