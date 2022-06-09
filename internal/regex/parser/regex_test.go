package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// nopConvs implements the converters interface for testing purposes.
type nopConvs struct{}

func (c *nopConvs) ToChar(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToNum(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToLetters(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToRepOp(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToUpperBound(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToRange(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToRepetition(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToQuantifier(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToCharRange(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToCharGroupItem(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToCharGroup(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToCharClass(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToASCIICharClass(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToAnyChar(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToMatchItem(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToMatch(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToBackref(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToAnchor(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToGroup(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToSubexprItem(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToSubexpr(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToExpr(res result) (any, bool) {
	return res.Val, true
}

func (c *nopConvs) ToRegex(res result) (any, bool) {
	return res.Val, true
}

func TestRegex(t *testing.T) {
	c := new(nopConvs)
	r := NewRegex(c)

	tests := []struct {
		name        string
		p           parser
		in          input
		expectedOK  bool
		expectedOut any
	}{
		{
			name:       "char_Successful",
			p:          r.char,
			in:         newStringInput(`!"#$%&'()*+,-./[\]^_{|}~`),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '!', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune(`"#$%&'()*+,-./[\]^_{|}~`),
				},
			},
		},
		{
			name:       "digit_Successful",
			p:          r.digit,
			in:         newStringInput("0123456789"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '0', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("123456789"),
				},
			},
		},
		{
			name:       "letter_Successful",
			p:          r.letter,
			in:         newStringInput("abcdefghijklmnopqrstuvwxyz"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: 'a', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("bcdefghijklmnopqrstuvwxyz"),
				},
			},
		},
		{
			name:       "num_Successful",
			p:          r.num,
			in:         newStringInput("2022tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '2', Pos: 0},
						result{Val: '0', Pos: 1},
						result{Val: '2', Pos: 2},
						result{Val: '2', Pos: 3},
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
			name:       "letters_Successful",
			p:          r.letters,
			in:         newStringInput("head2022"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: 'h', Pos: 0},
						result{Val: 'e', Pos: 1},
						result{Val: 'a', Pos: 2},
						result{Val: 'd', Pos: 3},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("2022"),
				},
			},
		},
		{
			name:       "repOp_ZeroOrOne_Successful",
			p:          r.repOp,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '?', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repOp_ZeroOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '*', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repOp_OneOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '+', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "upperBound_Unbounded_Successful",
			p:          r.upperBound,
			in:         newStringInput(",}"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: ',', Pos: 0},
						result{Val: empty{}},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("}"),
				},
			},
		},
		{
			name:       "upperBound_Bounded_Successful",
			p:          r.upperBound,
			in:         newStringInput(",4}"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: ',', Pos: 0},
						result{
							Val: list{
								result{Val: '4', Pos: 1},
							},
							Pos: 1,
						},
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
			p:          r.range_,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{Val: empty{}},
						result{Val: '}', Pos: 2},
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
			name:       "range_upper_Unbounded_Successful",
			p:          r.range_,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{
							Val: list{
								result{Val: ',', Pos: 2},
								result{Val: empty{}},
							},
							Pos: 2,
						},
						result{Val: '}', Pos: 3},
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
			name:       "range_upper_Bounded_Successful",
			p:          r.range_,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{
							Val: list{
								result{Val: ',', Pos: 2},
								result{
									Val: list{
										result{Val: '4', Pos: 3},
									},
									Pos: 3,
								},
							},
							Pos: 2,
						},
						result{Val: '}', Pos: 4},
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
			name:       "repetition_ZeroOrOne_Successful",
			p:          r.repetition,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '?',
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_ZeroOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '*',
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_OneOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '+',
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_range_Fixed_Successful",
			p:          r.repetition,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{Val: empty{}},
						result{Val: '}', Pos: 2},
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
			name:       "repetition_range_upper_Unbounded_Successful",
			p:          r.repetition,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{
							Val: list{
								result{Val: ',', Pos: 2},
								result{Val: empty{}},
							},
							Pos: 2,
						},
						result{Val: '}', Pos: 3},
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
			name:       "repetition_range_upper_Bounded_Successful",
			p:          r.repetition,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '{', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
						},
						result{
							Val: list{
								result{Val: ',', Pos: 2},
								result{
									Val: list{
										result{Val: '4', Pos: 3},
									},
									Pos: 3,
								},
							},
							Pos: 2,
						},
						result{Val: '}', Pos: 4},
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
			name:       "quantifier_ZeroOrOne_Successful",
			p:          r.quantifier,
			in:         newStringInput("??tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '?', Pos: 0},
						result{Val: '?', Pos: 1},
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
			name:       "quantifier_ZeroOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("*?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '*', Pos: 0},
						result{Val: '?', Pos: 1},
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
			name:       "quantifier_OneOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("+?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '+', Pos: 0},
						result{Val: '?', Pos: 1},
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
			p:          r.quantifier,
			in:         newStringInput("{2}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '{', Pos: 0},
								result{
									Val: list{
										result{Val: '2', Pos: 1},
									},
									Pos: 1,
								},
								result{Val: empty{}},
								result{Val: '}', Pos: 2},
							},
							Pos: 0,
						},
						result{Val: '?', Pos: 3},
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
			name:       "quantifier_range_upper_Unbounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '{', Pos: 0},
								result{
									Val: list{
										result{Val: '2', Pos: 1},
									},
									Pos: 1,
								},
								result{
									Val: list{
										result{Val: ',', Pos: 2},
										result{Val: empty{}},
									},
									Pos: 2,
								},
								result{Val: '}', Pos: 3},
							},
							Pos: 0,
						},
						result{Val: '?', Pos: 4},
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
			name:       "quantifier_range_upper_Bounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,4}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '{', Pos: 0},
								result{
									Val: list{
										result{Val: '2', Pos: 1},
									},
									Pos: 1,
								},
								result{
									Val: list{
										result{Val: ',', Pos: 2},
										result{
											Val: list{
												result{Val: '4', Pos: 3},
											},
											Pos: 3,
										},
									},
									Pos: 2,
								},
								result{Val: '}', Pos: 4},
							},
							Pos: 0,
						},
						result{Val: '?', Pos: 5},
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
			p:          r.charRange,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '0', Pos: 0},
						result{Val: '-', Pos: 1},
						result{Val: '9', Pos: 2},
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
			p:          r.charClass,
			in:         newStringInput(`\dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\d`,
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
			p:          r.charClass,
			in:         newStringInput(`\Dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\D`,
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
			p:          r.charClass,
			in:         newStringInput(`\stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\s`,
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
			p:          r.charClass,
			in:         newStringInput(`\Stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\S`,
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
			p:          r.charClass,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\w`,
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
			p:          r.charClass,
			in:         newStringInput(`\Wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\W`,
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:blank:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:blank:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:space:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:space:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:digit:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:digit:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:xdigit:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:xdigit:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:upper:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:upper:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:lower:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:lower:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:alpha:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:alpha:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:alnum:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:alnum:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:word:]",
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
			p:          r.asciiCharClass,
			in:         newStringInput("[:ascii:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:ascii:]",
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
			p:          r.charGroupItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\w`,
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
			p:          r.charGroupItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:word:]",
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
			p:          r.charGroupItem,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '0', Pos: 0},
						result{Val: '-', Pos: 1},
						result{Val: '9', Pos: 2},
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
			p:          r.charGroupItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '!',
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
			p:          r.charGroup,
			in:         newStringInput(`[\w]tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: empty{}},
						result{
							Val: list{
								result{Val: `\w`, Pos: 1},
							},
							Pos: 1,
						},
						result{Val: ']', Pos: 3},
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
			p:          r.charGroup,
			in:         newStringInput("[[:word:]]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: empty{}},
						result{
							Val: list{
								result{Val: "[:word:]", Pos: 1},
							},
							Pos: 1,
						},
						result{Val: ']', Pos: 9},
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
			p:          r.charGroup,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: empty{}},
						result{
							Val: list{
								result{
									Val: list{
										result{Val: '0', Pos: 1},
										result{Val: '-', Pos: 2},
										result{Val: '9', Pos: 3},
									},
									Pos: 1,
								},
							},
							Pos: 1,
						},
						result{Val: ']', Pos: 4},
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
			p:          r.charGroup,
			in:         newStringInput("[!]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: empty{}},
						result{
							Val: list{
								result{Val: '!', Pos: 1},
							},
							Pos: 1,
						},
						result{Val: ']', Pos: 2},
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
			p:          r.charGroup,
			in:         newStringInput("[^#$]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: '^', Pos: 1},
						result{
							Val: list{
								result{Val: '#', Pos: 2},
								result{Val: '$', Pos: 3},
							},
							Pos: 2,
						},
						result{Val: ']', Pos: 4},
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
			p:          r.anyChar,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '.',
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
			p:          r.matchItem,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '.',
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charClass_Successful",
			p:          r.matchItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: `\w`,
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_asciiCharClass_Successful",
			p:          r.matchItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "[:word:]",
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charGroup_charRange_Successful",
			p:          r.matchItem,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '[', Pos: 0},
						result{Val: empty{}},
						result{
							Val: list{
								result{
									Val: list{
										result{Val: '0', Pos: 1},
										result{Val: '-', Pos: 2},
										result{Val: '9', Pos: 3},
									},
									Pos: 1,
								},
							},
							Pos: 1,
						},
						result{Val: ']', Pos: 4},
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
			name:       "matchItem_char_Successful",
			p:          r.matchItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '!',
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_anyChar_Successful",
			p:          r.match,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '.', Pos: 0},
						result{Val: empty{}},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charClass_Successful",
			p:          r.match,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: `\w`, Pos: 0},
						result{Val: empty{}},
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
			name:       "match_asciiCharClass_Successful",
			p:          r.match,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: "[:word:]", Pos: 0},
						result{Val: empty{}},
					},
					Pos: 0,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charGroup_charRange_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("[0-9]{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '[', Pos: 0},
								result{Val: empty{}},
								result{
									Val: list{
										result{
											Val: list{
												result{Val: '0', Pos: 1},
												result{Val: '-', Pos: 2},
												result{Val: '9', Pos: 3},
											},
											Pos: 1,
										},
									},
									Pos: 1,
								},
								result{Val: ']', Pos: 4},
							},
							Pos: 0,
						},
						result{
							Val: list{
								result{
									Val: list{
										result{Val: '{', Pos: 5},
										result{
											Val: list{
												result{Val: '2', Pos: 6},
											},
											Pos: 6,
										},
										result{
											Val: list{
												result{Val: ',', Pos: 7},
												result{
													Val: list{
														result{Val: '4', Pos: 8},
													},
													Pos: 8,
												},
											},
											Pos: 7,
										},
										result{Val: '}', Pos: 9},
									},
									Pos: 5,
								},
								result{Val: empty{}},
							},
							Pos: 5,
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
			name:       "match_char_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("#{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '#', Pos: 0},
						result{
							Val: list{
								result{
									Val: list{
										result{Val: '{', Pos: 1},
										result{
											Val: list{
												result{Val: '2', Pos: 2},
											},
											Pos: 2,
										},
										result{
											Val: list{
												result{Val: ',', Pos: 3},
												result{
													Val: list{
														result{Val: '4', Pos: 4},
													},
													Pos: 4,
												},
											},
											Pos: 3,
										},
										result{Val: '}', Pos: 5},
									},
									Pos: 1,
								},
								result{Val: empty{}},
							},
							Pos: 1,
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
			name:       "group_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{result{Val: '(', Pos: 0},
						result{
							Val: list{
								result{
									Val: list{
										result{
											Val: list{
												result{Val: 'a', Pos: 1},
												result{Val: empty{}},
											},
											Pos: 1,
										},
									},
									Pos: 1,
								},
								result{
									Val: list{
										result{Val: '|', Pos: 2},
										result{
											Val: list{
												result{
													Val: list{
														result{
															Val: list{
																result{Val: 'b', Pos: 3},
																result{Val: empty{}},
															},
															Pos: 3,
														},
													},
													Pos: 3,
												},
												result{Val: empty{}},
											},
											Pos: 3,
										},
									},
									Pos: 2,
								},
							},
							Pos: 1,
						},
						result{Val: ')', Pos: 4},
						result{Val: empty{}},
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
			name:       "group_quantifier_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '(', Pos: 0},
						result{
							Val: list{
								result{
									Val: list{
										result{
											Val: list{
												result{Val: 'a', Pos: 1},
												result{Val: empty{}},
											},
											Pos: 1,
										},
									},
									Pos: 1,
								},
								result{
									Val: list{
										result{Val: '|', Pos: 2},
										result{
											Val: list{
												result{
													Val: list{
														result{
															Val: list{
																result{Val: 'b', Pos: 3},
																result{Val: empty{}},
															},
															Pos: 3,
														},
													},
													Pos: 3,
												},
												result{Val: empty{}},
											},
											Pos: 3,
										},
									},
									Pos: 2,
								},
							},
							Pos: 1,
						},
						result{Val: ')', Pos: 4},
						result{
							Val: list{
								result{Val: '+', Pos: 5},
								result{Val: empty{}},
							},
							Pos: 5,
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
			name:       "backref_Successful",
			p:          r.backref,
			in:         newStringInput(`\2tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '\\', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
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
			name:       "anchor_Successful",
			p:          r.anchor,
			in:         newStringInput("$tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{Val: '$', Pos: 0},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_group_Successful",
			p:          r.subexprItem,
			in:         newStringInput("(ab)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '(', Pos: 0},
						result{
							Val: list{
								result{
									Val: list{
										result{
											Val: list{
												result{Val: 'a', Pos: 1},
												result{Val: empty{}},
											},
											Pos: 1,
										},
										result{
											Val: list{
												result{Val: 'b', Pos: 2},
												result{Val: empty{}},
											},
											Pos: 2,
										},
									},
									Pos: 1,
								},
								result{Val: empty{}},
							},
							Pos: 1,
						},
						result{Val: ')', Pos: 3},
						result{
							Val: list{
								result{Val: '+', Pos: 4},
								result{Val: empty{}},
							},
							Pos: 4,
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
			name:       "subexprItem_anchor_Successful",
			p:          r.subexprItem,
			in:         newStringInput("$"),
			expectedOK: true,
			expectedOut: output{
				Result:    result{Val: '$', Pos: 0},
				Remaining: nil,
			},
		},
		{
			name:       "subexprItem_backref_Successful",
			p:          r.subexprItem,
			in:         newStringInput(`\2tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '\\', Pos: 0},
						result{
							Val: list{
								result{Val: '2', Pos: 1},
							},
							Pos: 1,
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
			p:          r.subexprItem,
			in:         newStringInput("[0-9]+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '[', Pos: 0},
								result{Val: empty{}},
								result{
									Val: list{
										result{
											Val: list{
												result{Val: '0', Pos: 1},
												result{Val: '-', Pos: 2},
												result{Val: '9', Pos: 3},
											},
											Pos: 1,
										},
									},
									Pos: 1,
								},
								result{Val: ']', Pos: 4},
							},
							Pos: 0,
						},
						result{
							Val: list{
								result{Val: '+', Pos: 5},
								result{Val: empty{}},
							},
							Pos: 5,
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
			p:          r.subexpr,
			in:         newStringInput("(ab)+[0-9]*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{Val: '(', Pos: 0},
								result{
									Val: list{
										result{
											Val: list{
												result{
													Val: list{
														result{Val: 'a', Pos: 1},
														result{Val: empty{}},
													},
													Pos: 1,
												},
												result{
													Val: list{
														result{Val: 'b', Pos: 2},
														result{Val: empty{}},
													},
													Pos: 2,
												},
											},
											Pos: 1,
										},
										result{Val: empty{}},
									},
									Pos: 1,
								},
								result{Val: ')', Pos: 3},
								result{
									Val: list{
										result{Val: '+', Pos: 4},
										result{Val: empty{}},
									},
									Pos: 4,
								},
							},
							Pos: 0,
						},
						result{
							Val: list{
								result{
									Val: list{
										result{Val: '[', Pos: 5},
										result{Val: empty{}},
										result{
											Val: list{
												result{
													Val: list{
														result{Val: '0', Pos: 6},
														result{Val: '-', Pos: 7},
														result{Val: '9', Pos: 8},
													},
													Pos: 6,
												},
											},
											Pos: 6,
										},
										result{Val: ']', Pos: 9},
									},
									Pos: 5,
								},
								result{
									Val: list{
										result{Val: '*', Pos: 10},
										result{Val: empty{}},
									},
									Pos: 10,
								},
							},
							Pos: 5,
						},
						result{
							Val: list{
								result{Val: 't', Pos: 11},
								result{Val: empty{}},
							},
							Pos: 11,
						},
						result{
							Val: list{
								result{Val: 'a', Pos: 12},
								result{Val: empty{}},
							},
							Pos: 12,
						},
						result{
							Val: list{
								result{Val: 'i', Pos: 13},
								result{Val: empty{}},
							},
							Pos: 13,
						},
						result{
							Val: list{
								result{Val: 'l', Pos: 14},
								result{Val: empty{}},
							},
							Pos: 14,
						},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "expr_Successful",
			p:          r.expr,
			in:         newStringInput(`[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{
							Val: list{
								result{
									Val: list{
										result{
											Val: list{
												result{Val: '[', Pos: 0},
												result{Val: empty{}},
												result{
													Val: list{
														result{
															Val: list{
																result{Val: '0', Pos: 1},
																result{Val: '-', Pos: 2},
																result{Val: '9', Pos: 3},
															},
															Pos: 1,
														},
														result{
															Val: list{
																result{Val: 'A', Pos: 4},
																result{Val: '-', Pos: 5},
																result{Val: 'Z', Pos: 6},
															},
															Pos: 4,
														},
														result{
															Val: list{
																result{Val: 'a', Pos: 7},
																result{Val: '-', Pos: 8},
																result{Val: 'z', Pos: 9},
															},
															Pos: 7,
														},
														result{Val: '_', Pos: 10},
													},
													Pos: 1,
												},
												result{Val: ']', Pos: 11},
											},
											Pos: 0,
										},
										result{Val: empty{}},
									},
									Pos: 0,
								},
								result{
									Val: list{
										result{
											Val: list{
												result{Val: '[', Pos: 12},
												result{Val: empty{}},
												result{
													Val: list{
														result{Val: `\d`, Pos: 13},
														result{Val: `\w`, Pos: 15},
													},
													Pos: 13,
												},
												result{Val: ']', Pos: 17},
											},
											Pos: 12,
										},
										result{
											Val: list{
												result{Val: '*', Pos: 18},
												result{Val: empty{}},
											},
											Pos: 18,
										},
									},
									Pos: 12,
								},
							},
							Pos: 0,
						},
						result{Val: empty{}},
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "regex_Successful",
			p:          r.regex,
			in:         newStringInput(`^package\s+[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: list{
						result{Val: '^', Pos: 0},
						result{
							Val: list{
								result{
									Val: list{
										result{
											Val: list{
												result{Val: 'p', Pos: 1},
												result{Val: empty{}},
											},
											Pos: 1,
										},
										result{
											Val: list{
												result{Val: 'a', Pos: 2},
												result{Val: empty{}},
											},
											Pos: 2,
										},
										result{
											Val: list{
												result{Val: 'c', Pos: 3},
												result{Val: empty{}},
											},
											Pos: 3,
										},
										result{
											Val: list{
												result{Val: 'k', Pos: 4},
												result{Val: empty{}},
											},
											Pos: 4,
										},
										result{
											Val: list{
												result{Val: 'a', Pos: 5},
												result{Val: empty{}},
											},
											Pos: 5,
										},
										result{
											Val: list{
												result{Val: 'g', Pos: 6},
												result{Val: empty{}},
											},
											Pos: 6,
										},
										result{
											Val: list{
												result{Val: 'e', Pos: 7},
												result{Val: empty{}},
											},
											Pos: 7,
										},
										result{
											Val: list{
												result{Val: `\s`, Pos: 8},
												result{
													Val: list{
														result{Val: '+', Pos: 10},
														result{Val: empty{}},
													},
													Pos: 10,
												},
											},
											Pos: 8,
										},
										result{
											Val: list{
												result{
													Val: list{
														result{Val: '[', Pos: 11},
														result{Val: empty{}},
														result{
															Val: list{
																result{
																	Val: list{
																		result{Val: '0', Pos: 12},
																		result{Val: '-', Pos: 13},
																		result{Val: '9', Pos: 14},
																	},
																	Pos: 12,
																},
																result{
																	Val: list{
																		result{Val: 'A', Pos: 15},
																		result{Val: '-', Pos: 16},
																		result{Val: 'Z', Pos: 17},
																	},
																	Pos: 15,
																},
																result{
																	Val: list{
																		result{Val: 'a', Pos: 18},
																		result{Val: '-', Pos: 19},
																		result{Val: 'z', Pos: 20},
																	},
																	Pos: 18,
																},
																result{Val: '_', Pos: 21},
															},
															Pos: 12,
														},
														result{Val: ']', Pos: 22},
													},
													Pos: 11,
												},
												result{Val: empty{}},
											},
											Pos: 11,
										},
										result{
											Val: list{
												result{
													Val: list{
														result{Val: '[', Pos: 23},
														result{Val: empty{}},
														result{
															Val: list{
																result{Val: "\\d", Pos: 24},
																result{Val: "\\w", Pos: 26},
															},
															Pos: 24,
														},
														result{Val: ']', Pos: 28},
													},
													Pos: 23,
												},
												result{
													Val: list{
														result{Val: '*', Pos: 29},
														result{Val: empty{}},
													},
													Pos: 29,
												},
											},
											Pos: 23,
										},
									},
									Pos: 1,
								},
								result{Val: empty{}},
							},
							Pos: 1,
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
