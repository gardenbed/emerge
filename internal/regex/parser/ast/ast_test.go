package ast

import (
	"fmt"
	"testing"

	auto "github.com/moorara/algo/automata"
	"github.com/stretchr/testify/assert"
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
			expectedError: "invalid regular expression",
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
													&Char{Val: 'a', Pos: 1},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{1},
													lastPos:  Poses{1},
												},
											},
											&Concat{
												Exprs: []Node{
													&Char{Val: 'b', Pos: 2},
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
								&Char{Val: 'a', Pos: 3},
								&Char{Val: 'b', Pos: 4},
								&Char{Val: 'b', Pos: 5},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1, 2, 3},
								lastPos:  Poses{5},
							},
						},
						&Char{Val: endMarker, Pos: 6},
					},
				},
				lastPos: 6,
				posToChar: map[Pos]rune{
					1: 'a',
					2: 'b',
					3: 'a',
					4: 'b',
					5: 'b',
					6: endMarker,
				},
				charToPos: map[rune]Poses{
					'a':       Poses{1, 3},
					'b':       Poses{2, 4, 5},
					endMarker: Poses{6},
				},
				follows: map[Pos]Poses{
					1: Poses{1, 2, 3},
					2: Poses{1, 2, 3},
					3: Poses{4},
					4: Poses{5},
					5: Poses{6},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2, 3},
			expectedLastPos:  Poses{6},
		},
		{
			name:  "Success_Complex",
			regex: `^[a-f][0-9a-f]*$`,
			expectedAST: &AST{
				Root: &Concat{
					Exprs: []Node{
						&Concat{
							Exprs: []Node{
								&Alt{
									Exprs: []Node{
										&Char{Val: 'a', Pos: 1},
										&Char{Val: 'b', Pos: 2},
										&Char{Val: 'c', Pos: 3},
										&Char{Val: 'd', Pos: 4},
										&Char{Val: 'e', Pos: 5},
										&Char{Val: 'f', Pos: 6},
									},
									comp: &computed{
										nullable: false,
										firstPos: Poses{1, 2, 3, 4, 5, 6},
										lastPos:  Poses{1, 2, 3, 4, 5, 6},
									},
								},
								&Star{
									Expr: &Alt{
										Exprs: []Node{
											&Char{Val: '0', Pos: 7},
											&Char{Val: '1', Pos: 8},
											&Char{Val: '2', Pos: 9},
											&Char{Val: '3', Pos: 10},
											&Char{Val: '4', Pos: 11},
											&Char{Val: '5', Pos: 12},
											&Char{Val: '6', Pos: 13},
											&Char{Val: '7', Pos: 14},
											&Char{Val: '8', Pos: 15},
											&Char{Val: '9', Pos: 16},
											&Char{Val: 'a', Pos: 17},
											&Char{Val: 'b', Pos: 18},
											&Char{Val: 'c', Pos: 19},
											&Char{Val: 'd', Pos: 20},
											&Char{Val: 'e', Pos: 21},
											&Char{Val: 'f', Pos: 22},
										},
										comp: &computed{
											nullable: false,
											firstPos: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
											lastPos:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
										},
									},
								},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1, 2, 3, 4, 5, 6},
								lastPos:  Poses{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
							},
						},
						&Char{Val: endMarker, Pos: 23},
					},
				},
				lastPos: 23,
				posToChar: map[Pos]rune{
					1:  'a',
					2:  'b',
					3:  'c',
					4:  'd',
					5:  'e',
					6:  'f',
					7:  '0',
					8:  '1',
					9:  '2',
					10: '3',
					11: '4',
					12: '5',
					13: '6',
					14: '7',
					15: '8',
					16: '9',
					17: 'a',
					18: 'b',
					19: 'c',
					20: 'd',
					21: 'e',
					22: 'f',
					23: endMarker,
				},
				charToPos: map[rune]Poses{
					'0':       Poses{7},
					'1':       Poses{8},
					'2':       Poses{9},
					'3':       Poses{10},
					'4':       Poses{11},
					'5':       Poses{12},
					'6':       Poses{13},
					'7':       Poses{14},
					'8':       Poses{15},
					'9':       Poses{16},
					'a':       Poses{1, 17},
					'b':       Poses{2, 18},
					'c':       Poses{3, 19},
					'd':       Poses{4, 20},
					'e':       Poses{5, 21},
					'f':       Poses{6, 22},
					endMarker: Poses{23},
				},
				follows: map[Pos]Poses{
					1:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					2:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					3:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					4:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					5:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					6:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					7:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					8:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					9:  Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					10: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					11: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					12: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					13: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					14: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					15: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					16: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					17: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					18: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					19: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					20: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					21: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					22: Poses{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
				},
			},
			expectedNullable: false,
			expectedFirstPos: Poses{1, 2, 3, 4, 5, 6},
			expectedLastPos:  Poses{23},
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
				assert.Equal(t, tc.expectedAST, a)
				assert.Equal(t, tc.expectedNullable, a.Root.nullable())
				assert.Equal(t, tc.expectedFirstPos, a.Root.firstPos())
				assert.Equal(t, tc.expectedLastPos, a.Root.lastPos())
			}
		})
	}
}

