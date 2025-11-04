package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/parser/combinator"

	"github.com/gardenbed/emerge/internal/char"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name             string
		regex            string
		expectedError    string
		expectedAST      *AST
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name:          "InvalidRegex",
			regex:         "[",
			expectedError: "invalid regular expression: [",
		},
		{
			name:          "InvalidCharRange",
			regex:         "[9-0]",
			expectedError: "invalid character range 9-0",
		},
		{
			name:          "InvalidRepRange",
			regex:         "[0-9]{4,2}",
			expectedError: "invalid repetition range {4,2}",
		},
		{
			name:  "Success_Simple",
			regex: `(a|b)*abb`,
			expectedAST: &AST{
				Root: &Concat{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Star{
									Expr: &Alt{
										Exprs: []Node{
											&Concat{
												Exprs: []Node{
													&Char{Lo: 'a', Hi: 'a', Pos: 1},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{1},
													lastPos:  Poses{1},
												},
											},
											&Concat{
												Exprs: []Node{
													&Char{Lo: 'b', Hi: 'b', Pos: 2},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{2},
													lastPos:  Poses{2},
												},
											},
										},
										comp: &computed{
											nullable: false,
											firstPos: Poses{1, 2},
											lastPos:  Poses{1, 2},
										},
									},
								},
								&Char{Lo: 'a', Hi: 'a', Pos: 3},
								&Char{Lo: 'b', Hi: 'b', Pos: 4},
								&Char{Lo: 'b', Hi: 'b', Pos: 5},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1, 2, 3},
								lastPos:  Poses{5},
							},
						},
						&Char{Lo: endMarker, Hi: endMarker, Pos: 6},
					},
				},
				lastPos: 6,
				posToChar: map[Pos]char.Range{
					1: {'a', 'a'},
					2: {'b', 'b'},
					3: {'a', 'a'},
					4: {'b', 'b'},
					5: {'b', 'b'},
					6: {endMarker, endMarker},
				},
				charToPos: map[char.Range]Poses{
					{'a', 'a'}:             {1, 3},
					{'b', 'b'}:             {2, 4, 5},
					{endMarker, endMarker}: {6},
				},
				follows: map[Pos]Poses{
					1: {1, 2, 3},
					2: {1, 2, 3},
					3: {4},
					4: {5},
					5: {6},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2, 3},
			expectedLastPos:  Poses{6},
		},
		{
			name:  "Success_EscapedChars",
			regex: `\n|\r|\r\n`,
			expectedAST: &AST{
				Root: &Concat{
					Exprs: []Node{
						&Alt{
							Exprs: []Node{
								&Concat{
									Exprs: []Node{
										&Char{Lo: '\n', Hi: '\n', Pos: 1},
									},
								},
								&Alt{
									Exprs: []Node{
										&Concat{
											Exprs: []Node{
												&Char{Lo: '\r', Hi: '\r', Pos: 2},
											},
										},
										&Concat{
											Exprs: []Node{
												&Char{Lo: '\r', Hi: '\r', Pos: 3},
												&Char{Lo: '\n', Hi: '\n', Pos: 4},
											},
										},
									},
								},
							},
						},
						&Char{Lo: endMarker, Hi: endMarker, Pos: 5},
					},
				},
				lastPos: 5,
				posToChar: map[Pos]char.Range{
					1: {'\n', '\n'},
					2: {'\r', '\r'},
					3: {'\r', '\r'},
					4: {'\n', '\n'},
					5: {endMarker, endMarker},
				},
				charToPos: map[char.Range]Poses{
					{'\n', '\n'}:           {1, 4},
					{'\r', '\r'}:           {2, 3},
					{endMarker, endMarker}: {5},
				},
				follows: map[Pos]Poses{
					1: {5},
					2: {5},
					3: {4},
					4: {5},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2, 3},
			expectedLastPos:  Poses{5},
		},
		{
			name:  "Success_CharRanges",
			regex: `^[a-f][0-9a-f]*$`,
			expectedAST: &AST{
				Root: &Concat{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Char{Lo: 'a', Hi: 'f', Pos: 1},
								&Star{
									Expr: &Alt{
										Exprs: []Node{
											&Char{Lo: '0', Hi: '9', Pos: 2},
											&Char{Lo: 'a', Hi: 'f', Pos: 3},
										},
										comp: &computed{
											nullable: false,
											firstPos: Poses{2, 3},
											lastPos:  Poses{2, 3},
										},
									},
								},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1},
								lastPos:  Poses{1, 2, 3},
							},
						},
						&Char{Lo: endMarker, Hi: endMarker, Pos: 4},
					},
				},
				lastPos: 4,
				posToChar: map[Pos]char.Range{
					1: {'a', 'f'},
					2: {'0', '9'},
					3: {'a', 'f'},
					4: {endMarker, endMarker},
				},
				charToPos: map[char.Range]Poses{
					{'0', '9'}:             {2},
					{'a', 'f'}:             {1, 3},
					{endMarker, endMarker}: {4},
				},
				follows: map[Pos]Poses{
					1: {2, 3, 4},
					2: {2, 3, 4},
					3: {2, 3, 4},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1},
			expectedLastPos:  Poses{4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a, err := Parse(tc.regex)

			if tc.expectedError != "" {
				assert.Nil(t, a)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, a)

				assert.True(t, a.Root.equal(tc.expectedAST.Root))
				assert.Equal(t, tc.expectedNullable, a.Root.nullable())
				assert.Equal(t, tc.expectedFirstPos, a.Root.firstPos())
				assert.Equal(t, tc.expectedLastPos, a.Root.lastPos())
				assert.Equal(t, tc.expectedAST.lastPos, a.lastPos)
				assert.Equal(t, tc.expectedAST.posToChar, a.posToChar)
				assert.Equal(t, tc.expectedAST.charToPos, a.charToPos)
				assert.Equal(t, tc.expectedAST.follows, a.follows)
			}
		})
	}
}

