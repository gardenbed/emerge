package command

import (
	"testing"

	"github.com/gardenbed/charm/ui"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ui := ui.NewNop()
	c := New(ui)

	assert.NotNil(t, c)
}

func TestNewFactory(t *testing.T) {
	ui := ui.NewNop()
	c, err := NewFactory(ui)()

	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestCommand_Synopsis(t *testing.T) {
	c := new(Command)
	synopsis := c.Synopsis()

	assert.NotEmpty(t, synopsis)
}

func TestCommand_Help(t *testing.T) {
	c := new(Command)
	help := c.Help()

	assert.NotEmpty(t, help)
}

func TestCommand_Run(t *testing.T) {
	t.Run("InvalidFlag", func(t *testing.T) {
		c := &Command{ui: ui.NewNop()}
		exitCode := c.Run([]string{"-undefined"})

		assert.Equal(t, FlagError, exitCode)
	})

	t.Run("OK", func(t *testing.T) {
		c := &Command{ui: ui.NewNop()}
		exitCode := c.Run([]string{})

		assert.Equal(t, Success, exitCode)
	})
}

func TestCommand_parseFlags(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedExitCode int
	}{
		{
			name:             "InvalidFlag",
			args:             []string{"-undefined"},
			expectedExitCode: FlagError,
		},
		{
			name:             "NoFlag",
			args:             []string{},
			expectedExitCode: Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{ui: ui.NewNop()}
			exitCode := c.parseFlags(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}

func TestCommand_exec(t *testing.T) {
	tests := []struct {
		name             string
		expectedExitCode int
	}{
		{
			name:             "Success",
			expectedExitCode: Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui: ui.NewNop(),
			}

			exitCode := c.exec()

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
