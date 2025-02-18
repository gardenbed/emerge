package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name                 string
		filename             string
		expectedSpec         *Spec
		expectedErrorStrings []string
	}{
		{
			name:     "Invalid",
			filename: "../../fixture/test.invalid.grammar",
			expectedErrorStrings: []string{
				`lexical error at ../../fixture/test.invalid.grammar:1:1:L`,
			},
		},
		{
			name:     "Error",
			filename: "../../fixture/test.error.grammar",
			expectedErrorStrings: []string{
				`3 errors occurred:`,
				`invalid predefined regex: $IDN`,
				`no definition for terminal "ID"`,
				`missing production rule with the start symbol: start`,
			},
		},
		{
			name:     "Success",
			filename: "../../fixture/test.success.grammar",
			expectedSpec: &Spec{
				Name:        "test",
				Definitions: []*TerminalDef{},
				Grammar:     grammars[1],
				Precedences: precedences[1],
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)
			defer f.Close()

			spec, err := Parse(tc.filename, f)

			if len(tc.expectedErrorStrings) > 0 {
				assert.Nil(t, spec)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			} else {
				assert.NotNil(t, spec)
				assert.NoError(t, err)

				assert.True(t, spec.Name == tc.expectedSpec.Name)
				assert.NotNil(t, spec.Definitions)
				assert.True(t, spec.Grammar.Equal(tc.expectedSpec.Grammar))
				assert.True(t, spec.Precedences.Equal(tc.expectedSpec.Precedences))
			}
		})
	}
}
