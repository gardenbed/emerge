package ebnf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInputBuffer(t *testing.T) {

}

func TestInputBuffer_getNextChar(t *testing.T) {
	tests := []struct {
		name string
		size int
		file string
	}{
		{
			name: "Please",
			size: 512,
			file: "./fixture/please.grammar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			_, err = newInputBuffer(tc.size, f)
			assert.NoError(t, err)
		})
	}
}

func TestTest(t *testing.T) {
	f, err := os.Open("./fixture/please.grammar")
	assert.NoError(t, err)
	defer f.Close()

	in, err := newInputBuffer(1024, f)
	assert.NoError(t, err)

	for in.Next() {
		fmt.Print(string(in.Char()))
	}

	assert.NoError(t, in.Err())
}
