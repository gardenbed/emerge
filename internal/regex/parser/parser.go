// Package parser provides combinators for parsing regular expressions.
package parser

import comb "github.com/moorara/algo/parser/combinator"

var escapedChars = []rune{'\\', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}', '$'}

// excludeRunes can be bound on a rune parser to exclude certain runes.
func excludeRunes(r ...rune) comb.BindFunc {
	return func(res comb.Result) comb.Parser {
		return func(in comb.Input) (comb.Output, bool) {
			if a, ok := res.Val.(rune); ok {
				for _, b := range r {
					if a == b {
						return comb.Output{}, false
					}
				}
			}

			return comb.Output{
				Result:    res,
				Remaining: in,
			}, true
		}
	}
}

func toDigit(r comb.Result) (comb.Result, bool) {
	v := r.Val.(rune)

	digit := int(v - '0')

	return comb.Result{
		Val: digit,
		Pos: r.Pos,
	}, true
}

func toHexDigit(r comb.Result) (comb.Result, bool) {
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
	}, true
}

func toNum(r comb.Result) (comb.Result, bool) {
	l := r.Val.(comb.List)

	var num int
	for _, r := range l {
		num = num*10 + r.Val.(int)
	}

	return comb.Result{
		Val: num,
		Pos: l[0].Pos,
	}, true
}

func toLetters(r comb.Result) (comb.Result, bool) {
	l := r.Val.(comb.List)

	var str string
	for _, r := range l {
		str += string(r.Val.(rune))
	}

	return comb.Result{
		Val: str,
		Pos: l[0].Pos,
	}, true
}

func toEscapedChar(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	c := r1.Val.(rune)
	switch c {
	case 't':
		c = '\t'
	case 'n':
		c = '\n'
	case 'v':
		c = '\v'
	case 'f':
		c = '\f'
	case 'r':
		c = '\r'
	}

	return comb.Result{
		Val: c,
		Pos: r0.Pos,
	}, true
}

func toASCIIChar(r comb.Result) (comb.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	d1 := r1.Val.(int)
	d2 := r2.Val.(int)
	c := d1<<4 + d2

	return comb.Result{
		Val: rune(c),
		Pos: r0.Pos,
	}, true
}

func toUnicodeChar(r comb.Result) (comb.Result, bool) {
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
	}, true
}

