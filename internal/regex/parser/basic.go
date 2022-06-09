package parser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	ast "github.com/gardenbed/emerge/internal/regex/ast/basic"
)

// All characters from 0x00 to 0x7f
var alphabet = []rune{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	':', ';', '<', '=', '>', '?', '@',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'[', '\\', ']', '^', '_', '`',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'{', '|', '}', '~', 0x7f,
}

type tuple[T, U any] struct {
	t T
	u U
}

func ParseBasic(in input) (ast.Node, error) {
	c := new(basicConverters)
	r := NewRegex(c)

	out, ok := r.regex(in)
	if !ok {
		return nil, errors.New("syntax error")
	}

	// Check for errors
	if c.errors != nil {
		return nil, c.errors
	}

	root := out.Result.Val.(ast.Node)

	// Backfill Pos fields
	pos := 0
	setPositions(root, &pos)

	return root, nil
}

// basicConverters implements the converters interface for the basic regex AST.
type basicConverters struct {
	sos    bool
	eos    bool
	errors error
	symTab []ast.Node
}

func containsRune(r rune, runes []rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func includesRune(r rune, ranges ...[2]rune) bool {
	for _, g := range ranges {
		if g[0] <= r && r <= g[1] {
			return true
		}
	}
	return false
}

func runeToChar(r rune) *ast.Char {
	// Pos will be set after the entire abstract syntax tree is constructed.
	return &ast.Char{
		Val: r,
	}
}

func runesToAlt(neg bool, runes ...rune) *ast.Alt {
	alt := new(ast.Alt)

	if neg {
		for _, r := range alphabet {
			if !containsRune(r, runes) {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
			}
		}
	} else {
		for _, r := range runes {
			alt.Exprs = append(alt.Exprs, runeToChar(r))
		}
	}

	return alt
}

func runeRangesToAlt(neg bool, ranges ...[2]rune) *ast.Alt {
	alt := new(ast.Alt)

	if neg {
		for _, r := range alphabet {
			if !includesRune(r, ranges...) {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
			}
		}
	} else {
		for _, g := range ranges {
			for r := g[0]; r <= g[1]; r++ {
				alt.Exprs = append(alt.Exprs, runeToChar(r))
			}
		}
	}

	return alt
}

func markAllChars(m []bool, n ast.Node) {
	switch v := n.(type) {
	case *ast.Alt:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *ast.Concat:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *ast.Star:
		markAllChars(m, v)

	case *ast.Char:
		m[v.Val] = true
	}
}

func quantifyNode(n ast.Node, q tuple[any, bool]) ast.Node {
	var node ast.Node

	switch rep := q.t.(type) {
	// Simple repetition
	case rune:
		switch rep {
		case '?':
			node = &ast.Alt{
				Exprs: []ast.Node{
					n,
					&ast.Empty{},
				},
			}

		case '*':
			node = &ast.Star{
				Expr: n,
			}

		case '+':
			node = &ast.Concat{
				Exprs: []ast.Node{
					n,
					&ast.Star{
						Expr: n,
					},
				},
			}
		}

	// Range repetition
	case tuple[int, *int]:
		low, up := rep.t, rep.u
		concat := new(ast.Concat)

		for i := 0; i < low; i++ {
			concat.Exprs = append(concat.Exprs, n)
		}

		if up == nil {
			concat.Exprs = append(concat.Exprs, &ast.Star{
				Expr: n,
			})
		} else {
			for i := 0; i < *up-low; i++ {
				concat.Exprs = append(concat.Exprs, &ast.Alt{
					Exprs: []ast.Node{
						n,
						&ast.Empty{},
					},
				})
			}
		}

		node = concat
	}

	// TODO:
	// lazy := q.u

	return node
}

// setPositions backfills Pos field for all characters in an abstract syntaxt tree from left to right.
func setPositions(n ast.Node, pos *int) {
	switch v := n.(type) {
	case *ast.Concat:
		for _, e := range v.Exprs {
			setPositions(e, pos)
		}

	case *ast.Alt:
		for _, e := range v.Exprs {
			setPositions(e, pos)
		}

	case *ast.Star:
		setPositions(v.Expr, pos)

	case *ast.Char:
		*pos++
		v.Pos = *pos
	}
}

func (c *basicConverters) ToChar(res result) (any, bool) {
	r, ok := res.Val.(rune)
	if !ok {
		return nil, false
	}

	return runeToChar(r), true
}

func (c *basicConverters) ToNum(res result) (any, bool) {
	digits, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	var num int
	for _, d := range digits {
		num = num*10 + int(d.Val.(rune)-'0')
	}

	return num, true
}

func (c *basicConverters) ToLetters(res result) (any, bool) {
	letters, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	var s string
	for _, l := range letters {
		s += string(l.Val.(rune))
	}

	return s, true
}

func (c *basicConverters) ToRepOp(res result) (any, bool) {
	// Passing the result up the parsing chain
	return res.Val, true
}

func (c *basicConverters) ToUpperBound(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	var num *int
	if v, ok := r1.Val.(int); ok {
		num = &v
	}

	// Passing the result up the parsing chain
	return num, true
}

func (c *basicConverters) ToRange(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	low, ok := r1.Val.(int)
	if !ok {
		return nil, false
	}

	// The upper bound is same as the lower bound if no upper bound is specified (default)
	up := &low

	// If an upper bound is specified, it can be either bounded or unbounded
	if v, ok := r2.Val.(*int); ok {
		up = v
	}

	if up != nil && low > *up {
		c.errors = multierror.Append(c.errors,
			fmt.Errorf("invalid repetition range {%d,%d}", low, *up),
		)
	}

	return tuple[int, *int]{
		t: low,
		u: up,
	}, true
}

func (c *basicConverters) ToRepetition(res result) (any, bool) {
	// Passing the result up the parsing chain
	return res.Val, true
}

func (c *basicConverters) ToQuantifier(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	// Check whether or not the lazy modifier is present
	_, lazy := r1.Val.(rune)

	return tuple[any, bool]{
		t: r0.Val,
		u: lazy,
	}, true
}

func (c *basicConverters) ToCharRange(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	low, ok := r0.Val.(*ast.Char)
	if !ok {
		return nil, false
	}

	up, ok := r2.Val.(*ast.Char)
	if !ok {
		return nil, false
	}

	if low.Val > up.Val {
		c.errors = multierror.Append(c.errors,
			fmt.Errorf("invalid character range %s-%s",
				string(low.Val),
				string(up.Val),
			),
		)
	}

	return runeRangesToAlt(false, [2]rune{low.Val, up.Val}), true
}

func (c *basicConverters) ToCharGroupItem(res result) (any, bool) {
	// Passing the result up the parsing chain
	return res.Val, true
}

func (c *basicConverters) ToCharGroup(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	// Check whether or not the negation modifier is present
	_, neg := r1.Val.(rune)

	l, ok := r2.Val.(list)
	if !ok {
		return nil, false
	}

	charMap := make([]bool, len(alphabet))
	for _, r := range l {
		if n, ok := r.Val.(ast.Node); ok {
			markAllChars(charMap, n)
		}
	}

	var exprs []ast.Node
	for i, exist := range charMap {
		if (!neg && exist) || (neg && !exist) {
			exprs = append(exprs, &ast.Char{
				Val: rune(i),
			})
		}
	}

	return &ast.Alt{
		Exprs: exprs,
	}, true
}

func (c *basicConverters) ToCharClass(res result) (any, bool) {
	val, ok := res.Val.(string)
	if !ok {
		return nil, false
	}

	var alt *ast.Alt

	switch val {
	case `\d`:
		alt = runeRangesToAlt(false, [2]rune{'0', '9'})
	case `\D`:
		alt = runeRangesToAlt(true, [2]rune{'0', '9'})
	case `\s`:
		alt = runesToAlt(false, ' ', '\t', '\n', '\r', '\f')
	case `\S`:
		alt = runesToAlt(true, ' ', '\t', '\n', '\r', '\f')
	case `\w`:
		alt = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case `\W`:
		alt = runeRangesToAlt(true, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	default:
		return nil, false
	}

	return alt, true
}

func (c *basicConverters) ToASCIICharClass(res result) (any, bool) {
	val, ok := res.Val.(string)
	if !ok {
		return nil, false
	}

	var alt *ast.Alt

	switch val {
	case "[:blank:]":
		alt = runesToAlt(false, ' ', '\t')
	case "[:space:]":
		alt = runesToAlt(false, ' ', '\t', '\n', '\r', '\f', '\v')
	case "[:digit:]":
		alt = runeRangesToAlt(false, [2]rune{'0', '9'})
	case "[:xdigit:]":
		alt = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'F'}, [2]rune{'a', 'f'})
	case "[:upper:]":
		alt = runeRangesToAlt(false, [2]rune{'A', 'Z'})
	case "[:lower:]":
		alt = runeRangesToAlt(false, [2]rune{'a', 'z'})
	case "[:alpha:]":
		alt = runeRangesToAlt(false, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:alnum:]":
		alt = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'})
	case "[:word:]":
		alt = runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'})
	case "[:ascii:]":
		alt = runeRangesToAlt(false, [2]rune{0x00, 0x7f})
	default:
		return nil, false
	}

	return alt, true
}

