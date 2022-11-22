// Package provides a combinator parser for parsing regular expression into a non-deterministic finite automata.
//
// It implements the McNaughton-Yamada-Thompson algorithm to convert a regular expression to an NFA.
// The algorithm is syntax-directed, in the sense that it works recursively up the parse tree for the regular expression.
package nfa

import auto "github.com/moorara/algo/automata"

// Empty returns an NFA accepting the empty string ε.
func empty() *auto.NFA {
	nfa := auto.NewNFA(0, auto.States{1})
	nfa.Add(0, auto.E, auto.States{1})

	return nfa
}

// Concat returns an NFA accepting the concatenation of a set of NFAs.
func concat(ns ...*auto.NFA) *auto.NFA {
	if len(ns) > 0 {
		return ns[0].Concat(ns[1:]...)
	}

	return empty()
}
