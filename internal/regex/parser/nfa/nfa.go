// Package nfa provides a combinator parser for parsing regular expression into a non-deterministic finite automata.
//
// It implements the McNaughton-Yamada-Thompson algorithm to convert a regular expression to an NFA.
// The algorithm is syntax-directed, in the sense that it works recursively up the parse tree for the regular expression.
package nfa

import "github.com/moorara/algo/automata"

// Empty returns an NFA accepting the empty string ε.
func empty() *automata.NFA {
	b := automata.NewNFABuilder().SetStart(0).SetFinal([]automata.State{1})
	b.AddTransition(0, automata.E, automata.E, []automata.State{1})

	return b.Build()
}
