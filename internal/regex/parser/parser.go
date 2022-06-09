// package parser provides all data types and primitive constructs for building a parser combinator.
package parser

type (
	// empty is the empty string ε.
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
	// input is the input to a parser function.
	input interface {
		// Current returns the current rune from input along with its position in the input.
		Current() (rune, int)
		// Remaining returns the remaining of input. If no input left, it returns nil.
		Remaining() input
	}

	// output is the output of a parser function.
	output struct {
		Result    result
		Remaining input
	}

	// parser is the type for a function that receives a parsing input and returns a parsing output.
	parser func(input) (output, bool)
)

// getAt returns the value of a symbol in the left-side of a production rule.
//
// Example:
//
// • Production Rule: range --> "{" num ( "," num? )? "}"
//
// • in = {2,4}
//
// • getAt(in, 1) = 2, get(in, 3) = 4
//
func getAt(v any, i int) (any, bool) {
	if l, ok := v.(list); ok {
		if 0 <= i && i < len(l) {
			return l[i].Val, true
		}
	}

	return nil, false
}

// expectRune creates a parser that returns a successful result only if the input starts with the given rune.
func expectRune(r rune) parser {
	return func(in input) (output, bool) {
		if in == nil {
			return output{}, false
		}

		if curr, pos := in.Current(); curr == r {
			return output{
				Result:    result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return output{}, false
	}
}

// expectRuneIn creates a parser that returns a successful result only if the input starts with one of the given runes.
func expectRuneIn(runes ...rune) parser {
	return func(in input) (output, bool) {
		if in == nil {
			return output{}, false
		}

		for _, r := range runes {
			if curr, pos := in.Current(); curr == r {
				return output{
					Result:    result{r, pos},
					Remaining: in.Remaining(),
				}, true
			}
		}

		return output{}, false
	}
}

// expectRuneInRange creates a parser that returns a successful result only if the input starts with a rune in the given range.
func expectRuneInRange(low, up rune) parser {
	return func(in input) (output, bool) {
		if in == nil {
			return output{}, false
		}

		if r, pos := in.Current(); low <= r && r <= up {
			return output{
				Result:    result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return output{}, false
	}
}

// expectRunes creates a parser that returns a successful result only if the input starts with the given runes in the given order.
func expectRunes(runes ...rune) parser {
	return func(in input) (output, bool) {
		var pos int

		for i, r := range runes {
			if in == nil {
				return output{}, false
			}

			curr, p := in.Current()
			if curr != r {
				return output{}, false
			}

			// Save only the first position
			if i == 0 {
				pos = p
			}

			in = in.Remaining()
		}

		return output{
			Result:    result{runes, pos},
			Remaining: in,
		}, true
	}
}

// expectString creates a parser that returns a successful result only if the input starts with the given string.
func expectString(s string) parser {
	return func(in input) (output, bool) {
		if out, ok := expectRunes([]rune(s)...)(in); ok {
			return output{
				Result:    result{s, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return output{}, false
	}
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

// CONCAT composes a parser that concats parser p to a sequence of parsers.
// It applies parser p to the input, then applies the next parser to the remaining of the input,
// and continues parsing to the last parser.
//
// • EBNF Operator: Concatenation
//
// • EBNF Notation: p q
func (p parser) CONCAT(q ...parser) parser {
	return func(in input) (output, bool) {
		var l list

		all := append([]parser{p}, q...)
		for _, parse := range all {
			out, ok := parse(in)
			if !ok {
				return output{}, false
			}

			l = append(l, out.Result)
			in = out.Remaining
		}

		return output{
			Result:    result{l, l[0].Pos},
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
	return func(in input) (output, bool) {
		all := append([]parser{p}, q...)
		for _, parse := range all {
			if out, ok := parse(in); ok {
				return out, true
			}
		}

		return output{}, false
	}
}

// OPT composes a parser that applies parser p zero or one time to the input.
// If the parser does not succeed, it will return an empty result.
//
// • EBNF Operator: Optional
//
// • EBNF Notation: [ p ] or p?
func (p parser) OPT() parser {
	return func(in input) (output, bool) {
		if out, ok := p(in); ok {
			return out, true
		}

		return output{
			Result: result{
				// Position for empty string ε is not defined
				Val: empty{},
			},
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
	return func(in input) (output, bool) {
		var l list

		for i := 0; in != nil; i++ {
			out, ok := p(in)
			if !ok {
				break
			}

			l = append(l, out.Result)
			in = out.Remaining
		}

		out := output{
			Remaining: in,
		}

		if len(l) == 0 {
			out.Result = result{
				// Position for empty string ε is not defined
				Val: empty{},
			}
		} else {
			out.Result = result{l, l[0].Pos}
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
	return func(in input) (output, bool) {
		if out, ok := p.REP()(in); ok {
			if res, ok := out.Result.Val.(list); ok && len(res) > 0 {
				return out, true
			}
		}

		return output{}, false
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
	return func(in input) (output, bool) {
		if out, ok := p(in); ok {
			val := flatten(out.Result)
			return output{
				Result:    result{val, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return output{}, false
	}
}

// Select composes a parser that applies parser p to the input and returns a list of symbols in the left-side of the production rule.
// This will not have any effect if the result of parsing is not a list.
// If indices are invalid, you will get an empty string (ε).
func (p parser) Select(i ...int) parser {
	return func(in input) (output, bool) {
		out, ok := p(in)
		if !ok {
			return output{}, false
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
			res = result{
				// Position for empty string ε is not defined
				Val: empty{},
			}
		}

		return output{
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
	return func(in input) (output, bool) {
		out, ok := p(in)
		if !ok {
			return output{}, false
		}

		l, ok := out.Result.Val.(list)
		if !ok {
			return out, true
		}

		var res result
		if 0 <= i && i < len(l) {
			res = l[i]
		} else {
			res = result{
				// Position for empty string ε is not defined
				Val: empty{},
			}
		}

		return output{
			Result:    res,
			Remaining: out.Remaining,
		}, ok
	}
}

// converter is the type for a function that receives a parsing result and returns a new value for the result.
type converter func(any) (any, bool)

// Convert composes a parser that uses parser p to parse the input and applies a converter function to the result of parsing.
// If the parser does not succeed, the converter function will not be applied.
func (p parser) Convert(f converter) parser {
	return func(in input) (output, bool) {
		if out, ok := p(in); ok {
			if val, ok := f(out.Result.Val); ok {
				return output{
					Result:    result{val, out.Result.Pos},
					Remaining: out.Remaining,
				}, true
			}
		}

		return output{}, false
	}
}

// constructor is the type for a function that receives a parsing result and returns a new one.
type constructor func(result) parser

// Bind composes a parser that uses parser p to parse the input and constructs a second parser from the result of parsing.
// It then applies the new parser to the remaining input from the first parser.
// You can use this to implement syntax annotations.
func (p parser) Bind(f constructor) parser {
	return func(in input) (output, bool) {
		out, ok := p(in)
		if !ok {
			return output{}, false
		}

		return f(out.Result)(out.Remaining)
	}
}
