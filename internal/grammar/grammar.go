// Package grammar implements data structures and algorithms for context-free grammars.
package grammar

import (
	"errors"
	"fmt"
	"strings"

	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/set"
)

// The empty string ε
var ε = String[Symbol]{}

// Symbol represents a grammar symbol (terminal or non-terminal).
type Symbol interface {
	fmt.Stringer

	Name() string
	Equals(Symbol) bool
}

// Terminal represents a terminal symbol.
// Terminals are the basic symbols from which strings of a language are formed.
// Token name or token for short are equivalent to terminal.
type Terminal string

// String implements the fmt.Stringer interface.
func (t Terminal) String() string {
	return t.Name()
}

// Name returns the name of terminal symbol.
func (t Terminal) Name() string {
	return string(t)
}

// Equals determines whether or not two terminal symbols are the same.
func (t Terminal) Equals(rhs Symbol) bool {
	if val, ok := rhs.(Terminal); ok {
		return t == val
	}
	return false
}

// NonTerminal represents a non-terminal symbol.
// Non-terminals are syntaxtic variables that denote sets of strings.
// Non-terminals impose a hierarchical structure on a language.
type NonTerminal string

// String implements the fmt.Stringer interface.
func (n NonTerminal) String() string {
	return n.Name()
}

// Name returns the name of non-terminal symbol.
func (n NonTerminal) Name() string {
	return string(n)
}

// Equals determines whether or not two non-terminal symbols are the same.
func (n NonTerminal) Equals(rhs Symbol) bool {
	if val, ok := rhs.(NonTerminal); ok {
		return n == val
	}
	return false
}

// String represent a string of grammar symbols.
type String[T Symbol] []T

// String implements the fmt.Stringer interface.
func (s String[T]) String() string {
	names := make([]string, len(s))
	for i, symbol := range s {
		names[i] = symbol.Name()
	}

	return strings.Join(names, " ")
}

// Equals determines whether or not two strings are the same.
func (s String[T]) Equals(rhs String[T]) bool {
	if len(s) != len(rhs) {
		return false
	}

	for i := range s {
		if !s[i].Equals(rhs[i]) {
			return false
		}
	}

	return true
}

// Terminals returns all terminal symbols of a string of symbols.
func (s String[Symbol]) Terminals() String[Terminal] {
	terms := String[Terminal]{}
	for _, symbol := range s {
		if term, ok := any(symbol).(Terminal); ok {
			terms = append(terms, term)
		}
	}
	return terms
}

// NonTerminals returns all non-terminal symbols of a string of symbols.
func (s String[Symbol]) NonTerminals() String[NonTerminal] {
	nonTerms := String[NonTerminal]{}
	for _, symbol := range s {
		if nonTerm, ok := any(symbol).(NonTerminal); ok {
			nonTerms = append(nonTerms, nonTerm)
		}
	}
	return nonTerms
}

// Production represents a production rule.
// The productions of a grammar determine how the terminals and non-terminals can be combined to form strings.
type Production struct {
	// Head or left side defines some of the strings denoted by the non-terminal symbol.
	Head NonTerminal
	// Body or right side describes one way in which strings of the non-terminal at the head can be constructed.
	Body String[Symbol]
}

// String implements the fmt.Stringer interface.
func (p Production) String() string {
	if len(p.Body) == 0 {
		return fmt.Sprintf("%s → ε", p.Head)
	}
	return fmt.Sprintf("%s → %s", p.Head, p.Body)
}

// Equals determines whether or not two production rules are the same.
func (p Production) Equals(rhs Production) bool {
	return p.Head.Equals(rhs.Head) && p.Body.Equals(rhs.Body)
}

// IsSingle determines whether or not a production rule is a single production (unit production).
//
// A single production (unit production) is a production whose body is a single non-terminal (A → B).
func (p Production) IsSingle() bool {
	if len(p.Body) == 1 {
		if _, ok := p.Body[0].(NonTerminal); ok {
			return true
		}
	}

	return false
}

// CFG represents a context-free grammar in formal language theory.
type CFG struct {
	Terminals    set.Set[Terminal]
	NonTerminals set.Set[NonTerminal]
	Productions  set.Set[Production]
	Start        NonTerminal
}

// New creates a new context-free grammar.
func New(terms []Terminal, nonTerms []NonTerminal, prods []Production, start NonTerminal) CFG {
	g := CFG{
		Terminals:    set.New[Terminal](generic.NewEqualFunc[Terminal]()),
		NonTerminals: set.New[NonTerminal](generic.NewEqualFunc[NonTerminal]()),
		Productions: set.New[Production](func(lhs, rhs Production) bool {
			return lhs.Equals(rhs)
		}),
		Start: start,
	}

	g.Terminals.Add(terms...)
	g.NonTerminals.Add(nonTerms...)
	g.Productions.Add(prods...)

	return g
}

