package compact

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/emerge/internal/regex/token"
)

func TestRegex(t *testing.T) {
	tests := []struct {
		name        string
		node        Regex
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Regex{
				Expr: Expr{
					Sub: Subexpr{
						Items: []SubexprItem{
							&Match{
								Item: &Char{0, 'a'},
							},
						},
					},
					Expr: &Expr{
						Sub: Subexpr{
							Items: []SubexprItem{
								&Match{
									Item: &Char{2, 'b'},
								},
							},
						},
					},
				},
			},
			expectedPos: 0,
		},
		{
			name: "WithStartOfString",
			node: Regex{
				SOS: true,
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
			expectedPos: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestExpr(t *testing.T) {
	tests := []struct {
		name        string
		node        Expr
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Expr{
				Sub: Subexpr{
					Items: []SubexprItem{
						&Match{
							Item: &Char{0, 'a'},
						},
					},
				},
				Expr: &Expr{
					Sub: Subexpr{
						Items: []SubexprItem{
							&Match{
								Item: &Char{2, 'b'},
							},
						},
					},
				},
			},
			expectedPos: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestSubexpr(t *testing.T) {
	tests := []struct {
		name        string
		node        Subexpr
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Subexpr{
				Items: []SubexprItem{
					&Match{
						Item: &Char{0, 'a'},
					},
					&Match{
						Item: &Char{1, 'b'},
					},
				},
			},
			expectedPos: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestAnchor(t *testing.T) {
	tests := []struct {
		name        string
		node        Anchor
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "EndOfString",
			node: Anchor{
				TokPos: 10,
			},
			expectedTag: token.END_OF_STRING,
			expectedPos: 10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implSubexprItem()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestBackref(t *testing.T) {
	tests := []struct {
		name        string
		node        Backref
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Backref{
				SlashPos: 0,
				Ref: Num{
					StartPos: 1,
					EndPos:   2,
					Val:      10,
				},
				Group: &Group{
					OpenPos:  0,
					ClosePos: 4,
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
			},
			expectedPos: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implSubexprItem()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestGroup(t *testing.T) {
	tests := []struct {
		name        string
		node        Group
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Group{
				OpenPos:  0,
				ClosePos: 4,
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
			expectedPos: 0,
		},
		{
			name: "Quantified",
			node: Group{
				OpenPos:  0,
				ClosePos: 4,
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
					Rep: &RepOp{
						TokPos: 5,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implSubexprItem()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name        string
		node        Match
		expectedPos token.Pos
	}{
		{
			name: "AnyChar",
			node: Match{
				Item: &AnyChar{
					TokPos: 2,
				},
			},
			expectedPos: 2,
		},
		{
			name: "CharClass",
			node: Match{
				Item: &CharClass{
					StartPos: 2,
					EndPos:   3,
					TokTag:   token.WORD,
				},
			},
			expectedPos: 2,
		},
		{
			name: "ASCIICharClass",
			node: Match{
				Item: &ASCIICharClass{
					StartPos: 2,
					EndPos:   10,
					TokTag:   token.ALNUM_CHARS,
				},
			},
			expectedPos: 2,
		},
		{
			name: "CharGroup",
			node: Match{
				Item: &CharGroup{
					OpenPos:  2,
					ClosePos: 6,
					Items: []CharGroupItem{
						&CharRange{
							Low: Char{3, 'a'},
							Up:  Char{5, 'z'},
						},
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "Char",
			node: Match{
				Item: &Char{2, 'a'},
			},
			expectedPos: 2,
		},
		{
			name: "QuantifiedAnyChar",
			node: Match{
				Item: &AnyChar{
					TokPos: 2,
				},
				Quant: &Quantifier{
					Rep: &RepOp{
						TokPos: 3,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "QuantifiedCharClass",
			node: Match{
				Item: &CharClass{
					StartPos: 2,
					EndPos:   3,
					TokTag:   token.WORD_CHARS,
				},
				Quant: &Quantifier{
					Rep: &RepOp{
						TokPos: 4,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "QuantifiedASCIICharClass",
			node: Match{
				Item: &ASCIICharClass{
					StartPos: 2,
					EndPos:   10,
					TokTag:   token.ALNUM_CHARS,
				},
				Quant: &Quantifier{
					Rep: &RepOp{
						TokPos: 11,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "QuantifiedCharGroup",
			node: Match{
				Item: &CharGroup{
					OpenPos:  2,
					ClosePos: 6,
					Items: []CharGroupItem{
						&CharRange{
							Low: Char{3, 'a'},
							Up:  Char{5, 'z'},
						},
					},
				},
				Quant: &Quantifier{
					Rep: &RepOp{
						TokPos: 7,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "QuantifiedChar",
			node: Match{
				Item: &Char{2, 'a'},
				Quant: &Quantifier{
					Rep: &RepOp{
						TokPos: 3,
						TokTag: token.ZERO_OR_MORE,
					},
				},
			},
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implSubexprItem()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestAnyChar(t *testing.T) {
	tests := []struct {
		name        string
		node        AnyChar
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: AnyChar{
				TokPos: 5,
			},
			expectedTag: token.ANY_CHAR,
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implMatchItem()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestCharClass(t *testing.T) {
	tests := []struct {
		name        string
		node        CharClass
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "DigitChars",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.DIGIT,
			},
			expectedTag: token.DIGIT,
			expectedPos: 2,
		},
		{
			name: "NonDigitChars",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.NON_DIGIT,
			},
			expectedTag: token.NON_DIGIT,
			expectedPos: 2,
		},
		{
			name: "Whitespace",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.WHITESPACE,
			},
			expectedTag: token.WHITESPACE,
			expectedPos: 2,
		},
		{
			name: "NonWhitespace",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.NON_WHITESPACE,
			},
			expectedTag: token.NON_WHITESPACE,
			expectedPos: 2,
		},
		{
			name: "WordChars",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.WORD_CHARS,
			},
			expectedTag: token.WORD_CHARS,
			expectedPos: 2,
		},
		{
			name: "NonWordChars",
			node: CharClass{
				StartPos: 2,
				EndPos:   3,
				TokTag:   token.NON_WORD,
			},
			expectedTag: token.NON_WORD,
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implMatchItem()
			tc.node.implCharGroupItem()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestASCIICharClass(t *testing.T) {
	tests := []struct {
		name        string
		node        ASCIICharClass
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "Blank",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.BLANK_CHARS,
			},
			expectedTag: token.BLANK_CHARS,
			expectedPos: 2,
		},
		{
			name: "Space",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.SPACE_CHARS,
			},
			expectedTag: token.SPACE_CHARS,
			expectedPos: 2,
		},
		{
			name: "Digit",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.DIGIT_CHARS,
			},
			expectedTag: token.DIGIT_CHARS,
			expectedPos: 2,
		},
		{
			name: "XDigit",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   11,
				TokTag:   token.XDIGIT_CHARS,
			},
			expectedTag: token.XDIGIT_CHARS,
			expectedPos: 2,
		},
		{
			name: "Upper",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.UPPER_CHARS,
			},
			expectedTag: token.UPPER_CHARS,
			expectedPos: 2,
		},
		{
			name: "Lower",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.LOWER_CHARS,
			},
			expectedTag: token.LOWER_CHARS,
			expectedPos: 2,
		},
		{
			name: "Alpha",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.ALPHA_CHARS,
			},
			expectedTag: token.ALPHA_CHARS,
			expectedPos: 2,
		},
		{
			name: "Alnum",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.ALNUM_CHARS,
			},
			expectedTag: token.ALNUM_CHARS,
			expectedPos: 2,
		},
		{
			name: "Word",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   9,
				TokTag:   token.WORD_CHARS,
			},
			expectedTag: token.WORD_CHARS,
			expectedPos: 2,
		},
		{
			name: "ASCII",
			node: ASCIICharClass{
				StartPos: 2,
				EndPos:   10,
				TokTag:   token.ASCII_CHARS,
			},
			expectedTag: token.ASCII_CHARS,
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implMatchItem()
			tc.node.implCharGroupItem()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestCharGroup(t *testing.T) {
	tests := []struct {
		name        string
		node        CharGroup
		expectedPos token.Pos
	}{
		{
			name: "CharClass",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 5,
				Items: []CharGroupItem{
					&CharClass{
						StartPos: 3,
						EndPos:   4,
						TokTag:   token.WORD,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "ASCIICharClass",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 12,
				Items: []CharGroupItem{
					&ASCIICharClass{
						StartPos: 3,
						EndPos:   11,
						TokTag:   token.ALNUM_CHARS,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "CharRange",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 6,
				Items: []CharGroupItem{
					&CharRange{
						Low: Char{3, 'a'},
						Up:  Char{5, 'z'},
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "Char",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 5,
				Items: []CharGroupItem{
					&Char{3, 'a'},
					&Char{4, 'b'},
				},
			},
			expectedPos: 2,
		},
		{
			name: "NegatedCharClass",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 6,
				Negated:  true,
				Items: []CharGroupItem{
					&CharClass{
						StartPos: 4,
						EndPos:   5,
						TokTag:   token.WORD,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "NegatedASCIICharClass",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 13,
				Negated:  true,
				Items: []CharGroupItem{
					&ASCIICharClass{
						StartPos: 4,
						EndPos:   12,
						TokTag:   token.ALNUM_CHARS,
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "NegatedCharRange",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 7,
				Negated:  true,
				Items: []CharGroupItem{
					&CharRange{
						Low: Char{4, 'a'},
						Up:  Char{6, 'z'},
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "NegatedChar",
			node: CharGroup{
				OpenPos:  2,
				ClosePos: 6,
				Negated:  true,
				Items: []CharGroupItem{
					&Char{4, 'a'},
					&Char{5, 'b'},
				},
			},
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implMatchItem()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestCharRange(t *testing.T) {
	tests := []struct {
		name        string
		node        CharRange
		expectedPos token.Pos
	}{
		{
			name: "Digits",
			node: CharRange{
				Low: Char{2, '0'},
				Up:  Char{4, '9'},
			},
			expectedPos: 2,
		},
		{
			name: "Letters",
			node: CharRange{
				Low: Char{2, 'a'},
				Up:  Char{4, 'z'},
			},
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implCharGroupItem()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestQuantifier(t *testing.T) {
	tests := []struct {
		name        string
		node        Quantifier
		expectedPos token.Pos
	}{
		{
			name: "ZeroOrOne",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ZERO_OR_ONE,
				},
			},
			expectedPos: 2,
		},
		{
			name: "ZeroOrMore",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ZERO_OR_MORE,
				},
			},
			expectedPos: 2,
		},
		{
			name: "OneOrMore",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ONE_OR_MORE,
				},
			},
			expectedPos: 2,
		},
		{
			name: "Range",
			node: Quantifier{
				Rep: &Range{
					OpenPos:  2,
					ClosePos: 6,
					Low: Num{
						StartPos: 3,
						EndPos:   3,
						Val:      2,
					},
					Up: &UpperBound{
						CommaPos: 4,
						Val: &Num{
							StartPos: 5,
							EndPos:   5,
							Val:      4,
						},
					},
				},
			},
			expectedPos: 2,
		},
		{
			name: "LazyZeroOrOne",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ZERO_OR_ONE,
				},
				Lazy: true,
			},
			expectedPos: 2,
		},
		{
			name: "LazyZeroOrMore",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ZERO_OR_MORE,
				},
				Lazy: true,
			},
			expectedPos: 2,
		},
		{
			name: "LazyOneOrMore",
			node: Quantifier{
				Rep: &RepOp{
					TokPos: 2,
					TokTag: token.ONE_OR_MORE,
				},
				Lazy: true,
			},
			expectedPos: 2,
		},
		{
			name: "LazyRange",
			node: Quantifier{
				Rep: &Range{
					OpenPos:  2,
					ClosePos: 6,
					Low: Num{
						StartPos: 3,
						EndPos:   3,
						Val:      2,
					},
					Up: &UpperBound{
						CommaPos: 4,
						Val: &Num{
							StartPos: 5,
							EndPos:   5,
							Val:      4,
						},
					},
				},
				Lazy: true,
			},
			expectedPos: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestRepOp(t *testing.T) {
	tests := []struct {
		name        string
		node        RepOp
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "ZeroOrOne",
			node: RepOp{
				TokPos: 5,
				TokTag: token.ZERO_OR_ONE,
			},
			expectedTag: token.ZERO_OR_ONE,
			expectedPos: 5,
		},
		{
			name: "ZeroOrMore",
			node: RepOp{
				TokPos: 5,
				TokTag: token.ZERO_OR_MORE,
			},
			expectedTag: token.ZERO_OR_MORE,
			expectedPos: 5,
		},
		{
			name: "OneOrMore",
			node: RepOp{
				TokPos: 5,
				TokTag: token.ONE_OR_MORE,
			},
			expectedTag: token.ONE_OR_MORE,
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implRepetition()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name        string
		node        Range
		expectedPos token.Pos
	}{
		{
			name: "Fixed",
			node: Range{
				OpenPos:  1,
				ClosePos: 4,
				Low: Num{
					StartPos: 2,
					EndPos:   3,
					Val:      10,
				},
			},
			expectedPos: 1,
		},
		{
			name: "Unbounded",
			node: Range{
				OpenPos:  1,
				ClosePos: 5,
				Low: Num{
					StartPos: 2,
					EndPos:   3,
					Val:      10,
				},
				Up: &UpperBound{
					CommaPos: 4,
				},
			},
			expectedPos: 1,
		},
		{
			name: "Bounded",
			node: Range{
				OpenPos:  1,
				ClosePos: 8,
				Low: Num{
					StartPos: 2,
					EndPos:   3,
					Val:      10,
				},
				Up: &UpperBound{
					CommaPos: 4,
					Val: &Num{
						StartPos: 6,
						EndPos:   7,
						Val:      20,
					},
				},
			},
			expectedPos: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implRepetition()

			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestUpperBound(t *testing.T) {
	tests := []struct {
		name        string
		node        UpperBound
		expectedPos token.Pos
	}{
		{
			name: "Unbounded",
			node: UpperBound{
				CommaPos: 5,
			},
			expectedPos: 5,
		},
		{
			name: "Bounded",
			node: UpperBound{
				CommaPos: 5,
				Val: &Num{
					StartPos: 7,
					EndPos:   8,
					Val:      10,
				},
			},
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestNum(t *testing.T) {
	tests := []struct {
		name        string
		node        Num
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Num{
				StartPos: 5,
				EndPos:   6,
				Val:      69,
			},
			expectedTag: token.NUM,
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestLetters(t *testing.T) {
	tests := []struct {
		name        string
		node        Letters
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name: "OK",
			node: Letters{
				StartPos: 5,
				EndPos:   15,
				Val:      "placeholder",
			},
			expectedTag: token.LETTERS,
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		name        string
		node        Char
		expectedTag token.Tag
		expectedPos token.Pos
	}{
		{
			name:        "OK",
			node:        Char{5, 'a'},
			expectedTag: token.CHAR,
			expectedPos: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.node.implMatchItem()
			tc.node.implCharGroupItem()

			assert.Equal(t, tc.expectedTag, tc.node.Tag())
			assert.Equal(t, tc.expectedPos, tc.node.Pos())
		})
	}
}
