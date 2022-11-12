package generator

import (
	"testing"

	auto "github.com/moorara/algo/automata"
	"github.com/stretchr/testify/assert"
)

func getTestDFAs() []*auto.DFA {
	// QUO
	d0 := auto.NewDFA(0, auto.States{1})
	d0.Add(0, '"', 1)

	// SOL
	d1 := auto.NewDFA(0, auto.States{1})
	d1.Add(0, '\\', 1)

	// DEF
	d2 := auto.NewDFA(0, auto.States{1})
	d2.Add(0, '=', 1)

	// ALT
	d3 := auto.NewDFA(0, auto.States{1})
	d3.Add(0, '|', 1)

	// LPAREN
	d4 := auto.NewDFA(0, auto.States{1})
	d4.Add(0, '(', 1)

	// RPAREN
	d5 := auto.NewDFA(0, auto.States{1})
	d5.Add(0, ')', 1)

	// LBRACK
	d6 := auto.NewDFA(0, auto.States{1})
	d6.Add(0, '[', 1)

	// RBRACK
	d7 := auto.NewDFA(0, auto.States{1})
	d7.Add(0, ']', 1)

	// LBRACE
	d8 := auto.NewDFA(0, auto.States{1})
	d8.Add(0, '{', 1)

	// RBRACE
	d9 := auto.NewDFA(0, auto.States{1})
	d9.Add(0, '}', 1)

	// LLBRACE
	d10 := auto.NewDFA(0, auto.States{2})
	d10.Add(0, '{', 1)
	d10.Add(1, '{', 2)

	// RRBRACE
	d11 := auto.NewDFA(0, auto.States{2})
	d11.Add(0, '}', 1)
	d11.Add(1, '}', 2)

	// GRAMMER
	d12 := auto.NewDFA(0, auto.States{7})
	d12.Add(0, 'g', 1)
	d12.Add(1, 'r', 2)
	d12.Add(2, 'a', 3)
	d12.Add(3, 'm', 4)
	d12.Add(4, 'm', 5)
	d12.Add(5, 'a', 6)
	d12.Add(6, 'r', 7)

	d13 := auto.NewDFA(0, auto.States{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 19})
	d13.Add(0, '"', 1)
	d13.Add(0, '\\', 2)
	d13.Add(0, '=', 3)
	d13.Add(0, '|', 4)
	d13.Add(0, '(', 5)
	d13.Add(0, ')', 6)
	d13.Add(0, '[', 7)
	d13.Add(0, ']', 8)
	d13.Add(0, '{', 9)
	d13.Add(9, '{', 10)
	d13.Add(0, '}', 11)
	d13.Add(11, '}', 12)
	d13.Add(0, 'g', 13)
	d13.Add(13, 'r', 14)
	d13.Add(14, 'a', 15)
	d13.Add(15, 'm', 16)
	d13.Add(16, 'm', 17)
	d13.Add(17, 'a', 18)
	d13.Add(18, 'r', 19)

	return []*auto.DFA{d0, d1, d2, d3, d4, d5, d6, d7, d8, d9, d10, d11, d12, d13}
}

func TestStringToDFA(t *testing.T) {
	dfas := getTestDFAs()

	tests := []struct {
		name        string
		s           string
		expectedDFA *auto.DFA
	}{
		{
			name:        "QUO",
			s:           "\"",
			expectedDFA: dfas[0],
		},
		{
			name:        "SOL",
			s:           "\\",
			expectedDFA: dfas[1],
		},
		{
			name:        "DEF",
			s:           "=",
			expectedDFA: dfas[2],
		},
		{
			name:        "ALT",
			s:           "|",
			expectedDFA: dfas[3],
		},
		{
			name:        "LPAREN",
			s:           "(",
			expectedDFA: dfas[4],
		},
		{
			name:        "RPAREN",
			s:           ")",
			expectedDFA: dfas[5],
		},
		{
			name:        "LBRACK",
			s:           "[",
			expectedDFA: dfas[6],
		},
		{
			name:        "RBRACK",
			s:           "]",
			expectedDFA: dfas[7],
		},
		{
			name:        "LBRACE",
			s:           "{",
			expectedDFA: dfas[8],
		},
		{
			name:        "RBRACE",
			s:           "}",
			expectedDFA: dfas[9],
		},
		{
			name:        "LLBRACE",
			s:           "{{",
			expectedDFA: dfas[10],
		},
		{
			name:        "RRBRACE",
			s:           "}}",
			expectedDFA: dfas[11],
		},
		{
			name:        "GRAMMER",
			s:           "grammar",
			expectedDFA: dfas[12],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := StringToDFA(tc.s)
			assert.True(t, dfa.Equals(tc.expectedDFA))
		})
	}
}
