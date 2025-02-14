package command

import (
	"testing"

	"github.com/gardenbed/charm/ui"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		u := ui.NewNop()
		cmd := New(u)

		assert.NotNil(t, cmd)
	})
}

func TestCommand_PrintHelp(t *testing.T) {
	tests := []struct {
		name          string
		c             *Command
		expectedError string
	}{
		{
			name: "OK",
			c: &Command{
				UI:      ui.NewNop(),
				Verbose: true,
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.PrintHelp()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCommand_Run(t *testing.T) {
	tests := []struct {
		name          string
		c             *Command
		args          []string
		expectedError string
	}{
		{
			name: "OK",
			c: &Command{
				UI:      ui.NewNop(),
				Verbose: true,
			},
			args:          []string{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.Run(tc.args)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
