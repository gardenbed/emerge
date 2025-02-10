package generator

import (
	"testing"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/stretchr/testify/assert"
)

func TestNewSymbolTable(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		st := NewSymbolTable()

		assert.NotNil(t, st.precedences.list)
		assert.NotNil(t, st.tokenDefs.strings)
		assert.NotNil(t, st.tokenDefs.regexes)
		assert.NotNil(t, st.terminals.strings)
		assert.NotNil(t, st.terminals.tokens)
		assert.NotNil(t, st.nonTerminals.table)
		assert.NotNil(t, st.productions.table)
		assert.NotNil(t, st.strings.table)
	})
}

func TestSymbolTable_Verify(t *testing.T) {
	st0 := NewSymbolTable()
	st0.AddTokenTerminal("QUOT", &lexer.Position{Filename: "test", Offset: 30, Line: 4, Column: 10})

	st1 := NewSymbolTable()
	st1.AddStringTokenDef("QUOT", "'", &lexer.Position{Filename: "test", Offset: 10, Line: 2, Column: 1})
	st1.AddStringTokenDef("QUOT", "\"", &lexer.Position{Filename: "test", Offset: 20, Line: 3, Column: 1})

	st2 := NewSymbolTable()
	st2.AddStringTokenDef("QUOT", "\"", &lexer.Position{Filename: "test", Offset: 10, Line: 2, Column: 1})
	st2.AddRegexTokenDef("NUM", "[0-9]+", &lexer.Position{Filename: "test", Offset: 20, Line: 3, Column: 1})
	st2.AddTokenTerminal("QUOT", &lexer.Position{Filename: "test", Offset: 30, Line: 4, Column: 10})
	st2.AddTokenTerminal("NUM", &lexer.Position{Filename: "test", Offset: 40, Line: 5, Column: 12})

	tests := []struct {
		name                 string
		st                   *SymbolTable
		expectedErrorStrings []string
	}{
		{
			name: "NoDefinition",
			st:   st0,
			expectedErrorStrings: []string{
				`1 error occurred:`,
				`no definition for terminal "QUOT"`,
			},
		},
		{
			name: "MultipleDefinitions",
			st:   st1,
			expectedErrorStrings: []string{
				`1 error occurred:`,
				`multiple definitions for terminal "QUOT":`,
				`test:2:1`,
				`test:3:1`,
			},
		},
		{
			name:                 "OK",
			st:                   st2,
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.st.Verify()

			if len(tc.expectedErrorStrings) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}

func TestSymbolTable_Precedences(t *testing.T) {
	st := NewSymbolTable()
	st.AddPrecedence(&lr.PrecedenceLevel{
		Associativity: lr.LEFT,
		Handles: lr.NewPrecedenceHandles(
			lr.PrecedenceHandleForTerminal("*"),
			lr.PrecedenceHandleForTerminal("/"),
		),
	})

	tests := []struct {
		name                string
		st                  *SymbolTable
		expectedPrecedences lr.PrecedenceLevels
	}{
		{
			name: "OK",
			st:   st,
			expectedPrecedences: lr.PrecedenceLevels{
				&lr.PrecedenceLevel{
					Associativity: lr.LEFT,
					Handles: lr.NewPrecedenceHandles(
						lr.PrecedenceHandleForTerminal("*"),
						lr.PrecedenceHandleForTerminal("/"),
					),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			precs := tc.st.Precedences()

			assert.Len(t, precs, len(tc.expectedPrecedences))
			for i, expectedPrecedence := range tc.expectedPrecedences {
				assert.True(t, precs[i].Equal(expectedPrecedence))
			}
		})
	}
}

func TestSymbolTable_Terminals(t *testing.T) {
	st := NewSymbolTable()
	st.AddStringTerminal(";", &lexer.Position{})
	st.AddTokenTerminal("ID", &lexer.Position{})

	tests := []struct {
		name              string
		st                *SymbolTable
		expectedTerminals []grammar.Terminal
	}{
		{
			name:              "OK",
			st:                st,
			expectedTerminals: []grammar.Terminal{";", "ID"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			terms := tc.st.Terminals()

			assert.Len(t, terms, len(tc.expectedTerminals))
			for _, expectedTerminal := range tc.expectedTerminals {
				assert.Contains(t, terms, expectedTerminal)
			}
		})
	}
}

func TestSymbolTable_NonTerminals(t *testing.T) {
	st := NewSymbolTable()
	st.AddNonTerminal("expr", &lexer.Position{})
	st.AddNonTerminal("stmt", &lexer.Position{})

	tests := []struct {
		name                 string
		st                   *SymbolTable
		expectedNonTerminals []grammar.NonTerminal
	}{
		{
			name:                 "OK",
			st:                   st,
			expectedNonTerminals: []grammar.NonTerminal{"expr", "stmt"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nonTerms := tc.st.NonTerminals()

			assert.Len(t, nonTerms, len(tc.expectedNonTerminals))
			for _, expectedNonTerminal := range tc.expectedNonTerminals {
				assert.Contains(t, nonTerms, expectedNonTerminal)
			}
		})
	}
}

func TestSymbolTable_Productions(t *testing.T) {
	st := NewSymbolTable()
	st.AddProduction(
		&grammar.Production{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")}},
		&lexer.Position{},
	)

	tests := []struct {
		name                string
		st                  *SymbolTable
		expectedProductions []*grammar.Production
	}{
		{
			name: "OK",
			st:   st,
			expectedProductions: []*grammar.Production{
				{Head: "E", Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prods := tc.st.Productions()

			assert.Len(t, prods, len(tc.expectedProductions))
			for i, expectedProduction := range tc.expectedProductions {
				assert.Equal(t, expectedProduction, prods[i])
			}
		})
	}
}

func TestSymbolTable_AddPrecedence(t *testing.T) {
	tests := []struct {
		name string
		st   *SymbolTable
		p    *lr.PrecedenceLevel
	}{
		{
			name: "OK",
			st:   NewSymbolTable(),
			p: &lr.PrecedenceLevel{
				Associativity: lr.LEFT,
				Handles: lr.NewPrecedenceHandles(
					lr.PrecedenceHandleForTerminal("*"),
					lr.PrecedenceHandleForTerminal("/"),
				),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddPrecedence(tc.p)

			l := len(tc.st.precedences.list) - 1
			assert.True(t, tc.st.precedences.list[l].Equal(tc.p))
		})
	}
}

func TestSymbolTable_AddStringTokenDef(t *testing.T) {
	tests := []struct {
		name  string
		st    *SymbolTable
		token grammar.Terminal
		value string
		pos   *lexer.Position
	}{
		{
			name:  "OK",
			st:    NewSymbolTable(),
			token: "QUOT",
			value: "\"",
			pos:   &lexer.Position{Filename: "test", Offset: 10, Line: 2, Column: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddStringTokenDef(tc.token, tc.value, tc.pos)

			l := len(tc.st.tokenDefs.strings) - 1
			entry := tc.st.tokenDefs.strings[l]

			assert.True(t, entry.token.Equal(tc.token))
			assert.Equal(t, tc.value, entry.value)
			assert.True(t, entry.occurrence.Equal(*tc.pos))
		})
	}
}

func TestSymbolTable_AddRegexTokenDef(t *testing.T) {
	tests := []struct {
		name  string
		st    *SymbolTable
		token grammar.Terminal
		value string
		pos   *lexer.Position
	}{
		{
			name:  "OK",
			st:    NewSymbolTable(),
			token: "NUM",
			value: "[0-9]+",
			pos:   &lexer.Position{Filename: "test", Offset: 20, Line: 3, Column: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddRegexTokenDef(tc.token, tc.value, tc.pos)

			l := len(tc.st.tokenDefs.regexes) - 1
			entry := tc.st.tokenDefs.regexes[l]

			assert.True(t, entry.token.Equal(tc.token))
			assert.Equal(t, tc.value, entry.value)
			assert.True(t, entry.occurrence.Equal(*tc.pos))
		})
	}
}

func TestSymbolTable_AddStringTerminal(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name string
		st   *SymbolTable
		a    grammar.Terminal
		pos  *lexer.Position
	}{
		{
			name: "FirstOccurrence",
			st:   st,
			a:    ";",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   50,
				Line:     5,
				Column:   10,
			},
		},
		{
			name: "SecondOccurrence",
			st:   st,
			a:    ";",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   64,
				Line:     6,
				Column:   12,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddStringTerminal(tc.a, tc.pos)

			e, ok := tc.st.terminals.strings.Get(tc.a)
			assert.True(t, ok)
			assert.NotZero(t, e.index)
			assert.Contains(t, e.occurrences, tc.pos)
		})
	}
}

func TestSymbolTable_AddTokenTerminal(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name string
		st   *SymbolTable
		a    grammar.Terminal
		pos  *lexer.Position
	}{
		{
			name: "FirstOccurrence",
			st:   st,
			a:    "ID",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   50,
				Line:     5,
				Column:   8,
			},
		},
		{
			name: "SecondOccurrence",
			st:   st,
			a:    "ID",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   64,
				Line:     6,
				Column:   16,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddTokenTerminal(tc.a, tc.pos)

			e, ok := tc.st.terminals.tokens.Get(tc.a)
			assert.True(t, ok)
			assert.NotZero(t, e.index)
			assert.Contains(t, e.occurrences, tc.pos)
		})
	}
}

