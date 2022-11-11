package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {
	tests := []struct {
		name string
		s    string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := NewInput(tc.s)
			assert.NotNil(t, in)

			for in != nil {
				r, pos := in.Current()
				assert.Equal(t, tc.s[pos], r)

				in = in.Remaining()
			}
		})
	}
}
