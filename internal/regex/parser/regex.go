package parser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/gardenbed/emerge/internal/regex/ast"
)

func Parse(in input) (ast.Node, error) {
	r := newRegex()

	out, ok := r.regex(in)
	if !ok {
		return nil, errors.New("syntax error")
	}

	// Check for errors
	if r.errors != nil {
		return nil, r.errors
	}

	root := out.Result.Val.(ast.Node)

	// Backfill Pos fields
	pos := 0
	setPositions(root, &pos)

	return root, nil
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

	char           parser
	charInGroup    parser
	charInMatch    parser
	digit          parser
	letter         parser
	num            parser
	letters        parser
	repOp          parser
	upperBound     parser
	range_         parser
	repetition     parser
	quantifier     parser
	charRange      parser
	charGroupItem  parser
	charGroup      parser
	charClass      parser
	asciiCharClass parser
	anyChar        parser
	matchItem      parser
	match          parser
	anchor         parser
	regex          parser
}

// newRegex creates a parser combinator for parsing regular expressions.
func newRegex() *regex {
	r := new(regex)

	r.char = expectRuneInRange(0x20, 0x7E).Map(r.toChar)                                     // char --> /* all valid characters */
	r.charInGroup = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(']')).Map(r.toChar)      // char --> /* all valid characters except ] */
	r.charInMatch = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(')', '|')).Map(r.toChar) // char --> /* all valid characters except ) and | */
	r.digit = expectRuneInRange('0', '9')                                                    // digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	r.letter = expectRuneInRange('A', 'Z').ALT(expectRuneInRange('a', 'z'))                  // letter --> "A" | ... | "Z" | "a" | ... | "z"
	r.num = r.digit.REP1().Map(r.toNum)                                                      // num --> digit+
	r.letters = r.letter.REP1().Map(r.toLetters)                                             // letters --> letter+

	// rep_op --> "?" | "*" | "+"
	r.repOp = expectRune('?').ALT(
		expectRune('*'),
		expectRune('+'),
	).Map(r.toRepOp)

	// upper_bound --> "," num?
	r.upperBound = expectRune(',').CONCAT(
		r.num.OPT(),
	).Map(r.toUpperBound)

	// range --> "{" num upper_bound? "}"
	r.range_ = expectRune('{').CONCAT(
		r.num,
		r.upperBound.OPT(),
		expectRune('}'),
	).Map(r.toRange)

	// repetition --> rep_op | range
	r.repetition = r.repOp.ALT(
		r.range_,
	).Map(r.toRepetition)

	// quantifier --> repetition lazy_modifier?
	r.quantifier = r.repetition.CONCAT(
		expectRune('?').OPT(),
	).Map(r.toQuantifier)

	// char_range --> char "-" char
	r.charRange = r.char.CONCAT(
		expectRune('-'),
		r.char,
	).Map(r.toCharRange)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	r.charClass = expectString(`\d`).ALT(
		expectString(`\D`),
		expectString(`\s`), expectString(`\S`),
		expectString(`\w`), expectString(`\W`),
	).Map(r.toCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	r.asciiCharClass = expectString("[:blank:]").ALT(
		expectString("[:space:]"),
		expectString("[:digit:]"), expectString("[:xdigit:]"),
		expectString("[:upper:]"), expectString("[:lower:]"),
		expectString("[:alpha:]"), expectString("[:alnum:]"),
		expectString("[:word:]"), expectString("[:ascii:]"),
	).Map(r.toASCIICharClass)

	// char_group_item -->  char_class | ascii_char_class | char_range | char /* excluding ] */
	r.charGroupItem = r.charClass.ALT(
		r.asciiCharClass,
		r.charRange,
		r.charInGroup,
	).Map(r.toCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	r.charGroup = expectRune('[').CONCAT(
		expectRune('^').OPT(),
		r.charGroupItem.REP1(),
		expectRune(']'),
	).Map(r.toCharGroup)

	// any_char --> "."
	r.anyChar = expectRune('.').Map(r.toAnyChar)

	// match_item --> any_char | char_class | ascii_char_class | char_group | char /* excluding | ) */
	r.matchItem = r.anyChar.ALT(
		r.charClass,
		r.asciiCharClass,
		r.charGroup,
		r.charInMatch,
	).Map(r.toMatchItem)

	// match --> match_item quantifier?
	r.match = r.matchItem.CONCAT(r.quantifier.OPT()).Map(r.toMatch)

	r.anchor = expectRune('$').Map(r.toAnchor) // anchor --> "$"

	// regex --> start_of_string? expr
	r.regex = expectRune('^').OPT().CONCAT(r.expr).Map(r.toRegex)

	return r
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (r *regex) group(in input) (output, bool) {
	return expectRune('(').CONCAT(
		r.expr,
		expectRune(')'),
		r.quantifier.OPT(),
	).Map(r.toGroup)(in)
}

// Recursive definition
// subexpr_item --> group | anchor | match
func (r *regex) subexprItem(in input) (output, bool) {
	return parser(r.group).ALT(r.anchor, r.match).Map(r.toSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (r *regex) subexpr(in input) (output, bool) {
	return parser(r.subexprItem).REP1().Map(r.toSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (r *regex) expr(in input) (output, bool) {
	return parser(r.subexpr).CONCAT(
		expectRune('|').CONCAT(r.expr).OPT(),
	).Map(r.toExpr)(in)
}

//================================================== MAPPERS ==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

func runeToChar(r rune) *ast.Char {
	// Pos will be set after the entire abstract syntax tree is constructed.
	return &ast.Char{
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

func includesRune(r rune, ranges ...[2]rune) bool {
	for _, g := range ranges {
		if g[0] <= r && r <= g[1] {
			return true
		}
	}
	return false
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
	case *ast.Concat:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *ast.Alt:
		for _, n := range v.Exprs {
			markAllChars(m, n)
		}

	case *ast.Star:
		markAllChars(m, v)

	case *ast.Char:
		m[v.Val] = true
	}
}

func cloneNode(n ast.Node) ast.Node {
	switch v := n.(type) {
	case *ast.Concat:
		concat := new(ast.Concat)
		for _, e := range v.Exprs {
			concat.Exprs = append(concat.Exprs, cloneNode(e))
		}
		return concat

	case *ast.Alt:
		alt := new(ast.Alt)
		for _, e := range v.Exprs {
			alt.Exprs = append(alt.Exprs, cloneNode(e))
		}
		return alt

	case *ast.Star:
		return &ast.Star{
			Expr: cloneNode(v.Expr),
		}

	case *ast.Empty:
		return new(ast.Empty)

	case *ast.Char:
		return &ast.Char{
			Val: v.Val,
			Pos: v.Pos,
		}

	default:
		return nil
	}
}

// TODO:
func quantifyNode(n ast.Node, t tuple[any, bool]) ast.Node {
	var node ast.Node

	switch rep := t.p.(type) {
	// Simple repetition
	case rune:
		switch rep {
		case '?':
			node = &ast.Alt{
				Exprs: []ast.Node{
					&ast.Empty{},
					cloneNode(n),
				},
			}

		case '*':
			node = &ast.Star{
				Expr: cloneNode(n),
			}

		case '+':
			node = &ast.Concat{
				Exprs: []ast.Node{
					cloneNode(n),
					&ast.Star{
						Expr: cloneNode(n),
					},
				},
			}
		}

	// Range repetition
	case tuple[int, *int]:
		low, up := rep.p, rep.q
		concat := new(ast.Concat)

		for i := 0; i < low; i++ {
			concat.Exprs = append(concat.Exprs, cloneNode(n))
		}

		if up == nil {
			concat.Exprs = append(concat.Exprs, &ast.Star{
				Expr: cloneNode(n),
			})
		} else {
			for i := 0; i < *up-low; i++ {
				concat.Exprs = append(concat.Exprs, &ast.Alt{
					Exprs: []ast.Node{
						&ast.Empty{},
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

func (r *regex) toChar(v any) (any, bool) {
	c := v.(rune)

	return runeToChar(c), true
}

func (r *regex) toNum(v any) (any, bool) {
	digits := v.(list)

	var num int
	for _, d := range digits {
		num = num*10 + int(d.Val.(rune)-'0')
	}

	return num, true
}

func (r *regex) toLetters(v any) (any, bool) {
	letters := v.(list)

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
	v1, _ := getAt(v, 1)

	var num *int
	if v, ok := v1.(int); ok {
		num = &v
	}

	// Passing the result up the parsing chain
	return num, true
}

func (r *regex) toRange(v any) (any, bool) {
	v1, _ := getAt(v, 1)
	v2, _ := getAt(v, 2)

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
	v0, _ := getAt(v, 0)
	v1, _ := getAt(v, 1)

	// Check whether or not the lazy modifier is present
	_, lazy := v1.(rune)

	return tuple[any, bool]{
		p: v0,
		q: lazy,
	}, true
}

func (r *regex) toCharRange(v any) (any, bool) {
	v0, _ := getAt(v, 0)
	v2, _ := getAt(v, 2)

	low := v0.(*ast.Char)
	up := v2.(*ast.Char)

	if low.Val > up.Val {
		r.errors = multierror.Append(r.errors,
			fmt.Errorf("invalid character range %s-%s",
				string(low.Val),
				string(up.Val),
			),
		)
	}

	return runeRangesToAlt(false, [2]rune{low.Val, up.Val}), true
}

func (r *regex) toCharGroupItem(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toCharGroup(v any) (any, bool) {
	v1, _ := getAt(v, 1)
	v2, _ := getAt(v, 2)

	// Check whether or not the negation modifier is present
	_, neg := v1.(rune)

	items := v2.(list)

	charMap := make([]bool, len(alphabet))
	for _, r := range items {
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

func (r *regex) toAnyChar(v any) (any, bool) {
	alt := new(ast.Alt)
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
	v0, _ := getAt(v, 0)
	v1, _ := getAt(v, 1)

	node := v0.(ast.Node)

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
	v1, _ := getAt(v, 1)
	v3, _ := getAt(v, 3)

	node := v1.(ast.Node)

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
	items := v.(list)

	concat := new(ast.Concat)
	for _, r := range items {
		// Anchor result value is not a node
		if n, ok := r.Val.(ast.Node); ok {
			concat.Exprs = append(concat.Exprs, n)
		}
	}

	return concat, true
}

func (r *regex) toExpr(v any) (any, bool) {
	v0, _ := getAt(v, 0)
	v1, _ := getAt(v, 1)

	subexpr := v0.(ast.Node)

	if _, ok := v1.(empty); ok {
		return subexpr, true
	}

	v11, _ := getAt(v1, 1)

	expr := v11.(ast.Node)

	alt := &ast.Alt{
		Exprs: []ast.Node{
			subexpr,
			expr,
		},
	}

	return alt, true
}

// TODO: handle start-of-string
func (r *regex) toRegex(v any) (any, bool) {
	// v0, _ := getAt(v, 0)
	v1, _ := getAt(v, 1)

	// Check whether or not the start-of-string is present
	// _, sos := v0.(rune)

	expr := v1.(ast.Node)

	return expr, true
}