func (c *basicConverters) ToAnyChar(res result) (any, bool) {
	if _, ok := res.Val.(rune); !ok {
		return nil, false
	}

	alt := new(ast.Alt)
	for _, r := range alphabet {
		alt.Exprs = append(alt.Exprs, runeToChar(r))
	}

	return alt, true
}

func (c *basicConverters) ToMatchItem(res result) (any, bool) {
	// Passing the result up the parsing chain
	return res.Val, true
}

func (c *basicConverters) ToMatch(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	node, ok := r0.Val.(ast.Node)
	if !ok {
		return nil, false
	}

	if q, ok := r1.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	return node, true
}

func (c *basicConverters) ToBackref(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	num, ok := r1.Val.(int)
	if !ok {
		return nil, false
	}

	var node ast.Node

	// Look up the symbol table
	if i := num - 1; 0 <= i && i < len(c.symTab) {
		node = c.symTab[i]
	} else {
		c.errors = multierror.Append(c.errors, fmt.Errorf("invalid back reference \\%d", num))
	}

	// Backref is successfully parsed, but it is invalid
	// If we return false, backref will be parsed by other parsers (match)
	return node, true
}

func (c *basicConverters) ToAnchor(res result) (any, bool) {
	a, ok := res.Val.(rune)
	if !ok {
		return nil, false
	}

	// Check whether or not the anchor is end-of-string
	c.eos = a == '$'

	return a, true
}