func TestAST_ToDFA(t *testing.T) {
	tests := []struct {
		name        string
		a           *AST
		expectedDFA *automata.DFA
	}{
		{
			name: "OK",
			a: &AST{ // (a|b)*abb
				Root: &Concat{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Star{
									Expr: &Alt{
										Exprs: []Node{
											&Concat{
												Exprs: []Node{
													&Char{Lo: 'a', Hi: 'a', Pos: 1},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{1},
													lastPos:  Poses{1},
												},
											},
											&Concat{
												Exprs: []Node{
													&Char{Lo: 'b', Hi: 'b', Pos: 2},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{2},
													lastPos:  Poses{2},
												},
											},
										},
										comp: &computed{
											nullable: false,
											firstPos: Poses{1, 2},
											lastPos:  Poses{1, 2},
										},
									},
								},
								&Char{Lo: 'a', Hi: 'a', Pos: 3},
								&Char{Lo: 'b', Hi: 'b', Pos: 4},
								&Char{Lo: 'b', Hi: 'b', Pos: 5},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1, 2, 3},
								lastPos:  Poses{5},
							},
						},
						&Char{Lo: endMarker, Hi: endMarker, Pos: 6},
					},
				},
				lastPos: 6,
				posToChar: map[Pos]char.Range{
					1: {'a', 'a'},
					2: {'b', 'b'},
					3: {'a', 'a'},
					4: {'b', 'b'},
					5: {'b', 'b'},
					6: {endMarker, endMarker},
				},
				charToPos: map[char.Range]Poses{
					{'a', 'a'}:             {1, 3},
					{'b', 'b'}:             {2, 4, 5},
					{endMarker, endMarker}: {6},
				},
				follows: map[Pos]Poses{
					1: {1, 2, 3},
					2: {1, 2, 3},
					3: {4},
					4: {5},
					5: {6},
				},
			},
			expectedDFA: automata.NewDFABuilder().
				SetStart(0).
				SetFinal([]automata.State{3}).
				AddTransition(0, 'a', 'a', 1).
				AddTransition(0, 'b', 'b', 0).
				AddTransition(1, 'a', 'a', 1).
				AddTransition(1, 'b', 'b', 2).
				AddTransition(2, 'a', 'a', 1).
				AddTransition(2, 'b', 'b', 3).
				AddTransition(3, 'a', 'a', 1).
				AddTransition(3, 'b', 'b', 0).
				Build(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := tc.a.ToDFA()
			assert.True(t, dfa.Equal(tc.expectedDFA))
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name             string
		node             *Concat
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "Nullable",
			node: &Concat{
				Exprs: []Node{
					&Star{
						Expr: &Char{
							Lo:  'a',
							Hi:  'a',
							Pos: 1,
						},
					},
					&Star{
						Expr: &Char{
							Lo:  'b',
							Hi:  'b',
							Pos: 2,
						},
					},
				},
			},
			expectedNullable: true,
			expectedFirstPos: Poses{1, 2},
			expectedLastPos:  Poses{1, 2},
		},
		{
			name: "Flat",
			node: &Concat{
				Exprs: []Node{
					&Char{
						Lo:  '_',
						Hi:  '_',
						Pos: 1,
					},
					&Char{
						Lo:  '0',
						Hi:  '9',
						Pos: 2,
					},
					&Char{
						Lo:  'A',
						Hi:  'Z',
						Pos: 3,
					},
					&Char{
						Lo:  'a',
						Hi:  'z',
						Pos: 4,
					},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1},
			expectedLastPos:  Poses{4},
		},
		{
			name: "Nested",
			node: &Concat{
				Exprs: []Node{
					&Alt{
						Exprs: []Node{
							&Char{
								Lo:  '_',
								Hi:  '_',
								Pos: 1,
							},
							&Char{
								Lo:  '0',
								Hi:  '9',
								Pos: 2,
							},
						},
					},
					&Alt{
						Exprs: []Node{
							&Char{
								Lo:  'A',
								Hi:  'Z',
								Pos: 3,
							},
							&Char{
								Lo:  'a',
								Hi:  'z',
								Pos: 4,
							},
						},
					},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2},
			expectedLastPos:  Poses{3, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.firstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.lastPos())
		})
	}
}

