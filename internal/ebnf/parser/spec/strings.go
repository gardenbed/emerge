package spec

import (
	"hash/fnv"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/sort"
)

var h = fnv.New64()

// Strings is a list of grammar strings, each representing a sequence of grammar symbols.
type Strings []grammar.String[grammar.Symbol]

// Contains returns true if the list of strings contains the given grammar string.
func (s Strings) Contains(α grammar.String[grammar.Symbol]) bool {
	for _, β := range s {
		if β.Equal(α) {
			return true
		}
	}

	return false
}

func eqStrings(lhs, rhs Strings) bool {
	for _, α := range lhs {
		if !rhs.Contains(α) {
			return false
		}
	}

	for _, α := range rhs {
		if !lhs.Contains(α) {
			return false
		}
	}

	return true
}

func hashStrings(s Strings) uint64 {
	h.Reset()

	sort.Quick(s, grammar.CmpString)
	for _, α := range s {
		// Hash.Write never returns an error
		_, _ = grammar.WriteString(h, α)
	}

	return h.Sum64()
}
