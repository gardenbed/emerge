package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	tests := []struct {
		name             string
		node             Concat
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
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
			expectedFirstPos: []int{1},
			expectedLastPos:  []int{4},
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
			expectedFirstPos: []int{1, 2},
			expectedLastPos:  []int{3, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.Nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.FirstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.LastPos())
		})
	}
}

func TestAlt(t *testing.T) {
	tests := []struct {
		name             string
		node             Alt
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
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
			expectedFirstPos: []int{1, 2, 3, 4},
			expectedLastPos:  []int{1, 2, 3, 4},
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
			expectedFirstPos: []int{1, 3},
			expectedLastPos:  []int{2, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.Nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.FirstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.LastPos())
		})
	}
}

func TestStar(t *testing.T) {
	tests := []struct {
		name             string
		node             Star
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
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
			expectedFirstPos: []int{1},
			expectedLastPos:  []int{1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.Nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.FirstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.LastPos())
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name             string
		node             Empty
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
	}{
		{
			name:             "OK",
			node:             Empty{},
			expectedNullable: true,
			expectedFirstPos: []int{},
			expectedLastPos:  []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.Nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.FirstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.LastPos())
		})
	}
}

func TestChar(t *testing.T) {
	tests := []struct {
		name             string
		node             Char
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
	}{
		{
			name: "OK",
			node: Char{
				Val: 'a',
				Pos: 1,
			},
			expectedNullable: false,
			expectedFirstPos: []int{1},
			expectedLastPos:  []int{1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNullable, tc.node.Nullable())
			assert.Equal(t, tc.expectedFirstPos, tc.node.FirstPos())
			assert.Equal(t, tc.expectedLastPos, tc.node.LastPos())
		})
	}
}
