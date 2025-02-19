package spec

import (
	"testing"

	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"
)

func TestSpec_DFA(t *testing.T) {
	dfa := getDFA()

	tests := []struct {
		name                    string
		s                       *Spec
		expectedDFA             *auto.DFA
		expectedTerminalMapping map[grammar.Terminal][]auto.State
		expectedErrorStrings    []string
	}{
		{
			name: "InvalidRegex",
			s: &Spec{
				Definitions: []*TerminalDef{
					{Terminal: "ID", Value: "[A-Z", IsRegex: true},
					{Terminal: "NUM", Value: "[0-9", IsRegex: true},
				},
			},
			expectedDFA:             nil,
			expectedTerminalMapping: nil,
			expectedErrorStrings: []string{
				`2 errors occurred:`,
				`"ID": invalid regular expression: [A-Z`,
				`"NUM": invalid regular expression: [0-9`,
			},
		},
		{
			name: "OverlappingDefinitions",
			s: &Spec{
				Definitions: []*TerminalDef{
					{Terminal: "NUM", Value: "[0-9]+", IsRegex: true, Pos: &lexer.Position{Filename: "test", Offset: 20, Line: 2, Column: 1}},
					{Terminal: "INT", Value: "[0-9]+", IsRegex: true, Pos: &lexer.Position{Filename: "test", Offset: 30, Line: 3, Column: 1}},
				},
			},
			expectedDFA:             nil,
			expectedTerminalMapping: nil,
			expectedErrorStrings: []string{
				`1 error occurred:`,
				`conflicting definitions capture the same string:`,
				`  test:2:1: "NUM"`,
				`  test:3:1: "INT"`,
			},
		},
		{
			name: "Success",
			s: &Spec{
				Definitions: []*TerminalDef{
					{Terminal: ";", Value: ";", IsRegex: false},
					{Terminal: "if", Value: "if", IsRegex: false},
					{Terminal: "ID", Value: "[A-Za-z_][0-9A-Za-z_]*", IsRegex: true},
					{Terminal: "NUM", Value: "[0-9]+", IsRegex: true},
				},
			},
			expectedDFA: dfa[4],
			expectedTerminalMapping: map[grammar.Terminal][]auto.State{
				"NUM": {1},
				";":   {2},
				"ID":  {3, 4},
				"if":  {5},
			},
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa, termMap, err := tc.s.DFA()

			if len(tc.expectedErrorStrings) == 0 {
				assert.True(t, dfa.Equal(tc.expectedDFA))
				assert.Equal(t, tc.expectedTerminalMapping, termMap)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, dfa)
				assert.Nil(t, termMap)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestSpec_Productions(t *testing.T) {
	tests := []struct {
		name                string
		s                   *Spec
		expectedProductions []*grammar.Production
	}{
		{
			name: "OK",
			s: &Spec{
				Grammar: grammars[0],
			},
			expectedProductions: []*grammar.Production{
				{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("*"), grammar.NonTerminal("E")}},
				{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")}},
				{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.Terminal("("), grammar.NonTerminal("E"), grammar.Terminal(")")}},
				{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.Terminal("id")}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prods := tc.s.Productions()

			assert.Len(t, prods, len(tc.expectedProductions))
			for i, expectedProduction := range tc.expectedProductions {
				assert.True(t, prods[i].Equal(expectedProduction))
			}
		})
	}
}

func TestSpec_SLRParsingTable(t *testing.T) {
	tests := []struct {
		name                 string
		s                    *Spec
		expectedErrorStrings []string
	}{
		{
			name: "Error",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: lr.PrecedenceLevels{},
			},
			expectedErrorStrings: []string{
				`error on building SLR(1) parsing table:`,
				`Error:      Ambiguous Grammar`,
				`Cause:      Multiple conflicts in the parsing table:`,
				`              1. Shift/Reduce conflict in ACTION[2, "*"]`,
				`              2. Shift/Reduce conflict in ACTION[2, "+"]`,
				`              3. Shift/Reduce conflict in ACTION[3, "*"]`,
				`              4. Shift/Reduce conflict in ACTION[3, "+"]`,
				`Resolution: Specify associativity and precedence for these Terminals/Productions:`,
				`              • "*" vs. "*", "+"`,
				`              • "+" vs. "*", "+"`,
				`            Terminals/Productions listed earlier will have higher precedence.`,
				`            Terminals/Productions in the same line will have the same precedence.`,
			},
		},
		{
			name: "Success",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: precedences[0],
			},
			expectedErrorStrings: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			T, err := tc.s.SLRParsingTable()

			if len(tc.expectedErrorStrings) == 0 {
				assert.NotNil(t, T)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, T)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestSpec_LALRParsingTable(t *testing.T) {
	tests := []struct {
		name                 string
		s                    *Spec
		expectedErrorStrings []string
	}{
		{
			name: "Error",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: lr.PrecedenceLevels{},
			},
			expectedErrorStrings: []string{
				`error on building LALR(1) parsing table:`,
				`Error:      Ambiguous Grammar`,
				`Cause:      Multiple conflicts in the parsing table:`,
				`              1. Shift/Reduce conflict in ACTION[2, "*"]`,
				`              2. Shift/Reduce conflict in ACTION[2, "+"]`,
				`              3. Shift/Reduce conflict in ACTION[3, "*"]`,
				`              4. Shift/Reduce conflict in ACTION[3, "+"]`,
				`Resolution: Specify associativity and precedence for these Terminals/Productions:`,
				`              • "*" vs. "*", "+"`,
				`              • "+" vs. "*", "+"`,
				`            Terminals/Productions listed earlier will have higher precedence.`,
				`            Terminals/Productions in the same line will have the same precedence.`,
			},
		},
		{
			name: "Success",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: precedences[0],
			},
			expectedErrorStrings: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			T, err := tc.s.LALRParsingTable()

			if len(tc.expectedErrorStrings) == 0 {
				assert.NotNil(t, T)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, T)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestSpec_GLRParsingTable(t *testing.T) {
	tests := []struct {
		name                 string
		s                    *Spec
		expectedErrorStrings []string
	}{
		{
			name: "Error",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: lr.PrecedenceLevels{},
			},
			expectedErrorStrings: []string{
				`error on building GLR(1) parsing table:`,
				`Error:      Ambiguous Grammar`,
				`Cause:      Multiple conflicts in the parsing table:`,
				`              1. Shift/Reduce conflict in ACTION[2, "*"]`,
				`              2. Shift/Reduce conflict in ACTION[2, "+"]`,
				`              3. Shift/Reduce conflict in ACTION[3, "*"]`,
				`              4. Shift/Reduce conflict in ACTION[3, "+"]`,
				`Resolution: Specify associativity and precedence for these Terminals/Productions:`,
				`              • "*" vs. "*", "+"`,
				`              • "+" vs. "*", "+"`,
				`            Terminals/Productions listed earlier will have higher precedence.`,
				`            Terminals/Productions in the same line will have the same precedence.`,
			},
		},
		{
			name: "Success",
			s: &Spec{
				Grammar:     grammars[0],
				Precedences: precedences[0],
			},
			expectedErrorStrings: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			T, err := tc.s.GLRParsingTable()

			if len(tc.expectedErrorStrings) == 0 {
				assert.NotNil(t, T)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, T)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}
