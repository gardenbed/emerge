// Package provides a combinator parser for parsing regular expression into a non-deterministic finite automata.
//
// It implements the McNaughton-Yamada-Thompson algorithm to convert a regular expression to an NFA.
// The algorithm is syntax-directed, in the sense that it works recursively up the parse tree for the regular expression.
package nfa

import auto "github.com/moorara/algo/automata"

// Empty returns an NFA accepting the empty string ε.
func Empty() *auto.NFA {
	nfa := auto.NewNFA(0, auto.States{1})
	nfa.Add(0, auto.E, auto.States{1})

	return nfa
}

// Star returns an NFA accepting the Kleene closure of an NFA.
// N(r) accepts L(s*)
func Star(n *auto.NFA) *auto.NFA {
	s, f := auto.State(0), auto.State(1)
	nfa := auto.NewNFA(s, auto.States{f})

	_, ss, ff := nfa.Join(n)
	nfa.Add(s, auto.E, auto.States{f})
	nfa.Add(s, auto.E, auto.States{ss})
	for _, t := range ff {
		nfa.Add(t, auto.E, auto.States{ss})
		nfa.Add(t, auto.E, auto.States{f})
	}

	return nfa
}

// Alt returns an NFA accepting the alternation of a set of NFAs.
// N(r) accepts L(s) ∪ L(t) ∪ ...
func Alt(ns ...*auto.NFA) *auto.NFA {
	s, f := auto.State(0), auto.State(1)
	nfa := auto.NewNFA(s, auto.States{f})

	for _, n := range ns {
		_, ss, ff := nfa.Join(n)
		nfa.Add(s, auto.E, auto.States{ss})
		for _, t := range ff {
			nfa.Add(t, auto.E, auto.States{f})
		}
	}

	return nfa
}

// Concat returns an NFA accepting the concatenation of a set of NFAs.
// N(r) accepts L(s)L(t)...
func Concat(ns ...*auto.NFA) *auto.NFA {
	s, f := auto.State(0), auto.State(1)
	nfa := auto.NewNFA(s, auto.States{f})

	prev := auto.States{s}
	for _, n := range ns {
		_, ss, ff := nfa.Join(n)
		for _, t := range prev {
			nfa.Add(t, auto.E, auto.States{ss})
		}
		prev = ff
	}
	for _, t := range prev {
		nfa.Add(t, auto.E, auto.States{f})
	}

	return nfa
}
