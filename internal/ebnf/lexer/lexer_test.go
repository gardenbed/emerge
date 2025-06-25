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
			name: "EOL",
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
			state: 2,
			expectedToken: lexer.Token{
				Terminal: EOL,
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
			state: 3,
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
			name: "SEMI",
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
				Terminal: SEMI,
				Lexeme:   ";",
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
			state: 5,
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
			state: 6,
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
			state: 7,
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
			state: 8,
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
			state: 9,
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
			state: 10,
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
			state: 11,
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
			state: 12,
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
			state: 13,
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
			name: "LANGLE",
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
			state: 14,
			expectedToken: lexer.Token{
				Terminal: LANGLE,
				Lexeme:   "<",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
			},
		},
		{
			name: "RANGLE",
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
			state: 15,
			expectedToken: lexer.Token{
				Terminal: RANGLE,
				Lexeme:   ">",
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
			state: 17,
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
			state: 22,
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
			state: 27,
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
			state: 31,
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
			state: 38,
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
			name: "IDENT_01",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "gr",
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
			state: 40,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "gr",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_02",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "a",
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
			state: 40,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "a",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_03",
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
			state: 40,
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
			name: "TOKEN_01",
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
			state: 42,
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
			name: "STRING_01",
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
			state: 46,
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
			name: "REGEX_01",
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
			state: 50,
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
			state: 51,
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
			state: 54,
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
			state: 20,
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
		// Whitespace
		{"0_Tab", 0, '\t', 1},
		{"0_Space", 0, ' ', 1},
		{"1_Tab", 1, '\t', 1},
		{"1_Space", 1, ' ', 1},

		// Newline
		{"0_Tab", 0, '\n', 2},
		{"0_Space", 0, '\r', 2},
		{"2_Tab", 2, '\n', 2},
		{"2_Space", 2, '\r', 2},

		{"0_DEF", 0, '=', 3},
		{"0_SEMI", 0, ';', 4},
		{"0_ALT", 0, '|', 5},
		{"0_LPAREN", 0, '(', 6},
		{"0_RPAREN", 0, ')', 7},
		{"0_LBRACK", 0, '[', 8},
		{"0_RBRACK", 0, ']', 9},
		{"0_LBRACE", 0, '{', 10},
		{"0_RBRACE", 0, '}', 11},
		{"10_LBRACE", 10, '{', 12},
		{"11_RBRACE", 11, '}', 13},
		{"0_LANGLE", 0, '<', 14},
		{"0_RANGLE", 0, '>', 15},

		// $
		{"0_$", 0, '$', 16},

		// $STRING
		{"16_S", 16, 'S', 17},
		{"17_T", 17, 'S', 17},
		{"17_R", 17, 'R', 17},
		{"17_I", 17, 'I', 17},
		{"17_N", 17, 'N', 17},
		{"17_G", 17, 'G', 17},

		// @
		{"0_@", 0, '@', 18},

		// @left
		{"18_l", 18, 'l', 19},
		{"19_e", 19, 'e', 20},
		{"20_f", 20, 'f', 21},
		{"21_t", 21, 't', 22},

		// @ right
		{"18_r", 18, 'r', 23},
		{"23_i", 23, 'i', 24},
		{"24_g", 24, 'g', 25},
		{"25_h", 25, 'h', 26},
		{"26_t", 26, 't', 27},

		// @none
		{"18_n", 18, 'n', 28},
		{"28_o", 28, 'o', 29},
		{"29_n", 29, 'n', 30},
		{"30_e", 30, 'e', 31},

		// grammar
		{"0_g", 0, 'g', 32},
		{"32_r", 32, 'r', 33},
		{"33_a", 33, 'a', 34},
		{"34_m", 34, 'm', 35},
		{"35_m", 35, 'm', 36},
		{"36_a", 36, 'a', 37},
		{"37_r", 37, 'r', 38},

		{"32_o", 32, 'o', 40}, // go
		{"33_e", 33, 'e', 40}, // gre
		{"34_n", 34, 'n', 40}, // gran
		{"35_s", 35, 's', 40}, // grams
		{"36_o", 36, 'o', 40}, // grammo
		{"37_t", 37, 't', 40}, // grammat
		{"38_i", 38, 'i', 40}, // grammari
		{"40_n", 40, 'n', 40}, // grammarian

		// name
		{"0_d", 0, 'n', 39},
		{"39_e", 39, 'a', 40},
		{"40_c", 40, 'm', 40},
		{"40_l", 40, 'e', 40},

		// NAME
		{"0_N", 0, 'N', 41},
		{"41_A", 41, 'A', 42},
		{"42_M", 42, 'M', 42},
		{"42_E", 42, 'E', 42},

		// "foo"
		{"0_QUOT", 0, '"', 43},
		{"43_f", 43, 'f', 45},
		{"45_o", 45, 'o', 45},
		{"45_o", 45, 'o', 45},
		{"45_QUOT", 45, '"', 46},

		// "\"\\"
		{"0_QUOT", 0, '"', 43},
		{"43_BSOL", 43, '\\', 44},
		{"44_QUOT", 44, '"', 45},
		{"45_BSOL", 45, '\\', 44},
		{"44_BSOL", 44, '\\', 45},
		{"45_QUOT", 45, '"', 46},

		// /[a-z]+/
		{"0_SOL", 0, '/', 47},
		{"47_LSQB", 47, '[', 49},
		{"49_a", 49, 'a', 49},
		{"49_MINUS", 49, '-', 49},
		{"49_z", 49, 'z', 49},
		{"49_RSQB", 49, ']', 49},
		{"49_PLUS", 49, '+', 49},
		{"49_SOL", 49, '/', 50},

		// /\{\}/
		{"0_SOL", 0, '/', 47},
		{"47_BSOL", 47, '\\', 48},
		{"48_LCUB", 48, '{', 49},
		{"49_BSOL", 49, '\\', 48},
		{"48_RCUB", 48, '}', 49},
		{"49_SOL", 49, '/', 50},

		// // comment
		{"0_SOL", 0, '/', 47},
		{"47_SOL", 47, '/', 51},
		{"51_Space", 51, ' ', 51},
		{"51_c", 51, 'c', 51},
		{"51_o", 51, 'o', 51},
		{"51_m", 51, 'm', 51},
		{"51_m", 51, 'm', 51},
		{"51_e", 51, 'e', 51},
		{"51_n", 51, 'n', 51},
		{"51_t", 51, 't', 51},

		// /* foo*bar */
		{"0_SOL", 0, '/', 47},
		{"47_AST", 47, '*', 52},
		{"52_Space", 52, ' ', 52},
		{"52_f", 52, 'f', 52},
		{"52_o", 52, 'o', 52},
		{"52_o", 52, 'o', 52},
		{"52_AST", 52, '*', 53},
		{"53_b", 53, 'b', 52},
		{"52_a", 52, 'a', 52},
		{"52_r", 52, 'r', 52},
		{"52_Space", 52, ' ', 52},
		{"52_AST", 52, '*', 53},
		{"53_SOL", 53, '/', 54},

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
		name     string
		filename string
	}{
		{
			name:     "Success",
			filename: "../fixture/please.grammar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)

			defer func() {
				assert.NoError(t, f.Close())
			}()

			lex, err := New(tc.filename, f)
			assert.NoError(t, err)

			for token, err := lex.NextToken(); err != io.EOF; token, err = lex.NextToken() {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			}
		})
	}
}
