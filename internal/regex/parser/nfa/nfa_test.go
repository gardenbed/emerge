package nfa

import (
	"testing"

	"github.com/stretchr/testify/assert"

	auto "github.com/moorara/algo/automata"
)

func getTestNFAs() []*auto.NFA {
	n0 := auto.NewNFA(0, auto.States{1})
	n0.Add(0, auto.E, auto.States{1})

	// L(a)
	n1 := auto.NewNFA(0, auto.States{1})
	n1.Add(0, 'a', auto.States{1})

	// L(b)
	n2 := auto.NewNFA(0, auto.States{1})
	n2.Add(0, 'b', auto.States{1})

	// L(a*)
	n3 := auto.NewNFA(0, auto.States{1})
	n3.Add(0, auto.E, auto.States{1})
	n3.Add(0, auto.E, auto.States{2})
	n3.Add(2, 'a', auto.States{3})
	n3.Add(3, auto.E, auto.States{2})
	n3.Add(3, auto.E, auto.States{1})

	// L(b*)
	n4 := auto.NewNFA(0, auto.States{1})
	n4.Add(0, auto.E, auto.States{1})
	n4.Add(0, auto.E, auto.States{2})
	n4.Add(2, 'b', auto.States{3})
	n4.Add(3, auto.E, auto.States{2})
	n4.Add(3, auto.E, auto.States{1})

	// L(a|b)
	n5 := auto.NewNFA(0, auto.States{1})
	n5.Add(0, auto.E, auto.States{2})
	n5.Add(0, auto.E, auto.States{4})
	n5.Add(2, 'a', auto.States{3})
	n5.Add(4, 'b', auto.States{5})
	n5.Add(3, auto.E, auto.States{1})
	n5.Add(5, auto.E, auto.States{1})

	// L(ab)
	n6 := auto.NewNFA(0, auto.States{2})
	n6.Add(0, 'a', auto.States{1})
	n6.Add(1, 'b', auto.States{2})

	return []*auto.NFA{n0, n1, n2, n3, n4, n5, n6}
}

func TestEmpty(t *testing.T) {
	nfa := empty()
	expectedNFA := getTestNFAs()[0]
	assert.True(t, nfa.Equals(expectedNFA))
}

func TestConcat(t *testing.T) {
	nfas := getTestNFAs()

	tests := []struct {
		name        string
		ns          []*auto.NFA
		expectedNFA *auto.NFA
	}{
		{
			name:        "OK",
			ns:          nfas[1:3],
			expectedNFA: nfas[6],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nfa := concat(tc.ns...)
			assert.True(t, nfa.Equals(tc.expectedNFA))
		})
	}
}
