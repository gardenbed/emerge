package parser

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		src           io.Reader
		expectedError string
	}{
		{
			name:          "Success",
			src:           strings.NewReader("Lorem ipsum"),
			expectedError: "",
		},
		{
			name:          "Failure",
			src:           iotest.ErrReader(errors.New("io error")),
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			par, err := New(tc.src)

			if tc.expectedError == "" {
				assert.NotNil(t, par)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, par)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
