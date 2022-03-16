package regex

type (
	// any is an alias for interface type.
	any = interface{}

	// empty is the empty string (ε).
	empty struct{}

	// result is the result of parsing a production rule.
	// It represents the left-side of a production rule.
	result struct {
		Val any
		Pos int
	}

	// list is the type for the result of concatenation or repetition.
	list []result
)

type (
	// parseInput is the input to a parser function.
	parseInput interface {
		// Current returns the current rune from input along with its position in the input.
		Current() (rune, int)

		// Remaining returns the remaining of input. If no input left, it returns nil.
		Remaining() parseInput
	}

	// parseInput is the output of a parser function.
	parseOutput struct {
		Result    result
		Remaining parseInput
	}

	// parser is the type for a function that receives a parsing input and returns a parsing output.
	parser func(parseInput) (parseOutput, bool)
)

// getVal returns the value of a symbol in the left-side of a production rule from a parsing result.
//
// Example:
//
// • Production Rule: range --> "{" num ( "," num? )? "}"
//
// • Input: "{2,4}"
//
// • getVal(res, 1) = 2, getVal(res, 3) = 4
func getVal(v any, i int) (result, bool) {
	if l, ok := v.(list); ok {
		if 0 <= i && i < len(l) {
			return l[i], true
		}
	}

	return result{}, false
}

