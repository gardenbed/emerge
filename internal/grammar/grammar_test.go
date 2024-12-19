package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var grammars = []Grammar{
	New(
		[]Terminal{"0", "1"},
		[]NonTerminal{"S", "X", "Y"},
		[]Production{
			{"S", String[Symbol]{NonTerminal("X"), NonTerminal("Y"), NonTerminal("X")}}, // S → XYX
			{"X", String[Symbol]{Terminal("0"), NonTerminal("X")}},                      // X → 0X
			{"X", ε}, // X → ε
			{"Y", String[Symbol]{Terminal("1"), NonTerminal("Y")}}, // Y → 1Y
			{"Y", ε}, // Y → ε
		},
		"S",
	),
	New(
		[]Terminal{"a", "b"},
		[]NonTerminal{"S"},
		[]Production{
			{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S → aSbS
			{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S → bSaS
			{"S", ε}, // S → ε
		},
		"S",
	),
	New(
		[]Terminal{"a", "b"},
		[]NonTerminal{"S", "A", "B"},
		[]Production{
			{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
			{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
			{"S", String[Symbol]{Terminal("a")}},                                  // S → a
			{"A", String[Symbol]{Terminal("b")}},                                  // A → b
			{"A", ε},                                                              // A → ε
			{"B", String[Symbol]{NonTerminal("A")}},                               // B → A
			{"B", String[Symbol]{Terminal("b")}},                                  // B → b
		},
		"S",
	),
	New(
		[]Terminal{"b", "c", "d", "s"},
		[]NonTerminal{"S", "A", "B", "C", "D"},
		[]Production{
			{"S", String[Symbol]{NonTerminal("A")}}, // S → A
			{"S", String[Symbol]{Terminal("s")}},    // S → s
			{"A", String[Symbol]{NonTerminal("B")}}, // A → B
			{"B", String[Symbol]{NonTerminal("C")}}, // B → C
			{"B", String[Symbol]{Terminal("b")}},    // B → b
			{"C", String[Symbol]{NonTerminal("D")}}, // C → D
			{"D", String[Symbol]{Terminal("d")}},    // D → d
		},
		"S",
	),
	New(
		[]Terminal{"a", "b", "c", "d"},
		[]NonTerminal{"S", "A", "B", "C", "D"},
		[]Production{
			{"S", String[Symbol]{NonTerminal("A"), NonTerminal("B")}}, // S → AB
			{"A", String[Symbol]{Terminal("a"), NonTerminal("A")}},    // A → aA
			{"A", String[Symbol]{Terminal("a")}},                      // A → a
			{"B", String[Symbol]{Terminal("b"), NonTerminal("B")}},    // B → bB
			{"B", String[Symbol]{Terminal("b")}},                      // B → b
			{"C", String[Symbol]{Terminal("c"), NonTerminal("C")}},    // C → cC
			{"C", String[Symbol]{Terminal("c")}},                      // C → c
			{"D", String[Symbol]{Terminal("d")}},                      // D → d
		},
		"S",
	),
	New(
		[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
		[]NonTerminal{"E", "S"},
		[]Production{
			{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
			{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // E → E + E
			{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // E → E - E
			{"E", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // E → E * E
			{"E", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // E → E / E
			{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
			{"E", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // E → - E
			{"E", String[Symbol]{Terminal("id")}},                                    // E → id
		},
		"S",
	),
	New(
		[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
		[]NonTerminal{"E", "T", "F", "S"},
		[]Production{
			{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
			{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
			{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
			{"E", String[Symbol]{NonTerminal("T")}},                                  // E → T
			{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
			{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
			{"T", String[Symbol]{NonTerminal("F")}},                                  // T → F
			{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
			{"F", String[Symbol]{Terminal("id")}},                                    // F → id
		},
		"S",
	),
	New(
		[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
		[]NonTerminal{"grammar", "name", "decls", "decl", "token", "rule", "lhs", "rhs", "nonterm", "term"},
		[]Production{
			{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}}, // grammar → name decls
			{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},       // name → GRAMMAR IDENT
			{"decls", String[Symbol]{NonTerminal("decls"), NonTerminal("decl")}},   // decls → decls decl
			{"decls", ε}, // decls → ε
			{"decl", String[Symbol]{NonTerminal("token")}},                                  // decl → token
			{"decl", String[Symbol]{NonTerminal("rule")}},                                   // decl → rule
			{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}}, // token → TOKEN "=" STRING
			{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},  // token → TOKEN "=" REGEX
			{"rule", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // rule → lhs "=" rhs
			{"lhs", String[Symbol]{NonTerminal("nonterm")}},                                 // lhs → nonterm
			{"rhs", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs")}},                 // rhs → rhs rhs
			{"rhs", String[Symbol]{NonTerminal("rhs"), Terminal("|"), NonTerminal("rhs")}},  // rhs → rhs "|" rhs
			{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")")}},       // rhs → "(" rhs ")"
			{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]")}},       // rhs → "[" rhs "]"
			{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}")}},       // rhs → "{" rhs "}"
			{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}")}},     // rhs → "{{" rhs "}}"
			{"rhs", String[Symbol]{NonTerminal("nonterm")}},                                 // rhs → nonterm
			{"rhs", String[Symbol]{NonTerminal("term")}},                                    // rhs → term
			{"nonterm", String[Symbol]{Terminal("IDENT")}},                                  // nonterm → IDENT
			{"term", String[Symbol]{Terminal("TOKEN")}},                                     // term → TOKEN
			{"term", String[Symbol]{Terminal("STRING")}},                                    // term → STRING
		},
		"grammar",
	),
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Terminal
		nonTerms []NonTerminal
		prods    []Production
		start    NonTerminal
	}{
		{
			name:     "MatchingPairs",
			terms:    []Terminal{"a", "b"},
			nonTerms: []NonTerminal{"S"},
			prods: []Production{
				{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}}, //  S → aSb
				{"S", ε}, //  S → ε
			},
			start: "S",
		},
		{
			name:     "WellformedParantheses",
			terms:    []Terminal{"(", ")"},
			nonTerms: []NonTerminal{"S"},
			prods: []Production{
				{"S", String[Symbol]{NonTerminal("S"), NonTerminal("S")}},             //  S → SS
				{"S", String[Symbol]{Terminal("("), NonTerminal("S"), Terminal(")")}}, //  S → (S)
				{"S", String[Symbol]{Terminal("("), Terminal(")")}},                   //  S → ()
			},
			start: "S",
		},
		{
			name:     "WellformedParanthesesAndBrackets",
			terms:    []Terminal{"(", ")", "[", "]"},
			nonTerms: []NonTerminal{"S"},
			prods: []Production{
				{"S", String[Symbol]{NonTerminal("S"), NonTerminal("S")}},             //  S → SS
				{"S", String[Symbol]{Terminal("("), NonTerminal("S"), Terminal(")")}}, //  S → (S)
				{"S", String[Symbol]{Terminal("["), NonTerminal("S"), Terminal("]")}}, //  S → [S]
				{"S", String[Symbol]{Terminal("("), Terminal(")")}},                   //  S → ()
				{"S", String[Symbol]{Terminal("["), Terminal("]")}},                   //  S → []
			},
			start: "S",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := New(tc.terms, tc.nonTerms, tc.prods, tc.start)
			assert.NotEmpty(t, g)
		})
	}
}

func TestGrammar_Verify(t *testing.T) {
	tests := []struct {
		name          string
		g             Grammar
		expectedError string
	}{
		{
			name: "StartSymbolNotDeclared",
			g: New(
				[]Terminal{},
				[]NonTerminal{},
				[]Production{},
				"S",
			),
			expectedError: "start symbol S not in the set of non-terminal symbols\nno production rule for start symbol S",
		},
		{
			name: "StartSymbolHasNoProduction",
			g: New(
				[]Terminal{},
				[]NonTerminal{"S"},
				[]Production{},
				"S",
			),
			expectedError: "no production rule for start symbol S\nno production rule for non-terminal symbol S",
		},
		{
			name: "NonTerminalHasNoProduction",
			g: New(
				[]Terminal{},
				[]NonTerminal{"A", "S"},
				[]Production{
					{"S", ε}, // S → ε
				},
				"S",
			),
			expectedError: "no production rule for non-terminal symbol A",
		},
		{
			name: "ProductionHeadNotDeclared",
			g: New(
				[]Terminal{},
				[]NonTerminal{"A", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A")}}, // S → A
					{"A", ε},                                // A → ε
					{"B", ε},                                // B → ε
				},
				"S",
			),
			expectedError: "production head B not in the set of non-terminal symbols",
		},
		{
			name: "TerminalNotDeclared",
			g: New(
				[]Terminal{},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A")}}, // S → A
					{"A", String[Symbol]{Terminal("a")}},    // A → a
					{"B", ε},                                // B → ε
				},
				"S",
			),
			expectedError: "terminal symbol \"a\" not in the set of terminal symbols",
		},
		{
			name: "NonTerminalNotDeclared",
			g: New(
				[]Terminal{"a"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A")}}, // S → A
					{"A", String[Symbol]{Terminal("a")}},    // A → a
					{"B", String[Symbol]{NonTerminal("C")}}, // B → C
				},
				"S",
			),
			expectedError: "non-terminal symbol C not in the set of non-terminal symbols",
		},
		{
			name: "Valid",
			g: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A")}}, // S → A
					{"S", String[Symbol]{NonTerminal("B")}}, // S → B
					{"A", String[Symbol]{Terminal("a")}},    // A → a
					{"B", String[Symbol]{Terminal("b")}},    // B → b
				},
				"S",
			),
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.Verify()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGrammar_Equals(t *testing.T) {
	tests := []struct {
		name           string
		lhs            Grammar
		rhs            Grammar
		expectedEquals bool
	}{
		{
			name: "TerminalsNotEqual",
			lhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{},
				"S",
			),
			rhs: New(
				[]Terminal{"a", "b", "c"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{},
				"S",
			),
			expectedEquals: false,
		},
		{
			name: "NonTerminalsNotEqual",
			lhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "C", "S"},
				[]Production{},
				"S",
			),
			rhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{},
				"S",
			),
			expectedEquals: false,
		},
		{
			name: "ProductionsNotEqual",
			lhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("A")}}, // S → aA
					{"S", String[Symbol]{Terminal("b"), NonTerminal("B")}}, // S → bB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("S")}}, // A → aS
					{"A", String[Symbol]{Terminal("b"), NonTerminal("A")}}, // A → bA
					{"A", ε}, // A → ε
					{"B", String[Symbol]{Terminal("b"), NonTerminal("S")}}, // B → bS
					{"B", String[Symbol]{Terminal("a"), NonTerminal("B")}}, // B → aB
					{"B", ε}, // B → ε
				},
				"S",
			),
			rhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("A")}}, // S → aA
					{"S", String[Symbol]{Terminal("b"), NonTerminal("B")}}, // S → bB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("S")}}, // A → aS
					{"A", String[Symbol]{Terminal("b"), NonTerminal("A")}}, // A → bA
					{"B", String[Symbol]{Terminal("b"), NonTerminal("S")}}, // B → bS
					{"B", String[Symbol]{Terminal("a"), NonTerminal("B")}}, // B → aB
					{"B", ε}, // B → ε
				},
				"S",
			),
			expectedEquals: false,
		},
		{
			name: "StartSymbolsNotEqual",
			lhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("A")}}, // S → aA
					{"S", String[Symbol]{Terminal("b"), NonTerminal("B")}}, // S → bB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("S")}}, // A → aS
					{"A", String[Symbol]{Terminal("b"), NonTerminal("A")}}, // A → bA
					{"A", ε}, // A → ε
					{"B", String[Symbol]{Terminal("b"), NonTerminal("S")}}, // B → bS
					{"B", String[Symbol]{Terminal("a"), NonTerminal("B")}}, // B → aB
					{"B", ε}, // B → ε
				},
				"S",
			),
			rhs: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"A", "B", "S"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("A")}}, // S → aA
					{"S", String[Symbol]{Terminal("b"), NonTerminal("B")}}, // S → bB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("S")}}, // A → aS
					{"A", String[Symbol]{Terminal("b"), NonTerminal("A")}}, // A → bA
					{"A", ε}, // A → ε
					{"B", String[Symbol]{Terminal("b"), NonTerminal("S")}}, // B → bS
					{"B", String[Symbol]{Terminal("a"), NonTerminal("B")}}, // B → aB
					{"B", ε}, // B → ε
				},
				"A",
			),
			expectedEquals: false,
		},
		{
			name: "Equal",
			lhs: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"S", "E", "T", "F"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"E", String[Symbol]{NonTerminal("T")}},                                  // E → T
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
					{"T", String[Symbol]{NonTerminal("F")}},                                  // T → F
					{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
					{"F", String[Symbol]{Terminal("id")}},                                    // F → id
				},
				"S",
			),
			rhs: New(
				[]Terminal{"id", "(", ")", "+", "-", "*", "/"},
				[]NonTerminal{"F", "T", "E", "S"},
				[]Production{
					{"F", String[Symbol]{Terminal("id")}},                                    // F → id
					{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
					{"T", String[Symbol]{NonTerminal("F")}},                                  // T → F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
					{"E", String[Symbol]{NonTerminal("T")}},                                  // E → T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
				},
				"S",
			),
			expectedEquals: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedEquals, tc.lhs.Equals(tc.rhs))
		})
	}
}

