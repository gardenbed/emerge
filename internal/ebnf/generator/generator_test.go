package generator

import (
	"testing"

	auto "github.com/moorara/algo/automata"
	"github.com/stretchr/testify/assert"
)

func TestStringToDFA(t *testing.T) {
	dfa := auto.NewDFA(0, auto.States{7})
	dfa.Add(0, 'g', 1)
	dfa.Add(1, 'r', 2)
	dfa.Add(2, 'a', 3)
	dfa.Add(3, 'm', 4)
	dfa.Add(4, 'm', 5)
	dfa.Add(5, 'a', 6)
	dfa.Add(6, 'r', 7)

	tests := []struct {
		name        string
		s           string
		expectedDFA *auto.DFA
	}{
		{
			name:        "OK",
			s:           "grammar",
			expectedDFA: dfa,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := StringToDFA(tc.s)
			assert.True(t, dfa.Equals(tc.expectedDFA))
		})
	}
}
