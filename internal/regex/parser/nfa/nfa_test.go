package nfa

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/automata"
)

func TestEmpty(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		nfa := empty()
		expectedNFA := automata.NewNFABuilder().
			SetStart(0).
			SetFinal([]automata.State{1}).
			AddTransition(0, automata.E, automata.E, []automata.State{1}).
			Build()

		assert.True(t, nfa.Equal(expectedNFA))
	})
}
