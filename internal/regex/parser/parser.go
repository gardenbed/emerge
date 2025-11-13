// Package parser provides combinators for parsing regular expressions.
package parser

import (
	"fmt"

	comb "github.com/moorara/algo/parser/combinator"
)

// excludeRunes can be bound on a rune parser to exclude certain runes.
func excludeRunes(r ...rune) comb.BindFunc {
	return func(res comb.Result) comb.Parser {
		return func(in comb.Input) (*comb.Output, error) {
			if a, ok := res.Val.(rune); ok {
				for _, b := range r {
					if a == b {
						return nil, fmt.Errorf("%d: unexpected rune %q", res.Pos, a)
					}
				}
			}

			return &comb.Output{
				Result:    res,
				Remaining: in,
			}, nil
		}
	}
}

func toDigit(r comb.Result) (comb.Result, error) {
	v := r.Val.(rune)

	digit := int(v - '0')

	return comb.Result{
		Val: digit,
		Pos: r.Pos,
	}, nil
}

func toHexDigit(r comb.Result) (comb.Result, error) {
	v := r.Val.(rune)

	var digit int
	if '0' <= v && v <= '9' {
		digit = int(v - '0')
	} else { // 'A' <= v && v <= 'F'
		digit = int(v - 55)
	}

	return comb.Result{
		Val: digit,
		Pos: r.Pos,
	}, nil
}

func toNum(r comb.Result) (comb.Result, error) {
	l := r.Val.(comb.List)

	var num int
	for _, r := range l {
		num = num*10 + r.Val.(int)
	}

	return comb.Result{
		Val: num,
		Pos: l[0].Pos,
	}, nil
}

func toLetters(r comb.Result) (comb.Result, error) {
	l := r.Val.(comb.List)

	var str string
	for _, r := range l {
		str += string(r.Val.(rune))
	}

	return comb.Result{
		Val: str,
		Pos: l[0].Pos,
	}, nil
}

func toEscapedChar(r comb.Result) (comb.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	c := r1.Val.(rune)

	switch c {
	case 't':
		c = '\t'
	case 'n':
		c = '\n'
	case 'r':
		c = '\r'
	}

	return comb.Result{
		Val: c,
		Pos: r0.Pos,
	}, nil
}

func toASCIIChar(r comb.Result) (comb.Result, error) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	d1 := r1.Val.(int)
	d2 := r2.Val.(int)
	c := d1<<4 + d2

	return comb.Result{
		Val: rune(c),
		Pos: r0.Pos,
	}, nil
}

func toUnicodeChar(r comb.Result) (comb.Result, error) {
	l := r.Val.(comb.List)

	var c int
	for _, r := range l[1:] {
		if d, ok := r.Val.(int); ok {
			c = c<<4 + d
		}
	}

	return comb.Result{
		Val: rune(c),
		Pos: l[0].Pos,
	}, nil
}