func (c *basicConverters) ToGroup(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r3, ok := res.Get(3)
	if !ok {
		return nil, false
	}

	node, ok := r1.Val.(ast.Node)
	if !ok {
		return nil, false
	}

	if q, ok := r3.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	// Adding the group to the symbol table
	c.symTab = append(c.symTab, node)

	return node, true
}

func (c *basicConverters) ToSubexprItem(res result) (any, bool) {
	// Passing the result up the parsing chain
	return res.Val, true
}

func (c *basicConverters) ToSubexpr(res result) (any, bool) {
	l, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	concat := new(ast.Concat)
	for _, r := range l {
		if n, ok := r.Val.(ast.Node); ok {
			concat.Exprs = append(concat.Exprs, n)
		}
	}

	return concat, true
}

func (c *basicConverters) ToExpr(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	subexpr, ok := r0.Val.(ast.Node)
	if !ok {
		return nil, false
	}

	if _, ok := r1.Val.(empty); ok {
		return subexpr, true
	}

	r11, ok := r1.Get(1)
	if !ok {
		return nil, false
	}

	expr, ok := r11.Val.(ast.Node)
	if !ok {
		return nil, false
	}

	alt := &ast.Alt{
		Exprs: []ast.Node{
			subexpr,
			expr,
		},
	}

	return alt, true
}

func (c *basicConverters) ToRegex(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	// Check whether or not the start-of-string is present
	_, c.sos = r0.Val.(rune)

	expr, ok := r1.Val.(ast.Node)
	if !ok {
		return nil, false
	}

	return expr, true
}
