// package combinator provides all data types and primitive constructs for building a parser combinator.
package combinator

type (
	// Empty is the empty string ε.
	Empty struct{}

	// Result is the result of parsing a production rule.
	// It represents a production rule result.
	Result struct {
		Val any
		Pos int
	}

	// List is the type for the result of concatenation or repetition.
	List []Result
)

type (
	// Input is the input to a parser function.
	Input interface {
		// Current returns the current rune from input along with its position in the input.
		Current() (rune, int)
		// Remaining returns the remaining of input. If no input left, it returns nil.
		Remaining() Input
	}

	// Output is the output of a parser function.
	Output struct {
		Result    Result
		Remaining Input
	}

	// Parser is the type for a function that receives a parsing input and returns a parsing output.
	Parser func(Input) (Output, bool)
)

type (
	// mapper is the type for a function that receives a parsing result and returns a new value for the result.
	Mapper func(any) (any, bool)

	// binder is the type for a function that receives a parsing result and returns a new parser.
	Binder func(Result) Parser
)

// GetAt returns the value of a symbol from the right-side of a production rule.
//
// Example:
//
// • Production Rule: range --> "{" num ( "," num? )? "}"
//
// • in = {2,4}
//
// • GetAt(in, 1) = 2, get(in, 3) = 4
func GetAt(v any, i int) (any, bool) {
	if l, ok := v.(List); ok {
		if 0 <= i && i < len(l) {
			return l[i].Val, true
		}
	}

	return nil, false
}

