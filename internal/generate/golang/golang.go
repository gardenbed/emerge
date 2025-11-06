// Package golang provides functionality for generating Go code that implements a full LALR parser based on an EBNF specification.
package golang

import (
	"bytes"
	"embed"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/range/disc"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

//go:embed templates/*.tmpl
var templates embed.FS

var (
	navajoWhite = ui.Fg256Color(223)
	darkOrange  = ui.Fg256Color(166)
	hotPink     = ui.Fg256Color(168)
	orchid      = ui.Fg256Color(170)
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

	if err := g.generateLexer(); err != nil {
		errs = errors.Append(errs, err)
	}

	if err := g.generateParser(); err != nil {
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

type coreData struct {
	Debug           bool
	Package         string
	GenerateCommand string
}

// generateCore generates essential data types and structures for the lexer and parser,
// ensuring no dependencies on third-party, non-built-in packages.
func (g *generator) generateCore() error {
	g.Infof(darkOrange, "     Generating core types ...")

	data := &coreData{
		Debug:           g.Debug,
		Package:         g.Spec.Name,
		GenerateCommand: strings.Join(os.Args, " "),
	}

	var errs error
	for _, filename := range []string{"core.go.tmpl", "grammar.go.tmpl"} {
		if err := g.renderTemplate(filename, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	return errs
}

type lexerData struct {
	Debug          bool
	Package        string
	Assocs         []spec.FinalTerminalAssociation
	DFATransitions iter.Seq2[automata.State, iter.Seq2[[]disc.Range[automata.Symbol], automata.State]]
}

// generateLexer generates the lexer code based on the provided terminal (token) definitions for the input language.
func (g *generator) generateLexer() error {
	g.Infof(hotPink, "     Generating the lexer ...")

	g.Infof(hotPink, "       Constructing finite automaton ...")
	dfa, assocs, err := g.Spec.BuildLexerDFA()
	if err != nil {
		return err
	}

	data := &lexerData{
		Debug:          g.Debug,
		Package:        g.Spec.Name,
		Assocs:         assocs,
		DFATransitions: dfa.Transitions(),
	}

	var errs error
	for _, filename := range []string{"input.go.tmpl", "lexer.go.tmpl"} {
		if err := g.renderTemplate(filename, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	// Generate the lexer graph if debugging is enabled.
	if err := g.generateLexerGraph(dfa, assocs); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs
}

// generateLexerGraph generates a DOT format graph of the lexer's DFA if debugging is enabled.
func (g *generator) generateLexerGraph(dfa *automata.DFA, assocs []spec.FinalTerminalAssociation) error {
	if !g.Params.Debug {
		return nil
	}

	g.Debugf(navajoWhite, "       Generating the lexer graph ...")

	// Write the DOT code to the file.
	filepath := filepath.Join(g.Path, g.Spec.Name, "lexer.dot")
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	// Generate the DOT code for the DFA.
	dot := dfa.DOT()

	// Modify the DOT code: colorize edges and their labels.
	nodeRE := regexp.MustCompile(`  node \[shape=circle\];`)
	dot = nodeRE.ReplaceAllString(dot, "  node [shape=circle];\n  edge [color=darkblue fontcolor=red];")

	// Modify the DOT code: annotate final states with terminal names.
	numRE := regexp.MustCompile(`\d+`)
	nodeStmtRE := regexp.MustCompile(`  \d+ \[label="\d+", shape=doublecircle\];`)

	dot = nodeStmtRE.ReplaceAllStringFunc(dot, func(m string) string {
		if i, err := strconv.Atoi(numRE.FindString(m)); err == nil {
			// Find the association with this final state.
			assoc, ok := generic.FirstMatch(assocs, func(assoc spec.FinalTerminalAssociation) bool {
				return assoc.Final.Contains(automata.State(i))
			})

			if ok {
				return fmt.Sprintf(`  %d [label="%d", shape=doublecircle style=filled color=skyblue xlabel=%q];`, i, i, assoc.Terminal)
			}
		}

		return m
	})

	if _, err := f.WriteString(dot); err != nil {
		return err
	}

	return nil
}

type parserData struct {
	Debug        bool
	Package      string
	Terminals    []grammar.Terminal
	NonTerminals []grammar.NonTerminal
	Productions  []*grammar.Production
	ParsingTable *lr.ParsingTable
}

// generateParser generates the parser code based on the provided grammar and precedence levels for the input language.
func (g *generator) generateParser() error {
	g.Infof(orchid, "     Generating the parser ...")

	g.Infof(orchid, "       Constructing LALR(1) Parsing Table ...")
	T, err := g.Spec.LALRParsingTable()
	if err != nil {
		return err
	}

	terminals := g.Spec.Grammar.OrderTerminals()
	_, _, nonTerminals := g.Spec.Grammar.OrderNonTerminals()
	productions := g.Spec.Grammar.OrderProductions()

	data := &parserData{
		Debug:        g.Debug,
		Package:      g.Spec.Name,
		Terminals:    terminals,
		NonTerminals: nonTerminals,
		Productions:  productions,
		ParsingTable: T,
	}

	var errs error
	for _, filename := range []string{"ast.go.tmpl", "parser.lalr.go.tmpl"} {
		if err := g.renderTemplate(filename, data); err != nil {
			errs = errors.Append(errs, err)
		}
	}

	// Generate the parsing table if debugging is enabled.
	if err := g.generateParsingTable(T); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs
}

// generateParsingTable generates a text file containing the parsing table if debugging is enabled.
func (g *generator) generateParsingTable(T *lr.ParsingTable) error {
	if !g.Params.Debug {
		return nil
	}

	g.Debugf(navajoWhite, "       Generating the parsing table ...")

	// Write the DOT code to the file.
	filepath := filepath.Join(g.Path, g.Spec.Name, "parser.txt")
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	// Generate the content for the parsing table.
	content := T.String()

	if _, err := f.WriteString(content); err != nil {
		return err
	}

	return nil
}

// renderTemplate renders an embedded template by name and
// writes the output to a file in the directory specified by Path and Package.
func (g *generator) renderTemplate(filename string, data any) error {
	g.Debugf(navajoWhite, "       Rendering %q ...", filename)

	content, err := templates.ReadFile(filepath.Join("templates", filename))
	if err != nil {
		return err
	}

	tmpl := template.New(filename).Funcs(template.FuncMap{
		"formatStates":        formatStates,
		"formatRanges":        formatRanges,
		"formatSymbolString":  formatSymbolString,
		"equalEndmarker":      equalEndmarker,
		"appendEndmarker":     appendEndmarker,
		"hasAnyACTION":        hasAnyACTION,
		"lookupACTION":        lookupACTION,
		"hasAnyGOTO":          hasAnyGOTO,
		"lookupGOTO":          lookupGOTO,
		"findProductionIndex": findProductionIndex,
	})

	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return err
	}

	filepath := filepath.Join(g.Path, g.Spec.Name, fmt.Sprintf("%s.go", g.Spec.Name))
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	return nil
}

func formatStates(states automata.States) string {
	var b bytes.Buffer

	for s := range states.All() {
		fmt.Fprintf(&b, "%d, ", s)
	}

	// Remove trailing comma and space
	if b.Len() >= 2 {
		b.Truncate(b.Len() - 2)
	}

	return b.String()
}

func formatRanges(ranges []disc.Range[automata.Symbol]) string {
	var b bytes.Buffer

	for _, r := range ranges {
		if r.Lo == r.Hi {
			fmt.Fprintf(&b, "r == %q, ", r.Lo)
		} else {
			fmt.Fprintf(&b, "%q <= r && r <= %q, ", r.Lo, r.Hi)
		}
	}

	// Remove trailing comma and space
	if b.Len() >= 2 {
		b.Truncate(b.Len() - 2)
	}

	return b.String()
}

func formatSymbolString(s grammar.String[grammar.Symbol]) string {
	var b bytes.Buffer

	for _, sym := range s {
		switch v := sym.(type) {
		case grammar.Terminal:
			fmt.Fprintf(&b, "Terminal(%s), ", v)
		case grammar.NonTerminal:
			fmt.Fprintf(&b, "NonTerminal(%q), ", v)
		}
	}

	// Remove trailing comma and space
	if b.Len() >= 2 {
		b.Truncate(b.Len() - 2)
	}

	return b.String()
}

func equalEndmarker(a grammar.Terminal) bool {
	return a == grammar.Endmarker
}

func appendEndmarker(terminals []grammar.Terminal) []grammar.Terminal {
	return append(terminals, grammar.Endmarker)
}

func hasAnyACTION(T *lr.ParsingTable, s lr.State, terminals []grammar.Terminal) bool {
	for _, a := range terminals {
		if _, err := T.ACTION(s, a); err == nil {
			return true
		}
	}

	return false
}

func lookupACTION(T *lr.ParsingTable, s lr.State, a grammar.Terminal) *lr.Action {
	if action, err := T.ACTION(s, a); err == nil {
		return action
	}

	return nil
}

func hasAnyGOTO(T *lr.ParsingTable, s lr.State, nonTerminals []grammar.NonTerminal) bool {
	for _, A := range nonTerminals {
		if _, err := T.GOTO(s, A); err == nil {
			return true
		}
	}

	return false
}

func lookupGOTO(T *lr.ParsingTable, s lr.State, A grammar.NonTerminal) string {
	if next, err := T.GOTO(s, A); err == nil {
		return fmt.Sprintf("%d", next)
	}

	return ""
}

func findProductionIndex(prods []*grammar.Production, prod *grammar.Production) int {
	for i, p := range prods {
		if p.Equal(prod) {
			return i
		}
	}

	return -1
}