func TestSymbolTable_AddNonTerminal(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name string
		st   *SymbolTable
		A    grammar.NonTerminal
		pos  *lexer.Position
	}{
		{
			name: "FirstOccurrence",
			st:   st,
			A:    "expr",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   40,
				Line:     4,
				Column:   1,
			},
		},
		{
			name: "SecondOccurrence",
			st:   st,
			A:    "expr",
			pos: &lexer.Position{
				Filename: "test",
				Offset:   49,
				Line:     4,
				Column:   10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddNonTerminal(tc.A, tc.pos)

			e, ok := tc.st.nonTerminals.table.Get(tc.A)
			assert.True(t, ok)
			assert.NotZero(t, e.index)
			assert.Contains(t, e.occurrences, tc.pos)
		})
	}
}

func TestSymbolTable_AddProduction(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name string
		st   *SymbolTable
		p    *grammar.Production
		pos  *lexer.Position
	}{
		{
			name: "FirstOccurrence",
			st:   st,
			p: &grammar.Production{
				Head: "E",
				Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")},
			},
			pos: &lexer.Position{
				Filename: "test",
				Offset:   40,
				Line:     4,
				Column:   1,
			},
		},
		{
			name: "SecondOccurrence",
			st:   st,
			p: &grammar.Production{
				Head: "E",
				Body: grammar.String[grammar.Symbol]{grammar.NonTerminal("E"), grammar.Terminal("+"), grammar.NonTerminal("E")},
			},
			pos: &lexer.Position{
				Filename: "test",
				Offset:   80,
				Line:     8,
				Column:   1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.st.AddProduction(tc.p, tc.pos)

			e, ok := tc.st.productions.table.Get(tc.p)
			assert.True(t, ok)
			assert.NotZero(t, e.index)
			assert.Contains(t, e.occurrences, tc.pos)
		})
	}
}

