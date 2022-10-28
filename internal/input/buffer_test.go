package input

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {
	tests := []struct {
		name          string
		n             int
		file          string
		expectedError string
	}{
		{
			name:          "Success",
			n:             4096,
			file:          "./fixture/lorem_ipsum",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			_, err = NewBuffer(tc.n, f)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestBuffer(t *testing.T) {
	tests := []struct {
		name string
		n    int
		file string
	}{
		{
			name: "Success",
			n:    1024,
			file: "./fixture/lorem_ipsum",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := NewBuffer(tc.n, f)
			assert.NoError(t, err)

			for in.Next() {
				b := in.Char()
				assert.NotEmpty(t, b)
			}

			assert.NoError(t, in.Err())
		})
	}
}