func TestGrammar_nullableNonTerminals(t *testing.T) {
	tests := []struct {
		name              string
		g                 Grammar
		expectedNullables []NonTerminal
	}{
		{
			name:              "1st",
			g:                 grammars[0],
			expectedNullables: []NonTerminal{"S", "X", "Y"},
		},
		{
			name:              "2nd",
			g:                 grammars[1],
			expectedNullables: []NonTerminal{"S"},
		},
		{
			name:              "3rd",
			g:                 grammars[2],
			expectedNullables: []NonTerminal{"A", "B"},
		},
		{
			name:              "4th",
			g:                 grammars[3],
			expectedNullables: []NonTerminal{},
		},
		{
			name:              "5th",
			g:                 grammars[4],
			expectedNullables: []NonTerminal{},
		},
		{
			name:              "6th",
			g:                 grammars[5],
			expectedNullables: []NonTerminal{},
		},
		{
			name:              "7th",
			g:                 grammars[6],
			expectedNullables: []NonTerminal{},
		},
		{
			name:              "8th",
			g:                 grammars[7],
			expectedNullables: []NonTerminal{"decls"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nullables := tc.g.nullableNonTerminals()

			for nullable := range nullables.All() {
				assert.Contains(t, tc.expectedNullables, nullable)
			}

			for _, expectedNullable := range tc.expectedNullables {
				assert.True(t, nullables.Contains(expectedNullable))
			}
		})
	}
}

