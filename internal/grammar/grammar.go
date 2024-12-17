// Package grammar implements data structures and algorithms for context-free grammars.
package grammar

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/hash"
	"github.com/moorara/algo/set"
	"github.com/moorara/algo/sort"
	"github.com/moorara/algo/symboltable"
)

// The empty string ε
var ε = String[Symbol]{}

var (
	eqTerminal      = generic.NewEqualFunc[Terminal]()
	eqNonTerminal   = generic.NewEqualFunc[NonTerminal]()
	cmpString       = generic.NewCompareFunc[string]()
	cmpNonTerminal  = generic.NewCompareFunc[NonTerminal]()
	hashNonTerminal = hash.HashFuncForString[NonTerminal](nil)

	eqProduction = func(lhs, rhs Production) bool {
		return lhs.Equals(rhs)
	}

	eqProductionSet = func(lhs, rhs set.Set[Production]) bool {
		return lhs.Equals(rhs)
	}
)

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

// Productions represents a mapping of non-terminal symbols to their corresponding production rules.
type Productions symboltable.SymbolTable[NonTerminal, set.Set[Production]]

// NewProductions creates a new instance of the Productions type.
func NewProductions() Productions {
	return symboltable.NewQuadraticHashTable[NonTerminal, set.Set[Production]](
		hashNonTerminal,
		eqNonTerminal,
		eqProductionSet,
		symboltable.HashOpts{},
	)
}

// Grammar represents a context-free grammar in formal language theory.
type Grammar struct {
	Terminals    set.Set[Terminal]
	NonTerminals set.Set[NonTerminal]
	Productions  Productions
	Start        NonTerminal
}

// New creates a new context-free grammar.
func New(terms []Terminal, nonTerms []NonTerminal, prods []Production, start NonTerminal) Grammar {
	g := Grammar{
		Terminals:    set.New(eqTerminal),
		NonTerminals: set.New(eqNonTerminal),
		Productions:  NewProductions(),
		Start:        start,
	}

	g.Terminals.Add(terms...)
	g.NonTerminals.Add(nonTerms...)
	g.AddProduction(prods...)

	return g
}

// AddProduction adds one or more productions to a context-free grammar.
func (g Grammar) AddProduction(ps ...Production) {
	for _, p := range ps {
		if _, ok := g.Productions.Get(p.Head); !ok {
			g.Productions.Put(p.Head, set.New[Production](eqProduction))
		}

		list, _ := g.Productions.Get(p.Head)
		list.Add(p)
	}
}

// AllProductions returns an iterator that yields all productions in the context-free grammar.
func (g Grammar) AllProductions() iter.Seq[Production] {
	return func(yield func(Production) bool) {
		for _, list := range g.Productions.All() {
			for p := range list.All() {
				if !yield(p) {
					return
				}
			}
		}
	}
}

