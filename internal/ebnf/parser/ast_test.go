package parser

import (
	"testing"

	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"
)

type EqualTest struct {
	rhs      Node
	expected bool
}

func TestTraverse(t *testing.T) {
	root := &Grammar{
		Name: "pascal",
		Decls: []Decl{
			&StringTokenDecl{
				Name:  "SEMI",
				Value: ";",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     2,
					Column:   1,
				},
			},
			&RegexTokenDecl{
				Name:  "ID",
				Regex: "[A-Za-z][0-9A-Za-z]*",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     3,
					Column:   1,
				},
			},
			&PrecedenceDecl{
				Associativity: lr.LEFT,
				Handles: []PrecedenceHandle{
					&TerminalHandle{
						Terminal: "*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     4,
							Column:   8,
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     4,
					Column:   1,
				},
			},
			&PrecedenceDecl{
				Associativity: lr.LEFT,
				Handles: []PrecedenceHandle{
					&TerminalHandle{
						Terminal: "+",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     5,
							Column:   8,
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     5,
					Column:   1,
				},
			},
			&RuleDecl{
				LHS: "expr",
				RHS: &ConcatRHS{
					Ops: []RHS{
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   0,
								Line:     6,
								Column:   10,
							},
						},
						&AltRHS{
							Ops: []RHS{
								&TerminalRHS{
									Terminal: "+",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     6,
										Column:   16,
									},
								},
								&TerminalRHS{
									Terminal: "+",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     6,
										Column:   20,
									},
								},
							},
						},
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   0,
								Line:     6,
								Column:   24,
							},
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     6,
					Column:   1,
				},
			},
		},
		Position: &lexer.Position{
			Filename: "program.code",
			Offset:   0,
			Line:     1,
			Column:   1,
		},
	}

	tests := []struct {
		name           string
		n              Node
		order          generic.TraverseOrder
		expectedVisits []string
	}{
		{
			name:  "VLR",
			n:     root,
			order: generic.VLR,
			expectedVisits: []string{
				`Grammar::pascal <program.code:1:1>`,
				`TokenDecl::SEMI=";" <program.code:2:1>`,
				`TokenDecl::ID=/[A-Za-z][0-9A-Za-z]*/ <program.code:3:1>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::* <program.code:4:8>] <program.code:4:1>`,
				`TerminalHandle::* <program.code:4:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::+ <program.code:5:8>] <program.code:5:1>`,
				`TerminalHandle::+ <program.code:5:8>`,
				`RuleDecl::expr → ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`NonTerminalRHS::expr <program.code:6:10>`,
				`AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20>`,
				`TerminalRHS::+ <program.code:6:16>`,
				`TerminalRHS::+ <program.code:6:20>`,
				"NonTerminalRHS::expr <program.code:6:24>",
			},
		},
		{
			name:  "VRL",
			n:     root,
			order: generic.VRL,
			expectedVisits: []string{
				`Grammar::pascal <program.code:1:1>`,
				`RuleDecl::expr → ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				"NonTerminalRHS::expr <program.code:6:24>",
				`AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20>`,
				`TerminalRHS::+ <program.code:6:20>`,
				`TerminalRHS::+ <program.code:6:16>`,
				`NonTerminalRHS::expr <program.code:6:10>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::+ <program.code:5:8>] <program.code:5:1>`,
				`TerminalHandle::+ <program.code:5:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::* <program.code:4:8>] <program.code:4:1>`,
				`TerminalHandle::* <program.code:4:8>`,
				`TokenDecl::ID=/[A-Za-z][0-9A-Za-z]*/ <program.code:3:1>`,
				`TokenDecl::SEMI=";" <program.code:2:1>`,
			},
		},
		{
			name:  "LRV",
			n:     root,
			order: generic.LRV,
			expectedVisits: []string{
				`TokenDecl::SEMI=";" <program.code:2:1>`,
				`TokenDecl::ID=/[A-Za-z][0-9A-Za-z]*/ <program.code:3:1>`,
				`TerminalHandle::* <program.code:4:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::* <program.code:4:8>] <program.code:4:1>`,
				`TerminalHandle::+ <program.code:5:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::+ <program.code:5:8>] <program.code:5:1>`,
				`NonTerminalRHS::expr <program.code:6:10>`,
				`TerminalRHS::+ <program.code:6:16>`,
				`TerminalRHS::+ <program.code:6:20>`,
				`AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20>`,
				"NonTerminalRHS::expr <program.code:6:24>",
				`ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`RuleDecl::expr → ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				"Grammar::pascal <program.code:1:1>",
			},
		},
		{
			name:  "RLV",
			n:     root,
			order: generic.RLV,
			expectedVisits: []string{
				"NonTerminalRHS::expr <program.code:6:24>",
				`TerminalRHS::+ <program.code:6:20>`,
				`TerminalRHS::+ <program.code:6:16>`,
				`AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20>`,
				`NonTerminalRHS::expr <program.code:6:10>`,
				`ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`RuleDecl::expr → ConcatRHS::NonTerminalRHS::expr <program.code:6:10> AltRHS::TerminalRHS::+ <program.code:6:16> "|" TerminalRHS::+ <program.code:6:20> NonTerminalRHS::expr <program.code:6:24>`,
				`TerminalHandle::+ <program.code:5:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::+ <program.code:5:8>] <program.code:5:1>`,
				`TerminalHandle::* <program.code:4:8>`,
				`PrecedenceDecl::LEFT=[TerminalHandle::* <program.code:4:8>] <program.code:4:1>`,
				`TokenDecl::ID=/[A-Za-z][0-9A-Za-z]*/ <program.code:3:1>`,
				`TokenDecl::SEMI=";" <program.code:2:1>`,
				"Grammar::pascal <program.code:1:1>",
			},
		},
		{
			name:           "InvalidOrder",
			n:              root,
			order:          generic.RVL,
			expectedVisits: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var visits []string
			Traverse(tc.n, tc.order, func(n Node) bool {
				visits = append(visits, n.String())
				return true
			})

			assert.Equal(t, tc.expectedVisits, visits)
		})
	}
}

func TestGrammar(t *testing.T) {
	tests := []struct {
		name           string
		n              *Grammar
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &Grammar{
				Name: "pascal",
				Decls: []Decl{
					&StringTokenDecl{
						Name:  "SEMI",
						Value: ";",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     2,
							Column:   1,
						},
					},
					&RegexTokenDecl{
						Name:  "ID",
						Regex: "[A-Za-z][0-9A-Za-z]*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     3,
							Column:   1,
						},
					},
					&PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     4,
									Column:   8,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     4,
							Column:   1,
						},
					},
					&PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "+",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     5,
									Column:   8,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     5,
							Column:   1,
						},
					},
					&RuleDecl{
						LHS: "expr",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     6,
										Column:   10,
									},
								},
								&TerminalRHS{
									Terminal: "+",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     6,
										Column:   16,
									},
								},
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     6,
										Column:   20,
									},
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     6,
							Column:   1,
						},
					},
					&RuleDecl{
						LHS: "expr",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     7,
										Column:   10,
									},
								},
								&TerminalRHS{
									Terminal: "*",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     7,
										Column:   16,
									},
								},
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   0,
										Line:     7,
										Column:   20,
									},
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     7,
							Column:   1,
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedString: `Grammar::pascal <program.code:1:1>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &Grammar{
						Name: "pascal",
						Decls: []Decl{
							&RuleDecl{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   10,
											},
										},
										&TerminalRHS{
											Terminal: "-",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   16,
											},
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   20,
											},
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     6,
									Column:   1,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: false,
				},
				{
					rhs: &Grammar{
						Name: "pascal",
						Decls: []Decl{
							&StringTokenDecl{
								Name:  "SEMI",
								Value: ";",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     2,
									Column:   1,
								},
							},
							&RegexTokenDecl{
								Name:  "ID",
								Regex: "[A-Za-z][0-9A-Za-z]*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     3,
									Column:   1,
								},
							},
							&PrecedenceDecl{
								Associativity: lr.LEFT,
								Handles: []PrecedenceHandle{
									&TerminalHandle{
										Terminal: "*",
										Position: &lexer.Position{
											Filename: "program.code",
											Offset:   0,
											Line:     4,
											Column:   8,
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     4,
									Column:   1,
								},
							},
							&PrecedenceDecl{
								Associativity: lr.LEFT,
								Handles: []PrecedenceHandle{
									&TerminalHandle{
										Terminal: "+",
										Position: &lexer.Position{
											Filename: "program.code",
											Offset:   0,
											Line:     5,
											Column:   8,
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     5,
									Column:   1,
								},
							},
							&RuleDecl{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   10,
											},
										},
										&TerminalRHS{
											Terminal: "-",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   16,
											},
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   20,
											},
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     6,
									Column:   1,
								},
							},
							&RuleDecl{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   10,
											},
										},
										&TerminalRHS{
											Terminal: "/",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   16,
											},
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   20,
											},
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     7,
									Column:   1,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: false,
				},
				{
					rhs: &Grammar{
						Name: "pascal",
						Decls: []Decl{
							&StringTokenDecl{
								Name:  "SEMI",
								Value: ";",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     2,
									Column:   1,
								},
							},
							&RegexTokenDecl{
								Name:  "ID",
								Regex: "[A-Za-z][0-9A-Za-z]*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     3,
									Column:   1,
								},
							},
							&PrecedenceDecl{
								Associativity: lr.LEFT,
								Handles: []PrecedenceHandle{
									&TerminalHandle{
										Terminal: "*",
										Position: &lexer.Position{
											Filename: "program.code",
											Offset:   0,
											Line:     4,
											Column:   8,
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     4,
									Column:   1,
								},
							},
							&PrecedenceDecl{
								Associativity: lr.LEFT,
								Handles: []PrecedenceHandle{
									&TerminalHandle{
										Terminal: "+",
										Position: &lexer.Position{
											Filename: "program.code",
											Offset:   0,
											Line:     5,
											Column:   8,
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     5,
									Column:   1,
								},
							},
							&RuleDecl{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   10,
											},
										},
										&TerminalRHS{
											Terminal: "+",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   16,
											},
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     6,
												Column:   20,
											},
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     6,
									Column:   1,
								},
							},
							&RuleDecl{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   10,
											},
										},
										&TerminalRHS{
											Terminal: "*",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   16,
											},
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
											Position: &lexer.Position{
												Filename: "program.code",
												Offset:   0,
												Line:     7,
												Column:   20,
											},
										},
									},
								},
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   0,
									Line:     7,
									Column:   1,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for i, child := range tc.n.Children() {
				assert.Equal(t, tc.n.Decls[i], child)
			}

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}
		})
	}
}

