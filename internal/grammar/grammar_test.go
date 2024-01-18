package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var grammars = []CFG{
	New(
		[]Terminal{"a", "b"},
		[]NonTerminal{"A", "B", "S"},
		[]Production{
			{"S", String[Symbol]{Terminal("a"), NonTerminal("S"), Terminal("b"), NonTerminal("S")}}, // S → aSbS
			{"S", String[Symbol]{Terminal("b"), NonTerminal("S"), Terminal("a"), NonTerminal("S")}}, // S → bSaS
			{"S", ε}, // S → ε
		},
		"S",
	),
	New(
		[]Terminal{"a", "b"},
		[]NonTerminal{"A", "B", "S"},
		[]Production{
			{"S", String[Symbol]{Terminal("a")}},                                  // S → a
			{"S", String[Symbol]{NonTerminal("A"), Terminal("b")}},                // S → Ab
			{"S", String[Symbol]{Terminal("a"), NonTerminal("B"), Terminal("a")}}, // S → aBa
			{"A", String[Symbol]{Terminal("b")}},                                  // A → b
			{"A", ε},                                                              // A → ε
			{"B", String[Symbol]{Terminal("b")}},                                  // B → b
			{"B", String[Symbol]{NonTerminal("A")}},                               // B → A
		},
		"S",
	),
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
		[]Terminal{"=", "|", "(", ")", "[", "]", "{", "}", "{{", "}}", "GRAMMAR", "IDENT", "TOKEN", "STRING", "REGEX"},
		[]NonTerminal{"grammar", "name", "decls", "decl", "token", "rule", "lhs", "rhs", "nonterm", "term"},
		[]Production{
			{"grammar", String[Symbol]{NonTerminal("name"), NonTerminal("decls")}},          // grammar → name decls
			{"name", String[Symbol]{Terminal("GRAMMAR"), Terminal("IDENT")}},                // name → GRAMMAR IDENT
			{"decls", String[Symbol]{}},                                                     // decls → ε
			{"decls", String[Symbol]{NonTerminal("decls"), NonTerminal("decl")}},            // decls → decls decl
			{"decl", String[Symbol]{NonTerminal("token")}},                                  // decl → token
			{"decl", String[Symbol]{NonTerminal("rule")}},                                   // decl → rule
			{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("STRING")}}, // token → TOKEN "=" STRING
			{"token", String[Symbol]{Terminal("TOKEN"), Terminal("="), Terminal("REGEX")}},  // token → TOKEN "=" REGEX
			{"rule", String[Symbol]{NonTerminal("lhs"), Terminal("="), NonTerminal("rhs")}}, // rule → lhs "=" rhs
			{"lhs", String[Symbol]{NonTerminal("nonterm")}},                                 // lhs → nonterm
			{"rhs", String[Symbol]{NonTerminal("nonterm")}},                                 // rhs → nonterm
			{"rhs", String[Symbol]{Terminal("term")}},                                       // rhs → term
			{"rhs", String[Symbol]{Terminal("("), NonTerminal("rhs"), Terminal(")")}},       // rhs → "(" rhs ")"
			{"rhs", String[Symbol]{Terminal("["), NonTerminal("rhs"), Terminal("]")}},       // rhs → "[" rhs "]"
			{"rhs", String[Symbol]{Terminal("{"), NonTerminal("rhs"), Terminal("}")}},       // rhs → "{" rhs "}"
			{"rhs", String[Symbol]{Terminal("{{"), NonTerminal("rhs"), Terminal("}}")}},     // rhs → "{{" rhs "}}"
			{"rhs", String[Symbol]{NonTerminal("rhs"), NonTerminal("rhs")}},                 // rhs → rhs rhs
			{"rhs", String[Symbol]{NonTerminal("rhs"), Terminal("|"), NonTerminal("rhs")}},  // rhs → rhs "|" rhs
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
		expectedIsSingle bool
	}{
		{
			name:             "First",
			p:                Production{"S", ε},
			expectedString:   "S → ε",
			expectedIsSingle: false,
		},
		{
			name:             "Second",
			p:                Production{"A", String[Symbol]{NonTerminal("B")}},
			expectedString:   "A → B",
			expectedIsSingle: true,
		},
		{
			name:             "Third",
			p:                Production{"stmt", String[Symbol]{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt")}},
			expectedString:   "stmt → if expr then stmt",
			expectedIsSingle: false,
		},
	}

	notEqual := Production{"😐", String[Symbol]{Terminal("🙂"), NonTerminal("🙃")}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.p.String())
			assert.True(t, tc.p.Equals(tc.p))
			assert.False(t, tc.p.Equals(notEqual))
			assert.Equal(t, tc.expectedIsSingle, tc.p.IsSingle())
		})
	}
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
				{"S", String[Symbol]{Terminal("a"), Terminal("S"), Terminal("b")}}, //  S → aSb
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
				{"S", String[Symbol]{Terminal("("), Terminal(")")}},                   //  S → ()
				{"S", String[Symbol]{Terminal("["), NonTerminal("S"), Terminal("]")}}, //  S → [S]
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

