package generator

import (
	"fmt"
	"io"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"

	"github.com/gardenbed/emerge/internal/ebnf/parser"
)

var predefs = map[string]string{
	"$DIGIT":  `[0-9]+`,
	"$LETTER": `[A-Za-z]`,
	"$ID":     `[A-Za-z_][0-9A-Za-z_]*`,
	"$INT":    `-?[0-9]+`,
	"$UINT":   `[0-9]+`,
	"$FLOAT":  `-?[0-9]+(\.[0-9]+)?`,
}

type Generator struct {
	parser *parser.Parser
	table  *SymbolTable

	artifacts struct {
		CFG         *grammar.CFG
		Precedences lr.PrecedenceLevels
	}
}

func New(filename string, src io.Reader) (*Generator, error) {
	parser, err := parser.New(filename, src)
	if err != nil {
		return nil, err
	}

	table := NewSymbolTable()

	return &Generator{
		parser: parser,
		table:  table,
	}, nil
}

func (g *Generator) parse() error {
	g.table.Reset()

	res, err := g.parser.ParseAndEvaluate(g.buildGrammar)
	if err != nil {
		return err
	}

	vals := res.Val.([]any)
	g.artifacts.CFG = vals[0].(*grammar.CFG)
	g.artifacts.Precedences = vals[1].(lr.PrecedenceLevels)

	return nil
}

