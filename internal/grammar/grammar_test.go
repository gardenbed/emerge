package grammar

import (
	"slices"
	"testing"

	"github.com/moorara/algo/set"
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
			{"E", String[Symbol]{NonTerminal("T")}},                                  // T → F
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

func TestTerminal(t *testing.T) {
	tests := []struct {
		value string
	}{
		{value: "a"},
		{value: "b"},
		{value: "c"},
		{value: "0"},
		{value: "1"},
		{value: "2"},
		{value: "+"},
		{value: "*"},
		{value: "("},
		{value: ")"},
		{value: "id"},
		{value: "if"},
	}

	notEqual := Terminal("🙂")

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			tr := Terminal(tc.value)
			assert.Equal(t, tc.value, tr.String())
			assert.Equal(t, tc.value, tr.Name())
			assert.True(t, tr.Equals(Terminal(tc.value)))
			assert.False(t, tr.Equals(NonTerminal(tc.value)))
			assert.False(t, tr.Equals(notEqual))
			assert.True(t, tr.IsTerminal())
		})
	}
}

func TestNonTerminal(t *testing.T) {
	tests := []struct {
		value string
	}{
		{value: "A"},
		{value: "B"},
		{value: "C"},
		{value: "S"},
		{value: "expr"},
		{value: "stmt"},
	}

	notEqual := NonTerminal("🙂")

	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			n := NonTerminal(tc.value)
			assert.Equal(t, tc.value, n.String())
			assert.Equal(t, tc.value, n.Name())
			assert.True(t, n.Equals(NonTerminal(tc.value)))
			assert.False(t, n.Equals(Terminal(tc.value)))
			assert.False(t, n.Equals(notEqual))
			assert.False(t, n.IsTerminal())
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name                 string
		s                    String[Symbol]
		expectedString       string
		expectedTerminals    String[Terminal]
		expectedNonTerminals String[NonTerminal]
	}{
		{
			name:                 "Empty",
			s:                    ε,
			expectedString:       "ε",
			expectedTerminals:    String[Terminal]{},
			expectedNonTerminals: String[NonTerminal]{},
		},
		{
			name:                 "AllTerminals",
			s:                    String[Symbol]{Terminal("a"), Terminal("b"), Terminal("c")},
			expectedString:       "a b c",
			expectedTerminals:    String[Terminal]{"a", "b", "c"},
			expectedNonTerminals: String[NonTerminal]{},
		},
		{
			name:                 "AllNonTerminals",
			s:                    String[Symbol]{NonTerminal("A"), NonTerminal("B"), NonTerminal("C")},
			expectedString:       "A B C",
			expectedTerminals:    String[Terminal]{},
			expectedNonTerminals: String[NonTerminal]{"A", "B", "C"},
		},
		{
			name:                 "TerminalsAndNonTerminals",
			s:                    String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c")},
			expectedString:       "a A b B c",
			expectedTerminals:    String[Terminal]{"a", "b", "c"},
			expectedNonTerminals: String[NonTerminal]{"A", "B"},
		},
	}

	notEqual := String[Symbol]{Terminal("🙂"), NonTerminal("🙃")}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.s.String())
			assert.Equal(t, tc.expectedTerminals, tc.s.Terminals())
			assert.Equal(t, tc.expectedNonTerminals, tc.s.NonTerminals())
			assert.True(t, tc.s.Equals(tc.s))
			assert.False(t, tc.s.Equals(notEqual))
		})
	}
}

func TestProduction(t *testing.T) {
	tests := []struct {
		name             string
		p                Production
		expectedString   string
		expectedIsEmpty  bool
		expectedIsSingle bool
	}{
		{
			name:             "First",
			p:                Production{"S", ε},
			expectedString:   "S → ε",
			expectedIsEmpty:  true,
			expectedIsSingle: false,
		},
		{
			name:             "Second",
			p:                Production{"A", String[Symbol]{NonTerminal("B")}},
			expectedString:   "A → B",
			expectedIsEmpty:  false,
			expectedIsSingle: true,
		},
		{
			name:             "Third",
			p:                Production{"stmt", String[Symbol]{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt")}},
			expectedString:   "stmt → if expr then stmt",
			expectedIsEmpty:  false,
			expectedIsSingle: false,
		},
	}

	notEqual := Production{"😐", String[Symbol]{Terminal("🙂"), NonTerminal("🙃")}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.p.String())
			assert.True(t, tc.p.Equals(tc.p))
			assert.False(t, tc.p.Equals(notEqual))
			assert.Equal(t, tc.expectedIsEmpty, tc.p.IsEmpty())
			assert.Equal(t, tc.expectedIsSingle, tc.p.IsSingle())
		})
	}
}

