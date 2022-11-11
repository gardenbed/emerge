package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	comb "github.com/gardenbed/emerge/internal/combinator"
)

func TestToNum(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedOK     bool
	}{
		{
			name: "Successful",
			r: comb.Result{
				Val: comb.List{
					{Val: '6', Pos: 1},
					{Val: '9', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: 69, Pos: 1},
			expectedOK:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := toNum(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestToLetters(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedOK     bool
	}{
		{
			name: "Successful",
			r: comb.Result{
				Val: comb.List{
					{Val: 'f', Pos: 1},
					{Val: 'o', Pos: 2},
					{Val: 'o', Pos: 3},
				},
			},
			expectedResult: comb.Result{Val: "foo", Pos: 1},
			expectedOK:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := toLetters(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestNew(t *testing.T) {
	m := new(mockMappers)
	p := New(m)

	assert.NotNil(t, p)
}

func TestParser_char(t *testing.T) {
	tests := []struct {
		name        string
		m           *mockMappers
		in          comb.Input
		expectedOut comb.Output
		expectedOK  bool
	}{
		{
			name:        "Failure",
			m:           &mockMappers{},
			in:          NewInput(`µ`),
			expectedOut: comb.Output{},
			expectedOK:  false,
		},
		{
			name: "Success_Low",
			m:    &mockMappers{},
			in:   NewInput(` `),
			expectedOut: comb.Output{
				Result: comb.Result{Val: ' ', Pos: 0},
			},
			expectedOK: true,
		},
		{
			name: "Success_High",
			m:    &mockMappers{},
			in:   NewInput(`~`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: '~', Pos: 0},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.char(tc.in)

			assert.Equal(t, tc.expectedOut, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestParser_unescapedChar(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`*`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`a`),
			expectedInResult: comb.Result{Val: 'a', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.unescapedChar(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToUnescapedCharMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_escapedChar(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`a`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToEscapedCharMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`\*`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 0},
					{Val: '*', Pos: 1},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.escapedChar(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToEscapedCharMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_digit(t *testing.T) {
	tests := []struct {
		name        string
		m           *mockMappers
		in          comb.Input
		expectedOut comb.Output
		expectedOK  bool
	}{
		{
			name:        "Failure",
			m:           &mockMappers{},
			in:          NewInput(`a`),
			expectedOut: comb.Output{},
			expectedOK:  false,
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   NewInput(`7`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: '7', Pos: 0},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.digit(tc.in)

			assert.Equal(t, tc.expectedOut, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestParser_letter(t *testing.T) {
	tests := []struct {
		name        string
		m           *mockMappers
		in          comb.Input
		expectedOut comb.Output
		expectedOK  bool
	}{
		{
			name:        "Failure",
			m:           &mockMappers{},
			in:          NewInput(`0`),
			expectedOut: comb.Output{},
			expectedOK:  false,
		},
		{
			name: "Success_Upper",
			m:    &mockMappers{},
			in:   NewInput(`A`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: 'A', Pos: 0},
			},
			expectedOK: true,
		},
		{
			name: "Success_Lower",
			m:    &mockMappers{},
			in:   NewInput(`a`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: 'a', Pos: 0},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.letter(tc.in)

			assert.Equal(t, tc.expectedOut, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestParser_num(t *testing.T) {
	tests := []struct {
		name        string
		m           *mockMappers
		in          comb.Input
		expectedOut comb.Output
		expectedOK  bool
	}{
		{
			name:        "Failure",
			m:           &mockMappers{},
			in:          NewInput(`a`),
			expectedOut: comb.Output{},
			expectedOK:  false,
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   NewInput(`69`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: 69, Pos: 0},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.num(tc.in)

			assert.Equal(t, tc.expectedOut, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestParser_letters(t *testing.T) {
	tests := []struct {
		name        string
		m           *mockMappers
		in          comb.Input
		expectedOut comb.Output
		expectedOK  bool
	}{
		{
			name:        "Failure",
			m:           &mockMappers{},
			in:          NewInput(`0`),
			expectedOut: comb.Output{},
			expectedOK:  false,
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   NewInput(`Symbol`),
			expectedOut: comb.Output{
				Result: comb.Result{Val: "Symbol", Pos: 0},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.letters(tc.in)

			assert.Equal(t, tc.expectedOut, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestParser_repOp(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`!`),
		},
		{
			name: "Success_ZeroOrOne",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`?`),
			expectedInResult: comb.Result{Val: '?', Pos: 0},
		},
		{
			name: "Success_ZeroOrMany",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`*`),
			expectedInResult: comb.Result{Val: '*', Pos: 0},
		},
		{
			name: "Success_OneOrMany",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`+`),
			expectedInResult: comb.Result{Val: '+', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.repOp(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToRepOpMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_upperBound(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`;`),
		},
		{
			name: "Success_Unbounded",
			m: &mockMappers{
				ToUpperBoundMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`,`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: ',', Pos: 0},
					{Val: comb.Empty{}},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Bounded",
			m: &mockMappers{
				ToUpperBoundMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`,4`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: ',', Pos: 0},
					{Val: 4, Pos: 1},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.upperBound(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToUpperBoundMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_range(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`{`),
		},
		{
			name: "Success_WithoutUpperBound",
			m: &mockMappers{
				ToRangeMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`{2}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 0},
					{Val: 2, Pos: 1},
					{Val: comb.Empty{}},
					{Val: '}', Pos: 2},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_WithUpperBound",
			m: &mockMappers{
				ToUpperBoundMocks: []MapperMock{
					{OutOK: true},
				},
				ToRangeMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`{2,}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 0},
					{Val: 2, Pos: 1},
					{},
					{Val: '}', Pos: 3},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.range_(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToRangeMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_repetition(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`!`),
		},
		{
			name: "Success_RepOp",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`*`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Range",
			m: &mockMappers{
				ToUpperBoundMocks: []MapperMock{
					{OutOK: true},
				},
				ToRangeMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`{2,4}`),
			expectedInResult: comb.Result{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.repetition(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToRepetitionMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_quantifier(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`!`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{OutOK: true},
				},
				ToQuantifierMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`*`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: comb.Empty{}},
				},
			},
		},
		{
			name: "Success_Lazy",
			m: &mockMappers{
				ToRepOpMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{OutOK: true},
				},
				ToQuantifierMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`*?`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: '?', Pos: 1},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.quantifier(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToQuantifierMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_charRange(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`a`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToCharRangeMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`a-z`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: 'a', Pos: 0},
					{Val: '-', Pos: 1},
					{Val: 'z', Pos: 2},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.charRange(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToCharRangeMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_charClass(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\a`),
		},
		{
			name: "Success_Digit",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\d`),
			expectedInResult: comb.Result{Val: "\\d", Pos: 0},
		},
		{
			name: "Success_NotDigit",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\D`),
			expectedInResult: comb.Result{Val: "\\D", Pos: 0},
		},
		{
			name: "Success_Whitespace",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\s`),
			expectedInResult: comb.Result{Val: "\\s", Pos: 0},
		},
		{
			name: "Success_NotWhitespace",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\S`),
			expectedInResult: comb.Result{Val: "\\S", Pos: 0},
		},
		{
			name: "Success_Word",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\w`),
			expectedInResult: comb.Result{Val: "\\w", Pos: 0},
		},
		{
			name: "Success_NotWord",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\W`),
			expectedInResult: comb.Result{Val: "\\W", Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.charClass(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToCharClassMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_asciiCharClass(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`[:invalid:]`),
		},
		{
			name: "Success_Blank",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:blank:]`),
			expectedInResult: comb.Result{Val: "[:blank:]", Pos: 0},
		},
		{
			name: "Success_Space",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:space:]`),
			expectedInResult: comb.Result{Val: "[:space:]", Pos: 0},
		},
		{
			name: "Success_Digit",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:digit:]`),
			expectedInResult: comb.Result{Val: "[:digit:]", Pos: 0},
		},
		{
			name: "Success_XDigit",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:xdigit:]`),
			expectedInResult: comb.Result{Val: "[:xdigit:]", Pos: 0},
		},
		{
			name: "Success_Upper",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:upper:]`),
			expectedInResult: comb.Result{Val: "[:upper:]", Pos: 0},
		},
		{
			name: "Success_Lower",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:lower:]`),
			expectedInResult: comb.Result{Val: "[:lower:]", Pos: 0},
		},
		{
			name: "Success_Alpha",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:alpha:]`),
			expectedInResult: comb.Result{Val: "[:alpha:]", Pos: 0},
		},
		{
			name: "Success_Alnum",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:alnum:]`),
			expectedInResult: comb.Result{Val: "[:alnum:]", Pos: 0},
		},
		{
			name: "Success_Word",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:word:]`),
			expectedInResult: comb.Result{Val: "[:word:]", Pos: 0},
		},
		{
			name: "Success_ASCII",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:ascii:]`),
			expectedInResult: comb.Result{Val: "[:ascii:]", Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.asciiCharClass(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToASCIICharClassMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_charGroupItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_CharClass",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\d`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_ASCIICharClass",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:digit:]`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharRange",
			m: &mockMappers{
				ToCharRangeMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`a-z`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_EscapedChar",
			m: &mockMappers{
				ToEscapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\*`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_UnescapedChar",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`a`),
			expectedInResult: comb.Result{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.charGroupItem(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToCharGroupItemMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_charGroup(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`[`),
		},
		{
			name: "Success_Chars",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToCharGroupMocks: []MapperMock{
					MapperMock{},
				},
			},
			in: NewInput(`[ab]`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '[', Pos: 0},
					{Val: comb.Empty{}},
					{
						Val: comb.List{
							{},
							{},
						},
					},
					comb.Result{Val: ']', Pos: 3},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Negated_Chars",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToCharGroupMocks: []MapperMock{
					MapperMock{},
				},
			},
			in: NewInput(`[^ab]`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '[', Pos: 0},
					{Val: '^', Pos: 1},
					{
						Val: comb.List{
							{},
							{},
						},
					},
					comb.Result{Val: ']', Pos: 4},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.charGroup(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToCharGroupMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_anyChar(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`:`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToAnyCharMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`.`),
			expectedInResult: comb.Result{Val: '.', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.anyChar(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToAnyCharMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_matchItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_AnyChar",
			m: &mockMappers{
				ToAnyCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`.`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_UnescapedChar",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`a`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_EscapedChar",
			m: &mockMappers{
				ToEscapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\*`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharClass",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\d`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_ASCIICharClass",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[:digit:]`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharGroup",
			m: &mockMappers{
				ToCharRangeMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToCharGroupMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`[a-z]`),
			expectedInResult: comb.Result{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.matchItem(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToMatchItemMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_match(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_WithoutQuantifier",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`\d`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: comb.Empty{}},
				},
			},
		},
		{
			name: "Success_WithQuantifier",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepOpMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{OutOK: true},
				},
				ToQuantifierMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					MapperMock{},
				},
			},
			in: NewInput(`[:digit:]+`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.match(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToMatchMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_anchor(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`#`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToAnchorMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`$`),
			expectedInResult: comb.Result{Val: '$', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.anchor(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToAnchorMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_group(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`(`),
		},
		{
			name: "Success_WithoutQuantifier",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
				},
				ToGroupMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`(a)`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 0},
					{},
					{Val: ')', Pos: 2},
					{Val: comb.Empty{}},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_WithQuantifier",
			m: &mockMappers{
				ToEscapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepOpMocks: []MapperMock{
					{OutOK: true},
				},
				ToRepetitionMocks: []MapperMock{
					{OutOK: true},
				},
				ToQuantifierMocks: []MapperMock{
					{OutOK: true},
				},
				ToGroupMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`(\*)?`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 0},
					{},
					{Val: ')', Pos: 3},
					{},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.group(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToGroupMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, tc.m.ToGroupMocks[0].InResult)
			}
		})
	}
}

func TestParser_subexprItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_Anchor",
			m: &mockMappers{
				ToAnchorMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`$`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Group",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
				},
				ToGroupMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
					{},
				},
			},
			in:               NewInput(`(a)`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Match",
			m: &mockMappers{
				ToCharClassMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{},
				},
			},
			in:               NewInput(`\d`),
			expectedInResult: comb.Result{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.subexprItem(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToSubexprItemMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_subexpr(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_UnescapedChar",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					MapperMock{},
				},
			},
			in: NewInput(`ab`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.subexpr(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToSubexprMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_expr(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`a`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: comb.Empty{}},
				},
			},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
					{},
				},
			},
			in: NewInput(`a|b`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{
						Val: comb.List{
							{Val: '|', Pos: 1},
							{},
						},
						Pos: 1,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.expr(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToExprMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_regex(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name: "Failure",
			m:    &mockMappers{},
			in:   NewInput(`\`),
		},
		{
			name: "Success_WithoutStartOfString",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
				},
				ToRegexMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`a`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: comb.Empty{}},
					{},
				},
			},
		},
		{
			name: "Success_WithStartOfString",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
				},
				ToRegexMocks: []MapperMock{
					{},
				},
			},
			in: NewInput(`^a`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: '^', Pos: 0},
					{},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.regex(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToRegexMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name           string
		m              *mockMappers
		in             comb.Input
		expectedOutput comb.Output
		expectedOK     bool
	}{
		{
			name: "Success",
			m: &mockMappers{
				ToUnescapedCharMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToMatchMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprItemMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToSubexprMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToExprMocks: []MapperMock{
					{OutOK: true},
					{OutOK: true},
				},
				ToRegexMocks: []MapperMock{
					{OutOK: true},
				},
			},
			in:             NewInput(`a`),
			expectedOutput: comb.Output{},
			expectedOK:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, ok := p.Parse(tc.in)

			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

//==================================================< HELPERS >==================================================

type (
	MapperMock struct {
		InResult  comb.Result
		OutResult comb.Result
		OutOK     bool
	}

	// mockMappers implements the Mapper interface for testing purposes.
	mockMappers struct {
		ToUnescapedCharIndex int
		ToUnescapedCharMocks []MapperMock

		ToEscapedCharIndex int
		ToEscapedCharMocks []MapperMock

		ToRepOpIndex int
		ToRepOpMocks []MapperMock

		ToUpperBoundIndex int
		ToUpperBoundMocks []MapperMock

		ToRangeIndex int
		ToRangeMocks []MapperMock

		ToRepetitionIndex int
		ToRepetitionMocks []MapperMock

		ToQuantifierIndex int
		ToQuantifierMocks []MapperMock

		ToCharRangeIndex int
		ToCharRangeMocks []MapperMock

		ToCharGroupItemIndex int
		ToCharGroupItemMocks []MapperMock

		ToCharGroupIndex int
		ToCharGroupMocks []MapperMock

		ToASCIICharClassIndex int
		ToASCIICharClassMocks []MapperMock

		ToCharClassIndex int
		ToCharClassMocks []MapperMock

		ToAnyCharIndex int
		ToAnyCharMocks []MapperMock

		ToMatchItemIndex int
		ToMatchItemMocks []MapperMock

		ToMatchIndex int
		ToMatchMocks []MapperMock

		ToAnchorIndex int
		ToAnchorMocks []MapperMock

		ToGroupIndex int
		ToGroupMocks []MapperMock

		ToSubexprItemIndex int
		ToSubexprItemMocks []MapperMock

		ToSubexprIndex int
		ToSubexprMocks []MapperMock

		ToExprIndex int
		ToExprMocks []MapperMock

		ToRegexIndex int
		ToRegexMocks []MapperMock
	}
)

func (m *mockMappers) ToUnescapedChar(r comb.Result) (comb.Result, bool) {
	i := m.ToUnescapedCharIndex
	m.ToUnescapedCharIndex++
	m.ToUnescapedCharMocks[i].InResult = r
	return m.ToUnescapedCharMocks[i].OutResult, m.ToUnescapedCharMocks[i].OutOK
}

func (m *mockMappers) ToEscapedChar(r comb.Result) (comb.Result, bool) {
	i := m.ToEscapedCharIndex
	m.ToEscapedCharIndex++
	m.ToEscapedCharMocks[i].InResult = r
	return m.ToEscapedCharMocks[i].OutResult, m.ToEscapedCharMocks[i].OutOK
}

func (m *mockMappers) ToRepOp(r comb.Result) (comb.Result, bool) {
	i := m.ToRepOpIndex
	m.ToRepOpIndex++
	m.ToRepOpMocks[i].InResult = r
	return m.ToRepOpMocks[i].OutResult, m.ToRepOpMocks[i].OutOK
}

func (m *mockMappers) ToUpperBound(r comb.Result) (comb.Result, bool) {
	i := m.ToUpperBoundIndex
	m.ToUpperBoundIndex++
	m.ToUpperBoundMocks[i].InResult = r
	return m.ToUpperBoundMocks[i].OutResult, m.ToUpperBoundMocks[i].OutOK
}

func (m *mockMappers) ToRange(r comb.Result) (comb.Result, bool) {
	i := m.ToRangeIndex
	m.ToRangeIndex++
	m.ToRangeMocks[i].InResult = r
	return m.ToRangeMocks[i].OutResult, m.ToRangeMocks[i].OutOK
}

func (m *mockMappers) ToRepetition(r comb.Result) (comb.Result, bool) {
	i := m.ToRepetitionIndex
	m.ToRepetitionIndex++
	m.ToRepetitionMocks[i].InResult = r
	return m.ToRepetitionMocks[i].OutResult, m.ToRepetitionMocks[i].OutOK
}

func (m *mockMappers) ToQuantifier(r comb.Result) (comb.Result, bool) {
	i := m.ToQuantifierIndex
	m.ToQuantifierIndex++
	m.ToQuantifierMocks[i].InResult = r
	return m.ToQuantifierMocks[i].OutResult, m.ToQuantifierMocks[i].OutOK
}

func (m *mockMappers) ToCharRange(r comb.Result) (comb.Result, bool) {
	i := m.ToCharRangeIndex
	m.ToCharRangeIndex++
	m.ToCharRangeMocks[i].InResult = r
	return m.ToCharRangeMocks[i].OutResult, m.ToCharRangeMocks[i].OutOK
}

func (m *mockMappers) ToCharGroupItem(r comb.Result) (comb.Result, bool) {
	i := m.ToCharGroupItemIndex
	m.ToCharGroupItemIndex++
	m.ToCharGroupItemMocks[i].InResult = r
	return m.ToCharGroupItemMocks[i].OutResult, m.ToCharGroupItemMocks[i].OutOK
}

func (m *mockMappers) ToCharGroup(r comb.Result) (comb.Result, bool) {
	i := m.ToCharGroupIndex
	m.ToCharGroupIndex++
	m.ToCharGroupMocks[i].InResult = r
	return m.ToCharGroupMocks[i].OutResult, m.ToCharGroupMocks[i].OutOK
}

func (m *mockMappers) ToASCIICharClass(r comb.Result) (comb.Result, bool) {
	i := m.ToASCIICharClassIndex
	m.ToASCIICharClassIndex++
	m.ToASCIICharClassMocks[i].InResult = r
	return m.ToASCIICharClassMocks[i].OutResult, m.ToASCIICharClassMocks[i].OutOK
}

func (m *mockMappers) ToCharClass(r comb.Result) (comb.Result, bool) {
	i := m.ToCharClassIndex
	m.ToCharClassIndex++
	m.ToCharClassMocks[i].InResult = r
	return m.ToCharClassMocks[i].OutResult, m.ToCharClassMocks[i].OutOK
}

func (m *mockMappers) ToAnyChar(r comb.Result) (comb.Result, bool) {
	i := m.ToAnyCharIndex
	m.ToAnyCharIndex++
	m.ToAnyCharMocks[i].InResult = r
	return m.ToAnyCharMocks[i].OutResult, m.ToAnyCharMocks[i].OutOK
}

func (m *mockMappers) ToMatchItem(r comb.Result) (comb.Result, bool) {
	i := m.ToMatchItemIndex
	m.ToMatchItemIndex++
	m.ToMatchItemMocks[i].InResult = r
	return m.ToMatchItemMocks[i].OutResult, m.ToMatchItemMocks[i].OutOK
}

func (m *mockMappers) ToMatch(r comb.Result) (comb.Result, bool) {
	i := m.ToMatchIndex
	m.ToMatchIndex++
	m.ToMatchMocks[i].InResult = r
	return m.ToMatchMocks[i].OutResult, m.ToMatchMocks[i].OutOK
}

func (m *mockMappers) ToAnchor(r comb.Result) (comb.Result, bool) {
	i := m.ToAnchorIndex
	m.ToAnchorIndex++
	m.ToAnchorMocks[i].InResult = r
	return m.ToAnchorMocks[i].OutResult, m.ToAnchorMocks[i].OutOK
}

func (m *mockMappers) ToGroup(r comb.Result) (comb.Result, bool) {
	i := m.ToGroupIndex
	m.ToGroupIndex++
	m.ToGroupMocks[i].InResult = r
	return m.ToGroupMocks[i].OutResult, m.ToGroupMocks[i].OutOK
}

func (m *mockMappers) ToSubexprItem(r comb.Result) (comb.Result, bool) {
	i := m.ToSubexprItemIndex
	m.ToSubexprItemIndex++
	m.ToSubexprItemMocks[i].InResult = r
	return m.ToSubexprItemMocks[i].OutResult, m.ToSubexprItemMocks[i].OutOK
}

func (m *mockMappers) ToSubexpr(r comb.Result) (comb.Result, bool) {
	i := m.ToSubexprIndex
	m.ToSubexprIndex++
	m.ToSubexprMocks[i].InResult = r
	return m.ToSubexprMocks[i].OutResult, m.ToSubexprMocks[i].OutOK
}

func (m *mockMappers) ToExpr(r comb.Result) (comb.Result, bool) {
	i := m.ToExprIndex
	m.ToExprIndex++
	m.ToExprMocks[i].InResult = r
	return m.ToExprMocks[i].OutResult, m.ToExprMocks[i].OutOK
}

func (m *mockMappers) ToRegex(r comb.Result) (comb.Result, bool) {
	i := m.ToRegexIndex
	m.ToRegexIndex++
	m.ToRegexMocks[i].InResult = r
	return m.ToRegexMocks[i].OutResult, m.ToRegexMocks[i].OutOK
}
