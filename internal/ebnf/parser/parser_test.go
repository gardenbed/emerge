package parser

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"
)

// MockLexer is an implementation of lexer.Lexer for testing purposes.
type MockLexer struct {
	NextTokenIndex int
	NextTokenMocks []NextTokenMock
}

type NextTokenMock struct {
	OutToken lexer.Token
	OutError error
}

func (m *MockLexer) NextToken() (lexer.Token, error) {
	i := m.NextTokenIndex
	m.NextTokenIndex++
	return m.NextTokenMocks[i].OutToken, m.NextTokenMocks[i].OutError
}

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
			par, err := New(tc.filename, tc.src)

			if tc.expectedError == "" {
				assert.NotNil(t, par)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, par)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name                 string
		p                    *Parser
		tokenF               parser.TokenFunc
		prodF                ProductionFunc
		expectedErrorStrings []string
	}{
		{
			name: "EmptyString",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						{OutError: io.EOF},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return nil },
			prodF:  func(int) error { return nil },
			expectedErrorStrings: []string{
				`unexpected string "": no action exists in the parsing table for ACTION[0, $]`,
			},
		},
		{
			name: "First_NextToken_Fails",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						{OutError: errors.New("cannot read rune")},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return nil },
			prodF:  func(int) error { return nil },
			expectedErrorStrings: []string{
				`cannot read rune`,
			},
		},
		{
			name: "Second_NextToken_Fails",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "a",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// EOF
						{OutError: errors.New("input stream failed")},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return nil },
			prodF:  func(int) error { return nil },
			expectedErrorStrings: []string{
				`input stream failed`,
			},
		},
		{
			name: "Invalid_Input",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return nil },
			prodF:  func(int) error { return nil },
			expectedErrorStrings: []string{
				`test:1:8: unexpected string "name": no action exists in the parsing table for ACTION[0, "IDENT"]`,
			},
		},
		{
			name: "TokenFuncError",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return errors.New("invalid semantic") },
			prodF:  func(int) error { return nil },
			expectedErrorStrings: []string{
				`test:1:1: invalid semantic`,
			},
		},
		{
			name: "ProductionFuncError",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			tokenF: func(*lexer.Token) error { return nil },
			prodF:  func(int) error { return errors.New("invalid semantic") },
			expectedErrorStrings: []string{
				`invalid semantic`,
			},
		},
		{
			name: "Success",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			tokenF:               func(*lexer.Token) error { return nil },
			prodF:                func(int) error { return nil },
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.p.Parse(tc.tokenF, tc.prodF)

			if len(tc.expectedErrorStrings) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestParser_ParseAndBuildAST(t *testing.T) {
	tests := []struct {
		name                 string
		p                    *Parser
		expectedAST          parser.Node
		expectedErrorStrings []string
	}{
		{
			name: "EmptyString",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						{OutError: io.EOF},
					},
				},
			},
			expectedAST: nil,
			expectedErrorStrings: []string{
				`unexpected string "": no action exists in the parsing table for ACTION[0, $]`,
			},
		},
		{
			name: "Success",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			expectedAST: &parser.InternalNode{
				NonTerminal: "grammar",
				Production: &grammar.Production{
					Head: "grammar",
					Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("name"), grammar.NonTerminal("decls")},
				},
				Children: []parser.Node{
					&parser.InternalNode{
						NonTerminal: "name",
						Production: &grammar.Production{
							Head: "name",
							Body: grammar.String[grammar.Symbol]{grammar.Terminal("grammar"), grammar.Terminal("IDENT"), grammar.NonTerminal("semi_opt")},
						},
						Children: []parser.Node{
							&parser.LeafNode{
								Terminal: "grammar",
								Lexeme:   "grammar",
								Position: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
							&parser.LeafNode{
								Terminal: "IDENT",
								Lexeme:   "name",
								Position: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
							&parser.InternalNode{
								NonTerminal: "semi_opt",
								Production: &grammar.Production{
									Head: "semi_opt",
									Body: grammar.String[grammar.Symbol]{grammar.Terminal(";")},
								},
								Children: []parser.Node{
									&parser.LeafNode{
										Terminal: ";",
										Lexeme:   ";",
										Position: lexer.Position{
											Filename: "test",
											Offset:   11,
											Line:     1,
											Column:   12,
										},
									},
								},
							},
						},
					},
					&parser.InternalNode{
						NonTerminal: "decls",
						Production: &grammar.Production{
							Head: "decls",
							Body: grammar.E,
						},
					},
				},
			},
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ast, err := tc.p.ParseAndBuildAST()

			if len(tc.expectedErrorStrings) == 0 {
				assert.True(t, ast.Equal(tc.expectedAST))
				assert.NoError(t, err)
			} else {
				assert.Nil(t, ast)
				assert.Error(t, err)
				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestParser_ParseAndEvaluate(t *testing.T) {
	tests := []struct {
		name                 string
		p                    *Parser
		eval                 EvaluateFunc
		expectedValue        *lr.Value
		expectedErrorStrings []string
	}{
		{
			name: "EmptyString",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						{OutError: io.EOF},
					},
				},
			},
			eval:          func(int, []*lr.Value) (any, error) { return nil, nil },
			expectedValue: nil,
			expectedErrorStrings: []string{
				`unexpected string "": no action exists in the parsing table for ACTION[0, $]`,
			},
		},
		{
			name: "EvaluateFuncError",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			eval:          func(int, []*lr.Value) (any, error) { return nil, errors.New("invalid semantic") },
			expectedValue: nil,
			expectedErrorStrings: []string{
				`invalid semantic`,
			},
		},
		{
			name: "Success",
			p: &Parser{
				L: &MockLexer{
					NextTokenMocks: []NextTokenMock{
						// First token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("grammar"),
								Lexeme:   "grammar",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   0,
									Line:     1,
									Column:   1,
								},
							},
						},
						// Second token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal("IDENT"),
								Lexeme:   "name",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						// Third token
						{
							OutToken: lexer.Token{
								Terminal: grammar.Terminal(";"),
								Lexeme:   ";",
								Pos: lexer.Position{
									Filename: "test",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						// EOF
						{OutError: io.EOF},
					},
				},
			},
			eval: func(int, []*lr.Value) (any, error) { return &grammar.CFG{}, nil },
			expectedValue: &lr.Value{
				Val: &grammar.CFG{},
				Pos: &lexer.Position{
					Filename: "test",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.p.ParseAndEvaluate(tc.eval)

			if len(tc.expectedErrorStrings) == 0 {
				assert.True(t, reflect.DeepEqual(val, tc.expectedValue))
				assert.NoError(t, err)
			} else {
				assert.Nil(t, val)
				assert.Error(t, err)
				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestParser_ParseAndBuildAST2(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectedError string
	}{
		{
			name:          "Invalid",
			filename:      "../fixture/invalid.grammar",
			expectedError: "lexical error at ../fixture/invalid.grammar:1:1:L",
		},
		{
			name:          "Success",
			filename:      "../fixture/test.grammar",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)
			defer f.Close()

			p, err := New(tc.filename, f)
			assert.NoError(t, err)

			root, err := p.ParseAndBuildAST2()

			if len(tc.expectedError) == 0 {
				assert.NotNil(t, root)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, root)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
