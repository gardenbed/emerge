// Package spec defines the models and methods for parsing EBNF specifications.
package spec

import (
	"fmt"
	"strings"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/parser/lr/canonical"
	"github.com/moorara/algo/parser/lr/lookahead"
	"github.com/moorara/algo/parser/lr/simple"
	"github.com/moorara/algo/sort"
	"github.com/moorara/algo/symboltable"

	"github.com/gardenbed/emerge/internal/regex/parser/nfa"
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
//
// The second return value associates each terminal to its set of final states in the DFA.
func (s *Spec) DFA() (*automata.DFA, []TerminalFinal, error) {
	errs := &errors.MultiError{
		Format: errors.BulletErrorFormat,
	}

	// Construct a DFA for each terminal.
	ds := make([]*automata.DFA, len(s.Definitions))
	for i, def := range s.Definitions {
		switch def.Kind {
		case StringDef:
			ds[i] = stringToDFA(def.Value)
		case RegexDef:
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
	dfa, finalMap := automata.UnionDFA(ds...)

	// Map each final state in the merged DFA to the terminal definition(s) it accepts.
	// A final state may correspond to multiple terminals if multiple DFAs recognize the same string.
	finalToDefs := make(map[automata.State][]*TerminalDef)
	for i, finals := range finalMap {
		for _, f := range finals {
			finalToDefs[f] = append(finalToDefs[f], s.Definitions[i])
		}
	}

	// Build a mapping from terminal definitions to their sets of final states,
	// ensuring each final state resolves to exactly one terminal definition.
	defToFinals := symboltable.NewRedBlack[*TerminalDef, []automata.State](
		func(lhs, rhs *TerminalDef) int {
			return grammar.CmpTerminal(lhs.Terminal, rhs.Terminal)
		},
		nil,
	)

	for f, defs := range finalToDefs {
		switch len(defs) {
		case 0:
		case 1:
			finals, _ := defToFinals.Get(defs[0])
			finals = append(finals, f)
			defToFinals.Put(defs[0], finals)
		default:
			// Prefer a string-based definition over regex-based definitions.
			// This ensures keywords take precedence over identifiers.
			stringDefs := generic.SelectMatch(defs, func(def *TerminalDef) bool {
				return def.Kind == StringDef
			})

			if len(stringDefs) == 1 {
				finals, _ := defToFinals.Get(stringDefs[0])
				finals = append(finals, f)
				defToFinals.Put(stringDefs[0], finals)
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

	// Associate final states back to a single terminal.
	termFinals := make([]TerminalFinal, 0, defToFinals.Size())
	for def, finals := range defToFinals.All() {
		termFinals = append(termFinals, TerminalFinal{
			Kind:     def.Kind,
			Terminal: def.Terminal,
			Final:    automata.NewStates(finals...),
		})
	}

	sort.Quick(termFinals, func(lhs, rhs TerminalFinal) int {
		return automata.CmpStates(lhs.Final, rhs.Final)
	})

	if err := errs.ErrorOrNil(); err != nil {
		return nil, nil, err
	}

	return dfa, termFinals, nil
}

// TerminalFinal associates a terminal with its set of final states in a DFA.
type TerminalFinal struct {
	Kind     TerminalDefKind
	Terminal grammar.Terminal
	Final    automata.States
}

func stringToDFA(value string) *automata.DFA {
	start := automata.State(0)
	b := automata.NewDFABuilder().SetStart(start)

	curr, next := start, start+1
	for _, r := range value {
		sym := automata.Symbol(r)
		b.AddTransition(curr, sym, sym, next)
		curr, next = next, next+1
	}

	final := []automata.State{curr}
	b.SetFinal(final)

	return b.Build()
}

func regexToDFA(regex string) (*automata.DFA, error) {
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
