package golang

import (
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
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

var definitions = []*spec.TerminalDef{
	{Terminal: "ID", Kind: spec.RegexDef, Value: "[A-Za-z_][0-9A-Za-z_]*"},
	{Terminal: "NUM", Kind: spec.RegexDef, Value: "[0-9]+"},
}
