package parser

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	ast "github.com/gardenbed/emerge/internal/regex/ast/compact"
	"github.com/gardenbed/emerge/internal/regex/token"
)

func ParseCompact(in input) (*ast.Regex, error) {
	c := new(compactConverters)
	r := NewRegex(c)

	out, ok := r.regex(in)
	if !ok {
		return nil, errors.New("syntax error")
	}

	// Check for errors
	if c.errors != nil {
		return nil, c.errors
	}

	root := out.Result.Val.(ast.Regex)

	return &root, nil
}

// compactConverters implements the converters interface for the compact regex AST.
type compactConverters struct {
	errors error
	symTab []ast.Group
}

func (c *compactConverters) ToChar(res result) (any, bool) {
	r, ok := res.Val.(rune)
	if !ok {
		return nil, false
	}

	return ast.Char{
		TokPos: token.Pos(res.Pos),
		Val:    r,
	}, true
}

func (c *compactConverters) ToNum(res result) (any, bool) {
	digits, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	start := res.Pos
	end := digits[len(digits)-1].Pos

	var num int
	for _, d := range digits {
		num = num*10 + int(d.Val.(rune)-'0')
	}

	return ast.Num{
		StartPos: token.Pos(start),
		EndPos:   token.Pos(end),
		Val:      num,
	}, true
}

func (c *compactConverters) ToLetters(res result) (any, bool) {
	letters, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	start := res.Pos
	end := letters[len(letters)-1].Pos

	var s string
	for _, l := range letters {
		s += string(l.Val.(rune))
	}

	return ast.Letters{
		StartPos: token.Pos(start),
		EndPos:   token.Pos(end),
		Val:      s,
	}, true
}

func (c *compactConverters) ToRepOp(res result) (any, bool) {
	val, ok := res.Val.(rune)
	if !ok {
		return nil, false
	}

	var tag token.Tag

	switch val {
	case '?':
		tag = token.ZERO_OR_ONE
	case '*':
		tag = token.ZERO_OR_MORE
	case '+':
		tag = token.ONE_OR_MORE
	}

	return ast.RepOp{
		TokPos: token.Pos(res.Pos),
		TokTag: tag,
	}, true
}

func (c *compactConverters) ToUpperBound(res result) (any, bool) {
	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	var num *ast.Num
	if v, ok := r1.Val.(ast.Num); ok {
		num = &v
	}

	return ast.UpperBound{
		CommaPos: token.Pos(r1.Pos),
		Val:      num,
	}, true
}

func (c *compactConverters) ToRange(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	r3, ok := res.Get(3)
	if !ok {
		return nil, false
	}

	low, ok := r1.Val.(ast.Num)
	if !ok {
		return nil, false
	}

	var up *ast.UpperBound
	if v, ok := r2.Val.(ast.UpperBound); ok {
		up = &v
	}

	if up != nil && up.Val != nil && low.Val > up.Val.Val {
		c.errors = multierror.Append(c.errors,
			fmt.Errorf("invalid repetition range {%d,%d}", low.Val, up.Val.Val),
		)
	}

	return ast.Range{
		OpenPos:  token.Pos(r0.Pos),
		ClosePos: token.Pos(r3.Pos),
		Low:      low,
		Up:       up,
	}, true
}

func (c *compactConverters) ToRepetition(res result) (any, bool) {
	switch v := res.Val.(type) {
	case ast.RepOp:
		return &v, true
	case ast.Range:
		return &v, true
	default:
		return nil, false
	}
}

func (c *compactConverters) ToQuantifier(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	rep, ok := r0.Val.(ast.Repetition)
	if !ok {
		return nil, false
	}

	// Check whether or not the lazy modifier is present
	_, lazy := r1.Val.(rune)

	return ast.Quantifier{
		Rep:  rep,
		Lazy: lazy,
	}, true
}

func (c *compactConverters) ToCharRange(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	low, ok := r0.Val.(ast.Char)
	if !ok {
		return nil, false
	}

	up, ok := r2.Val.(ast.Char)
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

	return ast.CharRange{
		Low: low,
		Up:  up,
	}, true
}

func (c *compactConverters) ToCharGroupItem(res result) (any, bool) {
	switch v := res.Val.(type) {
	case ast.CharClass:
		return &v, true
	case ast.ASCIICharClass:
		return &v, true
	case ast.CharRange:
		return &v, true
	case ast.Char:
		return &v, true
	default:
		return nil, false
	}
}

func (c *compactConverters) ToCharGroup(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	r3, ok := res.Get(3)
	if !ok {
		return nil, false
	}

	// Check whether or not the negation modifier is present
	_, negated := r1.Val.(rune)

	l, ok := r2.Val.(list)
	if !ok {
		return nil, false
	}

	var items []ast.CharGroupItem
	for _, r := range l {
		if i, ok := r.Val.(ast.CharGroupItem); ok {
			items = append(items, i)
		}
	}

	return ast.CharGroup{
		OpenPos:  token.Pos(r0.Pos),
		ClosePos: token.Pos(r3.Pos),
		Negated:  negated,
		Items:    items,
	}, true
}

func (c *compactConverters) ToCharClass(res result) (any, bool) {
	val, ok := res.Val.(string)
	if !ok {
		return nil, false
	}

	start := res.Pos
	end := start + len(val) - 1

	var tag token.Tag

	switch val {
	case `\d`:
		tag = token.DIGIT
	case `\D`:
		tag = token.NON_DIGIT
	case `\s`:
		tag = token.WHITESPACE
	case `\S`:
		tag = token.NON_WHITESPACE
	case `\w`:
		tag = token.WORD
	case `\W`:
		tag = token.NON_WORD
	default:
		return nil, false
	}

	return ast.CharClass{
		StartPos: token.Pos(start),
		EndPos:   token.Pos(end),
		TokTag:   tag,
	}, true
}