func TestCFG_String(t *testing.T) {
	tests := []struct {
		name           string
		g              CFG
		expectedString string
	}{
		{
			name:           "First",
			g:              grammars[0],
			expectedString: "Terminal Symbols: a b\nNon-Terminal Symbols: A B S\nStart Symbol: S\nProduction Rules:\n  S → a S b S\n  S → b S a S\n  S → ε\n",
		},
		{
			name:           "Second",
			g:              grammars[1],
			expectedString: "Terminal Symbols: a b\nNon-Terminal Symbols: A B S\nStart Symbol: S\nProduction Rules:\n  S → a\n  S → A b\n  S → a B a\n  A → b\n  A → ε\n  B → b\n  B → A\n",
		},
		{
			name:           "Third",
			g:              grammars[2],
			expectedString: "Terminal Symbols: 0 1\nNon-Terminal Symbols: S X Y\nStart Symbol: S\nProduction Rules:\n  S → X Y X\n  X → 0 X\n  X → ε\n  Y → 1 Y\n  Y → ε\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.g.String())
		})
	}
}

func TestCFG_verify(t *testing.T) {
	tests := []struct {
		name          string
		g             CFG
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
			err := tc.g.verify()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCFG_findNullableNonTerminals(t *testing.T) {
	tests := []struct {
		name              string
		g                 CFG
		expectedNullables []NonTerminal
	}{
		{
			name:              "First",
			g:                 grammars[0],
			expectedNullables: []NonTerminal{"S"},
		},
		{
			name:              "Second",
			g:                 grammars[1],
			expectedNullables: []NonTerminal{"A", "B"},
		},
		{
			name:              "Third",
			g:                 grammars[2],
			expectedNullables: []NonTerminal{"X", "Y", "S"},
		},
		{
			name:              "Fourth",
			g:                 grammars[3],
			expectedNullables: []NonTerminal{"decls"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nullables := tc.g.findNullableNonTerminals()
			assert.Equal(t, tc.expectedNullables, nullables.Members())
		})
	}
}

func TestCFG_EliminateEmptyProductions(t *testing.T) {
	tests := []struct {
		name            string
		g               CFG
		expectedGrammar string
	}{
		{
			name:            "First",
			g:               grammars[0],
			expectedGrammar: "Terminal Symbols: a b\nNon-Terminal Symbols: A B S\nStart Symbol: S\nProduction Rules:\n  S → a b\n  S → a b S\n  S → a S b\n  S → a S b S\n  S → b a\n  S → b a S\n  S → b S a\n  S → b S a S\n",
		},
		{
			name:            "Second",
			g:               grammars[1],
			expectedGrammar: "Terminal Symbols: a b\nNon-Terminal Symbols: A B S\nStart Symbol: S\nProduction Rules:\n  S → a\n  S → b\n  S → A b\n  S → a a\n  S → a B a\n  A → b\n  B → b\n  B → A\n",
		},
		{
			name:            "Third",
			g:               grammars[2],
			expectedGrammar: "Terminal Symbols: 0 1\nNon-Terminal Symbols: S X Y\nStart Symbol: S\nProduction Rules:\n  S → X\n  S → Y\n  S → Y X\n  S → X X\n  S → X Y\n  S → X Y X\n  X → 0\n  X → 0 X\n  Y → 1\n  Y → 1 Y\n",
		},
		{
			name:            "Fourth",
			g:               grammars[3],
			expectedGrammar: "Terminal Symbols: = | ( ) [ ] { } {{ }} GRAMMAR IDENT TOKEN STRING REGEX\nNon-Terminal Symbols: grammar name decls decl token rule lhs rhs nonterm term\nStart Symbol: grammar\nProduction Rules:\n  grammar → name\n  grammar → name decls\n  name → GRAMMAR IDENT\n  decls → decl\n  decls → decls decl\n  decl → token\n  decl → rule\n  token → TOKEN = STRING\n  token → TOKEN = REGEX\n  rule → lhs = rhs\n  lhs → nonterm\n  rhs → nonterm\n  rhs → term\n  rhs → ( rhs )\n  rhs → [ rhs ]\n  rhs → { rhs }\n  rhs → {{ rhs }}\n  rhs → rhs rhs\n  rhs → rhs | rhs\n  nonterm → IDENT\n  term → TOKEN\n  term → STRING\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.g.EliminateEmptyProductions()
			assert.Equal(t, tc.expectedGrammar, g.String())
		})
	}
}
