package nfa

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	auto "github.com/moorara/algo/automata"

	comb "github.com/gardenbed/emerge/internal/combinator"
	"github.com/gardenbed/emerge/internal/regex/parser"
)

// Anchor is the type for a regular expression anchor.
type Anchor int

const (
	StartOfString Anchor = iota + 1
	EndOfString
)

const (
	bagKeyChars          comb.BagKey = "chars"
	bagKeyLazyQuantifier comb.BagKey = "lazy_quantifier"
	BagKeyStartOfString  comb.BagKey = "start_of_string"
)

func Parse(regex string) (*auto.NFA, error) {
	m := new(mappers)
	p := parser.New(m)

	out, ok := p.Parse(regex)
	if !ok {
		return nil, errors.New("invalid regular expression")
	}

	if m.errors != nil {
		return nil, m.errors
	}

	nfa := out.Result.Val.(*auto.NFA)

	return nfa, nil
}

//==================================================< MAPPERS >==================================================

// mappers implements the parser.Mappers interface.
type mappers struct {
	errors error
}

func (m *mappers) ToAnyChar(r comb.Result) (comb.Result, bool) {
	nfa := auto.NewNFA(0, auto.States{1})
	for _, r := range parser.Alphabet {
		nfa.Add(0, auto.Symbol(r), auto.States{1})
	}

	return comb.Result{
		Val: nfa,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToSingleChar(r comb.Result) (comb.Result, bool) {
	c := r.Val.(rune)

	return comb.Result{
		Val: runeToNFA(c),
		Pos: r.Pos,
		Bag: comb.Bag{
			bagKeyChars: []rune{c},
		},
	}, true
}

func (m *mappers) ToCharClass(r comb.Result) (comb.Result, bool) {
	class := r.Val.(string)

	var nfa *auto.NFA
	var chars []rune

	switch class {
	case `\d`:
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'})
	case `\D`:
		nfa, chars = runeRangesToNFA(true, [2]rune{'0', '9'})
	case `\s`:
		nfa, chars = runesToNFA(false, ' ', '\t', '\n', '\r', '\f')
	case `\S`:
		nfa, chars = runesToNFA(true, ' ', '\t', '\n', '\r', '\f')
	case `\w`:
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case `\W`:
		nfa, chars = runeRangesToNFA(true, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	default:
		return comb.Result{}, false
	}

	return comb.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: comb.Bag{
			bagKeyChars: chars,
		},
	}, true
}

func (m *mappers) ToASCIICharClass(r comb.Result) (comb.Result, bool) {
	class := r.Val.(string)

	var nfa *auto.NFA
	var chars []rune

	switch class {
	case "[:blank:]":
		nfa, chars = runesToNFA(false, ' ', '\t')
	case "[:space:]":
		nfa, chars = runesToNFA(false, ' ', '\t', '\n', '\r', '\f', '\v')
	case "[:digit:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'})
	case "[:xdigit:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'}, [2]rune{'A', 'F'}, [2]rune{'a', 'f'})
	case "[:upper:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'A', 'Z'})
	case "[:lower:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'a', 'z'})
	case "[:alpha:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:alnum:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:word:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case "[:ascii:]":
		nfa, chars = runeRangesToNFA(false, [2]rune{0x00, 0x7f})
	default:
		return comb.Result{}, false
	}

	return comb.Result{
		Val: nfa,
		Pos: r.Pos,
		Bag: comb.Bag{
			bagKeyChars: chars,
		},
	}, true
}

func (m *mappers) ToRepOp(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToUpperBound(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	var num *int
	if v, ok := r1.Val.(int); ok {
		num = &v
	}

	return comb.Result{
		Val: num,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRange(r comb.Result) (comb.Result, bool) {
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
		// The input syntax is correct while its semantic is incorrect
		// We continue parsing the rest of input to find more errors
		m.errors = multierror.Append(m.errors, fmt.Errorf("invalid repetition range {%d,%d}", low, *up))
	}

	return comb.Result{
		Val: tuple[int, *int]{
			p: low,
			q: up,
		},
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRepetition(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToQuantifier(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	// Check whether or not the lazy modifier is present
	_, lazy := r1.Val.(rune)

	return comb.Result{
		Val: tuple[any, bool]{
			p: r0.Val,
			q: lazy,
		},
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToCharInRange(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToCharRange(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r2, _ := r.Get(2)

	low, up := r0.Val.(rune), r2.Val.(rune)

	if low > up {
		// The input syntax is correct while its semantic is incorrect
		// We continue parsing the rest of input to find more errors
		m.errors = multierror.Append(m.errors, fmt.Errorf("invalid character range %s-%s", string(low), string(up)))
	}

	nfa, chars := runeRangesToNFA(false, [2]rune{low, up})

	return comb.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: comb.Bag{
			bagKeyChars: chars,
		},
	}, true
}

func (m *mappers) ToCharGroupItem(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToCharGroup(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	// Check whether or not the negation modifier is present
	_, neg := r1.Val.(rune)

	items := r2.Val.(comb.List)

	charMap := make([]bool, len(parser.Alphabet))
	for _, r := range items {
		if chars, ok := r.Bag[bagKeyChars].([]rune); ok {
			for _, c := range chars {
				charMap[c] = true
			}
		}
	}

	nfa := auto.NewNFA(0, auto.States{1})
	for i, marked := range charMap {
		if (!neg && marked) || (neg && !marked) {
			nfa.Add(0, auto.Symbol(rune(i)), auto.States{1})
		}
	}

	return comb.Result{
		Val: nfa,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToMatchItem(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToMatch(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	nfa := r0.Val.(*auto.NFA)
	var bag comb.Bag

	if t, ok := r1.Val.(tuple[any, bool]); ok {
		nfa = quantifyNFA(nfa, t.p)
		if lazy := t.q; lazy {
			bag = comb.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return comb.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: bag,
	}, true
}

func (m *mappers) ToGroup(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r3, _ := r.Get(3)

	nfa := r1.Val.(*auto.NFA)
	var bag comb.Bag

	if t, ok := r3.Val.(tuple[any, bool]); ok {
		nfa = quantifyNFA(nfa, t.p)
		if lazy := t.q; lazy {
			bag = comb.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return comb.Result{
		Val: nfa,
		Pos: r0.Pos,
		Bag: bag,
	}, true
}

func (m *mappers) ToAnchor(r comb.Result) (comb.Result, bool) {
	c := r.Val.(rune)

	var anchor Anchor
	switch c {
	case '$': // end-of-string
		anchor = EndOfString
	}

	return comb.Result{
		Val: anchor,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToSubexprItem(r comb.Result) (comb.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToSubexpr(r comb.Result) (comb.Result, bool) {
	items := r.Val.(comb.List)

	ns := []*auto.NFA{}
	for _, r := range items {
		// TODO: Anchor result value is not a node
		if n, ok := r.Val.(*auto.NFA); ok {
			ns = append(ns, n)
		}
	}

	nfa := concat(ns...)

	return comb.Result{
		Val: nfa,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToExpr(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	nfa := r0.Val.(*auto.NFA)

	if _, ok := r1.Val.(comb.List); ok {
		r11, _ := r1.Get(1)
		expr := r11.Val.(*auto.NFA)
		nfa = nfa.Union(expr)
	}

	return comb.Result{
		Val: nfa,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRegex(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	pos := r1.Pos
	var bag comb.Bag

	// Check whether or not the start-of-string is present
	if _, sos := r0.Val.(rune); sos {
		pos = r0.Pos
		bag = comb.Bag{
			BagKeyStartOfString: true,
		}
	}

	nfa := r1.Val.(*auto.NFA)

	return comb.Result{
		Val: nfa,
		Pos: pos,
		Bag: bag,
	}, true
}

//==================================================< HELPERS >==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

func runeToNFA(r rune) *auto.NFA {
	nfa := auto.NewNFA(0, auto.States{1})
	nfa.Add(0, auto.Symbol(r), auto.States{1})

	return nfa
}

func runesToNFA(neg bool, runes ...rune) (*auto.NFA, []rune) {
	nfa := auto.NewNFA(0, auto.States{1})
	chars := []rune{}

	if neg {
		for _, r := range parser.Alphabet {
			if !containsRune(r, runes) {
				nfa.Add(0, auto.Symbol(r), auto.States{1})
				chars = append(chars, r)
			}
		}
	} else {
		for _, r := range runes {
			nfa.Add(0, auto.Symbol(r), auto.States{1})
			chars = append(chars, r)
		}
	}

	return nfa, chars
}

func containsRune(r rune, runes []rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func runeRangesToNFA(neg bool, ranges ...[2]rune) (*auto.NFA, []rune) {
	nfa := auto.NewNFA(0, auto.States{1})
	chars := []rune{}

	if neg {
		for _, r := range parser.Alphabet {
			if !includesRune(r, ranges...) {
				nfa.Add(0, auto.Symbol(r), auto.States{1})
				chars = append(chars, r)
			}
		}
	} else {
		for _, g := range ranges {
			for r := g[0]; r <= g[1]; r++ {
				nfa.Add(0, auto.Symbol(r), auto.States{1})
				chars = append(chars, r)
			}
		}
	}

	return nfa, chars
}

func includesRune(r rune, ranges ...[2]rune) bool {
	for _, g := range ranges {
		if g[0] <= r && r <= g[1] {
			return true
		}
	}
	return false
}

func quantifyNFA(n *auto.NFA, q any) *auto.NFA {
	var nfa *auto.NFA

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
		ns := []*auto.NFA{}

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

		nfa = concat(ns...)
	}

	return nfa
}