func (c *compactConverters) ToASCIICharClass(res result) (any, bool) {
	val, ok := res.Val.(string)
	if !ok {
		return nil, false
	}

	start := res.Pos
	end := start + len(val) - 1

	var tag token.Tag

	switch val {
	case "[:blank:]":
		tag = token.BLANK_CHARS
	case "[:space:]":
		tag = token.SPACE_CHARS
	case "[:digit:]":
		tag = token.DIGIT_CHARS
	case "[:xdigit:]":
		tag = token.XDIGIT_CHARS
	case "[:upper:]":
		tag = token.UPPER_CHARS
	case "[:lower:]":
		tag = token.LOWER_CHARS
	case "[:alpha:]":
		tag = token.ALPHA_CHARS
	case "[:alnum:]":
		tag = token.ALNUM_CHARS
	case "[:word:]":
		tag = token.WORD_CHARS
	case "[:ascii:]":
		tag = token.ASCII_CHARS
	default:
		return nil, false
	}

	return ast.ASCIICharClass{
		StartPos: token.Pos(start),
		EndPos:   token.Pos(end),
		TokTag:   tag,
	}, true
}

func (c *compactConverters) ToAnyChar(res result) (any, bool) {
	if _, ok := res.Val.(rune); !ok {
		return nil, false
	}

	return ast.AnyChar{
		TokPos: token.Pos(res.Pos),
	}, true
}

func (c *compactConverters) ToMatchItem(res result) (any, bool) {
	switch v := res.Val.(type) {
	case ast.AnyChar:
		return &v, true
	case ast.CharGroup:
		return &v, true
	case ast.CharClass:
		return &v, true
	case ast.ASCIICharClass:
		return &v, true
	case ast.Char:
		return &v, true
	default:
		return nil, false
	}
}

func (c *compactConverters) ToMatch(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	item, ok := r0.Val.(ast.MatchItem)
	if !ok {
		return nil, false
	}

	var quant *ast.Quantifier
	if v, ok := r1.Val.(ast.Quantifier); ok {
		quant = &v
	}

	return ast.Match{
		Item:  item,
		Quant: quant,
	}, true
}

func (c *compactConverters) ToBackref(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	num, ok := r1.Val.(ast.Num)
	if !ok {
		return nil, false
	}

	backref := ast.Backref{
		SlashPos: token.Pos(r0.Pos),
		Ref:      num,
	}

	// Look up the symbol table
	if i := num.Val - 1; 0 <= i && i < len(c.symTab) {
		backref.Group = &c.symTab[i]
	} else {
		c.errors = multierror.Append(c.errors, fmt.Errorf("invalid back reference \\%d", num.Val))
	}

	// Backref is successfully parsed, but it is invalid
	// If we return false, backref will be parsed by other parsers (match)
	return backref, true
}

func (c *compactConverters) ToAnchor(res result) (any, bool) {
	if _, ok := res.Val.(rune); !ok {
		return nil, false
	}

	return ast.Anchor{
		TokPos: token.Pos(res.Pos),
	}, true
}

func (c *compactConverters) ToGroup(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	r2, ok := res.Get(2)
	if !ok {
		return nil, false
	}

	r3, ok := res.Get(3)
	if !ok {
		return nil, false
	}

	expr, ok := r1.Val.(ast.Expr)
	if !ok {
		return nil, false
	}

	var quant *ast.Quantifier
	if v, ok := r3.Val.(ast.Quantifier); ok {
		quant = &v
	}

	group := ast.Group{
		OpenPos:  token.Pos(r0.Pos),
		ClosePos: token.Pos(r2.Pos),
		Expr:     expr,
		Quant:    quant,
	}

	// Adding the group to the symbol table
	c.symTab = append(c.symTab, group)

	return group, true
}

func (c *compactConverters) ToSubexprItem(res result) (any, bool) {
	switch v := res.Val.(type) {
	case ast.Group:
		return &v, true
	case ast.Anchor:
		return &v, true
	case ast.Backref:
		return &v, true
	case ast.Match:
		return &v, true
	default:
		return nil, false
	}
}

func (c *compactConverters) ToSubexpr(res result) (any, bool) {
	l, ok := res.Val.(list)
	if !ok {
		return nil, false
	}

	var items []ast.SubexprItem
	for _, r := range l {
		if item, ok := r.Val.(ast.SubexprItem); ok {
			items = append(items, item)
		}
	}

	return ast.Subexpr{
		Items: items,
	}, true
}

func (c *compactConverters) ToExpr(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	sub, ok := r0.Val.(ast.Subexpr)
	if !ok {
		return nil, false
	}

	var expr *ast.Expr

	if _, ok := r1.Val.(empty); !ok {
		r11, ok := r1.Get(1)
		if !ok {
			return nil, false
		}

		v, ok := r11.Val.(ast.Expr)
		if !ok {
			return nil, false
		}

		expr = &v
	}

	return ast.Expr{
		Sub:  sub,
		Expr: expr,
	}, true
}

func (c *compactConverters) ToRegex(res result) (any, bool) {
	r0, ok := res.Get(0)
	if !ok {
		return nil, false
	}

	r1, ok := res.Get(1)
	if !ok {
		return nil, false
	}

	// Check whether or not the start-of-string is present
	_, sos := r0.Val.(rune)

	expr, ok := r1.Val.(ast.Expr)
	if !ok {
		return nil, false
	}

	return ast.Regex{
		SOS:  sos,
		Expr: expr,
	}, true
}
