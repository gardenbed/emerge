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
	IsTerminal() bool
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
	if v, ok := rhs.(Terminal); ok {
		return t == v
	}
	return false
}

// IsTerminal always returns true for terminal symbols.
func (t Terminal) IsTerminal() bool {
	return true
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
	if v, ok := rhs.(NonTerminal); ok {
		return n == v
	}
	return false
}

// IsTerminal always returns false for non-terminal symbols.
func (n NonTerminal) IsTerminal() bool {
	return false
}

// String represent a string of grammar symbols.
type String[T Symbol] []T

// String implements the fmt.Stringer interface.
func (s String[T]) String() string {
	if len(s) == 0 {
		return "ε"
	}

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
	for _, sym := range s {
		if v, ok := any(sym).(Terminal); ok {
			terms = append(terms, v)
		}
	}
	return terms
}

// NonTerminals returns all non-terminal symbols of a string of symbols.
func (s String[Symbol]) NonTerminals() String[NonTerminal] {
	nonTerms := String[NonTerminal]{}
	for _, sym := range s {
		if v, ok := any(sym).(NonTerminal); ok {
			nonTerms = append(nonTerms, v)
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

// IsEmpty determines whether or not a production rule is an empty production (ε-production).
//
// An empty production (ε-production) is any production of the form A → ε.
func (p Production) IsEmpty() bool {
	return len(p.Body) == 0
}

// IsSingle determines whether or not a production rule is a single production (unit production).
//
// A single production (unit production) is a production whose body is a single non-terminal (A → B).
func (p Production) IsSingle() bool {
	return len(p.Body) == 1 && !p.Body[0].IsTerminal()
}

// Grammar represents a context-free grammar in formal language theory.
type Grammar struct {
	Terminals    set.Set[Terminal]
	NonTerminals set.Set[NonTerminal]
	Productions  set.Set[Production]
	Start        NonTerminal
}

// New creates a new context-free grammar.
func New(terms []Terminal, nonTerms []NonTerminal, prods []Production, start NonTerminal) Grammar {
	g := Grammar{
		Terminals:    set.New(generic.NewEqualFunc[Terminal]()),
		NonTerminals: set.New(generic.NewEqualFunc[NonTerminal]()),
		Productions: set.New(func(lhs, rhs Production) bool {
			return lhs.Equals(rhs)
		}),
		Start: start,
	}

	g.Terminals.Add(terms...)
	g.NonTerminals.Add(nonTerms...)
	g.Productions.Add(prods...)

	// TODO: Verify the grammar
	// if err := g.verify(); err != nil {
	// 	return err
	// }

	return g
}

// verify takes a context-free grammar and determines whether or not it is valid.
// If the given grammar is invalid, an error with a descriptive message will be returned.
func (g Grammar) verify() error {
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

// String implements the fmt.Stringer interface.
func (g Grammar) String() string {
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

// Equals determines whether or not two context-free grammars are the same.
func (g Grammar) Equals(rhs Grammar) bool {
	return g.Terminals.Equals(rhs.Terminals) &&
		g.NonTerminals.Equals(rhs.NonTerminals) &&
		g.Productions.Equals(rhs.Productions) &&
		g.Start.Equals(rhs.Start)
}

// nullableNonTerminals finds all non-terminal symbols in a context-free grammar
// that can derive the empty string ε in one or more steps (A ⇒* ε for some non-terminal A).
func (g Grammar) nullableNonTerminals() set.Set[NonTerminal] {
	// Define a set for all non-terminals that can derive the empty string ε
	nullable := set.New(generic.NewEqualFunc[NonTerminal]())

	for updated := true; updated; {
		updated = false

		// Iterate through each production rule of the form A → α,
		// where A is a non-terminal symbol and α is a string of terminals and non-terminals.
		for _, p := range g.Productions.Members() {
			// Skip the production rule if A is already in the nullable set.
			if nullable.Contains(p.Head) {
				continue
			}

			if p.IsEmpty() {
				// α is the empty string ε, add A to the nullable set.
				nullable.Add(p.Head)
				updated = true
			} else if n := p.Body.NonTerminals(); len(n) == len(p.Body) && nullable.Contains(n...) {
				// α consists of only non-terminal symbols already in the nullable set, add A to the nullable set.
				nullable.Add(p.Head)
				updated = true
			}
		}
	}

	return nullable
}

// EliminateEmptyProductions converts a context-free grammar into an equivalent ε-free grammar.
//
// An empty production (ε-production) is any production of the form A → ε.
func (g Grammar) EliminateEmptyProductions() Grammar {
	nullable := g.nullableNonTerminals()

	// Create a set for the new production rules.
	newProds := set.New(func(lhs, rhs Production) bool {
		return lhs.Equals(rhs)
	})

	// Iterate through each production rule in the input grammar.
	// For each production rule of the form A → α,
	//   generate all possible combinations of α, excluding symbols that are in the nullable set.
	for _, p := range g.Productions.Members() {
		// Ignore ε-production rules (A → ε)
		// Only consider the production rules of the form A → α
		if p.IsEmpty() {
			continue
		}

		// bodies holds all possible combinations of the right-hand side of a production rule.
		bodies, aux := []String[Symbol]{ε}, []String[Symbol]{}

		// Every nullable non-terminal symbol creates two possibilities, once by including and once by excluding it.
		for _, sym := range p.Body {
			v, ok := sym.(NonTerminal)
			nonTermNullable := ok && nullable.Contains(v)

			for _, β := range bodies {
				if nonTermNullable {
					aux = append(aux, β)
				}
				aux = append(aux, append(β, sym))
			}

			bodies, aux = aux, nil
		}

		for _, β := range bodies {
			// Skip ε-production rules (A → ε)
			if len(β) > 0 {
				newProds.Add(Production{p.Head, β})
			}
		}
	}

	// The set data structure automatically prevents duplicate items from being added.
	// Therefore, we don't need to worry about deduplicating the new production rules at this stage.

	start := g.Start
	terms := g.Terminals.Clone()
	nonTerms := g.NonTerminals.Clone()

	// If the start symbol of the grammer is nullable (S ⇒* ε),
	//   a new start symbol with an ε-production rule must be introduced (S′ → S | ε).
	// This guarantees that the resulting grammar generates the same language as the original grammar.
	if nullable.Contains(g.Start) {
		// TODO: Make sure the new start symbol does not already exist as a non-terminal
		start = NonTerminal(g.Start + "′")

		// TODO: It may be preferable to add the new elements to the beginning of the sets.
		nonTerms.Add(start)
		newProds.Add(Production{start, String[Symbol]{g.Start}}) // S′ → S
		newProds.Add(Production{start, ε})                       // S′ → ε
	}

	return Grammar{
		Terminals:    terms,
		NonTerminals: nonTerms,
		Productions:  newProds,
		Start:        start,
	}
}

// EliminateSingleProductions converts a context-free grammar into an equivalent single-production-free grammar.
//
// A single production a.k.a. unit production is a production rule whose body is a single non-terminal symbol (A → B).
func (g Grammar) EliminateSingleProductions() Grammar {
	// Map each non-terminal symbol to its production bodies for efficient access.
	prods := map[NonTerminal][]String[Symbol]{}
	for _, p := range g.Productions.Members() {
		prods[p.Head] = append(prods[p.Head], p.Body)
	}

	// Identify all single productions.
	singleProds := map[NonTerminal][]NonTerminal{}
	for _, p := range g.Productions.Members() {
		if p.IsSingle() {
			singleProds[p.Head] = append(singleProds[p.Head], p.Body[0].(NonTerminal))
		}
	}

	// Compute the transitive closure for all non-terminal symbols.
	// The transitive closure of a non-terminal A is the the set of all non-terminals B
	//   such that there exists a sequence of single productions starting from A and reaching B (i.e., A → B₁ → B₂ → ... → B).

	closure := make(map[NonTerminal]map[NonTerminal]bool, g.NonTerminals.Cardinality())

	// Initially, each non-terminal symbol is reachable from itself.
	for _, A := range g.NonTerminals.Members() {
		closure[A] = map[NonTerminal]bool{A: true}
	}

	// Next, add directly reachable non-terminal symbols from single productions.
	for A, nonTerms := range singleProds {
		for _, B := range nonTerms {
			closure[A][B] = true
		}
	}

	// Repeat until no new non-terminal symbols can be added to the closure set.
	for updated := true; updated; {
		updated = false

		for A, closureA := range closure {
			for B := range closureA {
				for next := range closure[B] {
					if !closureA[next] {
						closure[A][next] = true
						updated = true
					}
				}
			}
		}
	}

	// Create a set for the new production rules.
	newProds := set.New(func(lhs, rhs Production) bool {
		return lhs.Equals(rhs)
	})

	// For each production rule p of the form B → α, add a new production rule A → α
	//   if p is not a single production and B is in the transitive closure set of A.
	for A, closureA := range closure {
		for B := range closureA {
			for _, body := range prods[B] {
				// Skip single productions
				if len(body) != 1 || body[0].IsTerminal() {
					newProds.Add(Production{A, body})
				}
			}
		}
	}

	return Grammar{
		Terminals:    g.Terminals.Clone(),
		NonTerminals: g.NonTerminals.Clone(),
		Productions:  newProds,
		Start:        g.Start,
	}
}

// EliminateUnreachableProductions converts a context-free grammar into an equivalent grammar
//   with all unreachable productions and their associated non-terminal symbols removed.
//
// An unreachable production refers to a production rule in a grammar that cannot be used to derive any string starting from the start symbol.
func (g Grammar) EliminateUnreachableProductions() Grammar {
	return Grammar{}
}

// EliminateCycles converts a context-free grammar into an equivalent cycle-free grammar.
//
// A grammar is cyclic if it has derivations of one or more steps in which A ⇒* A for some non-terminal A.
func (g Grammar) EliminateCycles() Grammar {
	// TODO: Explain why
	return g.EliminateEmptyProductions().EliminateSingleProductions().EliminateUnreachableProductions()
}

// EliminateLeftRecursion converts a context-free grammar into an equivalent grammar with no left recursion.
//
// A grammar is left-recursive if it has a non-terminal A such that there is a derivation A ⇒+ Aα for some string.
// For top-down parsers, left recursion causes the parser to loop forever.
// Many bottom-up parsers also will not accept left-recursive grammars.
func (g Grammar) EliminateLeftRecursion() Grammar {
	return Grammar{}
}

// LeftFactor converts a context-free grammar into an equivalent left-factored grammar.
//
// Left factoring is a grammar transformation for producing a grammar suitable for top-down parsing.
// When the choice between two alternative A-productions is not clear, we may be able to rewrite the productions
// to defer the decision until enough of the input has been seen that we can make the right choice.
func (g Grammar) LeftFactor() Grammar {
	return Grammar{}
}

// ChomskyNormalForm converts a context-free grammar into an equivalent grammar in Chomsky Normal Form.
//
// A grammar is in Chomsky Normal Form (CNF) if every production is either of the form A → BC or A → a,
//   where A, B, and C are non-terminal symbols, and a is a terminal symbol
//   (with the possible exception of the empty string derived from the start symbol, S → ε).
func (g Grammar) ChomskyNormalForm() Grammar {
	return Grammar{}
}
