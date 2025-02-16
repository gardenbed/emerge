package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/grammar"
)

func TestStrings_Contains(t *testing.T) {
	tests := []struct {
		name             string
		s                Strings
		α                grammar.String[grammar.Symbol]
		expectedContains bool
	}{
		{
			name: "Contained",
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			α:                grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			expectedContains: true,
		},
		{
			name: "NotContained",
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			α:                grammar.String[grammar.Symbol]{grammar.Terminal("c"), grammar.NonTerminal("C")},
			expectedContains: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedContains, tc.s.Contains(tc.α))
		})
	}
}

func TestEqStrings(t *testing.T) {
	tests := []struct {
		name          string
		lhs           Strings
		rhs           Strings
		expectedEqual bool
	}{
		{
			name: "Equal",
			lhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			rhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
			},
			expectedEqual: true,
		},
		{
			name: "NotEqual",
			lhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			rhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
			},
			expectedEqual: false,
		},
		{
			name: "NotEqual",
			lhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			rhs: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
			},
			expectedEqual: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedEqual, eqStrings(tc.lhs, tc.rhs))
		})
	}
}

func TestHashStrings(t *testing.T) {
	tests := []struct {
		name         string
		s            Strings
		expectedHash uint64
	}{
		{
			name: "OK",
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("a"), grammar.NonTerminal("A")},
				grammar.String[grammar.Symbol]{grammar.Terminal("b"), grammar.NonTerminal("B")},
			},
			expectedHash: 0x27fdf99569a907d1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedHash, hashStrings(tc.s))
		})
	}
}