func TestGrammar_EliminateEmptyProductions(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{
		{
			name: "1st",
			g:    grammars[0],
			expectedGrammar: New(
				[]Terminal{"0", "1"},
				[]NonTerminal{"S′", "S", "X", "Y"},
				[]Production{
					{"S′", String[Symbol]{NonTerminal("S")}}, // S′ → S
					{"S′", ε}, // S′ → ε
					{"S", String[Symbol]{NonTerminal("X"), NonTerminal("Y"), NonTerminal("X")}}, // S → XYX
					{"S", String[Symbol]{NonTerminal("X"), NonTerminal("X")}},                   // S → XX
					{"S", String[Symbol]{NonTerminal("X"), NonTerminal("Y")}},                   // S → XY
					{"S", String[Symbol]{NonTerminal("Y"), NonTerminal("X")}},                   // S → YX
					{"S", String[Symbol]{NonTerminal("X")}},                                     // S → X
					{"S", String[Symbol]{NonTerminal("Y")}},                                     // S → Y
					{"X", String[Symbol]{Terminal("0"), NonTerminal("X")}},                      // X → 0X
					{"X", String[Symbol]{Terminal("0")}},                                        // X → 0
					{"Y", String[Symbol]{Terminal("1"), NonTerminal("Y")}},                      // Y → 1Y
					{"Y", String[Symbol]{Terminal("1")}},                                        // Y → 1
				},
				"S′",
			),
		},
		{
			name: "2nd",
			g:    grammars[1],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S′", "S"},
				[]Production{
					{"S′", String[Symbol]{NonTerminal("S")}}, // S′ → S
					{"S′", ε}, // S′ → ε
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S → aSbS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S → bSaS
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}},                   // S → aSb
					{"S", String[Symbol]{Terminal("a"), Terminal("b"), NonTerminal("S")}},                   // S → abS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a")}},                   // S → bSa
					{"S", String[Symbol]{Terminal("b"), Terminal("a"), NonTerminal("S")}},                   // S → baS
					{"S", String[Symbol]{Terminal("a"), Terminal("b")}},                                     // S → ab
					{"S", String[Symbol]{Terminal("b"), Terminal("a")}},                                     // S → ba
				},
				"S′",
			),
		},
		{
			name: "3rd",
			g:    grammars[2],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
					{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
					{"S", String[Symbol]{Terminal("a"), Terminal("a")}},                   // S → aa
					{"S", String[Symbol]{Terminal("a")}},                                  // S → a
					{"S", String[Symbol]{Terminal("b")}},                                  // S → b
					{"A", String[Symbol]{Terminal("b")}},                                  // A → b
					{"B", String[Symbol]{NonTerminal("A")}},                               // B → A
					{"B", String[Symbol]{Terminal("b")}},                                  // B → b
				},
				"S",
			),
		},
		{
			name:            "4th",
			g:               grammars[3],
			expectedGrammar: grammars[3],
		},
		{
			name:            "5th",
			g:               grammars[4],
			expectedGrammar: grammars[4],
		},
		{
			name:            "6th",
			g:               grammars[5],
			expectedGrammar: grammars[5],
		},
		{
			name:            "7th",
			g:               grammars[6],
			expectedGrammar: grammars[6],
		},
		{
			name: "8th",
			g:    grammars[7],
			expectedGrammar: New(
				[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
				[]NonTerminal{"grammar", "name", "decls", "decl", "token", "rule", "lhs", "rhs", "nonterm", "term"},
				[]Production{
					{"grammar", String[Symbol]{NonTerminal("name")}},                                // grammar → name
					{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}},          // grammar → name decls
					{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},                // name → GRAMMAR IDENT
					{"decls", String[Symbol]{NonTerminal("decls"), NonTerminal("decl")}},            // decls → decls decl
					{"decls", String[Symbol]{NonTerminal("decl")}},                                  // decls → decl
					{"decl", String[Symbol]{NonTerminal("token")}},                                  // decl → token
					{"decl", String[Symbol]{NonTerminal("rule")}},                                   // decl → rule
					{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}}, // token → TOKEN "=" STRING
					{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},  // token → TOKEN "=" REGEX
					{"rule", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // rule → lhs "=" rhs
					{"lhs", String[Symbol]{NonTerminal("nonterm")}},                                 // lhs → nonterm
					{"rhs", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs")}},                 // rhs → rhs rhs
					{"rhs", String[Symbol]{NonTerminal("rhs"), Terminal("|"), NonTerminal("rhs")}},  // rhs → rhs "|" rhs
					{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")")}},       // rhs → "(" rhs ")"
					{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]")}},       // rhs → "[" rhs "]"
					{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}")}},       // rhs → "{" rhs "}"
					{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}")}},     // rhs → "{{" rhs "}}"
					{"rhs", String[Symbol]{NonTerminal("nonterm")}},                                 // rhs → nonterm
					{"rhs", String[Symbol]{NonTerminal("term")}},                                    // rhs → term
					{"nonterm", String[Symbol]{Terminal("IDENT")}},                                  // nonterm → IDENT
					{"term", String[Symbol]{Terminal("TOKEN")}},                                     // term → TOKEN
					{"term", String[Symbol]{Terminal("STRING")}},                                    // term → STRING
				},
				"grammar",
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateEmptyProductions()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_EliminateSingleProductions(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{
		{
			name:            "1st",
			g:               grammars[0],
			expectedGrammar: grammars[0],
		},
		{
			name:            "2nd",
			g:               grammars[1],
			expectedGrammar: grammars[1],
		},
		{
			name: "3rd",
			g:    grammars[2],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
					{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
					{"S", String[Symbol]{Terminal("a")}},                                  // S → a
					{"A", String[Symbol]{Terminal("b")}},                                  // A → b
					{"A", ε},                                                              // A → ε
					{"B", String[Symbol]{Terminal("b")}},                                  // B → b
					{"B", ε},                                                              // B → ε
				},
				"S",
			),
		},
		{
			name: "4th",
			g:    grammars[3],
			expectedGrammar: New(
				[]Terminal{"b", "c", "d", "s"},
				[]NonTerminal{"S", "A", "B", "C", "D"},
				[]Production{
					{"S", String[Symbol]{Terminal("b")}}, // S → b
					{"S", String[Symbol]{Terminal("d")}}, // S → d
					{"S", String[Symbol]{Terminal("s")}}, // S → s
					{"A", String[Symbol]{Terminal("b")}}, // A → b
					{"A", String[Symbol]{Terminal("d")}}, // A → d
					{"B", String[Symbol]{Terminal("b")}}, // B → b
					{"B", String[Symbol]{Terminal("d")}}, // B → d
					{"C", String[Symbol]{Terminal("d")}}, // C → d
					{"D", String[Symbol]{Terminal("d")}}, // D → d
				},
				"S",
			),
		},
		{
			name:            "5th",
			g:               grammars[4],
			expectedGrammar: grammars[4],
		},
		{
			name: "6th",
			g:    grammars[5],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"E", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // S → E + E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // S → E - E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // S → E * E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // S → E / E
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // S → ( E )
					{"S", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // S → - E
					{"S", String[Symbol]{Terminal("id")}},                                    // S → id
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // E → E + E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // E → E - E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // E → E * E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // E → E / E
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
					{"E", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // E → - E
					{"E", String[Symbol]{Terminal("id")}},                                    // E → id
				},
				"S",
			),
		},
		{
			name: "7th",
			g:    grammars[6],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"E", "T", "F", "S"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // S → E + T
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // S → E - T
					{"S", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // S → T * F
					{"S", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // S → T / F
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // S → ( E )
					{"S", String[Symbol]{Terminal("id")}},                                    // S → id
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"E", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // E → T * F
					{"E", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // E → T / F
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
					{"E", String[Symbol]{Terminal("id")}},                                    // E → id
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
					{"T", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // T → ( E )
					{"T", String[Symbol]{Terminal("id")}},                                    // T → id
					{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
					{"F", String[Symbol]{Terminal("id")}},                                    // F → id
				},
				"S",
			),
		},
		{
			name: "8th",
			g:    grammars[7],
			expectedGrammar: New(
				[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
				[]NonTerminal{"grammar", "name", "decls", "decl", "token", "rule", "lhs", "rhs", "nonterm", "term"},
				[]Production{
					{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}}, // grammar → name decls
					{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},       // name → GRAMMAR IDENT
					{"decls", String[Symbol]{NonTerminal("decls"), NonTerminal("decl")}},   // decls → decls decl
					{"decls", ε}, // decls → ε
					{"decl", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // decl → lhs "=" rhs
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}},  // decl → TOKEN "=" STRING
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},   // decl → TOKEN "=" REGEX
					{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}}, // token → TOKEN "=" STRING
					{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},  // token → TOKEN "=" REGEX
					{"rule", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // rule → lhs "=" rhs
					{"lhs", String[Symbol]{Terminal("IDENT")}},                                      // lhs → IDENT
					{"rhs", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs")}},                 // rhs → rhs rhs
					{"rhs", String[Symbol]{NonTerminal("rhs"), Terminal("|"), NonTerminal("rhs")}},  // rhs → rhs "|" rhs
					{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")")}},       // rhs → "(" rhs ")"
					{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]")}},       // rhs → "[" rhs "]"
					{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}")}},       // rhs → "{" rhs "}"
					{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}")}},     // rhs → "{{" rhs "}}"
					{"rhs", String[Symbol]{Terminal("IDENT")}},                                      // rhs → IDENT
					{"rhs", String[Symbol]{Terminal("TOKEN")}},                                      // rhs → TOKEN
					{"rhs", String[Symbol]{Terminal("STRING")}},                                     // rhs → STRING
					{"nonterm", String[Symbol]{Terminal("IDENT")}},                                  // nonterm → IDENT
					{"term", String[Symbol]{Terminal("TOKEN")}},                                     // term → TOKEN
					{"term", String[Symbol]{Terminal("STRING")}},                                    // term → STRING
				},
				"grammar",
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateSingleProductions()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_EliminateUnreachableProductions(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{
		{
			name:            "1st",
			g:               grammars[0],
			expectedGrammar: grammars[0],
		},
		{
			name:            "2nd",
			g:               grammars[1],
			expectedGrammar: grammars[1],
		},
		{
			name:            "3rd",
			g:               grammars[2],
			expectedGrammar: grammars[2],
		},
		{
			name:            "4th",
			g:               grammars[3],
			expectedGrammar: grammars[3],
		},
		{
			name: "5th",
			g:    grammars[4],
			expectedGrammar: New(
				[]Terminal{"a", "b", "c", "d"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A"), NonTerminal("B")}}, // S → AB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("A")}},    // A → aA
					{"A", String[Symbol]{Terminal("a")}},                      // A → a
					{"B", String[Symbol]{Terminal("b"), NonTerminal("B")}},    // B → bB
					{"B", String[Symbol]{Terminal("b")}},                      // B → b
				},
				"S",
			),
		},
		{
			name:            "6th",
			g:               grammars[5],
			expectedGrammar: grammars[5],
		},
		{
			name:            "7th",
			g:               grammars[6],
			expectedGrammar: grammars[6],
		},
		{
			name:            "8th",
			g:               grammars[7],
			expectedGrammar: grammars[7],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateUnreachableProductions()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_EliminateCycles(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{
		{
			name: "1st",
			g:    grammars[0],
			expectedGrammar: New(
				[]Terminal{"0", "1"},
				[]NonTerminal{"S′", "X", "Y"},
				[]Production{
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("Y"), NonTerminal("X")}}, // S′ → XYX
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("X")}},                   // S′ → XX
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("Y")}},                   // S′ → XY
					{"S′", String[Symbol]{NonTerminal("Y"), NonTerminal("X")}},                   // S′ → YX
					{"S′", String[Symbol]{Terminal("0"), NonTerminal("X")}},                      // S′ → 0X
					{"S′", String[Symbol]{Terminal("1"), NonTerminal("Y")}},                      // S′ → 1Y
					{"S′", String[Symbol]{Terminal("0")}},                                        // S′ → 0
					{"S′", String[Symbol]{Terminal("1")}},                                        // S′ → 1
					{"S′", ε},                                                                    // S′ → ε
					{"X", String[Symbol]{Terminal("0"), NonTerminal("X")}},                       // X → 0X
					{"X", String[Symbol]{Terminal("0")}},                                         // X → 0
					{"Y", String[Symbol]{Terminal("1"), NonTerminal("Y")}},                       // Y → 1Y
					{"Y", String[Symbol]{Terminal("1")}},                                         // Y → 1
				},
				"S′",
			),
		},
		{
			name: "2nd",
			g:    grammars[1],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S′", "S"},
				[]Production{
					{"S′", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S′ → aSbS
					{"S′", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S′ → bSaS
					{"S′", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}},                   // S′ → aSb
					{"S′", String[Symbol]{Terminal("a"), Terminal("b"), NonTerminal("S")}},                   // S′ → abS
					{"S′", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a")}},                   // S′ → bSa
					{"S′", String[Symbol]{Terminal("b"), Terminal("a"), NonTerminal("S")}},                   // S′ → baS
					{"S′", String[Symbol]{Terminal("a"), Terminal("b")}},                                     // S′ → ab
					{"S′", String[Symbol]{Terminal("b"), Terminal("a")}},                                     // S′ → ba
					{"S′", ε}, // S′ → ε
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S → aSbS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S → bSaS
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}},                   // S → aSb
					{"S", String[Symbol]{Terminal("a"), Terminal("b"), NonTerminal("S")}},                   // S → abS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a")}},                   // S → bSa
					{"S", String[Symbol]{Terminal("b"), Terminal("a"), NonTerminal("S")}},                   // S → baS
					{"S", String[Symbol]{Terminal("a"), Terminal("b")}},                                     // S → ab
					{"S", String[Symbol]{Terminal("b"), Terminal("a")}},                                     // S → ba
				},
				"S′",
			),
		},
		{
			name: "3rd",
			g:    grammars[2],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
					{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
					{"S", String[Symbol]{Terminal("a"), Terminal("a")}},                   // S → aa
					{"S", String[Symbol]{Terminal("a")}},                                  // S → a
					{"S", String[Symbol]{Terminal("b")}},                                  // S → b
					{"A", String[Symbol]{Terminal("b")}},                                  // A → b
					{"B", String[Symbol]{Terminal("b")}},                                  // B → b
				},
				"S",
			),
		},
		{
			name: "4th",
			g:    grammars[3],
			expectedGrammar: New(
				[]Terminal{"b", "c", "d", "s"},
				[]NonTerminal{"S"},
				[]Production{
					{"S", String[Symbol]{Terminal("b")}}, // S → b
					{"S", String[Symbol]{Terminal("d")}}, // S → d
					{"S", String[Symbol]{Terminal("s")}}, // S → s
				},
				"S",
			),
		},
		{
			name: "5th",
			g:    grammars[4],
			expectedGrammar: New(
				[]Terminal{"a", "b", "c", "d"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A"), NonTerminal("B")}}, // S → AB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("A")}},    // A → aA
					{"A", String[Symbol]{Terminal("a")}},                      // A → a
					{"B", String[Symbol]{Terminal("b"), NonTerminal("B")}},    // B → bB
					{"B", String[Symbol]{Terminal("b")}},                      // B → b
				},
				"S",
			),
		},
		{
			name: "6th",
			g:    grammars[5],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"S", "E"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // S → E + E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // S → E - E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // S → E * E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // S → E / E
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // S → ( E )
					{"S", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // S → - E
					{"S", String[Symbol]{Terminal("id")}},                                    // S → id
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // E → E + E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // E → E - E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // E → E * E
					{"E", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // E → E / E
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
					{"E", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // E → - E
					{"E", String[Symbol]{Terminal("id")}},                                    // E → id
				},
				"S",
			),
		},
		{
			name: "7th",
			g:    grammars[6],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"S", "E", "T", "F"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // S → E + T
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // S → E - T
					{"S", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // S → T * F
					{"S", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // S → T / F
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // S → ( E )
					{"S", String[Symbol]{Terminal("id")}},                                    // S → id
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"E", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // E → T * F
					{"E", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // E → T / F
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
					{"E", String[Symbol]{Terminal("id")}},                                    // E → id
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
					{"T", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // T → ( E )
					{"T", String[Symbol]{Terminal("id")}},                                    // T → id
					{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
					{"F", String[Symbol]{Terminal("id")}},                                    // F → id
				},
				"S",
			),
		},
		{
			name: "8th",
			g:    grammars[7],
			expectedGrammar: New(
				[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
				[]NonTerminal{"grammar", "name", "decls", "decl", "lhs", "rhs"},
				[]Production{
					{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}},           // grammar → name decls
					{"grammar", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},              // grammar → GRAMMAR IDENT
					{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},                 // name → GRAMMAR IDENT
					{"decls", String[Symbol]{NonTerminal("decls"), NonTerminal("decl")}},             // decls → decls decl
					{"decls", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // decls → lhs "=" rhs
					{"decls", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}},  // decls → TOKEN "=" STRING
					{"decls", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},   // decls → TOKEN "=" REGEX
					{"decl", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}},  // decl → lhs "=" rhs
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}},   // decl → TOKEN "=" STRING
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},    // decl → TOKEN "=" REGEX
					{"lhs", String[Symbol]{Terminal("IDENT")}},                                       // lhs → IDENT
					{"rhs", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs")}},                  // rhs → rhs rhs
					{"rhs", String[Symbol]{NonTerminal("rhs"), Terminal("|"), NonTerminal("rhs")}},   // rhs → rhs "|" rhs
					{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")")}},        // rhs → "(" rhs ")"
					{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]")}},        // rhs → "[" rhs "]"
					{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}")}},        // rhs → "{" rhs "}"
					{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}")}},      // rhs → "{{" rhs "}}"
					{"rhs", String[Symbol]{Terminal("IDENT")}},                                       // rhs → IDENT
					{"rhs", String[Symbol]{Terminal("TOKEN")}},                                       // rhs → TOKEN
					{"rhs", String[Symbol]{Terminal("STRING")}},                                      // rhs → STRING
				},
				"grammar",
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateCycles()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_ChomskyNormalForm(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.ChomskyNormalForm()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_EliminateLeftRecursion(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{
		{
			name: "1st",
			g:    grammars[0],
			expectedGrammar: New(
				[]Terminal{"0", "1"},
				[]NonTerminal{"S′", "X", "Y"},
				[]Production{
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("Y"), NonTerminal("X")}}, // S′ → XYX
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("X")}},                   // S′ → XX
					{"S′", String[Symbol]{NonTerminal("X"), NonTerminal("Y")}},                   // S′ → XY
					{"S′", String[Symbol]{NonTerminal("Y"), NonTerminal("X")}},                   // S′ → YX
					{"S′", String[Symbol]{Terminal("0"), NonTerminal("X")}},                      // S′ → 0X
					{"S′", String[Symbol]{Terminal("1"), NonTerminal("Y")}},                      // S′ → 1Y
					{"S′", String[Symbol]{Terminal("0")}},                                        // S′ → 0
					{"S′", String[Symbol]{Terminal("1")}},                                        // S′ → 1
					{"S′", ε},                                                                    // S′ → ε
					{"X", String[Symbol]{Terminal("0"), NonTerminal("X")}},                       // X → 0X
					{"X", String[Symbol]{Terminal("0")}},                                         // X → 0
					{"Y", String[Symbol]{Terminal("1"), NonTerminal("Y")}},                       // Y → 1Y
					{"Y", String[Symbol]{Terminal("1")}},                                         // Y → 1
				},
				"S′",
			),
		},
		{
			name: "2nd",
			g:    grammars[1],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S′", "S"},
				[]Production{
					{"S′", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S′ → aSbS
					{"S′", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S′ → bSaS
					{"S′", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}},                   // S′ → aSb
					{"S′", String[Symbol]{Terminal("a"), Terminal("b"), NonTerminal("S")}},                   // S′ → abS
					{"S′", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a")}},                   // S′ → bSa
					{"S′", String[Symbol]{Terminal("b"), Terminal("a"), NonTerminal("S")}},                   // S′ → baS
					{"S′", String[Symbol]{Terminal("a"), Terminal("b")}},                                     // S′ → ab
					{"S′", String[Symbol]{Terminal("b"), Terminal("a")}},                                     // S′ → ba
					{"S′", ε}, // S′ → ε
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S → aSbS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S → bSaS
					{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b")}},                   // S → aSb
					{"S", String[Symbol]{Terminal("a"), Terminal("b"), NonTerminal("S")}},                   // S → abS
					{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a")}},                   // S → bSa
					{"S", String[Symbol]{Terminal("b"), Terminal("a"), NonTerminal("S")}},                   // S → baS
					{"S", String[Symbol]{Terminal("a"), Terminal("b")}},                                     // S → ab
					{"S", String[Symbol]{Terminal("b"), Terminal("a")}},                                     // S → ba
				},
				"S′",
			),
		},
		{
			name: "3rd",
			g:    grammars[2],
			expectedGrammar: New(
				[]Terminal{"a", "b"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
					{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
					{"S", String[Symbol]{Terminal("a"), Terminal("a")}},                   // S → aa
					{"S", String[Symbol]{Terminal("a")}},                                  // S → a
					{"S", String[Symbol]{Terminal("b")}},                                  // S → b
					{"A", String[Symbol]{Terminal("b")}},                                  // A → b
					{"B", String[Symbol]{Terminal("b")}},                                  // B → b
				},
				"S",
			),
		},
		{
			name: "4th",
			g:    grammars[3],
			expectedGrammar: New(
				[]Terminal{"b", "c", "d", "s"},
				[]NonTerminal{"S"},
				[]Production{
					{"S", String[Symbol]{Terminal("b")}}, // S → b
					{"S", String[Symbol]{Terminal("d")}}, // S → d
					{"S", String[Symbol]{Terminal("s")}}, // S → s
				},
				"S",
			),
		},
		{
			name: "5th",
			g:    grammars[4],
			expectedGrammar: New(
				[]Terminal{"a", "b", "c", "d"},
				[]NonTerminal{"S", "A", "B"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("A"), NonTerminal("B")}}, // S → AB
					{"A", String[Symbol]{Terminal("a"), NonTerminal("A")}},    // A → aA
					{"A", String[Symbol]{Terminal("a")}},                      // A → a
					{"B", String[Symbol]{Terminal("b"), NonTerminal("B")}},    // B → bB
					{"B", String[Symbol]{Terminal("b")}},                      // B → b
				},
				"S",
			),
		},
		{
			name: "6th",
			g:    grammars[5],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"S", "E", "E′"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}},                 // S → E + E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}},                 // S → E - E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}},                 // S → E * E
					{"S", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}},                 // S → E / E
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},                    // S → ( E )
					{"S", String[Symbol]{Terminal("-"), NonTerminal("E")}},                                   // S → - E
					{"S", String[Symbol]{Terminal("id")}},                                                    // S → id
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")"), NonTerminal("E′")}}, // E → ( E ) E′
					{"E", String[Symbol]{Terminal("-"), NonTerminal("E"), NonTerminal("E′")}},                // E → - E E′
					{"E", String[Symbol]{Terminal("id"), NonTerminal("E′")}},                                 // E → id E′
					{"E′", String[Symbol]{Terminal("+"), NonTerminal("E"), NonTerminal("E′")}},               // E′ → + E E′
					{"E′", String[Symbol]{Terminal("-"), NonTerminal("E"), NonTerminal("E′")}},               // E′ → - E E′
					{"E′", String[Symbol]{Terminal("*"), NonTerminal("E"), NonTerminal("E′")}},               // E′ → * E E′
					{"E′", String[Symbol]{Terminal("/"), NonTerminal("E"), NonTerminal("E′")}},               // E′ → / E E′
					{"E′", ε}, // E′ → ε
				},
				"S",
			),
		},
		{
			name: "7th",
			g:    grammars[6],
			expectedGrammar: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"S", "E", "E′", "T", "T′", "F"},
				[]Production{
					{"S", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}},                    // S → E + T
					{"S", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}},                    // S → E - T
					{"S", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}},                    // S → T * F
					{"S", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}},                    // S → T / F
					{"S", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},                       // S → ( E )
					{"S", String[Symbol]{Terminal("id")}},                                                       // S → id
					{"E", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F"), NonTerminal("E′")}}, // E → T * F E′
					{"E", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F"), NonTerminal("E′")}}, // E → T / F E′
					{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")"), NonTerminal("E′")}},    // E → ( E ) E′
					{"E", String[Symbol]{Terminal("id"), NonTerminal("E′")}},                                    // E → id E′
					{"E′", String[Symbol]{Terminal("+"), NonTerminal("T"), NonTerminal("E′")}},                  // E′ → + T E′
					{"E′", String[Symbol]{Terminal("-"), NonTerminal("T"), NonTerminal("E′")}},                  // E′ → - T E′
					{"E′", ε}, // E′ → ε
					{"T", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")"), NonTerminal("T′")}}, // T → ( E ) T′
					{"T", String[Symbol]{Terminal("id"), NonTerminal("T′")}},                                 // T → id T′
					{"T′", String[Symbol]{Terminal("*"), NonTerminal("F"), NonTerminal("T′")}},               // T′ → * F T′
					{"T′", String[Symbol]{Terminal("/"), NonTerminal("F"), NonTerminal("T′")}},               // T′ → / F T′
					{"T′", ε}, // T′ → ε
					{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}}, // F → ( E )
					{"F", String[Symbol]{Terminal("id")}},                                 // F → id
				},
				"S",
			),
		},
		{
			name: "8th",
			g:    grammars[7],
			expectedGrammar: New(
				[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
				[]NonTerminal{"grammar", "name", "decls", "decls′", "decl", "lhs", "rhs", "rhs′"},
				[]Production{
					{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}},                                  // grammar → name decls
					{"grammar", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},                                     // grammar → GRAMMAR IDENT
					{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},                                        // name → GRAMMAR IDENT
					{"decls", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs"), NonTerminal("decls′")}}, // decls → lhs "=" rhs decls′
					{"decls", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX"), NonTerminal("decls′")}},   // decls → TOKEN "=" REGEX decls′
					{"decls", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING"), NonTerminal("decls′")}},  // decls → TOKEN "=" STRING decls′
					{"decls′", String[Symbol]{NonTerminal("decl"), NonTerminal("decls′")}},                                  // decls′ → decl decls′
					{"decls′", ε}, // decls′ → ε
					{"decl", String[Symbol]{Terminal("IDENT"), Terminal("="), NonTerminal("rhs")}},                   // decl → IDENT "=" rhs
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},                    // decl → TOKEN "=" REGEX
					{"decl", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}},                   // decl → TOKEN "=" STRING
					{"lhs", String[Symbol]{Terminal("IDENT")}},                                                       // lhs → IDENT
					{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")"), NonTerminal("rhs′")}},   // rhs → "(" rhs ")" rhs′
					{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]"), NonTerminal("rhs′")}},   // rhs → "[" rhs "]" rhs′
					{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}"), NonTerminal("rhs′")}},   // rhs → "{" rhs "}" rhs′
					{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}"), NonTerminal("rhs′")}}, // rhs → "{{" rhs "}}" rhs′
					{"rhs", String[Symbol]{Terminal("IDENT"), NonTerminal("rhs′")}},                                  // rhs → IDENT rhs′
					{"rhs", String[Symbol]{Terminal("TOKEN"), NonTerminal("rhs′")}},                                  // rhs → TOKEN rhs′
					{"rhs", String[Symbol]{Terminal("STRING"), NonTerminal("rhs′")}},                                 // rhs → STRING rhs′
					{"rhs′", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs′")}},                                // rhs′ → rhs rhs′
					{"rhs′", String[Symbol]{Terminal("|"), NonTerminal("rhs"), NonTerminal("rhs′")}},                 // rhs′ → "|" rhs rhs′
					{"rhs′", ε}, // rhs′ → ε
				},
				"grammar",
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateLeftRecursion()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_LeftFactor(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar Grammar
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.LeftFactor()
			assert.True(t, g.Equals(tc.expectedGrammar))
		})
	}
}