// verify takes a context-free grammar and determines whether or not it is valid.
// If the given grammar is invalid, an error with a descriptive message will be returned.
func (g Grammar) Verify() error {
	var err error

	getPredicate := func(n NonTerminal) generic.Predicate2[NonTerminal, set.Set[Production]] {
		return func(head NonTerminal, _ set.Set[Production]) bool {
			return head.Equals(n)
		}
	}

	// Check if the start symbol is in the set of non-terminal symbols.
	if !g.NonTerminals.Contains(g.Start) {
		err = errors.Join(err, fmt.Errorf("start symbol %q not in the set of non-terminal symbols", g.Start))
	}

	// Check if there is at least one production rule for the start symbol.
	if !g.Productions.AnyMatch(getPredicate(g.Start)) {
		err = errors.Join(err, fmt.Errorf("no production rule for start symbol %q", g.Start))
	}

	// Check if there is at least one prodcution rule for every non-terminal symbol.
	for n := range g.NonTerminals.All() {
		if !g.Productions.AnyMatch(getPredicate(n)) {
			err = errors.Join(err, fmt.Errorf("no production rule for non-terminal symbol %q", n))
		}
	}

	for head, list := range g.Productions.All() {
		// Check if the head of production rule is in the set of non-terminal symbols.
		if !g.NonTerminals.Contains(head) {
			err = errors.Join(err, fmt.Errorf("production head %q not in the set of non-terminal symbols", head))
		}

		for p := range list.All() {
			if !p.Head.Equals(head) {
				err = errors.Join(err, fmt.Errorf("production head %q not matching %q", p.Head, head))
			}

			// Check if every symbol in the body of production rule is either in the set of terminal or non-terminal symbols.
			for _, s := range p.Body {
				if v, ok := s.(Terminal); ok && !g.Terminals.Contains(v) {
					err = errors.Join(err, fmt.Errorf("terminal symbol %q not in the set of terminal symbols", v))
				}

				if v, ok := s.(NonTerminal); ok && !g.NonTerminals.Contains(v) {
					err = errors.Join(err, fmt.Errorf("non-terminal symbol %q not in the set of non-terminal symbols", v))
				}
			}
		}
	}

	return err
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
	nullable := set.New(eqNonTerminal)

	for updated := true; updated; {
		updated = false

		// Iterate through each production rule of the form A → α,
		// where A is a non-terminal symbol and α is a string of terminals and non-terminals.
		for head, list := range g.Productions.All() {
			// Skip the production rule if A is already in the nullable set.
			if nullable.Contains(head) {
				continue
			}

			for p := range list.All() {
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
	}

	return nullable
}

// EliminateEmptyProductions converts a context-free grammar into an equivalent ε-free grammar.
//
// An empty production (ε-production) is any production of the form A → ε.
func (g Grammar) EliminateEmptyProductions() Grammar {
	nullable := g.nullableNonTerminals()

	newGrammar := Grammar{
		Terminals:    g.Terminals.Clone(),
		NonTerminals: g.NonTerminals.Clone(),
		Productions:  NewProductions(),
		Start:        g.Start,
	}

	// Iterate through each production rule in the input grammar.
	// For each production rule of the form A → α,
	//   generate all possible combinations of α, excluding symbols that are in the nullable set.
	for p := range g.AllProductions() {
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
				newGrammar.AddProduction(Production{p.Head, β})
			}
		}
	}

	// The set data structure automatically prevents duplicate items from being added.
	// Therefore, we don't need to worry about deduplicating the new production rules at this stage.

	// If the start symbol of the grammer is nullable (S ⇒* ε),
	//   a new start symbol with an ε-production rule must be introduced (S′ → S | ε).
	// This guarantees that the resulting grammar generates the same language as the original grammar.
	if nullable.Contains(g.Start) {
		newStart, _ := g.generateNewNonTerminal(g.Start, "′", "_new")

		newGrammar.Start = newStart
		newGrammar.NonTerminals.Add(newStart)
		newGrammar.AddProduction(Production{newStart, String[Symbol]{g.Start}}) // S′ → S
		newGrammar.AddProduction(Production{newStart, ε})                       // S′ → ε
	}

	return newGrammar
}

