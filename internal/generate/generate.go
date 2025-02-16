package generate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gardenbed/charm/ui"
	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Params contains the configuration and data required for generating the parser code.
type Params struct {
	Debug        bool
	Path         string
	Package      string
	DFA          *auto.DFA
	Productions  []*grammar.Production
	ParsingTable *lr.ParsingTable
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

	if err := g.generateLexer(); err != nil {
		errs = errors.Append(errs, err)
	}

	if err := g.generateParser(); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs
}

// generator holds the data and objects that are shared and needed by the generator methods.
type generator struct {
	ui.UI
	*Params
}

// prepare validates the params and ensures the required directory structure exists before generating package code.
func (g *generator) prepare() error {
	g.Path = filepath.Clean(g.Path)

	g.Debugf(ui.Cyan, "     Checking output path %q ...", g.Path)

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

	g.Debugf(ui.Cyan, "     Checking package directory %q ...", g.Package)

	if !isIDValid(g.Package) {
		return fmt.Errorf("invalid package name: %s", g.Package)
	}

	packageDir := filepath.Join(g.Path, g.Package)

	// Create the package directory.
	if err := os.Mkdir(packageDir, os.ModePerm); err != nil {
		return fmt.Errorf("error on creating package directory: %s", err)
	}

	return nil
}

// generateCore generates essential data types and structures for the lexer and parser,
// ensuring no dependencies on third-party, non-built-in packages.
func (g *generator) generateCore() error {
	var errs error

	g.Infof(ui.Yellow, "     Generating core types ...")

	data := map[string]any{
		"Package": g.Package,
	}

	for _, name := range []string{"errors.go", "types.go", "stack.go"} {
		g.Debugf(ui.Cyan, "       Generating %q ...", name)

		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// generateLexer generates the lexer code based on the provided terminal (token) definitions for the input language.
func (g *generator) generateLexer() error {
	var errs error

	g.Infof(ui.Yellow, "     Generating the lexer ...")

	data := map[string]any{
		"Package": g.Package,
	}

	for _, name := range []string{"input.go", "lexer.go"} {
		g.Debugf(ui.Cyan, "       Generating %q ...", name)

		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// generateParser generates the parser code based on the provided grammar and precedence levels for the input language.
func (g *generator) generateParser() error {
	var errs error

	g.Infof(ui.Yellow, "     Generating the parser ...")

	data := map[string]any{
		"Package": g.Package,
	}

	for _, name := range []string{"parser.go"} {
		g.Debugf(ui.Cyan, "       Generating %q ...", name)

		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// renderTemplate renders an embedded template by name and
// writes the output to a file in the directory specified by Path and Package.
func (g *generator) renderTemplate(filename string, data any) error {
	content, err := templates.ReadFile(filepath.Join("templates", filename+".tmpl"))
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(string(content))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(g.Path, g.Package, filename), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
}
