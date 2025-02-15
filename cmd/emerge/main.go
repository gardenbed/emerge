package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gardenbed/charm/flagit"
	"github.com/gardenbed/charm/ui"

	"github.com/gardenbed/emerge/internal/command"
	"github.com/gardenbed/emerge/metadata"
)

func main() {
	u := ui.New(ui.Info)
	fs := flag.NewFlagSet("emerge", flag.ContinueOnError)

	cmd, err := command.New(u)
	if err != nil {
		u.Errorf(ui.Red, "%s", err)
		os.Exit(1)
	}

	if err := flagit.Register(fs, cmd, false); err != nil {
		u.Errorf(ui.Red, "%s", err)
		os.Exit(1)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	// Update the verbosity level.
	if cmd.Verbose {
		u.SetLevel(ui.Debug)
	}

	switch {
	case cmd.Help:
		if err := cmd.PrintHelp(); err != nil {
			u.Errorf(ui.Red, "%s", err)
			os.Exit(1)
		}

	case cmd.Version:
		fmt.Println(metadata.String())

	default:
		if err := cmd.Run(fs.Args()); err != nil {
			u.Errorf(ui.Red, "\n%s\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
