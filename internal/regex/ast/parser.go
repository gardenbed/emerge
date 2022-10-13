package ast

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	comb "github.com/gardenbed/emerge/internal/combinator"
)

func Parse(in comb.Input) (Node, error) {
	r := newRegex()

	out, ok := r.regex(in)
	if !ok {
		return nil, errors.New("invalid regular expression")
	}

	// Check for errors
	if r.errors != nil {
		return nil, r.errors
	}

	root := out.Result.Val.(Node)

	// Backfill Pos fields
	pos := 0
	setPositions(root, &pos)

	return root, nil
}

// setPositions backfills Pos field for all characters in an abstract syntaxt tree from left to right.
func setPositions(n Node, pos *int) {
	switch v := n.(type) {
	case *Concat:
		for _, e := range v.Exprs {
			setPositions(e, pos)
		}

	case *Alt:
		for _, e := range v.Exprs {
			setPositions(e, pos)
		}

	case *Star:
		setPositions(v.Expr, pos)

	case *Char:
		*pos++
		v.Pos = *pos
	}
}

//================================================== PARSER COMBINATOR ==================================================

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

// regex is a parser combinator for parsing regular expressions.
type regex struct {
	errors error

	char           comb.Parser
	unescapedChar  comb.Parser
	escapedChar    comb.Parser
	digit          comb.Parser
	letter         comb.Parser
	num            comb.Parser
	letters        comb.Parser
	repOp          comb.Parser
	upperBound     comb.Parser
	range_         comb.Parser
	repetition     comb.Parser
	quantifier     comb.Parser
	charRange      comb.Parser
	charGroupItem  comb.Parser
	charGroup      comb.Parser
	charClass      comb.Parser
	asciiCharClass comb.Parser
	anyChar        comb.Parser
	matchItem      comb.Parser
	match          comb.Parser
	anchor         comb.Parser
	regex          comb.Parser
}

