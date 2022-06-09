package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// empty --> ε
var _empty = func(in input) (output, bool) {
	return output{
		Result: result{
			Val: empty{},
		},
		Remaining: in,
	}, true
}

// stringInput implements the input interface for strings.
type stringInput struct {
	pos   int
	runes []rune
}

func newStringInput(s string) input {
	return &stringInput{
		pos:   0,
		runes: []rune(s),
	}
}

func (s *stringInput) Current() (rune, int) {
	return s.runes[0], s.pos
}

func (s *stringInput) Remaining() input {
	if len(s.runes) == 1 {
		return nil
	}

	return &stringInput{
		pos:   s.pos + 1,
		runes: s.runes[1:],
	}
}

func TestGetAt(t *testing.T) {
	tests := []struct {
		name          string
		v             any
		i             int
		expectedOK    bool
		expectedValue any
	}{
		{
			name:       "Input_Not_List",
			v:          'c',
			i:          2,
			expectedOK: false,
		},
		{
			name: "Index_LT_Zero",
			v: list{
				result{'a', 0},
				result{'b', 1},
				result{'c', 2},
				result{'d', 3},
			},
			i:          -1,
			expectedOK: false,
		},
		{
			name: "Index_GEQ_Len",
			v: list{
				result{'a', 0},
				result{'b', 1},
				result{'c', 2},
				result{'d', 3},
			},
			i:          4,
			expectedOK: false,
		},
		{
			name: "Successful",
			v: list{
				result{'a', 0},
				result{'b', 1},
				result{'c', 2},
				result{'d', 3},
			},
			i:             2,
			expectedOK:    true,
			expectedValue: 'c',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, ok := getAt(tc.v, tc.i)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedValue, v)
		})
	}
}