func TestGrammar_String(t *testing.T) {
	tests := []struct {
		name           string
		g              Grammar
		expectedString string
	}{
		{
			name:           "1st",
			g:              grammars[0],
			expectedString: "Terminal Symbols: \"0\" \"1\"\nNon-Terminal Symbols: S X Y\nStart Symbol: S\nProduction Rules:\n  S → X Y X\n  X → \"0\" X | ε\n  Y → \"1\" Y | ε\n",
		},
		{
			name:           "2nd",
			g:              grammars[1],
			expectedString: "Terminal Symbols: \"a\" \"b\"\nNon-Terminal Symbols: S\nStart Symbol: S\nProduction Rules:\n  S → \"a\" S \"b\" S | \"b\" S \"a\" S | ε\n",
		},
		{
			name:           "3rd",
			g:              grammars[2],
			expectedString: "Terminal Symbols: \"a\" \"b\"\nNon-Terminal Symbols: S B A\nStart Symbol: S\nProduction Rules:\n  S → \"a\" B \"a\" | A \"b\" | \"a\"\n  B → A | \"b\"\n  A → \"b\" | ε\n",
		},
		{
			name:           "4th",
			g:              grammars[3],
			expectedString: "Terminal Symbols: \"b\" \"c\" \"d\" \"s\"\nNon-Terminal Symbols: S A B C D\nStart Symbol: S\nProduction Rules:\n  S → A | \"s\"\n  A → B\n  B → C | \"b\"\n  C → D\n  D → \"d\"\n",
		},
		{
			name:           "5th",
			g:              grammars[4],
			expectedString: "Terminal Symbols: \"a\" \"b\" \"c\" \"d\"\nNon-Terminal Symbols: S A B C D\nStart Symbol: S\nProduction Rules:\n  S → A B\n  A → \"a\" A | \"a\"\n  B → \"b\" B | \"b\"\n  C → \"c\" C | \"c\"\n  D → \"d\"\n",
		},
		{
			name:           "6th",
			g:              grammars[5],
			expectedString: "Terminal Symbols: \"(\" \")\" \"*\" \"+\" \"-\" \"/\" \"id\"\nNon-Terminal Symbols: S E\nStart Symbol: S\nProduction Rules:\n  S → E\n  E → E \"*\" E | E \"+\" E | E \"-\" E | E \"/\" E | \"(\" E \")\" | \"-\" E | \"id\"\n",
		},
		{
			name:           "7th",
			g:              grammars[6],
			expectedString: "Terminal Symbols: \"(\" \")\" \"*\" \"+\" \"-\" \"/\" \"id\"\nNon-Terminal Symbols: S E T F\nStart Symbol: S\nProduction Rules:\n  S → E\n  E → E \"+\" T | E \"-\" T | T\n  T → T \"*\" F | T \"/\" F | F\n  F → \"(\" E \")\" | \"id\"\n",
		},
		{
			name:           "8th",
			g:              grammars[7],
			expectedString: "Terminal Symbols: \"(\" \")\" \"=\" \"GRAMMAR\" \"IDENT\" \"REGEX\" \"STRING\" \"TOKEN\" \"[\" \"]\" \"{\" \"{{\" \"|\" \"}\" \"}}\"\nNon-Terminal Symbols: grammar name decls decl rule token lhs rhs nonterm term\nStart Symbol: grammar\nProduction Rules:\n  grammar → name decls\n  name → \"GRAMMAR\" \"IDENT\"\n  decls → decls decl | ε\n  decl → rule | token\n  rule → lhs \"=\" rhs\n  token → \"TOKEN\" \"=\" \"REGEX\" | \"TOKEN\" \"=\" \"STRING\"\n  lhs → nonterm\n  rhs → rhs \"|\" rhs | rhs rhs | \"(\" rhs \")\" | \"[\" rhs \"]\" | \"{\" rhs \"}\" | \"{{\" rhs \"}}\" | nonterm | term\n  nonterm → \"IDENT\"\n  term → \"STRING\" | \"TOKEN\"\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.g.String())
		})
	}
}

