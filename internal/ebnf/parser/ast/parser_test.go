package ast

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectedError string
	}{
		{
			name:          "Invalid",
			filename:      "../../fixture/test.invalid.grammar",
			expectedError: "lexical error at ../../fixture/test.invalid.grammar:1:1:L",
		},
		{
			name:          "Success",
			filename:      "../../fixture/test.success.grammar",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)
			defer f.Close()

			root, err := Parse(tc.filename, f)

			if len(tc.expectedError) == 0 {
				assert.NotNil(t, root)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, root)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
