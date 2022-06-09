package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ast "github.com/gardenbed/emerge/internal/regex/ast/compact"
	"github.com/gardenbed/emerge/internal/regex/token"
)

func TestParseCompact(t *testing.T) {
	tests := []struct {
		name           string
		in             input
		expectedError  string
		expectedResult *ast.Regex
	}{
		{
			name:          "InvalidBackref",
			in:            newStringInput(`\1`),
			expectedError: "1 error occurred:\n\t* invalid back reference \\1\n\n",
		},
		{
			name:          "InvalidCharRange",
			in:            newStringInput("[9-0]"),
			expectedError: "1 error occurred:\n\t* invalid character range 9-0\n\n",
		},
		{
			name:          "InvalidRepRange",
			in:            newStringInput("[0-9]{4,2}"),
			expectedError: "1 error occurred:\n\t* invalid repetition range {4,2}\n\n",
		},
		{
			name: "Successful",
			in:   newStringInput("[A-Z][0-9A-Za-z_]*"),
			expectedResult: &ast.Regex{
				SOS: false,
				Expr: ast.Expr{
					Sub: ast.Subexpr{
						Items: []ast.SubexprItem{
							&ast.Match{
								Item: &ast.CharGroup{
									OpenPos:  0,
									ClosePos: 4,
									Negated:  false,
									Items: []ast.CharGroupItem{
										&ast.CharRange{
											Low: ast.Char{
												TokPos: 1,
												Val:    'A',
											},
											Up: ast.Char{
												TokPos: 3,
												Val:    'Z',
											},
										},
									},
								},
							},
							&ast.Match{
								Item: &ast.CharGroup{
									OpenPos:  5,
									ClosePos: 16,
									Negated:  false,
									Items: []ast.CharGroupItem{
										&ast.CharRange{
											Low: ast.Char{
												TokPos: 6,
												Val:    '0',
											},
											Up: ast.Char{
												TokPos: 8,
												Val:    '9',
											},
										},
										&ast.CharRange{
											Low: ast.Char{
												TokPos: 9,
												Val:    'A',
											},
											Up: ast.Char{
												TokPos: 11,
												Val:    'Z',
											},
										},
										&ast.CharRange{
											Low: ast.Char{
												TokPos: 12,
												Val:    'a',
											},
											Up: ast.Char{
												TokPos: 14,
												Val:    'z',
											},
										},
										&ast.Char{
											TokPos: 15,
											Val:    '_',
										},
									},
								},
								Quant: &ast.Quantifier{
									Rep: &ast.RepOp{
										TokPos: 17,
										TokTag: token.ZERO_OR_MORE,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseCompact(tc.in)

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

func TestCompactConverters(t *testing.T) {
	c := new(compactConverters)
	r := NewRegex(c)

	tests := []struct {
		name        string
		p           parser
		in          input
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "char_Successful",
			p:          r.char,
			in:         newStringInput(`!"#$%&'()*+,-./[\]^_{|}~`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Char{
						TokPos: token.Pos(0),
						Val:    '!',
					},
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
			p:          r.digit,
			in:         newStringInput("0123456789"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '0',
					Pos: 0,
				},
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
				Result: result{
					Val: 'a',
					Pos: 0,
				},
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
					Val: ast.Num{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(3),
						Val:      2022,
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
					Val: ast.Letters{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(3),
						Val:      "head",
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
				Result: result{
					Val: ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ZERO_OR_ONE,
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
			name:       "repOp_ZeroOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ZERO_OR_MORE,
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
			name:       "repOp_OneOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ONE_OR_MORE,
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
			name:       "upperBound_Unbounded_Successful",
			p:          r.upperBound,
			in:         newStringInput(",}"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.UpperBound{
						CommaPos: 0,
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
					Val: ast.UpperBound{
						CommaPos: 1,
						Val: &ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      4,
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
					Val: ast.Range{
						OpenPos:  0,
						ClosePos: 2,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
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
			name:       "range_upper_Unbounded_Successful",
			p:          r.range_,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Range{
						OpenPos:  0,
						ClosePos: 3,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Up: &ast.UpperBound{},
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
					Val: ast.Range{
						OpenPos:  0,
						ClosePos: 4,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Up: &ast.UpperBound{
							CommaPos: 3,
							Val: &ast.Num{
								StartPos: token.Pos(3),
								EndPos:   token.Pos(3),
								Val:      4,
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
			name:       "repetition_ZeroOrOne_Successful",
			p:          r.repetition,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ZERO_OR_ONE,
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
			name:       "repetition_ZeroOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ZERO_OR_MORE,
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
			name:       "repetition_OneOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.RepOp{
						TokPos: token.Pos(0),
						TokTag: token.ONE_OR_MORE,
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
			name:       "repetition_range_Fixed_Successful",
			p:          r.repetition,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Range{
						OpenPos:  0,
						ClosePos: 2,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
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
			name:       "repetition_range_upper_Unbounded_Successful",
			p:          r.repetition,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Range{
						OpenPos:  0,
						ClosePos: 3,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Up: &ast.UpperBound{},
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
					Val: &ast.Range{
						OpenPos:  0,
						ClosePos: 4,
						Low: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Up: &ast.UpperBound{
							CommaPos: 3,
							Val: &ast.Num{
								StartPos: token.Pos(3),
								EndPos:   token.Pos(3),
								Val:      4,
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
			name:       "quantifier_ZeroOrOne_Successful",
			p:          r.quantifier,
			in:         newStringInput("??tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.RepOp{
							TokPos: token.Pos(0),
							TokTag: token.ZERO_OR_ONE,
						},
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
			name:       "quantifier_ZeroOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("*?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.RepOp{
							TokPos: token.Pos(0),
							TokTag: token.ZERO_OR_MORE,
						},
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
			name:       "quantifier_OneOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("+?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.RepOp{
							TokPos: token.Pos(0),
							TokTag: token.ONE_OR_MORE,
						},
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
			p:          r.quantifier,
			in:         newStringInput("{2}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.Range{
							OpenPos:  0,
							ClosePos: 2,
							Low: ast.Num{
								StartPos: token.Pos(1),
								EndPos:   token.Pos(1),
								Val:      2,
							},
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
			name:       "quantifier_range_upper_Unbounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.Range{
							OpenPos:  0,
							ClosePos: 3,
							Low: ast.Num{
								StartPos: token.Pos(1),
								EndPos:   token.Pos(1),
								Val:      2,
							},
							Up: &ast.UpperBound{},
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
			name:       "quantifier_range_upper_Bounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,4}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Quantifier{
						Rep: &ast.Range{
							OpenPos:  0,
							ClosePos: 4,
							Low: ast.Num{
								StartPos: token.Pos(1),
								EndPos:   token.Pos(1),
								Val:      2,
							},
							Up: &ast.UpperBound{
								CommaPos: 3,
								Val: &ast.Num{
									StartPos: token.Pos(3),
									EndPos:   token.Pos(3),
									Val:      4,
								},
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
			p:          r.charRange,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharRange{
						Low: ast.Char{
							TokPos: token.Pos(0),
							Val:    '0',
						},
						Up: ast.Char{
							TokPos: token.Pos(2),
							Val:    '9',
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
			name:       "charClass_Digit_Successful",
			p:          r.charClass,
			in:         newStringInput(`\dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.DIGIT,
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
			name:       "charClass_NotDigit_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.NON_DIGIT,
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
			name:       "charClass_Whitespace_Successful",
			p:          r.charClass,
			in:         newStringInput(`\stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.WHITESPACE,
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
			name:       "charClass_NotWhitespace_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.NON_WHITESPACE,
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
			name:       "charClass_Word_Successful",
			p:          r.charClass,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.WORD,
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
			name:       "charClass_NotWord_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.NON_WORD,
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
			name:       "asciiCharClass_Blank_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:blank:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.BLANK_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.SPACE_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.DIGIT_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(9),
						TokTag:   token.XDIGIT_CHARS,
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
			name:       "asciiCharClass_Upper_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:upper:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.UPPER_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.LOWER_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.ALPHA_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.ALNUM_CHARS,
					},
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
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(7),
						TokTag:   token.WORD_CHARS,
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
			name:       "asciiCharClass_ASCII_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:ascii:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(8),
						TokTag:   token.ASCII_CHARS,
					},
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
					Val: &ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.WORD,
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
			name:       "charGroupItem_asciiCharClass_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(7),
						TokTag:   token.WORD_CHARS,
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
			name:       "charGroupItem_charRange_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.CharRange{
						Low: ast.Char{
							TokPos: token.Pos(0),
							Val:    '0',
						},
						Up: ast.Char{
							TokPos: token.Pos(2),
							Val:    '9',
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
			name:       "charGroupItem_char_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Char{
						TokPos: token.Pos(0),
						Val:    '!',
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
			name:       "charGroup_charClass_Successful",
			p:          r.charGroup,
			in:         newStringInput(`[\w]tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharGroup{
						OpenPos:  0,
						ClosePos: 3,
						Items: []ast.CharGroupItem{
							&ast.CharClass{
								StartPos: token.Pos(1),
								EndPos:   token.Pos(2),
								TokTag:   token.WORD,
							},
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
			p:          r.charGroup,
			in:         newStringInput("[[:word:]]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharGroup{
						OpenPos:  0,
						ClosePos: 9,
						Items: []ast.CharGroupItem{
							&ast.ASCIICharClass{
								StartPos: token.Pos(1),
								EndPos:   token.Pos(8),
								TokTag:   token.WORD_CHARS,
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
			name:       "charGroup_charRange_Successful",
			p:          r.charGroup,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharGroup{
						OpenPos:  0,
						ClosePos: 4,
						Items: []ast.CharGroupItem{
							&ast.CharRange{
								Low: ast.Char{
									TokPos: token.Pos(1),
									Val:    '0',
								},
								Up: ast.Char{
									TokPos: token.Pos(3),
									Val:    '9',
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
			name:       "charGroup_char_Successful",
			p:          r.charGroup,
			in:         newStringInput("[!]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharGroup{
						OpenPos:  0,
						ClosePos: 2,
						Items: []ast.CharGroupItem{
							&ast.Char{
								TokPos: token.Pos(1),
								Val:    '!',
							},
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
			p:          r.charGroup,
			in:         newStringInput("[^#$]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.CharGroup{
						OpenPos:  0,
						ClosePos: 4,
						Negated:  true,
						Items: []ast.CharGroupItem{
							&ast.Char{
								TokPos: token.Pos(2),
								Val:    '#',
							},
							&ast.Char{
								TokPos: token.Pos(3),
								Val:    '$',
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
			name:       "anyChar_Successful",
			p:          r.anyChar,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.AnyChar{
						TokPos: token.Pos(0),
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
			name:       "matchItem_anyChar_Successful",
			p:          r.matchItem,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.AnyChar{
						TokPos: token.Pos(0),
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
			name:       "matchItem_charClass_Successful",
			p:          r.matchItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.CharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(1),
						TokTag:   token.WORD,
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
			name:       "matchItem_asciiCharClass_Successful",
			p:          r.matchItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.ASCIICharClass{
						StartPos: token.Pos(0),
						EndPos:   token.Pos(7),
						TokTag:   token.WORD_CHARS,
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
			name:       "matchItem_charGroup_charRange_Successful",
			p:          r.matchItem,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.CharGroup{
						OpenPos:  0,
						ClosePos: 4,
						Items: []ast.CharGroupItem{
							&ast.CharRange{
								Low: ast.Char{
									TokPos: token.Pos(1),
									Val:    '0',
								},
								Up: ast.Char{
									TokPos: token.Pos(3),
									Val:    '9',
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
			name:       "matchItem_char_Successful",
			p:          r.matchItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Char{
						TokPos: token.Pos(0),
						Val:    '!',
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
					Val: ast.Match{
						Item: &ast.CharClass{
							StartPos: token.Pos(0),
							EndPos:   token.Pos(1),
							TokTag:   token.WORD,
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
			name:       "match_asciiCharClass_Successful",
			p:          r.match,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Match{
						Item: &ast.ASCIICharClass{
							StartPos: token.Pos(0),
							EndPos:   token.Pos(7),
							TokTag:   token.WORD_CHARS,
						},
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
					Val: ast.Match{
						Item: &ast.CharGroup{
							OpenPos:  0,
							ClosePos: 4,
							Items: []ast.CharGroupItem{
								&ast.CharRange{
									Low: ast.Char{
										TokPos: token.Pos(1),
										Val:    '0',
									},
									Up: ast.Char{
										TokPos: token.Pos(3),
										Val:    '9',
									},
								},
							},
						},
						Quant: &ast.Quantifier{
							Rep: &ast.Range{
								OpenPos:  5,
								ClosePos: 9,
								Low: ast.Num{
									StartPos: token.Pos(6),
									EndPos:   token.Pos(6),
									Val:      2,
								},
								Up: &ast.UpperBound{
									CommaPos: 8,
									Val: &ast.Num{
										StartPos: token.Pos(8),
										EndPos:   token.Pos(8),
										Val:      4,
									},
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
			name:       "match_char_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("#{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Match{
						Item: &ast.Char{
							Val:    '#',
							TokPos: token.Pos(0),
						},
						Quant: &ast.Quantifier{
							Rep: &ast.Range{
								OpenPos:  1,
								ClosePos: 5,
								Low: ast.Num{
									StartPos: token.Pos(2),
									EndPos:   token.Pos(2),
									Val:      2,
								},
								Up: &ast.UpperBound{
									CommaPos: 4,
									Val: &ast.Num{
										StartPos: token.Pos(4),
										EndPos:   token.Pos(4),
										Val:      4,
									},
								},
							},
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
					Val: ast.Group{
						OpenPos:  0,
						ClosePos: 4,
						Expr: ast.Expr{
							Sub: ast.Subexpr{
								Items: []ast.SubexprItem{
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(1),
											Val:    'a',
										},
									},
								},
							},
							Expr: &ast.Expr{
								Sub: ast.Subexpr{
									Items: []ast.SubexprItem{
										&ast.Match{
											Item: &ast.Char{
												TokPos: token.Pos(3),
												Val:    'b',
											},
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
			name:       "group_quantifier_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Group{
						OpenPos:  0,
						ClosePos: 4,
						Expr: ast.Expr{
							Sub: ast.Subexpr{
								Items: []ast.SubexprItem{
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(1),
											Val:    'a',
										},
									},
								},
							},
							Expr: &ast.Expr{
								Sub: ast.Subexpr{
									Items: []ast.SubexprItem{
										&ast.Match{
											Item: &ast.Char{
												TokPos: token.Pos(3),
												Val:    'b',
											},
										},
									},
								},
							},
						},
						Quant: &ast.Quantifier{
							Rep: &ast.RepOp{
								TokPos: token.Pos(5),
								TokTag: token.ONE_OR_MORE,
							},
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
			name:       "backref_Unsuccessful",
			p:          r.backref,
			in:         newStringInput(`\2tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Backref{
						SlashPos: token.Pos(0),
						Ref: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Group: nil,
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
				Result: result{
					Val: ast.Anchor{
						TokPos: token.Pos(0),
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
			name:       "subexprItem_group_Successful",
			p:          r.subexprItem,
			in:         newStringInput("(ab)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Group{
						OpenPos:  0,
						ClosePos: 3,
						Expr: ast.Expr{
							Sub: ast.Subexpr{
								Items: []ast.SubexprItem{
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(1),
											Val:    'a',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(2),
											Val:    'b',
										},
									},
								},
							},
						},
						Quant: &ast.Quantifier{
							Rep: &ast.RepOp{
								TokPos: token.Pos(4),
								TokTag: token.ONE_OR_MORE,
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
			name:       "subexprItem_anchor_Successful",
			p:          r.subexprItem,
			in:         newStringInput("$"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Anchor{
						TokPos: token.Pos(0),
					},
					Pos: 0,
				},
				Remaining: nil,
			},
		},
		{
			name:       "subexprItem_backref_Unsuccessful",
			p:          r.subexprItem,
			in:         newStringInput(`\2tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Backref{
						SlashPos: token.Pos(0),
						Ref: ast.Num{
							StartPos: token.Pos(1),
							EndPos:   token.Pos(1),
							Val:      2,
						},
						Group: nil,
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
					Val: &ast.Match{
						Item: &ast.CharGroup{
							OpenPos:  0,
							ClosePos: 4,
							Items: []ast.CharGroupItem{
								&ast.CharRange{
									Low: ast.Char{
										TokPos: token.Pos(1),
										Val:    '0',
									},
									Up: ast.Char{
										TokPos: token.Pos(3),
										Val:    '9',
									},
								},
							},
						},
						Quant: &ast.Quantifier{
							Rep: &ast.RepOp{
								TokPos: token.Pos(5),
								TokTag: token.ONE_OR_MORE,
							},
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
					Val: ast.Subexpr{
						Items: []ast.SubexprItem{
							&ast.Group{
								OpenPos:  0,
								ClosePos: 3,
								Expr: ast.Expr{
									Sub: ast.Subexpr{
										Items: []ast.SubexprItem{
											&ast.Match{
												Item: &ast.Char{
													TokPos: token.Pos(1),
													Val:    'a',
												},
											},
											&ast.Match{
												Item: &ast.Char{
													TokPos: token.Pos(2),
													Val:    'b',
												},
											},
										},
									},
								},
								Quant: &ast.Quantifier{
									Rep: &ast.RepOp{
										TokPos: token.Pos(4),
										TokTag: token.ONE_OR_MORE,
									},
								},
							},
							&ast.Match{
								Item: &ast.CharGroup{
									OpenPos:  5,
									ClosePos: 9,
									Items: []ast.CharGroupItem{
										&ast.CharRange{
											Low: ast.Char{
												TokPos: token.Pos(6),
												Val:    '0',
											},
											Up: ast.Char{
												TokPos: token.Pos(8),
												Val:    '9',
											},
										},
									},
								},
								Quant: &ast.Quantifier{
									Rep: &ast.RepOp{
										TokPos: token.Pos(10),
										TokTag: token.ZERO_OR_MORE,
									},
								},
							},
							&ast.Match{
								Item: &ast.Char{
									TokPos: token.Pos(11),
									Val:    't',
								},
							},
							&ast.Match{
								Item: &ast.Char{
									TokPos: token.Pos(12),
									Val:    'a',
								},
							},
							&ast.Match{
								Item: &ast.Char{
									TokPos: token.Pos(13),
									Val:    'i',
								},
							},
							&ast.Match{
								Item: &ast.Char{
									TokPos: token.Pos(14),
									Val:    'l',
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
			name:       "expr_Successful",
			p:          r.expr,
			in:         newStringInput(`[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Expr{
						Sub: ast.Subexpr{
							Items: []ast.SubexprItem{
								&ast.Match{
									Item: &ast.CharGroup{
										OpenPos:  token.Pos(0),
										ClosePos: token.Pos(11),
										Items: []ast.CharGroupItem{
											&ast.CharRange{
												Low: ast.Char{
													TokPos: token.Pos(1),
													Val:    '0',
												},
												Up: ast.Char{
													TokPos: token.Pos(3),
													Val:    '9',
												},
											},
											&ast.CharRange{
												Low: ast.Char{
													TokPos: token.Pos(4),
													Val:    'A',
												},
												Up: ast.Char{
													TokPos: token.Pos(6),
													Val:    'Z',
												},
											},
											&ast.CharRange{
												Low: ast.Char{
													TokPos: token.Pos(7),
													Val:    'a',
												},
												Up: ast.Char{
													TokPos: token.Pos(9),
													Val:    'z',
												},
											},
											&ast.Char{
												TokPos: token.Pos(10),
												Val:    '_',
											},
										},
									},
								},
								&ast.Match{
									Item: &ast.CharGroup{
										OpenPos:  token.Pos(12),
										ClosePos: token.Pos(17),
										Items: []ast.CharGroupItem{
											&ast.CharClass{
												StartPos: token.Pos(13),
												EndPos:   token.Pos(14),
												TokTag:   token.DIGIT,
											},
											&ast.CharClass{
												StartPos: token.Pos(15),
												EndPos:   token.Pos(16),
												TokTag:   token.WORD,
											},
										},
									},
									Quant: &ast.Quantifier{
										Rep: &ast.RepOp{
											TokPos: token.Pos(18),
											TokTag: token.ZERO_OR_MORE,
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
		{
			name:       "regex_Successful",
			p:          r.regex,
			in:         newStringInput(`^package\s+[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ast.Regex{
						SOS: true,
						Expr: ast.Expr{
							Sub: ast.Subexpr{
								Items: []ast.SubexprItem{
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(1),
											Val:    'p',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(2),
											Val:    'a',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(3),
											Val:    'c',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(4),
											Val:    'k',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(5),
											Val:    'a',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(6),
											Val:    'g',
										},
									},
									&ast.Match{
										Item: &ast.Char{
											TokPos: token.Pos(7),
											Val:    'e',
										},
									},
									&ast.Match{
										Item: &ast.CharClass{
											StartPos: token.Pos(8),
											EndPos:   token.Pos(9),
											TokTag:   token.WHITESPACE},
										Quant: &ast.Quantifier{
											Rep: &ast.RepOp{
												TokPos: token.Pos(10),
												TokTag: token.ONE_OR_MORE,
											},
										},
									},
									&ast.Match{
										Item: &ast.CharGroup{
											OpenPos:  token.Pos(11),
											ClosePos: token.Pos(22),
											Items: []ast.CharGroupItem{
												&ast.CharRange{
													Low: ast.Char{
														TokPos: token.Pos(12),
														Val:    '0',
													},
													Up: ast.Char{
														TokPos: token.Pos(14),
														Val:    '9',
													},
												},
												&ast.CharRange{
													Low: ast.Char{
														TokPos: token.Pos(15),
														Val:    'A',
													},
													Up: ast.Char{
														TokPos: token.Pos(17),
														Val:    'Z',
													},
												},
												&ast.CharRange{
													Low: ast.Char{
														TokPos: token.Pos(18),
														Val:    'a',
													},
													Up: ast.Char{
														TokPos: token.Pos(20),
														Val:    'z',
													},
												},
												&ast.Char{
													TokPos: token.Pos(21),
													Val:    '_',
												},
											},
										},
									},
									&ast.Match{
										Item: &ast.CharGroup{
											OpenPos:  token.Pos(23),
											ClosePos: token.Pos(28),
											Items: []ast.CharGroupItem{
												&ast.CharClass{
													StartPos: token.Pos(24),
													EndPos:   token.Pos(25),
													TokTag:   token.DIGIT,
												},
												&ast.CharClass{
													StartPos: token.Pos(26),
													EndPos:   token.Pos(27),
													TokTag:   token.WORD,
												},
											},
										},
										Quant: &ast.Quantifier{
											Rep: &ast.RepOp{
												TokPos: token.Pos(29),
												TokTag: token.ZERO_OR_MORE,
											},
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
			// Reset symbol table and errors
			c.symTab, c.errors = nil, nil

			out, ok := tc.p(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
