package lexer

import (
	"errors"
	"fmt"
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
			name: "IDENT_g",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "g",
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
			state: 32,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "g",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_gr",
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
			state: 33,
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
			name: "IDENT_gra",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "gra",
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
			state: 34,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "gra",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_gram",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "gram",
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
			state: 35,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "gram",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_gramm",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "gramm",
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
				Lexeme:   "gramm",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_gramma",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "gramma",
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
			state: 37,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "gramma",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   8,
					Line:     1,
					Column:   9,
				},
			},
		},
		{
			name: "IDENT_stmt",
			l: &Lexer{
				in: &mockInputBuffer{
					LexemeMocks: []LexemeMock{
						{
							OutVal: "stmt",
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
			state: 39,
			expectedToken: lexer.Token{
				Terminal: IDENT,
				Lexeme:   "stmt",
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
			state: 40,
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
			state: 61,
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
			state: 65,
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
			state: 66,
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
			state: 69,
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
		state         int
		r             rune
		expectedState int
	}{
		// Whitespaces
		{0, '\t', 1},
		{0, ' ', 1},
		{1, '\t', 1},
		{1, ' ', 1},

		// Newlines
		{0, '\n', 2},
		{0, '\r', 2},
		{2, '\n', 2},
		{2, '\r', 2},

		{0, '=', 3},
		{0, ';', 4},
		{0, '|', 5},
		{0, '(', 6},
		{0, ')', 7},
		{0, '[', 8},
		{0, ']', 9},
		{0, '{', 10},
		{0, '}', 11},
		{10, '{', 12},
		{11, '}', 13},
		{0, '<', 14},
		{0, '>', 15},

		// $
		{0, '$', 16},

		// $STRING
		{16, 'S', 17},
		{17, 'T', 17},
		{17, 'R', 17},
		{17, 'I', 17},
		{17, 'N', 17},
		{17, 'G', 17},

		// @
		{0, '@', 18},

		// @left
		{18, 'l', 19},
		{19, 'e', 20},
		{20, 'f', 21},
		{21, 't', 22},

		// @ right
		{18, 'r', 23},
		{23, 'i', 24},
		{24, 'g', 25},
		{25, 'h', 26},
		{26, 't', 27},

		// @none
		{18, 'n', 28},
		{28, 'o', 29},
		{29, 'n', 30},
		{30, 'e', 31},

		// grammar
		{0, 'g', 32},
		{32, 'r', 33},
		{33, 'a', 34},
		{34, 'm', 35},
		{35, 'm', 36},
		{36, 'a', 37},
		{37, 'r', 38},

		{32, 'o', 39}, // go
		{33, 'e', 39}, // gre
		{34, 'n', 39}, // gran
		{35, 's', 39}, // grams
		{36, 'o', 39}, // grammo
		{37, 't', 39}, // grammat
		{38, 'i', 39}, // grammari
		{39, 'n', 39}, // grammarian

		// name
		{0, 'n', 39},
		{39, 'a', 39},
		{39, 'm', 39},
		{39, 'e', 39},

		// NAME
		{0, 'N', 40},
		{40, 'A', 40},
		{40, 'M', 40},
		{40, 'E', 40},

		// "
		{0, '"', 41},

		// "foo"
		{41, 'f', 41},
		{41, 'o', 41},
		{41, 'o', 41},
		{41, '"', 61},

		// \
		{41, '\\', 42},

		// \" \' \\ \n \r \t
		{42, '"', 43},
		{42, '\'', 43},
		{42, '\\', 43},
		{42, 'n', 43},
		{42, 'r', 43},
		{42, 't', 43},
		{43, 's', 41},
		{43, '\\', 42},
		{43, '"', 61},

		// \xHH
		{42, 'x', 44},
		{44, '7', 45},
		{45, 'F', 46},
		{46, 's', 41},
		{46, '\\', 42},
		{46, '"', 61},

		// \uHHHH
		{42, 'u', 47},
		{47, '0', 48},
		{48, '1', 49},
		{49, 'F', 50},
		{50, 'D', 51},
		{51, 's', 41},
		{51, '\\', 42},
		{51, '"', 61},

		// \UHHHHHH
		{42, 'U', 52},
		{52, '0', 53},
		{53, '0', 54},
		{54, '0', 55},
		{55, '1', 56},
		{56, 'F', 57},
		{57, '4', 58},
		{58, '0', 59},
		{59, '9', 60},
		{60, 's', 41},
		{60, '\\', 42},
		{60, '"', 61},

		// /
		{0, '/', 62},

		// /[a-z]+/
		{62, '[', 64},
		{64, 'a', 64},
		{64, '-', 64},
		{64, 'z', 64},
		{64, ']', 64},
		{64, '+', 64},
		{64, '/', 65},

		// /\{\}/
		{62, '\\', 63},
		{63, '{', 64},
		{64, '\\', 63},
		{63, '}', 64},
		{64, '/', 65},

		// // comment
		{62, '/', 66},
		{66, ' ', 66},
		{66, 'c', 66},
		{66, 'o', 66},
		{66, 'm', 66},
		{66, 'm', 66},
		{66, 'e', 66},
		{66, 'n', 66},
		{66, 't', 66},

		// /* comment */
		{62, '*', 67},
		{67, ' ', 67},
		{67, 'c', 67},
		{67, 'o', 67},
		{67, 'm', 67},
		{67, 'm', 67},
		{67, 'e', 67},
		{67, 'n', 67},
		{67, 't', 67},
		{67, ' ', 67},
		{67, '*', 68},
		{68, '/', 69},

		// ERR
		{0, '0', -1},
		{0, '1', -1},
		{0, '2', -1},
		{0, '3', -1},
		{0, '4', -1},
		{0, '5', -1},
		{0, '6', -1},
		{0, '7', -1},
		{0, '8', -1},
		{0, '9', -1},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d_%c", tc.state, tc.r), func(t *testing.T) {
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
			name:     "EBNF",
			filename: "../fixture/ebnf.grammar",
		},
		{
			name:     "Please",
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