func TestAST_ToDFA(t *testing.T) {
	dfa := auto.NewDFA(0, auto.States{3})
	dfa.Add(0, 'a', 1)
	dfa.Add(0, 'b', 0)
	dfa.Add(1, 'a', 1)
	dfa.Add(1, 'b', 2)
	dfa.Add(2, 'a', 1)
	dfa.Add(2, 'b', 3)
	dfa.Add(3, 'a', 1)
	dfa.Add(3, 'b', 0)

	tests := []struct {
		name        string
		a           *AST
		expectedDFA *auto.DFA
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
													&Char{Val: 'a', Pos: 1},
												},
												comp: &computed{
													nullable: false,
													firstPos: Poses{1},
													lastPos:  Poses{1},
												},
											},
											&Concat{
												Exprs: []Node{
													&Char{Val: 'b', Pos: 2},
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
								&Char{Val: 'a', Pos: 3},
								&Char{Val: 'b', Pos: 4},
								&Char{Val: 'b', Pos: 5},
							},
							comp: &computed{
								nullable: false,
								firstPos: Poses{1, 2, 3},
								lastPos:  Poses{5},
							},
						},
						&Char{Val: endMarker, Pos: 6},
					},
				},
				lastPos: 6,
				posToChar: map[Pos]rune{
					1: 'a',
					2: 'b',
					3: 'a',
					4: 'b',
					5: 'b',
					6: endMarker,
				},
				charToPos: map[rune]Poses{
					'a':       Poses{1, 3},
					'b':       Poses{2, 4, 5},
					endMarker: Poses{6},
				},
				follows: map[Pos]Poses{
					1: Poses{1, 2, 3},
					2: Poses{1, 2, 3},
					3: Poses{4},
					4: Poses{5},
					5: Poses{6},
				},
			},
			expectedDFA: dfa,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := tc.a.ToDFA()
			fmt.Println(dfa.Graphviz())
			assert.True(t, dfa.Equals(tc.expectedDFA))
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name             string
		node             Concat
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "Flat",
			node: Concat{
				Exprs: []Node{
					&Char{
						Val: 'a',
						Pos: 1,
					},
					&Char{
						Val: 'b',
						Pos: 2,
					},
					&Char{
						Val: 'c',
						Pos: 3,
					},
					&Char{
						Val: 'd',
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
			node: Concat{
				Exprs: []Node{
					&Alt{
						Exprs: []Node{
							&Char{
								Val: 'a',
								Pos: 1,
							},
							&Char{
								Val: 'b',
								Pos: 2,
							},
						},
					},
					&Alt{
						Exprs: []Node{
							&Char{
								Val: 'c',
								Pos: 3,
							},
							&Char{
								Val: 'd',
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
		node             Alt
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "Flat",
			node: Alt{
				Exprs: []Node{
					&Char{
						Val: 'a',
						Pos: 1,
					},
					&Char{
						Val: 'b',
						Pos: 2,
					},
					&Char{
						Val: 'c',
						Pos: 3,
					},
					&Char{
						Val: 'd',
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
			node: Alt{
				Exprs: []Node{
					&Concat{
						Exprs: []Node{
							&Char{
								Val: 'a',
								Pos: 1,
							},
							&Char{
								Val: 'b',
								Pos: 2,
							},
						},
					},
					&Concat{
						Exprs: []Node{
							&Char{
								Val: 'c',
								Pos: 3,
							},
							&Char{
								Val: 'd',
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
		node             Star
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "OK",
			node: Star{
				Expr: &Char{
					Val: 'a',
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
		node             Empty
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name:             "OK",
			node:             Empty{},
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
		node             Char
		expectedNullable bool
		expectedFirstPos Poses
		expectedLastPos  Poses
	}{
		{
			name: "OK",
			node: Char{
				Val: 'a',
				Pos: 1,
			},
			expectedNullable: false,
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

	type EqualsTests struct {
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
		EqualsTests   []EqualsTests
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
			EqualsTests: []EqualsTests{
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

			assert.Equal(t, tc.expectedLen, p.Len())

			for _, tc := range tc.LessTests {
				assert.Equal(t, tc.expected, p.Less(tc.i, tc.j))
			}

			for _, tc := range tc.SwapTests {
				p.Swap(tc.i, tc.j)
				assert.Equal(t, tc.expected, p)
			}

			for _, tc := range tc.ContainsTests {
				assert.Equal(t, tc.expected, p.Contains(tc.q))
			}

			for _, tc := range tc.EqualsTests {
				assert.Equal(t, tc.expected, p.Equals(tc.q))
			}

			for _, tc := range tc.UnionTests {
				u := p.Union(tc.q)
				assert.True(t, u.Equals(tc.expected))
			}
		})
	}
}
