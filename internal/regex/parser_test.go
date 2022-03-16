package regex

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// stringInput implements the parseInput interface for strings.
type stringInput struct {
	pos   int
	runes []rune
}

func newStringInput(s string) parseInput {
	return &stringInput{
		pos:   0,
		runes: []rune(s),
	}
}

func (s *stringInput) Current() (rune, int) {
	return s.runes[0], s.pos
}

func (s *stringInput) Remaining() parseInput {
	if len(s.runes) == 1 {
		return nil
	}

	return &stringInput{
		pos:   s.pos + 1,
		runes: s.runes[1:],
	}
}

func TestGetVal(t *testing.T) {
	tests := []struct {
		name           string
		res            any
		i              int
		expectedOK     bool
		expectedResult result
	}{
		{
			name:       "Input_Not_List",
			res:        'c',
			i:          2,
			expectedOK: false,
		},
		{
			name: "Index_LT_Zero",
			res: list{
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
			res: list{
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
			res: list{
				result{'a', 0},
				result{'b', 1},
				result{'c', 2},
				result{'d', 3},
			},
			i:              2,
			expectedOK:     true,
			expectedResult: result{'c', 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := getVal(tc.res, tc.i)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestExpectRune(t *testing.T) {
	tests := []struct {
		name        string
		in          parseInput
		r           rune
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
				Result:    result{'a', 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("ab"),
			r:          'a',
			expectedOK: true,
			expectedOut: parseOutput{
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
		in          parseInput
		runes       []rune
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
				Result:    result{'a', 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("ab"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: parseOutput{
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
		in          parseInput
		low, up     rune
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		runes       []rune
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
				Result:    result{[]rune{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Witouth_Remaining",
			in:         newStringInput("ab"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{[]rune{'a', 'b'}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			runes:      []rune{'a', 'b'},
			expectedOK: true,
			expectedOut: parseOutput{
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
		in          parseInput
		s           string
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
				Result:    result{"", 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("ab"),
			s:          "ab",
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			s:          "ab",
			expectedOK: true,
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		q           []parser
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		q           []parser
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		expectedOK  bool
		expectedOut parseOutput
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{empty{}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_Empty_Result",
			in:         newStringInput("ab"),
			p:          expectString("00"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{empty{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_Without_Remaining",
			in:         newStringInput("ab"),
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{"ab", 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_With_Remaining",
			in:         newStringInput("abcd"),
			p:          expectString("ab"),
			expectedOK: true,
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		expectedOK  bool
		expectedOut parseOutput
	}{
		{
			name:       "Input_Empty",
			in:         nil,
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{empty{}, 0},
				Remaining: nil,
			},
		},
		{
			name:       "Successful_Zero",
			in:         newStringInput("ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{empty{}, 0},
				Remaining: newStringInput("ab"),
			},
		},
		{
			name:       "Successful_One",
			in:         newStringInput("1ab"),
			p:          expectRuneInRange('0', '9'),
			expectedOK: true,
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		pos         []int
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		i           int
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
	toUpper := func(r result) (any, bool) {
		return strings.ToUpper(r.Val.(string)), true
	}

	tests := []struct {
		name        string
		in          parseInput
		p           parser
		f           converter
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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
		in          parseInput
		p           parser
		f           constructor
		expectedOK  bool
		expectedOut parseOutput
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
			expectedOut: parseOutput{
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
			expectedOut: parseOutput{
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

func TestParsers(t *testing.T) {
	tests := []struct {
		name        string
		p           parser
		in          parseInput
		expectedOK  bool
		expectedOut parseOutput
	}{
		{
			name:       "space_Successful",
			p:          _space,
			in:         newStringInput(" "),
			expectedOK: true,
			expectedOut: parseOutput{
				Result:    result{' ', 0},
				Remaining: nil,
			},
		},
		{
			name:       "char_Successful",
			p:          _char,
			in:         newStringInput(`!"#$%&'()*+,-./[\]^_{|}~`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Char{0, '!'},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune(`"#$%&'()*+,-./[\]^_{|}~`),
				},
			},
		},
		{
			name:       "digit_Successful",
			p:          _digit,
			in:         newStringInput("0123456789"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{'0', 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("123456789"),
				},
			},
		},
		{
			name:       "letter_Successful",
			p:          _letter,
			in:         newStringInput("abcdefghijklmnopqrstuvwxyz"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{'a', 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("bcdefghijklmnopqrstuvwxyz"),
				},
			},
		},
		{
			name:       "num_Successful",
			p:          _num,
			in:         newStringInput("2022tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Num{0, 2022},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "letters_Successful",
			p:          _letters,
			in:         newStringInput("head2022"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Letters{0, "head"},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("2022"),
				},
			},
		},
		{
			name:       "zeroOrOne_Successful",
			p:          _zeroOrOne,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ZeroOrOne{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "zeroOrMore_Successful",
			p:          _zeroOrMore,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ZeroOrMore{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "oneOrMore_Successful",
			p:          _oneOrMore,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: OneOrMore{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "upperBound_Unbounded_Successful",
			p:          _upperBound,
			in:         newStringInput(", }"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: UpperBound{},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("}"),
				},
			},
		},
		{
			name:       "upperBound_Bounded_Successful",
			p:          _upperBound,
			in:         newStringInput(",4}"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: UpperBound{
						CommaPos: 1,
						Val:      &Num{1, 4},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("}"),
				},
			},
		},
		{
			name:       "range_Fixed_Successful",
			p:          _range,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Range{
						OpenPos: 0,
						Low:     Num{1, 2},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "range_Upper_Unbounded_Successful",
			p:          _range,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Range{
						OpenPos: 0,
						Low:     Num{1, 2},
						Up:      &UpperBound{},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "range_Upper_Bounded_Successful",
			p:          _range,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Range{
						OpenPos: 0,
						Low:     Num{1, 2},
						Up: &UpperBound{
							CommaPos: 3,
							Val:      &Num{3, 4},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_zeroOrOne_Successful",
			p:          _cardinality,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &ZeroOrOne{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_zeroOrMore_Successful",
			p:          _cardinality,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &ZeroOrMore{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_oneOrMore_Successful",
			p:          _cardinality,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &OneOrMore{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_range_Fixed_Successful",
			p:          _cardinality,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Range{
						OpenPos: 0,
						Low:     Num{1, 2},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_range_Upper_Unbounded_Successful",
			p:          _cardinality,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Range{
						OpenPos: 0,
						Low:     Num{1, 2},
						Up:      &UpperBound{},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "cardinality_range_Upper_Bounded_Successful",
			p:          _cardinality,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Range{
						OpenPos: 0,
						Low:     Num{1, 2},
						Up: &UpperBound{
							CommaPos: 3,
							Val:      &Num{3, 4},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_zeroOrOne_Successful",
			p:          _quantifier,
			in:         newStringInput("??tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &ZeroOrOne{0},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_zeroOrMore_Successful",
			p:          _quantifier,
			in:         newStringInput("*?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &ZeroOrMore{0},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_oneOrMore_Successful",
			p:          _quantifier,
			in:         newStringInput("+?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &OneOrMore{0},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_Fixed_Successful",
			p:          _quantifier,
			in:         newStringInput("{2}?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &Range{
							OpenPos: 0,
							Low:     Num{1, 2},
						},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_Upper_Unbounded_Successful",
			p:          _quantifier,
			in:         newStringInput("{2,}?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &Range{
							OpenPos: 0,
							Low:     Num{1, 2},
							Up:      &UpperBound{},
						},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_Upper_Bounded_Successful",
			p:          _quantifier,
			in:         newStringInput("{2,4}?tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Quantifier{
						Card: &Range{
							OpenPos: 0,
							Low:     Num{1, 2},
							Up: &UpperBound{
								CommaPos: 3,
								Val:      &Num{3, 4},
							},
						},
						Lazy: true,
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charRange_Range_Successful",
			p:          _charRange,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharRange{
						Low: Char{0, '0'},
						Up:  Char{2, '9'},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Digit_Successful",
			p:          _charClass,
			in:         newStringInput(`\dtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, DIGIT_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotDigit_Successful",
			p:          _charClass,
			in:         newStringInput(`\Dtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, NOT_DIGIT_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Whitespace_Successful",
			p:          _charClass,
			in:         newStringInput(`\stail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, WHITESPACE},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotWhitespace_Successful",
			p:          _charClass,
			in:         newStringInput(`\Stail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, NOT_WHITESPACE},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Word_Successful",
			p:          _charClass,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, WORD_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotWord_Successful",
			p:          _charClass,
			in:         newStringInput(`\Wtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharClass{0, NOT_WORD_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Blank_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:blank:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, BLANK},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Space_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:space:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, SPACE},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Digit_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:digit:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, DIGIT},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_XDigit_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:xdigit:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, XDIGIT},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Upper_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:upper:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, UPPER},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Lower_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:lower:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, LOWER},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Alpha_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:alpha:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, ALPHA},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Alnum_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:alnum:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, ALNUM},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Word_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, WORD},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_ASCII_Successful",
			p:          _asciiCharClass,
			in:         newStringInput("[:ascii:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: ASCIICharClass{0, ASCII},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_charClass_Successful",
			p:          _charGroupItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &CharClass{0, WORD_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_asciiCharClass_Successful",
			p:          _charGroupItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &ASCIICharClass{0, WORD},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_charRange_Successful",
			p:          _charGroupItem,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &CharRange{
						Low: Char{0, '0'},
						Up:  Char{2, '9'},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_char_Successful",
			p:          _charGroupItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Char{0, '!'},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_charClass_Successful",
			p:          _charGroup,
			in:         newStringInput(`[\w]tail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharGroup{
						OpenPos: 0,
						Neg:     false,
						Items: []CharGroupItem{
							&CharClass{1, WORD_CHARS},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_asciiCharClass_Successful",
			p:          _charGroup,
			in:         newStringInput("[[:word:]]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharGroup{
						OpenPos: 0,
						Neg:     false,
						Items: []CharGroupItem{
							&ASCIICharClass{1, WORD},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_charRange_Successful",
			p:          _charGroup,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharGroup{
						OpenPos: 0,
						Neg:     false,
						Items: []CharGroupItem{
							&CharRange{
								Low: Char{1, '0'},
								Up:  Char{3, '9'},
							},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_char_Successful",
			p:          _charGroup,
			in:         newStringInput("[!]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharGroup{
						OpenPos: 0,
						Neg:     false,
						Items: []CharGroupItem{
							&Char{1, '!'},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_Negated_Successful",
			p:          _charGroup,
			in:         newStringInput("[^#$]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: CharGroup{
						OpenPos: 0,
						Neg:     true,
						Items: []CharGroupItem{
							&Char{2, '#'},
							&Char{3, '$'},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anyChar_Successful",
			p:          _anyChar,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: AnyChar{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_anyChar_Successful",
			p:          _matchItem,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &AnyChar{0},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charGroup_charRange_Successful",
			p:          _matchItem,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &CharGroup{
						OpenPos: 0,
						Neg:     false,
						Items: []CharGroupItem{
							&CharRange{
								Low: Char{1, '0'},
								Up:  Char{3, '9'},
							},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charClass_Word_Successful",
			p:          _matchItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &CharClass{0, WORD_CHARS},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_asciiCharClass_Word_Successful",
			p:          _matchItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &ASCIICharClass{0, WORD},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_char_Successful",
			p:          _matchItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Char{0, '!'},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charClass_Whitespace_Successful",
			p:          _match,
			in:         newStringInput(`\stail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Match{
						Item: &CharClass{0, WHITESPACE},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charGroup_charRange_Quantifier_Successful",
			p:          _match,
			in:         newStringInput("[0-9]{2,4}tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Match{
						Item: &CharGroup{
							OpenPos: 0,
							Neg:     false,
							Items: []CharGroupItem{
								&CharRange{
									Low: Char{1, '0'},
									Up:  Char{3, '9'},
								},
							},
						},
						Quant: &Quantifier{
							Card: &Range{
								OpenPos: 5,
								Low:     Num{6, 2},
								Up: &UpperBound{
									CommaPos: 8,
									Val:      &Num{8, 4},
								},
							},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "backref_Successful",
			p:          _backref,
			in:         newStringInput(`\2tail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Backref{
						SlashPos: 0,
						Ref:      Num{1, 2},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anchor_End_Of_String_Successful",
			p:          _anchor,
			in:         newStringInput("$tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Anchor{0, END_OF_STRING},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anchor_Word_Boundary_Successful",
			p:          _anchor,
			in:         newStringInput(`\btail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Anchor{0, WORD_BOUNDARY},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anchor_Not_Word_Boundary_Successful",
			p:          _anchor,
			in:         newStringInput(`\Btail`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Anchor{0, NOT_WORD_BOUNDARY},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "group_Successful",
			p:          _group,
			in:         newStringInput("(a|b)tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Group{
						OpenPos: 0,
						Expr: Expr{
							Sub: Subexpr{
								Items: []SubexprItem{
									&Match{
										Item: &Char{1, 'a'},
									},
								},
							},
							Expr: &Expr{
								Sub: Subexpr{
									Items: []SubexprItem{
										&Match{
											Item: &Char{3, 'b'},
										},
									},
								},
							},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "group_Non_Capturing_Successful",
			p:          _group,
			in:         newStringInput("(?:a|b)tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Group{
						OpenPos: 0,
						NonCap:  true,
						Expr: Expr{
							Sub: Subexpr{
								Items: []SubexprItem{
									&Match{
										Item: &Char{3, 'a'},
									},
								},
							},
							Expr: &Expr{
								Sub: Subexpr{
									Items: []SubexprItem{
										&Match{
											Item: &Char{5, 'b'},
										},
									},
								},
							},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   7,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "group_Quantifier_Successful",
			p:          _group,
			in:         newStringInput("(a|b)+tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Group{
						OpenPos: 0,
						Expr: Expr{
							Sub: Subexpr{
								Items: []SubexprItem{
									&Match{
										Item: &Char{1, 'a'},
									},
								},
							},
							Expr: &Expr{
								Sub: Subexpr{
									Items: []SubexprItem{
										&Match{
											Item: &Char{3, 'b'},
										},
									},
								},
							},
						},
						Quant: &Quantifier{
							Card: &OneOrMore{5},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_group_Successful",
			p:          _subexprItem,
			in:         newStringInput("(ab)+tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Group{
						OpenPos: 0,
						Expr: Expr{
							Sub: Subexpr{
								Items: []SubexprItem{
									&Match{
										Item: &Char{1, 'a'},
									},
									&Match{
										Item: &Char{2, 'b'},
									},
								},
							},
						},
						Quant: &Quantifier{
							Card: &OneOrMore{4},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_match_anyChar_Successful",
			p:          _subexprItem,
			in:         newStringInput(".*tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Match{
						Item: &AnyChar{0},
						Quant: &Quantifier{
							Card: &ZeroOrMore{1},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_match_charGroup_Successful",
			p:          _subexprItem,
			in:         newStringInput("[0-9]+tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: &Match{
						Item: &CharGroup{
							OpenPos: 0,
							Items: []CharGroupItem{
								&CharRange{
									Low: Char{1, '0'},
									Up:  Char{3, '9'},
								},
							},
						},
						Quant: &Quantifier{
							Card: &OneOrMore{5},
						},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexpr_group_matches_Successful",
			p:          _subexpr,
			in:         newStringInput("(ab)+[0-9]*tail"),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Subexpr{
						Items: []SubexprItem{
							&Group{
								OpenPos: 0,
								Expr: Expr{
									Sub: Subexpr{
										Items: []SubexprItem{
											&Match{
												Item: &Char{1, 'a'},
											},
											&Match{
												Item: &Char{2, 'b'},
											},
										},
									},
								},
								Quant: &Quantifier{
									Card: &OneOrMore{4},
								},
							},
							&Match{
								Item: &CharGroup{
									OpenPos: 5,
									Items: []CharGroupItem{
										&CharRange{
											Low: Char{6, '0'},
											Up:  Char{8, '9'},
										},
									},
								},
								Quant: &Quantifier{
									Card: &ZeroOrMore{10},
								},
							},
							&Match{
								Item: &Char{11, 't'},
							},
							&Match{
								Item: &Char{12, 'a'},
							},
							&Match{
								Item: &Char{13, 'i'},
							},
							&Match{
								Item: &Char{14, 'l'},
							},
						},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "expr_Successful",
			p:          _expr,
			in:         newStringInput(`[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Expr{
						Sub: Subexpr{
							Items: []SubexprItem{
								&Match{
									Item: &CharGroup{
										OpenPos: 0,
										Items: []CharGroupItem{
											&CharRange{
												Low: Char{1, '0'},
												Up:  Char{3, '9'},
											},
											&CharRange{
												Low: Char{4, 'A'},
												Up:  Char{6, 'Z'},
											},
											&CharRange{
												Low: Char{7, 'a'},
												Up:  Char{9, 'z'},
											},
											&Char{10, '_'},
										},
									},
								},
								&Match{
									Item: &CharGroup{
										OpenPos: 12,
										Items: []CharGroupItem{
											&CharClass{13, DIGIT_CHARS},
											&CharClass{15, WORD_CHARS},
										},
									},
									Quant: &Quantifier{
										Card: &ZeroOrMore{18},
									},
								},
							},
						},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "regex_Successful",
			p:          _regex,
			in:         newStringInput(`^package\s+[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: parseOutput{
				Result: result{
					Val: Regex{
						Begin: true,
						Expr: Expr{
							Sub: Subexpr{
								Items: []SubexprItem{
									&Match{
										Item: &Char{1, 'p'},
									},
									&Match{
										Item: &Char{2, 'a'},
									},
									&Match{
										Item: &Char{3, 'c'},
									},
									&Match{
										Item: &Char{4, 'k'},
									},
									&Match{
										Item: &Char{5, 'a'},
									},
									&Match{
										Item: &Char{6, 'g'},
									},
									&Match{
										Item: &Char{7, 'e'},
									},
									&Match{
										Item: &CharClass{8, WHITESPACE},
										Quant: &Quantifier{
											Card: &OneOrMore{10},
										},
									},
									&Match{
										Item: &CharGroup{
											OpenPos: 11,
											Items: []CharGroupItem{
												&CharRange{
													Low: Char{12, '0'},
													Up:  Char{14, '9'},
												},
												&CharRange{
													Low: Char{15, 'A'},
													Up:  Char{17, 'Z'},
												},
												&CharRange{
													Low: Char{18, 'a'},
													Up:  Char{20, 'z'},
												},
												&Char{21, '_'},
											},
										},
									},
									&Match{
										Item: &CharGroup{
											OpenPos: 23,
											Items: []CharGroupItem{
												&CharClass{24, DIGIT_CHARS},
												&CharClass{26, WORD_CHARS},
											},
										},
										Quant: &Quantifier{
											Card: &ZeroOrMore{29},
										},
									},
								},
							},
						},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := tc.p(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
