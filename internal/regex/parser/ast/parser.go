package ast

import (
	"errors"
	"fmt"

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

// mappers implements the parser.Mappers interface.
type mappers struct {
	errors error
}

func (m *mappers) ToAnyChar(r comb.Result) (comb.Result, bool) {
	alt := new(Alt)
	for _, r := range parser.Alphabet {
		alt.Exprs = append(alt.Exprs, runeToChar(r))
	}

	return comb.Result{
		Val: alt,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToSingleChar(r comb.Result) (comb.Result, bool) {
	c := r.Val.(rune)

	return comb.Result{
		Val: runeToChar(c),
		Pos: r.Pos,
		Bag: comb.Bag{
			bagKeyChars: []rune{c},
		},
	}, true
}

func (m *mappers) ToCharClass(r comb.Result) (comb.Result, bool) {
	class := r.Val.(string)

	var node *Alt
	var chars []rune

	switch class {
	case `\d`:
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'})
	case `\D`:
		node, chars = runeRangesToAlt(true, [2]rune{'0', '9'})
	case `\s`:
		node, chars = runesToAlt(false, ' ', '\t', '\n', '\r', '\f')
	case `\S`:
		node, chars = runesToAlt(true, ' ', '\t', '\n', '\r', '\f')
	case `\w`:
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case `\W`:
		node, chars = runeRangesToAlt(true, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	default:
		return comb.Result{}, false
	}

	return comb.Result{
		Val: node,
		Pos: r.Pos,
		Bag: comb.Bag{
			bagKeyChars: chars,
		},
	}, true
}

func (m *mappers) ToASCIICharClass(r comb.Result) (comb.Result, bool) {
	class := r.Val.(string)

	var node *Alt
	var chars []rune

	switch class {
	case "[:blank:]":
		node, chars = runesToAlt(false, ' ', '\t')
	case "[:space:]":
		node, chars = runesToAlt(false, ' ', '\t', '\n', '\r', '\f', '\v')
	case "[:digit:]":
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'})
	case "[:xdigit:]":
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'F'}, [2]rune{'a', 'f'})
	case "[:upper:]":
		node, chars = runeRangesToAlt(false, [2]rune{'A', 'Z'})
	case "[:lower:]":
		node, chars = runeRangesToAlt(false, [2]rune{'a', 'z'})
	case "[:alpha:]":
		node, chars = runeRangesToAlt(false, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:alnum:]":
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:word:]":
		node, chars = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case "[:ascii:]":
		node, chars = runeRangesToAlt(false, [2]rune{0x00, 0x7f})
	default:
		return comb.Result{}, false
	}

	return comb.Result{
		Val: node,
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
		m.errors = errors.Join(m.errors, fmt.Errorf("invalid repetition range {%d,%d}", low, *up))
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
		m.errors = errors.Join(m.errors, fmt.Errorf("invalid character range %s-%s", string(low), string(up)))
	}

	node, chars := runeRangesToAlt(false, [2]rune{low, up})

	return comb.Result{
		Val: node,
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

	alt := new(Alt)
	for i, marked := range charMap {
		if (!neg && marked) || (neg && !marked) {
			alt.Exprs = append(alt.Exprs, &Char{
				Val: rune(i),
			})
		}
	}

	return comb.Result{
		Val: alt,
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

	node := r0.Val.(Node)
	var bag comb.Bag

	if t, ok := r1.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, t.p)
		if lazy := t.q; lazy {
			bag = comb.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return comb.Result{
		Val: node,
		Pos: r0.Pos,
		Bag: bag,
	}, true
}

func (m *mappers) ToGroup(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r3, _ := r.Get(3)

	node := r1.Val.(Node)
	var bag comb.Bag

	if t, ok := r3.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, t.p)
		if lazy := t.q; lazy {
			bag = comb.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return comb.Result{
		Val: node,
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

	concat := new(Concat)
	for _, r := range items {
		// TODO: Anchor result value is not a node
		if n, ok := r.Val.(Node); ok {
			concat.Exprs = append(concat.Exprs, n)
		}
	}

	return comb.Result{
		Val: concat,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToExpr(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	node := r0.Val.(Node)

	if _, ok := r1.Val.(comb.List); ok {
		r11, _ := r1.Get(1)
		expr := r11.Val.(Node)
		node = &Alt{
			Exprs: []Node{node, expr},
		}
	}

	return comb.Result{
		Val: node,
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

	expr := r1.Val.(Node)

	return comb.Result{
		Val: expr,
		Pos: pos,
		Bag: bag,
	}, true
}

//==================================================< HELPERS >==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

func runeToChar(r rune) *Char {
	return &Char{
		Val: r,
		// Pos will be set after the entire abstract syntax tree is constructed.
	}
}

func runesToAlt(neg bool, runes ...rune) (*Alt, []rune) {
	alt := new(Alt)
	chars := []rune{}

	if neg {
		for _, r := range parser.Alphabet {
			if !containsRune(r, runes) {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
				chars = append(chars, r)
			}
		}
	} else {
		for _, r := range runes {
			alt.Exprs = append(alt.Exprs, runeToChar(r))
			chars = append(chars, r)
		}
	}

	return alt, chars
}

func containsRune(r rune, runes []rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func runeRangesToAlt(neg bool, ranges ...[2]rune) (*Alt, []rune) {
	alt := new(Alt)
	chars := []rune{}

	if neg {
		for _, r := range parser.Alphabet {
			if !includesRune(r, ranges...) {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
				chars = append(chars, r)
			}
		}
	} else {
		for _, g := range ranges {
			for r := g[0]; r <= g[1]; r++ {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
				chars = append(chars, r)
			}
		}
	}

	return alt, chars
}

func includesRune(r rune, ranges ...[2]rune) bool {
	for _, g := range ranges {
		if g[0] <= r && r <= g[1] {
			return true
		}
	}
	return false
}

func quantifyNode(n Node, q any) Node {
	var node Node

	switch rep := q.(type) {
	// Simple repetition
	case rune:
		switch rep {
		case '?':
			node = &Alt{
				Exprs: []Node{
					&Empty{},
					cloneNode(n),
				},
			}

		case '*':
			node = &Star{
				Expr: cloneNode(n),
			}

		case '+':
			node = &Concat{
				Exprs: []Node{
					cloneNode(n),
					&Star{
						Expr: cloneNode(n),
					},
				},
			}
		}

	// Range repetition
	case tuple[int, *int]:
		low, up := rep.p, rep.q
		concat := new(Concat)

		for i := 0; i < low; i++ {
			concat.Exprs = append(concat.Exprs, cloneNode(n))
		}

		if up == nil {
			concat.Exprs = append(concat.Exprs, &Star{
				Expr: cloneNode(n),
			})
		} else {
			for i := 0; i < *up-low; i++ {
				concat.Exprs = append(concat.Exprs, &Alt{
					Exprs: []Node{
						&Empty{},
						cloneNode(n),
					},
				})
			}
		}

		node = concat
	}

	return node
}

// Clonig is required since each instance of Char will have a distinct value for Pos field.
func cloneNode(n Node) Node {
	switch v := n.(type) {
	case *Concat:
		concat := new(Concat)
		for _, e := range v.Exprs {
			concat.Exprs = append(concat.Exprs, cloneNode(e))
		}
		return concat

	case *Alt:
		alt := new(Alt)
		for _, e := range v.Exprs {
			alt.Exprs = append(alt.Exprs, cloneNode(e))
		}
		return alt

	case *Star:
		return &Star{
			Expr: cloneNode(v.Expr),
		}

	case *Empty:
		return new(Empty)

	case *Char:
		return &Char{
			Val: v.Val,
			Pos: v.Pos,
		}

	default:
		return nil
	}
}
