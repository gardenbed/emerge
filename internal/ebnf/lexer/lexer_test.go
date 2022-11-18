package lexer

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	auto "github.com/moorara/algo/automata"
	"github.com/stretchr/testify/assert"
)

func getLexerDFA() *auto.DFA {
	dfa := auto.NewDFA(0, auto.States{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 18, 20, 22, 26, 30, 31, 34})

	dfa.Add(0, ' ', 1)
	dfa.Add(0, '\t', 1)
	dfa.Add(0, '\n', 1)
	dfa.Add(0, '\r', 1)
	dfa.Add(0, '\f', 1)
	dfa.Add(0, '\v', 1)
	dfa.Add(1, ' ', 1)
	dfa.Add(1, '\t', 1)
	dfa.Add(1, '\n', 1)
	dfa.Add(1, '\r', 1)
	dfa.Add(1, '\f', 1)
	dfa.Add(1, '\v', 1)

	dfa.Add(0, '=', 2)
	dfa.Add(0, '|', 3)
	dfa.Add(0, '(', 4)
	dfa.Add(0, ')', 5)
	dfa.Add(0, '[', 6)
	dfa.Add(0, ']', 7)
	dfa.Add(0, '{', 8)
	dfa.Add(8, '{', 9)
	dfa.Add(0, '}', 10)
	dfa.Add(10, '}', 11)

	dfa.Add(0, 'g', 12)
	dfa.Add(12, 'r', 13)
	dfa.Add(13, 'a', 14)
	dfa.Add(14, 'm', 15)
	dfa.Add(15, 'm', 16)
	dfa.Add(16, 'a', 17)
	dfa.Add(17, 'r', 18)

	//==================================================< IDENTIFIER >==================================================

	for r := 'a'; r <= 'z'; r++ {
		if r != 'g' {
			dfa.Add(0, auto.Symbol(r), 19)
		}

		if r != 'r' {
			dfa.Add(12, auto.Symbol(r), 20)
			dfa.Add(17, auto.Symbol(r), 20)
		}

		if r != 'a' {
			dfa.Add(13, auto.Symbol(r), 20)
			dfa.Add(16, auto.Symbol(r), 20)
		}

		if r != 'm' {
			dfa.Add(14, auto.Symbol(r), 20)
			dfa.Add(15, auto.Symbol(r), 20)
		}

		dfa.Add(18, auto.Symbol(r), 20)
		dfa.Add(19, auto.Symbol(r), 20)
		dfa.Add(20, auto.Symbol(r), 20)
	}

	for r := '0'; r <= '9'; r++ {
		dfa.Add(12, auto.Symbol(r), 20)
		dfa.Add(13, auto.Symbol(r), 20)
		dfa.Add(14, auto.Symbol(r), 20)
		dfa.Add(15, auto.Symbol(r), 20)
		dfa.Add(16, auto.Symbol(r), 20)
		dfa.Add(17, auto.Symbol(r), 20)
		dfa.Add(18, auto.Symbol(r), 20)
		dfa.Add(19, auto.Symbol(r), 20)
		dfa.Add(20, auto.Symbol(r), 20)
	}

	dfa.Add(12, '_', 20)
	dfa.Add(13, '_', 20)
	dfa.Add(14, '_', 20)
	dfa.Add(15, '_', 20)
	dfa.Add(16, '_', 20)
	dfa.Add(17, '_', 20)
	dfa.Add(18, '_', 20)
	dfa.Add(19, '_', 20)
	dfa.Add(20, '_', 20)

	//==================================================< TOKEN >==================================================

	for r := 'A'; r <= 'Z'; r++ {
		dfa.Add(0, auto.Symbol(r), 21)
		dfa.Add(21, auto.Symbol(r), 22)
		dfa.Add(22, auto.Symbol(r), 22)
	}

	for r := '0'; r <= '9'; r++ {
		dfa.Add(21, auto.Symbol(r), 22)
		dfa.Add(22, auto.Symbol(r), 22)
	}

	dfa.Add(21, '_', 22)
	dfa.Add(22, '_', 22)

	//==================================================< STRING >==================================================

	dfa.Add(0, '"', 23)
	dfa.Add(23, '\\', 24)
	dfa.Add(25, '\\', 24)
	dfa.Add(25, '"', 26)

	for r := 0x21; r <= 0x7E; r++ {
		dfa.Add(24, auto.Symbol(r), 25)
		if r != '"' && r != '\\' {
			dfa.Add(23, auto.Symbol(r), 25)
			dfa.Add(25, auto.Symbol(r), 25)
		}
	}

	//==================================================< REGEX >==================================================

	dfa.Add(0, '/', 27)
	dfa.Add(27, '\\', 28)
	dfa.Add(29, '\\', 28)
	dfa.Add(29, '/', 30)

	for r := 0x20; r <= 0x7E; r++ {
		if r != '*' && r != '/' && r != '\\' {
			dfa.Add(27, auto.Symbol(r), 29)
		}

		dfa.Add(28, auto.Symbol(r), 29)

		if r != '/' && r != '\\' {
			dfa.Add(29, auto.Symbol(r), 29)
		}
	}

	//==================================================< SINGLE-LINE COMMENT >==================================================

	dfa.Add(27, '/', 31)

	for r := 0x20; r <= 0x7E; r++ {
		dfa.Add(31, auto.Symbol(r), 31)
	}

	//==================================================< MULTI-LINE COMMENT >==================================================

	dfa.Add(27, '*', 32)
	dfa.Add(32, '*', 33)
	dfa.Add(33, '/', 34)

	for _, r := range []rune{'\t', '\n', '\r'} {
		dfa.Add(32, auto.Symbol(r), 32)
		dfa.Add(33, auto.Symbol(r), 32)
	}

	for r := 0x20; r <= 0x7E; r++ {
		if r != '*' {
			dfa.Add(32, auto.Symbol(r), 32)
		}

		if r != '/' {
			dfa.Add(33, auto.Symbol(r), 32)
		}
	}

	return dfa
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		src           io.Reader
		expectedError string
	}{
		{
			name:          "Success",
			src:           strings.NewReader("Lorem ipsum"),
			expectedError: "",
		},
		{
			name:          "Failure",
			src:           iotest.ErrReader(errors.New("io error")),
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lex, err := New(tc.src)

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
		expectedToken Token
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
			expectedToken: Token{},
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
			expectedToken: Token{},
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
						{OutVal: "#", OutPos: 2},
					},
				},
			},
			expectedToken: Token{},
			expectedError: "lexical error at 2:#",
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
						{OutVal: "statement", OutPos: 4},
					},
				},
			},
			expectedToken: Token{IDENT, "statement", 4},
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
						{OutPos: 2},
					},
					LexemeMocks: []LexemeMock{
						{OutVal: "expression", OutPos: 4},
					},
				},
			},
			expectedToken: Token{IDENT, "expression", 4},
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
						{OutPos: 2},
						{OutPos: 4},
						{OutPos: 8},
					},
				},
			},
			expectedToken: Token{GRAMMER, "grammar", 8},
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
						{OutPos: 2},
						{OutPos: 4},
						{OutPos: 8},
					},
				},
			},
			expectedToken: Token{GRAMMER, "grammar", 8},
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
		expectedToken Token
	}{
		{
			name: "WS",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 2},
					},
				},
			},
			state:         1,
			expectedToken: Token{WS, "", 2},
		},
		{
			name: "DEF",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         2,
			expectedToken: Token{DEF, "=", 4},
		},
		{
			name: "ALT",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         3,
			expectedToken: Token{ALT, "|", 4},
		},
		{
			name: "LPAREN",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         4,
			expectedToken: Token{LPAREN, "(", 4},
		},
		{
			name: "RPAREN",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         5,
			expectedToken: Token{RPAREN, ")", 4},
		},
		{
			name: "LBRACK",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         6,
			expectedToken: Token{LBRACK, "[", 4},
		},
		{
			name: "RBRACK",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         7,
			expectedToken: Token{RBRACK, "]", 4},
		},
		{
			name: "LBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         8,
			expectedToken: Token{LBRACE, "{", 4},
		},
		{
			name: "LLBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         9,
			expectedToken: Token{LLBRACE, "{{", 4},
		},
		{
			name: "RBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         10,
			expectedToken: Token{RBRACE, "}", 4},
		},
		{
			name: "RRBRACE",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         11,
			expectedToken: Token{RRBRACE, "}}", 4},
		},
		{
			name: "GRAMMER",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 4},
					},
				},
			},
			state:         18,
			expectedToken: Token{GRAMMER, "grammar", 4},
		},
		{
			name: "IDENT",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{OutVal: "statement", OutPos: 8},
					},
				},
			},
			state:         20,
			expectedToken: Token{IDENT, "statement", 8},
		},
		{
			name: "TOKEN",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{OutVal: "NUM", OutPos: 8},
					},
				},
			},
			state:         22,
			expectedToken: Token{TOKEN, "NUM", 8},
		},
		{
			name: "STRING",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{OutVal: `"foo"`, OutPos: 16},
					},
				},
			},
			state:         26,
			expectedToken: Token{STRING, `foo`, 16},
		},
		{
			name: "REGEX",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{OutVal: `/[a-z]+/`, OutPos: 16},
					},
				},
			},
			state:         30,
			expectedToken: Token{REGEX, `[a-z]+`, 16},
		},
		{
			name: "COMMENT_SingleLine",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 32},
					},
				},
			},
			state:         31,
			expectedToken: Token{COMMENT, "", 32},
		},
		{
			name: "COMMENT_MultiLine",
			l: &Lexer{
				in: &mockInputBuffer{
					SkipMocks: []SkipMock{
						{OutPos: 32},
					},
				},
			},
			state:         34,
			expectedToken: Token{COMMENT, "", 32},
		},
		{
			name: "ERR",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{OutVal: "foo", OutPos: 64},
					},
				},
			},
			state:         12,
			expectedToken: Token{ERR, "lexical error at 64:foo", 64},
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
