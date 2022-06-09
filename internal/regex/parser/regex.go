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
	sos    bool
	eos    bool
	errors error
	symTab []ast.Node

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
	backref        parser
	anchor         parser
	regex          parser
}

// newRegex creates a parser combinator for parsing regular expressions.
func newRegex() *regex {
	r := new(regex)

	r.char = expectRuneInRange(0x20, 0x7E).Convert(r.toChar)                                     // char --> /* all valid characters */
	r.charInGroup = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(']')).Convert(r.toChar)      // char --> /* all valid characters except ] */
	r.charInMatch = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(')', '|')).Convert(r.toChar) // char --> /* all valid characters except ) and | */
	r.digit = expectRuneInRange('0', '9')                                                        // digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	r.letter = expectRuneInRange('A', 'Z').ALT(expectRuneInRange('a', 'z'))                      // letter --> "A" | ... | "Z" | "a" | ... | "z"
	r.num = r.digit.REP1().Convert(r.toNum)                                                      // num --> digit+
	r.letters = r.letter.REP1().Convert(r.toLetters)                                             // letters --> letter+

	// rep_op --> "?" | "*" | "+"
	r.repOp = expectRune('?').ALT(
		expectRune('*'),
		expectRune('+'),
	).Convert(r.toRepOp)

	// upper_bound --> "," num?
	r.upperBound = expectRune(',').CONCAT(
		r.num.OPT(),
	).Convert(r.toUpperBound)

	// range --> "{" num upper_bound? "}"
	r.range_ = expectRune('{').CONCAT(
		r.num,
		r.upperBound.OPT(),
		expectRune('}'),
	).Convert(r.toRange)

	// repetition --> rep_op | range
	r.repetition = r.repOp.ALT(
		r.range_,
	).Convert(r.toRepetition)

	// quantifier --> repetition lazy_modifier?
	r.quantifier = r.repetition.CONCAT(
		expectRune('?').OPT(),
	).Convert(r.toQuantifier)

	// char_range --> char "-" char
	r.charRange = r.char.CONCAT(
		expectRune('-'),
		r.char,
	).Convert(r.toCharRange)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	r.charClass = expectString(`\d`).ALT(
		expectString(`\D`),
		expectString(`\s`), expectString(`\S`),
		expectString(`\w`), expectString(`\W`),
	).Convert(r.toCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	r.asciiCharClass = expectString("[:blank:]").ALT(
		expectString("[:space:]"),
		expectString("[:digit:]"), expectString("[:xdigit:]"),
		expectString("[:upper:]"), expectString("[:lower:]"),
		expectString("[:alpha:]"), expectString("[:alnum:]"),
		expectString("[:word:]"), expectString("[:ascii:]"),
	).Convert(r.toASCIICharClass)

	// char_group_item -->  char_class | ascii_char_class | char_range | char /* excluding ] */
	r.charGroupItem = r.charClass.ALT(
		r.asciiCharClass,
		r.charRange,
		r.charInGroup,
	).Convert(r.toCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	r.charGroup = expectRune('[').CONCAT(
		expectRune('^').OPT(),
		r.charGroupItem.REP1(),
		expectRune(']'),
	).Convert(r.toCharGroup)

	// any_char --> "."
	r.anyChar = expectRune('.').Convert(r.toAnyChar)

	// match_item --> any_char | char_class | ascii_char_class | char_group | char /* excluding | ) */
	r.matchItem = r.anyChar.ALT(
		r.charClass,
		r.asciiCharClass,
		r.charGroup,
		r.charInMatch,
	).Convert(r.toMatchItem)

	// match --> match_item quantifier?
	r.match = r.matchItem.CONCAT(r.quantifier.OPT()).Convert(r.toMatch)

	r.backref = expectRune('\\').CONCAT(r.num).Convert(r.toBackref) // backref --> "\" num
	r.anchor = expectRune('$').Convert(r.toAnchor)                  // anchor --> "$"

	// regex --> start_of_string? expr
	r.regex = expectRune('^').OPT().CONCAT(r.expr).Convert(r.toRegex)

	return r
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (r *regex) group(in input) (output, bool) {
	return expectRune('(').CONCAT(
		r.expr,
		expectRune(')'),
		r.quantifier.OPT(),
	).Convert(r.toGroup)(in)
}

// Recursive definition
// subexpr_item --> group | anchor | backref | match
func (r *regex) subexprItem(in input) (output, bool) {
	return parser(r.group).ALT(r.anchor, r.backref, r.match).Convert(r.toSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (r *regex) subexpr(in input) (output, bool) {
	return parser(r.subexprItem).REP1().Convert(r.toSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (r *regex) expr(in input) (output, bool) {
	return parser(r.subexpr).CONCAT(
		expectRune('|').CONCAT(r.expr).OPT(),
	).Convert(r.toExpr)(in)
}

//================================================== CONVERTERS ==================================================

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

func quantifyNode(n ast.Node, t tuple[any, bool]) ast.Node {
	var node ast.Node

	switch rep := t.p.(type) {
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
		low, up := rep.p, rep.q
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

func (r *regex) toChar(v any) (any, bool) {
	c, ok := v.(rune)
	if !ok {
		return nil, false
	}

	return runeToChar(c), true
}

func (r *regex) toNum(v any) (any, bool) {
	digits, ok := v.(list)
	if !ok {
		return nil, false
	}

	var num int
	for _, d := range digits {
		num = num*10 + int(d.Val.(rune)-'0')
	}

	return num, true
}

func (r *regex) toLetters(v any) (any, bool) {
	letters, ok := v.(list)
	if !ok {
		return nil, false
	}

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
	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	var num *int
	if v, ok := v1.(int); ok {
		num = &v
	}

	// Passing the result up the parsing chain
	return num, true
}

func (r *regex) toRange(v any) (any, bool) {
	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	v2, ok := getAt(v, 2)
	if !ok {
		return nil, false
	}

	low, ok := v1.(int)
	if !ok {
		return nil, false
	}

	// The upper bound is same as the lower bound if no upper bound is specified (default)
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
	v0, ok := getAt(v, 0)
	if !ok {
		return nil, false
	}

	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	// Check whether or not the lazy modifier is present
	_, lazy := v1.(rune)

	return tuple[any, bool]{
		p: v0,
		q: lazy,
	}, true
}

func (r *regex) toCharRange(v any) (any, bool) {
	v0, ok := getAt(v, 0)
	if !ok {
		return nil, false
	}

	v2, ok := getAt(v, 2)
	if !ok {
		return nil, false
	}

	low, ok := v0.(*ast.Char)
	if !ok {
		return nil, false
	}

	up, ok := v2.(*ast.Char)
	if !ok {
		return nil, false
	}

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
	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	v2, ok := getAt(v, 2)
	if !ok {
		return nil, false
	}

	// Check whether or not the negation modifier is present
	_, neg := v1.(rune)

	l, ok := v2.(list)
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

func (r *regex) toCharClass(v any) (any, bool) {
	class, ok := v.(string)
	if !ok {
		return nil, false
	}

	var alt *ast.Alt

	switch class {
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

func (r *regex) toASCIICharClass(v any) (any, bool) {
	class, ok := v.(string)
	if !ok {
		return nil, false
	}

	var alt *ast.Alt

	switch class {
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

func (r *regex) toAnyChar(v any) (any, bool) {
	if _, ok := v.(rune); !ok {
		return nil, false
	}

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
	v0, ok := getAt(v, 0)
	if !ok {
		return nil, false
	}

	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	node, ok := v0.(ast.Node)
	if !ok {
		return nil, false
	}

	if q, ok := v1.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	return node, true
}

func (r *regex) toBackref(v any) (any, bool) {
	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	num, ok := v1.(int)
	if !ok {
		return nil, false
	}

	var node ast.Node

	// Look up the symbol table
	if i := num - 1; 0 <= i && i < len(r.symTab) {
		node = r.symTab[i]
	} else {
		r.errors = multierror.Append(r.errors, fmt.Errorf("invalid back reference \\%d", num))
	}

	// Backref is successfully parsed, but it is invalid
	// If we return false, backref will be parsed by other parsers (match)
	return node, true
}

func (r *regex) toAnchor(v any) (any, bool) {
	a, ok := v.(rune)
	if !ok {
		return nil, false
	}

	// Check whether or not the anchor is end-of-string
	r.eos = a == '$'

	return a, true
}

func (r *regex) toGroup(v any) (any, bool) {
	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	v3, ok := getAt(v, 3)
	if !ok {
		return nil, false
	}

	node, ok := v1.(ast.Node)
	if !ok {
		return nil, false
	}

	if q, ok := v3.(tuple[any, bool]); ok {
		node = quantifyNode(node, q)
	}

	// Adding the group to the symbol table
	r.symTab = append(r.symTab, node)

	return node, true
}

func (r *regex) toSubexprItem(v any) (any, bool) {
	// Passing the result up the parsing chain
	return v, true
}

func (r *regex) toSubexpr(v any) (any, bool) {
	l, ok := v.(list)
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

func (r *regex) toExpr(v any) (any, bool) {
	v0, ok := getAt(v, 0)
	if !ok {
		return nil, false
	}

	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	subexpr, ok := v0.(ast.Node)
	if !ok {
		return nil, false
	}

	if _, ok := v1.(empty); ok {
		return subexpr, true
	}

	v11, ok := getAt(v1, 1)
	if !ok {
		return nil, false
	}

	expr, ok := v11.(ast.Node)
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

func (r *regex) toRegex(v any) (any, bool) {
	v0, ok := getAt(v, 0)
	if !ok {
		return nil, false
	}

	v1, ok := getAt(v, 1)
	if !ok {
		return nil, false
	}

	// Check whether or not the start-of-string is present
	_, r.sos = v0.(rune)

	expr, ok := v1.(ast.Node)
	if !ok {
		return nil, false
	}

	return expr, true
}