// EliminateSingleProductions converts a context-free grammar into an equivalent single-production-free grammar.
//
// A single production a.k.a. unit production is a production rule whose body is a single non-terminal symbol (A → B).
func (g Grammar) EliminateSingleProductions() Grammar {
	// Identify all single productions.
	singleProds := map[NonTerminal][]NonTerminal{}
	for p := range g.AllProductions() {
		if p.IsSingle() {
			singleProds[p.Head] = append(singleProds[p.Head], p.Body[0].(NonTerminal))
		}
	}

	// Compute the transitive closure for all non-terminal symbols.
	// The transitive closure of a non-terminal A is the the set of all non-terminals B
	//   such that there exists a sequence of single productions starting from A and reaching B (i.e., A → B₁ → B₂ → ... → B).

	closure := make(map[NonTerminal]map[NonTerminal]bool, g.NonTerminals.Cardinality())

	// Initially, each non-terminal symbol is reachable from itself.
	for A := range g.NonTerminals.All() {
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

	newGrammar := Grammar{
		Terminals:    g.Terminals.Clone(),
		NonTerminals: g.NonTerminals.Clone(),
		Productions:  NewProductions(),
		Start:        g.Start,
	}

	// For each production rule p of the form B → α, add a new production rule A → α
	//   if p is not a single production and B is in the transitive closure set of A.
	for A, closureA := range closure {
		for B := range closureA {
			list, _ := g.Productions.Get(B)
			for p := range list.All() {
				// Skip single productions
				if !p.IsSingle() {
					newGrammar.AddProduction(Production{A, p.Body})
				}
			}
		}
	}

	return newGrammar
}

// EliminateUnreachableProductions converts a context-free grammar into an equivalent grammar
// with all unreachable productions and their associated non-terminal symbols removed.
//
// An unreachable production refers to a production rule in a grammar
// that cannot be used to derive any string starting from the start symbol.
func (g Grammar) EliminateUnreachableProductions() Grammar {
	reachable := set.New(eqNonTerminal)
	reachable.Add(g.Start)

	// Reppeat until no new non-terminal is added to reachable:
	//   For each production rule of the form A → α:
	//     If A is in reachable, add all non-terminal in α to reachable.
	for updated := true; updated; {
		updated = false

		for p := range g.AllProductions() {
			if reachable.Contains(p.Head) {
				for _, n := range p.Body.NonTerminals() {
					if !reachable.Contains(n) {
						reachable.Add(n)
						updated = true
					}
				}
			}
		}
	}

	newGrammar := Grammar{
		Terminals:    g.Terminals.Clone(),
		NonTerminals: reachable,
		Productions:  NewProductions(),
		Start:        g.Start,
	}

	// Only consider the reachable production rules.
	for p := range g.AllProductions() {
		if reachable.Contains(p.Head) {
			newGrammar.AddProduction(p)
		}
	}

	return newGrammar
}

// EliminateCycles converts a context-free grammar into an equivalent cycle-free grammar.
//
// A grammar is cyclic if it has derivations of one or more steps in which A ⇒* A for some non-terminal A.
func (g Grammar) EliminateCycles() Grammar {
	// Single productions (unit productions) can create cycles in a grammar.
	// Eliminating empty productions (ε-productions) may introduce additional single productions,
	// so it is necessary to eliminate empty productions first, followed by single productions.
	// After removing single productions, some productions may become unreachable.
	// These unreachable productions should then be removed from the grammar.
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
//
//	where A, B, and C are non-terminal symbols, and a is a terminal symbol
//	(with the possible exception of the empty string derived from the start symbol, S → ε).
func (g Grammar) ChomskyNormalForm() Grammar {
	return Grammar{}
}

// String returns a string representation of a context-free grammar.
func (g Grammar) String() string {
	var b strings.Builder

	terms := g.orderTerminals()
	visited, unvisited, nonTerms := g.orderNonTerminals()

	fmt.Fprintf(&b, "Terminal Symbols: %s\n", strings.Join(terms, " "))
	fmt.Fprintf(&b, "Non-Terminal Symbols: %s\n", strings.Join(nonTerms, " "))
	fmt.Fprintf(&b, "Start Symbol: %s\n", g.Start)
	fmt.Fprintln(&b, "Production Rules:")

	for _, n := range visited {
		list, _ := g.Productions.Get(n)
		prods := orderProductions(list)

		for _, p := range prods {
			fmt.Fprintf(&b, "  %s\n", p)
		}
	}

	for _, n := range unvisited {
		list, _ := g.Productions.Get(n)
		prods := orderProductions(list)

		for _, p := range prods {
			fmt.Fprintf(&b, "  %s\n", p)
		}
	}

	return b.String()
}

// generateNewNonTerminal generates a new non-terminal symbol by appending a list of suffixes to a given prefix.
// It returns the first non-terminal that does not already exist in the set of non-terminals, along with a boolean indicating success.
func (g Grammar) generateNewNonTerminal(prefix NonTerminal, suffixes ...string) (NonTerminal, bool) {
	for _, suffix := range suffixes {
		nonTerm := NonTerminal(string(prefix) + suffix)
		if !g.NonTerminals.Contains(nonTerm) {
			return nonTerm, true
		}
	}

	return NonTerminal(""), false
}

// orderTerminals orders the unordered set of grammar terminals in a deterministic way.
//
// The goal of this function is to ensure a consistent and deterministic order for any given set of terminals.
func (g Grammar) orderTerminals() []string {
	terms := make([]string, 0)
	for t := range g.Terminals.All() {
		terms = append(terms, t.String())
	}

	// Sort terminals alphabetically based on the string representation of them.
	sort.Quick[string](terms, cmpString)

	return terms
}

// orderTerminals orders the unordered set of grammar non-terminals in a deterministic way.
//
// The goal of this function is to ensure a consistent and deterministic order for any given set of non-terminals.
func (g Grammar) orderNonTerminals() ([]NonTerminal, []NonTerminal, []string) {
	visited := make([]NonTerminal, 0)
	isVisited := func(n NonTerminal) bool {
		for _, v := range visited {
			if v == n {
				return true
			}
		}
		return false
	}

	visited = append(visited, g.Start)

	// Reppeat until no new non-terminal is added to visited:
	//   For each production rule of the form A → α:
	//     If A is in visited, add all non-terminal in α to visited.
	for updated := true; updated; {
		updated = false
		for _, list := range g.Productions.All() {
			for _, p := range orderProductions(list) {
				if isVisited(p.Head) {
					for _, n := range p.Body.NonTerminals() {
						if !isVisited(n) {
							visited = append(visited, n)
							updated = true
						}
					}
				}
			}
		}
	}

	// Identify any unvisited non-terminals in the grammar.
	unvisited := make([]NonTerminal, 0)
	for n := range g.NonTerminals.All() {
		if !isVisited(n) {
			unvisited = append(unvisited, n)
		}
	}

	// Sort unvisited non-terminals alphabetically based on the string representation of them.
	sort.Quick[NonTerminal](unvisited, cmpNonTerminal)

	allNonTerms := make([]string, 0)
	for _, n := range visited {
		allNonTerms = append(allNonTerms, n.String())
	}
	for _, n := range unvisited {
		allNonTerms = append(allNonTerms, n.String())
	}

	return visited, unvisited, allNonTerms
}

// orderProductions orders an unordered set of production rules in a deterministic way.
// This method assumes that all provided productions belong to the same head non-terminal.
//
// The ordering criteria are as follows:
//
//  1. Productions whose bodies contain more non-terminal symbols are prioritized first.
//  2. If two productions have the same number of non-terminals, those with more terminal symbols in the body come first.
//  3. If two productions have the same number of non-terminals and terminals, they are ordered alphabetically based on the symbols in their bodies.
//
// The goal of this function is to ensure a consistent and deterministic order for any given set of production rules.
func orderProductions(set set.Set[Production]) []Production {
	// Collect all production rules into a slice from the set iterator.
	prods := slices.Collect(set.All())

	// Sort the productions using a custom comparison function.
	sort.Quick[Production](prods, func(lhs, rhs Production) int {
		// First, compare based on the number of non-terminal symbols in the body.
		lhsNonTermsLen, rhsNonTermsLen := len(lhs.Body.NonTerminals()), len(rhs.Body.NonTerminals())
		if lhsNonTermsLen > rhsNonTermsLen {
			return -1
		} else if rhsNonTermsLen > lhsNonTermsLen {
			return 1
		}

		// Next, if the number of non-terminals is the same,
		//   compare based on the number of terminal symbols.
		lhsTermsLen, rhsTermsLen := len(lhs.Body.Terminals()), len(rhs.Body.Terminals())
		if lhsTermsLen > rhsTermsLen {
			return -1
		} else if rhsTermsLen > lhsTermsLen {
			return 1
		}

		// Then, if the number of terminals is also the same,
		//   compare alphabetically based on the string representation of the bodies.
		lhsString, rhsString := lhs.String(), rhs.String()
		if lhsString < rhsString {
			return -1
		} else if rhsString < lhsString {
			return 1
		}

		return 0
	})

	return prods
}
