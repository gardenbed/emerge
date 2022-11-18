package lexer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_String(t *testing.T) {
	tests := []struct {
		name           string
		t              Tag
		expectedString string
	}{
		{
			name:           "ERR",
			t:              ERR,
			expectedString: "ERR",
		},
		{
			name:           "WS",
			t:              WS,
			expectedString: "WS",
		},
		{
			name:           "DEF",
			t:              DEF,
			expectedString: "DEF",
		},
		{
			name:           "ALT",
			t:              ALT,
			expectedString: "ALT",
		},
		{
			name:           "LPAREN",
			t:              LPAREN,
			expectedString: "LPAREN",
		},
		{
			name:           "RPAREN",
			t:              RPAREN,
			expectedString: "RPAREN",
		},
		{
			name:           "LBRACK",
			t:              LBRACK,
			expectedString: "LBRACK",
		},
		{
			name:           "RBRACK",
			t:              RBRACK,
			expectedString: "RBRACK",
		},
		{
			name:           "LBRACE",
			t:              LBRACE,
			expectedString: "LBRACE",
		},
		{
			name:           "RBRACE",
			t:              RBRACE,
			expectedString: "RBRACE",
		},
		{
			name:           "LLBRACE",
			t:              LLBRACE,
			expectedString: "LLBRACE",
		},
		{
			name:           "RRBRACE",
			t:              RRBRACE,
			expectedString: "RRBRACE",
		},
		{
			name:           "GRAMMER",
			t:              GRAMMER,
			expectedString: "GRAMMER",
		},
		{
			name:           "IDENT",
			t:              IDENT,
			expectedString: "IDENT",
		},
		{
			name:           "TOKEN",
			t:              TOKEN,
			expectedString: "TOKEN",
		},
		{
			name:           "STRING",
			t:              STRING,
			expectedString: "STRING",
		},
		{
			name:           "REGEX",
			t:              REGEX,
			expectedString: "REGEX",
		},
		{
			name:           "Unknown",
			t:              Tag(99),
			expectedString: "Tag(99)",
		},
		{
			name:           "COMMENT",
			t:              COMMENT,
			expectedString: "COMMENT",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := fmt.Sprintf("%v", tc.t)
			assert.Equal(t, tc.expectedString, s)
		})
	}
}

func TestToken_String(t *testing.T) {
	tests := []struct {
		name           string
		t              Token
		expectedString string
	}{
		{
			name: "IDENT",
			t: Token{
				Tag:    IDENT,
				Lexeme: "foo",
				Pos:    4,
			},
			expectedString: "IDENT<foo,4>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := fmt.Sprintf("%v", tc.t)
			assert.Equal(t, tc.expectedString, s)
		})
	}
}
