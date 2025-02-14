package result

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name                 string
		filename             string
		expectedResult       *Result
		expectedErrorStrings []string
	}{
		{
			name:     "Invalid",
			filename: "../../fixture/invalid.grammar",
			expectedErrorStrings: []string{
				`lexical error at ../../fixture/invalid.grammar:1:1:L`,
			},
		},
		{
			name:     "Error",
			filename: "../../fixture/test.error.grammar",
			expectedErrorStrings: []string{
				`5 errors occurred:`,
				`invalid predefined regex: $IDN`,
				`"NUMBER": invalid regular expression: [0-9+`,
				`no definition for terminal "NUMBER"`,
				`no definition for terminal "ID"`,
				`missing production rule with the start symbol: start`,
			},
		},
		{
			name:     "Success",
			filename: "../../fixture/test.success.grammar",
			expectedResult: &Result{
				Name:        "test",
				Definitions: []*TerminalDef{},
				Grammar:     grammars[0],
				Precedences: precedences[0],
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)
			defer f.Close()

			res, err := Parse(tc.filename, f)

			if len(tc.expectedErrorStrings) > 0 {
				assert.Nil(t, res)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			} else {
				assert.NotNil(t, res)
				assert.NoError(t, err)

				assert.True(t, res.Name == tc.expectedResult.Name)
				assert.NotNil(t, res.Definitions)
				assert.True(t, res.Grammar.Equal(tc.expectedResult.Grammar))
				assert.True(t, res.Precedences.Equal(tc.expectedResult.Precedences))
			}
		})
	}
}