// buildGrammar is a parser.EvaluateFunc for building a context-free grammar.
func (g *Generator) buildGrammar(i int, rhs []*lr.Value) (any, error) {
	switch i {
	// term → STRING
	case 34:
		a := grammar.Terminal(rhs[0].Val.(string))
		g.table.AddStringTerminal(a, rhs[0].Pos)
		return a, nil

	// term → TOKEN
	case 33:
		a := grammar.Terminal(rhs[0].Val.(string))
		g.table.AddTokenTerminal(a, rhs[0].Pos)
		return a, nil

	// nonterm → IDENT
	case 32:
		A := grammar.NonTerminal(rhs[0].Val.(string))
		g.table.AddNonTerminal(A, rhs[0].Pos)
		return A, nil

	// rhs → term
	case 31:
		a := rhs[0].Val.(grammar.Terminal)
		α := grammar.String[grammar.Symbol]{a}
		return Strings{α}, nil

	// rhs → nonterm
	case 30:
		A := rhs[0].Val.(grammar.NonTerminal)
		α := grammar.String[grammar.Symbol]{A}
		return Strings{α}, nil

	// rhs → rhs "|"
	case 29:
		s := rhs[0].Val.(Strings)

		var all Strings
		all = append(all, s...)
		all = append(all, grammar.E)

		return all, nil

	// rhs → rhs "|" rhs
	case 28:
		s1 := rhs[0].Val.(Strings)
		s2 := rhs[2].Val.(Strings)

		var all Strings
		all = append(all, s1...)
		all = append(all, s2...)

		return all, nil

	// rhs → "{{" rhs "}}"
	case 27:
		s := rhs[1].Val.(Strings)
		plus := g.table.GetPlus(s)
		g.table.AddNonTerminal(plus, rhs[1].Pos)

		for _, α := range s {
			g.table.AddProduction(
				&grammar.Production{Head: plus, Body: α.Prepend(plus)},
				rhs[0].Pos,
			)

			g.table.AddProduction(
				&grammar.Production{Head: plus, Body: α},
				rhs[0].Pos,
			)
		}

		return Strings{{plus}}, nil

	// rhs → "{" rhs "}"
	case 26:
		s := rhs[1].Val.(Strings)
		star := g.table.GetStar(s)
		g.table.AddNonTerminal(star, rhs[1].Pos)

		for _, α := range s {
			g.table.AddProduction(
				&grammar.Production{Head: star, Body: α.Prepend(star)},
				rhs[0].Pos,
			)
		}

		g.table.AddProduction(
			&grammar.Production{Head: star, Body: grammar.E},
			rhs[0].Pos,
		)

		return Strings{{star}}, nil

	// rhs → "[" rhs "]"
	case 25:
		s := rhs[1].Val.(Strings)
		opt := g.table.GetOpt(s)
		g.table.AddNonTerminal(opt, rhs[1].Pos)

		for _, α := range s {
			g.table.AddProduction(
				&grammar.Production{Head: opt, Body: α},
				rhs[0].Pos,
			)
		}

		g.table.AddProduction(
			&grammar.Production{Head: opt, Body: grammar.E},
			rhs[0].Pos,
		)

		return Strings{{opt}}, nil

	// rhs → "(" rhs ")"
	case 24:
		s := rhs[1].Val.(Strings)
		group := g.table.GetGroup(s)
		g.table.AddNonTerminal(group, rhs[1].Pos)

		for _, α := range s {
			g.table.AddProduction(
				&grammar.Production{Head: group, Body: α},
				rhs[0].Pos,
			)
		}

		return Strings{{group}}, nil

	// rhs → rhs rhs
	case 23:
		s1 := rhs[0].Val.(Strings)
		s2 := rhs[1].Val.(Strings)

		var all Strings
		for _, α := range s1 {
			for _, β := range s2 {
				all = append(all, α.Concat(β))
			}
		}

		return all, nil

	// lhs → nonterm
	case 22:
		return rhs[0].Val, nil

	// rule → lhs "="
	case 21:
		p := &grammar.Production{
			Head: rhs[0].Val.(grammar.NonTerminal),
			Body: grammar.E,
		}

		g.table.AddProduction(p, rhs[0].Pos)
		prods := []*grammar.Production{p}

		return prods, nil

	// rule → lhs "=" rhs
	case 20:
		head := rhs[0].Val.(grammar.NonTerminal)
		s := rhs[2].Val.(Strings)

		prods := []*grammar.Production{}
		for _, α := range s {
			p := &grammar.Production{Head: head, Body: α}
			g.table.AddProduction(p, rhs[0].Pos)
			prods = append(prods, p)
		}

		return prods, nil

	// rule_handle → "<" rule ">"
	case 19:
		return rhs[1].Val, nil

	// handles → rule_handle
	case 18:
		prods := rhs[0].Val.([]*grammar.Production)

		var handles []*lr.PrecedenceHandle
		for _, p := range prods {
			handles = append(handles, &lr.PrecedenceHandle{
				Production: p,
			})
		}

		return handles, nil

	// handles → term
	case 17:
		term := rhs[0].Val.(grammar.Terminal)

		handles := []*lr.PrecedenceHandle{
			{Terminal: &term},
		}

		return handles, nil

	// handles → handles rule_handle
	case 16:
		handles := rhs[0].Val.([]*lr.PrecedenceHandle)
		prods := rhs[1].Val.([]*grammar.Production)

		for _, p := range prods {
			handles = append(handles, &lr.PrecedenceHandle{
				Production: p,
			})
		}

		return handles, nil

	// handles → handles term
	case 15:
		handles := rhs[0].Val.([]*lr.PrecedenceHandle)
		term := rhs[1].Val.(grammar.Terminal)

		handles = append(handles, &lr.PrecedenceHandle{
			Terminal: &term,
		})

		return handles, nil

	// directive → "@none" handles
	case 14:
		handles := rhs[1].Val.([]*lr.PrecedenceHandle)

		p := &lr.PrecedenceLevel{
			Associativity: lr.NONE,
			Handles:       lr.NewPrecedenceHandles(handles...),
		}

		g.table.AddPrecedence(p)

		return p, nil

	// directive → "@right" handles
	case 13:
		handles := rhs[1].Val.([]*lr.PrecedenceHandle)

		p := &lr.PrecedenceLevel{
			Associativity: lr.RIGHT,
			Handles:       lr.NewPrecedenceHandles(handles...),
		}

		g.table.AddPrecedence(p)

		return p, nil

	// directive → "@left" handles
	case 12:
		handles := rhs[1].Val.([]*lr.PrecedenceHandle)

		p := &lr.PrecedenceLevel{
			Associativity: lr.LEFT,
			Handles:       lr.NewPrecedenceHandles(handles...),
		}

		g.table.AddPrecedence(p)

		return p, nil

	// token → TOKEN "=" PREDEF
	case 11:
		token := grammar.Terminal(rhs[0].Val.(string))
		value := rhs[2].Val.(string)

		regex, ok := predefs[value]
		if !ok {
			return nil, fmt.Errorf("invalid predefined regex: %s", value)
		}

		g.table.AddRegexTokenDef(token, regex, rhs[0].Pos)

		return nil, nil

	// token → TOKEN "=" REGEX
	case 10:
		token := grammar.Terminal(rhs[0].Val.(string))
		value := rhs[2].Val.(string)

		g.table.AddRegexTokenDef(token, value, rhs[0].Pos)

		return nil, nil

	// token → TOKEN "=" STRING
	case 9:
		token := grammar.Terminal(rhs[0].Val.(string))
		value := rhs[2].Val.(string)

		g.table.AddStringTokenDef(token, value, rhs[0].Pos)

		return nil, nil

	// semi_opt → ε
	case 8:
		return nil, nil

	// semi_opt → ";"
	case 7:
		return nil, nil

	// decl → rule ";"
	case 6:
		return nil, nil

	// decl → directive semi_opt
	case 5:
		return nil, nil

	// decl → token semi_opt
	case 4:
		return nil, nil

	// decls → ε
	case 3:
		return nil, nil

	// decls → decls decl
	case 2:
		return nil, nil

	// name → "grammar" IDENT semi_opt
	case 1:
		return rhs[1].Val, nil

	// grammar → name decls
	case 0:
		if err := g.table.Verify(); err != nil {
			return nil, err
		}

		P := g.table.Precedences()
		G := grammar.NewCFG(
			g.table.Terminals(),
			g.table.NonTerminals(),
			g.table.Productions(),
			"start",
		)

		return []any{G, P}, nil
	}

	return nil, fmt.Errorf("invalid production index: %d", i)
}
