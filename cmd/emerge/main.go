package main

import (
	"flag"
	"fmt"
	"os"

	// "net/http"
	// _ "net/http/pprof"

	"github.com/gardenbed/charm/flagit"
	"github.com/gardenbed/charm/ui"

	"github.com/gardenbed/emerge/internal/command"
	"github.com/gardenbed/emerge/metadata"
)

func main() {
	// Start pprof server.
	// Visit http://localhost:6060/debug/pprof/
	/* go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}() */

	if code := run(); code != 0 {
		os.Exit(code)
	}
}

// run is the main entry point for the emerge command.
func run() int {
	u := ui.New(ui.Info)
	fs := flag.NewFlagSet("emerge", flag.ContinueOnError)

	cmd, err := command.New(u)
	if err != nil {
		u.Errorf(ui.Red, "%s", err)
		return 1
	}

	if err := flagit.Register(fs, cmd, false); err != nil {
		u.Errorf(ui.Red, "%s", err)
		return 1
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		u.Errorf(ui.Red, "%s", err)
		return 1
	}

	// Update the verbosity level.
	if cmd.Verbose {
		u.SetLevel(ui.Debug)
	}

	switch {
	case cmd.Help:
		if err := cmd.PrintHelp(); err != nil {
			u.Errorf(ui.Red, "%s", err)
			return 1
		}

	case cmd.Version:
		fmt.Println(metadata.String())

	default:
		if err := cmd.Run(fs.Args()); err != nil {
			u.Errorf(ui.Red, "\n%s\n", err)
			return 1
		}
	}

	return 0
}
