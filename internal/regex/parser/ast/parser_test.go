package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	comb "github.com/gardenbed/emerge/internal/combinator"
)

func TestMappers_ToAnyChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: '.',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: ascii,
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnyChar(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSingleChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: 'x',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Char{Val: 'x'},
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: []rune{'x'},
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSingleChar(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Digit",
			r: comb.Result{
				Val: `\d`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotDigit",
			r: comb.Result{
				Val: `\D`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: notDigit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: notDigitChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Whitespace",
			r: comb.Result{
				Val: `\s`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: whitespace,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: whitespaceChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWhitespace",
			r: comb.Result{
				Val: `\S`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: notWhitespace,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: notWhitespaceChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: comb.Result{
				Val: `\w`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: word,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: wordChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWord",
			r: comb.Result{
				Val: `\W`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: notWord,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: notWordChars,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharClass(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToASCIICharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Blank",
			r: comb.Result{
				Val: "[:blank:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: blank,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: blankChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Space",
			r: comb.Result{
				Val: "[:space:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: space,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: spaceChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Digit",
			r: comb.Result{
				Val: "[:digit:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_XDigit",
			r: comb.Result{
				Val: "[:xdigit:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: xdigit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: xdigitChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Upper",
			r: comb.Result{
				Val: "[:upper:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: upper,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: upperChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Lower",
			r: comb.Result{
				Val: "[:lower:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: lower,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: lowerChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alpha",
			r: comb.Result{
				Val: "[:alpha:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: alpha,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: alphaChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alnum",
			r: comb.Result{
				Val: "[:alnum:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: alnum,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: alnumChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: comb.Result{
				Val: "[:word:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: word,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: wordChars,
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_ASCII",
			r: comb.Result{
				Val: "[:ascii:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: ascii,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: asciiChars,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToASCIICharClass(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepOp(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepOp(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUpperBound(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Unbounded",
			r: comb.Result{
				Val: comb.List{
					{Val: ',', Pos: 2},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: (*int)(nil),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: comb.Result{
				Val: comb.List{
					{Val: ',', Pos: 2},
					{Val: 4, Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: intPtr(4),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUpperBound(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Fixed",
			r: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{Val: comb.Empty{}},
					{Val: '}', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Unbounded",
			r: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: (*int)(nil),
						Pos: 4,
					},
					{Val: '}', Pos: 5},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: (*int)(nil),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: intPtr(6),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(6),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: comb.Result{
				Val: comb.List{
					{Val: '{', Pos: 2},
					{Val: 6, Pos: 3},
					{
						Val: intPtr(2),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[int, *int]{
					p: 6,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK:    true,
			expectedError: "1 error occurred:\n\t* invalid repetition range {6,2}\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepetition(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepetition(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToQuantifier(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_NonLazy",
			r: comb.Result{
				Val: comb.List{
					{Val: '*', Pos: 2},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: false,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy",
			r: comb.Result{
				Val: comb.List{
					{Val: '*', Pos: 2},
					{Val: '?', Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: true,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToQuantifier(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharInRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharInRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: 'a', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'f', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Char{Val: 'a'},
						&Char{Val: 'b'},
						&Char{Val: 'c'},
						&Char{Val: 'd'},
						&Char{Val: 'e'},
						&Char{Val: 'f'},
					},
				},
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: []rune{'a', 'b', 'c', 'd', 'e', 'f'},
				},
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: comb.Result{
				Val: comb.List{
					{Val: 'f', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'a', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{},
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: []rune{},
				},
			},
			expectedOK:    true,
			expectedError: "1 error occurred:\n\t* invalid character range f-a\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroupItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedResult: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroupItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: '[', Pos: 2},
					{Val: comb.Empty{}},
					{
						Val: comb.List{
							{
								Val: xdigit,
								Pos: 3,
								Bag: comb.Bag{
									bagKeyChars: xdigitChars,
								},
							},
							{
								Val: &Char{Val: '-'},
								Pos: 12,
								Bag: comb.Bag{
									bagKeyChars: []rune{'-'},
								},
							},
						},
						Pos: 3,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Char{Val: '-'},
						&Char{Val: '0'},
						&Char{Val: '1'},
						&Char{Val: '2'},
						&Char{Val: '3'},
						&Char{Val: '4'},
						&Char{Val: '5'},
						&Char{Val: '6'},
						&Char{Val: '7'},
						&Char{Val: '8'},
						&Char{Val: '9'},
						&Char{Val: 'A'},
						&Char{Val: 'B'},
						&Char{Val: 'C'},
						&Char{Val: 'D'},
						&Char{Val: 'E'},
						&Char{Val: 'F'},
						&Char{Val: 'a'},
						&Char{Val: 'b'},
						&Char{Val: 'c'},
						&Char{Val: 'd'},
						&Char{Val: 'e'},
						&Char{Val: 'f'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Negated",
			r: comb.Result{
				Val: comb.List{
					{Val: '[', Pos: 2},
					{Val: '^', Pos: 3},
					{
						Val: comb.List{
							{
								Val: alnum,
								Pos: 4,
								Bag: comb.Bag{
									bagKeyChars: alnumChars,
								},
							},
						},
						Pos: 4,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Char{Val: 0},
						&Char{Val: 1},
						&Char{Val: 2},
						&Char{Val: 3},
						&Char{Val: 4},
						&Char{Val: 5},
						&Char{Val: 6},
						&Char{Val: 7},
						&Char{Val: 8},
						&Char{Val: 9},
						&Char{Val: 10},
						&Char{Val: 11},
						&Char{Val: 12},
						&Char{Val: 13},
						&Char{Val: 14},
						&Char{Val: 15},
						&Char{Val: 16},
						&Char{Val: 17},
						&Char{Val: 18},
						&Char{Val: 19},
						&Char{Val: 20},
						&Char{Val: 21},
						&Char{Val: 22},
						&Char{Val: 23},
						&Char{Val: 24},
						&Char{Val: 25},
						&Char{Val: 26},
						&Char{Val: 27},
						&Char{Val: 28},
						&Char{Val: 29},
						&Char{Val: 30},
						&Char{Val: 31},
						&Char{Val: 32},
						&Char{Val: 33},
						&Char{Val: 34},
						&Char{Val: 35},
						&Char{Val: 36},
						&Char{Val: 37},
						&Char{Val: 38},
						&Char{Val: 39},
						&Char{Val: 40},
						&Char{Val: 41},
						&Char{Val: 42},
						&Char{Val: 43},
						&Char{Val: 44},
						&Char{Val: 45},
						&Char{Val: 46},
						&Char{Val: 47},
						&Char{Val: 58},
						&Char{Val: 59},
						&Char{Val: 60},
						&Char{Val: 61},
						&Char{Val: 62},
						&Char{Val: 63},
						&Char{Val: 64},
						&Char{Val: 91},
						&Char{Val: 92},
						&Char{Val: 93},
						&Char{Val: 94},
						&Char{Val: 95},
						&Char{Val: 96},
						&Char{Val: 123},
						&Char{Val: 124},
						&Char{Val: 125},
						&Char{Val: 126},
						&Char{Val: 127},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroup(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatchItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedResult: comb.Result{
				Val: digit,
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatchItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatch(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Char{Val: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Empty{},
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Star{
					Expr: &Char{Val: 'x'},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Val: 'x'},
							},
						},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Val: 'x'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
				Bag: comb.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatch(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Char{Val: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Empty{},
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Star{
					Expr: &Char{Val: 'x'},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Char{Val: 'x'},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Val: 'x'},
							},
						},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Val: 'x'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Val: 'x'},
						Pos: 3,
						Bag: comb.Bag{
							bagKeyChars: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
						&Star{
							Expr: &Char{Val: 'x'},
						},
					},
				},
				Pos: 2,
				Bag: comb.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToGroup(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToAnchor(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: '$',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: EndOfString,
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnchor(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexprItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: &Char{Val: 'x'},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Char{Val: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexprItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Char{Val: 'x'},
						Pos: 2,
					},
					{
						Val: EndOfString,
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexpr(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToExpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Val: 'x'},
							},
						},
						Pos: 2,
					},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithExpr",
			r: comb.Result{
				Val: comb.List{
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Val: 'x'},
							},
						},
						Pos: 2,
					},
					{
						Val: comb.List{
							{Val: '|', Pos: 3},
							{
								Val: &Concat{
									Exprs: []Node{
										&Char{Val: 'y'},
									},
								},
								Pos: 4,
							},
						},
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: &Alt{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Char{Val: 'x'},
							},
						},
						&Concat{
							Exprs: []Node{
								&Char{Val: 'y'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToExpr(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRegex(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: comb.Empty{}},
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Val: 's'},
							},
						},
						Pos: 0,
					},
				},
				Pos: 0,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 's'},
					},
				},
				Pos: 0,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithStartOfString",
			r: comb.Result{
				Val: comb.List{
					{Val: '^', Pos: 0},
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Val: 's'},
							},
						},
						Pos: 1,
					},
				},
				Pos: 0,
			},
			expectedResult: comb.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Val: 's'},
					},
				},
				Pos: 0,
				Bag: comb.Bag{
					BagKeyStartOfString: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRegex(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

//==================================================< HELPERS >==================================================

var (
	digit = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
		},
	}

	digitChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	notDigit = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	notDigitChars = []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		58, 59, 60, 61, 62, 63, 64,
		65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
		91, 92, 93, 94, 95, 96,
		97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
		123, 124, 125, 126, 127,
	}

	whitespace = &Alt{
		Exprs: []Node{
			&Char{Val: ' '},
			&Char{Val: '\t'},
			&Char{Val: '\n'},
			&Char{Val: '\r'},
			&Char{Val: '\f'},
		},
	}

	whitespaceChars = []rune{
		' ', '\t', '\n', '\r', '\f',
	}

	notWhitespace = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 11},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 48},
			&Char{Val: 49},
			&Char{Val: 50},
			&Char{Val: 51},
			&Char{Val: 52},
			&Char{Val: 53},
			&Char{Val: 54},
			&Char{Val: 55},
			&Char{Val: 56},
			&Char{Val: 57},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	notWhitespaceChars = []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8,
		11,
		14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
		58, 59, 60, 61, 62, 63, 64,
		65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
		91, 92, 93, 94, 95, 96,
		97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
		123, 124, 125, 126, 127,
	}

	word = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: '_'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	wordChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'_',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	notWord = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 96},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	notWordChars = []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		58, 59, 60, 61, 62, 63, 64,
		91, 92, 93, 94, 96,
		123, 124, 125, 126, 127,
	}

	blank = &Alt{
		Exprs: []Node{
			&Char{Val: ' '},
			&Char{Val: '\t'},
		},
	}

	blankChars = []rune{
		' ', '\t',
	}

	space = &Alt{
		Exprs: []Node{
			&Char{Val: ' '},
			&Char{Val: '\t'},
			&Char{Val: '\n'},
			&Char{Val: '\r'},
			&Char{Val: '\f'},
			&Char{Val: '\v'},
		},
	}

	spaceChars = []rune{
		' ', '\t', '\n', '\r', '\f', '\v',
	}

	xdigit = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
		},
	}

	xdigitChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F',
		'a', 'b', 'c', 'd', 'e', 'f',
	}

	upper = &Alt{
		Exprs: []Node{
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
		},
	}

	upperChars = []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}

	lower = &Alt{
		Exprs: []Node{
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	lowerChars = []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	alpha = &Alt{
		Exprs: []Node{
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	alphaChars = []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	alnum = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	alnumChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	ascii = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 48},
			&Char{Val: 49},
			&Char{Val: 50},
			&Char{Val: 51},
			&Char{Val: 52},
			&Char{Val: 53},
			&Char{Val: 54},
			&Char{Val: 55},
			&Char{Val: 56},
			&Char{Val: 57},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	asciiChars = []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
		58, 59, 60, 61, 62, 63, 64,
		65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
		91, 92, 93, 94, 95, 96,
		97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
		123, 124, 125, 126, 127,
	}
)

func intPtr(v int) *int {
	return &v
}

type MapperTest struct {
	name           string
	r              comb.Result
	expectedResult comb.Result
	expectedOK     bool
	expectedError  string
}