func TestAlt(t *testing.T) {
	tests := []struct {
		name             string
		node             *Alt
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "Nullable",
			node: &Alt{
				Exprs: []Node{
					&Star{
						Expr: &Char{
							Lo:  'a',
							Hi:  'a',
							Pos: 1,
						},
					},
					&Star{
						Expr: &Char{
							Lo:  'b',
							Hi:  'b',
							Pos: 2,
						},
					},
				},
			},
			expectedNullable: true,
			expectedFirstPos: Poses{1, 2},
			expectedLastPos:  Poses{1, 2},
		},
		{
			name: "Flat",
			node: &Alt{
				Exprs: []Node{
					&Char{
						Lo:  '_',
						Hi:  '_',
						Pos: 1,
					},
					&Char{
						Lo:  '0',
						Hi:  '9',
						Pos: 2,
					},
					&Char{
						Lo:  'A',
						Hi:  'Z',
						Pos: 3,
					},
					&Char{
						Lo:  'a',
						Hi:  'z',
						Pos: 4,
					},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2, 3, 4},
			expectedLastPos:  Poses{1, 2, 3, 4},
		},
		{
			name: "Nested",
			node: &Alt{
				Exprs: []Node{
					&Concat{
						Exprs: []Node{
							&Char{
								Lo:  '_',
								Hi:  '_',
								Pos: 1,
							},
							&Char{
								Lo:  '0',
								Hi:  '9',
								Pos: 2,
							},
						},
					},
					&Concat{
						Exprs: []Node{
							&Char{
								Lo:  'A',
								Hi:  'Z',
								Pos: 3,
							},
							&Char{
								Lo:  'a',
								Hi:  'z',
								Pos: 4,
							},
						},
					},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 3},
			expectedLastPos:  Poses{2, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.firstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.lastPos())
		})
	}
}

func TestStar(t *testing.T) {
	tests := []struct {
		name             string
		node             *Star
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "OK",
			node: &Star{
				Expr: &Char{
					Lo:  'a',
					Hi:  'a',
					Pos: 1,
				},
			},
			expectedNullable: true,
			expectedFirstPos: Poses{1},
			expectedLastPos:  Poses{1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.firstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.lastPos())
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name             string
		node             *Empty
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name:             "OK",
			node:             &Empty{},
			expectedNullable: true,
			expectedFirstPos: Poses{},
			expectedLastPos:  Poses{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.firstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.lastPos())
		})
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		name             string
		node             *Char
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name:             "Single",
			node:             &Char{Lo: 'a', Hi: 'a', Pos: 1},
			expectedNullable: false,
			expectedFirstPos: Poses{1},
			expectedLastPos:  Poses{1},
		},
		{
			name:             "Range",
			node:             &Char{Lo: 'a', Hi: 'z', Pos: 2},
			expectedNullable: false,
			expectedFirstPos: Poses{2},
			expectedLastPos:  Poses{2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.firstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.lastPos())
		})
	}
}