func TestSymbolTable_GetOpt(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name                string
		st                  *SymbolTable
		s                   Strings
		expectedNonTerminal grammar.NonTerminal
	}{
		{
			name: "New",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal(";")},
			},
			expectedNonTerminal: "gen_semi_opt",
		},
		{
			name: "Existent",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal(";")},
			},
			expectedNonTerminal: "gen_semi_opt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNonTerminal, tc.st.GetOpt(tc.s))
		})
	}
}

func TestSymbolTable_GetGroup(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name                string
		st                  *SymbolTable
		s                   Strings
		expectedNonTerminal grammar.NonTerminal
	}{
		{
			name: "New",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("BOOLEAN")},
				grammar.String[grammar.Symbol]{grammar.Terminal("INTEGER")},
				grammar.String[grammar.Symbol]{grammar.Terminal("REAL")},
			},
			expectedNonTerminal: "gen1_group",
		},
		{
			name: "Existent",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.Terminal("BOOLEAN")},
				grammar.String[grammar.Symbol]{grammar.Terminal("INTEGER")},
				grammar.String[grammar.Symbol]{grammar.Terminal("REAL")},
			},
			expectedNonTerminal: "gen1_group",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNonTerminal, tc.st.GetGroup(tc.s))
		})
	}
}

func TestSymbolTable_GetStar(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name                string
		st                  *SymbolTable
		s                   Strings
		expectedNonTerminal grammar.NonTerminal
	}{
		{
			name: "New",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.NonTerminal("decl")},
			},
			expectedNonTerminal: "gen_decl_star",
		},
		{
			name: "Existent",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.NonTerminal("decl")},
			},
			expectedNonTerminal: "gen_decl_star",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNonTerminal, tc.st.GetStar(tc.s))
		})
	}
}

func TestSymbolTable_GetPlus(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name                string
		st                  *SymbolTable
		s                   Strings
		expectedNonTerminal grammar.NonTerminal
	}{
		{
			name: "New",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.NonTerminal("op")},
			},
			expectedNonTerminal: "gen_op_plus",
		},
		{
			name: "Existent",
			st:   st,
			s: Strings{
				grammar.String[grammar.Symbol]{grammar.NonTerminal("op")},
			},
			expectedNonTerminal: "gen_op_plus",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedNonTerminal, tc.st.GetPlus(tc.s))
		})
	}
}
