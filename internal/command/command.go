package command

import (
	"os"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/gardenbed/charm/ui"
)

const helpTemplate = `
  {{green "Emerge"}} is a parser generator that produces an {{magenta "LALR(1)"}} parser for Go from an {{magenta "EBNF"}} specification of a context-free grammar.
  It also generates a lexical analyzer and all necessary auxiliary types.
  The generated code is {{magenta "self-contained"}} and does not rely on any external third-party modules.

  {{yellow "What is EBNF?"}}

  EBNF (Extended Backus-Naur Form) is an extension of BNF (Backus-Naur Form) used to specify context-free grammars (CFGs).
  A context-free grammar describes Type-2 languages in the Chomsky hierarchy,
  making it more expressive than Type-3 languages (regular languages, a.k.a. those described by regular expressions).

  EBNF provides additional expressive features such as:

    • {{blue "Alternation"}} ...|...
    • {{blue "Grouping"}} (...)
    • {{blue "Optional"}} [...]
    • {{blue "Kleene star closure"}} {...}
    • {{blue "Kleene plus closure"}} {{"{{"}}...{{"}}"}}

  {{yellow "Lexer and Tokenization"}}

  {{green "Emerge"}} generates a lexer (lexical analyzer or scanner) that tokenizes input source code.
  A lexer processes a stream of characters and converts them into a stream of tokens—the fundamental units of a language.
  These tokens are then passed to the parser, which performs syntax analysis and constructs parse trees.

  Tokens can be defined in two ways:

    1. {{blue "Implicitly"}}—by specifying string values.
    2. {{blue "Explicitly"}}—by defining a token name and a token definition using strings or regular expressions.

  {{yellow "LALR(1) Parsing"}}

  {{green "Emerge"}} generates an {{blue "LALR(1)"}} parser, a type of {{blue "bottom-up"}} parser designed for {{blue "LR(1)"}} languages.
  LR parsing methods can handle a larger class of grammars than LL (predictive top-down) parsing methods.
  An LR parser is expressive enough to recognize almost all programming language constructs described by context-free grammars.
  The generated parser uses a precomputed parsing table and operates in linear time {{blue "O(n)"}}.

  {{yellow "Handling Ambiguity"}}

  {{green "Emerge"}} supports ambiguous grammars and resolves ambiguity using explicit {{blue "associativity"}} and {{blue "precedence"}} directives.
  Unlike some parser generators, {{green "Emerge"}} does not rely on implicit precedence assignment via ordered choice.
  The grammar definition follows standard BNF semantics, while precedence and associativity are handled separately from the grammar itself.

  {{yellow "Usage:"}}  {{ green "emerge [flags] FILE_PATH"}}

  {{yellow "Flags:"}}

    -out=path    Generate the parser in the specified directory.
    -name=foo    Generate the parser with the specified name and ignore the name in the grammar specification.
    -debug       Generate the parser with extra types and methods for debugging and troubleshooting purposes.

    -help        Show the help text
    -version     Print the version number
    -verbose     Show the vervbosity logs (default: {{.Verbose}})

  {{yellow "Examples:"}}

    emerge grammar.ebnf
    emerge -out="~/src/project/internal" grammar.ebnf
    emerge -name="parser" grammar.ebnf
    emerge -debug grammar.ebnf

`

// Command represents the "emerge" command and its associated flags.
type Command struct {
	ui.UI

	Help    bool `flag:"help"`
	Version bool `flag:"version"`
	Verbose bool `flag:"verbose"`

	Out   string `flag:"out"`
	Name  string `flag:"name"`
	Debug bool   `flag:"debug"`
}

// New creates a new instance of the command.
func New(u ui.UI) *Command {
	return &Command{
		UI: u,
	}
}

// PrintHelp prints the help text for the command.
func (c *Command) PrintHelp() error {
	tmpl := template.New("help")
	tmpl = tmpl.Funcs(template.FuncMap{
		"join":    strings.Join,
		"blue":    color.New(color.FgBlue).Sprintf,
		"green":   color.New(color.FgGreen).Sprintf,
		"cyan":    color.New(color.FgCyan).Sprintf,
		"magenta": color.New(color.FgMagenta).Sprintf,
		"red":     color.New(color.FgRed).Sprintf,
		"white":   color.New(color.FgWhite).Sprintf,
		"yellow":  color.New(color.FgYellow).Sprintf,
	})

	tmpl, err := tmpl.Parse(helpTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, c)
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) error {
	c.Infof(ui.Yellow, "🚧 WIP")

	// TODO: Implement!

	return nil
}