// String implements the fmt.Stringer interface.
func (g CFG) String() string {
	var b strings.Builder

	terms := make([]string, g.Terminals.Cardinality())
	for i, t := range g.Terminals.Members() {
		terms[i] = t.String()
	}

	nonTerms := make([]string, g.NonTerminals.Cardinality())
	for i, n := range g.NonTerminals.Members() {
		nonTerms[i] = n.String()
	}

	fmt.Fprintf(&b, "Terminal Symbols: %s\n", strings.Join(terms, " "))
	fmt.Fprintf(&b, "Non-Terminal Symbols: %s\n", strings.Join(nonTerms, " "))
	fmt.Fprintf(&b, "Start Symbol: %s\n", g.Start)
	fmt.Fprintln(&b, "Production Rules:")

	for _, p := range g.Productions.Members() {
		fmt.Fprintf(&b, "  %s\n", p)
	}

	return b.String()
}

// verify receives a context-free grammar and determines whether or not it is valid.
// If the given grammar is invalid, an error with a descriptive message will be returned.
func (g CFG) verify() error {
	var err error

	getPredicate := func(n NonTerminal) set.Predicate[Production] {
		return func(p Production) bool {
			return p.Head.Equals(n)
		}
	}

	// Check if the start symbol is in the set of non-terminal symbols
	if !g.NonTerminals.Contains(g.Start) {
		err = errors.Join(err, fmt.Errorf("start symbol %q not in the set of non-terminal symbols", g.Start))
	}

	// Check if there is at least one production rule for the start symbol
	if !g.Productions.Any(getPredicate(g.Start)) {
		err = errors.Join(err, fmt.Errorf("no production rule for start symbol %q", g.Start))
	}

	// Check if there is at least one prodcution rule for every non-terminal symbol
	for _, n := range g.NonTerminals.Members() {
		if !g.Productions.Any(getPredicate(n)) {
			err = errors.Join(err, fmt.Errorf("no production rule for non-terminal symbol %q", n))
		}
	}

	for _, p := range g.Productions.Members() {
		// Check if the head of production rule is in the set of non-terminal symbols
		if !g.NonTerminals.Contains(p.Head) {
			err = errors.Join(err, fmt.Errorf("production head %q not in the set of non-terminal symbols", p.Head))
		}

		// Check if every symbol in the body of production rule is either in the set of terminal or non-terminal symbols
		for _, s := range p.Body {
			if v, ok := s.(Terminal); ok && !g.Terminals.Contains(v) {
				err = errors.Join(err, fmt.Errorf("terminal symbol %q not in the set of terminal symbols", v))
			}

			if v, ok := s.(NonTerminal); ok && !g.NonTerminals.Contains(v) {
				err = errors.Join(err, fmt.Errorf("non-terminal symbol %q not in the set of non-terminal symbols", v))
			}
		}
	}

	return err
}

// findNullableNonTerminals finds and returns all non-terminal symbols in a context-free grammar
// that can derive the empty string ε in one or more steps in which A ⇒* ε for some non-terminal A.
func (g CFG) findNullableNonTerminals() set.Set[NonTerminal] {
	// Define a set for all non-terminals that can derive the empty string ε
	nullable := set.New[NonTerminal](generic.NewEqualFunc[NonTerminal]())

	// We need to clone the production rules set since we want to modify it
	prods := g.Productions.Clone()

	for updated := true; updated; {
		updated = false

		// Iterate through each production rule of the form A → α,
		// where A is a non-terminal symbol and α is a string of terminals and non-terminals.
		for _, p := range prods.Members() {
			if len(p.Body) == 0 {
				// α is the empty string ε, add A to the nullable set.
				nullable.Add(p.Head)
				prods.Remove(p)
				updated = true
			} else if n := p.Body.NonTerminals(); len(n) == len(p.Body) && nullable.Contains(n...) {
				// α consists of only non-terminal symbols already in the nullable set, add A to the nullable set.
				nullable.Add(p.Head)
				prods.Remove(p)
				updated = true
			}
		}
	}

	return nullable
}

// EliminateEmptyProductions receives a context-free grammar and returns an equivalent ε-free grammar.
//
// An ε-production is any production of the form A → ε.
func (g CFG) EliminateEmptyProductions() CFG {
	nullable := g.findNullableNonTerminals()

	// The new set of production rules
	prods := set.New[Production](func(lhs, rhs Production) bool {
		return lhs.Equals(rhs)
	})

	// Iterate through each production rule in the input grammar.
	// For each production rule of the form A → α,
	// generate all possible combinations of α, excluding symbols that are in the nullable set.
	for _, p := range g.Productions.Members() {
		// Ignore ε-production rules of the form A → ε
		// Only consider the production rules of the form A → α
		if len(p.Body) > 0 {
			// alt holds all possible combinations of the right-hand side of a production rule.
			alt := []String[Symbol]{ε}
			var aux []String[Symbol]

			// Every nullable non-terminal symbol creates two possibilities, once by including it and once by excluding it.
			for _, sym := range p.Body {
				for _, β := range alt {
					if v, ok := sym.(NonTerminal); ok && nullable.Contains(v) {
						aux = append(aux, β)
					}
					aux = append(aux, append(β, sym))
				}
				alt, aux = aux, nil
			}

			for _, β := range alt {
				// Skip ε-production rules of the form A → ε
				if len(β) > 0 {
					prods.Add(Production{p.Head, β})
				}
			}
		}
	}

	return CFG{
		Terminals:    g.Terminals.Clone(),
		NonTerminals: g.NonTerminals.Clone(),
		Productions:  prods,
		Start:        g.Start,
	}
}