// Mappers defines mapping functions for various non-terminals in the regex grammar.
type Mappers interface {
	ToAnyChar(comb.Result) (comb.Result, error)          // any_char --> "."
	ToSingleChar(comb.Result) (comb.Result, error)       // single_char --> unicode_char | ascii_char | escaped_char | raw_char
	ToCharClass(comb.Result) (comb.Result, error)        // char_class --> "\s" | "\S" | "\d" | "\D" | "\w" | "\W"
	ToASCIICharClass(comb.Result) (comb.Result, error)   // ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | ...
	ToUnicodeCategory(comb.Result) (comb.Result, error)  // unicode_category --> ...
	ToUnicodeCharClass(comb.Result) (comb.Result, error) // unicode_char_class --> ("\p" | "\P") "{" unicode_category "}"
	ToRepOp(comb.Result) (comb.Result, error)            // rep_op --> "?" | "*" | "+"
	ToUpperBound(comb.Result) (comb.Result, error)       // upper_bound --> "," num?
	ToRange(comb.Result) (comb.Result, error)            // range --> "{" num upper_bound? "}"
	ToRepetition(comb.Result) (comb.Result, error)       // repetition --> rep_op | range
	ToQuantifier(comb.Result) (comb.Result, error)       // quantifier --> repetition lazy_modifier?
	ToCharInGroup(comb.Result) (comb.Result, error)      // char_in_group --> unicode_char | ascii_char | escaped_char | raw_char_in_group
	ToCharRange(comb.Result) (comb.Result, error)        // char_range --> char_in_range "-" char_in_range
	ToCharGroupItem(comb.Result) (comb.Result, error)    // char_group_item --> unicode_char_class | ascii_char_class | char_class | char_range | char_in_group
	ToCharGroup(comb.Result) (comb.Result, error)        // char_group --> "[" "^"? char_group_item+ "]"
	ToMatchItem(comb.Result) (comb.Result, error)        // match_item --> any_char | single_char | char_class | ascii_char_class | unicode_char_class | char_group
	ToMatch(comb.Result) (comb.Result, error)            // match --> match_item quantifier?
	ToGroup(comb.Result) (comb.Result, error)            // group --> "(" expr ")" quantifier?
	ToSubexprItem(comb.Result) (comb.Result, error)      // subexpr_item --> group | match
	ToSubexpr(comb.Result) (comb.Result, error)          // subexpr --> subexpr_item+
	ToExpr(comb.Result) (comb.Result, error)             // expr --> subexpr ("|" expr)?
	ToRegex(comb.Result) (comb.Result, error)            // regex --> expr
}

// Parser is a parser combinator for regular expressions.
type Parser struct {
	m Mappers

	// Combinators
	digit            comb.Parser
	hexDigit         comb.Parser
	letter           comb.Parser
	num              comb.Parser
	letters          comb.Parser
	char             comb.Parser
	rawCharInGroup   comb.Parser
	rawChar          comb.Parser
	escapedChar      comb.Parser
	asciiChar        comb.Parser
	unicodeChar      comb.Parser
	anyChar          comb.Parser
	singleChar       comb.Parser
	charClass        comb.Parser
	asciiCharClass   comb.Parser
	unicodeCategory  comb.Parser
	unicodeCharClass comb.Parser
	repOp            comb.Parser
	upperBound       comb.Parser
	range_           comb.Parser
	repetition       comb.Parser
	quantifier       comb.Parser
	charInGroup      comb.Parser
	charRange        comb.Parser
	charGroupItem    comb.Parser
	charGroup        comb.Parser
	matchItem        comb.Parser
	match            comb.Parser
	regex            comb.Parser
}