// Mappers defines mapping functions for various non-terminals in the regex grammar.
type Mappers interface {
	ToAnyChar(comb.Result) (comb.Result, bool)          // any_char --> "."
	ToSingleChar(comb.Result) (comb.Result, bool)       // single_char --> unicode_char | ascii_char | escaped_char | unescaped_char
	ToCharClass(comb.Result) (comb.Result, bool)        // char_class --> "\s" | "\S" | "\d" | "\D" | "\w" | "\W"
	ToASCIICharClass(comb.Result) (comb.Result, bool)   // ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | ...
	ToUnicodeCategory(comb.Result) (comb.Result, bool)  // unicode_category --> ...
	ToUnicodeCharClass(comb.Result) (comb.Result, bool) // unicode_char_class --> "\p" "{" unicode_category "}"
	ToRepOp(comb.Result) (comb.Result, bool)            // rep_op --> "?" | "*" | "+"
	ToUpperBound(comb.Result) (comb.Result, bool)       // upper_bound --> "," num?
	ToRange(comb.Result) (comb.Result, bool)            // range --> "{" num upper_bound? "}"
	ToRepetition(comb.Result) (comb.Result, bool)       // repetition --> rep_op | range
	ToQuantifier(comb.Result) (comb.Result, bool)       // quantifier --> repetition lazy_modifier?
	ToCharInRange(comb.Result) (comb.Result, bool)      // char_in_range --> unicode_char | ascii_char | char
	ToCharRange(comb.Result) (comb.Result, bool)        // char_range --> char_in_range "-" char_in_range
	ToCharGroupItem(comb.Result) (comb.Result, bool)    // char_group_item --> unicode_char_class | ascii_char_class | char_class | char_range | single_char
	ToCharGroup(comb.Result) (comb.Result, bool)        // char_group --> "[" "^"? char_group_item+ "]"
	ToMatchItem(comb.Result) (comb.Result, bool)        // match_item --> any_char | single_char | char_class | ascii_char_class | unicode_char_class | char_group
	ToMatch(comb.Result) (comb.Result, bool)            // match --> match_item quantifier?
	ToGroup(comb.Result) (comb.Result, bool)            // group --> "(" expr ")" quantifier?
	ToAnchor(comb.Result) (comb.Result, bool)           // anchor --> "$"
	ToSubexprItem(comb.Result) (comb.Result, bool)      // subexpr_item --> anchor | group | match
	ToSubexpr(comb.Result) (comb.Result, bool)          // subexpr --> subexpr_item+
	ToExpr(comb.Result) (comb.Result, bool)             // expr --> subexpr ("|" expr)?
	ToRegex(comb.Result) (comb.Result, bool)            // regex --> start_of_string? expr
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
	unescapedChar    comb.Parser
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
	charInRange      comb.Parser
	charRange        comb.Parser
	charGroupItem    comb.Parser
	charGroup        comb.Parser
	matchItem        comb.Parser
	match            comb.Parser
	anchor           comb.Parser
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

	p.char = comb.ExpectRuneInRange(0x20, 0x10FFFF)                                                     // char --> all Unicode characters
	p.unescapedChar = p.char.Bind(excludeRunes(escapedChars...))                                        // unescaped_char --> all characters excluding the escaped ones
	p.escapedChar = comb.ExpectRune('\\').CONCAT(comb.ExpectRuneIn(escapedChars...)).Map(toEscapedChar) // escaped_char --> "\" ( "\" | ... | "$" )
	p.asciiChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit).Map(toASCIIChar)               // ascii_char --> "\x" hex_digit{2}

	// unicode_char --> "\x" hex_digit{4,8}
	p.unicodeChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit, p.hexDigit, p.hexDigit,
		p.hexDigit.OPT(), p.hexDigit.OPT(), p.hexDigit.OPT(), p.hexDigit.OPT(),
	).Map(toUnicodeChar)

	// any_char --> "."
	p.anyChar = comb.ExpectRune('.').Map(p.m.ToAnyChar)

	// single_char --> unicode_char | ascii_char | escaped_char | unescaped_char
	p.singleChar = comb.ALT(p.unicodeChar, p.asciiChar, p.escapedChar, p.unescapedChar).Map(p.m.ToSingleChar)

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
	p.unicodeCategory = comb.ExpectString("Letter").ALT(
		// Derived
		comb.ExpectString("Math"), comb.ExpectString("Emoji"),
		// Scripts
		comb.ExpectString("Latin"), comb.ExpectString("Greek"), comb.ExpectString("Cyrillic"), comb.ExpectString("Han"), comb.ExpectString("Persian"),
		// General
		comb.ExpectString("Lu"), comb.ExpectString("Ll"), comb.ExpectString("Lt"), comb.ExpectString("Lm"), comb.ExpectString("Lo"), comb.ExpectString("L"),
		comb.ExpectString("Mark"), comb.ExpectString("Mn"), comb.ExpectString("Mc"), comb.ExpectString("Me"), comb.ExpectString("M"),
		comb.ExpectString("Number"), comb.ExpectString("Nd"), comb.ExpectString("Nl"), comb.ExpectString("No"), comb.ExpectString("N"),
		comb.ExpectString("Punctuation"), comb.ExpectString("Pc"), comb.ExpectString("Pd"), comb.ExpectString("Ps"),
		comb.ExpectString("Pe"), comb.ExpectString("Pi"), comb.ExpectString("Pf"), comb.ExpectString("Po"), comb.ExpectString("P"),
		comb.ExpectString("Separator"), comb.ExpectString("Zs"), comb.ExpectString("Zl"), comb.ExpectString("Zp"), comb.ExpectString("Z"),
		comb.ExpectString("Symbol"), comb.ExpectString("Sm"), comb.ExpectString("Sc"), comb.ExpectString("Sk"), comb.ExpectString("So"), comb.ExpectString("S"),
	).Map(p.m.ToUnicodeCategory)

	// unicode_char_class --> "\p" "{" unicode_category "}"
	p.unicodeCharClass = comb.ExpectString(`\p`).ALT(comb.ExpectString(`\P`)).CONCAT(
		comb.ExpectRune('{'),
		p.unicodeCategory,
		comb.ExpectRune('}'),
	).Map(p.m.ToUnicodeCharClass)

	// rep_op --> "?" | "*" | "+"
	p.repOp = comb.ExpectRune('?').ALT(
		comb.ExpectRune('*'),
		comb.ExpectRune('+'),
	).Map(p.m.ToRepOp)

	// upper_bound --> "," num?
	p.upperBound = comb.ExpectRune(',').CONCAT(
		p.num.OPT(),
	).Map(p.m.ToUpperBound)

	// range --> "{" num upper_bound? "}"
	p.range_ = comb.ExpectRune('{').CONCAT(
		p.num,
		p.upperBound.OPT(),
		comb.ExpectRune('}'),
	).Map(p.m.ToRange)

	// repetition --> rep_op | range
	p.repetition = p.repOp.ALT(
		p.range_,
	).Map(p.m.ToRepetition)

	// quantifier --> repetition lazy_modifier?
	p.quantifier = p.repetition.CONCAT(
		comb.ExpectRune('?').OPT(),
	).Map(p.m.ToQuantifier)

	// char_in_range --> unicode_char | ascii_char | char
	p.charInRange = p.unicodeChar.ALT(p.asciiChar, p.char).Map(p.m.ToCharInRange)

	// char_range --> char_in_range "-" char_in_range
	p.charRange = comb.CONCAT(
		p.charInRange,
		comb.ExpectRune('-'),
		p.charInRange,
	).Map(p.m.ToCharRange)

	// char_group_item --> unicode_char_class | ascii_char_class | char_class | char_range | single_char
	p.charGroupItem = comb.ALT(
		p.unicodeCharClass,
		p.asciiCharClass,
		p.charClass,
		p.charRange,
		p.singleChar,
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

	// match --> match_item quantifier?
	p.match = p.matchItem.CONCAT(p.quantifier.OPT()).Map(p.m.ToMatch)

	// anchor --> "$"
	p.anchor = comb.ExpectRune('$').Map(p.m.ToAnchor)

	// regex --> start_of_string? expr
	p.regex = comb.ExpectRune('^').OPT().CONCAT(p.expr).Map(p.m.ToRegex)

	return p
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (p *Parser) group(in comb.Input) (comb.Output, bool) {
	return comb.ExpectRune('(').CONCAT(
		p.expr,
		comb.ExpectRune(')'),
		p.quantifier.OPT(),
	).Map(p.m.ToGroup)(in)
}

// Recursive definition
// subexpr_item --> anchor | group | match
func (p *Parser) subexprItem(in comb.Input) (comb.Output, bool) {
	return comb.Parser(p.anchor).ALT(p.group, p.match).Map(p.m.ToSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (p *Parser) subexpr(in comb.Input) (comb.Output, bool) {
	return comb.Parser(p.subexprItem).REP1().Map(p.m.ToSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (p *Parser) expr(in comb.Input) (comb.Output, bool) {
	return comb.Parser(p.subexpr).CONCAT(
		comb.ExpectRune('|').CONCAT(p.expr).OPT(),
	).Map(p.m.ToExpr)(in)
}

// Parse is the topmost parser combinator for parsing a regular expression read from the input.
func (p *Parser) Parse(regex string) (comb.Output, bool) {
	in := newStringInput(regex)
	return p.regex(in)
}