func TestGrammar_DOT(t *testing.T) {
	tests := []struct {
		name        string
		n           *Grammar
		expectedDOT string
	}{
		{
			name: "OK",
			n: &Grammar{
				Name: "pascal",
				Decls: []Decl{
					&StringTokenDecl{
						Name:  "SEMI",
						Value: ";",
					},
					&RegexTokenDecl{
						Name:  "ID",
						Regex: "[A-Za-z][0-9A-Za-z]*",
					},
					&PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
							},
							&TerminalHandle{
								Terminal: "/",
							},
						},
					},
					&PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "+",
							},
							&TerminalHandle{
								Terminal: "-",
							},
						},
					},
					&PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&ProductionHandle{
								LHS: "expr",
								RHS: &ConcatRHS{
									Ops: []RHS{
										&NonTerminalRHS{
											NonTerminal: "expr",
										},
										&NonTerminalRHS{
											NonTerminal: "logop",
										},
										&NonTerminalRHS{
											NonTerminal: "expr",
										},
									},
								},
							},
						},
					},
					&RuleDecl{
						LHS: "program",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&TerminalRHS{
									Terminal: "PROGRAM",
								},
								&TerminalRHS{
									Terminal: "ID",
								},
								&TerminalRHS{
									Terminal: ";",
								},
								&OptRHS{
									Op: &NonTerminalRHS{
										NonTerminal: "block",
									},
								},
							},
						},
					},
					&RuleDecl{
						LHS: "block",
						RHS: &AltRHS{
							Ops: []RHS{
								&ConcatRHS{
									Ops: []RHS{
										&StarRHS{
											Op: &NonTerminalRHS{
												NonTerminal: "decl",
											},
										},
										&PlusRHS{
											Op: &NonTerminalRHS{
												NonTerminal: "stmt",
											},
										},
									},
								},
								&EmptyRHS{},
							},
						},
					},
					&RuleDecl{
						LHS: "expr",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&NonTerminalRHS{
									NonTerminal: "expr",
								},
								&AltRHS{
									Ops: []RHS{
										&TerminalRHS{
											Terminal: "+",
										},
										&TerminalRHS{
											Terminal: "-",
										},
										&TerminalRHS{
											Terminal: "*",
										},
										&TerminalRHS{
											Terminal: "/",
										},
									},
								},
								&NonTerminalRHS{
									NonTerminal: "expr",
								},
							},
						},
					},
				},
			},
			expectedDOT: `strict digraph "AST" {
  concentrate=false;
  node [];

  1 [label="Grammar::pascal", color=gold, style=filled, shape=square];
  2 [label="TokenDecl::SEMI", color=skyblue, style=filled, shape=box];
  3 [label="TokenDecl::ID", color=skyblue, style=filled, shape=box];
  4 [label="PrecedenceDecl::LEFT", color=burlywood, style=filled, shape=box];
  5 [label="TerminalHandle::*", color=burlywood, style=filled, shape=box];
  6 [label="TerminalHandle::/", color=burlywood, style=filled, shape=box];
  7 [label="PrecedenceDecl::LEFT", color=burlywood, style=filled, shape=box];
  8 [label="TerminalHandle::+", color=burlywood, style=filled, shape=box];
  9 [label="TerminalHandle::-", color=burlywood, style=filled, shape=box];
  10 [label="PrecedenceDecl::LEFT", color=burlywood, style=filled, shape=box];
  11 [label="ProductionHandle::expr →", color=burlywood, style=filled, shape=box];
  12 [label="CONCAT", color=lavender, style=filled, shape=box];
  13 [label="NonTerminal::expr", color=turquoise, style=filled, shape=box];
  14 [label="NonTerminal::logop", color=turquoise, style=filled, shape=box];
  15 [label="NonTerminal::expr", color=turquoise, style=filled, shape=box];
  16 [label="RuleDecl::program →", color=lightpink, style=filled, shape=box];
  17 [label="CONCAT", color=lavender, style=filled, shape=box];
  18 [label="Terminal::PROGRAM", color=springgreen, style=filled, shape=oval];
  19 [label="Terminal::ID", color=springgreen, style=filled, shape=oval];
  20 [label="Terminal::;", color=springgreen, style=filled, shape=oval];
  21 [label="ZERO OR ONE", color=lavender, style=filled, shape=box];
  22 [label="NonTerminal::block", color=turquoise, style=filled, shape=box];
  23 [label="RuleDecl::block →", color=lightpink, style=filled, shape=box];
  24 [label="ALT", color=lavender, style=filled, shape=box];
  25 [label="CONCAT", color=lavender, style=filled, shape=box];
  26 [label="ZERO OR MORE", color=lavender, style=filled, shape=box];
  27 [label="NonTerminal::decl", color=turquoise, style=filled, shape=box];
  28 [label="ONE OR MORE", color=lavender, style=filled, shape=box];
  29 [label="NonTerminal::stmt", color=turquoise, style=filled, shape=box];
  30 [label="ε", color=violet, style=filled, shape=circle];
  31 [label="RuleDecl::expr →", color=lightpink, style=filled, shape=box];
  32 [label="CONCAT", color=lavender, style=filled, shape=box];
  33 [label="NonTerminal::expr", color=turquoise, style=filled, shape=box];
  34 [label="ALT", color=lavender, style=filled, shape=box];
  35 [label="Terminal::+", color=springgreen, style=filled, shape=oval];
  36 [label="Terminal::-", color=springgreen, style=filled, shape=oval];
  37 [label="Terminal::*", color=springgreen, style=filled, shape=oval];
  38 [label="Terminal::/", color=springgreen, style=filled, shape=oval];
  39 [label="NonTerminal::expr", color=turquoise, style=filled, shape=box];

  1 -> 2 [];
  1 -> 3 [];
  1 -> 4 [];
  1 -> 7 [];
  1 -> 10 [];
  1 -> 16 [];
  1 -> 23 [];
  1 -> 31 [];
  4 -> 5 [];
  4 -> 6 [];
  7 -> 8 [];
  7 -> 9 [];
  10 -> 11 [];
  11 -> 12 [];
  12 -> 13 [];
  12 -> 14 [];
  12 -> 15 [];
  16 -> 17 [];
  17 -> 18 [];
  17 -> 19 [];
  17 -> 20 [];
  17 -> 21 [];
  21 -> 22 [];
  23 -> 24 [];
  24 -> 25 [];
  24 -> 30 [];
  25 -> 26 [];
  25 -> 28 [];
  26 -> 27 [];
  28 -> 29 [];
  31 -> 32 [];
  32 -> 33 [];
  32 -> 34 [];
  32 -> 39 [];
  34 -> 35 [];
  34 -> 36 [];
  34 -> 37 [];
  34 -> 38 [];
}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedDOT, tc.n.DOT())
		})
	}
}

func TestStringTokenDecl(t *testing.T) {
	tests := []struct {
		name           string
		n              *StringTokenDecl
		expectedString string
		equalTests     []EqualTest
		expectedPos    *lexer.Position
	}{
		{
			name: "OK",
			n: &StringTokenDecl{
				Name:  "SEMI",
				Value: ";",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedString: `TokenDecl::SEMI=";" <program.code:1:1>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &StringTokenDecl{
						Name:  "SEMI",
						Value: ";",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.decl()
		})
	}
}

