package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/parser/lr/lookahead"
)

func main() {
	code, err := GenerateParsingTable(G, precedences)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile("parsing_table.go", code, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("OK")
}

func GenerateParsingTable(G *grammar.CFG, precedences lr.PrecedenceLevels) ([]byte, error) {
	T, err := lookahead.BuildParsingTable(G, precedences)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	b.WriteString(`//go:generate go run ./generate

package parser

import (
	"fmt"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
)

var (
	// terminals is an ordered list of terminal symbols for the EBNF grammar.
	terminals = []grammar.Terminal{
		"=", ";", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "<", ">",
		"grammar", "@left", "@right", "@none",
		"IDENT", "TOKEN", "STRING", "REGEX", "PREDEF",
	}

	// nonTerminals is an ordered list of non-terminal symbols for the EBNF grammar.
	nonTerminals = []grammar.NonTerminal{
		"grammar", "name", "decls", "decl", "semi_opt", "token",
		"directive", "handles", "rule_handle", "rule", "lhs", "rhs", "nonterm", "term",
	}

	// productions is an ordered list of productions rules for the EBNF grammar.
	productions = []*grammar.Production{
		/*  0: grammar → name decls */ {Head: "grammar", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("name"), grammar.NonTerminal("decls")}},
		/*  1: name → "grammar" IDENT semi_opt */ {Head: "name", Body: grammar.String[grammar.Symbol]{grammar.Terminal("grammar"), grammar.Terminal("IDENT"), grammar.NonTerminal("semi_opt")}},
		/*  2: decls → decls decl */ {Head: "decls", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("decls"), grammar.NonTerminal("decl")}},
		/*  3: decls → ε */ {Head: "decls", Body: grammar.E},
		/*  4: decl → token semi_opt */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("token"), grammar.NonTerminal("semi_opt")}},
		/*  5: decl → directive semi_opt */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("directive"), grammar.NonTerminal("semi_opt")}},
		/*  6: decl → rule ";" */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rule"), grammar.Terminal(";")}},
		/*  7: semi_opt → ";" */ {Head: "semi_opt", Body: grammar.String[grammar.Symbol]{grammar.Terminal(";")}},
		/*  8: semi_opt → ε */ {Head: "semi_opt", Body: grammar.E},
		/*  9: token → TOKEN "=" STRING */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("STRING")}},
		/* 10: token → TOKEN "=" REGEX */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("REGEX")}},
		/* 11: token → TOKEN "=" PREDEF */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("PREDEF")}},
		/* 12: directive → "@left" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@left"), grammar.NonTerminal("handles")}},
		/* 13: directive → "@right" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@right"), grammar.NonTerminal("handles")}},
		/* 14: directive → "@none" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@none"), grammar.NonTerminal("handles")}},
		/* 15: handles → handles term */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("handles"), grammar.NonTerminal("term")}},
		/* 16: handles → handles rule_handle */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("handles"), grammar.NonTerminal("rule_handle")}},
		/* 17: handles → term */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("term")}},
		/* 18: handles → rule_handle */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rule_handle")}},
		/* 19: rule_handle → "<" rule ">" */ {Head: "rule_handle", Body: grammar.String[grammar.Symbol]{grammar.Terminal("<"), grammar.NonTerminal("rule"), grammar.Terminal(">")}},
		/* 20: rule → lhs "=" rhs */ {Head: "rule", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("lhs"), grammar.Terminal("="), grammar.NonTerminal("rhs")}},
		/* 21: rule → lhs "=" */ {Head: "rule", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("lhs"), grammar.Terminal("=")}},
		/* 22: lhs → nonterm */ {Head: "lhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("nonterm")}},
		/* 23: rhs → rhs rhs */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.NonTerminal("rhs")}},
		/* 24: rhs → "(" rhs ")" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("("), grammar.NonTerminal("rhs"), grammar.Terminal(")")}},
		/* 25: rhs → "[" rhs "]" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("["), grammar.NonTerminal("rhs"), grammar.Terminal("]")}},
		/* 26: rhs → "{" rhs "}" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("{"), grammar.NonTerminal("rhs"), grammar.Terminal("}")}},
		/* 27: rhs → "{{" rhs "}}" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("{{"), grammar.NonTerminal("rhs"), grammar.Terminal("}}")}},
		/* 28: rhs → rhs "|" rhs */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.Terminal("|"), grammar.NonTerminal("rhs")}},
		/* 29: rhs → rhs "|" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.Terminal("|")}},
		/* 30: rhs → nonterm */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("nonterm")}},
		/* 31: rhs → term */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("term")}},
		/* 32: nonterm → IDENT */ {Head: "nonterm", Body: grammar.String[grammar.Symbol]{grammar.Terminal("IDENT")}},
		/* 33: term → TOKEN */ {Head: "term", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN")}},
		/* 34: term → STRING */ {Head: "term", Body: grammar.String[grammar.Symbol]{grammar.Terminal("STRING")}},
	}

	// G is the EBNF grammar.
	G = grammar.NewCFG(
		terminals,
		nonTerminals,
		productions,
		"grammar",
	)

	// precedences define the associativity and precedence for the EBNF grammar.
	precedences = lr.PrecedenceLevels{
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForProduction(&grammar.Production{
					Head: "rhs",
					Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.NonTerminal("rhs")},
				}),
			),
		},
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("("),
				lr.PrecedenceHandleForTerminal("["),
				lr.PrecedenceHandleForTerminal("{"),
				lr.PrecedenceHandleForTerminal("{{"),
				lr.PrecedenceHandleForTerminal("IDENT"),
				lr.PrecedenceHandleForTerminal("TOKEN"),
				lr.PrecedenceHandleForTerminal("STRING"),
			),
		},
		{
			Associativity: lr.RIGHT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("|"),
			),
		},
		{
			Associativity: lr.NONE,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("="),
			),
		},
		{
			Associativity: lr.NONE,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("@left"),
				lr.PrecedenceHandleForTerminal("@right"),
				lr.PrecedenceHandleForTerminal("@none"),
			),
		},
	}
)

`)

	b.Write(generateACTION(T))
	b.WriteString("\n")
	b.Write(generateGOTO(T))

	return b.Bytes(), nil
}

func generateACTION(T *lr.ParsingTable) []byte {
	var b, c bytes.Buffer

	b.WriteString("// ACTION looks up and returns the action for state s and terminal a.\n")
	b.WriteString("func ACTION(s int, a grammar.Terminal) (lr.ActionType, int, error) {\n")
	b.WriteString("	switch s {\n")

	for _, s := range T.States {
		c.Reset()

		for _, a := range append(terminals, grammar.Endmarker) {
			if action, err := T.ACTION(s, a); err == nil {
				if a == grammar.Endmarker {
					fmt.Fprintf(&c, "		case grammar.Endmarker:\n")
				} else {
					fmt.Fprintf(&c, "		case %q:\n", string(a))
				}

				switch action.Type {
				case lr.SHIFT:
					fmt.Fprintf(&c, "			return lr.SHIFT, %d, nil // %s\n", action.State, action)
				case lr.REDUCE:
					if i := productionIndex(productions, action.Production); i >= 0 {
						fmt.Fprintf(&c, "			return lr.REDUCE, %d, nil // %s\n", i, action)
					} else {
						fmt.Fprintf(&c, "			return lr.REDUCE, , nil // %s\n", action)
					}
				case lr.ACCEPT:
					fmt.Fprintf(&c, "			return lr.ACCEPT, 0, nil // ACCEPT\n")
				}
			}
		}

		if c.Len() > 0 {
			fmt.Fprintf(&b, "	case %d:\n", s)
			fmt.Fprintf(&b, "		switch a {\n")
			b.Write(c.Bytes())
			b.WriteString("		}\n")
			b.WriteString("\n")
		}
	}

	b.WriteString("	}\n")
	b.WriteString("\n")
	b.WriteString("	return lr.ERROR, -1, fmt.Errorf(\"no action exists in the parsing table for ACTION[%d, %q]\", s, a)\n")
	b.WriteString("}\n")

	return b.Bytes()
}

func generateGOTO(T *lr.ParsingTable) []byte {
	var b, c bytes.Buffer

	b.WriteString("// GOTO looks up and returns the next state for state s and non-terminal A.\n")
	b.WriteString("func GOTO(s int, A grammar.NonTerminal) int {\n")
	b.WriteString("	switch s {\n")

	for _, s := range T.States {
		c.Reset()

		for _, A := range nonTerminals {
			if next, err := T.GOTO(s, A); err == nil {
				fmt.Fprintf(&c, "		case %q:\n", string(A))
				fmt.Fprintf(&c, "			return %d\n", next)
			}
		}

		if c.Len() > 0 {
			fmt.Fprintf(&b, "	case %d:\n", s)
			fmt.Fprintf(&b, "		switch A {\n")
			b.Write(c.Bytes())
			b.WriteString("		}\n")
			b.WriteString("\n")
		}
	}

	b.WriteString("	}\n")
	b.WriteString("\n")
	b.WriteString("	return -1\n")
	b.WriteString("}\n")

	return b.Bytes()
}

func productionIndex(prods []*grammar.Production, p *grammar.Production) int {
	for i, prod := range prods {
		if prod.Equal(p) {
			return i
		}
	}

	return -1
}

var (
	// terminals is an ordered list of terminal symbols for the EBNF grammar.
	terminals = []grammar.Terminal{
		"=", ";", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "<", ">",
		"grammar", "@left", "@right", "@none",
		"IDENT", "TOKEN", "STRING", "REGEX", "PREDEF",
	}

	// nonTerminals is an ordered list of non-terminal symbols for the EBNF grammar.
	nonTerminals = []grammar.NonTerminal{
		"grammar", "name", "decls", "decl", "semi_opt", "token",
		"directive", "handles", "rule_handle", "rule", "lhs", "rhs", "nonterm", "term",
	}

	// productions is an ordered list of productions rules for the EBNF grammar.
	productions = []*grammar.Production{
		/*  0: grammar → name decls */ {Head: "grammar", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("name"), grammar.NonTerminal("decls")}},
		/*  1: name → "grammar" IDENT semi_opt */ {Head: "name", Body: grammar.String[grammar.Symbol]{grammar.Terminal("grammar"), grammar.Terminal("IDENT"), grammar.NonTerminal("semi_opt")}},
		/*  2: decls → decls decl */ {Head: "decls", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("decls"), grammar.NonTerminal("decl")}},
		/*  3: decls → ε */ {Head: "decls", Body: grammar.E},
		/*  4: decl → token semi_opt */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("token"), grammar.NonTerminal("semi_opt")}},
		/*  5: decl → directive semi_opt */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("directive"), grammar.NonTerminal("semi_opt")}},
		/*  6: decl → rule ";" */ {Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rule"), grammar.Terminal(";")}},
		/*  7: semi_opt → ";" */ {Head: "semi_opt", Body: grammar.String[grammar.Symbol]{grammar.Terminal(";")}},
		/*  8: semi_opt → ε */ {Head: "semi_opt", Body: grammar.E},
		/*  9: token → TOKEN "=" STRING */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("STRING")}},
		/* 10: token → TOKEN "=" REGEX */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("REGEX")}},
		/* 11: token → TOKEN "=" PREDEF */ {Head: "token", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN"), grammar.Terminal("="), grammar.Terminal("PREDEF")}},
		/* 12: directive → "@left" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@left"), grammar.NonTerminal("handles")}},
		/* 13: directive → "@right" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@right"), grammar.NonTerminal("handles")}},
		/* 14: directive → "@none" handles */ {Head: "directive", Body: grammar.String[grammar.Symbol]{grammar.Terminal("@none"), grammar.NonTerminal("handles")}},
		/* 15: handles → handles term */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("handles"), grammar.NonTerminal("term")}},
		/* 16: handles → handles rule_handle */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("handles"), grammar.NonTerminal("rule_handle")}},
		/* 17: handles → term */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("term")}},
		/* 18: handles → rule_handle */ {Head: "handles", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rule_handle")}},
		/* 19: rule_handle → "<" rule ">" */ {Head: "rule_handle", Body: grammar.String[grammar.Symbol]{grammar.Terminal("<"), grammar.NonTerminal("rule"), grammar.Terminal(">")}},
		/* 20: rule → lhs "=" rhs */ {Head: "rule", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("lhs"), grammar.Terminal("="), grammar.NonTerminal("rhs")}},
		/* 21: rule → lhs "=" */ {Head: "rule", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("lhs"), grammar.Terminal("=")}},
		/* 22: lhs → nonterm */ {Head: "lhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("nonterm")}},
		/* 23: rhs → rhs rhs */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.NonTerminal("rhs")}},
		/* 24: rhs → "(" rhs ")" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("("), grammar.NonTerminal("rhs"), grammar.Terminal(")")}},
		/* 25: rhs → "[" rhs "]" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("["), grammar.NonTerminal("rhs"), grammar.Terminal("]")}},
		/* 26: rhs → "{" rhs "}" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("{"), grammar.NonTerminal("rhs"), grammar.Terminal("}")}},
		/* 27: rhs → "{{" rhs "}}" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.Terminal("{{"), grammar.NonTerminal("rhs"), grammar.Terminal("}}")}},
		/* 28: rhs → rhs "|" rhs */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.Terminal("|"), grammar.NonTerminal("rhs")}},
		/* 29: rhs → rhs "|" */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.Terminal("|")}},
		/* 30: rhs → nonterm */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("nonterm")}},
		/* 31: rhs → term */ {Head: "rhs", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("term")}},
		/* 32: nonterm → IDENT */ {Head: "nonterm", Body: grammar.String[grammar.Symbol]{grammar.Terminal("IDENT")}},
		/* 33: term → TOKEN */ {Head: "term", Body: grammar.String[grammar.Symbol]{grammar.Terminal("TOKEN")}},
		/* 34: term → STRING */ {Head: "term", Body: grammar.String[grammar.Symbol]{grammar.Terminal("STRING")}},
	}

	// G is the EBNF grammar.
	G = grammar.NewCFG(
		terminals,
		nonTerminals,
		productions,
		"grammar",
	)

	// precedences define the associativity and precedence for the EBNF grammar.
	precedences = lr.PrecedenceLevels{
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForProduction(&grammar.Production{
					Head: "rhs",
					Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("rhs"), grammar.NonTerminal("rhs")},
				}),
			),
		},
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("("),
				lr.PrecedenceHandleForTerminal("["),
				lr.PrecedenceHandleForTerminal("{"),
				lr.PrecedenceHandleForTerminal("{{"),
				lr.PrecedenceHandleForTerminal("IDENT"),
				lr.PrecedenceHandleForTerminal("TOKEN"),
				lr.PrecedenceHandleForTerminal("STRING"),
			),
		},
		{
			Associativity: lr.RIGHT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("|"),
			),
		},
		{
			Associativity: lr.NONE,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("="),
			),
		},
		{
			Associativity: lr.NONE,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("@left"),
				lr.PrecedenceHandleForTerminal("@right"),
				lr.PrecedenceHandleForTerminal("@none"),
			),
		},
	}
)
