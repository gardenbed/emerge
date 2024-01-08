// Package grammar implements data structures and algorithms for context-free grammars.
package grammar

import (
	"fmt"
	"strings"

	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/set"
)

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

// String implements the fmt.Stringer interface.
func (t Terminal) String() string {
	return t.Name()
}

// NonTerminal represents a non-terminal symbol.
// Non-terminals are syntaxtic variables that denote sets of strings.
// Non-terminals impose a hierarchical structure on a language.
type NonTerminal string

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

// String implements the fmt.Stringer interface.
func (n NonTerminal) String() string {
	return n.Name()
}

// String represent a string of grammar symbols.
type String[T Symbol] []T

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

// String implements the fmt.Stringer interface.
func (s String[T]) String() string {
	names := make([]string, len(s))
	for i, symbol := range s {
		names[i] = symbol.Name()
	}

	return strings.Join(names, " ")
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

// Equals determines whether or not two production rules are the same.
func (p Production) Equals(rhs Production) bool {
	return p.Head.Equals(rhs.Head) && p.Body.Equals(rhs.Body)
}

// String implements the fmt.Stringer interface.
func (p Production) String() string {
	return fmt.Sprintf("%s → %s", p.Head, p.Body)
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