func TestRegexTokenDecl(t *testing.T) {
	tests := []struct {
		name           string
		n              *RegexTokenDecl
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &RegexTokenDecl{
				Name:  "ID",
				Regex: "[A-Za-z][0-9A-Za-z]*",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedString: `TokenDecl::ID=/[A-Za-z][0-9A-Za-z]*/ <program.code:1:1>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &RegexTokenDecl{
						Name:  "ID",
						Regex: "[A-Za-z][0-9A-Za-z]*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.decl()
		})
	}
}

func TestPrecedenceDecl(t *testing.T) {
	tests := []struct {
		name           string
		n              *PrecedenceDecl
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &PrecedenceDecl{
				Associativity: lr.LEFT,
				Handles: []PrecedenceHandle{
					&TerminalHandle{
						Terminal: "*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   7,
							Line:     1,
							Column:   8,
						},
					},
					&TerminalHandle{
						Terminal: "/",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   11,
							Line:     1,
							Column:   12,
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedString: `PrecedenceDecl::LEFT=[TerminalHandle::* <program.code:1:8> TerminalHandle::/ <program.code:1:12>] <program.code:1:1>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &PrecedenceDecl{
						Associativity: lr.RIGHT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
							&TerminalHandle{
								Terminal: "/",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: false,
				},
				{
					rhs: &PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: false,
				},
				{
					rhs: &PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
							&TerminalHandle{
								Terminal: "%",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: false,
				},
				{
					rhs: &PrecedenceDecl{
						Associativity: lr.LEFT,
						Handles: []PrecedenceHandle{
							&TerminalHandle{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   7,
									Line:     1,
									Column:   8,
								},
							},
							&TerminalHandle{
								Terminal: "/",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   11,
									Line:     1,
									Column:   12,
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for i, child := range tc.n.Children() {
				assert.Equal(t, tc.n.Handles[i], child)
			}

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.decl()
		})
	}
}

func TestTerminalHandle(t *testing.T) {
	tests := []struct {
		name           string
		n              *TerminalHandle
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &TerminalHandle{
				Terminal: "*",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   7,
					Line:     1,
					Column:   8,
				},
			},
			expectedString: `TerminalHandle::* <program.code:1:8>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   7,
				Line:     1,
				Column:   8,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &TerminalHandle{
						Terminal: "*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   7,
							Line:     1,
							Column:   8,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.precedenceHandle()
		})
	}
}

func TestProductionHandle(t *testing.T) {
	tests := []struct {
		name           string
		n              *ProductionHandle
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &ProductionHandle{
				LHS: "expr",
				RHS: &ConcatRHS{
					Ops: []RHS{
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   15,
								Line:     1,
								Column:   16,
							},
						},
						&TerminalRHS{
							Terminal: "*",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   21,
								Line:     1,
								Column:   22,
							},
						},
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   25,
								Line:     1,
								Column:   26,
							},
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   7,
					Line:     1,
					Column:   8,
				},
			},
			expectedString: `ProductionHandle::expr → ConcatRHS::NonTerminalRHS::expr <program.code:1:16> TerminalRHS::* <program.code:1:22> NonTerminalRHS::expr <program.code:1:26>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   7,
				Line:     1,
				Column:   8,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &ProductionHandle{
						LHS: "expr",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   15,
										Line:     1,
										Column:   16,
									},
								},
								&TerminalRHS{
									Terminal: "*",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   21,
										Line:     1,
										Column:   22,
									},
								},
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   25,
										Line:     1,
										Column:   26,
									},
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   7,
							Line:     1,
							Column:   8,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			expectedChildren := []Node{tc.n.RHS}
			assert.Equal(t, expectedChildren, tc.n.Children())

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.precedenceHandle()
		})
	}
}

