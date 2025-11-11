// Package nfa provides a combinator parser for parsing regular expression into a non-deterministic finite automata.
//
// It implements the McNaughton-Yamada-Thompson algorithm to convert a regular expression to an NFA.
// The algorithm is syntax-directed, in the sense that it works recursively up the parse tree for the regular expression.
package nfa

import (
	"fmt"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/parser/combinator"

	"github.com/gardenbed/emerge/internal/char"
	"github.com/gardenbed/emerge/internal/regex/parser"
)

func Parse(regex string) (*automata.NFA, error) {
	m := new(mappers)
	p := parser.New(m)

	out, err := p.Parse(regex)
	if err != nil {
		return nil, fmt.Errorf("invalid regular expression: %s: %w", regex, err)
	}

	nfa := out.Result.Val.(*automata.NFA)

	return nfa, nil
}

//==================================================< MAPPERS >==================================================

// Anchor is the type for a regular expression anchor.
type Anchor int

const (
	StartOfString Anchor = iota + 1
	EndOfString

	bagKeyChar           combinator.BagKey = "char"
	bagKeyCharRanges     combinator.BagKey = "char_ranges"
	bagKeyLazyQuantifier combinator.BagKey = "lazy_quantifier"
	BagKeyStartOfString  combinator.BagKey = "start_of_string"
)

// mappers implements the parser.Mappers interface.
type mappers struct{}

func (m *mappers) ToAnyChar(r combinator.Result) (combinator.Result, error) {
	nfa, _ := charRangesToNFA(false, char.Classes["UNICODE"])

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
	}, nil
}

