package nfa

import (
	"testing"

	auto "github.com/moorara/algo/automata"
	comb "github.com/moorara/algo/parser/combinator"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	nfas := createTestNFAs()

	tests := []struct {
		name          string
		regex         string
		expectedError string
		expectedNFA   *auto.NFA
	}{
		{
			name:          "InvalidRegex",
			regex:         "[",
			expectedError: "invalid regular expression",
		},
		{
			name:          "InvalidCharRange",
			regex:         "[9-0]",
			expectedError: "invalid character range 9-0",
		},
		{
			name:          "InvalidRepRange",
			regex:         "[0-9]{4,2}",
			expectedError: "invalid repetition range {4,2}",
		},
		{
			name:  "Success",
			regex: `^[A-Z]?[a-z][0-9A-Za-z]{1,}$`,
			expectedNFA: empty().Union(nfas["upper"]).Concat(
				nfas["lower"],
				nfas["alnum"].Concat(nfas["alnum"].Star()),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nfa, err := Parse(tc.regex)

			if tc.expectedError != "" {
				assert.Nil(t, nfa)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.True(t, nfa.Equals(tc.expectedNFA))
			}
		})
	}
}

func TestMappers_ToAnyChar(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: '.',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["ascii"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnyChar(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSingleChar(t *testing.T) {
	xNFA := auto.NewNFA(0, auto.States{1})
	xNFA.Add(0, 'x', auto.States{1})

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: 'x',
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: xNFA,
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharClass(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success_Digit",
			r: comb.Result{
				Val: `\d`,
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
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
				Val: nfas["notDigit"],
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
				Val: nfas["whitespace"],
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
				Val: nfas["notWhitespace"],
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
				Val: nfas["word"],
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
				Val: nfas["notWord"],
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToASCIICharClass(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success_Blank",
			r: comb.Result{
				Val: "[:blank:]",
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["blank"],
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
				Val: nfas["space"],
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
				Val: nfas["digit"],
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
				Val: nfas["xdigit"],
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
				Val: nfas["upper"],
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
				Val: nfas["lower"],
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
				Val: nfas["alpha"],
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
				Val: nfas["alnum"],
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
				Val: nfas["word"],
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
				Val: nfas["ascii"],
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

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
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
			expectedError: "invalid repetition range {6,2}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRange(tc.r)

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharRange(t *testing.T) {
	aTofNFA := auto.NewNFA(0, auto.States{1})
	aTofNFA.Add(0, 'a', auto.States{1})
	aTofNFA.Add(0, 'b', auto.States{1})
	aTofNFA.Add(0, 'c', auto.States{1})
	aTofNFA.Add(0, 'd', auto.States{1})
	aTofNFA.Add(0, 'e', auto.States{1})
	aTofNFA.Add(0, 'f', auto.States{1})

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
				Val: aTofNFA,
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
				Val: auto.NewNFA(0, auto.States{1}),
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: []rune{},
				},
			},
			expectedOK:    true,
			expectedError: "invalid character range f-a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharRange(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroupItem(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: nfas["digit"],
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroup(t *testing.T) {
	nfas := createTestNFAs()

	hyphenNFA := auto.NewNFA(0, auto.States{1})
	hyphenNFA.Add(0, '-', auto.States{1})

	uuidNFA := auto.NewNFA(0, auto.States{1})
	uuidNFA.Add(0, '-', auto.States{1})
	for r := '0'; r <= '9'; r++ {
		uuidNFA.Add(0, auto.Symbol(r), auto.States{1})
	}
	for r := 'A'; r <= 'F'; r++ {
		uuidNFA.Add(0, auto.Symbol(r), auto.States{1})
	}
	for r := 'a'; r <= 'f'; r++ {
		uuidNFA.Add(0, auto.Symbol(r), auto.States{1})
	}

	notAlnumNFA := auto.NewNFA(0, auto.States{1})
	for r := 0; r <= 47; r++ {
		notAlnumNFA.Add(0, auto.Symbol(r), auto.States{1})
	}
	for r := 58; r <= 64; r++ {
		notAlnumNFA.Add(0, auto.Symbol(r), auto.States{1})
	}
	for r := 91; r <= 96; r++ {
		notAlnumNFA.Add(0, auto.Symbol(r), auto.States{1})
	}
	for r := 123; r <= 127; r++ {
		notAlnumNFA.Add(0, auto.Symbol(r), auto.States{1})
	}

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
								Val: nfas["xdigit"],
								Pos: 3,
								Bag: comb.Bag{
									bagKeyChars: xdigitChars,
								},
							},
							{
								Val: hyphenNFA,
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
				Val: uuidNFA,
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
								Val: nfas["alnum"],
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
				Val: notAlnumNFA,
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroup(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatchItem(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: nfas["digit"],
				Pos: 2,
				Bag: comb.Bag{
					bagKeyChars: digitChars,
				},
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatch(t *testing.T) {
	xNFA := auto.NewNFA(0, auto.States{1})
	xNFA.Add(0, 'x', auto.States{1})

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA,
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: empty().Union(xNFA),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Star(),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA.Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA, xNFA.Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Concat(
					xNFA,
					empty().Union(xNFA),
					empty().Union(xNFA),
				),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: comb.Result{
				Val: comb.List{
					{
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA.Star()),
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToGroup(t *testing.T) {
	xNFA := auto.NewNFA(0, auto.States{1})
	xNFA.Add(0, 'x', auto.States{1})

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: '(', Pos: 2},
					{
						Val: xNFA,
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
				Val: xNFA,
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
						Val: xNFA,
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
				Val: empty().Union(xNFA),
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
						Val: xNFA,
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
				Val: xNFA.Star(),
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
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA.Star()),
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
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA),
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
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA, xNFA.Star()),
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
						Val: xNFA,
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
				Val: xNFA.Concat(
					xNFA,
					empty().Union(xNFA),
					empty().Union(xNFA),
				),
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
						Val: xNFA,
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
				Val: xNFA.Concat(xNFA.Star()),
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

			EqualResults(t, tc.expectedResult, res)
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexprItem(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: nfas["digit"],
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexprItem(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexpr(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: nfas["digit"],
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
				Val: nfas["digit"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexpr(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToExpr(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{
						Val: nfas["upper"],
						Pos: 2,
					},
					{Val: comb.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["upper"],
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithExpr",
			r: comb.Result{
				Val: comb.List{
					{
						Val: nfas["upper"],
						Pos: 2,
					},
					{
						Val: comb.List{
							{Val: '|', Pos: 3},
							{
								Val: nfas["lower"],
								Pos: 4,
							},
						},
					},
				},
				Pos: 2,
			},
			expectedResult: comb.Result{
				Val: nfas["upper"].Union(nfas["lower"]),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToExpr(tc.r)

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRegex(t *testing.T) {
	nfas := createTestNFAs()

	tests := []MapperTest{
		{
			name: "Success",
			r: comb.Result{
				Val: comb.List{
					{Val: comb.Empty{}},
					{
						Val: nfas["digit"],
						Pos: 0,
					},
				},
				Pos: 0,
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
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
						Val: nfas["digit"],
						Pos: 1,
					},
				},
				Pos: 0,
			},
			expectedResult: comb.Result{
				Val: nfas["digit"],
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

			EqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

//==================================================< HELPERS >==================================================

var (
	digitChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
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

	whitespaceChars = []rune{
		' ', '\t', '\n', '\r', '\f',
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

	wordChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'_',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	notWordChars = []rune{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		58, 59, 60, 61, 62, 63, 64,
		91, 92, 93, 94, 96,
		123, 124, 125, 126, 127,
	}

	blankChars = []rune{
		' ', '\t',
	}

	spaceChars = []rune{
		' ', '\t', '\n', '\r', '\f', '\v',
	}

	xdigitChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F',
		'a', 'b', 'c', 'd', 'e', 'f',
	}

	upperChars = []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}

	lowerChars = []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	alphaChars = []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	alnumChars = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
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

func createTestNFAs() map[string]*auto.NFA {
	digit := auto.NewNFA(0, auto.States{1})
	for _, r := range digitChars {
		digit.Add(0, auto.Symbol(r), auto.States{1})
	}

	notDigit := auto.NewNFA(0, auto.States{1})
	for _, r := range notDigitChars {
		notDigit.Add(0, auto.Symbol(r), auto.States{1})
	}

	whitespace := auto.NewNFA(0, auto.States{1})
	for _, r := range whitespaceChars {
		whitespace.Add(0, auto.Symbol(r), auto.States{1})
	}

	notWhitespace := auto.NewNFA(0, auto.States{1})
	for _, r := range notWhitespaceChars {
		notWhitespace.Add(0, auto.Symbol(r), auto.States{1})
	}

	word := auto.NewNFA(0, auto.States{1})
	for _, r := range wordChars {
		word.Add(0, auto.Symbol(r), auto.States{1})
	}

	notWord := auto.NewNFA(0, auto.States{1})
	for _, r := range notWordChars {
		notWord.Add(0, auto.Symbol(r), auto.States{1})
	}

	blank := auto.NewNFA(0, auto.States{1})
	for _, r := range blankChars {
		blank.Add(0, auto.Symbol(r), auto.States{1})
	}

	space := auto.NewNFA(0, auto.States{1})
	for _, r := range spaceChars {
		space.Add(0, auto.Symbol(r), auto.States{1})
	}

	xdigit := auto.NewNFA(0, auto.States{1})
	for _, r := range xdigitChars {
		xdigit.Add(0, auto.Symbol(r), auto.States{1})
	}

	upper := auto.NewNFA(0, auto.States{1})
	for _, r := range upperChars {
		upper.Add(0, auto.Symbol(r), auto.States{1})
	}

	lower := auto.NewNFA(0, auto.States{1})
	for _, r := range lowerChars {
		lower.Add(0, auto.Symbol(r), auto.States{1})
	}

	alpha := auto.NewNFA(0, auto.States{1})
	for _, r := range alphaChars {
		alpha.Add(0, auto.Symbol(r), auto.States{1})
	}

	alnum := auto.NewNFA(0, auto.States{1})
	for _, r := range alnumChars {
		alnum.Add(0, auto.Symbol(r), auto.States{1})
	}

	ascii := auto.NewNFA(0, auto.States{1})
	for _, r := range asciiChars {
		ascii.Add(0, auto.Symbol(r), auto.States{1})
	}

	return map[string]*auto.NFA{
		"digit":         digit,
		"notDigit":      notDigit,
		"whitespace":    whitespace,
		"notWhitespace": notWhitespace,
		"word":          word,
		"notWord":       notWord,
		"blank":         blank,
		"space":         space,
		"xdigit":        xdigit,
		"upper":         upper,
		"lower":         lower,
		"alpha":         alpha,
		"alnum":         alnum,
		"ascii":         ascii,
	}
}

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

func EqualResults(t *testing.T, expectedResult, res comb.Result) {
	expectedNFA, ok := expectedResult.Val.(*auto.NFA)
	if !ok {
		assert.Equal(t, expectedResult, res)
		return
	}

	nfa, ok := res.Val.(*auto.NFA)
	if !ok {
		assert.Equal(t, expectedResult, res)
		return
	}

	assert.True(t, nfa.Equals(expectedNFA))
	assert.Equal(t, expectedResult.Pos, res.Pos)
	assert.Equal(t, expectedResult.Bag, res.Bag)
}
