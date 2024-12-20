package grammar

import "strings"

// The empty string ε
var ε = String[Symbol]{}

// String represent a string of grammar symbols.
type String[T Symbol] []T

// String returns a string representation of a string of symbols.
func (s String[T]) String() string {
	if len(s) == 0 {
		return "ε"
	}

	names := make([]string, len(s))
	for i, symbol := range s {
		names[i] = symbol.String()
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

// HasPrefix checks whether a string starts with the given prefix.
func (s String[T]) HasPrefix(prefix String[T]) bool {
	ls, lp := len(s), len(prefix)
	return ls >= lp && s[:lp].Equals(prefix)
}

// HasSuffix checks whether a string ends with the given suffix.
func (s String[T]) HasSuffix(prefix String[T]) bool {
	ls, lp := len(s), len(prefix)
	return ls >= lp && s[ls-lp:].Equals(prefix)
}

// Append appends new symbols to the current string and returns a new string
func (s String[T]) Append(symbols ...T) String[T] {
	newS := make(String[T], len(s)+len(symbols))

	copy(newS, s)
	copy(newS[len(s):], symbols)

	return newS
}

// Concat concatenates the current string with one or more strings and returns a new string.
func (s String[T]) Concat(ss ...String[T]) String[T] {
	l := len(s)
	for _, t := range ss {
		l += len(t)
	}

	newS := make(String[T], l)

	copy(newS, s)
	i := len(s)
	for _, t := range ss {
		copy(newS[i:], t)
		i += len(t)
	}

	return newS
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

// LongestCommonPrefixOf computes the longest common prefix of a list of strings.
// If the input is empty or there is no common prefix, it returns the empty string ε.
func LongestCommonPrefixOf(ss ...String[Symbol]) String[Symbol] {
	if len(ss) == 0 {
		return ε
	}

	lcp := ss[0]

	for i := 1; i < len(ss); i++ {
		for !ss[i].HasPrefix(lcp) {
			lcp = lcp[:len(lcp)-1]
			if len(lcp) == 0 {
				return ε
			}
		}
	}

	return lcp
}
