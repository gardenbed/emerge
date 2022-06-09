package parser

// converters are used for converting a parse result into a node in an abstract syntax tree and pass it up the parser chain.
// They can be used for constructing an abstract syntax tree in a bottom-up manner.
type converters interface {
	ToChar(res result) (any, bool)
	ToNum(res result) (any, bool)
	ToLetters(res result) (any, bool)
	ToRepOp(res result) (any, bool)
	ToUpperBound(res result) (any, bool)
	ToRange(res result) (any, bool)
	ToRepetition(res result) (any, bool)
	ToQuantifier(res result) (any, bool)
	ToCharRange(res result) (any, bool)
	ToCharGroupItem(res result) (any, bool)
	ToCharGroup(res result) (any, bool)
	ToCharClass(res result) (any, bool)
	ToASCIICharClass(res result) (any, bool)
	ToAnyChar(res result) (any, bool)
	ToMatchItem(res result) (any, bool)
	ToMatch(res result) (any, bool)
	ToBackref(res result) (any, bool)
	ToAnchor(res result) (any, bool)
	ToGroup(res result) (any, bool)
	ToSubexprItem(res result) (any, bool)
	ToSubexpr(res result) (any, bool)
	ToExpr(res result) (any, bool)
	ToRegex(res result) (any, bool)
}

// regex is a parser combinator for parsing regular expressions.
type regex struct {
	c converters

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

// NewRegex creates a parser combinator for parsing regular expressions.
func NewRegex(c converters) *regex {
	r := &regex{
		c: c,
	}

	r.char = expectRuneInRange(0x20, 0x7E).Convert(r.c.ToChar)                                     // char --> /* all valid characters */
	r.charInGroup = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(']')).Convert(r.c.ToChar)      // char --> /* all valid characters except ] */
	r.charInMatch = expectRuneInRange(0x20, 0x7E).Bind(excludeRunes(')', '|')).Convert(r.c.ToChar) // char --> /* all valid characters except ) and | */
	r.digit = expectRuneInRange('0', '9')                                                          // digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	r.letter = expectRuneInRange('A', 'Z').ALT(expectRuneInRange('a', 'z'))                        // letter --> "A" | ... | "Z" | "a" | ... | "z"
	r.num = r.digit.REP1().Convert(r.c.ToNum)                                                      // num --> digit+
	r.letters = r.letter.REP1().Convert(r.c.ToLetters)                                             // letters --> letter+

	// rep_op --> "?" | "*" | "+"
	r.repOp = expectRune('?').ALT(
		expectRune('*'),
		expectRune('+'),
	).Convert(r.c.ToRepOp)

	// upper_bound --> "," num?
	r.upperBound = expectRune(',').CONCAT(
		r.num.OPT(),
	).Convert(r.c.ToUpperBound)

	// range --> "{" num upper_bound? "}"
	r.range_ = expectRune('{').CONCAT(
		r.num,
		r.upperBound.OPT(),
		expectRune('}'),
	).Convert(r.c.ToRange)

	// repetition --> rep_op | range
	r.repetition = r.repOp.ALT(
		r.range_,
	).Convert(r.c.ToRepetition)

	// quantifier --> repetition lazy_modifier?
	r.quantifier = r.repetition.CONCAT(
		expectRune('?').OPT(),
	).Convert(r.c.ToQuantifier)

	// char_range --> char "-" char
	r.charRange = r.char.CONCAT(
		expectRune('-'),
		r.char,
	).Convert(r.c.ToCharRange)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	r.charClass = expectString(`\d`).ALT(
		expectString(`\D`),
		expectString(`\s`), expectString(`\S`),
		expectString(`\w`), expectString(`\W`),
	).Convert(r.c.ToCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	r.asciiCharClass = expectString("[:blank:]").ALT(
		expectString("[:space:]"),
		expectString("[:digit:]"), expectString("[:xdigit:]"),
		expectString("[:upper:]"), expectString("[:lower:]"),
		expectString("[:alpha:]"), expectString("[:alnum:]"),
		expectString("[:word:]"), expectString("[:ascii:]"),
	).Convert(r.c.ToASCIICharClass)

	// char_group_item -->  char_class | ascii_char_class | char_range | char /* excluding ] */
	r.charGroupItem = r.charClass.ALT(
		r.asciiCharClass,
		r.charRange,
		r.charInGroup,
	).Convert(r.c.ToCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	r.charGroup = expectRune('[').CONCAT(
		expectRune('^').OPT(),
		r.charGroupItem.REP1(),
		expectRune(']'),
	).Convert(r.c.ToCharGroup)

	// any_char --> "."
	r.anyChar = expectRune('.').Convert(r.c.ToAnyChar)

	// match_item --> any_char | char_class | ascii_char_class | char_group | char /* excluding | ) */
	r.matchItem = r.anyChar.ALT(
		r.charClass,
		r.asciiCharClass,
		r.charGroup,
		r.charInMatch,
	).Convert(r.c.ToMatchItem)

	// match --> match_item quantifier?
	r.match = r.matchItem.CONCAT(r.quantifier.OPT()).Convert(r.c.ToMatch)

	r.backref = expectRune('\\').CONCAT(r.num).Convert(r.c.ToBackref) // backref --> "\" num
	r.anchor = expectRune('$').Convert(r.c.ToAnchor)                  // anchor --> "$"

	// regex --> start_of_string? expr
	r.regex = expectRune('^').OPT().CONCAT(r.expr).Convert(r.c.ToRegex)

	return r
}

// Recursive definition
// group --> "(" expr ")" quantifier?
func (r *regex) group(in input) (output, bool) {
	return expectRune('(').CONCAT(
		r.expr,
		expectRune(')'),
		r.quantifier.OPT(),
	).Convert(r.c.ToGroup)(in)
}

// Recursive definition
// subexpr_item --> group | anchor | backref | match
func (r *regex) subexprItem(in input) (output, bool) {
	return parser(r.group).ALT(r.anchor, r.backref, r.match).Convert(r.c.ToSubexprItem)(in)
}

// Recursive definition
// subexpr --> subexpr_item+
func (r *regex) subexpr(in input) (output, bool) {
	return parser(r.subexprItem).REP1().Convert(r.c.ToSubexpr)(in)
}

// Recursive definition
// expr --> subexpr ("|" expr)?
func (r *regex) expr(in input) (output, bool) {
	return parser(r.subexpr).CONCAT(
		expectRune('|').CONCAT(r.expr).OPT(),
	).Convert(r.c.ToExpr)(in)
}

// excludeRunes can be bound on a rune parser to exclude certain runes.
func excludeRunes(r ...rune) constructor {
	return func(res result) parser {
		return func(in input) (output, bool) {
			if a, ok := res.Val.(rune); ok {
				for _, b := range r {
					if a == b {
						return output{}, false
					}
				}
			}

			return output{
				Result:    res,
				Remaining: in,
			}, true
		}
	}
}
