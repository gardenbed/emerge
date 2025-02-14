//go:generate go run ./generate

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

// ACTION looks up and returns the action for state s and terminal a.
func ACTION(s int, a grammar.Terminal) (lr.ActionType, int, error) {
	switch s {
	case 0:
		switch a {
		case "grammar":
			return lr.SHIFT, 43, nil // SHIFT 43
		}

	case 1:
		switch a {
		case grammar.Endmarker:
			return lr.ACCEPT, 0, nil // ACCEPT
		}

	case 2:
		switch a {
		case "@left":
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		case "@right":
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		case "@none":
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		case "IDENT":
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		case "TOKEN":
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		case grammar.Endmarker:
			return lr.REDUCE, 1, nil // REDUCE name → "grammar" "IDENT" semi_opt
		}

	case 3:
		switch a {
		case ";":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case ")":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "]":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "}":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "}}":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case ">":
			return lr.REDUCE, 28, nil // REDUCE rhs → rhs "|" rhs
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 4:
		switch a {
		case ";":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "|":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "(":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case ")":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "[":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "]":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "{":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "}":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "{{":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "}}":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case ">":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "IDENT":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "TOKEN":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		case "STRING":
			return lr.REDUCE, 24, nil // REDUCE rhs → "(" rhs ")"
		}

	case 5:
		switch a {
		case ";":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "|":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "(":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case ")":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "[":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "]":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "{":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "}":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "{{":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "}}":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case ">":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "IDENT":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "TOKEN":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		case "STRING":
			return lr.REDUCE, 25, nil // REDUCE rhs → "[" rhs "]"
		}

	case 6:
		switch a {
		case ";":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "|":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "(":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case ")":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "[":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "]":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "{":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "}":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "{{":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "}}":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case ">":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "IDENT":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "TOKEN":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		case "STRING":
			return lr.REDUCE, 26, nil // REDUCE rhs → "{" rhs "}"
		}

	case 7:
		switch a {
		case ";":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "|":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "(":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case ")":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "[":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "]":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "{":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "}":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "{{":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "}}":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case ">":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "IDENT":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "TOKEN":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		case "STRING":
			return lr.REDUCE, 27, nil // REDUCE rhs → "{{" rhs "}}"
		}

	case 8:
		switch a {
		case ";":
			return lr.REDUCE, 20, nil // REDUCE rule → lhs "=" rhs
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case ">":
			return lr.REDUCE, 20, nil // REDUCE rule → lhs "=" rhs
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 9:
		switch a {
		case ";":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "<":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "@left":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "@right":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "@none":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "IDENT":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "TOKEN":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case "STRING":
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		case grammar.Endmarker:
			return lr.REDUCE, 19, nil // REDUCE rule_handle → "<" rule ">"
		}

	case 10:
		switch a {
		case ";":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case "@left":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case "@right":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case "@none":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case "IDENT":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case "TOKEN":
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		case grammar.Endmarker:
			return lr.REDUCE, 11, nil // REDUCE token → "TOKEN" "=" "PREDEF"
		}

	case 11:
		switch a {
		case ";":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case "@left":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case "@right":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case "@none":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case "IDENT":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case "TOKEN":
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		case grammar.Endmarker:
			return lr.REDUCE, 10, nil // REDUCE token → "TOKEN" "=" "REGEX"
		}

	case 12:
		switch a {
		case ";":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case "@left":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case "@right":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case "@none":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case "IDENT":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case "TOKEN":
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		case grammar.Endmarker:
			return lr.REDUCE, 9, nil // REDUCE token → "TOKEN" "=" "STRING"
		}

	case 13:
		switch a {
		case "@left":
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		case "@right":
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		case "@none":
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		case "IDENT":
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		case "TOKEN":
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		case grammar.Endmarker:
			return lr.REDUCE, 5, nil // REDUCE decl → directive semi_opt
		}

	case 14:
		switch a {
		case "@left":
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		case "@right":
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		case "@none":
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		case "IDENT":
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		case "TOKEN":
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		case grammar.Endmarker:
			return lr.REDUCE, 4, nil // REDUCE decl → token semi_opt
		}

	case 15:
		switch a {
		case "@left":
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		case "@right":
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		case "@none":
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		case "IDENT":
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		case "TOKEN":
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		case grammar.Endmarker:
			return lr.REDUCE, 6, nil // REDUCE decl → rule ";"
		}

	case 16:
		switch a {
		case "@left":
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		case "@right":
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		case "@none":
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		case "IDENT":
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		case "TOKEN":
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		case grammar.Endmarker:
			return lr.REDUCE, 2, nil // REDUCE decls → decls decl
		}

	case 17:
		switch a {
		case ";":
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "@left":
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		case "@right":
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		case "@none":
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		case "IDENT":
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		case grammar.Endmarker:
			return lr.REDUCE, 12, nil // REDUCE directive → "@left" handles
		}

	case 18:
		switch a {
		case ";":
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "@left":
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		case "@right":
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		case "@none":
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		case "IDENT":
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		case grammar.Endmarker:
			return lr.REDUCE, 14, nil // REDUCE directive → "@none" handles
		}

	case 19:
		switch a {
		case ";":
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "@left":
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		case "@right":
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		case "@none":
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		case "IDENT":
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		case grammar.Endmarker:
			return lr.REDUCE, 13, nil // REDUCE directive → "@right" handles
		}

	case 20:
		switch a {
		case "@left":
			return lr.SHIFT, 36, nil // SHIFT 36
		case "@right":
			return lr.SHIFT, 38, nil // SHIFT 38
		case "@none":
			return lr.SHIFT, 37, nil // SHIFT 37
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 56, nil // SHIFT 56
		case grammar.Endmarker:
			return lr.REDUCE, 0, nil // REDUCE grammar → name decls
		}

	case 21:
		switch a {
		case ";":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "<":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "@left":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "@right":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "@none":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "IDENT":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "TOKEN":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case "STRING":
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		case grammar.Endmarker:
			return lr.REDUCE, 16, nil // REDUCE handles → handles rule_handle
		}

	case 22:
		switch a {
		case ";":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "<":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "@left":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "@right":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "@none":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "IDENT":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "TOKEN":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case "STRING":
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		case grammar.Endmarker:
			return lr.REDUCE, 15, nil // REDUCE handles → handles term
		}

	case 23:
		switch a {
		case ";":
			return lr.SHIFT, 53, nil // SHIFT 53
		case "@left":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@right":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@none":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "IDENT":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "TOKEN":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case grammar.Endmarker:
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		}

	case 24:
		switch a {
		case ";":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "|":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case ")":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "]":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "}":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "}}":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case ">":
			return lr.REDUCE, 29, nil // REDUCE rhs → rhs "|"
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 25:
		switch a {
		case ";":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "|":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "(":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case ")":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "[":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "]":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "{":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "}":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "{{":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "}}":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case ">":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "IDENT":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "TOKEN":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		case "STRING":
			return lr.REDUCE, 23, nil // REDUCE rhs → rhs rhs
		}

	case 26:
		switch a {
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case ")":
			return lr.SHIFT, 4, nil // SHIFT 4
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 27:
		switch a {
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "]":
			return lr.SHIFT, 5, nil // SHIFT 5
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 28:
		switch a {
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "}":
			return lr.SHIFT, 6, nil // SHIFT 6
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 29:
		switch a {
		case "|":
			return lr.SHIFT, 24, nil // SHIFT 24
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "}}":
			return lr.SHIFT, 7, nil // SHIFT 7
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 30:
		switch a {
		case ";":
			return lr.REDUCE, 21, nil // REDUCE rule → lhs "="
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case ">":
			return lr.REDUCE, 21, nil // REDUCE rule → lhs "="
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 31:
		switch a {
		case ">":
			return lr.SHIFT, 9, nil // SHIFT 9
		}

	case 32:
		switch a {
		case "STRING":
			return lr.SHIFT, 12, nil // SHIFT 12
		case "REGEX":
			return lr.SHIFT, 11, nil // SHIFT 11
		case "PREDEF":
			return lr.SHIFT, 10, nil // SHIFT 10
		}

	case 33:
		switch a {
		case ";":
			return lr.SHIFT, 53, nil // SHIFT 53
		case "@left":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@right":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@none":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "IDENT":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "TOKEN":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case grammar.Endmarker:
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		}

	case 34:
		switch a {
		case ";":
			return lr.SHIFT, 53, nil // SHIFT 53
		case "@left":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@right":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "@none":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "IDENT":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case "TOKEN":
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		case grammar.Endmarker:
			return lr.REDUCE, 8, nil // REDUCE semi_opt → ε
		}

	case 35:
		switch a {
		case ";":
			return lr.SHIFT, 15, nil // SHIFT 15
		}

	case 36:
		switch a {
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 37:
		switch a {
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 38:
		switch a {
		case "<":
			return lr.SHIFT, 52, nil // SHIFT 52
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 39:
		switch a {
		case "@left":
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		case "@right":
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		case "@none":
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		case "IDENT":
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		case "TOKEN":
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		case grammar.Endmarker:
			return lr.REDUCE, 3, nil // REDUCE decls → ε
		}

	case 40:
		switch a {
		case ";":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "<":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "@left":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "@right":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "@none":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "IDENT":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "TOKEN":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case "STRING":
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		case grammar.Endmarker:
			return lr.REDUCE, 18, nil // REDUCE handles → rule_handle
		}

	case 41:
		switch a {
		case ";":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "<":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "@left":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "@right":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "@none":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "IDENT":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "TOKEN":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case "STRING":
			return lr.REDUCE, 17, nil // REDUCE handles → term
		case grammar.Endmarker:
			return lr.REDUCE, 17, nil // REDUCE handles → term
		}

	case 42:
		switch a {
		case "=":
			return lr.REDUCE, 22, nil // REDUCE lhs → nonterm
		}

	case 43:
		switch a {
		case "IDENT":
			return lr.SHIFT, 23, nil // SHIFT 23
		}

	case 44:
		switch a {
		case "=":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case ";":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "|":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "(":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case ")":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "[":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "]":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "{":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "}":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "{{":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "}}":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case ">":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "IDENT":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "TOKEN":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		case "STRING":
			return lr.REDUCE, 32, nil // REDUCE nonterm → "IDENT"
		}

	case 45:
		switch a {
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 46:
		switch a {
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 47:
		switch a {
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 48:
		switch a {
		case "(":
			return lr.SHIFT, 45, nil // SHIFT 45
		case "[":
			return lr.SHIFT, 46, nil // SHIFT 46
		case "{":
			return lr.SHIFT, 47, nil // SHIFT 47
		case "{{":
			return lr.SHIFT, 48, nil // SHIFT 48
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		case "TOKEN":
			return lr.SHIFT, 55, nil // SHIFT 55
		case "STRING":
			return lr.SHIFT, 54, nil // SHIFT 54
		}

	case 49:
		switch a {
		case ";":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "|":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "(":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case ")":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "[":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "]":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "{":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "}":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "{{":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "}}":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case ">":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "IDENT":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "TOKEN":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		case "STRING":
			return lr.REDUCE, 30, nil // REDUCE rhs → nonterm
		}

	case 50:
		switch a {
		case ";":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "|":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "(":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case ")":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "[":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "]":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "{":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "}":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "{{":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "}}":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case ">":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "IDENT":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "TOKEN":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		case "STRING":
			return lr.REDUCE, 31, nil // REDUCE rhs → term
		}

	case 51:
		switch a {
		case "=":
			return lr.SHIFT, 30, nil // SHIFT 30
		}

	case 52:
		switch a {
		case "IDENT":
			return lr.SHIFT, 44, nil // SHIFT 44
		}

	case 53:
		switch a {
		case "@left":
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		case "@right":
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		case "@none":
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		case "IDENT":
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		case "TOKEN":
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		case grammar.Endmarker:
			return lr.REDUCE, 7, nil // REDUCE semi_opt → ";"
		}

	case 54:
		switch a {
		case ";":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "|":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "(":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case ")":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "[":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "]":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "{":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "}":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "{{":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "}}":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "<":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case ">":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "@left":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "@right":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "@none":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "IDENT":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "TOKEN":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case "STRING":
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		case grammar.Endmarker:
			return lr.REDUCE, 34, nil // REDUCE term → "STRING"
		}

	case 55:
		switch a {
		case ";":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "|":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "(":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case ")":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "[":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "]":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "{":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "}":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "{{":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "}}":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "<":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case ">":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "@left":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "@right":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "@none":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "IDENT":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "TOKEN":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case "STRING":
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		case grammar.Endmarker:
			return lr.REDUCE, 33, nil // REDUCE term → "TOKEN"
		}

	case 56:
		switch a {
		case "=":
			return lr.SHIFT, 32, nil // SHIFT 32
		}

	}

	return lr.ERROR, -1, fmt.Errorf("no action exists in the parsing table for ACTION[%d, %s]", s, a)
}

// GOTO looks up and returns the next state for state s and non-terminal A.
func GOTO(s int, A grammar.NonTerminal) int {
	switch s {
	case 0:
		switch A {
		case "grammar":
			return 1
		case "name":
			return 39
		}

	case 3:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 8:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 17:
		switch A {
		case "rule_handle":
			return 21
		case "term":
			return 22
		}

	case 18:
		switch A {
		case "rule_handle":
			return 21
		case "term":
			return 22
		}

	case 19:
		switch A {
		case "rule_handle":
			return 21
		case "term":
			return 22
		}

	case 20:
		switch A {
		case "decl":
			return 16
		case "token":
			return 34
		case "directive":
			return 33
		case "rule":
			return 35
		case "lhs":
			return 51
		case "nonterm":
			return 42
		}

	case 23:
		switch A {
		case "semi_opt":
			return 2
		}

	case 24:
		switch A {
		case "rhs":
			return 3
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 25:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 26:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 27:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 28:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 29:
		switch A {
		case "rhs":
			return 25
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 30:
		switch A {
		case "rhs":
			return 8
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 33:
		switch A {
		case "semi_opt":
			return 13
		}

	case 34:
		switch A {
		case "semi_opt":
			return 14
		}

	case 36:
		switch A {
		case "handles":
			return 17
		case "rule_handle":
			return 40
		case "term":
			return 41
		}

	case 37:
		switch A {
		case "handles":
			return 18
		case "rule_handle":
			return 40
		case "term":
			return 41
		}

	case 38:
		switch A {
		case "handles":
			return 19
		case "rule_handle":
			return 40
		case "term":
			return 41
		}

	case 39:
		switch A {
		case "decls":
			return 20
		}

	case 45:
		switch A {
		case "rhs":
			return 26
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 46:
		switch A {
		case "rhs":
			return 27
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 47:
		switch A {
		case "rhs":
			return 28
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 48:
		switch A {
		case "rhs":
			return 29
		case "nonterm":
			return 49
		case "term":
			return 50
		}

	case 52:
		switch A {
		case "rule":
			return 31
		case "lhs":
			return 51
		case "nonterm":
			return 42
		}

	}

	return -1
}
