package lexer

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/lexer"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		src           io.Reader
		expectedError string
	}{
		{
			name:          "Success",
			filename:      "lorem_ipsum",
			src:           strings.NewReader("Lorem ipsum"),
			expectedError: "",
		},
		{
			name:          "Failure",
			filename:      "lorem_ipsum",
			src:           iotest.ErrReader(errors.New("io error")),
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lex, err := New(tc.filename, tc.src)

			if tc.expectedError == "" {
				assert.NotNil(t, lex)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, lex)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name          string
		l             *Lexer
		expectedToken lexer.Token
		expectedError string
	}{
		{
			name: "IOError",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutError: errors.New("io error")},
					},
				},
			},
			expectedToken: lexer.Token{},
			expectedError: "io error",
		},
		{
			name: "EOF",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutError: io.EOF},
					},
				},
			},
			expectedToken: lexer.Token{},
			expectedError: "EOF",
		},
		{
			name: "LexicalError",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutRune: '#'},
					},
					LexemeMocks: []LexemeMock{
						{
							OutVal: "#",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   2,
								Line:     1,
								Column:   3,
							},
						},
					},
				},
			},
			expectedToken: lexer.Token{},
			expectedError: "lexical error at test:1:3:#",
		},
		{
			name: "Identifier",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutRune: 's'},
						{OutRune: 't'},
						{OutRune: 'a'},
						{OutRune: 't'},
						{OutRune: 'e'},
						{OutRune: 'm'},
						{OutRune: 'e'},
						{OutRune: 'n'},
						{OutRune: 't'},
						{OutRune: '='},
					},
					LexemeMocks: []LexemeMock{
						{
							OutVal: "statement",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "statement",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
			expectedError: "",
		},
		{
			name: "Identifier_After_Whitespace",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutRune: '\n'},
						{OutRune: '\t'},
						{OutRune: 'e'},
						{OutRune: 'x'},
						{OutRune: 'p'},
						{OutRune: 'r'},
						{OutRune: 'e'},
						{OutRune: 's'},
						{OutRune: 's'},
						{OutRune: 'i'},
						{OutRune: 'o'},
						{OutRune: 'n'},
						{OutRune: '='},
					},
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   2,
								Line:     1,
								Column:   3,
							},
						},
					},
					LexemeMocks: []LexemeMock{
						{
							OutVal: "expression",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "expression",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
			expectedError: "",
		},
		{
			name: "Keyword_After_SingleLineComment",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutRune: '/'},
						{OutRune: '/'},
						{OutRune: ' '},
						{OutRune: 'C'},
						{OutRune: 'o'},
						{OutRune: 'm'},
						{OutRune: 'm'},
						{OutRune: 'e'},
						{OutRune: 'n'},
						{OutRune: 't'},
						{OutRune: '\n'},
						{OutRune: '\n'}, // Repeat after retract
						{OutRune: 'g'},
						{OutRune: 'g'}, // Repeat after retract
						{OutRune: 'r'},
						{OutRune: 'a'},
						{OutRune: 'm'},
						{OutRune: 'm'},
						{OutRune: 'a'},
						{OutRune: 'r'},
						{OutRune: ' '},
					},
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   2,
								Line:     1,
								Column:   3,
							},
						},
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   8,
								Line:     1,
								Column:   9,
							},
						},
					},
				},
			},
			expectedToken: lexer.Token{
				Terminal: GRAMMER,
				Lexeme:   "grammar",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
			expectedError: "",
		},
		{
			name: "Keyword_After_MultiLineComment",
			l: &Lexer{
				in: &mockInputBuffer{
					NextMocks: []NextMock{
						{OutRune: '/'},
						{OutRune: '*'},
						{OutRune: 'C'},
						{OutRune: 'o'},
						{OutRune: 'm'},
						{OutRune: 'm'},
						{OutRune: 'e'},
						{OutRune: 'n'},
						{OutRune: 't'},
						{OutRune: '*'},
						{OutRune: '/'},
						{OutRune: '\n'},
						{OutRune: '\n'}, // Repeat after retract
						{OutRune: 'g'},
						{OutRune: 'g'}, // Repeat after retract
						{OutRune: 'r'},
						{OutRune: 'a'},
						{OutRune: 'm'},
						{OutRune: 'm'},
						{OutRune: 'a'},
						{OutRune: 'r'},
						{OutRune: ' '},
					},
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   2,
								Line:     1,
								Column:   3,
							},
						},
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   8,
								Line:     1,
								Column:   9,
							},
						},
					},
				},
			},
			expectedToken: lexer.Token{
				Terminal: GRAMMER,
				Lexeme:   "grammar",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, err := tc.l.NextToken()

			if tc.expectedError == "" {
				assert.Equal(t, tc.expectedToken, token)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectedToken, token)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestLexer_evalDFA(t *testing.T) {
	tests := []struct {
		name          string
		l             *Lexer
		state         int
		expectedToken lexer.Token
	}{
		{
			name: "WS",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   2,
								Line:     1,
								Column:   3,
							},
						},
					},
				},
			},
			state: 1,
			expectedToken: lexer.Token{
				Terminal: WS,
				Lexeme:   "",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   2,
					Line:     1,
					Column:   3,
				},
			},
		},
		{
			name: "DEF",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 2,
			expectedToken: lexer.Token{
				Terminal: DEF,
				Lexeme:   "=",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "ALT",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 3,
			expectedToken: lexer.Token{
				Terminal: ALT,
				Lexeme:   "|",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "LPAREN",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 4,
			expectedToken: lexer.Token{
				Terminal: LPAREN,
				Lexeme:   "(",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RPAREN",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 5,
			expectedToken: lexer.Token{
				Terminal: RPAREN,
				Lexeme:   ")",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "LBRACK",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 6,
			expectedToken: lexer.Token{
				Terminal: LBRACK,
				Lexeme:   "[",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RBRACK",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 7,
			expectedToken: lexer.Token{
				Terminal: RBRACK,
				Lexeme:   "]",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "LBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 8,
			expectedToken: lexer.Token{
				Terminal: LBRACE,
				Lexeme:   "{",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 9,
			expectedToken: lexer.Token{
				Terminal: RBRACE,
				Lexeme:   "}",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "LLBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 10,
			expectedToken: lexer.Token{
				Terminal: LLBRACE,
				Lexeme:   "{{",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RRBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 11,
			expectedToken: lexer.Token{
				Terminal: RRBRACE,
				Lexeme:   "}}",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "LASSOC",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 16,
			expectedToken: lexer.Token{
				Terminal: LASSOC,
				Lexeme:   "@left",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RASSOC",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 21,
			expectedToken: lexer.Token{
				Terminal: RASSOC,
				Lexeme:   "@right",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "NOASSOC",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 25,
			expectedToken: lexer.Token{
				Terminal: NOASSOC,
				Lexeme:   "@none",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "PREDEF",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "$STRING",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 27,
			expectedToken: lexer.Token{
				Terminal: PREDEF,
				Lexeme:   "$STRING",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "GRAMMER",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   4,
								Line:     1,
								Column:   5,
							},
						},
					},
				},
			},
			state: 34,
			expectedToken: lexer.Token{
				Terminal: GRAMMER,
				Lexeme:   "grammar",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "IDENT",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "statement",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   8,
								Line:     1,
								Column:   9,
							},
						},
					},
				},
			},
			state: 36,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "statement",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "TOKEN",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "NUM",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   8,
								Line:     1,
								Column:   9,
							},
						},
					},
				},
			},
			state: 38,
			expectedToken: lexer.Token{
				Terminal: TOKEN,
				Lexeme:   "NUM",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "STRING",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: `"foo"`,
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   16,
								Line:     1,
								Column:   17,
							},
						},
					},
				},
			},
			state: 42,
			expectedToken: lexer.Token{
				Terminal: STRING,
				Lexeme:   "foo",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   16,
					Line:     1,
					Column:   17,
				},
			},
		},
		{
			name: "REGEX",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: `/[a-z]+/`,
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   16,
								Line:     1,
								Column:   17,
							},
						},
					},
				},
			},
			state: 46,
			expectedToken: lexer.Token{
				Terminal: REGEX,
				Lexeme:   "[a-z]+",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   16,
					Line:     1,
					Column:   17,
				},
			},
		},
		{
			name: "COMMENT_SingleLine",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   32,
								Line:     1,
								Column:   33,
							},
						},
					},
				},
			},
			state: 47,
			expectedToken: lexer.Token{
				Terminal: COMMENT,
				Lexeme:   "",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   32,
					Line:     1,
					Column:   33,
				},
			},
		},
		{
			name: "COMMENT_MultiLine",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   32,
								Line:     1,
								Column:   33,
							},
						},
					},
				},
			},
			state: 50,
			expectedToken: lexer.Token{
				Terminal: COMMENT,
				Lexeme:   "",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   32,
					Line:     1,
					Column:   33,
				},
			},
		},
		{
			name: "ERR",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "foo",
							OutPos: lexer.Position{
								Filename: "test",
								Offset:   64,
								Line:     1,
								Column:   65,
							},
						},
					},
				},
			},
			state: 12,
			expectedToken: lexer.Token{
				Terminal: ERR,
				Lexeme:   "lexical error at test:1:65:foo",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   64,
					Line:     1,
					Column:   65,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := tc.l.evalDFA(tc.state)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestAdvanceDFA(t *testing.T) {
	tests := []struct {
		name          string
		state         int
		r             rune
		expectedState int
	}{
		{"0_Space", 0, ' ', 1},
		{"0_Tab", 0, '\t', 1},
		{"0_NewLine", 0, '\n', 1},

		{"1_Space", 1, ' ', 1},
		{"1_Tab", 1, '\t', 1},
		{"1_NewLine", 1, '\n', 1},

		{"0_DEF", 0, '=', 2},
		{"0_ALT", 0, '|', 3},
		{"0_LPAREN", 0, '(', 4},
		{"0_RPAREN", 0, ')', 5},
		{"0_LBRACK", 0, '[', 6},
		{"0_RBRACK", 0, ']', 7},
		{"0_LBRACE", 0, '{', 8},
		{"0_RBRACE", 0, '}', 9},
		{"8_LBRACE", 8, '{', 10},
		{"9_RBRACE", 9, '}', 11},

		// @
		{"0_@", 0, '@', 12},

		// @left
		{"12_l", 12, 'l', 13},
		{"13_e", 13, 'e', 14},
		{"14_f", 14, 'f', 15},
		{"15_t", 15, 't', 16},

		// @ right
		{"12_r", 12, 'r', 17},
		{"17_i", 17, 'i', 18},
		{"18_g", 18, 'g', 19},
		{"19_h", 19, 'h', 20},
		{"20_t", 20, 't', 21},

		// @none
		{"12_n", 12, 'n', 22},
		{"22_o", 22, 'o', 23},
		{"23_n", 23, 'n', 24},
		{"24_e", 24, 'e', 25},

		// $
		{"0_$", 0, '$', 26},

		// $STRING
		{"26_S", 26, 'S', 27},
		{"27_T", 27, 'S', 27},
		{"27_R", 27, 'R', 27},
		{"27_I", 27, 'I', 27},
		{"27_N", 27, 'N', 27},
		{"27_G", 27, 'G', 27},

		// grammar
		{"0_g", 0, 'g', 28},
		{"28_r", 28, 'r', 29},
		{"29_a", 29, 'a', 30},
		{"30_m", 30, 'm', 31},
		{"31_m", 31, 'm', 32},
		{"32_a", 32, 'a', 33},
		{"33_r", 33, 'r', 34},

		{"28_o", 28, 'o', 36}, // go
		{"29_e", 29, 'e', 36}, // gre
		{"30_n", 30, 'n', 36}, // gran
		{"31_s", 31, 's', 36}, // grams
		{"32_o", 32, 'o', 36}, // grammo
		{"33_t", 33, 't', 36}, // grammat
		{"34_i", 34, 'i', 36}, // grammari
		{"36_n", 36, 'n', 36}, // grammarian

		// decl
		{"0_d", 0, 'd', 35},
		{"35_e", 35, 'e', 36},
		{"36_c", 36, 'c', 36},
		{"36_l", 36, 'l', 36},

		// NAME
		{"0_N", 0, 'N', 37},
		{"37_A", 37, 'A', 38},
		{"38_M", 38, 'M', 38},
		{"38_E", 38, 'E', 38},

		// "foo"
		{"0_QUOT", 0, '"', 39},
		{"39_f", 39, 'f', 41},
		{"41_o", 41, 'o', 41},
		{"41_o", 41, 'o', 41},
		{"41_QUOT", 41, '"', 42},

		// "\"\\"
		{"0_QUOT", 0, '"', 39},
		{"39_BSOL", 39, '\\', 40},
		{"40_QUOT", 40, '"', 41},
		{"41_BSOL", 41, '\\', 40},
		{"40_BSOL", 40, '\\', 41},
		{"41_QUOT", 41, '"', 42},

		// /[a-z]+/
		{"0_SOL", 0, '/', 43},
		{"43_LSQB", 43, '[', 45},
		{"45_a", 45, 'a', 45},
		{"45_MINUS", 45, '-', 45},
		{"45_z", 45, 'z', 45},
		{"45_RSQB", 45, ']', 45},
		{"45_PLUS", 45, '+', 45},
		{"45_SOL", 45, '/', 46},

		// /\{\}/
		{"0_SOL", 0, '/', 43},
		{"43_BSOL", 43, '\\', 44},
		{"44_LCUB", 44, '{', 45},
		{"45_BSOL", 45, '\\', 44},
		{"44_RCUB", 44, '}', 45},
		{"45_SOL", 45, '/', 46},

		// // comment
		{"0_SOL", 0, '/', 43},
		{"43_SOL", 43, '/', 47},
		{"47_Space", 47, ' ', 47},
		{"47_c", 47, 'c', 47},
		{"47_o", 47, 'o', 47},
		{"47_m", 47, 'm', 47},
		{"47_m", 47, 'm', 47},
		{"47_e", 47, 'e', 47},
		{"47_n", 47, 'n', 47},
		{"47_t", 47, 't', 47},

		// /* foo*bar */
		{"0_SOL", 0, '/', 43},
		{"43_AST", 43, '*', 48},
		{"48_Space", 48, ' ', 48},
		{"48_f", 48, 'f', 48},
		{"48_o", 48, 'o', 48},
		{"48_o", 48, 'o', 48},
		{"48_AST", 48, '*', 49},
		{"33_b", 49, 'b', 48},
		{"48_a", 48, 'a', 48},
		{"48_r", 48, 'r', 48},
		{"48_Space", 48, ' ', 48},
		{"48_AST", 48, '*', 49},
		{"49_SOL", 49, '/', 50},

		// ERR
		{"0_0", 0, '0', -1},
		{"0_1", 0, '1', -1},
		{"0_2", 0, '2', -1},
		{"0_3", 0, '3', -1},
		{"0_4", 0, '4', -1},
		{"0_5", 0, '5', -1},
		{"0_6", 0, '6', -1},
		{"0_7", 0, '7', -1},
		{"0_8", 0, '8', -1},
		{"0_9", 0, '9', -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := advanceDFA(tc.state, tc.r)
			assert.Equal(t, tc.expectedState, state)
		})
	}
}

func TestLexer(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "Success",
			file: "../fixture/please.grammar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			lex, err := New(tc.file, f)
			assert.NoError(t, err)

			for token, err := lex.NextToken(); err != io.EOF; token, err = lex.NextToken() {
				assert.NotEmpty(t, token)
			}
		})
	}
}