func TestRuleDecl(t *testing.T) {
	tests := []struct {
		name           string
		n              *RuleDecl
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &RuleDecl{
				LHS: "expr",
				RHS: &ConcatRHS{
					Ops: []RHS{
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   9,
								Line:     1,
								Column:   10,
							},
						},
						&TerminalRHS{
							Terminal: "+",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   15,
								Line:     1,
								Column:   16,
							},
						},
						&NonTerminalRHS{
							NonTerminal: "expr",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   19,
								Line:     1,
								Column:   20,
							},
						},
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			expectedString: `RuleDecl::expr → ConcatRHS::NonTerminalRHS::expr <program.code:1:10> TerminalRHS::+ <program.code:1:16> NonTerminalRHS::expr <program.code:1:20>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   0,
				Line:     1,
				Column:   1,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &RuleDecl{
						LHS: "expr",
						RHS: &ConcatRHS{
							Ops: []RHS{
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   9,
										Line:     1,
										Column:   10,
									},
								},
								&TerminalRHS{
									Terminal: "+",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   15,
										Line:     1,
										Column:   16,
									},
								},
								&NonTerminalRHS{
									NonTerminal: "expr",
									Position: &lexer.Position{
										Filename: "program.code",
										Offset:   19,
										Line:     1,
										Column:   20,
									},
								},
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   0,
							Line:     1,
							Column:   1,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			expectedChildren := []Node{tc.n.RHS}
			assert.Equal(t, expectedChildren, tc.n.Children())

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.decl()
		})
	}
}

func TestConcatRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *ConcatRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &ConcatRHS{
				Ops: []RHS{
					&TerminalRHS{
						Terminal: "PROGRAM",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   9,
							Line:     1,
							Column:   10,
						},
					},
					&TerminalRHS{
						Terminal: "ID",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   19,
							Line:     1,
							Column:   20,
						},
					},
				},
			},
			expectedString: `ConcatRHS::TerminalRHS::PROGRAM <program.code:1:10> TerminalRHS::ID <program.code:1:20>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   9,
				Line:     1,
				Column:   10,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &ConcatRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "ID",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   19,
									Line:     1,
									Column:   20,
								},
							},
						},
					},
					expected: false,
				},
				{
					rhs: &ConcatRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "PROG",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   9,
									Line:     1,
									Column:   10,
								},
							},
							&TerminalRHS{
								Terminal: "ID",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   19,
									Line:     1,
									Column:   20,
								},
							},
						},
					},
					expected: false,
				},
				{
					rhs: &ConcatRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "PROGRAM",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   9,
									Line:     1,
									Column:   10,
								},
							},
							&TerminalRHS{
								Terminal: "ID",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   19,
									Line:     1,
									Column:   20,
								},
							},
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for i, child := range tc.n.Children() {
				assert.Equal(t, tc.n.Ops[i], child)
			}

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.rhs()
		})
	}
}

func TestAltRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *AltRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &AltRHS{
				Ops: []RHS{
					&TerminalRHS{
						Terminal: "*",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   9,
							Line:     1,
							Column:   10,
						},
					},
					&TerminalRHS{
						Terminal: "+",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   15,
							Line:     1,
							Column:   16,
						},
					},
				},
			},
			expectedString: `AltRHS::TerminalRHS::* <program.code:1:10> "|" TerminalRHS::+ <program.code:1:16>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   9,
				Line:     1,
				Column:   10,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &AltRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   9,
									Line:     1,
									Column:   10,
								},
							},
						},
					},
					expected: false,
				},
				{
					rhs: &AltRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   9,
									Line:     1,
									Column:   10,
								},
							},
							&TerminalRHS{
								Terminal: "-",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   15,
									Line:     1,
									Column:   16,
								},
							},
						},
					},
					expected: false,
				},
				{
					rhs: &AltRHS{
						Ops: []RHS{
							&TerminalRHS{
								Terminal: "*",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   9,
									Line:     1,
									Column:   10,
								},
							},
							&TerminalRHS{
								Terminal: "+",
								Position: &lexer.Position{
									Filename: "program.code",
									Offset:   15,
									Line:     1,
									Column:   16,
								},
							},
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for i, child := range tc.n.Children() {
				assert.Equal(t, tc.n.Ops[i], child)
			}

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.rhs()
		})
	}
}

func TestOptRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *OptRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &OptRHS{
				Op: &NonTerminalRHS{
					NonTerminal: "decl",
					Position: &lexer.Position{
						Filename: "program.code",
						Offset:   7,
						Line:     1,
						Column:   8,
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
			},
			expectedString: `OptRHS::NonTerminalRHS::decl <program.code:1:8>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   6,
				Line:     1,
				Column:   7,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &OptRHS{
						Op: &NonTerminalRHS{
							NonTerminal: "decl",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   7,
								Line:     1,
								Column:   8,
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   6,
							Line:     1,
							Column:   7,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			expectedChildren := []Node{tc.n.Op}
			assert.Equal(t, expectedChildren, tc.n.Children())

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.rhs()
		})
	}
}

func TestStarRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *StarRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &StarRHS{
				Op: &NonTerminalRHS{
					NonTerminal: "decl",
					Position: &lexer.Position{
						Filename: "program.code",
						Offset:   7,
						Line:     1,
						Column:   8,
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
			},
			expectedString: `StarRHS::NonTerminalRHS::decl <program.code:1:8>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   6,
				Line:     1,
				Column:   7,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &StarRHS{
						Op: &NonTerminalRHS{
							NonTerminal: "decl",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   7,
								Line:     1,
								Column:   8,
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   6,
							Line:     1,
							Column:   7,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			expectedChildren := []Node{tc.n.Op}
			assert.Equal(t, expectedChildren, tc.n.Children())

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.rhs()
		})
	}
}

func TestPlusRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *PlusRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &PlusRHS{
				Op: &NonTerminalRHS{
					NonTerminal: "decl",
					Position: &lexer.Position{
						Filename: "program.code",
						Offset:   7,
						Line:     1,
						Column:   8,
					},
				},
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
			},
			expectedString: `PlusRHS::NonTerminalRHS::decl <program.code:1:8>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   6,
				Line:     1,
				Column:   7,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &PlusRHS{
						Op: &NonTerminalRHS{
							NonTerminal: "decl",
							Position: &lexer.Position{
								Filename: "program.code",
								Offset:   7,
								Line:     1,
								Column:   8,
							},
						},
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   6,
							Line:     1,
							Column:   7,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			expectedChildren := []Node{tc.n.Op}
			assert.Equal(t, expectedChildren, tc.n.Children())

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.rhs()
		})
	}
}

func TestNonTerminalRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *NonTerminalRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &NonTerminalRHS{
				NonTerminal: "expr",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   9,
					Line:     1,
					Column:   10,
				},
			},
			expectedString: `NonTerminalRHS::expr <program.code:1:10>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   9,
				Line:     1,
				Column:   10,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &NonTerminalRHS{
						NonTerminal: "expr",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   9,
							Line:     1,
							Column:   10,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.rhs()
		})
	}
}

func TestTerminalRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *TerminalRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name: "OK",
			n: &TerminalRHS{
				Terminal: "++",
				Position: &lexer.Position{
					Filename: "program.code",
					Offset:   9,
					Line:     1,
					Column:   10,
				},
			},
			expectedString: `TerminalRHS::++ <program.code:1:10>`,
			expectedPos: &lexer.Position{
				Filename: "program.code",
				Offset:   9,
				Line:     1,
				Column:   10,
			},
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs: &TerminalRHS{
						Terminal: "++",
						Position: &lexer.Position{
							Filename: "program.code",
							Offset:   9,
							Line:     1,
							Column:   10,
						},
					},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.rhs()
		})
	}
}

func TestEmptyRHS(t *testing.T) {
	tests := []struct {
		name           string
		n              *EmptyRHS
		expectedString string
		expectedPos    *lexer.Position
		equalTests     []EqualTest
	}{
		{
			name:           "OK",
			n:              &EmptyRHS{},
			expectedString: `EmptyRHS::ε`,
			expectedPos:    nil,
			equalTests: []EqualTest{
				{
					rhs:      nil,
					expected: false,
				},
				{
					rhs:      &EmptyRHS{},
					expected: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.n.String())
			assert.True(t, equalPositions(tc.expectedPos, tc.n.Pos()))

			for _, equalTest := range tc.equalTests {
				assert.Equal(t, equalTest.expected, tc.n.Equal(equalTest.rhs))
			}

			tc.n.leaf()
			tc.n.rhs()
		})
	}
}
