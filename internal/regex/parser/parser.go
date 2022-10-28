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

	escapedChars = []rune{'\\', '/', '|', '.', '?', '*', '+', '(', ')', '[', ']', '{', '}', '$'}
)

//==================================================< MAPPERS >==================================================

type Mappers interface {
	// ToUnescapedChar corresponds to the production rule for all characters excluding the escaped ones.
	ToUnescapedChar(comb.Result) (comb.Result, bool)
	// ToEscapedChar corresponds to escaped_char --> "\" ( "\" | "/" | "|" | "." | "?" | "*" | "+" | "(" | ")" | "[" | "]" | "{" | "}" )
	ToEscapedChar(comb.Result) (comb.Result, bool)
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
	// ToCharRange corresponds to char_range --> char "-" char
	ToCharRange(comb.Result) (comb.Result, bool)
	// ToCharGroupItem corresponds to char_group_item --> char_class | ascii_char_class | char_range | escaped_char | unescaped_char
	ToCharGroupItem(comb.Result) (comb.Result, bool)
	// ToCharGroup corresponds to char_group --> "[" "^"? char_group_item+ "]"
	ToCharGroup(comb.Result) (comb.Result, bool)
	// ToASCIICharClass corresponds to ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	ToASCIICharClass(comb.Result) (comb.Result, bool)
	// ToCharClass corresponds to char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	ToCharClass(comb.Result) (comb.Result, bool)
	// ToAnyChar corresponds to any_char --> "."
	ToAnyChar(comb.Result) (comb.Result, bool)
	// ToMatchItem corresponds to match_item --> any_char | unescaped_char | escaped_char | char_class | ascii_char_class | char_group
	ToMatchItem(comb.Result) (comb.Result, bool)
	// ToMatch corresponds to match --> match_item quantifier?
	ToMatch(comb.Result) (comb.Result, bool)
	// ToAnchor corresponds to anchor --> "$"
	ToAnchor(comb.Result) (comb.Result, bool)
	// ToGroup corresponds to group --> "(" expr ")" quantifier?
	ToGroup(comb.Result) (comb.Result, bool)
	// ToSubexprItem corresponds to subexpr_item --> anchor | group | match
	ToSubexprItem(comb.Result) (comb.Result, bool)
	// ToSubexpr corresponds to subexpr --> subexpr_item+
	ToSubexpr(comb.Result) (comb.Result, bool)
	// ToExpr corresponds to expr --> subexpr ("|" expr)?
	ToExpr(comb.Result) (comb.Result, bool)
	// ToRegex corresponds to regex --> start_of_string? expr
	ToRegex(comb.Result) (comb.Result, bool)
}

func toNum(r comb.Result) (comb.Result, bool) {
	l := r.Val.(comb.List)

	var num int
	for _, r := range l {
		num = num*10 + int(r.Val.(rune)-'0')
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

//==================================================< COMBINATORS >==================================================

// Parser is a parser combinator for regular expressions.
type Parser struct {
	m Mappers

	// Combinators
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

// New creates a parser combinator for regular expressions.
func New(m Mappers) *Parser {
	p := &Parser{
		m: m,
	}

	// char --> /* all valid characters */
	p.char = comb.ExpectRuneInRange(0x20, 0x7E)
	// all characters excluding the escaped ones
	p.unescapedChar = p.char.Bind(comb.ExcludeRunes(escapedChars...)).Map(p.m.ToUnescapedChar)
	// escaped_char --> "\" ( "\" | "/" | "|" | "." | "?" | "*" | "+" | "(" | ")" | "[" | "]" | "{" | "}" | "$" )
	p.escapedChar = comb.ExpectRune('\\').CONCAT(comb.ExpectRuneIn(escapedChars...)).Map(p.m.ToEscapedChar)

	// digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	p.digit = comb.ExpectRuneInRange('0', '9')
	// letter --> "A" | ... | "Z" | "a" | ... | "z"
	p.letter = comb.ExpectRuneInRange('A', 'Z').ALT(comb.ExpectRuneInRange('a', 'z'))

	p.num = p.digit.REP1().Map(toNum)          // num --> digit+
	p.letters = p.letter.REP1().Map(toLetters) // letters --> letter+

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

	// char_range --> char "-" char
	p.charRange = p.char.CONCAT(
		comb.ExpectRune('-'),
		p.char,
	).Map(p.m.ToCharRange)

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

	// char_group_item --> char_class | ascii_char_class | char_range | escaped_char | unescaped_char
	p.charGroupItem = p.charClass.ALT(
		p.asciiCharClass,
		p.charRange,
		p.escapedChar,
		p.unescapedChar,
	).Map(p.m.ToCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	p.charGroup = comb.ExpectRune('[').CONCAT(
		comb.ExpectRune('^').OPT(),
		p.charGroupItem.REP1(),
		comb.ExpectRune(']'),
	).Map(p.m.ToCharGroup)

	// any_char --> "."
	p.anyChar = comb.ExpectRune('.').Map(p.m.ToAnyChar)

	// match_item --> any_char | unescaped_char | escaped_char | char_class | ascii_char_class | char_group
	p.matchItem = p.anyChar.ALT(
		p.unescapedChar,
		p.escapedChar,
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
func (p *Parser) Parse(in comb.Input) (comb.Output, bool) {
	return p.regex(in)
}