// newRegex creates a parser combinator for parsing regular expressions.
func newRegex() *regex {
	r := new(regex)

	// char --> /* all valid characters */
	r.char = comb.ExpectRuneInRange(0x20, 0x7E)
	// all characters excluding the escaped ones
	r.unescapedChar = r.char.Bind(comb.ExcludeRunes('\\', '/', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}')).Map(r.toUnescapedChar)
	// escaped_char --> "\" ( "\" | "/" | "|" | "." | "?" | "*" | "+" | "(" | ")" | "[" | "]" | "{" | "}" )
	r.escapedChar = comb.ExpectRune('\\').CONCAT(comb.ExpectRuneIn('\\', '/', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}')).Map(r.toEscapedChar)

	r.digit = comb.ExpectRuneInRange('0', '9')                                        // digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	r.letter = comb.ExpectRuneInRange('A', 'Z').ALT(comb.ExpectRuneInRange('a', 'z')) // letter --> "A" | ... | "Z" | "a" | ... | "z"

	r.num = r.digit.REP1().Map(r.toNum)          // num --> digit+
	r.letters = r.letter.REP1().Map(r.toLetters) // letters --> letter+

	// rep_op --> "?" | "*" | "+"
	r.repOp = comb.ExpectRune('?').ALT(
		comb.ExpectRune('*'),
		comb.ExpectRune('+'),
	).Map(r.toRepOp)

	// upper_bound --> "," num?
	r.upperBound = comb.ExpectRune(',').CONCAT(
		r.num.OPT(),
	).Map(r.toUpperBound)

	// range --> "{" num upper_bound? "}"
	r.range_ = comb.ExpectRune('{').CONCAT(
		r.num,
		r.upperBound.OPT(),
		comb.ExpectRune('}'),
	).Map(r.toRange)

	// repetition --> rep_op | range
	r.repetition = r.repOp.ALT(
		r.range_,
	).Map(r.toRepetition)

	// quantifier --> repetition lazy_modifier?
	r.quantifier = r.repetition.CONCAT(
		comb.ExpectRune('?').OPT(),
	).Map(r.toQuantifier)

	// char_range --> char "-" char
	r.charRange = r.char.CONCAT(
		comb.ExpectRune('-'),
		r.char,
	).Map(r.toCharRange)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	r.charClass = comb.ExpectString(`\d`).ALT(
		comb.ExpectString(`\D`),
		comb.ExpectString(`\s`), comb.ExpectString(`\S`),
		comb.ExpectString(`\w`), comb.ExpectString(`\W`),
	).Map(r.toCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	r.asciiCharClass = comb.ExpectString("[:blank:]").ALT(
		comb.ExpectString("[:space:]"),
		comb.ExpectString("[:digit:]"), comb.ExpectString("[:xdigit:]"),
		comb.ExpectString("[:upper:]"), comb.ExpectString("[:lower:]"),
		comb.ExpectString("[:alpha:]"), comb.ExpectString("[:alnum:]"),
		comb.ExpectString("[:word:]"), comb.ExpectString("[:ascii:]"),
	).Map(r.toASCIICharClass)

	// char_group_item --> char_class | ascii_char_class | char_range | escaped_char | unescaped_char
	r.charGroupItem = r.charClass.ALT(
		r.asciiCharClass,
		r.charRange,
		r.escapedChar,
		r.unescapedChar,
	).Map(r.toCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	r.charGroup = comb.ExpectRune('[').CONCAT(
		comb.ExpectRune('^').OPT(),
		r.charGroupItem.REP1(),
		comb.ExpectRune(']'),
	).Map(r.toCharGroup)

	// any_char --> "."
	r.anyChar = comb.ExpectRune('.').Map(r.toAnyChar)

	// match_item --> any_char | unescaped_char | escaped_char | char_class | ascii_char_class | char_group
	r.matchItem = r.anyChar.ALT(
		r.unescapedChar,
		r.escapedChar,
		r.charClass,
		r.asciiCharClass,
		r.charGroup,
	).Map(r.toMatchItem)

	// match --> match_item quantifier?
	r.match = r.matchItem.CONCAT(r.quantifier.OPT()).Map(r.toMatch)

	r.anchor = comb.ExpectRune('$').Map(r.toAnchor) // anchor --> "$"

	// regex --> start_of_string? expr
	r.regex = comb.ExpectRune('^').OPT().CONCAT(r.expr).Map(r.toRegex)

	return r
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (r *regex) group(in comb.Input) (comb.Output, bool) {
	return comb.ExpectRune('(').CONCAT(
		r.expr,
		comb.ExpectRune(')'),
		r.quantifier.OPT(),
	).Map(r.toGroup)(in)
}

// Recursive definition
// subexpr_item --> anchor | group | match
func (r *regex) subexprItem(in comb.Input) (comb.Output, bool) {
	return comb.Parser(r.anchor).ALT(r.group, r.match).Map(r.toSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (r *regex) subexpr(in comb.Input) (comb.Output, bool) {
	return comb.Parser(r.subexprItem).REP1().Map(r.toSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (r *regex) expr(in comb.Input) (comb.Output, bool) {
	return comb.Parser(r.subexpr).CONCAT(
		comb.ExpectRune('|').CONCAT(r.expr).OPT(),
	).Map(r.toExpr)(in)
}

//================================================== MAPPERS ==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

func runeToChar(r rune) *Char {
	// Pos will be set after the entire abstract syntax tree is constructed.
	return &Char{
		Val: r,
	}
}

func containsRune(r rune, runes []rune) bool {
	for _, v := range runes {
		if v == r {
			return true
		}
	}
	return false
}

func runesToAlt(neg bool, runes ...rune) *Alt {
	alt := new(Alt)

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

func includesRune(r rune, ranges ...[2]rune) bool {
	for _, g := range ranges {
		if g[0] <= r && r <= g[1] {
			return true
		}
	}
	return false
}

func runeRangesToAlt(neg bool, ranges ...[2]rune) *Alt {
	alt := new(Alt)

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

func markAllChars(m []bool, n Node) {
	switch v := n.(type) {
	case *Concat:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *Alt:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *Star:
		markAllChars(m, v)

	case *Char:
		m[v.Val] = true
	}
}

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

// TODO:
func quantifyNode(n Node, t tuple[any, bool]) Node {
	var node Node

	switch rep := t.p.(type) {
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

	// lazy := q.u

	return node
}

func (r *regex) toUnescapedChar(v any) (any, bool) {
	c := v.(rune)

	return runeToChar(c), true
}

func (r *regex) toEscapedChar(v any) (any, bool) {
	v1, _ := comb.GetAt(v, 1)

	c := v1.(rune)

	return runeToChar(c), true
}

func (r *regex) toNum(v any) (any, bool) {
	digits := v.(comb.List)

	var num int
	for _, d := range digits {
		num = num*10 + int(d.Val.(rune)-'0')
	}

	return num, true
}

func (r *regex) toLetters(v any) (any, bool) {
	letters := v.(comb.List)

	var s string
	for _, l := range letters {
		s += string(l.Val.(rune))
	}

	return s, true
}

func (r *regex) toRepOp(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toUpperBound(v any) (any, bool) {
	v1, _ := comb.GetAt(v, 1)

	var num *int
	if v, ok := v1.(int); ok {
		num = &v
	}

	// Passing the result up the parsing chain
	return num, true
}

func (r *regex) toRange(v any) (any, bool) {
	v1, _ := comb.GetAt(v, 1)
	v2, _ := comb.GetAt(v, 2)

	// The upper bound is same as the lower bound if no upper bound is specified (default)
	low := v1.(int)
	up := &low

	// If an upper bound is specified, it can be either bounded or unbounded
	if v, ok := v2.(*int); ok {
		up = v
	}

	if up != nil && low > *up {
		r.errors = multierror.Append(r.errors,
			fmt.Errorf("invalid repetition range {%d,%d}", low, *up),
		)
	}

	return tuple[int, *int]{
		p: low,
		q: up,
	}, true
}

func (r *regex) toRepetition(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toQuantifier(v any) (any, bool) {
	v0, _ := comb.GetAt(v, 0)
	v1, _ := comb.GetAt(v, 1)

	// Check whether or not the lazy modifier is present
	_, lazy := v1.(rune)

	return tuple[any, bool]{
		p: v0,
		q: lazy,
	}, true
}

func (r *regex) toCharRange(v any) (any, bool) {
	v0, _ := comb.GetAt(v, 0)
	v2, _ := comb.GetAt(v, 2)

	low, up := v0.(rune), v2.(rune)

	if low > up {
		r.errors = multierror.Append(r.errors,
			fmt.Errorf("invalid character range %s-%s",
				string(low),
				string(up),
			),
		)
	}

	return runeRangesToAlt(false, [2]rune{low, up}), true
}

func (r *regex) toCharGroupItem(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toCharGroup(v any) (any, bool) {
	v1, _ := comb.GetAt(v, 1)
	v2, _ := comb.GetAt(v, 2)

	// Check whether or not the negation modifier is present
	_, neg := v1.(rune)

	items := v2.(comb.List)

	charMap := make([]bool, len(alphabet))
	for _, r := range items {
		if n, ok := r.Val.(Node); ok {
			markAllChars(charMap, n)
		}
	}

	var exprs []Node
	for i, exist := range charMap {
		if (!neg && exist) || (neg && !exist) {
			exprs = append(exprs, &Char{
				Val: rune(i),
			})
		}
	}

	return &Alt{
		Exprs: exprs,
	}, true
}

func (r *regex) toASCIICharClass(v any) (any, bool) {
	class := v.(string)

	switch class {
	case "[:blank:]":
		return runesToAlt(false, ' ', '\t'), true
	case "[:space:]":
		return runesToAlt(false, ' ', '\t', '\n', '\r', '\f', '\v'), true
	case "[:digit:]":
		return runeRangesToAlt(false, [2]rune{'0', '9'}), true
	case "[:xdigit:]":
		return runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'F'}, [2]rune{'a', 'f'}), true
	case "[:upper:]":
		return runeRangesToAlt(false, [2]rune{'A', 'Z'}), true
	case "[:lower:]":
		return runeRangesToAlt(false, [2]rune{'a', 'z'}), true
	case "[:alpha:]":
		return runeRangesToAlt(false, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'}), true
	case "[:alnum:]":
		return runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'a', 'z'}), true
	case "[:word:]":
		return runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'}), true
	case "[:ascii:]":
		return runeRangesToAlt(false, [2]rune{0x00, 0x7f}), true
	default:
		return nil, false
	}
}

func (r *regex) toCharClass(v any) (any, bool) {
	class := v.(string)

	switch class {
	case `\d`:
		return runeRangesToAlt(false, [2]rune{'0', '9'}), true
	case `\D`:
		return runeRangesToAlt(true, [2]rune{'0', '9'}), true
	case `\s`:
		return runesToAlt(false, ' ', '\t', '\n', '\r', '\f'), true
	case `\S`:
		return runesToAlt(true, ' ', '\t', '\n', '\r', '\f'), true
	case `\w`:
		return runeRangesToAlt(false, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'}), true
	case `\W`:
		return runeRangesToAlt(true, [2]rune{'0', '9'}, [2]rune{'A', 'Z'}, [2]rune{'_', '_'}, [2]rune{'a', 'z'}), true
	default:
		return nil, false
	}
}

func (r *regex) toAnyChar(v any) (any, bool) {
	alt := new(Alt)
	for _, r := range alphabet {
		alt.Exprs = append(alt.Exprs, runeToChar(r))
	}

	return alt, true
}

func (r *regex) toMatchItem(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toMatch(v any) (any, bool) {
	v0, _ := comb.GetAt(v, 0)
	v1, _ := comb.GetAt(v, 1)

	node := v0.(Node)

	if q, ok := v1.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	return node, true
}

// TODO: handle of end-of-string
func (r *regex) toAnchor(v any) (any, bool) {
	a := v.(rune)

	// Check whether or not the anchor is end-of-string
	// eos := a == '$'

	return a, true
}

func (r *regex) toGroup(v any) (any, bool) {
	v1, _ := comb.GetAt(v, 1)
	v3, _ := comb.GetAt(v, 3)

	node := v1.(Node)

	if q, ok := v3.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	return node, true
}

func (r *regex) toSubexprItem(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toSubexpr(v any) (any, bool) {
	items := v.(comb.List)

	concat := new(Concat)
	for _, r := range items {
		// Anchor result value is not a node
		if n, ok := r.Val.(Node); ok {
			concat.Exprs = append(concat.Exprs, n)
		}
	}

	return concat, true
}

func (r *regex) toExpr(v any) (any, bool) {
	v0, _ := comb.GetAt(v, 0)
	v1, _ := comb.GetAt(v, 1)

	subexpr := v0.(Node)

	if _, ok := v1.(comb.Empty); ok {
		return subexpr, true
	}

	v11, _ := comb.GetAt(v1, 1)

	expr := v11.(Node)

	alt := &Alt{
		Exprs: []Node{
			subexpr,
			expr,
		},
	}

	return alt, true
}

// TODO: handle start-of-string
func (r *regex) toRegex(v any) (any, bool) {
	// v0, _ := getAt(v, 0)
	v1, _ := comb.GetAt(v, 1)

	// Check whether or not the start-of-string is present
	// _, sos := v0.(rune)

	expr := v1.(Node)

	return expr, true
}
