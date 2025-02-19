package generate

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gardenbed/charm/ui"
	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/symboltable"

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

	data := &coreData{
		Package: g.Spec.Name,
	}

	var errs error
	for _, name := range []string{"errors.go", "types.go", "stack.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

type (
	coreData struct {
		Package string
	}
)

// generateLexer generates the lexer code based on the provided terminal (token) definitions for the input language.
func (g *generator) generateLexer() error {
	g.Infof(hotPink, "     Generating the lexer ...")

	g.Infof(hotPink, "       Constructing DFA ...")
	dfa, termMap, err := g.Spec.DFA()
	if err != nil {
		return err
	}

	// Group all the transitions with the same from and to states and combine their symbols.
	groups := groupDFAStates(dfa)

	data := &lexerData{
		Package: g.Spec.Name,
		DFA: &DFA{
			Transitions: make([]*DFATransition, 0),
			FinalStates: make([]*DFAFinalStates, len(g.Params.Spec.Definitions)),
		},
	}

	// Populate data.DFA.Trans
	for from, group := range groups.All() {
		t := &DFATransition{
			From:  from,
			Trans: make([]*DFAStateTransition, 0, group.Size()),
		}

		for to, syms := range group.All() {
			t.Trans = append(t.Trans, &DFAStateTransition{
				Symbols: syms,
				Next:    to,
			})
		}

		data.DFA.Transitions = append(data.DFA.Transitions, t)
	}

	// Populate data.DFA.Final
	for i, def := range g.Params.Spec.Definitions {
		term := def.Terminal
		states := termMap[term]

		data.DFA.FinalStates[i] = &DFAFinalStates{
			Terminal: string(term),
			States: generic.Transform(states, func(s auto.State) int {
				return int(s)
			}),
		}
	}

	var errs error
	for _, name := range []string{"input.go", "lexer.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

type (
	lexerData struct {
		Package string
		DFA     *DFA
	}

	DFA struct {
		Transitions []*DFATransition
		FinalStates []*DFAFinalStates
	}

	DFATransition struct {
		From  int
		Trans []*DFAStateTransition
	}

	DFAStateTransition struct {
		Symbols []rune
		Next    int
	}

	DFAFinalStates struct {
		Terminal string
		States   []int
	}
)

func groupDFAStates(dfa *auto.DFA) symboltable.SymbolTable[int, symboltable.SymbolTable[int, []rune]] {
	cmpState := generic.NewCompareFunc[int]()
	groups := symboltable.NewRedBlack[int, symboltable.SymbolTable[int, []rune]](cmpState, nil)

	for from, ftrans := range dfa.Trans.All() {
		group, ok := groups.Get(int(from))
		if !ok {
			group = symboltable.NewRedBlack[int, []rune](cmpState, nil)
			groups.Put(int(from), group)
		}

		for sym, to := range ftrans.All() {
			symbols, _ := group.Get(int(to))
			symbols = append(symbols, rune(sym))
			group.Put(int(to), symbols)
		}
	}

	return groups
}

// generateParser generates the parser code based on the provided grammar and precedence levels for the input language.
func (g *generator) generateParser() error {
	g.Infof(orchid, "     Generating the parser ...")

	g.Infof(orchid, "       Constructing LALR(1) Parsing Table ...")
	_, err := g.Spec.LALRParsingTable()
	if err != nil {
		return err
	}

	data := &parserData{
		Package: g.Spec.Name,
	}

	var errs error
	for _, name := range []string{"parser.go"} {
		if err := g.renderTemplate(name, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

type (
	parserData struct {
		Package string
	}
)

// renderTemplate renders an embedded template by name and
// writes the output to a file in the directory specified by Path and Package.
func (g *generator) renderTemplate(filename string, data any) error {
	g.Debugf(navajoWhite, "       Rendering %q ...", filename)

	content, err := templates.ReadFile(filepath.Join("templates", filename+".tmpl"))
	if err != nil {
		return err
	}

	tmpl := template.New(filename).Funcs(template.FuncMap{
		"formatInts":  formatInts,
		"formatRunes": formatRunes,
	})

	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return err
	}

	filepath := filepath.Join(g.Path, g.Spec.Name, filename)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
}

func formatInts(vals []int) string {
	var b bytes.Buffer

	for _, v := range vals {
		fmt.Fprintf(&b, "%d, ", v)
	}

	if len(vals) > 0 {
		b.Truncate(b.Len() - 2)
	}

	return b.String()
}

func formatRunes(runes []rune) string {
	var b bytes.Buffer

	for _, r := range runes {
		fmt.Fprintf(&b, "'%c', ", r)
	}

	if len(runes) > 0 {
		b.Truncate(b.Len() - 2)
	}

	return b.String()
}
