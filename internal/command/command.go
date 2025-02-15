package command

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/gardenbed/charm/ui"
	"github.com/moorara/algo/generic"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
	"github.com/gardenbed/emerge/internal/generate"
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

  {{yellow "How To Use?"}}

  You can use the generated parser in different ways:

    • {{blue "Bottom-up Processing:"}} Handle each production rule as it is recognized in a bottom-up fashion.
    • {{blue "Abstract Syntax Tree (AST) Traversal:"}} Alternatively, you can traverse the AST generated by the parser.

  {{yellow "Documentation"}}

  Explore Emerge, learn EBNF, and browse examples at:

    {{cyan "🔗 https://gardenbed.github.io/emerge"}}

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

type (
	// Command represents the "emerge" command and its associated flags.
	Command struct {
		ui.UI
		funcs

		Help    bool `flag:"help"`
		Version bool `flag:"version"`
		Verbose bool `flag:"verbose"`

		Out   string `flag:"out"`
		Name  string `flag:"name"`
		Debug bool   `flag:"debug"`
	}

	// funcs defines the function types required by the command.
	// This abstraction allows these functions to be mocked for testing purposes.
	funcs struct {
		Parse    func(string, io.Reader) (*spec.Spec, error)
		Generate func(ui.UI, *generate.Params) error
	}
)

// New creates a new instance of the command.
func New(u ui.UI) (*Command, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	c := &Command{
		UI:  u,
		Out: path,
	}

	c.funcs.Parse = spec.Parse
	c.funcs.Generate = generate.Generate

	return c, nil
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
	path, ok := generic.FirstMatch(args, func(a string) bool {
		return !strings.HasPrefix(a, "-")
	})

	if !ok {
		return errors.New("no input file specified, please provide a file path")
	}

	filename := filepath.Base(path)

	c.Infof(nil, "%c Opening %q ...", getPlant(), filename)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	c.Infof(nil, "%c Reading %q ...", getFruit(), filename)

	spec, err := c.funcs.Parse(filename, f)
	if err != nil {
		return err
	}

	c.Infof(ui.Magenta, "%c Generating parser ...", getAnimal())

	err = c.funcs.Generate(c.UI, &generate.Params{
		Path:    c.Out,
		Package: c.Name,
		Debug:   c.Debug,
		Spec:    spec,
	})

	if err != nil {
		return err
	}

	c.Infof(ui.Green, "%c Successful!", getFood())

	return nil
}

var emojis = map[string]string{
	"animals": "🐶🐱🐭🐹🐰🦊🐻🐼🐻‍❄️🐨🐯🦁🐮🐷🐸🐵🐔🐧🐦🐤🐴🦄🐝🐛🦋🐌🐞🐙🦞🐠🐬🦧🦚🦜🦢🦩🐿️🐲🐦‍🔥",
	"plants":  "🌵🌲🌳🌴🌱🌿☘️🍀🪴🎋🍃🍄🍄‍🟫🌾💐🌷🌹🥀🪻🪷🌺🌸🌼🌻",
	"fruits":  "🍏🍎🍐🍊🍋🍋‍🟩🍌🍉🍇🍓🫐🍈🍒🍑🥭🍍🥥🥝🥑🥒🌶️🫑🌽🥕🫒",
	"food":    "🥐🥯🥨🧀🥞🧇🍕🥙🧆🌮🌯🫔🥗🍝🍜🍣🥟🍤🍚🍥🥠🥮🍡🍧🍨🍦🧁🍰🍯",
}

func getAnimal() rune {
	animals := []rune(emojis["animals"])
	i := rand.Intn(len(animals))
	return animals[i]
}

func getPlant() rune {
	plants := []rune(emojis["plants"])
	i := rand.Intn(len(plants))
	return plants[i]
}

func getFruit() rune {
	fruits := []rune(emojis["fruits"])
	i := rand.Intn(len(fruits))
	return fruits[i]
}

func getFood() rune {
	food := []rune(emojis["food"])
	i := rand.Intn(len(food))
	return food[i]
}