func TestGrammar_addNewNonTerminal(t *testing.T) {
	tests := []struct {
		name                string
		g                   Grammar
		prefix              NonTerminal
		suffixes            []string
		expectedOK          bool
		expectedNonTerminal NonTerminal
	}{
		{
			name:                "OK",
			g:                   grammars[0],
			prefix:              NonTerminal("S"),
			suffixes:            []string{"_new"},
			expectedOK:          true,
			expectedNonTerminal: NonTerminal("S_new"),
		},
		{
			name:                "NotOK",
			g:                   grammars[0],
			prefix:              NonTerminal("S"),
			suffixes:            []string{""},
			expectedOK:          false,
			expectedNonTerminal: NonTerminal(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nonTerm, ok := tc.g.addNewNonTerminal(tc.prefix, tc.suffixes...)
			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedNonTerminal, nonTerm)
		})
	}
}

func TestGrammar_orderTerminals(t *testing.T) {
	tests := []struct {
		name              string
		g                 Grammar
		expectedTerminals String[Terminal]
	}{
		{
			name:              "OK",
			g:                 grammars[4],
			expectedTerminals: String[Terminal]{"a", "b", "c", "d"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			terms := tc.g.orderTerminals()
			assert.Equal(t, tc.expectedTerminals, terms)
		})
	}
}

func TestGrammar_orderNonTerminals(t *testing.T) {
	tests := []struct {
		name                 string
		g                    Grammar
		expectedVisited      []NonTerminal
		expectedUnvisited    []NonTerminal
		expectedNonTerminals String[NonTerminal]
	}{
		{
			name:                 "OK",
			g:                    grammars[4],
			expectedVisited:      []NonTerminal{"S", "A", "B"},
			expectedUnvisited:    []NonTerminal{"C", "D"},
			expectedNonTerminals: String[NonTerminal]{"S", "A", "B", "C", "D"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			visited, unvisited, nonTerms := tc.g.orderNonTerminals()
			assert.Equal(t, tc.expectedVisited, visited)
			assert.Equal(t, tc.expectedUnvisited, unvisited)
			assert.Equal(t, tc.expectedNonTerminals, nonTerms)
		})
	}
}
