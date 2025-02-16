package spec

import (
	"fmt"

	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/parser/lr/canonical"
	"github.com/moorara/algo/parser/lr/lookahead"
	"github.com/moorara/algo/parser/lr/simple"
	"github.com/moorara/algo/sort"
)

// Spec contains the result of a successful input parsing.
type Spec struct {
	Name        string
	Definitions []*TerminalDef
	Grammar     *grammar.CFG
	Precedences lr.PrecedenceLevels
}

// DFA constructs a deterministic finite automaton (DFA)
// for recognizing all terminal symbols (tokens) in the grammar of the spec.
func (s *Spec) DFA() *auto.DFA {
	n0 := auto.NewNFA(0, nil)

	nfa := make([]*auto.NFA, len(s.Definitions))
	for i, def := range s.Definitions {
		nfa[i] = def.ToNFA()
	}

	return n0.Union(nfa...).ToDFA().EliminateDeadStates().ReindexStates()
}

// Productions returns an ordered list of all production rules in the grammar of the spec.
func (s *Spec) Productions() []*grammar.Production {
	prods := generic.Collect1(s.Grammar.Productions.All())
	sort.Quick(prods, grammar.CmpProduction)

	return prods
}

// SLRParsingTable builds and returns the SLR(1) (Simple LR) parsing table
// for the grammar and precedences in the spec.
func (s *Spec) SLRParsingTable() (*lr.ParsingTable, error) {
	T, err := simple.BuildParsingTable(s.Grammar, s.Precedences)
	if err != nil {
		return nil, fmt.Errorf("error on building SLR(1) parsing table:\n%s", err)
	}

	return T, nil
}

// LALRParsingTable builds and returns the LALR(1) (Lookahead LR) parsing table
// for the grammar and precedences in the spec.
func (s *Spec) LALRParsingTable() (*lr.ParsingTable, error) {
	T, err := lookahead.BuildParsingTable(s.Grammar, s.Precedences)
	if err != nil {
		return nil, fmt.Errorf("error on building LALR(1) parsing table:\n%s", err)
	}

	return T, nil
}

// GLRParsingTable builds and returns the GLR(1) (Canonical LR a.k.a. Generalized LR) parsing table
// for the grammar and precedences in the spec.
func (s *Spec) GLRParsingTable() (*lr.ParsingTable, error) {
	T, err := canonical.BuildParsingTable(s.Grammar, s.Precedences)
	if err != nil {
		return nil, fmt.Errorf("error on building GLR(1) parsing table:\n%s", err)
	}

	return T, nil
}
