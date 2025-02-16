package command

import (
	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
)

var grammars = []*grammar.CFG{
	// G0
	grammar.NewCFG(
		[]grammar.Terminal{"+", "*", "(", ")", "id"},
		[]grammar.NonTerminal{"E"},
		[]*grammar.Production{
			{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")}}, // E → E + E
			{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("*"), grammar.NonTerminal("E")}}, // E → E * E
			{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.Terminal("("), grammar.NonTerminal("E"), grammar.Terminal(")")}},    // E → ( E )
			{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.Terminal("id")}},                                                    // E → id
		},
		"E",
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
	},
}

func getDFA() []*auto.DFA {
	d0 := auto.NewDFA(0, []auto.State{1})
	d0.Add(0, ';', 1)

	d1 := auto.NewDFA(0, []auto.State{2})
	d1.Add(0, 'i', 1)
	d1.Add(1, 'f', 2)

	d2 := auto.NewDFA(0, []auto.State{1})
	for _, r := range "0123456789" {
		d2.Add(0, auto.Symbol(r), 1)
		d2.Add(1, auto.Symbol(r), 1)
	}

	d3 := auto.NewDFA(0, []auto.State{1})
	for _, r := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz" {
		if r >= 'A' {
			d3.Add(0, auto.Symbol(r), 1)
		}
		d3.Add(1, auto.Symbol(r), 1)
	}

	d4 := auto.NewDFA(0, []auto.State{1, 2, 3, 4, 5})

	for _, r := range "0123456789" {
		d4.Add(0, auto.Symbol(r), 1)
		d4.Add(1, auto.Symbol(r), 1)
	}

	d4.Add(0, ';', 2)

	for _, r := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz" {
		if r == 'i' {
			d4.Add(0, auto.Symbol(r), 4)
		} else if r >= 'A' {
			d4.Add(0, auto.Symbol(r), 3)
		}

		if r == 'f' {
			d4.Add(4, auto.Symbol(r), 5)
		} else {
			d4.Add(4, auto.Symbol(r), 3)
		}

		d4.Add(3, auto.Symbol(r), 3)
		d4.Add(5, auto.Symbol(r), 3)
	}

	return []*auto.DFA{d0, d1, d2, d3, d4}
}