func TestExpectRune(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		r           rune
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			r:          'a',
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("a"),
			r:          'b',
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("a"),
			r:          'a',
			expectedOK: true,
			expectedOut: output{
				Result:    result{'a', 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("ab"),
			r:          'a',
			expectedOK: true,
			expectedOut: output{
				Result: result{'a', 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("b"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := expectRune(tc.r)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestExpectRuneIn(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		runes       []rune
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			runes:      []rune{'a', 'b'},
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("a"),
			runes:      []rune{'0', '1'},
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("a"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: output{
				Result:    result{'a', 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("ab"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: output{
				Result: result{'a', 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("b"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := expectRuneIn(tc.runes...)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestExpectRuneInRange(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		low, up     rune
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			low:        'a',
			up:         'z',
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("a"),
			low:        '0',
			up:         '9',
			expectedOK: false,
		},
		{
			name:       "Invalid_Range",
			in:         newStringInput("a"),
			low:        'z',
			up:         'a',
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("a"),
			low:        'a',
			up:         'z',
			expectedOK: true,
			expectedOut: output{
				Result:    result{'a', 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("ab"),
			low:        'a',
			up:         'z',
			expectedOK: true,
			expectedOut: output{
				Result: result{'a', 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("b"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := expectRuneInRange(tc.low, tc.up)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestExpectRunes(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		runes       []rune
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			runes:      []rune{'a', 'b'},
			expectedOK: false,
		},
		{
			name:       "Input_Not_Enough",
			in:         newStringInput("a"),
			runes:      []rune{'a', 'b'},
			expectedOK: false,
		},
		{
			name:       "Input_Not_Matching",
			in:         newStringInput("ab"),
			runes:      []rune{'0', '9'},
			expectedOK: false,
		},
		{
			name:       "Successful_Empty_Runes",
			in:         newStringInput("ab"),
			runes:      []rune{},
			expectedOK: true,
			expectedOut: output{
				Result:    result{[]rune{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Witouth_Remaining",
			in:         newStringInput("ab"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: output{
				Result:    result{[]rune{'a', 'b'}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: output{
				Result: result{[]rune{'a', 'b'}, 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := expectRunes(tc.runes...)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestExpectString(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		s           string
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			s:          "ab",
			expectedOK: false,
		},
		{
			name:       "Input_Not_Enough",
			in:         newStringInput("a"),
			s:          "ab",
			expectedOK: false,
		},
		{
			name:       "Input_Not_Matching",
			in:         newStringInput("ab"),
			s:          "09",
			expectedOK: false,
		},
		{
			name:       "Successful_Empty_String",
			in:         newStringInput("ab"),
			s:          "",
			expectedOK: true,
			expectedOut: output{
				Result:    result{"", 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("ab"),
			s:          "ab",
			expectedOK: true,
			expectedOut: output{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			s:          "ab",
			expectedOK: true,
			expectedOut: output{
				Result: result{"ab", 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := expectString(tc.s)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_CONCAT(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		q           []parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectString("a"),
			q:          []parser{expectString("b")},
			expectedOK: false,
		},
		{
			name:       "Input_Not_Enough",
			in:         newStringInput("a"),
			p:          expectString("a"),
			q:          []parser{expectString("b")},
			expectedOK: false,
		},
		{
			name:       "First_Parser_Unsuccessful",
			in:         newStringInput("abcd"),
			p:          expectString("00"),
			q:          []parser{expectString("cd")},
			expectedOK: false,
		},
		{
			name:       "Second_Parser_Unsuccessful",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			q:          []parser{expectString("00")},
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			q:          []parser{expectString("cd")},
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{"ab", 0},
						result{"cd", 2},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcdef"),
			p:          expectString("ab"),
			q:          []parser{expectString("cd")},
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{"ab", 0},
						result{"cd", 2},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("ef"),
				},
			},
		},
		{
			name:       "Unuccessful_Multiple_Parsers",
			in:         newStringInput("abcdefghijklmn"),
			p:          expectString("ab"),
			q:          []parser{expectString("cd"), expectString("ef"), expectString("ij")},
			expectedOK: false,
		},
		{
			name:       "Successful_Multiple_Parsers",
			in:         newStringInput("abcdefghijklmn"),
			p:          expectString("ab"),
			q:          []parser{expectString("cd"), expectString("ef"), expectString("gh"), expectString("ij")},
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{"ab", 0},
						result{"cd", 2},
						result{"ef", 4},
						result{"gh", 6},
						result{"ij", 8},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("klmn"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.CONCAT(tc.q...)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_ALT(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		q           []parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectString("ab"),
			q:          []parser{expectString("00")},
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("ab"),
			p:          expectString("00"),
			q:          []parser{expectString("11")},
			expectedOK: false,
		},
		{
			name:       "First_Parser_Successful",
			in:         newStringInput("ab"),
			p:          expectString("ab"),
			q:          []parser{expectString("00")},
			expectedOK: true,
			expectedOut: output{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Second_Parser_Successful",
			in:         newStringInput("ab"),
			p:          expectString("00"),
			q:          []parser{expectString("ab")},
			expectedOK: true,
			expectedOut: output{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			q:          []parser{expectString("cd")},
			expectedOK: true,
			expectedOut: output{
				Result: result{"ab", 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
		{
			name:       "Unsuccessful_Multiple_Parsers",
			in:         newStringInput("abcd"),
			p:          expectString("00"),
			q:          []parser{expectString("11"), expectString("22"), expectString("33"), expectString("44")},
			expectedOK: false,
		},
		{
			name:       "Successful_Multiple_Parsers",
			in:         newStringInput("abcd"),
			p:          expectString("00"),
			q:          []parser{expectString("11"), expectString("22"), expectString("33"), expectString("ab")},
			expectedOK: true,
			expectedOut: output{
				Result: result{"ab", 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.ALT(tc.q...)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_OPT(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: output{
				Result:    result{empty{}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_Empty_Result",
			in:         newStringInput("ab"),
			p:          expectString("00"),
			expectedOK: true,
			expectedOut: output{
				Result:    result{empty{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("ab"),
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: output{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: output{
				Result: result{"ab", 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.OPT()(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_REP(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result:    result{empty{}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_Zero",
			in:         newStringInput("ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result:    result{empty{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_One",
			in:         newStringInput("1ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("ab"),
				},
			},
		},
		{
			name:       "Successful_Many",
			in:         newStringInput("1234ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
						result{'2', 1},
						result{'3', 2},
						result{'4', 3},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("ab"),
				},
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("1234"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
						result{'2', 1},
						result{'3', 2},
						result{'4', 3},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.REP()(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_REP1(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRuneInRange('0', '9'),
			expectedOK: false,
		},
		{
			name:       "Unsuccessful_Zero",
			in:         newStringInput("ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: false,
		},
		{
			name:       "Successful_One",
			in:         newStringInput("1ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("ab"),
				},
			},
		},
		{
			name:       "Successful_Many",
			in:         newStringInput("1234ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
						result{'2', 1},
						result{'3', 2},
						result{'4', 3},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("ab"),
				},
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("1234"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'1', 0},
						result{'2', 1},
						result{'3', 2},
						result{'4', 3},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.REP1()(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_Flatten(t *testing.T) {
	rangeParser := expectRune('{').CONCAT(
		expectRuneInRange('0', '9'),
		expectRune(','),
		expectRune(' ').OPT(),
		expectRuneInRange('0', '9'),
		expectRune('}'),
	)

	tests := []struct {
		name        string
		in          input
		p           parser
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRune('!'),
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("{2,4}"),
			p:          expectRune('!'),
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("{2,4}"),
			p:          rangeParser,
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'{', 0},
						result{'2', 1},
						result{',', 2},
						result{'4', 3},
						result{'}', 4},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("{2,4}ab"),
			p:          rangeParser,
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'{', 0},
						result{'2', 1},
						result{',', 2},
						result{'4', 3},
						result{'}', 4},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("ab"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.Flatten()(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_Select(t *testing.T) {
	rangeParser := expectRune('{').CONCAT(
		expectRuneInRange('0', '9'),
		expectRune(','),
		expectRuneInRange('0', '9'),
		expectRune('}'),
	)

	tests := []struct {
		name        string
		in          input
		p           parser
		pos         []int
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRune('!'),
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("{2,4}"),
			p:          expectRune('!'),
			expectedOK: false,
		},
		{
			name:       "Result_Not_List",
			in:         newStringInput("{2,4}"),
			p:          expectString("{2,4}"),
			expectedOK: true,
			expectedOut: output{
				Result:    result{"{2,4}", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Indices_Invalid",
			in:         newStringInput("{2,4}"),
			p:          rangeParser,
			pos:        []int{-1, 5},
			expectedOK: true,
			expectedOut: output{
				Result:    result{Val: empty{}},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("{2,4}"),
			p:          rangeParser,
			pos:        []int{1, 3},
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'2', 1},
						result{'4', 3},
					},
					Pos: 1,
				},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("{2,4}ab"),
			p:          rangeParser,
			pos:        []int{1, 3},
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{'2', 1},
						result{'4', 3},
					},
					Pos: 1,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("ab"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.Select(tc.pos...)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_Get(t *testing.T) {
	tests := []struct {
		name        string
		in          input
		p           parser
		i           int
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRune('!'),
			i:          0,
			expectedOK: false,
		},
		{
			name:       "Parser_Unuccessful",
			in:         newStringInput("ab"),
			p:          expectRune('!'),
			i:          0,
			expectedOK: false,
		},
		{
			name:       "Result_Not_List",
			in:         newStringInput("abcd"),
			p:          expectString("abcd"),
			i:          -1,
			expectedOK: true,
			expectedOut: output{
				Result:    result{"abcd", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Index_LT_Zero",
			in:         newStringInput("abcd"),
			p:          expectRuneInRange('a', 'z').REP(),
			i:          -1,
			expectedOK: true,
			expectedOut: output{
				Result:    result{Val: empty{}},
				Remaining: nil,
			},
		},
		{
			name:       "Index_GEQ_Len",
			in:         newStringInput("abcd"),
			p:          expectRuneInRange('a', 'z').REP(),
			i:          4,
			expectedOK: true,
			expectedOut: output{
				Result:    result{Val: empty{}},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_CONCAT",
			in:         newStringInput("abcd"),
			p:          expectString("ab").CONCAT(expectString("cd")),
			i:          1,
			expectedOK: true,
			expectedOut: output{
				Result:    result{"cd", 2},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_REP",
			in:         newStringInput("abcd"),
			p:          expectRuneIn('a', 'b', 'c', 'd').REP(),
			i:          2,
			expectedOK: true,
			expectedOut: output{
				Result:    result{'c', 2},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_REP1",
			in:         newStringInput("abcd"),
			p:          expectRuneInRange('a', 'z').REP(),
			i:          3,
			expectedOK: true,
			expectedOut: output{
				Result:    result{'d', 3},
				Remaining: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.Get(tc.i)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_Convert(t *testing.T) {
	toUpper := func(v any) (any, bool) {
		return strings.ToUpper(v.(string)), true
	}

	tests := []struct {
		name        string
		in          input
		p           parser
		f           converter
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRune('!'),
			f:          toUpper,
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("ab"),
			p:          expectRune('!'),
			f:          toUpper,
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("ab"),
			p:          expectString("ab"),
			f:          toUpper,
			expectedOK: true,
			expectedOut: output{
				Result:    result{"AB", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			f:          toUpper,
			expectedOK: true,
			expectedOut: output{
				Result: result{"AB", 0},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("cd"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.Convert(tc.f)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}

func TestParser_Bind(t *testing.T) {
	annotate := func(r result) parser {
		if r.Val.(rune) == '(' {
			return expectRune(' ').REP().CONCAT(expectRune(')')).Get(1)
		}
		return _empty
	}

	tests := []struct {
		name        string
		in          input
		p           parser
		f           constructor
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRune('('),
			f:          annotate,
			expectedOK: false,
		},
		{
			name:       "Parser_Unsuccessful",
			in:         newStringInput("(  )"),
			p:          expectRune('['),
			f:          annotate,
			expectedOK: false,
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("(  )"),
			p:          expectRune('('),
			f:          annotate,
			expectedOK: true,
			expectedOut: output{
				Result:    result{')', 3},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("(  )tail"),
			p:          expectRune('('),
			f:          annotate,
			expectedOK: true,
			expectedOut: output{
				Result: result{')', 3},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p.Bind(tc.f)(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
