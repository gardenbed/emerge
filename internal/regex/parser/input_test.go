package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInput(t *testing.T) {
	tests := []struct {
		name string
		s    string
	}{
		{
			name: "OK",
			s:    "Lorem ipsum",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := newStringInput(tc.s)
			assert.NotNil(t, in)

			runes := []rune(tc.s)

			for in != nil {
				r, pos := in.Current()
				assert.Equal(t, runes[pos], r)

				in = in.Remaining()
			}
		})
	}
}
