package spec

import (
	"testing"

	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"
)

func TestSpec_DFA(t *testing.T) {
	dfa := getDFA()

	tests := []struct {
		name        string
		s           *Spec
		expectedDFA *auto.DFA
	}{
		{
			name: "OK",
			s: &Spec{
				Definitions: []*TerminalDef{
					{Terminal: ";", DFA: dfa[0]},
					{Terminal: "ID", DFA: dfa[3]},
					{Terminal: "if", DFA: dfa[1]},
					{Terminal: "NUM", DFA: dfa[2]},
				},
			},
			expectedDFA: dfa[4],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dfa := tc.s.DFA()

			assert.True(t, dfa.Equal(tc.expectedDFA))
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