// expectRune creates a parser that returns a successful result only if the input starts with the given rune.
func expectRune(r rune) parser {
	return func(in parseInput) (parseOutput, bool) {
		if in == nil {
			return parseOutput{}, false
		}

		if curr, pos := in.Current(); curr == r {
			return parseOutput{
				Result:    result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return parseOutput{}, false
	}
}

// expectRuneIn creates a parser that returns a successful result only if the input starts with one of the given runes.
func expectRuneIn(runes ...rune) parser {
	return func(in parseInput) (parseOutput, bool) {
		if in == nil {
			return parseOutput{}, false
		}

		for _, r := range runes {
			if curr, pos := in.Current(); curr == r {
				return parseOutput{
					Result:    result{r, pos},
					Remaining: in.Remaining(),
				}, true
			}
		}

		return parseOutput{}, false
	}
}

// expectRuneInRange creates a parser that returns a successful result only if the input starts with a rune in the given range.
func expectRuneInRange(low, up rune) parser {
	return func(in parseInput) (parseOutput, bool) {
		if in == nil {
			return parseOutput{}, false
		}

		if r, pos := in.Current(); low <= r && r <= up {
			return parseOutput{
				Result:    result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return parseOutput{}, false
	}
}

// expectRunes creates a parser that returns a successful result only if the input starts with the given runes in the given order.
func expectRunes(runes ...rune) parser {
	return func(in parseInput) (parseOutput, bool) {
		var pos int

		for i, r := range runes {
			if in == nil {
				return parseOutput{}, false
			}

			curr, p := in.Current()
			if curr != r {
				return parseOutput{}, false
			}

			// Save only the first position
			if i == 0 {
				pos = p
			}

			in = in.Remaining()
		}

		return parseOutput{
			Result:    result{runes, pos},
			Remaining: in,
		}, true
	}
}

// expectString creates a parser that returns a successful result only if the input starts with the given string.
func expectString(s string) parser {
	return func(in parseInput) (parseOutput, bool) {
		if out, ok := expectRunes([]rune(s)...)(in); ok {
			return parseOutput{
				Result:    result{s, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return parseOutput{}, false
	}
}

// CONCAT composes a parser that concats parser p to a sequence of parsers.
// It applies parser p to the input, then applies the next parser to the remaining of the input,
// and continues parsing to the last parser.
//
// • EBNF Operator: Concatenation
//
// • EBNF Notation: p q
func (p parser) CONCAT(q ...parser) parser {
	return func(in parseInput) (parseOutput, bool) {
		var res list

		all := append([]parser{p}, q...)
		for _, parse := range all {
			out, ok := parse(in)
			if !ok {
				return parseOutput{}, false
			}

			res = append(res, out.Result)
			in = out.Remaining
		}

		return parseOutput{
			Result:    result{res, res[0].Pos},
			Remaining: in,
		}, true
	}
}

// ALT composes a parser that alternate parser p with a sequence of parsers.
// It applies parser p to the input and if it does not succeed,
// it applies the next parser to the same input, and continues parsing to the last parser.
// It stops at the first successful parsing and returns its result.
//
// • EBNF Operator: Alternation
//
// • EBNF Notation: p | q
func (p parser) ALT(q ...parser) parser {
	return func(in parseInput) (parseOutput, bool) {
		all := append([]parser{p}, q...)
		for _, parse := range all {
			if out, ok := parse(in); ok {
				return out, true
			}
		}

		return parseOutput{}, false
	}
}

// OPT composes a parser that applies parser p zero or one time to the input.
// If the parser does not succeed, it will return an empty result.
//
// • EBNF Operator: Optional
//
// • EBNF Notation: [ p ] or p?
func (p parser) OPT() parser {
	return func(in parseInput) (parseOutput, bool) {
		if out, ok := p(in); ok {
			return out, true
		}

		return parseOutput{
			Result:    result{Val: empty{}}, // Position for empty string (ε) is undefined!
			Remaining: in,
		}, true
	}
}

// REP composes a parser that applies parser p zero or more times to the input and accumulates the results.
// If the parser does not succeed, it will return an empty result.
//
// • EBNF Operator: Repetition (Kleene Star)
//
// • EBNF Notation: { p } or p*
func (p parser) REP() parser {
	return func(in parseInput) (parseOutput, bool) {
		var res list

		for i := 0; in != nil; i++ {
			out, ok := p(in)
			if !ok {
				break
			}

			res = append(res, out.Result)
			in = out.Remaining
		}

		out := parseOutput{
			Remaining: in,
		}

		if len(res) == 0 {
			out.Result = result{Val: empty{}} // Position for empty string (ε) is undefined!
		} else {
			out.Result = result{res, res[0].Pos}
		}

		return out, true
	}
}

// REP1 composes a parser that applies parser p one or more times to the input and accumulates the results.
// This does not allow parsing zero times (empty result).
//
// • EBNF Operator: Kleene Plus
//
// • EBNF Notation: p+
func (p parser) REP1() parser {
	return func(in parseInput) (parseOutput, bool) {
		if out, ok := p.REP()(in); ok {
			if res, ok := out.Result.Val.(list); ok && len(res) > 0 {
				return out, true
			}
		}

		return parseOutput{}, false
	}
}

func flatten(res result) list {
	switch v := res.Val.(type) {
	case empty:
		return list{}

	case list:
		var l list
		for _, v := range v {
			l = append(l, flatten(v)...)
		}
		return l

	default:
		return list{res}
	}
}

// Flatten composes a parser that applies parser p to the input and flattens all results into a single list.
// This can be used for accessing the values of symbols in the left-side of a production rule more intuitively.
func (p parser) Flatten() parser {
	return func(in parseInput) (parseOutput, bool) {
		if out, ok := p(in); ok {
			val := flatten(out.Result)
			return parseOutput{
				Result:    result{val, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return parseOutput{}, false
	}
}

// Select composes a parser that applies parser p to the input and returns a list of symbols in the left-side of the production rule.
// This will not have any effect if the result of parsing is not a list.
// If indices are invalid, you will get an empty string (ε).
func (p parser) Select(i ...int) parser {
	return func(in parseInput) (parseOutput, bool) {
		out, ok := p(in)
		if !ok {
			return parseOutput{}, false
		}

		l, ok := out.Result.Val.(list)
		if !ok {
			return out, true
		}

		var sub list
		for _, j := range i {
			if 0 <= j && j < len(l) {
				sub = append(sub, l[j])
			}
		}

		var res result
		if len(sub) > 0 {
			res = result{sub, sub[0].Pos}
		} else {
			res = result{Val: empty{}} // Position for empty string (ε) is undefined!
		}

		return parseOutput{
			Result:    res,
			Remaining: out.Remaining,
		}, true
	}
}

// Get composes a parser that applies parser p to the input and returns the value of a symbol in the left-side of the production rule.
// This can be used after CONCAT, REP, REP1, Flatten, and/or Select.
// It will not have any effect if used after other operators and the result of parsing is not a list.
// If index is invalid, you will get an empty string (ε).
func (p parser) Get(i int) parser {
	return func(in parseInput) (parseOutput, bool) {
		out, ok := p(in)
		if !ok {
			return parseOutput{}, false
		}

		l, ok := out.Result.Val.(list)
		if !ok {
			return out, true
		}

		var res result
		if 0 <= i && i < len(l) {
			res = l[i]
		} else {
			res = result{Val: empty{}} // Position for empty string (ε) is undefined!
		}

		return parseOutput{
			Result:    res,
			Remaining: out.Remaining,
		}, ok
	}
}

// converter is the type for a function that receives a parsing result and returns a new value for the result.
type converter func(result) (any, bool)

// Convert composes a parser that uses parser p to parse the input and applies a converter function to the result of parsing.
// If the parser does not succeed, the converter function will not be applied.
func (p parser) Convert(f converter) parser {
	return func(in parseInput) (parseOutput, bool) {
		if out, ok := p(in); ok {
			if val, ok := f(out.Result); ok {
				return parseOutput{
					Result:    result{val, out.Result.Pos},
					Remaining: out.Remaining,
				}, true
			}
		}

		return parseOutput{}, false
	}
}

// converter is the type for a function that receives a parsing result and returns a new one.
type constructor func(result) parser

// excludeChars can be bound on a _char parser to exclude certain runes.
func excludeChars(rs ...rune) constructor {
	return func(res result) parser {
		return func(in parseInput) (parseOutput, bool) {
			if ch, ok := res.Val.(Char); ok {
				for _, r := range rs {
					if ch.Val == r {
						return parseOutput{}, false
					}
				}
			}

			return parseOutput{
				Result:    res,
				Remaining: in,
			}, true
		}
	}
}

// Bind composes a parser that uses parser p to parse the input and constructs a second parser from the result of parsing.
// It then applies the new parser to the remaining input from the first parser.
// You can use this to implement syntax annotations.
func (p parser) Bind(f constructor) parser {
	return func(in parseInput) (parseOutput, bool) {
		out, ok := p(in)
		if !ok {
			return parseOutput{}, false
		}

		q := f(out.Result)
		return q(out.Remaining)
	}
}

// _empty is a no-op parser that returns the empty string (ε).
//
// empty --> ε
var _empty = func(in parseInput) (parseOutput, bool) {
	return parseOutput{
		Result: result{
			Val: empty{},
		},
		Remaining: in,
	}, true
}

// GRAMMAR
var (
	_space   parser = expectRuneIn(' ')                                            // space --> " "
	_char    parser = expectRuneInRange(0x20, 0x7E).Convert(toChar)                // char --> /* all valid characters */
	_digit   parser = expectRuneInRange('0', '9')                                  // digit --> "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
	_letter  parser = expectRuneInRange('A', 'Z').ALT(expectRuneInRange('a', 'z')) // letter --> "A" | ... | "Z" | "a" | ... | "z"
	_num     parser = _digit.REP1().Convert(toNum)                                 // num --> digit+
	_letters parser = _letter.REP1().Convert(toLetters)                            // letters --> letter+

	_zeroOrOne  parser = expectRune('?').Convert(toCardNode) // zero_or_one --> "?"
	_zeroOrMore parser = expectRune('*').Convert(toCardNode) // zero_or_more --> "*"
	_oneOrMore  parser = expectRune('+').Convert(toCardNode) // one_or_more --> "+"

	// upper_bound --> "," " "? num?
	_upperBound parser = expectRune(',').CONCAT(
		_space.OPT(),
		_num.OPT(),
	).Convert(toUpperBound)

	// range --> "{" num upper_bound? "}"
	_range parser = expectRune('{').CONCAT(
		_num,
		_upperBound.OPT(),
		expectRune('}'),
	).Convert(toRange)

	// cardinality --> zero_or_one | zero_or_more | one_or_more | range
	_cardinality parser = _zeroOrOne.ALT(_zeroOrMore, _oneOrMore, _range).Convert(toCardinality)

	// quantifier --> cardinality "?"?
	_quantifier parser = _cardinality.CONCAT(
		expectRune('?').OPT(),
	).Convert(toQuantifier)

	// char_range --> char ("-" char)?
	_charRange parser = _char.CONCAT(
		expectRune('-'),
		_char,
	).Convert(toCharRange)

	// char_group_item -->  char_class | ascii_char_class | char_range | char /* excluding ] */
	_charGroupItem parser = _charClass.ALT(
		_asciiCharClass,
		_charRange,
		_char.Bind(excludeChars(']')),
	).Convert(toCharGroupItem)

	// char_group --> "[" "^"? char_group_item+ "]"
	_charGroup parser = expectRune('[').CONCAT(
		expectRune('^').OPT(),
		_charGroupItem.REP1(),
		expectRune(']'),
	).Convert(toCharGroup)

	// char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
	_charClass parser = expectString(`\d`).ALT(
		expectString(`\D`),
		expectString(`\s`), expectString(`\S`),
		expectString(`\w`), expectString(`\W`),
	).Convert(toCharClass)

	// ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
	_asciiCharClass parser = expectString("[:blank:]").ALT(
		expectString("[:space:]"),
		expectString("[:digit:]"), expectString("[:xdigit:]"),
		expectString("[:upper:]"), expectString("[:lower:]"),
		expectString("[:alpha:]"), expectString("[:alnum:]"),
		expectString("[:word:]"), expectString("[:ascii:]"),
	).Convert(toASCIICharClass)

	// any_char --> "."
	_anyChar parser = expectRune('.').Convert(toAnyChar)

	// match_item --> any_char | char_class | ascii_char_class | char_group | char /* excluding | ) */
	_matchItem parser = _anyChar.ALT(
		_charClass,
		_asciiCharClass,
		_charGroup,
		_char.Bind(excludeChars('|', ')')),
	).Convert(toMatchItem)

	// match --> match_item quantifier?
	_match parser = _matchItem.CONCAT(_quantifier.OPT()).Convert(toMatch)

	// backref --> "\" num
	_backref parser = expectRune('\\').CONCAT(_num).Convert(toBackref)

	// anchor --> "$" | "\b" | "\B"
	_anchor parser = expectString("$").ALT(
		expectString("\\b"),
		expectString("\\B"),
	).Convert(toAnchor)

	// regex --> "^"? expr
	_regex parser = expectRune('^').OPT().CONCAT(_expr).Convert(toRegex)
)

// group --> "(" "?:"? expr ")" quantifier?
func _group(in parseInput) (parseOutput, bool) {
	return expectRune('(').CONCAT(
		expectString("?:").OPT(),
		_expr,
		expectRune(')'),
		_quantifier.OPT(),
	).Convert(toGroup)(in)
}

// subexpr_item --> group | anchor | backref | match
func _subexprItem(in parseInput) (parseOutput, bool) {
	return parser(_group).ALT(_anchor, _backref, _match).Convert(toSubexprItem)(in)
}

// subexpr --> subexpr_item+
func _subexpr(in parseInput) (parseOutput, bool) {
	return parser(_subexprItem).REP1().Convert(toSubexpr)(in)
}

// expr --> subexpr ("|" expr)?
func _expr(in parseInput) (parseOutput, bool) {
	return parser(_subexpr).CONCAT(
		expectRune('|').CONCAT(_expr).OPT(),
	).Convert(toExpr)(in)
}

var (
	toChar converter = func(res result) (any, bool) {
		r, ok := res.Val.(rune)
		if !ok {
			return nil, false
		}

		return Char{
			TokPos: res.Pos,
			Val:    r,
		}, true
	}

	toNum converter = func(res result) (any, bool) {
		l, ok := res.Val.(list)
		if !ok {
			return nil, false
		}

		var num int
		for _, r := range l {
			num = num*10 + int(r.Val.(rune)-'0')
		}

		return Num{res.Pos, num}, true
	}

	toLetters converter = func(res result) (any, bool) {
		l, ok := res.Val.(list)
		if !ok {
			return nil, false
		}

		var s string
		for _, r := range l {
			s += string(r.Val.(rune))
		}

		return Letters{res.Pos, s}, true
	}

	toCardNode converter = func(res result) (any, bool) {
		r, ok := res.Val.(rune)
		if !ok {
			return nil, false
		}

		switch r {
		case '?':
			return ZeroOrOne{res.Pos}, true
		case '*':
			return ZeroOrMore{res.Pos}, true
		case '+':
			return OneOrMore{res.Pos}, true
		default:
			return nil, false
		}
	}

	toUpperBound converter = func(res result) (any, bool) {
		r, ok := getVal(res.Val, 2)
		if !ok {
			return nil, false
		}

		var val *Num
		switch v := r.Val.(type) {
		case empty:
		case Num:
			val = &v
		}

		return UpperBound{
			CommaPos: r.Pos,
			Val:      val,
		}, true
	}

	toRange converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		r2, ok := getVal(res.Val, 2)
		if !ok {
			return nil, false
		}

		low, ok := r1.Val.(Num)
		if !ok {
			return nil, false
		}

		var up *UpperBound
		switch v := r2.Val.(type) {
		case empty:
		case UpperBound:
			up = &v
		}

		return Range{
			OpenPos: r0.Pos,
			Low:     low,
			Up:      up,
		}, true
	}

	toCardinality converter = func(res result) (any, bool) {
		switch v := res.Val.(type) {
		case ZeroOrOne:
			return &v, true
		case ZeroOrMore:
			return &v, true
		case OneOrMore:
			return &v, true
		case Range:
			return &v, true
		default:
			return nil, false
		}
	}

	toQuantifier converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		card, ok := r0.Val.(Cardinality)
		if !ok {
			return nil, false
		}

		// Check whether or not lazy modifier is present
		r, ok := r1.Val.(rune)
		lazy := ok && r == '?'

		return Quantifier{
			Card: card,
			Lazy: lazy,
		}, true
	}

	toCharRange converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r2, ok := getVal(res.Val, 2)
		if !ok {
			return nil, false
		}

		low, ok := r0.Val.(Char)
		if !ok {
			return nil, false
		}

		up, ok := r2.Val.(Char)
		if !ok {
			return nil, false
		}

		return CharRange{
			Low: low,
			Up:  up,
		}, true
	}

	toCharGroupItem converter = func(res result) (any, bool) {
		switch v := res.Val.(type) {
		case CharClass:
			return &v, true
		case ASCIICharClass:
			return &v, true
		case CharRange:
			return &v, true
		case Char:
			return &v, true
		default:
			return nil, false
		}
	}

	toCharGroup converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		r2, ok := getVal(res.Val, 2)
		if !ok {
			return nil, false
		}

		// Check whether or not the negation modifier is present
		r, ok := r1.Val.(rune)
		neg := ok && r == '^'

		l, ok := r2.Val.(list)
		if !ok {
			return nil, false
		}

		var items []CharGroupItem
		for _, r := range l {
			if i, ok := r.Val.(CharGroupItem); ok {
				items = append(items, i)
			}
		}

		return CharGroup{
			OpenPos: r0.Pos,
			Neg:     neg,
			Items:   items,
		}, true
	}

	toCharClass converter = func(res result) (any, bool) {
		val, ok := res.Val.(string)
		if !ok {
			return nil, false
		}

		var tok Token

		switch val {
		case `\d`:
			tok = DIGIT_CHARS
		case `\D`:
			tok = NOT_DIGIT_CHARS
		case `\s`:
			tok = WHITESPACE
		case `\S`:
			tok = NOT_WHITESPACE
		case `\w`:
			tok = WORD_CHARS
		case `\W`:
			tok = NOT_WORD_CHARS
		default:
			return nil, false
		}

		return CharClass{
			TokPos: res.Pos,
			Tok:    tok,
		}, true
	}

	toASCIICharClass converter = func(res result) (any, bool) {
		val, ok := res.Val.(string)
		if !ok {
			return nil, false
		}

		var tok Token

		switch val {
		case "[:blank:]":
			tok = BLANK
		case "[:space:]":
			tok = SPACE
		case "[:digit:]":
			tok = DIGIT
		case "[:xdigit:]":
			tok = XDIGIT
		case "[:upper:]":
			tok = UPPER
		case "[:lower:]":
			tok = LOWER
		case "[:alpha:]":
			tok = ALPHA
		case "[:alnum:]":
			tok = ALNUM
		case "[:word:]":
			tok = WORD
		case "[:ascii:]":
			tok = ASCII
		default:
			return nil, false
		}

		return ASCIICharClass{
			TokPos: res.Pos,
			Tok:    tok,
		}, true
	}

	toAnyChar converter = func(res result) (any, bool) {
		if r, ok := res.Val.(rune); ok && r == '.' {
			return AnyChar{res.Pos}, true
		}

		return nil, false
	}

	toMatchItem converter = func(res result) (any, bool) {
		switch v := res.Val.(type) {
		case AnyChar:
			return &v, true
		case CharGroup:
			return &v, true
		case CharClass:
			return &v, true
		case ASCIICharClass:
			return &v, true
		case Char:
			return &v, true
		default:
			return nil, false
		}
	}

	toMatch converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		item, ok := r0.Val.(MatchItem)
		if !ok {
			return nil, false
		}

		var quant *Quantifier
		if v, ok := r1.Val.(Quantifier); ok {
			quant = &v
		}

		return Match{
			Item:  item,
			Quant: quant,
		}, true
	}

	toBackref converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		num, ok := r1.Val.(Num)
		if !ok {
			return nil, false
		}

		return Backref{
			SlashPos: r0.Pos,
			Ref:      num,
		}, true
	}

	toAnchor converter = func(res result) (any, bool) {
		r, ok := res.Val.(string)
		if !ok {
			return nil, false
		}

		switch r {
		case "$":
			return Anchor{res.Pos, END_OF_STRING}, true
		case "\\b":
			return Anchor{res.Pos, WORD_BOUNDARY}, true
		case "\\B":
			return Anchor{res.Pos, NOT_WORD_BOUNDARY}, true
		default:
			return nil, false
		}
	}

	toGroup converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		r2, ok := getVal(res.Val, 2)
		if !ok {
			return nil, false
		}

		r4, ok := getVal(res.Val, 4)
		if !ok {
			return nil, false
		}

		// Check whether or not the group non-capturing modifier is present
		s, ok := r1.Val.(string)
		nonCap := ok && s == "?:"

		expr, ok := r2.Val.(Expr)
		if !ok {
			return nil, false
		}

		var quant *Quantifier
		if v, ok := r4.Val.(Quantifier); ok {
			quant = &v
		}

		return Group{
			OpenPos: r0.Pos,
			NonCap:  nonCap,
			Expr:    expr,
			Quant:   quant,
		}, true
	}

	toSubexprItem converter = func(res result) (any, bool) {
		switch v := res.Val.(type) {
		case Group:
			return &v, true
		case Anchor:
			return &v, true
		case Backref:
			return &v, true
		case Match:
			return &v, true
		default:
			return nil, false
		}
	}

	toSubexpr converter = func(res result) (any, bool) {
		l, ok := res.Val.(list)
		if !ok {
			return nil, false
		}

		var items []SubexprItem
		for _, r := range l {
			if item, ok := r.Val.(SubexprItem); ok {
				items = append(items, item)
			}
		}

		return Subexpr{
			Items: items,
		}, true
	}

	toExpr converter = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		sub, ok := r0.Val.(Subexpr)
		if !ok {
			return nil, false
		}

		var expr *Expr

		if _, ok := r1.Val.(empty); !ok {
			r11, ok := getVal(r1.Val, 1)
			if !ok {
				return nil, false
			}

			v, ok := r11.Val.(Expr)
			if !ok {
				return nil, false
			}

			expr = &v
		}

		return Expr{
			Sub:  sub,
			Expr: expr,
		}, true
	}

	toRegex = func(res result) (any, bool) {
		r0, ok := getVal(res.Val, 0)
		if !ok {
			return nil, false
		}

		r1, ok := getVal(res.Val, 1)
		if !ok {
			return nil, false
		}

		// Check whether or not the start-of-string anchor is present
		r, ok := r0.Val.(rune)
		begin := ok && r == '^'

		expr, ok := r1.Val.(Expr)
		if !ok {
			return nil, false
		}

		return Regex{
			Begin: begin,
			Expr:  expr,
		}, true
	}
)
