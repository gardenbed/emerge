package parser

import comb "github.com/gardenbed/emerge/internal/combinator"

var (
	// All characters from 0x00 to 0x7f
	Alphabet = []rune{
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

	escapedChars = []rune{'\\', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}', '$'}
)

//==================================================< MAPPERS >==================================================

type Mappers interface {
	// ToAnyChar corresponds to any_char --> "."
	ToAnyChar(comb.Result) (comb.Result, bool)
	// ToSingleChar corresponds to single_char --> unicode_char | ascii_char | escaped_char | unescaped_char
	ToSingleChar(comb.Result) (comb.Result, bool)
	// ToCharClass corresponds to char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	ToCharClass(comb.Result) (comb.Result, bool)
	// ToASCIICharClass corresponds to ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	ToASCIICharClass(comb.Result) (comb.Result, bool)
	// ToRepOp corresponds to rep_op --> "?" | "*" | "+"
	ToRepOp(comb.Result) (comb.Result, bool)
	// ToUpperBound corresponds to upper_bound --> "," num?
	ToUpperBound(comb.Result) (comb.Result, bool)
	// ToRange corresponds to range --> "{" num upper_bound? "}"
	ToRange(comb.Result) (comb.Result, bool)
	// ToRepetition corresponds to repetition --> rep_op | range
	ToRepetition(comb.Result) (comb.Result, bool)
	// ToQuantifier corresponds to quantifier --> repetition lazy_modifier?
	ToQuantifier(comb.Result) (comb.Result, bool)
	// ToCharInRange corresponds to char_in_range --> unicode_char | ascii_char | char
	ToCharInRange(comb.Result) (comb.Result, bool)
	// ToCharRange corresponds to char_range --> char_in_range "-" char_in_range
	ToCharRange(comb.Result) (comb.Result, bool)
	// ToCharGroupItem corresponds to char_group_item --> char_class | ascii_char_class | char_range | single_char
	ToCharGroupItem(comb.Result) (comb.Result, bool)
	// ToCharGroup corresponds to char_group --> "[" "^"? char_group_item+ "]"
	ToCharGroup(comb.Result) (comb.Result, bool)
	// ToMatchItem corresponds to match_item --> any_char | single_char | char_class | ascii_char_class | char_group
	ToMatchItem(comb.Result) (comb.Result, bool)
	// ToMatch corresponds to match --> match_item quantifier?
	ToMatch(comb.Result) (comb.Result, bool)
	// ToGroup corresponds to group --> "(" expr ")" quantifier?
	ToGroup(comb.Result) (comb.Result, bool)
	// ToAnchor corresponds to anchor --> "$"
	ToAnchor(comb.Result) (comb.Result, bool)
	// ToSubexprItem corresponds to subexpr_item --> anchor | group | match
	ToSubexprItem(comb.Result) (comb.Result, bool)
	// ToSubexpr corresponds to subexpr --> subexpr_item+
	ToSubexpr(comb.Result) (comb.Result, bool)
	// ToExpr corresponds to expr --> subexpr ("|" expr)?
	ToExpr(comb.Result) (comb.Result, bool)
	// ToRegex corresponds to regex --> start_of_string? expr
	ToRegex(comb.Result) (comb.Result, bool)
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

//==================================================< COMBINATORS >==================================================

// Parser is a parser combinator for regular expressions.
type Parser struct {
	m Mappers

	// Combinators
	digit          comb.Parser
	hexDigit       comb.Parser
	letter         comb.Parser
	num            comb.Parser
	letters        comb.Parser
	char           comb.Parser
	unescapedChar  comb.Parser
	escapedChar    comb.Parser
	asciiChar      comb.Parser
	unicodeChar    comb.Parser
	anyChar        comb.Parser
	singleChar     comb.Parser
	charClass      comb.Parser
	asciiCharClass comb.Parser
	repOp          comb.Parser
	upperBound     comb.Parser
	range_         comb.Parser
	repetition     comb.Parser
	quantifier     comb.Parser
	charInRange    comb.Parser
	charRange      comb.Parser
	charGroupItem  comb.Parser
	charGroup      comb.Parser
	matchItem      comb.Parser
	match          comb.Parser
	anchor         comb.Parser
	regex          comb.Parser
}

// New creates a parser combinator for regular expressions.
func New(m Mappers) *Parser {
	p := &Parser{
		m: m,
	}

	// digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	p.digit = comb.ExpectRuneInRange('0', '9').Map(toDigit)
	// hex_digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" | "A" | "B" | "C" | "D" | "E" | "F"
	p.hexDigit = comb.ExpectRuneInRange('0', '9').ALT(comb.ExpectRuneInRange('A', 'F')).Map(toHexDigit)
	// letter --> "A" | ... | "Z" | "a" | ... | "z"
	p.letter = comb.ExpectRuneInRange('A', 'Z').ALT(comb.ExpectRuneInRange('a', 'z'))

	p.num = p.digit.REP1().Map(toNum)          // num --> digit+
	p.letters = p.letter.REP1().Map(toLetters) // letters --> letter+

	// char --> all valid characters
	p.char = comb.ExpectRuneInRange(0x20, 0x7E)
	// unescaped_char --> all characters excluding the escaped ones
	p.unescapedChar = p.char.Bind(comb.ExcludeRunes(escapedChars...))
	// escaped_char --> "\" ( "\" | "|" | "." | "?" | "*" | "+" | "(" | ")" | "[" | "]" | "{" | "}" | "$" )
	p.escapedChar = comb.ExpectRune('\\').CONCAT(comb.ExpectRuneIn(escapedChars...)).Map(toEscapedChar)

	// ascii_char --> "\x" hex_digit{2}
	p.asciiChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit).Map(toASCIIChar)

	// unicode_char --> "\x" hex_digit{4,8}
	p.unicodeChar = comb.ExpectString(`\x`).CONCAT(p.hexDigit, p.hexDigit, p.hexDigit, p.hexDigit,
		p.hexDigit.OPT(),
		p.hexDigit.OPT(),
		p.hexDigit.OPT(),
		p.hexDigit.OPT(),
	).Map(toUnicodeChar)

	// any_char --> "."
	p.anyChar = comb.ExpectRune('.').Map(p.m.ToAnyChar)

	// single_char --> unicode_char | ascii_char | escaped_char | unescaped_char
	p.singleChar = p.unicodeChar.ALT(p.asciiChar, p.escapedChar, p.unescapedChar).Map(p.m.ToSingleChar)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	p.charClass = comb.ExpectString(`\d`).ALT(
		comb.ExpectString(`\D`),
		comb.ExpectString(`\s`), comb.ExpectString(`\S`),
		comb.ExpectString(`\w`), comb.ExpectString(`\W`),
	).Map(p.m.ToCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	p.asciiCharClass = comb.ExpectString("[:blank:]").ALT(
		comb.ExpectString("[:space:]"),
		comb.ExpectString("[:digit:]"), comb.ExpectString("[:xdigit:]"),
		comb.ExpectString("[:upper:]"), comb.ExpectString("[:lower:]"),
		comb.ExpectString("[:alpha:]"), comb.ExpectString("[:alnum:]"),
		comb.ExpectString("[:word:]"), comb.ExpectString("[:ascii:]"),
	).Map(p.m.ToASCIICharClass)

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
	p.charRange = p.charInRange.CONCAT(
		comb.ExpectRune('-'),
		p.charInRange,
	).Map(p.m.ToCharRange)

	// char_group_item --> char_class | ascii_char_class | char_range | single_char
	p.charGroupItem = p.asciiCharClass.ALT(
		p.charClass,
		p.charRange,
		p.singleChar,
	).Map(p.m.ToCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	p.charGroup = comb.ExpectRune('[').CONCAT(
		comb.ExpectRune('^').OPT(),
		p.charGroupItem.REP1(),
		comb.ExpectRune(']'),
	).Map(p.m.ToCharGroup)

	// match_item --> any_char | single_char | char_class | ascii_char_class | char_group
	p.matchItem = p.anyChar.ALT(
		p.singleChar,
		p.charClass,
		p.asciiCharClass,
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
