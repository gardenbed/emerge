package generate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/errors"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

//go:embed templates/*.tmpl
var templates embed.FS

var (
	darkOrange  = ui.Fg256Color(166)
	hotPink     = ui.Fg256Color(168)
	orchid      = ui.Fg256Color(170)
	navajoWhite = ui.Fg256Color(223)
)

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

	g.Debugf(navajoWhite, "     Checking package directory %q ...", g.Spec.Name)

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

	data := map[string]any{
		"Package": g.Spec.Name,
	}

	var errs error

	for _, name := range []string{"errors.go", "types.go", "stack.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// generateLexer generates the lexer code based on the provided terminal (token) definitions for the input language.
func (g *generator) generateLexer() error {
	g.Infof(hotPink, "     Generating the lexer ...")
	g.Infof(hotPink, "       Constructing DFA ...")

	_, _, err := g.Spec.DFA()
	if err != nil {
		return err
	}

	data := map[string]any{
		"Package": g.Spec.Name,
	}

	var errs error

	for _, name := range []string{"input.go", "lexer.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// generateParser generates the parser code based on the provided grammar and precedence levels for the input language.
func (g *generator) generateParser() error {
	g.Infof(orchid, "     Generating the parser ...")
	g.Infof(orchid, "       Constructing LALR(1) Parsing Table ...")

	_, err := g.Spec.LALRParsingTable()
	if err != nil {
		return err
	}

	data := map[string]any{
		"Package": g.Spec.Name,
	}

	var errs error

	for _, name := range []string{"parser.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

// renderTemplate renders an embedded template by name and
// writes the output to a file in the directory specified by Path and Package.
func (g *generator) renderTemplate(filename string, data any) error {
	g.Debugf(navajoWhite, "       Rendering %q ...", filename)

	content, err := templates.ReadFile(filepath.Join("templates", filename+".tmpl"))
	if err != nil {
		return err
	}

	tmpl, err := template.New(filename).Parse(string(content))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(g.Path, g.Spec.Name, filename), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
}
