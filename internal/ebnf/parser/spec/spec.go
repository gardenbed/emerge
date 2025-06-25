// Package spec defines the models and methods for parsing EBNF specifications.
package spec

import (
	"fmt"
	"strings"

	"github.com/gardenbed/emerge/internal/regex/parser/nfa"
	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
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
func (s *Spec) DFA() (*auto.DFA, map[grammar.Terminal][]auto.State, error) {
	errs := &errors.MultiError{
		Format: errors.BulletErrorFormat,
	}

	// Construct a DFA for each terminal.
	ds := make([]*auto.DFA, len(s.Definitions))
	for i, def := range s.Definitions {
		if !def.IsRegex {
			ds[i] = stringToDFA(def.Value)
		} else {
			var err error
			ds[i], err = regexToDFA(def.Value)
			if err != nil {
				errs = errors.Append(errs, fmt.Errorf("%s: %s", def.Terminal, err))
			}
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, nil, err
	}

	// Combine multiple DFAs into one, preserving state mappings.
	dfa, stateMap := auto.CombineDFA(ds...)

	// Map final states in the DFA to their corresponding terminal definitions.
	stateDefs := make(map[auto.State][]*TerminalDef)
	for i, finals := range stateMap {
		for _, f := range finals {
			stateDefs[f] = append(stateDefs[f], s.Definitions[i])
		}
	}

	// Map each terminal to a set of final states while ensuring each final state identifies a single terminal.
	termMap := make(map[grammar.Terminal][]auto.State)
	for f, defs := range stateDefs {
		switch len(defs) {
		case 0:
		case 1:
			a := defs[0].Terminal
			termMap[a] = append(termMap[a], f)
		default:
			// Prefer a string-based definition over regex-based definitions.
			// This ensures keywords take precedence over identifiers.
			strDefs := generic.SelectMatch(defs, func(def *TerminalDef) bool {
				return !def.IsRegex
			})

			if len(strDefs) == 1 {
				a := strDefs[0].Terminal
				termMap[a] = append(termMap[a], f)
			} else {
				poses := generic.Transform(defs, func(def *TerminalDef) string {
					return fmt.Sprintf("  %s: %s", def.Pos, def.Terminal)
				})

				errs = errors.Append(errs,
					fmt.Errorf("conflicting definitions capture the same string:\n%s", strings.Join(poses, "\n")),
				)
			}
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, nil, err
	}

	return dfa, termMap, nil
}

func stringToDFA(value string) *auto.DFA {
	start := auto.State(0)
	d := auto.NewDFA(start, nil)

	curr, next := start, start+1
	for _, r := range value {
		d.Add(curr, auto.Symbol(r), next)
		curr, next = next, next+1
	}

	d.Final = auto.NewStates(curr)

	return d
}

func regexToDFA(regex string) (*auto.DFA, error) {
	n, err := nfa.Parse(regex)
	if err != nil {
		return nil, err
	}

	d := n.ToDFA().Minimize().EliminateDeadStates().ReindexStates()

	return d, nil
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
