package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name                 string
		s                    String[Symbol]
		expectedString       string
		prefix               String[Symbol]
		expectedHasPrefix    bool
		suffix               String[Symbol]
		expectedHasSuffix    bool
		append               []Symbol
		expectedAppend       String[Symbol]
		concat               []String[Symbol]
		expectedConcat       String[Symbol]
		expectedTerminals    String[Terminal]
		expectedNonTerminals String[NonTerminal]
	}{
		{
			name:                 "Empty",
			s:                    ε,
			expectedString:       `ε`,
			prefix:               String[Symbol]{},
			expectedHasPrefix:    true,
			suffix:               String[Symbol]{},
			expectedHasSuffix:    true,
			append:               []Symbol{},
			expectedAppend:       ε,
			concat:               []String[Symbol]{ε},
			expectedConcat:       ε,
			expectedTerminals:    String[Terminal]{},
			expectedNonTerminals: String[NonTerminal]{},
		},
		{
			name:                 "AllTerminals",
			s:                    String[Symbol]{Terminal("a"), Terminal("b"), Terminal("c")},
			expectedString:       `"a" "b" "c"`,
			prefix:               String[Symbol]{Terminal("a"), Terminal("b")},
			expectedHasPrefix:    true,
			suffix:               String[Symbol]{Terminal("a"), Terminal("c")},
			expectedHasSuffix:    false,
			append:               []Symbol{Terminal("d")},
			expectedAppend:       String[Symbol]{Terminal("a"), Terminal("b"), Terminal("c"), Terminal("d")},
			concat:               []String[Symbol]{{Terminal("d"), Terminal("e"), Terminal("f")}},
			expectedConcat:       String[Symbol]{Terminal("a"), Terminal("b"), Terminal("c"), Terminal("d"), Terminal("e"), Terminal("f")},
			expectedTerminals:    String[Terminal]{"a", "b", "c"},
			expectedNonTerminals: String[NonTerminal]{},
		},
		{
			name:                 "AllNonTerminals",
			s:                    String[Symbol]{NonTerminal("A"), NonTerminal("B"), NonTerminal("C")},
			expectedString:       `A B C`,
			prefix:               String[Symbol]{NonTerminal("A"), NonTerminal("C")},
			expectedHasPrefix:    false,
			suffix:               String[Symbol]{NonTerminal("B"), NonTerminal("C")},
			expectedHasSuffix:    true,
			append:               []Symbol{NonTerminal("D")},
			expectedAppend:       String[Symbol]{NonTerminal("A"), NonTerminal("B"), NonTerminal("C"), NonTerminal("D")},
			concat:               []String[Symbol]{{NonTerminal("D"), NonTerminal("E"), NonTerminal("F")}},
			expectedConcat:       String[Symbol]{NonTerminal("A"), NonTerminal("B"), NonTerminal("C"), NonTerminal("D"), NonTerminal("E"), NonTerminal("F")},
			expectedTerminals:    String[Terminal]{},
			expectedNonTerminals: String[NonTerminal]{"A", "B", "C"},
		},
		{
			name:                 "TerminalsAndNonTerminals",
			s:                    String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c")},
			expectedString:       `"a" A "b" B "c"`,
			prefix:               String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c")},
			expectedHasPrefix:    true,
			suffix:               String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c")},
			expectedHasSuffix:    true,
			append:               []Symbol{NonTerminal("C"), Terminal("d"), NonTerminal("D")},
			expectedAppend:       String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c"), NonTerminal("C"), Terminal("d"), NonTerminal("D")},
			concat:               []String[Symbol]{{NonTerminal("C")}, {Terminal("d"), NonTerminal("D")}},
			expectedConcat:       String[Symbol]{Terminal("a"), NonTerminal("A"), Terminal("b"), NonTerminal("B"), Terminal("c"), NonTerminal("C"), Terminal("d"), NonTerminal("D")},
			expectedTerminals:    String[Terminal]{"a", "b", "c"},
			expectedNonTerminals: String[NonTerminal]{"A", "B"},
		},
	}

	notEqual := String[Symbol]{Terminal("🙂"), NonTerminal("🙃")}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.s.String())
			assert.True(t, tc.s.Equals(tc.s))
			assert.False(t, tc.s.Equals(notEqual))
			assert.Equal(t, tc.expectedHasPrefix, tc.s.HasPrefix(tc.prefix))
			assert.Equal(t, tc.expectedHasSuffix, tc.s.HasSuffix(tc.suffix))
			assert.Equal(t, tc.expectedAppend, tc.s.Append(tc.append...))
			assert.Equal(t, tc.expectedConcat, tc.s.Concat(tc.concat...))
			assert.Equal(t, tc.expectedTerminals, tc.s.Terminals())
			assert.Equal(t, tc.expectedNonTerminals, tc.s.NonTerminals())
		})
	}
}

func TestLongestCommonPrefixOf(t *testing.T) {
	tests := []struct {
		name                        string
		ss                          []String[Symbol]
		expectedLongestCommonPrefix String[Symbol]
	}{
		{
			name:                        "Empty",
			ss:                          []String[Symbol]{},
			expectedLongestCommonPrefix: ε,
		},
		{
			name: "NoCommonPrefix",
			ss: []String[Symbol]{
				{NonTerminal("expr"), Terminal("?"), NonTerminal("stmt"), Terminal(":"), NonTerminal("stmt")},
				{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt"), Terminal("else"), NonTerminal("stmt")},
			},
			expectedLongestCommonPrefix: ε,
		},
		{
			name: "CommonPrefix",
			ss: []String[Symbol]{
				{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt"), Terminal("else"), NonTerminal("stmt")},
				{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt")},
			},
			expectedLongestCommonPrefix: String[Symbol]{Terminal("if"), NonTerminal("expr"), Terminal("then"), NonTerminal("stmt")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lcp := LongestCommonPrefixOf(tc.ss...)
			assert.Equal(t, tc.expectedLongestCommonPrefix, lcp)
		})
	}
}
