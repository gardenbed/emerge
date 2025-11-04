package ast

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name                 string
		filename             string
		expectedErrorStrings []string
	}{
		{
			name:     "Invalid",
			filename: "../../fixture/test.invalid.grammar",
			expectedErrorStrings: []string{
				`unexpected string "L": no action exists in the parsing table for ACTION[0, "TOKEN"]`,
			},
		},
		{
			name:                 "Success",
			filename:             "../../fixture/test.success.grammar",
			expectedErrorStrings: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)

			defer func() {
				assert.NoError(t, f.Close())
			}()

			root, err := Parse(tc.filename, f)

			if len(tc.expectedErrorStrings) == 0 {
				assert.NotNil(t, root)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, root)
				assert.Error(t, err)

				s := err.Error()
				for _, expectedErrorString := range tc.expectedErrorStrings {
					assert.Contains(t, s, expectedErrorString)
				}
			}
		})
	}
}