// New creates a parser combinator for regular expressions.
func New(m Mappers) *Parser {
	p := &Parser{
		m: m,
	}

	p.digit = comb.ExpectRuneInRange('0', '9').Map(toDigit)                                                   // digit --> "0" | ... | "9"
	p.hexDigit = comb.ALT(comb.ExpectRuneInRange('0', '9'), comb.ExpectRuneInRange('A', 'F')).Map(toHexDigit) // hex_digit --> "0" | ... | "9" | "A" | ... | "F"
	p.letter = comb.ALT(comb.ExpectRuneInRange('A', 'Z'), comb.ExpectRuneInRange('a', 'z'))                   // letter --> "A" | ... | "Z" | "a" | ... | "z"
	p.num = p.digit.REP1().Map(toNum)                                                                         // num --> digit+
	p.letters = p.letter.REP1().Map(toLetters)                                                                // letters --> letter+

	// char --> all Unicode characters
	p.char = comb.ExpectRuneInRange(0x20, 0x10FFFF)

	// raw_char_in_group --> all characters except ...
	p.rawCharInGroup = p.char.Bind(
		excludeRunes('/', '\\', '\t', '\n', '\r', '[', ']'),
	)

	// raw_char --> all characters except the escaped ones
	p.rawChar = p.char.Bind(
		excludeRunes('/', '\\', '\t', '\n', '\r', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}'),
	)

	// escaped_char --> "\" (...)
	p.escapedChar = comb.ExpectRune('\\').CONCAT(
		comb.ExpectRuneIn('/', '\\', 't', 'n', 'r', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}'),
	).Map(toEscapedChar)

	// ascii_char --> "\x" hex_digit{2}
	p.asciiChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit).Map(toASCIIChar)

	// unicode_char --> "\x" hex_digit{4,8}
	p.unicodeChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit, p.hexDigit, p.hexDigit,
		p.hexDigit.OPT(), p.hexDigit.OPT(), p.hexDigit.OPT(), p.hexDigit.OPT(),
	).Map(toUnicodeChar)

	// any_char --> "."
	p.anyChar = comb.ExpectRune('.').Map(p.m.ToAnyChar)

	// single_char --> unicode_char | ascii_char | escaped_char | raw_char
	p.singleChar = comb.ALT(
		p.unicodeChar,
		p.asciiChar,
		p.escapedChar,
		p.rawChar,
	).Map(p.m.ToSingleChar)

	// char_class --> "\s" | "\S" | "\d" | "\D" | "\w" | "\W"
	p.charClass = comb.ALT(
		comb.ExpectString(`\s`), comb.ExpectString(`\S`),
		comb.ExpectString(`\d`), comb.ExpectString(`\D`),
		comb.ExpectString(`\w`), comb.ExpectString(`\W`),
	).Map(p.m.ToCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	p.asciiCharClass = comb.ALT(
		comb.ExpectString("[:blank:]"), comb.ExpectString("[:space:]"),
		comb.ExpectString("[:digit:]"), comb.ExpectString("[:xdigit:]"),
		comb.ExpectString("[:upper:]"), comb.ExpectString("[:lower:]"),
		comb.ExpectString("[:alpha:]"), comb.ExpectString("[:alnum:]"),
		comb.ExpectString("[:word:]"), comb.ExpectString("[:ascii:]"),
	).Map(p.m.ToASCIICharClass)

	// unicode_category --> ...
	p.unicodeCategory = comb.ALT(
		// Derived
		comb.ExpectString("Math"), comb.ExpectString("Emoji"),
		// Scripts
		comb.ExpectString("Latin"), comb.ExpectString("Greek"), comb.ExpectString("Cyrillic"), comb.ExpectString("Han"), comb.ExpectString("Persian"),
		// General
		comb.ExpectString("Letter"), comb.ExpectString("Lu"), comb.ExpectString("Ll"), comb.ExpectString("Lt"),
		comb.ExpectString("Lm"), comb.ExpectString("Lo"), comb.ExpectString("L"),
		comb.ExpectString("Mark"), comb.ExpectString("Mn"), comb.ExpectString("Mc"), comb.ExpectString("Me"), comb.ExpectString("M"),
		comb.ExpectString("Number"), comb.ExpectString("Nd"), comb.ExpectString("Nl"), comb.ExpectString("No"), comb.ExpectString("N"),
		comb.ExpectString("Punctuation"), comb.ExpectString("Pc"), comb.ExpectString("Pd"), comb.ExpectString("Ps"), comb.ExpectString("Pe"),
		comb.ExpectString("Pi"), comb.ExpectString("Pf"), comb.ExpectString("Po"), comb.ExpectString("P"),
		comb.ExpectString("Separator"), comb.ExpectString("Zs"), comb.ExpectString("Zl"), comb.ExpectString("Zp"), comb.ExpectString("Z"),
		comb.ExpectString("Symbol"), comb.ExpectString("Sm"), comb.ExpectString("Sc"), comb.ExpectString("Sk"), comb.ExpectString("So"), comb.ExpectString("S"),
	).Map(p.m.ToUnicodeCategory)

	// unicode_char_class --> ("\p" | "\P") "{" unicode_category "}"
	p.unicodeCharClass = comb.ExpectString(`\p`).ALT(comb.ExpectString(`\P`)).CONCAT(
		comb.ExpectRune('{'),
		p.unicodeCategory,
		comb.ExpectRune('}'),
	).Map(p.m.ToUnicodeCharClass)

	// rep_op --> "?" | "*" | "+"
	p.repOp = comb.ALT(
		comb.ExpectRune('?'),
		comb.ExpectRune('*'),
		comb.ExpectRune('+'),
	).Map(p.m.ToRepOp)

	// upper_bound --> "," num?
	p.upperBound = comb.ExpectRune(',').CONCAT(p.num.OPT()).Map(p.m.ToUpperBound)

	// range --> "{" num upper_bound? "}"
	p.range_ = comb.CONCAT(
		comb.ExpectRune('{'),
		p.num,
		p.upperBound.OPT(),
		comb.ExpectRune('}'),
	).Map(p.m.ToRange)

	p.repetition = p.repOp.ALT(p.range_).Map(p.m.ToRepetition)                           // repetition --> rep_op | range
	p.quantifier = p.repetition.CONCAT(comb.ExpectRune('?').OPT()).Map(p.m.ToQuantifier) // quantifier --> repetition lazy_modifier?

	// char_in_group --> unicode_char | ascii_char | escaped_char | raw_char_in_group
	p.charInGroup = comb.ALT(
		p.unicodeChar,
		p.asciiChar,
		p.escapedChar,
		p.rawCharInGroup,
	).Map(p.m.ToCharInGroup)

	// char_range --> char_in_group "-" char_in_group
	p.charRange = comb.CONCAT(
		p.charInGroup,
		comb.ExpectRune('-'),
		p.charInGroup,
	).Map(p.m.ToCharRange)

	// char_group_item --> unicode_char_class | ascii_char_class | char_class | char_range | char_in_group
	p.charGroupItem = comb.ALT(
		p.unicodeCharClass,
		p.asciiCharClass,
		p.charClass,
		p.charRange,
		p.charInGroup,
	).Map(p.m.ToCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	p.charGroup = comb.CONCAT(
		comb.ExpectRune('['),
		comb.ExpectRune('^').OPT(),
		p.charGroupItem.REP1(),
		comb.ExpectRune(']'),
	).Map(p.m.ToCharGroup)

	// match_item --> any_char | single_char | char_class | ascii_char_class | unicode_char_class | char_group
	p.matchItem = comb.ALT(
		p.anyChar,
		p.singleChar,
		p.charClass,
		p.asciiCharClass,
		p.unicodeCharClass,
		p.charGroup,
	).Map(p.m.ToMatchItem)

	p.match = p.matchItem.CONCAT(p.quantifier.OPT()).Map(p.m.ToMatch) // match --> match_item quantifier?
	p.regex = comb.Parser(p.expr).Map(p.m.ToRegex)                    // regex --> expr

	return p
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (p *Parser) group(in comb.Input) (*comb.Output, error) {
	return comb.CONCAT(
		comb.ExpectRune('('),
		p.expr,
		comb.ExpectRune(')'),
		p.quantifier.OPT(),
	).Map(p.m.ToGroup)(in)
}

// Recursive definition
// subexpr_item --> group | match
func (p *Parser) subexprItem(in comb.Input) (*comb.Output, error) {
	return comb.ALT(p.group, p.match).Map(p.m.ToSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (p *Parser) subexpr(in comb.Input) (*comb.Output, error) {
	return comb.Parser(p.subexprItem).REP1().Map(p.m.ToSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (p *Parser) expr(in comb.Input) (*comb.Output, error) {
	return comb.Parser(p.subexpr).CONCAT(
		comb.ExpectRune('|').CONCAT(p.expr).OPT(),
	).Map(p.m.ToExpr)(in)
}

// Parse is the topmost parser combinator for parsing a regular expression read from the input.
func (p *Parser) Parse(regex string) (*comb.Output, error) {
	in := newStringInput(regex)

	out, err := p.regex(in)
	if err != nil {
		return nil, err
	}

	// Ensure that the entire input has been matched.
	if out.Remaining != nil {
		curr, pos := out.Remaining.Current()
		return nil, fmt.Errorf("%d: unexpected rune %q", pos, curr)
	}

	return out, nil
}