func (m *mappers) ToSingleChar(r combinator.Result) (combinator.Result, error) {
	c := r.Val.(rune)
	nfa, ranges := charRangesToNFA(false, char.RangeList{{c, c}})

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyChar:       c,
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToCharClass(r combinator.Result) (combinator.Result, error) {
	class := r.Val.(string)

	var nfa *automata.NFA
	var ranges char.RangeList

	switch class {
	case `\s`:
		nfa, ranges = charRangesToNFA(false, char.Classes[`\s`])
	case `\S`:
		nfa, ranges = charRangesToNFA(true, char.Classes[`\s`])
	case `\d`:
		nfa, ranges = charRangesToNFA(false, char.Classes[`\d`])
	case `\D`:
		nfa, ranges = charRangesToNFA(true, char.Classes[`\d`])
	case `\w`:
		nfa, ranges = charRangesToNFA(false, char.Classes[`\w`])
	case `\W`:
		nfa, ranges = charRangesToNFA(true, char.Classes[`\w`])
	default:
		return combinator.Result{}, fmt.Errorf("invalid character class: %s", class)
	}

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToASCIICharClass(r combinator.Result) (combinator.Result, error) {
	class := r.Val.(string)

	ranges, ok := char.Classes[class]
	if !ok {
		return combinator.Result{}, fmt.Errorf("invalid ASCII character class: %s", class)
	}

	nfa, ranges := charRangesToNFA(false, ranges)

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToUnicodeCategory(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToUnicodeCharClass(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r2, _ := r.Get(2)

	prop := r0.Val.(string)
	class := r2.Val.(string)

	ranges, ok := char.Classes[class]
	if !ok {
		return combinator.Result{}, fmt.Errorf("invalid Unicode character class: %s", class)
	}

	nfa, ranges := charRangesToNFA(prop == `\P`, ranges)

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToRepOp(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToUpperBound(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	var num *int
	if v, ok := r1.Val.(int); ok {
		num = &v
	}

	return combinator.Result{
		Val: num,
		Pos: r0.Pos,
	}, nil
}

func (m *mappers) ToRange(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	// The upper bound is same as the lower bound if no upper bound is specified (default)
	low := r1.Val.(int)
	up := &low

	// If an upper bound is specified, it can be either bounded or unbounded
	if v, ok := r2.Val.(*int); ok {
		up = v
	}

	if up != nil && low > *up {
		return combinator.Result{}, fmt.Errorf("invalid repetition range {%d,%d}", low, *up)
	}

	return combinator.Result{
		Val: tuple[int, *int]{
			p: low,
			q: up,
		},
		Pos: r0.Pos,
	}, nil
}

func (m *mappers) ToRepetition(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToQuantifier(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	// Check whether or not the lazy modifier is present
	_, lazy := r1.Val.(rune)

	return combinator.Result{
		Val: tuple[any, bool]{
			p: r0.Val,
			q: lazy,
		},
		Pos: r0.Pos,
	}, nil
}

func (m *mappers) ToCharInGroup(r combinator.Result) (combinator.Result, error) {
	c := r.Val.(rune)
	nfa, ranges := charRangesToNFA(false, char.RangeList{{c, c}})

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyChar:       c,
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToCharRange(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r2, _ := r.Get(2)

	lo, hi := r0.Bag[bagKeyChar].(rune), r2.Bag[bagKeyChar].(rune)

	if lo > hi {
		return combinator.Result{}, fmt.Errorf("invalid character range %c-%c", lo, hi)
	}

	nfa, ranges := charRangesToNFA(false, char.RangeList{{lo, hi}})

	return combinator.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, nil
}

func (m *mappers) ToCharGroupItem(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToCharGroup(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	// Check whether or not the negation modifier is present
	_, neg := r1.Val.(rune)

	// Collect all character ranges from the character group items
	var all char.RangeList
	for _, r := range r2.Val.(combinator.List) {
		if ranges, ok := r.Bag[bagKeyCharRanges].(char.RangeList); ok {
			all = append(all, ranges...)
		}
	}

	nfa, _ := charRangesToNFA(neg, all.Dedup())

	return combinator.Result{
		Val: nfa,
		Pos: r0.Pos,
	}, nil
}

func (m *mappers) ToMatchItem(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToMatch(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	nfa := r0.Val.(*automata.NFA)
	var bag combinator.Bag

	if t, ok := r1.Val.(tuple[any, bool]); ok {
		nfa = quantifyNFA(nfa, t.p)
		if lazy := t.q; lazy {
			bag = combinator.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return combinator.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: bag,
	}, nil
}

func (m *mappers) ToGroup(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r3, _ := r.Get(3)

	nfa := r1.Val.(*automata.NFA)
	var bag combinator.Bag

	if t, ok := r3.Val.(tuple[any, bool]); ok {
		nfa = quantifyNFA(nfa, t.p)
		if lazy := t.q; lazy {
			bag = combinator.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return combinator.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: bag,
	}, nil
}

func (m *mappers) ToAnchor(r combinator.Result) (combinator.Result, error) {
	c := r.Val.(rune)

	var anchor Anchor
	switch c {
	case '$': // end-of-string
		anchor = EndOfString
	}

	return combinator.Result{
		Val: anchor,
		Pos: r.Pos,
	}, nil
}

func (m *mappers) ToSubexprItem(r combinator.Result) (combinator.Result, error) {
	// Passing the result up the parsing chain
	return r, nil
}

func (m *mappers) ToSubexpr(r combinator.Result) (combinator.Result, error) {
	items := r.Val.(combinator.List)

	ns := []*automata.NFA{}
	for _, r := range items {
		// TODO: Anchor result value is not a node
		if n, ok := r.Val.(*automata.NFA); ok {
			ns = append(ns, n)
		}
	}

	nfa := automata.ConcatNFA(ns...)

	return combinator.Result{
		Val: nfa,
		Pos: r.Pos,
	}, nil
}

func (m *mappers) ToExpr(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	nfa := r0.Val.(*automata.NFA)

	if _, ok := r1.Val.(combinator.List); ok {
		r11, _ := r1.Get(1)
		expr := r11.Val.(*automata.NFA)
		nfa = nfa.Union(expr)
	}

	return combinator.Result{
		Val: nfa,
		Pos: r0.Pos,
	}, nil
}

func (m *mappers) ToRegex(r combinator.Result) (combinator.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	pos := r1.Pos
	var bag combinator.Bag

	// Check whether or not the start-of-string is present
	if _, sos := r0.Val.(rune); sos {
		pos = r0.Pos
		bag = combinator.Bag{
			BagKeyStartOfString: true,
		}
	}

	nfa := r1.Val.(*automata.NFA)

	return combinator.Result{
		Val: nfa,
		Pos: pos,
		Bag: bag,
	}, nil
}

//==================================================< HELPERS >==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

// empty returns an NFA accepting the empty string ε.
func empty() *automata.NFA {
	b := automata.NewNFABuilder().SetStart(0).SetFinal([]automata.State{1})
	b.AddTransition(0, automata.E, automata.E, []automata.State{1})

	return b.Build()
}

// charRangesToNFA converts a list of character ranges into an NFA.
//
// If neg is true, the NFA accepts all runes except those in the given ranges.
// It also returns the list of rune ranges that the NFA accepts.
func charRangesToNFA(neg bool, ranges char.RangeList) (*automata.NFA, char.RangeList) {
	b := automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1})

	if neg {
		ranges = char.Classes["UNICODE"].Exclude(ranges)
	}

	for _, r := range ranges {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(0, lo, hi, []automata.State{1})
	}

	return b.Build(), ranges
}

// quantifyNFA applies a quantifier to an NFA and returns the resulting NFA.
func quantifyNFA(n *automata.NFA, q any) *automata.NFA {
	var nfa *automata.NFA

	switch rep := q.(type) {
	// Simple repetition
	case rune:
		switch rep {
		case '?':
			nfa = empty().Union(n)
		case '*':
			nfa = n.Star()
		case '+':
			nfa = n.Concat(n.Star())
		}

	// Range repetition
	case tuple[int, *int]:
		low, up := rep.p, rep.q
		ns := []*automata.NFA{}

		for i := 0; i < low; i++ {
			ns = append(ns, n)
		}

		if up == nil {
			ns = append(ns, n.Star())
		} else {
			for i := 0; i < *up-low; i++ {
				ns = append(ns, empty().Union(n))
			}
		}

		nfa = automata.ConcatNFA(ns...)
	}

	return nfa
}