func TestProductions(t *testing.T) {
	prods := NewProductions()
	assert.NotNil(t, prods)
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

func TestGrammar_AddProduction(t *testing.T) {
	tests := []struct {
		name  string
		g     Grammar
		prods []Production
	}{
		{
			name: "OK",
			g: New(
				[]Terminal{"+", "-", "*", "/", "(", ")", "id"},
				[]NonTerminal{"E", "T", "F", "S"},
				[]Production{},
				"S",
			),
			prods: []Production{
				{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
				{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
				{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
				{"E", String[Symbol]{NonTerminal("T")}},                                  // E → T
				{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
				{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
				{"E", String[Symbol]{NonTerminal("T")}},                                  // T → F
				{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
				{"F", String[Symbol]{Terminal("id")}},                                    // F → id
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.g.AddProduction(tc.prods...)

			for _, expectedProduction := range tc.prods {
				list, ok := tc.g.Productions.Get(expectedProduction.Head)
				assert.True(t, ok)
				assert.True(t, list.Contains(expectedProduction))
			}
		})
	}
}

func TestGrammar_AllProductions(t *testing.T) {
	tests := []struct {
		name                string
		g                   Grammar
		expectedProductions []Production
	}{
		{
			name: "OK",
			g:    grammars[6],
			expectedProductions: []Production{
				{"S", String[Symbol]{NonTerminal("E")}},                                  // S → E
				{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
				{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
				{"E", String[Symbol]{NonTerminal("T")}},                                  // E → T
				{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
				{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
				{"E", String[Symbol]{NonTerminal("T")}},                                  // T → F
				{"F", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // F → ( E )
				{"F", String[Symbol]{Terminal("id")}},                                    // F → id
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prods := slices.Collect[Production](tc.g.AllProductions())

			for _, expectedProduction := range tc.expectedProductions {
				assert.Contains(t, prods, expectedProduction)
			}
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
			expectedError: "start symbol \"S\" not in the set of non-terminal symbols\nno production rule for start symbol \"S\"",
		},
		{
			name: "StartSymbolHasNoProduction",
			g: New(
				[]Terminal{},
				[]NonTerminal{"S"},
				[]Production{},
				"S",
			),
			expectedError: "no production rule for start symbol \"S\"\nno production rule for non-terminal symbol \"S\"",
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
			expectedError: "no production rule for non-terminal symbol \"A\"",
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
			expectedError: "production head \"B\" not in the set of non-terminal symbols",
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
			expectedError: "non-terminal symbol \"C\" not in the set of non-terminal symbols",
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
					{"E", String[Symbol]{NonTerminal("T")}},                                  // T → F
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
					{"E", String[Symbol]{NonTerminal("T")}},                                  // T → F
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
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"E", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // E → T * F
					{"E", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // E → T / F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
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
					{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("T")}}, // E → E + T
					{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("T")}}, // E → E - T
					{"E", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // E → T * F
					{"E", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // E → T / F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("*"), NonTerminal("F")}}, // T → T * F
					{"T", String[Symbol]{NonTerminal("T"), Terminal("/"), NonTerminal("F")}}, // T → T / F
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

func TestGrammar_EliminateLeftRecursion(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateLeftRecursion()
			assert.Equal(t, tc.expectedGrammar, g.String())
		})
	}
}

func TestGrammar_LeftFactor(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.LeftFactor()
			assert.Equal(t, tc.expectedGrammar, g.String())
		})
	}
}

func TestGrammar_ChomskyNormalForm(t *testing.T) {
	tests := []struct {
		name            string
		g               Grammar
		expectedGrammar string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.ChomskyNormalForm()
			assert.Equal(t, tc.expectedGrammar, g.String())
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
			expectedString: "Terminal Symbols: 0 1\nNon-Terminal Symbols: S X Y\nStart Symbol: S\nProduction Rules:\n  S → X Y X\n  X → 0 X\n  X → ε\n  Y → 1 Y\n  Y → ε\n",
		},
		{
			name:           "2nd",
			g:              grammars[1],
			expectedString: "Terminal Symbols: a b\nNon-Terminal Symbols: S\nStart Symbol: S\nProduction Rules:\n  S → a S b S\n  S → b S a S\n  S → ε\n",
		},
		{
			name:           "3rd",
			g:              grammars[2],
			expectedString: "Terminal Symbols: a b\nNon-Terminal Symbols: S B A\nStart Symbol: S\nProduction Rules:\n  S → a B a\n  S → A b\n  S → a\n  B → A\n  B → b\n  A → b\n  A → ε\n",
		},
		{
			name:           "4th",
			g:              grammars[3],
			expectedString: "Terminal Symbols: b c d s\nNon-Terminal Symbols: S A B C D\nStart Symbol: S\nProduction Rules:\n  S → A\n  S → s\n  A → B\n  B → C\n  B → b\n  C → D\n  D → d\n",
		},
		{
			name:           "5th",
			g:              grammars[4],
			expectedString: "Terminal Symbols: a b c d\nNon-Terminal Symbols: S A B C D\nStart Symbol: S\nProduction Rules:\n  S → A B\n  A → a A\n  A → a\n  B → b B\n  B → b\n  C → c C\n  C → c\n  D → d\n",
		},
		{
			name:           "6th",
			g:              grammars[5],
			expectedString: "Terminal Symbols: ( ) * + - / id\nNon-Terminal Symbols: S E\nStart Symbol: S\nProduction Rules:\n  S → E\n  E → E * E\n  E → E + E\n  E → E - E\n  E → E / E\n  E → ( E )\n  E → - E\n  E → id\n",
		},
		{
			name:           "7th",
			g:              grammars[6],
			expectedString: "Terminal Symbols: ( ) * + - / id\nNon-Terminal Symbols: S E T F\nStart Symbol: S\nProduction Rules:\n  S → E\n  E → E + T\n  E → E - T\n  E → T\n  T → T * F\n  T → T / F\n  F → ( E )\n  F → id\n",
		},
		{
			name:           "8th",
			g:              grammars[7],
			expectedString: "Terminal Symbols: ( ) = GRAMMAR IDENT REGEX STRING TOKEN [ ] { {{ | } }}\nNon-Terminal Symbols: grammar name decls decl rule token lhs rhs nonterm term\nStart Symbol: grammar\nProduction Rules:\n  grammar → name decls\n  name → GRAMMAR IDENT\n  decls → decls decl\n  decls → ε\n  decl → rule\n  decl → token\n  rule → lhs = rhs\n  token → TOKEN = REGEX\n  token → TOKEN = STRING\n  lhs → nonterm\n  rhs → rhs | rhs\n  rhs → rhs rhs\n  rhs → ( rhs )\n  rhs → [ rhs ]\n  rhs → { rhs }\n  rhs → {{ rhs }}\n  rhs → nonterm\n  rhs → term\n  nonterm → IDENT\n  term → STRING\n  term → TOKEN\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.g.String())
		})
	}
}

func TestGrammar_generateNewNonTerminal(t *testing.T) {
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
			nonTerm, ok := tc.g.generateNewNonTerminal(tc.prefix, tc.suffixes...)
			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedNonTerminal, nonTerm)
		})
	}
}

func TestGrammar_orderTerminals(t *testing.T) {
	tests := []struct {
		name              string
		g                 Grammar
		expectedTerminals []string
	}{
		{
			name:              "OK",
			g:                 grammars[4],
			expectedTerminals: []string{"a", "b", "c", "d"},
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
		expectedNonTerminals []string
	}{
		{
			name:                 "OK",
			g:                    grammars[4],
			expectedVisited:      []NonTerminal{"S", "A", "B"},
			expectedUnvisited:    []NonTerminal{"C", "D"},
			expectedNonTerminals: []string{"S", "A", "B", "C", "D"},
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

func TestGrammar_orderProductions(t *testing.T) {
	s := set.New[Production](eqProduction)
	s.Add(
		Production{"E", String[Symbol]{Terminal("id")}},                                    // E → id
		Production{"E", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // E → - E
		Production{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
		Production{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // E → E + E
		Production{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // E → E - E
		Production{"E", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // E → E * E
		Production{"E", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // E → E / E
	)

	tests := []struct {
		name                string
		g                   Grammar
		set                 set.Set[Production]
		expectedProductions []Production
	}{
		{
			name: "OK",
			set:  s,
			expectedProductions: []Production{
				{"E", String[Symbol]{NonTerminal("E"), Terminal("*"), NonTerminal("E")}}, // E → E * E
				{"E", String[Symbol]{NonTerminal("E"), Terminal("+"), NonTerminal("E")}}, // E → E + E
				{"E", String[Symbol]{NonTerminal("E"), Terminal("-"), NonTerminal("E")}}, // E → E - E
				{"E", String[Symbol]{NonTerminal("E"), Terminal("/"), NonTerminal("E")}}, // E → E / E
				{"E", String[Symbol]{Terminal("("), NonTerminal("E"), Terminal(")")}},    // E → ( E )
				{"E", String[Symbol]{Terminal("-"), NonTerminal("E")}},                   // E → - E
				{"E", String[Symbol]{Terminal("id")}},                                    // E → id
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prods := orderProductions(tc.set)
			assert.Equal(t, tc.expectedProductions, prods)
		})
	}
}
