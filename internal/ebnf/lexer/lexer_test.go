package lexer

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/moorara/algo/lexer"
	"github.com/stretchr/testify/assert"
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
			state: 9,
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
			state: 10,
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
			state: 18,
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
			state: 20,
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
			state: 22,
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
			state: 26,
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
			state: 30,
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
			state: 31,
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
			state: 34,
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
		{"0_RBRACE", 0, '}', 10},
		{"8_LLBRACE", 8, '{', 9},
		{"10_RRBRACE", 10, '}', 11},

		{"0_g", 0, 'g', 12},
		{"12_r", 12, 'r', 13},
		{"13_a", 13, 'a', 14},
		{"14_m", 14, 'm', 15},
		{"15_m", 15, 'm', 16},
		{"16_a", 16, 'a', 17},
		{"17_r", 17, 'r', 18},

		{"12_a", 12, 'a', 20},
		{"13_b", 13, 'b', 20},
		{"14_c", 14, 'c', 20},
		{"15_d", 15, 'd', 20},
		{"16_e", 16, 'e', 20},
		{"17_f", 17, 'f', 20},
		{"18_g", 18, 'g', 20},
		{"20_h", 20, 'h', 20},

		{"0_d", 0, 'd', 19},
		{"19_e", 19, 'e', 20},
		{"20_c", 20, 'c', 20},
		{"20_l", 20, 'l', 20},

		{"0_N", 0, 'N', 21},
		{"21_A", 21, 'A', 22},
		{"22_M", 22, 'M', 22},
		{"22_E", 22, 'E', 22},

		{"0_QUOT", 0, '"', 23},
		{"23_f", 23, 'f', 25},
		{"25_o", 25, 'o', 25},
		{"25_o", 25, 'o', 25},
		{"25_QUOT", 25, '"', 26},

		{"0_QUOT", 0, '"', 23},
		{"23_BSOL", 23, '\\', 24},
		{"24_QUOT", 24, '"', 25},
		{"25_BSOL", 25, '\\', 24},
		{"24_BSOL", 24, '\\', 25},
		{"25_QUOT", 25, '"', 26},

		{"0_SOL", 0, '/', 27},
		{"27_LSQB", 27, '[', 29},
		{"29_a", 29, 'a', 29},
		{"29_MINUS", 29, '-', 29},
		{"29_z", 29, 'z', 29},
		{"29_RSQB", 29, ']', 29},
		{"29_PLUS", 29, '+', 29},
		{"29_SOL", 29, '/', 30},

		{"0_SOL", 0, '/', 27},
		{"27_BSOL", 27, '\\', 28},
		{"28_LCUB", 28, '{', 29},
		{"29_BSOL", 29, '\\', 28},
		{"28_RCUB", 28, '}', 29},
		{"29_SOL", 29, '/', 30},

		{"0_SOL", 0, '/', 27},
		{"27_SOL", 27, '/', 31},
		{"31_C", 31, 'C', 31},
		{"31_o", 31, 'o', 31},
		{"31_m", 31, 'm', 31},
		{"31_m", 31, 'm', 31},
		{"31_e", 31, 'e', 31},
		{"31_n", 31, 'n', 31},
		{"31_t", 31, 't', 31},

		{"0_SOL", 0, '/', 27},
		{"27_AST", 27, '*', 32},
		{"32_F", 32, 'F', 32},
		{"32_i", 32, 'i', 32},
		{"32_r", 32, 'r', 32},
		{"32_s", 32, 's', 32},
		{"32_t", 32, 't', 32},
		{"32_AST", 32, '*', 33},
		{"33_S", 33, 'S', 32},
		{"32_e", 32, 's', 32},
		{"32_c", 32, 'c', 32},
		{"32_o", 32, 'o', 32},
		{"32_n", 32, 'n', 32},
		{"32_d", 32, 'd', 32},
		{"32_AST", 32, '*', 33},
		{"33_SOL", 33, '/', 34},

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
