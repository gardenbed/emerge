package generator

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		src           io.Reader
		expectedError string
	}{
		{
			name:          "Success",
			filename:      "lorem_ipsum",
			src:           strings.NewReader("Lorem ipsum"),
			expectedError: "",
		},
		{
			name:          "Failure",
			filename:      "lorem_ipsum",
			src:           iotest.ErrReader(errors.New("io error")),
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := New(tc.filename, tc.src)

			if tc.expectedError == "" {
				assert.NotNil(t, gen)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, gen)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGenerator_parse(t *testing.T) {
	tests := []struct {
		name                 string
		filename             string
		expectedResult       *result
		expectedErrorStrings []string
	}{
		{
			name:     "Invalid",
			filename: "../fixture/invalid.grammar",
			expectedErrorStrings: []string{
				`lexical error at ../fixture/invalid.grammar:1:1:L`,
			},
		},
		{
			name:     "Error",
			filename: "../fixture/test.error.grammar",
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
			filename: "../fixture/test.success.grammar",
			expectedResult: &result{
				Name:        "test",
				Definitions: []*terminalDef{},
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

			g, err := New(tc.filename, f)
			assert.NoError(t, err)

			res, err := g.parse()

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