func TestPoses(t *testing.T) {
	type LessTest struct {
		i, j     int
		expected bool
	}

	type SwapTest struct {
		i, j     int
		expected Poses
	}

	type ContainsTest struct {
		q        Pos
		expected bool
	}

	type EqualTests struct {
		q        Poses
		expected bool
	}

	type UnionTest struct {
		q        Poses
		expected Poses
	}

	tests := []struct {
		name          string
		p             Poses
		expectedLen   int
		LessTests     []LessTest
		SwapTests     []SwapTest
		ContainsTests []ContainsTest
		EqualTests    []EqualTests
		UnionTests    []UnionTest
	}{
		{
			name:        "OK",
			p:           Poses{1, 2, 3, 5, 8},
			expectedLen: 5,
			LessTests: []LessTest{
				{0, 0, false},
				{1, 2, true},
				{4, 3, false},
			},
			SwapTests: []SwapTest{
				{0, 1, Poses{2, 1, 3, 5, 8}},
			},
			ContainsTests: []ContainsTest{
				{8, true},
				{13, false},
			},
			EqualTests: []EqualTests{
				{Poses{1, 2, 3, 5}, false},
				{Poses{1, 2, 3, 5, 8}, true},
				{Poses{2, 1, 3, 5, 8}, true},
				{Poses{1, 2, 3, 5, 8, 13}, false},
			},
			UnionTests: []UnionTest{
				{
					Poses{5, 8, 13, 21, 34, 55},
					Poses{1, 2, 3, 5, 8, 13, 21, 34, 55},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.p

			t.Run("Len", func(t *testing.T) {
				assert.Equal(t, tc.expectedLen, p.Len())
			})

			t.Run("Less", func(t *testing.T) {
				for _, tc := range tc.LessTests {
					assert.Equal(t, tc.expected, p.Less(tc.i, tc.j))
				}
			})

			t.Run("Swap", func(t *testing.T) {
				for _, tc := range tc.SwapTests {
					p.Swap(tc.i, tc.j)
					assert.Equal(t, tc.expected, p)
				}
			})

			t.Run("Contains", func(t *testing.T) {
				for _, tc := range tc.ContainsTests {
					assert.Equal(t, tc.expected, p.Contains(tc.q))
				}
			})

			t.Run("Equal", func(t *testing.T) {
				for _, tc := range tc.EqualTests {
					assert.Equal(t, tc.expected, p.Equal(tc.q))
				}
			})

			t.Run("Union", func(t *testing.T) {
				for _, tc := range tc.UnionTests {
					u := p.Union(tc.q)
					assert.True(t, u.Equal(tc.expected))
				}
			})
		})
	}
}

//==================================================< MAPPERS >==================================================

var testNodes = map[string]Node{
	/* CHAR CLASSES */

	"ws": &Alt{
		Exprs: []Node{
			&Char{Lo: ' ', Hi: ' '},
			&Char{Lo: '\t', Hi: '\r'},
		},
	},

	"not_ws": &Alt{
		Exprs: []Node{
			&Char{Lo: 0x00, Hi: 0x08},
			&Char{Lo: 0x0E, Hi: 0x1F},
			&Char{Lo: 0x21, Hi: 0x10FFFF},
		},
	},

	"digit": &Char{Lo: '0', Hi: '9'},

	"not_digit": &Alt{
		Exprs: []Node{
			&Char{Lo: 0x00, Hi: 0x2F},
			&Char{Lo: 0x3A, Hi: 0x10FFFF},
		},
	},

	"word": &Alt{
		Exprs: []Node{
			&Char{Lo: '0', Hi: '9'},
			&Char{Lo: 'A', Hi: 'Z'},
			&Char{Lo: '_', Hi: '_'},
			&Char{Lo: 'a', Hi: 'z'},
		},
	},

	"not_word": &Alt{
		Exprs: []Node{
			&Char{Lo: 0x00, Hi: 0x2F},
			&Char{Lo: 0x3A, Hi: 0x40},
			&Char{Lo: 0x5B, Hi: 0x5E},
			&Char{Lo: 0x60, Hi: 0x60},
			&Char{Lo: 0x7B, Hi: 0x10FFFF},
		},
	},

	/* ASCII CLASSES */

	"blank": &Alt{
		Exprs: []Node{
			&Char{Lo: ' ', Hi: ' '},
			&Char{Lo: '\t', Hi: '\t'},
		},
	},

	"space": &Alt{
		Exprs: []Node{
			&Char{Lo: ' ', Hi: ' '},
			&Char{Lo: '\t', Hi: '\r'},
		},
	},

	// "digit" is already added.

	"xdigit": &Alt{
		Exprs: []Node{
			&Char{Lo: '0', Hi: '9'},
			&Char{Lo: 'A', Hi: 'F'},
			&Char{Lo: 'a', Hi: 'f'},
		},
	},

	"upper": &Char{Lo: 'A', Hi: 'Z'},

	"lower": &Char{Lo: 'a', Hi: 'z'},

	"alpha": &Alt{
		Exprs: []Node{
			&Char{Lo: 'A', Hi: 'Z'},
			&Char{Lo: 'a', Hi: 'z'},
		},
	},

	"alnum": &Alt{
		Exprs: []Node{
			&Char{Lo: '0', Hi: '9'},
			&Char{Lo: 'A', Hi: 'Z'},
			&Char{Lo: 'a', Hi: 'z'},
		},
	},

	// "word" is already added.

	"ascii": &Char{Lo: 0x00, Hi: 0x7F},

	/* UNICODE CLASSES */

	"Number": &Alt{
		Exprs: []Node{
			&Char{Lo: 0x30, Hi: 0x39},
			&Char{Lo: 0xB2, Hi: 0xB3},
			&Char{Lo: 0xB9, Hi: 0xB9},
			&Char{Lo: 0xBC, Hi: 0xBC},
			&Char{Lo: 0xBD, Hi: 0xBE},
			&Char{Lo: 0x0660, Hi: 0x0669},
			&Char{Lo: 0x06F0, Hi: 0x06F9},
			&Char{Lo: 0x07C0, Hi: 0x07C9},
			&Char{Lo: 0x0966, Hi: 0x096F},
			&Char{Lo: 0x09E6, Hi: 0x09EF},
			&Char{Lo: 0x09F4, Hi: 0x09F9},
			&Char{Lo: 0x0A66, Hi: 0x0A6F},
			&Char{Lo: 0x0AE6, Hi: 0x0AEF},
			&Char{Lo: 0x0B66, Hi: 0x0B6F},
			&Char{Lo: 0x0B72, Hi: 0x0B77},
			&Char{Lo: 0x0BE6, Hi: 0x0BF2},
			&Char{Lo: 0x0C66, Hi: 0x0C6F},
			&Char{Lo: 0x0C78, Hi: 0x0C7E},
			&Char{Lo: 0x0CE6, Hi: 0x0CEF},
			&Char{Lo: 0x0D58, Hi: 0x0D5E},
			&Char{Lo: 0x0D66, Hi: 0x0D78},
			&Char{Lo: 0x0DE6, Hi: 0x0DEF},
			&Char{Lo: 0x0E50, Hi: 0x0E59},
			&Char{Lo: 0x0ED0, Hi: 0x0ED9},
			&Char{Lo: 0x0F20, Hi: 0x0F33},
			&Char{Lo: 0x1040, Hi: 0x1049},
			&Char{Lo: 0x1090, Hi: 0x1099},
			&Char{Lo: 0x1369, Hi: 0x137C},
			&Char{Lo: 0x16EE, Hi: 0x16F0},
			&Char{Lo: 0x17E0, Hi: 0x17E9},
			&Char{Lo: 0x17F0, Hi: 0x17F9},
			&Char{Lo: 0x1810, Hi: 0x1819},
			&Char{Lo: 0x1946, Hi: 0x194F},
			&Char{Lo: 0x19D0, Hi: 0x19DA},
			&Char{Lo: 0x1A80, Hi: 0x1A89},
			&Char{Lo: 0x1A90, Hi: 0x1A99},
			&Char{Lo: 0x1B50, Hi: 0x1B59},
			&Char{Lo: 0x1BB0, Hi: 0x1BB9},
			&Char{Lo: 0x1C40, Hi: 0x1C49},
			&Char{Lo: 0x1C50, Hi: 0x1C59},
			&Char{Lo: 0x2070, Hi: 0x2070},
			&Char{Lo: 0x2074, Hi: 0x2074},
			&Char{Lo: 0x2075, Hi: 0x2079},
			&Char{Lo: 0x2080, Hi: 0x2089},
			&Char{Lo: 0x2150, Hi: 0x2182},
			&Char{Lo: 0x2185, Hi: 0x2189},
			&Char{Lo: 0x2460, Hi: 0x249B},
			&Char{Lo: 0x24EA, Hi: 0x24FF},
			&Char{Lo: 0x2776, Hi: 0x2793},
			&Char{Lo: 0x2CFD, Hi: 0x2CFD},
			&Char{Lo: 0x3007, Hi: 0x3007},
			&Char{Lo: 0x3021, Hi: 0x3029},
			&Char{Lo: 0x3038, Hi: 0x303A},
			&Char{Lo: 0x3192, Hi: 0x3195},
			&Char{Lo: 0x3220, Hi: 0x3229},
			&Char{Lo: 0x3248, Hi: 0x324F},
			&Char{Lo: 0x3251, Hi: 0x325F},
			&Char{Lo: 0x3280, Hi: 0x3289},
			&Char{Lo: 0x32B1, Hi: 0x32BF},
			&Char{Lo: 0xA620, Hi: 0xA629},
			&Char{Lo: 0xA6E6, Hi: 0xA6EF},
			&Char{Lo: 0xA830, Hi: 0xA835},
			&Char{Lo: 0xA8D0, Hi: 0xA8D9},
			&Char{Lo: 0xA900, Hi: 0xA909},
			&Char{Lo: 0xA9D0, Hi: 0xA9D9},
			&Char{Lo: 0xA9F0, Hi: 0xA9F9},
			&Char{Lo: 0xAA50, Hi: 0xAA59},
			&Char{Lo: 0xABF0, Hi: 0xABF9},
			&Char{Lo: 0xFF10, Hi: 0xFF19},
			&Char{Lo: 0x010107, Hi: 0x010133},
			&Char{Lo: 0x010140, Hi: 0x010178},
			&Char{Lo: 0x01018A, Hi: 0x01018B},
			&Char{Lo: 0x0102E1, Hi: 0x0102FB},
			&Char{Lo: 0x010320, Hi: 0x010323},
			&Char{Lo: 0x010341, Hi: 0x010341},
			&Char{Lo: 0x01034A, Hi: 0x01034A},
			&Char{Lo: 0x0103D1, Hi: 0x0103D5},
			&Char{Lo: 0x0104A0, Hi: 0x0104A9},
			&Char{Lo: 0x010858, Hi: 0x01085F},
			&Char{Lo: 0x010879, Hi: 0x01087F},
			&Char{Lo: 0x0108A7, Hi: 0x0108AF},
			&Char{Lo: 0x0108FB, Hi: 0x0108FF},
			&Char{Lo: 0x010916, Hi: 0x01091B},
			&Char{Lo: 0x0109BC, Hi: 0x0109BD},
			&Char{Lo: 0x0109C0, Hi: 0x0109CF},
			&Char{Lo: 0x0109D2, Hi: 0x0109FF},
			&Char{Lo: 0x010A40, Hi: 0x010A48},
			&Char{Lo: 0x010A7D, Hi: 0x010A7E},
			&Char{Lo: 0x010A9D, Hi: 0x010A9F},
			&Char{Lo: 0x010AEB, Hi: 0x010AEF},
			&Char{Lo: 0x010B58, Hi: 0x010B5F},
			&Char{Lo: 0x010B78, Hi: 0x010B7F},
			&Char{Lo: 0x010BA9, Hi: 0x010BAF},
			&Char{Lo: 0x010CFA, Hi: 0x010CFF},
			&Char{Lo: 0x010D30, Hi: 0x010D39},
			&Char{Lo: 0x010E60, Hi: 0x010E7E},
			&Char{Lo: 0x010F1D, Hi: 0x010F26},
			&Char{Lo: 0x010F51, Hi: 0x010F54},
			&Char{Lo: 0x010FC5, Hi: 0x010FCB},
			&Char{Lo: 0x011052, Hi: 0x01106F},
			&Char{Lo: 0x0110F0, Hi: 0x0110F9},
			&Char{Lo: 0x011136, Hi: 0x01113F},
			&Char{Lo: 0x0111D0, Hi: 0x0111D9},
			&Char{Lo: 0x0111E1, Hi: 0x0111F4},
			&Char{Lo: 0x0112F0, Hi: 0x0112F9},
			&Char{Lo: 0x011450, Hi: 0x011459},
			&Char{Lo: 0x0114D0, Hi: 0x0114D9},
			&Char{Lo: 0x011650, Hi: 0x011659},
			&Char{Lo: 0x0116C0, Hi: 0x0116C9},
			&Char{Lo: 0x011730, Hi: 0x01173B},
			&Char{Lo: 0x0118E0, Hi: 0x0118F2},
			&Char{Lo: 0x011950, Hi: 0x011959},
			&Char{Lo: 0x011C50, Hi: 0x011C6C},
			&Char{Lo: 0x011D50, Hi: 0x011D59},
			&Char{Lo: 0x011DA0, Hi: 0x011DA9},
			&Char{Lo: 0x011F50, Hi: 0x011F59},
			&Char{Lo: 0x011FC0, Hi: 0x011FD4},
			&Char{Lo: 0x012400, Hi: 0x01246E},
			&Char{Lo: 0x016A60, Hi: 0x016A69},
			&Char{Lo: 0x016AC0, Hi: 0x016AC9},
			&Char{Lo: 0x016B50, Hi: 0x016B59},
			&Char{Lo: 0x016B5B, Hi: 0x016B61},
			&Char{Lo: 0x016E80, Hi: 0x016E96},
			&Char{Lo: 0x01D2C0, Hi: 0x01D2D3},
			&Char{Lo: 0x01D2E0, Hi: 0x01D2F3},
			&Char{Lo: 0x01D360, Hi: 0x01D378},
			&Char{Lo: 0x01D7CE, Hi: 0x01D7FF},
			&Char{Lo: 0x01E140, Hi: 0x01E149},
			&Char{Lo: 0x01E2F0, Hi: 0x01E2F9},
			&Char{Lo: 0x01E4F0, Hi: 0x01E4F9},
			&Char{Lo: 0x01E8C7, Hi: 0x01E8CF},
			&Char{Lo: 0x01E950, Hi: 0x01E959},
			&Char{Lo: 0x01EC71, Hi: 0x01ECAB},
			&Char{Lo: 0x01ECAD, Hi: 0x01ECAF},
			&Char{Lo: 0x01ECB1, Hi: 0x01ECB4},
			&Char{Lo: 0x01ED01, Hi: 0x01ED2D},
			&Char{Lo: 0x01ED2F, Hi: 0x01ED3D},
			&Char{Lo: 0x01F100, Hi: 0x01F10C},
			&Char{Lo: 0x01FBF0, Hi: 0x01FBF9},
		},
	},

	"Unicode": &Char{Lo: 0x00, Hi: 0x10FFFF},
}

var testRanges = map[string]char.RangeList{
	/* CHAR CLASSES */

	"ws":        {{' ', ' '}, {'\t', '\r'}},
	"not_ws":    {{0x00, 0x08}, {0x0E, 0x1F}, {0x21, 0x10FFFF}},
	"digit":     {{'0', '9'}},
	"not_digit": {{0x00, 0x2F}, {0x3A, 0x10FFFF}},
	"word":      {{'0', '9'}, {'A', 'Z'}, {'_', '_'}, {'a', 'z'}},
	"not_word":  {{0x00, 0x2F}, {0x3A, 0x40}, {0x5B, 0x5E}, {0x60, 0x60}, {0x7B, 0x10FFFF}},

	/* ASCII CLASSES */

	"blank": {{' ', ' '}, {'\t', '\t'}},
	"space": {{' ', ' '}, {'\t', '\r'}},
	// "digit" is already added.
	"xdigit": {{'0', '9'}, {'A', 'F'}, {'a', 'f'}},
	"upper":  {{'A', 'Z'}},
	"lower":  {{'a', 'z'}},
	"alpha":  {{'A', 'Z'}, {'a', 'z'}},
	"alnum":  {{'0', '9'}, {'A', 'Z'}, {'a', 'z'}},
	// "word" is already added.
	"ascii": {{0x00, 0x7F}},

	/* UNICODE CLASSES */

	"Number": {
		{0x30, 0x39},
		{0xB2, 0xB3},
		{0xB9, 0xB9},
		{0xBC, 0xBC},
		{0xBD, 0xBE},
		{0x0660, 0x0669},
		{0x06F0, 0x06F9},
		{0x07C0, 0x07C9},
		{0x0966, 0x096F},
		{0x09E6, 0x09EF},
		{0x09F4, 0x09F9},
		{0x0A66, 0x0A6F},
		{0x0AE6, 0x0AEF},
		{0x0B66, 0x0B6F},
		{0x0B72, 0x0B77},
		{0x0BE6, 0x0BF2},
		{0x0C66, 0x0C6F},
		{0x0C78, 0x0C7E},
		{0x0CE6, 0x0CEF},
		{0x0D58, 0x0D5E},
		{0x0D66, 0x0D78},
		{0x0DE6, 0x0DEF},
		{0x0E50, 0x0E59},
		{0x0ED0, 0x0ED9},
		{0x0F20, 0x0F33},
		{0x1040, 0x1049},
		{0x1090, 0x1099},
		{0x1369, 0x137C},
		{0x16EE, 0x16F0},
		{0x17E0, 0x17E9},
		{0x17F0, 0x17F9},
		{0x1810, 0x1819},
		{0x1946, 0x194F},
		{0x19D0, 0x19DA},
		{0x1A80, 0x1A89},
		{0x1A90, 0x1A99},
		{0x1B50, 0x1B59},
		{0x1BB0, 0x1BB9},
		{0x1C40, 0x1C49},
		{0x1C50, 0x1C59},
		{0x2070, 0x2070},
		{0x2074, 0x2074},
		{0x2075, 0x2079},
		{0x2080, 0x2089},
		{0x2150, 0x2182},
		{0x2185, 0x2189},
		{0x2460, 0x249B},
		{0x24EA, 0x24FF},
		{0x2776, 0x2793},
		{0x2CFD, 0x2CFD},
		{0x3007, 0x3007},
		{0x3021, 0x3029},
		{0x3038, 0x303A},
		{0x3192, 0x3195},
		{0x3220, 0x3229},
		{0x3248, 0x324F},
		{0x3251, 0x325F},
		{0x3280, 0x3289},
		{0x32B1, 0x32BF},
		{0xA620, 0xA629},
		{0xA6E6, 0xA6EF},
		{0xA830, 0xA835},
		{0xA8D0, 0xA8D9},
		{0xA900, 0xA909},
		{0xA9D0, 0xA9D9},
		{0xA9F0, 0xA9F9},
		{0xAA50, 0xAA59},
		{0xABF0, 0xABF9},
		{0xFF10, 0xFF19},
		{0x010107, 0x010133},
		{0x010140, 0x010178},
		{0x01018A, 0x01018B},
		{0x0102E1, 0x0102FB},
		{0x010320, 0x010323},
		{0x010341, 0x010341},
		{0x01034A, 0x01034A},
		{0x0103D1, 0x0103D5},
		{0x0104A0, 0x0104A9},
		{0x010858, 0x01085F},
		{0x010879, 0x01087F},
		{0x0108A7, 0x0108AF},
		{0x0108FB, 0x0108FF},
		{0x010916, 0x01091B},
		{0x0109BC, 0x0109BD},
		{0x0109C0, 0x0109CF},
		{0x0109D2, 0x0109FF},
		{0x010A40, 0x010A48},
		{0x010A7D, 0x010A7E},
		{0x010A9D, 0x010A9F},
		{0x010AEB, 0x010AEF},
		{0x010B58, 0x010B5F},
		{0x010B78, 0x010B7F},
		{0x010BA9, 0x010BAF},
		{0x010CFA, 0x010CFF},
		{0x010D30, 0x010D39},
		{0x010E60, 0x010E7E},
		{0x010F1D, 0x010F26},
		{0x010F51, 0x010F54},
		{0x010FC5, 0x010FCB},
		{0x011052, 0x01106F},
		{0x0110F0, 0x0110F9},
		{0x011136, 0x01113F},
		{0x0111D0, 0x0111D9},
		{0x0111E1, 0x0111F4},
		{0x0112F0, 0x0112F9},
		{0x011450, 0x011459},
		{0x0114D0, 0x0114D9},
		{0x011650, 0x011659},
		{0x0116C0, 0x0116C9},
		{0x011730, 0x01173B},
		{0x0118E0, 0x0118F2},
		{0x011950, 0x011959},
		{0x011C50, 0x011C6C},
		{0x011D50, 0x011D59},
		{0x011DA0, 0x011DA9},
		{0x011F50, 0x011F59},
		{0x011FC0, 0x011FD4},
		{0x012400, 0x01246E},
		{0x016A60, 0x016A69},
		{0x016AC0, 0x016AC9},
		{0x016B50, 0x016B59},
		{0x016B5B, 0x016B61},
		{0x016E80, 0x016E96},
		{0x01D2C0, 0x01D2D3},
		{0x01D2E0, 0x01D2F3},
		{0x01D360, 0x01D378},
		{0x01D7CE, 0x01D7FF},
		{0x01E140, 0x01E149},
		{0x01E2F0, 0x01E2F9},
		{0x01E4F0, 0x01E4F9},
		{0x01E8C7, 0x01E8CF},
		{0x01E950, 0x01E959},
		{0x01EC71, 0x01ECAB},
		{0x01ECAD, 0x01ECAF},
		{0x01ECB1, 0x01ECB4},
		{0x01ED01, 0x01ED2D},
		{0x01ED2F, 0x01ED3D},
		{0x01F100, 0x01F10C},
		{0x01FBF0, 0x01FBF9},
	},

	"Unicode": {
		{0x00, 0x10FFFF},
	},
}

type MapperTest struct {
	name           string
	r              combinator.Result
	expectedResult combinator.Result
	expectedOK     bool
	expectedError  string
}

func intPtr(v int) *int {
	return &v
}

func TestMappers_ToAnyChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '.',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["Unicode"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnyChar(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSingleChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: 'x',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Char{Lo: 'x', Hi: 'x'},
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: char.RangeList{
						{'x', 'x'},
					},
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSingleChar(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Whitespace",
			r: combinator.Result{
				Val: `\s`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["ws"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["ws"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWhitespace",
			r: combinator.Result{
				Val: `\S`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["not_ws"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_ws"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Digit",
			r: combinator.Result{
				Val: `\d`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotDigit",
			r: combinator.Result{
				Val: `\D`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["not_digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: combinator.Result{
				Val: `\w`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["word"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWord",
			r: combinator.Result{
				Val: `\W`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["not_word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_word"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharClass(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToASCIICharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Blank",
			r: combinator.Result{
				Val: "[:blank:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["blank"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["blank"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Space",
			r: combinator.Result{
				Val: "[:space:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["space"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["space"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Digit",
			r: combinator.Result{
				Val: "[:digit:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_XDigit",
			r: combinator.Result{
				Val: "[:xdigit:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["xdigit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["xdigit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Upper",
			r: combinator.Result{
				Val: "[:upper:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["upper"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["upper"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Lower",
			r: combinator.Result{
				Val: "[:lower:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["lower"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["lower"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alpha",
			r: combinator.Result{
				Val: "[:alpha:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["alpha"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["alpha"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alnum",
			r: combinator.Result{
				Val: "[:alnum:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["alnum"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["alnum"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: combinator.Result{
				Val: "[:word:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["word"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_ASCII",
			r: combinator.Result{
				Val: "[:ascii:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["ascii"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["ascii"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToASCIICharClass(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUnicodeCategory(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: "Letter",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: "Letter",
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUnicodeCategory(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUnicodeCharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "InvalidClass",
			r: combinator.Result{
				Val: combinator.List{
					{Val: `\p`, Pos: 2},
					{Val: '{', Pos: 4},
					{Val: "Runic", Pos: 5},
					{Val: '}', Pos: 11},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{},
			expectedOK:     false,
		},
		{
			name: "Success_Number",
			r: combinator.Result{
				Val: combinator.List{
					{Val: `\p`, Pos: 2},
					{Val: '{', Pos: 4},
					{Val: "Number", Pos: 5},
					{Val: '}', Pos: 11},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNodes["Number"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["Number"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUnicodeCharClass(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepOp(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepOp(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUpperBound(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Unbounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: ',', Pos: 2},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: (*int)(nil),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: ',', Pos: 2},
					{Val: 4, Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: intPtr(4),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUpperBound(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Fixed",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{Val: combinator.Empty{}},
					{Val: '}', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Unbounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: (*int)(nil),
						Pos: 4,
					},
					{Val: '}', Pos: 5},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: (*int)(nil),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: intPtr(6),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(6),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 6, Pos: 3},
					{
						Val: intPtr(2),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 6,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK:    true,
			expectedError: "invalid repetition range {6,2}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepetition(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepetition(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToQuantifier(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_NonLazy",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '*', Pos: 2},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: false,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '*', Pos: 2},
					{Val: '?', Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: true,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToQuantifier(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharInRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharInRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: 'a', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'f', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Char{Lo: 'a', Hi: 'f'},
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: char.RangeList{
						{'a', 'f'},
					},
				},
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: 'f', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'a', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: nil,
				Pos: 2,
				Bag: nil,
			},
			expectedOK:    true,
			expectedError: "invalid character range f-a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharRange(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroupItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedResult: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroupItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '[', Pos: 2},
					{Val: combinator.Empty{}},
					{
						Val: combinator.List{
							{
								Val: testNodes["xdigit"],
								Pos: 3,
								Bag: combinator.Bag{
									bagKeyCharRanges: testRanges["xdigit"],
								},
							},
							{
								Val: &Char{Lo: '-', Hi: '-'},
								Pos: 12,
								Bag: combinator.Bag{
									bagKeyCharRanges: char.RangeList{
										{'-', '-'},
									},
								},
							},
						},
						Pos: 3,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Alt{
					Exprs: []Node{
						&Char{Lo: '-', Hi: '-'},
						&Char{Lo: '0', Hi: '9'},
						&Char{Lo: 'A', Hi: 'F'},
						&Char{Lo: 'a', Hi: 'f'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Negated",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '[', Pos: 2},
					{Val: '^', Pos: 3},
					{
						Val: combinator.List{
							{
								Val: testNodes["alnum"],
								Pos: 4,
								Bag: combinator.Bag{
									bagKeyCharRanges: testRanges["alnum"],
								},
							},
						},
						Pos: 4,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Alt{
					Exprs: []Node{
						&Char{Lo: 0x00, Hi: 0x2F},
						&Char{Lo: 0x3A, Hi: 0x40},
						&Char{Lo: 0x5B, Hi: 0x60},
						&Char{Lo: 0x7B, Hi: 0x10FFFF},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroup(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatchItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedResult: combinator.Result{
				Val: testNodes["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatchItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatch(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Char{Lo: 'x', Hi: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Alt{
					Exprs: []Node{
						&Empty{},
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Star{
					Expr: &Char{Lo: 'x', Hi: 'x'},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatch(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Char{Lo: 'x', Hi: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Alt{
					Exprs: []Node{
						&Empty{},
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Star{
					Expr: &Char{Lo: 'x', Hi: 'x'},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Char{Lo: 'x', Hi: 'x'},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
						&Alt{
							Exprs: []Node{
								&Empty{},
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
						&Star{
							Expr: &Char{Lo: 'x', Hi: 'x'},
						},
					},
				},
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToGroup(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToAnchor(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '$',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: EndOfString,
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnchor(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexprItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: &Char{Lo: 'x', Hi: 'x'},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Char{Lo: 'x', Hi: 'x'},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexprItem(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Char{Lo: 'x', Hi: 'x'},
						Pos: 2,
					},
					{
						Val: EndOfString,
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexpr(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToExpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
						Pos: 2,
					},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 'x', Hi: 'x'},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithExpr",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
						Pos: 2,
					},
					{
						Val: combinator.List{
							{Val: '|', Pos: 3},
							{
								Val: &Concat{
									Exprs: []Node{
										&Char{Lo: 'y', Hi: 'y'},
									},
								},
								Pos: 4,
							},
						},
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: &Alt{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Char{Lo: 'x', Hi: 'x'},
							},
						},
						&Concat{
							Exprs: []Node{
								&Char{Lo: 'y', Hi: 'y'},
							},
						},
					},
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToExpr(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRegex(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: combinator.Empty{}},
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Lo: 's', Hi: 's'},
							},
						},
						Pos: 0,
					},
				},
				Pos: 0,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 's', Hi: 's'},
					},
				},
				Pos: 0,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithStartOfString",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '^', Pos: 0},
					{
						Val: &Concat{
							Exprs: []Node{
								&Char{Lo: 's', Hi: 's'},
							},
						},
						Pos: 1,
					},
				},
				Pos: 0,
			},
			expectedResult: combinator.Result{
				Val: &Concat{
					Exprs: []Node{
						&Char{Lo: 's', Hi: 's'},
					},
				},
				Pos: 0,
				Bag: combinator.Bag{
					BagKeyStartOfString: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRegex(tc.r)

			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}
