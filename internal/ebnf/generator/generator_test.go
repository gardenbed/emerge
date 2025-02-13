package generator

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
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
		name                string
		filename            string
		expectedName        string
		expectedGrammar     *grammar.CFG
		expectedPrecedences lr.PrecedenceLevels
		expectedError       string
	}{
		{
			name:          "Invalid",
			filename:      "../fixture/invalid.grammar",
			expectedError: "lexical error at ../fixture/invalid.grammar:1:1:L",
		},
		{
			name:                "Success",
			filename:            "../fixture/test.grammar",
			expectedName:        "test",
			expectedGrammar:     grammars[0],
			expectedPrecedences: precedences[0],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filename)
			assert.NoError(t, err)
			defer f.Close()

			g, err := New(tc.filename, f)
			assert.NoError(t, err)

			artifacts, err := g.parse()

			if tc.expectedError != "" {
				assert.Nil(t, artifacts)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NotNil(t, artifacts)
				assert.NoError(t, err)

				assert.True(t, artifacts.Name == tc.expectedName)
				assert.True(t, artifacts.CFG.Equal(tc.expectedGrammar))
				assert.True(t, artifacts.Precedences.Equal(tc.expectedPrecedences))
			}
		})
	}
}
