package command

import (
	"flag"

	"github.com/gardenbed/charm/ui"
	"github.com/mitchellh/cli"
)

const (
	// Success is the exit code when a command execution is successful.
	Success int = iota
	// GenericError is the generic exit code when something fails.
	GenericError
	// FlagError is the exit code when an undefined or invalid flag is provided to a command.
	FlagError
)

const (
	synopsis = `Generate a parser from a grammar specification`
	help     = `
  Use this command for ...

  Usage:

  Examples:
  `
)

// Command is the cli.Command implementation for emerge command.
type Command struct {
	ui ui.UI
}

// New creates a new command.
func New(ui ui.UI) *Command {
	return &Command{
		ui: ui,
	}
}

// NewFactory returns a cli.CommandFactory for creating a new command.
func NewFactory(ui ui.UI) cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(ui), nil
	}
}

// Synopsis returns a short one-line synopsis for the command.
func (c *Command) Synopsis() string {
	return synopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return help
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	if code := c.parseFlags(args); code != Success {
		return code
	}

	return c.exec()
}

func (c *Command) parseFlags(args []string) int {
	fs := flag.NewFlagSet("emerge", flag.ContinueOnError)

	fs.Usage = func() {
		c.ui.Printf(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		// In case of error, the error and help will be printed by the Parse method
		return FlagError
	}

	return Success
}

// exec in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) exec() int {
	c.ui.Warnf(ui.Yellow, "🚧 WIP")

	return Success
}
