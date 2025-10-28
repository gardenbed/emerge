// Package golang provides functionality for generating Go code that implements a full LALR parser based on an EBNF specification.
package golang

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

//go:embed templates/*.tmpl
var templates embed.FS

var (
	navajoWhite = ui.Fg256Color(223)
	darkOrange  = ui.Fg256Color(166)
	// hotPink     = ui.Fg256Color(168)
	// orchid      = ui.Fg256Color(170)
)

var (
	idRegex = regexp.MustCompile(`^[\p{L}_][\p{L}\p{Nd}_]*$`)

	builtin = []string{
		// Keywords
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
		// Types
		"any", "bool", "byte", "comparable",
		"complex64", "complex128", "error", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		// Constants
		"true", "false", "iota",
		// Zero value
		"nil",
		// Functions
		"append", "cap", "clear", "close", "complex", "copy", "delete", "imag", "len",
		"make", "max", "min", "new", "panic", "print", "println", "real", "recover",
	}
)

// isIDValid checks if a name is a valid identifier in Go.
func isIDValid(name string) bool {
	return idRegex.MatchString(name) && !generic.AnyMatch(builtin, func(s string) bool {
		return s == name
	})
}

// generator holds the data and objects that are shared and needed by the generator methods.
type generator struct {
	ui.UI
	*Params
}

// Params contains the configuration and data required for generating the parser code.
type Params struct {
	Debug bool
	Path  string
	Spec  *spec.Spec
}

// Generate creates a self-contained, complete package that implements a full LALR parser for the input language,
// including all necessary data types, data structures, and a lexer (a.k.a. scanner).
func Generate(u ui.UI, params *Params) error {
	g := &generator{
		UI:     u,
		Params: params,
	}

	if err := g.prepare(); err != nil {
		return err
	}

	var errs error

	if err := g.generateCore(); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs
}

// prepare validates the params and ensures the required directory structure exists before generating package code.
func (g *generator) prepare() error {
	g.Path = filepath.Clean(g.Path)

	g.Debugf(navajoWhite, "     Checking output path %q ...", g.Path)

	// Ensure the output path exists.
	info, err := os.Stat(g.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("output path does not exist: %q", g.Path)
		}
		return fmt.Errorf("error on checking output path: %s", err)
	}

	// Ensure the output path is a directory.
	if !info.IsDir() {
		return fmt.Errorf("output path is not a directory: %q", g.Path)
	}

	g.Debugf(navajoWhite, "     Creating package directory %q ...", g.Spec.Name)

	if !isIDValid(g.Spec.Name) {
		return fmt.Errorf("invalid package name: %s", g.Spec.Name)
	}

	packageDir := filepath.Join(g.Path, g.Spec.Name)

	// Create the package directory.
	if err := os.Mkdir(packageDir, os.ModePerm); err != nil {
		return fmt.Errorf("error on creating package directory: %s", err)
	}

	return nil
}

// generateCore generates essential data types and structures for the lexer and parser,
// ensuring no dependencies on third-party, non-built-in packages.
func (g *generator) generateCore() error {
	g.Infof(darkOrange, "     Generating core types ...")

	data := &coreData{
		GenerateCommand: strings.Join(os.Args, " "),
		Package:         g.Spec.Name,
	}

	var errs error
	if err := g.renderTemplate("core.go.tmpl", data); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs
}

type coreData struct {
	GenerateCommand string
	Package         string
}

// renderTemplate renders an embedded template by name and
// writes the output to a file in the directory specified by Path and Package.
func (g *generator) renderTemplate(filename string, data any) error {
	g.Debugf(navajoWhite, "       Rendering %q ...", filename)

	content, err := templates.ReadFile(filepath.Join("templates", filename))
	if err != nil {
		return err
	}

	tmpl := template.New(filename).Funcs(template.FuncMap{})

	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return err
	}

	filepath := filepath.Join(g.Path, g.Spec.Name, fmt.Sprintf("%s.go", g.Spec.Name))
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
}
