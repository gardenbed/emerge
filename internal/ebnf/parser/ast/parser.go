package ast

import (
	"fmt"
	"io"

	"github.com/moorara/algo/parser/lr"

	"github.com/gardenbed/emerge/internal/ebnf/parser"
)

// Parse implements the LR parsing algorithm.
// It analyzes a sequence of input tokens (terminal symbols) provided by a lexical analyzer.
// It attempts to parse the input according to the production rules of the EBNF grammar.
//
// If the input string is valid, the root node of the EBNF AST is returned,
// representing the syntactic structure of the input string.
//
// An error is returned if the input fails to conform to the grammar rules, indicating a syntax issue.
func Parse(filename string, src io.Reader) (*Grammar, error) {
	p, err := parser.New(filename, src)
	if err != nil {
		return nil, err
	}

	res, err := p.ParseAndEvaluate(func(i int, rhs []*lr.Value) (any, error) {
		switch i {
		// term → STRING
		case 34:
			return fmt.Sprintf("%q", rhs[0].Val), nil

		// term → TOKEN
		case 33:
			return rhs[0].Val.(string), nil

		// nonterm → IDENT
		case 32:
			return rhs[0].Val.(string), nil

		// rhs → term
		case 31:
			return &TerminalRHS{
				Terminal: rhs[0].Val.(string),
				Position: rhs[0].Pos,
			}, nil

		// rhs → nonterm
		case 30:
			return &NonTerminalRHS{
				NonTerminal: rhs[0].Val.(string),
				Position:    rhs[0].Pos,
			}, nil

		// rhs → rhs "|"
		case 29:
			var ops []RHS

			if c, ok := rhs[0].Val.(*AltRHS); ok {
				ops = append(ops, c.Ops...)
			} else {
				ops = append(ops, rhs[0].Val.(RHS))
			}

			ops = append(ops, &EmptyRHS{})

			return &AltRHS{
				Ops: ops,
			}, nil

		// rhs → rhs "|" rhs
		case 28:
			var ops []RHS

			if c, ok := rhs[0].Val.(*AltRHS); ok {
				ops = append(ops, c.Ops...)
			} else {
				ops = append(ops, rhs[0].Val.(RHS))
			}

			if c, ok := rhs[2].Val.(*AltRHS); ok {
				ops = append(ops, c.Ops...)
			} else {
				ops = append(ops, rhs[2].Val.(RHS))
			}

			return &AltRHS{
				Ops: ops,
			}, nil

		// rhs → "{{" rhs "}}"
		case 27:
			return &PlusRHS{
				Op:       rhs[1].Val.(RHS),
				Position: rhs[0].Pos,
			}, nil

		// rhs → "{" rhs "}"
		case 26:
			return &StarRHS{
				Op:       rhs[1].Val.(RHS),
				Position: rhs[0].Pos,
			}, nil

		// rhs → "[" rhs "]"
		case 25:
			return &OptRHS{
				Op:       rhs[1].Val.(RHS),
				Position: rhs[0].Pos,
			}, nil

		// rhs → "(" rhs ")"
		case 24:
			return rhs[1].Val, nil

		// rhs → rhs rhs
		case 23:
			var ops []RHS

			if c, ok := rhs[0].Val.(*ConcatRHS); ok {
				ops = append(ops, c.Ops...)
			} else {
				ops = append(ops, rhs[0].Val.(RHS))
			}

			if c, ok := rhs[1].Val.(*ConcatRHS); ok {
				ops = append(ops, c.Ops...)
			} else {
				ops = append(ops, rhs[1].Val.(RHS))
			}

			return &ConcatRHS{
				Ops: ops,
			}, nil

		// lhs → nonterm
		case 22:
			return rhs[0].Val, nil

		// rule → lhs "="
		case 21:
			return &RuleDecl{
				LHS:      rhs[0].Val.(string),
				RHS:      &EmptyRHS{},
				Position: rhs[0].Pos,
			}, nil

		// rule → lhs "=" rhs
		case 20:
			return &RuleDecl{
				LHS:      rhs[0].Val.(string),
				RHS:      rhs[2].Val.(RHS),
				Position: rhs[0].Pos,
			}, nil

		// rule_handle → "<" rule ">"
		case 19:
			return rhs[1].Val, nil

		// handles → rule_handle
		case 18:
			rule := rhs[0].Val.(*RuleDecl)
			return []PrecedenceHandle{
				&ProductionHandle{
					LHS:      rule.LHS,
					RHS:      rule.RHS,
					Position: rhs[0].Pos,
				},
			}, nil

		// handles → term
		case 17:
			return []PrecedenceHandle{
				&TerminalHandle{
					Terminal: rhs[0].Val.(string),
					Position: rhs[0].Pos,
				},
			}, nil

		// handles → handles rule_handle
		case 16:
			handles := rhs[0].Val.([]PrecedenceHandle)
			rule := rhs[1].Val.(*RuleDecl)

			handles = append(handles, &ProductionHandle{
				LHS:      rule.LHS,
				RHS:      rule.RHS,
				Position: rhs[1].Pos,
			})

			return handles, nil

		// handles → handles term
		case 15:
			handles := rhs[0].Val.([]PrecedenceHandle)

			handles = append(handles, &TerminalHandle{
				Terminal: rhs[1].Val.(string),
				Position: rhs[1].Pos,
			})

			return handles, nil

		// directive → "@none" handles
		case 14:
			return &PrecedenceDecl{
				Associativity: lr.NONE,
				Handles:       rhs[1].Val.([]PrecedenceHandle),
				Position:      rhs[0].Pos,
			}, nil

		// directive → "@right" handles
		case 13:
			return &PrecedenceDecl{
				Associativity: lr.RIGHT,
				Handles:       rhs[1].Val.([]PrecedenceHandle),
				Position:      rhs[0].Pos,
			}, nil

		// directive → "@left" handles
		case 12:
			return &PrecedenceDecl{
				Associativity: lr.LEFT,
				Handles:       rhs[1].Val.([]PrecedenceHandle),
				Position:      rhs[0].Pos,
			}, nil

		// token → TOKEN "=" PREDEF
		case 11:
			value := rhs[2].Val.(string)

			regex, ok := parser.Predefs[value]
			if !ok {
				return nil, fmt.Errorf("invalid predefined regex: %s", value)
			}

			return &RegexTokenDecl{
				Name:     rhs[0].Val.(string),
				Regex:    regex,
				Position: rhs[0].Pos,
			}, nil

		// token → TOKEN "=" REGEX
		case 10:
			return &RegexTokenDecl{
				Name:     rhs[0].Val.(string),
				Regex:    rhs[2].Val.(string),
				Position: rhs[0].Pos,
			}, nil

		// token → TOKEN "=" STRING
		case 9:
			return &StringTokenDecl{
				Name:     rhs[0].Val.(string),
				Value:    rhs[2].Val.(string),
				Position: rhs[0].Pos,
			}, nil

		// semi_opt → ε
		case 8:
			// Discard
			return nil, nil

		// semi_opt → ";"
		case 7:
			// Discard
			return nil, nil

		// decl → rule ";"
		case 6:
			return rhs[0].Val, nil

		// decl → directive semi_opt
		case 5:
			return rhs[0].Val, nil

		// decl → token semi_opt
		case 4:
			return rhs[0].Val, nil

		// decls → ε
		case 3:
			// Discard
			return nil, nil

		// decls → decls decl
		case 2:
			var decls []Decl

			if rhs[0].Val != nil {
				decls = rhs[0].Val.([]Decl)
			}

			decls = append(decls, rhs[1].Val.(Decl))

			return decls, nil

		// name → "grammar" IDENT semi_opt
		case 1:
			return rhs[1].Val.(string), nil

		// grammar → name decls
		case 0:
			return &Grammar{
				Name:     rhs[0].Val.(string),
				Decls:    rhs[1].Val.([]Decl),
				Position: rhs[0].Pos,
			}, nil
		}

		return nil, fmt.Errorf("invalid production index: %d", i)
	})

	if err != nil {
		return nil, err
	}

	return res.Val.(*Grammar), nil
}
