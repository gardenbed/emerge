package spec

import (
	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
)

var grammars = []*grammar.CFG{
	// G0
	grammar.NewCFG(
		[]grammar.Terminal{
			"(", ")", "=", "+", "-", "*", "/", "!", "|", "&", "^", "==", "!=", "<", ">", "<=", ">=",
			"int", "float", "void", "if", "else", "OR", "AND", "XOR", "SEMI", "ID", "NUMBER",
		},
		[]grammar.NonTerminal{
			"start", "decl", "type", "stmt", "assign", "if_stmt", "else_stmt", "expr", "bitop", "logop", "empty",
			"gen_decl_star", "gen_stmt_plus", "gen1_opt", "gen2_group",
		},
		[]*grammar.Production{
			{Head: "start", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("gen_decl_star"), grammar.NonTerminal("gen_stmt_plus")}},
			{Head: "gen_decl_star", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("gen_decl_star"), grammar.NonTerminal("decl")}},
			{Head: "gen_decl_star", Body: grammar.E},
			{Head: "gen_stmt_plus", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("gen_stmt_plus"), grammar.NonTerminal("stmt")}},
			{Head: "gen_stmt_plus", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("stmt")}},
			{Head: "decl", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("type"), grammar.Terminal("ID"), grammar.NonTerminal("gen1_opt"), grammar.Terminal("SEMI")}},
			{Head: "type", Body: grammar.String[grammar.Symbol]{grammar.Terminal("int")}},
			{Head: "type", Body: grammar.String[grammar.Symbol]{grammar.Terminal("float")}},
			{Head: "type", Body: grammar.String[grammar.Symbol]{grammar.Terminal("void")}},
			{Head: "gen1_opt", Body: grammar.String[grammar.Symbol]{grammar.Terminal("="), grammar.NonTerminal("expr")}},
			{Head: "gen1_opt", Body: grammar.E},
			{Head: "stmt", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("assign")}},
			{Head: "stmt", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("if_stmt")}},
			{Head: "assign", Body: grammar.String[grammar.Symbol]{grammar.Terminal("ID"), grammar.Terminal("="), grammar.NonTerminal("expr")}},
			{Head: "if_stmt", Body: grammar.String[grammar.Symbol]{grammar.Terminal("if"), grammar.Terminal("("), grammar.NonTerminal("expr"), grammar.Terminal(")"), grammar.NonTerminal("else_stmt")}},
			{Head: "else_stmt", Body: grammar.String[grammar.Symbol]{grammar.Terminal("else"), grammar.NonTerminal("stmt")}},
			{Head: "else_stmt", Body: grammar.E},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("expr"), grammar.NonTerminal("gen2_group"), grammar.NonTerminal("expr")}},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("expr"), grammar.NonTerminal("bitop"), grammar.NonTerminal("expr")}},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("expr"), grammar.NonTerminal("logop"), grammar.NonTerminal("expr")}},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.Terminal("!"), grammar.NonTerminal("expr")}},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.Terminal("NUMBER")}},
			{Head: "expr", Body: grammar.String[grammar.Symbol]{grammar.Terminal("ID")}},
			{Head: "gen2_group", Body: grammar.String[grammar.Symbol]{grammar.Terminal("+")}},
			{Head: "gen2_group", Body: grammar.String[grammar.Symbol]{grammar.Terminal("-")}},
			{Head: "gen2_group", Body: grammar.String[grammar.Symbol]{grammar.Terminal("*")}},
			{Head: "gen2_group", Body: grammar.String[grammar.Symbol]{grammar.Terminal("/")}},
			{Head: "bitop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("|")}},
			{Head: "bitop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("&")}},
			{Head: "bitop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("^")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("==")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("!=")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("<")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal(">")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("<=")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal(">=")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("OR")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("AND")}},
			{Head: "logop", Body: grammar.String[grammar.Symbol]{grammar.Terminal("XOR")}},
			{Head: "empty", Body: grammar.E},
		},
		"start",
	),
}

var precedences = []lr.PrecedenceLevels{
	{ // 0
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("*"),
				lr.PrecedenceHandleForTerminal("/"),
			),
		},
		{
			Associativity: lr.LEFT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("+"),
				lr.PrecedenceHandleForTerminal("-"),
			),
		},
		{
			Associativity: lr.RIGHT,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForProduction(&grammar.Production{
					Head: "expr",
					Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("expr"), grammar.NonTerminal("bitop"), grammar.NonTerminal("expr")},
				}),
				lr.PrecedenceHandleForProduction(&grammar.Production{
					Head: "expr",
					Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("expr"), grammar.NonTerminal("logop"), grammar.NonTerminal("expr")},
				}),
			),
		},
		{
			Associativity: lr.NONE,
			Handles: lr.NewPrecedenceHandles(
				lr.PrecedenceHandleForTerminal("="),
			),
		},
	},
}

func getDFA() []*auto.DFA {
	d0 := auto.NewDFA(0, []auto.State{1})
	d0.Add(0, ';', 1)

	d1 := auto.NewDFA(0, []auto.State{1})
	d1.Add(0, '0', 1)
	d1.Add(0, '1', 1)
	d1.Add(0, '2', 1)
	d1.Add(0, '3', 1)
	d1.Add(0, '4', 1)
	d1.Add(0, '5', 1)
	d1.Add(0, '6', 1)
	d1.Add(0, '7', 1)
	d1.Add(0, '8', 1)
	d1.Add(0, '9', 1)
	d1.Add(1, '0', 1)
	d1.Add(1, '1', 1)
	d1.Add(1, '2', 1)
	d1.Add(1, '3', 1)
	d1.Add(1, '4', 1)
	d1.Add(1, '5', 1)
	d1.Add(1, '6', 1)
	d1.Add(1, '7', 1)
	d1.Add(1, '8', 1)
	d1.Add(1, '9', 1)

	return []*auto.DFA{d0, d1}
}