// ExpectRune creates a parser that returns a successful result only if the input starts with the given rune.
func ExpectRune(r rune) Parser {
	return func(in Input) (Output, bool) {
		if in == nil {
			return Output{}, false
		}

		if curr, pos := in.Current(); curr == r {
			return Output{
				Result:    Result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return Output{}, false
	}
}

// ExpectRuneIn creates a parser that returns a successful result only if the input starts with one of the given runes.
func ExpectRuneIn(runes ...rune) Parser {
	return func(in Input) (Output, bool) {
		if in == nil {
			return Output{}, false
		}

		for _, r := range runes {
			if curr, pos := in.Current(); curr == r {
				return Output{
					Result:    Result{r, pos},
					Remaining: in.Remaining(),
				}, true
			}
		}

		return Output{}, false
	}
}

// ExpectRuneInRange creates a parser that returns a successful result only if the input starts with a rune in the given range.
func ExpectRuneInRange(low, up rune) Parser {
	return func(in Input) (Output, bool) {
		if in == nil {
			return Output{}, false
		}

		if r, pos := in.Current(); low <= r && r <= up {
			return Output{
				Result:    Result{r, pos},
				Remaining: in.Remaining(),
			}, true
		}

		return Output{}, false
	}
}

// ExpectRunes creates a parser that returns a successful result only if the input starts with the given runes in the given order.
func ExpectRunes(runes ...rune) Parser {
	return func(in Input) (Output, bool) {
		var pos int

		for i, r := range runes {
			if in == nil {
				return Output{}, false
			}

			curr, p := in.Current()
			if curr != r {
				return Output{}, false
			}

			// Save only the first position
			if i == 0 {
				pos = p
			}

			in = in.Remaining()
		}

		return Output{
			Result:    Result{runes, pos},
			Remaining: in,
		}, true
	}
}

// ExpectString creates a parser that returns a successful result only if the input starts with the given string.
func ExpectString(s string) Parser {
	return func(in Input) (Output, bool) {
		if out, ok := ExpectRunes([]rune(s)...)(in); ok {
			return Output{
				Result:    Result{s, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return Output{}, false
	}
}

// ExcludeRunes can be bound on a rune parser to exclude certain runes.
func ExcludeRunes(r ...rune) Binder {
	return func(res Result) Parser {
		return func(in Input) (Output, bool) {
			if a, ok := res.Val.(rune); ok {
				for _, b := range r {
					if a == b {
						return Output{}, false
					}
				}
			}

			return Output{
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
func (p Parser) CONCAT(q ...Parser) Parser {
	return func(in Input) (Output, bool) {
		var l List

		all := append([]Parser{p}, q...)
		for _, parse := range all {
			out, ok := parse(in)
			if !ok {
				return Output{}, false
			}

			l = append(l, out.Result)
			in = out.Remaining
		}

		return Output{
			Result:    Result{l, l[0].Pos},
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
func (p Parser) ALT(q ...Parser) Parser {
	return func(in Input) (Output, bool) {
		all := append([]Parser{p}, q...)
		for _, parse := range all {
			if out, ok := parse(in); ok {
				return out, true
			}
		}

		return Output{}, false
	}
}

// OPT composes a parser that applies parser p zero or one time to the input.
// If the parser does not succeed, it will return an empty result.
//
// • EBNF Operator: Optional
//
// • EBNF Notation: [ p ] or p?
func (p Parser) OPT() Parser {
	return func(in Input) (Output, bool) {
		if out, ok := p(in); ok {
			return out, true
		}

		return Output{
			Result: Result{
				// Position for empty string ε is not defined
				Val: Empty{},
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
func (p Parser) REP() Parser {
	return func(in Input) (Output, bool) {
		var l List

		for i := 0; in != nil; i++ {
			out, ok := p(in)
			if !ok {
				break
			}

			l = append(l, out.Result)
			in = out.Remaining
		}

		out := Output{
			Remaining: in,
		}

		if len(l) == 0 {
			out.Result = Result{
				// Position for empty string ε is not defined
				Val: Empty{},
			}
		} else {
			out.Result = Result{l, l[0].Pos}
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
func (p Parser) REP1() Parser {
	return func(in Input) (Output, bool) {
		if out, ok := p.REP()(in); ok {
			if res, ok := out.Result.Val.(List); ok && len(res) > 0 {
				return out, true
			}
		}

		return Output{}, false
	}
}

func flatten(res Result) List {
	switch v := res.Val.(type) {
	case Empty:
		return List{}

	case List:
		var l List
		for _, v := range v {
			l = append(l, flatten(v)...)
		}
		return l

	default:
		return List{res}
	}
}

// Flatten composes a parser that applies parser p to the input and flattens all results into a single list.
// This can be used for accessing the values of symbols in the right-side of a production rule more intuitively.
func (p Parser) Flatten() Parser {
	return func(in Input) (Output, bool) {
		if out, ok := p(in); ok {
			val := flatten(out.Result)
			return Output{
				Result:    Result{val, out.Result.Pos},
				Remaining: out.Remaining,
			}, true
		}

		return Output{}, false
	}
}

// Select composes a parser that applies parser p to the input and returns a list of symbols from the right-side of the production rule.
// This will not have any effect if the result of parsing is not a list.
// If indices are invalid, you will get an empty string (ε).
func (p Parser) Select(i ...int) Parser {
	return func(in Input) (Output, bool) {
		out, ok := p(in)
		if !ok {
			return Output{}, false
		}

		l, ok := out.Result.Val.(List)
		if !ok {
			return out, true
		}

		var sub List
		for _, j := range i {
			if 0 <= j && j < len(l) {
				sub = append(sub, l[j])
			}
		}

		var res Result
		if len(sub) > 0 {
			res = Result{sub, sub[0].Pos}
		} else {
			res = Result{
				// Position for empty string ε is not defined
				Val: Empty{},
			}
		}

		return Output{
			Result:    res,
			Remaining: out.Remaining,
		}, true
	}
}

// Get composes a parser that applies parser p to the input and returns the value of a symbol from the right-side of the production rule.
// This can be used after CONCAT, REP, REP1, Flatten, and/or Select.
// It will not have any effect if used after other operators and the result of parsing is not a list.
// If index is invalid, you will get an empty string (ε).
func (p Parser) Get(i int) Parser {
	return func(in Input) (Output, bool) {
		out, ok := p(in)
		if !ok {
			return Output{}, false
		}

		l, ok := out.Result.Val.(List)
		if !ok {
			return out, true
		}

		var res Result
		if 0 <= i && i < len(l) {
			res = l[i]
		} else {
			res = Result{
				// Position for empty string ε is not defined
				Val: Empty{},
			}
		}

		return Output{
			Result:    res,
			Remaining: out.Remaining,
		}, ok
	}
}

// Map composes a parser that uses parser p to parse the input and applies a mapper function to the result of parsing.
// If the parser does not succeed, the mapper function will not be applied.
func (p Parser) Map(f Mapper) Parser {
	return func(in Input) (Output, bool) {
		if out, ok := p(in); ok {
			if val, ok := f(out.Result.Val); ok {
				return Output{
					Result:    Result{val, out.Result.Pos},
					Remaining: out.Remaining,
				}, true
			}
		}

		return Output{}, false
	}
}

// Bind composes a parser that uses parser p to parse the input and builds a second parser from the result of parsing.
// It then applies the new parser to the remaining input from the first parser.
// You can use this to implement syntax annotations.
func (p Parser) Bind(f Binder) Parser {
	return func(in Input) (Output, bool) {
		out, ok := p(in)
		if !ok {
			return Output{}, false
		}

		return f(out.Result)(out.Remaining)
	}
}
