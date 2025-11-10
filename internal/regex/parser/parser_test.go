package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	comb "github.com/moorara/algo/parser/combinator"
)

func TestToDigit(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name:           "OK",
			r:              comb.Result{Val: '7', Pos: 1},
			expectedResult: comb.Result{Val: 7, Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toDigit(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToHexDigit(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name:           "Digit",
			r:              comb.Result{Val: '7', Pos: 1},
			expectedResult: comb.Result{Val: 7, Pos: 1},
			expectedError:  "",
		},
		{
			name:           "Hex",
			r:              comb.Result{Val: 'F', Pos: 1},
			expectedResult: comb.Result{Val: 15, Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toHexDigit(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToNum(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name: "OK",
			r: comb.Result{
				Val: comb.List{
					{Val: 6, Pos: 1},
					{Val: 9, Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: 69, Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toNum(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToLetters(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name: "OK",
			r: comb.Result{
				Val: comb.List{
					{Val: 'L', Pos: 1},
					{Val: 'o', Pos: 2},
					{Val: 'r', Pos: 3},
					{Val: 'e', Pos: 4},
					{Val: 'm', Pos: 5},
				},
			},
			expectedResult: comb.Result{Val: "Lorem", Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toLetters(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToEscapedChar(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name: "Backslash",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '\\', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '\\', Pos: 1},
			expectedError:  "",
		},
		{
			name: "HorizontalTab",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: 't', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '\t', Pos: 1},
			expectedError:  "",
		},
		{
			name: "NewLine",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: 'n', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '\n', Pos: 1},
			expectedError:  "",
		},
		{
			name: "CarriageReturn",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: 'r', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '\r', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Bar",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '|', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '|', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Dot",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '.', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '.', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Question",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '?', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '?', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Asterisk",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '*', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '*', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Plus",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '+', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '+', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Hyphen",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '-', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '-', Pos: 1},
			expectedError:  "",
		},
		{
			name: "OpenningParenthesis",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '(', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '(', Pos: 1},
			expectedError:  "",
		},
		{
			name: "ClosingParenthesis",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: ')', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: ')', Pos: 1},
			expectedError:  "",
		},
		{
			name: "OpenningBracket",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '[', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '[', Pos: 1},
			expectedError:  "",
		},
		{
			name: "ClosingBracket",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: ']', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: ']', Pos: 1},
			expectedError:  "",
		},
		{
			name: "OpenningBrace",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '{', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '{', Pos: 1},
			expectedError:  "",
		},
		{
			name: "ClosingBrace",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '}', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '}', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Caret",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '^', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '^', Pos: 1},
			expectedError:  "",
		},
		{
			name: "Dollar",
			r: comb.Result{
				Val: comb.List{
					{Val: '\\', Pos: 1},
					{Val: '$', Pos: 2},
				},
			},
			expectedResult: comb.Result{Val: '$', Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toEscapedChar(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToASCIIChar(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name: "OK",
			r: comb.Result{
				Val: comb.List{
					{Val: "\\x", Pos: 1},
					{Val: 0x4, Pos: 3},
					{Val: 0xD, Pos: 4},
				},
			},
			expectedResult: comb.Result{Val: 'M', Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toASCIIChar(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestToUnicodeChar(t *testing.T) {
	tests := []struct {
		name           string
		r              comb.Result
		expectedResult comb.Result
		expectedError  string
	}{
		{
			name: "OK",
			r: comb.Result{
				Val: comb.List{
					{Val: "\\x", Pos: 1},
					{Val: 0x0, Pos: 3},
					{Val: 0x1, Pos: 4},
					{Val: 0xA, Pos: 5},
					{Val: 0x9, Pos: 6},
				},
			},
			expectedResult: comb.Result{Val: 'Ʃ', Pos: 1},
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := toUnicodeChar(tc.r)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestNew(t *testing.T) {
	m := new(mockMappers)
	p := New(m)

	assert.NotNil(t, p)
}

func TestParser_digit(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`a`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune 'a'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`7`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 7, Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.digit(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_hexDigit(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`a`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune 'a'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`A`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 10, Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.hexDigit(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_letter(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`0`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune '0'",
		},
		{
			name: "Success_Upper",
			m:    &mockMappers{},
			in:   newStringInput(`A`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 'A', Pos: 0},
			},
			expectedError: "",
		},
		{
			name: "Success_Lower",
			m:    &mockMappers{},
			in:   newStringInput(`a`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 'a', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.letter(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_num(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`a`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune 'a'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`69`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 69, Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.num(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_letters(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`0`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune '0'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`Symbol`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: "Symbol", Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.letters(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_char(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput("\x00"),
			expectedOut:   nil,
			expectedError: "0: unexpected rune '\\x00'",
		},
		{
			name: "Success_Low",
			m:    &mockMappers{},
			in:   newStringInput(` `),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: ' ', Pos: 0},
			},
			expectedError: "",
		},
		{
			name: "Success_High",
			m:    &mockMappers{},
			in:   newStringInput(`~`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: '~', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.char(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_unescapedChar(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`*`),
			expectedOut:   nil,
			expectedError: "unexpected rune '*' at position 0",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`a`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 'a', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.unescapedChar(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_escapedChar(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`a`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune 'a'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`\*`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: '*', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.escapedChar(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_asciiChar(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`4D`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune '4'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`\x4D`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 'M', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.asciiChar(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_unicodeChar(t *testing.T) {
	tests := []struct {
		name          string
		m             *mockMappers
		in            comb.Input
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name:          "Failure",
			m:             &mockMappers{},
			in:            newStringInput(`01A9`),
			expectedOut:   nil,
			expectedError: "0: unexpected rune '0'",
		},
		{
			name: "Success",
			m:    &mockMappers{},
			in:   newStringInput(`\x01A9`),
			expectedOut: &comb.Output{
				Result: comb.Result{Val: 'Ʃ', Pos: 0},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.unicodeChar(tc.in)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`:`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToAnyCharMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`.`),
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

func TestParser_singleChar(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput("\x00"),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`a`),
			expectedInResult: comb.Result{Val: 'a', Pos: 0},
		},
		{
			name: "Success_ASCII",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\x40`),
			expectedInResult: comb.Result{Val: '@', Pos: 0},
		},
		{
			name: "Success_Unicode",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\x01A9`),
			expectedInResult: comb.Result{Val: 'Ʃ', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.singleChar(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToSingleCharMocks; len(m) > 0 {
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\a`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Whitespace",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\s`),
			expectedInResult: comb.Result{Val: "\\s", Pos: 0},
		},
		{
			name: "Success_NotWhitespace",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\S`),
			expectedInResult: comb.Result{Val: "\\S", Pos: 0},
		},
		{
			name: "Success_Digit",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\d`),
			expectedInResult: comb.Result{Val: "\\d", Pos: 0},
		},
		{
			name: "Success_NotDigit",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\D`),
			expectedInResult: comb.Result{Val: "\\D", Pos: 0},
		},
		{
			name: "Success_Word",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\w`),
			expectedInResult: comb.Result{Val: "\\w", Pos: 0},
		},
		{
			name: "Success_NotWord",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\W`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`[:invalid:]`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Blank",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:blank:]`),
			expectedInResult: comb.Result{Val: "[:blank:]", Pos: 0},
		},
		{
			name: "Success_Space",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:space:]`),
			expectedInResult: comb.Result{Val: "[:space:]", Pos: 0},
		},
		{
			name: "Success_Digit",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:digit:]`),
			expectedInResult: comb.Result{Val: "[:digit:]", Pos: 0},
		},
		{
			name: "Success_XDigit",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:xdigit:]`),
			expectedInResult: comb.Result{Val: "[:xdigit:]", Pos: 0},
		},
		{
			name: "Success_Upper",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:upper:]`),
			expectedInResult: comb.Result{Val: "[:upper:]", Pos: 0},
		},
		{
			name: "Success_Lower",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:lower:]`),
			expectedInResult: comb.Result{Val: "[:lower:]", Pos: 0},
		},
		{
			name: "Success_Alpha",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:alpha:]`),
			expectedInResult: comb.Result{Val: "[:alpha:]", Pos: 0},
		},
		{
			name: "Success_Alnum",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:alnum:]`),
			expectedInResult: comb.Result{Val: "[:alnum:]", Pos: 0},
		},
		{
			name: "Success_Word",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:word:]`),
			expectedInResult: comb.Result{Val: "[:word:]", Pos: 0},
		},
		{
			name: "Success_ASCII",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:ascii:]`),
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

func TestParser_unicodeCategory(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput("Invalid"),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Math",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Math"),
			expectedInResult: comb.Result{Val: "Math", Pos: 0},
		},
		{
			name: "Success_Emoji",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Emoji"),
			expectedInResult: comb.Result{Val: "Emoji", Pos: 0},
		},
		{
			name: "Success_Latin",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Latin"),
			expectedInResult: comb.Result{Val: "Latin", Pos: 0},
		},
		{
			name: "Success_Greek",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Greek"),
			expectedInResult: comb.Result{Val: "Greek", Pos: 0},
		},
		{
			name: "Success_Cyrillic",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Cyrillic"),
			expectedInResult: comb.Result{Val: "Cyrillic", Pos: 0},
		},
		{
			name: "Success_Han",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Han"),
			expectedInResult: comb.Result{Val: "Han", Pos: 0},
		},
		{
			name: "Success_Persian",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Persian"),
			expectedInResult: comb.Result{Val: "Persian", Pos: 0},
		},
		{
			name: "Success_Letter",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Letter"),
			expectedInResult: comb.Result{Val: "Letter", Pos: 0},
		},
		{
			name: "Success_Lu",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Lu"),
			expectedInResult: comb.Result{Val: "Lu", Pos: 0},
		},
		{
			name: "Success_Ll",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Ll"),
			expectedInResult: comb.Result{Val: "Ll", Pos: 0},
		},
		{
			name: "Success_Lt",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Lt"),
			expectedInResult: comb.Result{Val: "Lt", Pos: 0},
		},
		{
			name: "Success_Lm",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Lm"),
			expectedInResult: comb.Result{Val: "Lm", Pos: 0},
		},
		{
			name: "Success_Lo",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Lo"),
			expectedInResult: comb.Result{Val: "Lo", Pos: 0},
		},
		{
			name: "Success_L",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("L"),
			expectedInResult: comb.Result{Val: "L", Pos: 0},
		},
		{
			name: "Success_Mark",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Mark"),
			expectedInResult: comb.Result{Val: "Mark", Pos: 0},
		},
		{
			name: "Success_Mn",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Mn"),
			expectedInResult: comb.Result{Val: "Mn", Pos: 0},
		},
		{
			name: "Success_Mc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Mc"),
			expectedInResult: comb.Result{Val: "Mc", Pos: 0},
		},
		{
			name: "Success_Me",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Me"),
			expectedInResult: comb.Result{Val: "Me", Pos: 0},
		},
		{
			name: "Success_M",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("M"),
			expectedInResult: comb.Result{Val: "M", Pos: 0},
		},
		{
			name: "Success_Number",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Number"),
			expectedInResult: comb.Result{Val: "Number", Pos: 0},
		},
		{
			name: "Success_Nd",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Nd"),
			expectedInResult: comb.Result{Val: "Nd", Pos: 0},
		},
		{
			name: "Success_Nl",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Nl"),
			expectedInResult: comb.Result{Val: "Nl", Pos: 0},
		},
		{
			name: "Success_No",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("No"),
			expectedInResult: comb.Result{Val: "No", Pos: 0},
		},
		{
			name: "Success_N",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("N"),
			expectedInResult: comb.Result{Val: "N", Pos: 0},
		},
		{
			name: "Success_Punctuation",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Punctuation"),
			expectedInResult: comb.Result{Val: "Punctuation", Pos: 0},
		},
		{
			name: "Success_Pc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Pc"),
			expectedInResult: comb.Result{Val: "Pc", Pos: 0},
		},
		{
			name: "Success_Pd",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Pd"),
			expectedInResult: comb.Result{Val: "Pd", Pos: 0},
		},
		{
			name: "Success_Ps",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Ps"),
			expectedInResult: comb.Result{Val: "Ps", Pos: 0},
		},
		{
			name: "Success_Pe",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Pe"),
			expectedInResult: comb.Result{Val: "Pe", Pos: 0},
		},
		{
			name: "Success_Pi",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Pi"),
			expectedInResult: comb.Result{Val: "Pi", Pos: 0},
		},
		{
			name: "Success_Pf",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Pf"),
			expectedInResult: comb.Result{Val: "Pf", Pos: 0},
		},
		{
			name: "Success_Po",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Po"),
			expectedInResult: comb.Result{Val: "Po", Pos: 0},
		},
		{
			name: "Success_P",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("P"),
			expectedInResult: comb.Result{Val: "P", Pos: 0},
		},
		{
			name: "Success_Separator",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Separator"),
			expectedInResult: comb.Result{Val: "Separator", Pos: 0},
		},
		{
			name: "Success_Zs",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Zs"),
			expectedInResult: comb.Result{Val: "Zs", Pos: 0},
		},
		{
			name: "Success_Zl",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Zl"),
			expectedInResult: comb.Result{Val: "Zl", Pos: 0},
		},
		{
			name: "Success_Zp",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Zp"),
			expectedInResult: comb.Result{Val: "Zp", Pos: 0},
		},
		{
			name: "Success_Z",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Z"),
			expectedInResult: comb.Result{Val: "Z", Pos: 0},
		},
		{
			name: "Success_Symbol",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Symbol"),
			expectedInResult: comb.Result{Val: "Symbol", Pos: 0},
		},
		{
			name: "Success_Sm",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Sm"),
			expectedInResult: comb.Result{Val: "Sm", Pos: 0},
		},
		{
			name: "Success_Sc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Sc"),
			expectedInResult: comb.Result{Val: "Sc", Pos: 0},
		},
		{
			name: "Success_Sk",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("Sk"),
			expectedInResult: comb.Result{Val: "Sk", Pos: 0},
		},
		{
			name: "Success_So",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("So"),
			expectedInResult: comb.Result{Val: "So", Pos: 0},
		},
		{
			name: "Success_S",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput("S"),
			expectedInResult: comb.Result{Val: "S", Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.unicodeCategory(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToUnicodeCategoryMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
		})
	}
}

func TestParser_unicodeCharClass(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\p{Invalid}`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Math",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Math", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Math}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Math", Pos: 3},
					{Val: '}', Pos: 7},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Math_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Math", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Math}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Math", Pos: 3},
					{Val: '}', Pos: 7},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Emoji",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Emoji", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Emoji}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Emoji", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Emoji_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Emoji", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Emoji}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Emoji", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Latin",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Latin", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Latin}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Latin", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Latin_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Latin", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Latin}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Latin", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Greek",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Greek", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Greek}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Greek", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Greek_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Greek", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Greek}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Greek", Pos: 3},
					{Val: '}', Pos: 8},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Cyrillic",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Cyrillic", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Cyrillic}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Cyrillic", Pos: 3},
					{Val: '}', Pos: 11},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Cyrillic_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Cyrillic", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Cyrillic}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Cyrillic", Pos: 3},
					{Val: '}', Pos: 11},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Han",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Han", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Han}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Han", Pos: 3},
					{Val: '}', Pos: 6},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Han_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Han", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Han}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Han", Pos: 3},
					{Val: '}', Pos: 6},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Persian",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Persian", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Persian}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Persian", Pos: 3},
					{Val: '}', Pos: 10},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Persian_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Persian", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Persian}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Persian", Pos: 3},
					{Val: '}', Pos: 10},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Letter",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Letter", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Letter}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Letter", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Letter_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Letter", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Letter}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Letter", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lu",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lu", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Lu}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lu", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lu_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lu", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Lu}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lu", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Ll",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Ll", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Ll}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Ll", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Ll_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Ll", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Ll}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Ll", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lt",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lt", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Lt}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lt", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lt_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lt", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Lt}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lt", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lm",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lm", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Lm}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lm", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lm_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lm", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Lm}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lm", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lo",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lo", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Lo}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lo", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Lo_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Lo", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Lo}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Lo", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_L",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "L", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{L}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "L", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_L_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "L", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{L}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "L", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mark",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mark", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Mark}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mark", Pos: 3},
					{Val: '}', Pos: 7},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mark_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mark", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Mark}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mark", Pos: 3},
					{Val: '}', Pos: 7},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mn",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mn", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Mn}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mn", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mn_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mn", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Mn}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mn", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Mc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Mc_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Mc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Mc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Mc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Me",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Me", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Me}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Me", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Me_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Me", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Me}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Me", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_M",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "M", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{M}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "M", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_M_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "M", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{M}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "M", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Number",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Number", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Number}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Number", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Number_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Number", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Number}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Number", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Nd",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Nd", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Nd}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Nd", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Nd_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Nd", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Nd}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Nd", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Nl",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Nl", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Nl}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Nl", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Nl_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Nl", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Nl}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Nl", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_No",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "No", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{No}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "No", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_No_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "No", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{No}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "No", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_N",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "N", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{N}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "N", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_N_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "N", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{N}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "N", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Punctuation",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Punctuation", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Punctuation}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Punctuation", Pos: 3},
					{Val: '}', Pos: 14},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Punctuation_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Punctuation", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Punctuation}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Punctuation", Pos: 3},
					{Val: '}', Pos: 14},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Pc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pc_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Pc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pd",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pd", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Pd}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pd", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pd_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pd", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Pd}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pd", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Ps",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Ps", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Ps}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Ps", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Ps_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Ps", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Ps}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Ps", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pe",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pe", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Pe}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pe", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pe_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pe", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Pe}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pe", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pi",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pi", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Pi}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pi", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pi_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pi", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Pi}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pi", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pf",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pf", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Pf}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pf", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Pf_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Pf", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Pf}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Pf", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Po",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Po", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Po}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Po", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Po_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Po", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Po}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Po", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_P",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "P", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{P}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "P", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_P_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "P", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{P}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "P", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Separator",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Separator", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Separator}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Separator", Pos: 3},
					{Val: '}', Pos: 12},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Separator_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Separator", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Separator}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Separator", Pos: 3},
					{Val: '}', Pos: 12},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zs",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zs", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Zs}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zs", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zs_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zs", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Zs}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zs", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zl",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zl", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Zl}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zl", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zl_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zl", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Zl}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zl", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zp",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zp", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Zp}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zp", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Zp_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Zp", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Zp}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Zp", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Z",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Z", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Z}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Z", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Z_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Z", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Z}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Z", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Symbol",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Symbol", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Symbol}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Symbol", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Symbol_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Symbol", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Symbol}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Symbol", Pos: 3},
					{Val: '}', Pos: 9},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sm",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sm", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Sm}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sm", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sm_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sm", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Sm}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sm", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sc",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Sc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sc_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sc", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Sc}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sc", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sk",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sk", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{Sk}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sk", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Sk_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "Sk", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{Sk}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "Sk", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_So",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "So", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{So}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "So", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_So_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "So", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{So}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "So", Pos: 3},
					{Val: '}', Pos: 5},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_S",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "S", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\p{S}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\p`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "S", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_S_Negated",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{
						OutResult: comb.Result{Val: "S", Pos: 3, Bag: nil},
						OutError:  nil,
					},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\P{S}`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{Val: `\P`, Pos: 0},
					{Val: '{', Pos: 2},
					{Val: "S", Pos: 3},
					{Val: '}', Pos: 4},
				},
				Pos: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.unicodeCharClass(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToUnicodeCharClassMocks; len(m) > 0 {
				assert.Equal(t, tc.expectedInResult, m[len(m)-1].InResult)
			}
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`!`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_ZeroOrOne",
			m: &mockMappers{
				ToRepOpMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`?`),
			expectedInResult: comb.Result{Val: '?', Pos: 0},
		},
		{
			name: "Success_ZeroOrMany",
			m: &mockMappers{
				ToRepOpMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`*`),
			expectedInResult: comb.Result{Val: '*', Pos: 0},
		},
		{
			name: "Success_OneOrMany",
			m: &mockMappers{
				ToRepOpMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`+`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`;`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Unbounded",
			m: &mockMappers{
				ToUpperBoundMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`,`),
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
				ToUpperBoundMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`,4`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`{`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_WithoutUpperBound",
			m: &mockMappers{
				ToRangeMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`{2}`),
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
				ToUpperBoundMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRangeMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`{2,}`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`!`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_RepOp",
			m: &mockMappers{
				ToRepOpMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`*`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Range",
			m: &mockMappers{
				ToUpperBoundMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRangeMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`{2,4}`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`!`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToRepOpMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToQuantifierMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`*`),
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
				ToRepOpMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToQuantifierMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`*?`),
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

func TestParser_charInRange(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput("\x00"),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`a`),
			expectedInResult: comb.Result{Val: 'a', Pos: 0},
		},
		{
			name: "Success_ASCII",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\x61`),
			expectedInResult: comb.Result{Val: 'a', Pos: 0},
		},
		{
			name: "Success_Unicode",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\x01A9`),
			expectedInResult: comb.Result{Val: 'Ʃ', Pos: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			p.charInRange(tc.in)

			// Verify the expected result has been passed to the mapper function
			if m := tc.m.ToCharInRangeMocks; len(m) > 0 {
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
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
				},
			},
			in:               newStringInput(`a-`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharRangeMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`a-z`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: '-', Pos: 1},
					{},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_ASCII",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharRangeMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\x61-\x7A`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: '-', Pos: 4},
					{},
				},
				Pos: 0,
			},
		},
		{
			name: "Success_Unicode",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharRangeMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\x03F0-\x03FF`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: '-', Pos: 6},
					{},
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

func TestParser_charGroupItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput("\x00"),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_UnicodeCharClass",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\p{Letter}`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_ASCIICharClass",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:digit:]`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharClass",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\d`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharRange",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharRangeMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`a-z`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_SingleChar",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`a`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`[`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Chars",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
					{OutError: nil}, // ']'
				},
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharGroupMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`[ab]`),
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
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
					{OutError: nil}, // ']'
				},
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToCharGroupMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`[^ab]`),
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

func TestParser_matchItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_AnyChar",
			m: &mockMappers{
				ToAnyCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`.`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_SingleChar",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`a`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharClass",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\d`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_ASCIICharClass",
			m: &mockMappers{
				ToASCIICharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[:digit:]`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_UnicodeCharClass",
			m: &mockMappers{
				ToUnicodeCategoryMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToUnicodeCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\p{Letter}`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_CharGroup",
			m: &mockMappers{
				ToCharInRangeMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
					{OutError: nil}, // ]
				},
				ToCharRangeMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToCharGroupMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`[a-z]`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_WithoutQuantifier",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`\d`),
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
				ToASCIICharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepOpMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToQuantifierMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`[:digit:]+`),
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

func TestParser_group(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`(`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_WithoutQuantifier",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToGroupMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`(a)`),
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
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepOpMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRepetitionMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToQuantifierMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToGroupMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`(\*)?`),
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

func TestParser_anchor(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`#`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToAnchorMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`$`),
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

func TestParser_subexprItem(t *testing.T) {
	tests := []struct {
		name             string
		m                *mockMappers
		in               comb.Input
		expectedInResult comb.Result
	}{
		{
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Anchor",
			m: &mockMappers{
				ToAnchorMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`$`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Group",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToGroupMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
					{},
				},
			},
			in:               newStringInput(`(a)`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_Match",
			m: &mockMappers{
				ToCharClassMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{},
				},
			},
			in:               newStringInput(`\d`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_UnescapedChar",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`ab`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`a`),
			expectedInResult: comb.Result{
				Val: comb.List{
					{},
					{Val: comb.Empty{}},
				},
			},
		},
		{
			name: "Success_Alt",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
					{},
				},
			},
			in: newStringInput(`a|b`),
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
			name:             "Failure",
			m:                &mockMappers{},
			in:               newStringInput(`\`),
			expectedInResult: comb.Result{},
		},
		{
			name: "Success_WithoutStartOfString",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRegexMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`a`),
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
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
				},
				ToRegexMocks: []MapFuncMock{
					{},
				},
			},
			in: newStringInput(`^a`),
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
		name          string
		m             *mockMappers
		regex         string
		expectedOut   *comb.Output
		expectedError string
	}{
		{
			name: "Success",
			m: &mockMappers{
				ToSingleCharMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToMatchMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprItemMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToSubexprMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToExprMocks: []MapFuncMock{
					{OutError: nil},
					{OutError: nil},
				},
				ToRegexMocks: []MapFuncMock{
					{OutError: nil},
				},
			},
			regex: `a`,
			expectedOut: &comb.Output{
				Result: comb.Result{
					Val: nil,
					Pos: 0,
					Bag: nil,
				},
				Remaining: nil,
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.m)
			out, err := p.Parse(tc.regex)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOut, out)
			} else {
				assert.Nil(t, out)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

//==================================================< MAPPERS >==================================================

type MapFuncMock struct {
	InResult  comb.Result
	OutResult comb.Result
	OutError  error
}

// mockMappers implements the Mapper interface for testing purposes.
type mockMappers struct {
	ToAnyCharIndex int
	ToAnyCharMocks []MapFuncMock

	ToSingleCharIndex int
	ToSingleCharMocks []MapFuncMock

	ToCharClassIndex int
	ToCharClassMocks []MapFuncMock

	ToASCIICharClassIndex int
	ToASCIICharClassMocks []MapFuncMock

	ToUnicodeCategoryIndex int
	ToUnicodeCategoryMocks []MapFuncMock

	ToUnicodeCharClassIndex int
	ToUnicodeCharClassMocks []MapFuncMock

	ToRepOpIndex int
	ToRepOpMocks []MapFuncMock

	ToUpperBoundIndex int
	ToUpperBoundMocks []MapFuncMock

	ToRangeIndex int
	ToRangeMocks []MapFuncMock

	ToRepetitionIndex int
	ToRepetitionMocks []MapFuncMock

	ToQuantifierIndex int
	ToQuantifierMocks []MapFuncMock

	ToCharInRangeIndex int
	ToCharInRangeMocks []MapFuncMock

	ToCharRangeIndex int
	ToCharRangeMocks []MapFuncMock

	ToCharGroupItemIndex int
	ToCharGroupItemMocks []MapFuncMock

	ToCharGroupIndex int
	ToCharGroupMocks []MapFuncMock

	ToMatchItemIndex int
	ToMatchItemMocks []MapFuncMock

	ToMatchIndex int
	ToMatchMocks []MapFuncMock

	ToGroupIndex int
	ToGroupMocks []MapFuncMock

	ToAnchorIndex int
	ToAnchorMocks []MapFuncMock

	ToSubexprItemIndex int
	ToSubexprItemMocks []MapFuncMock

	ToSubexprIndex int
	ToSubexprMocks []MapFuncMock

	ToExprIndex int
	ToExprMocks []MapFuncMock

	ToRegexIndex int
	ToRegexMocks []MapFuncMock
}

func (m *mockMappers) ToAnyChar(r comb.Result) (comb.Result, error) {
	i := m.ToAnyCharIndex
	m.ToAnyCharIndex++
	m.ToAnyCharMocks[i].InResult = r
	return m.ToAnyCharMocks[i].OutResult, m.ToAnyCharMocks[i].OutError
}

func (m *mockMappers) ToSingleChar(r comb.Result) (comb.Result, error) {
	i := m.ToSingleCharIndex
	m.ToSingleCharIndex++
	m.ToSingleCharMocks[i].InResult = r
	return m.ToSingleCharMocks[i].OutResult, m.ToSingleCharMocks[i].OutError
}

func (m *mockMappers) ToCharClass(r comb.Result) (comb.Result, error) {
	i := m.ToCharClassIndex
	m.ToCharClassIndex++
	m.ToCharClassMocks[i].InResult = r
	return m.ToCharClassMocks[i].OutResult, m.ToCharClassMocks[i].OutError
}

func (m *mockMappers) ToASCIICharClass(r comb.Result) (comb.Result, error) {
	i := m.ToASCIICharClassIndex
	m.ToASCIICharClassIndex++
	m.ToASCIICharClassMocks[i].InResult = r
	return m.ToASCIICharClassMocks[i].OutResult, m.ToASCIICharClassMocks[i].OutError
}

func (m *mockMappers) ToUnicodeCategory(r comb.Result) (comb.Result, error) {
	i := m.ToUnicodeCategoryIndex
	m.ToUnicodeCategoryIndex++
	m.ToUnicodeCategoryMocks[i].InResult = r
	return m.ToUnicodeCategoryMocks[i].OutResult, m.ToUnicodeCategoryMocks[i].OutError
}

func (m *mockMappers) ToUnicodeCharClass(r comb.Result) (comb.Result, error) {
	i := m.ToUnicodeCharClassIndex
	m.ToUnicodeCharClassIndex++
	m.ToUnicodeCharClassMocks[i].InResult = r
	return m.ToUnicodeCharClassMocks[i].OutResult, m.ToUnicodeCharClassMocks[i].OutError
}

func (m *mockMappers) ToRepOp(r comb.Result) (comb.Result, error) {
	i := m.ToRepOpIndex
	m.ToRepOpIndex++
	m.ToRepOpMocks[i].InResult = r
	return m.ToRepOpMocks[i].OutResult, m.ToRepOpMocks[i].OutError
}

func (m *mockMappers) ToUpperBound(r comb.Result) (comb.Result, error) {
	i := m.ToUpperBoundIndex
	m.ToUpperBoundIndex++
	m.ToUpperBoundMocks[i].InResult = r
	return m.ToUpperBoundMocks[i].OutResult, m.ToUpperBoundMocks[i].OutError
}

func (m *mockMappers) ToRange(r comb.Result) (comb.Result, error) {
	i := m.ToRangeIndex
	m.ToRangeIndex++
	m.ToRangeMocks[i].InResult = r
	return m.ToRangeMocks[i].OutResult, m.ToRangeMocks[i].OutError
}

func (m *mockMappers) ToRepetition(r comb.Result) (comb.Result, error) {
	i := m.ToRepetitionIndex
	m.ToRepetitionIndex++
	m.ToRepetitionMocks[i].InResult = r
	return m.ToRepetitionMocks[i].OutResult, m.ToRepetitionMocks[i].OutError
}

func (m *mockMappers) ToQuantifier(r comb.Result) (comb.Result, error) {
	i := m.ToQuantifierIndex
	m.ToQuantifierIndex++
	m.ToQuantifierMocks[i].InResult = r
	return m.ToQuantifierMocks[i].OutResult, m.ToQuantifierMocks[i].OutError
}

func (m *mockMappers) ToCharInRange(r comb.Result) (comb.Result, error) {
	i := m.ToCharInRangeIndex
	m.ToCharInRangeIndex++
	m.ToCharInRangeMocks[i].InResult = r
	return m.ToCharInRangeMocks[i].OutResult, m.ToCharInRangeMocks[i].OutError
}

func (m *mockMappers) ToCharRange(r comb.Result) (comb.Result, error) {
	i := m.ToCharRangeIndex
	m.ToCharRangeIndex++
	m.ToCharRangeMocks[i].InResult = r
	return m.ToCharRangeMocks[i].OutResult, m.ToCharRangeMocks[i].OutError
}

func (m *mockMappers) ToCharGroupItem(r comb.Result) (comb.Result, error) {
	i := m.ToCharGroupItemIndex
	m.ToCharGroupItemIndex++
	m.ToCharGroupItemMocks[i].InResult = r
	return m.ToCharGroupItemMocks[i].OutResult, m.ToCharGroupItemMocks[i].OutError
}

func (m *mockMappers) ToCharGroup(r comb.Result) (comb.Result, error) {
	i := m.ToCharGroupIndex
	m.ToCharGroupIndex++
	m.ToCharGroupMocks[i].InResult = r
	return m.ToCharGroupMocks[i].OutResult, m.ToCharGroupMocks[i].OutError
}

func (m *mockMappers) ToMatchItem(r comb.Result) (comb.Result, error) {
	i := m.ToMatchItemIndex
	m.ToMatchItemIndex++
	m.ToMatchItemMocks[i].InResult = r
	return m.ToMatchItemMocks[i].OutResult, m.ToMatchItemMocks[i].OutError
}

func (m *mockMappers) ToMatch(r comb.Result) (comb.Result, error) {
	i := m.ToMatchIndex
	m.ToMatchIndex++
	m.ToMatchMocks[i].InResult = r
	return m.ToMatchMocks[i].OutResult, m.ToMatchMocks[i].OutError
}

func (m *mockMappers) ToGroup(r comb.Result) (comb.Result, error) {
	i := m.ToGroupIndex
	m.ToGroupIndex++
	m.ToGroupMocks[i].InResult = r
	return m.ToGroupMocks[i].OutResult, m.ToGroupMocks[i].OutError
}

func (m *mockMappers) ToAnchor(r comb.Result) (comb.Result, error) {
	i := m.ToAnchorIndex
	m.ToAnchorIndex++
	m.ToAnchorMocks[i].InResult = r
	return m.ToAnchorMocks[i].OutResult, m.ToAnchorMocks[i].OutError
}

func (m *mockMappers) ToSubexprItem(r comb.Result) (comb.Result, error) {
	i := m.ToSubexprItemIndex
	m.ToSubexprItemIndex++
	m.ToSubexprItemMocks[i].InResult = r
	return m.ToSubexprItemMocks[i].OutResult, m.ToSubexprItemMocks[i].OutError
}

func (m *mockMappers) ToSubexpr(r comb.Result) (comb.Result, error) {
	i := m.ToSubexprIndex
	m.ToSubexprIndex++
	m.ToSubexprMocks[i].InResult = r
	return m.ToSubexprMocks[i].OutResult, m.ToSubexprMocks[i].OutError
}

func (m *mockMappers) ToExpr(r comb.Result) (comb.Result, error) {
	i := m.ToExprIndex
	m.ToExprIndex++
	m.ToExprMocks[i].InResult = r
	return m.ToExprMocks[i].OutResult, m.ToExprMocks[i].OutError
}

func (m *mockMappers) ToRegex(r comb.Result) (comb.Result, error) {
	i := m.ToRegexIndex
	m.ToRegexIndex++
	m.ToRegexMocks[i].InResult = r
	return m.ToRegexMocks[i].OutResult, m.ToRegexMocks[i].OutError
}
