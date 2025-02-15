package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		name         string
		expectedBool bool
	}{
		{
			name:         "4ever",
			expectedBool: false,
		},
		{
			name:         "complex",
			expectedBool: false,
		},
		{
			name:         "pascal",
			expectedBool: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedBool, isIDValid(tc.name))
		})
	}
}
