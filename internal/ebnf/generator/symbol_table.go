package generator

import (
	"fmt"
	"strings"
	"sync"

	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/symboltable"

	"github.com/gardenbed/emerge/internal/regex/parser/nfa"
)

type terminalDef interface {
	Pos() *lexer.Position
	DFA() (*auto.DFA, error)
}

type stringTerminalDef struct {
	value string
	pos   *lexer.Position
}

func (d *stringTerminalDef) Pos() *lexer.Position {
	return d.pos
}

func (d *stringTerminalDef) DFA() (*auto.DFA, error) {
	start := auto.State(0)
	dfa := auto.NewDFA(start, nil)

	curr, next := start, start+1
	for _, r := range d.value {
		dfa.Add(curr, auto.Symbol(r), next)
		curr, next = next, next+1
	}

	dfa.Final = auto.NewStates(curr)

	return dfa, nil
}

type regexTerminalDef struct {
	regex string
	pos   *lexer.Position
}

func (d *regexTerminalDef) Pos() *lexer.Position {
	return d.pos
}

func (d *regexTerminalDef) DFA() (*auto.DFA, error) {
	nfa, err := nfa.Parse(d.regex)
	if err != nil {
		return nil, err
	}

	return nfa.ToDFA().Minimize().EliminateDeadStates(), nil
}

type (
	// SymbolTable is used by an EBNF parser during parsing.
	// It keeps track of grammar symbols encountered, their occurrences, and other relevant information.
	// It helps validate aspects of a grammar definition beyond the syntactic structure, such as identifier definitions.
	SymbolTable struct {
		sync.Mutex

		precedences struct {
			list lr.PrecedenceLevels
		}

		terminals struct {
			counter int
			table   symboltable.SymbolTable[grammar.Terminal, *terminalEntry]
		}

		nonTerminals struct {
			counter int
			table   symboltable.SymbolTable[grammar.NonTerminal, *nonTerminalEntry]
		}

		productions struct {
			counter int
			table   symboltable.SymbolTable[*grammar.Production, *productionEntry]
		}

		strings struct {
			counter int
			table   symboltable.SymbolTable[Strings, *stringsEntry]
		}
	}

	// terminalEntry is the table entry for a terminal.
	// term → TOKEN | STRING
	terminalEntry struct {
		index       int
		definitions []terminalDef
		occurrences []*lexer.Position
	}

	// nonTerminalEntry is the table entry for a non-terminal.
	// nonterm → IDENT
	nonTerminalEntry struct {
		index       int
		occurrences []*lexer.Position
	}

	// productionEntry is the table entry for a production rule.
	// rule → lhs "=" rhs | lhs "="
	productionEntry struct {
		index       int
		occurrences []*lexer.Position
	}

	// stringsEntry is the table entry for a list of strings of grammar symbols.
	stringsEntry struct {
		Group grammar.NonTerminal
		Opt   grammar.NonTerminal
		Star  grammar.NonTerminal
		Plus  grammar.NonTerminal
	}
)

// NewSymbolTable creates a new SymbolTable for an EBNF parser.
func NewSymbolTable() *SymbolTable {
	st := new(SymbolTable)

	opts := symboltable.HashOpts{
		InitialCap: 89,
	}

	st.precedences.list = make(lr.PrecedenceLevels, 0)

	st.terminals.table = symboltable.NewQuadraticHashTable[grammar.Terminal, *terminalEntry](
		grammar.HashTerminal,
		grammar.EqTerminal,
		nil,
		opts,
	)

	st.nonTerminals.table = symboltable.NewQuadraticHashTable[grammar.NonTerminal, *nonTerminalEntry](
		grammar.HashNonTerminal,
		grammar.EqNonTerminal,
		nil,
		opts,
	)

	st.productions.table = symboltable.NewQuadraticHashTable[*grammar.Production, *productionEntry](
		grammar.HashProduction,
		grammar.EqProduction,
		nil,
		opts,
	)

	st.strings.table = symboltable.NewQuadraticHashTable[Strings, *stringsEntry](
		hashStrings,
		eqStrings,
		nil,
		opts,
	)

	return st
}

// Reset clears all entries from the symbol table, making it empty.
func (t *SymbolTable) Reset() {
	t.Lock()
	defer t.Unlock()

	t.precedences.list = make(lr.PrecedenceLevels, 0)

	t.terminals.table.DeleteAll()
	t.nonTerminals.table.DeleteAll()
	t.productions.table.DeleteAll()
	t.strings.table.DeleteAll()
}

// Verify is called after parsing is complete and the symbol table is populated.
// It checks for errors beyond the syntactic structure, such as missing or duplicate identifier definitions.
// If any issues are found, it returns an error with a descriptive message.
func (t *SymbolTable) Verify() error {
	t.Lock()
	defer t.Unlock()

	err := &errors.MultiError{
		Format: errors.BulletErrorFormat,
	}

	// Ensure every terminal has one and only one definition.
	for a, e := range t.terminals.table.All() {
		if count := len(e.definitions); count == 0 {
			err = errors.Append(err, fmt.Errorf("no definition for terminal %s", a))
		} else if count > 1 {
			// Aggregate token definitions by their names.
			poss := make([]string, count)
			for i, def := range e.definitions {
				poss[i] = fmt.Sprintf("  %s", def.Pos())
			}

			err = errors.Append(err,
				fmt.Errorf("multiple definitions for terminal %s:\n%s", a, strings.Join(poss, "\n")),
			)
		}
	}

	return err.ErrorOrNil()
}

// Precedences returns the set of precedence levels added to the symbol table.
func (t *SymbolTable) Precedences() lr.PrecedenceLevels {
	t.Lock()
	defer t.Unlock()

	return t.precedences.list
}

// Terminals returns the set of terminal symbols added to the symbol table.
func (t *SymbolTable) Terminals() []grammar.Terminal {
	t.Lock()
	defer t.Unlock()

	var all []grammar.Terminal
	for a := range t.terminals.table.All() {
		all = append(all, a)
	}

	return all
}

// NonTerminals returns the set of non-terminal symbols added to the symbol table.
func (t *SymbolTable) NonTerminals() []grammar.NonTerminal {
	t.Lock()
	defer t.Unlock()

	var all []grammar.NonTerminal
	for A := range t.nonTerminals.table.All() {
		all = append(all, A)
	}

	return all
}

// Productions returns the set of production rules added to the symbol table.
func (t *SymbolTable) Productions() []*grammar.Production {
	t.Lock()
	defer t.Unlock()

	var all []*grammar.Production
	for p := range t.productions.table.All() {
		all = append(all, p)
	}

	return all
}

// AddPrecedence
func (t *SymbolTable) AddPrecedence(p *lr.PrecedenceLevel) {
	t.Lock()
	defer t.Unlock()

	t.precedences.list = append(t.precedences.list, p)
}

// AddStringTokenDef adds a new token definition based on a string value to the symbol table.
func (t *SymbolTable) AddStringTokenDef(token grammar.Terminal, value string, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	e, ok := t.terminals.table.Get(token)

	if !ok {
		t.terminals.counter++
		e = &terminalEntry{
			index:       t.terminals.counter,
			definitions: []terminalDef{},
			occurrences: []*lexer.Position{},
		}

		t.terminals.table.Put(token, e)
	}

	e.definitions = append(e.definitions, &stringTerminalDef{
		value: value,
		pos:   pos,
	})
}

// AddRegexTokenDef adds a new token definition based on a regex value to the symbol table.
func (t *SymbolTable) AddRegexTokenDef(token grammar.Terminal, regex string, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	e, ok := t.terminals.table.Get(token)

	if !ok {
		t.terminals.counter++
		e = &terminalEntry{
			index:       t.terminals.counter,
			definitions: []terminalDef{},
			occurrences: []*lexer.Position{},
		}

		t.terminals.table.Put(token, e)
	}

	e.definitions = append(e.definitions, &regexTerminalDef{
		regex: regex,
		pos:   pos,
	})
}

// AddStringTerminal adds a terminal symbol, defined by its string value, to the symbol table.
// If the terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddStringTerminal(a grammar.Terminal, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	if e, ok := t.terminals.table.Get(a); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.terminals.counter++

	def := &stringTerminalDef{
		value: string(a),
	}

	t.terminals.table.Put(a, &terminalEntry{
		index:       t.terminals.counter,
		definitions: []terminalDef{def},
		occurrences: []*lexer.Position{pos},
	})
}

// AddTokenTerminal adds a terminal symbol, referenced by its token name, to the symbol table.
// If the terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddTokenTerminal(a grammar.Terminal, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	if e, ok := t.terminals.table.Get(a); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.terminals.counter++

	t.terminals.table.Put(a, &terminalEntry{
		index:       t.terminals.counter,
		definitions: []terminalDef{},
		occurrences: []*lexer.Position{pos},
	})
}

// AddNonTerminal adds a non-terminal symbol to the symbol table.
// If the non-terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddNonTerminal(A grammar.NonTerminal, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	if e, ok := t.nonTerminals.table.Get(A); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.nonTerminals.counter++
	t.nonTerminals.table.Put(A, &nonTerminalEntry{
		index:       t.nonTerminals.counter,
		occurrences: []*lexer.Position{pos},
	})
}

// AddProduction adds a production rule to the symbol table.
// If the production rule already exists, the position is added to its occurrences.
func (t *SymbolTable) AddProduction(p *grammar.Production, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	if e, ok := t.productions.table.Get(p); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.productions.counter++
	t.productions.table.Put(p, &productionEntry{
		index:       t.productions.counter,
		occurrences: []*lexer.Position{pos},
	})
}

// GetOpt generates a new non-terminal symbol for an optional (zero or one) occurrence of a list of grammar strings.
// If a name was previously generated for the same strings and purpose, it will be reused.
func (t *SymbolTable) GetOpt(s Strings) grammar.NonTerminal {
	t.Lock()
	defer t.Unlock()

	e, ok := t.strings.table.Get(s)
	if ok {
		return e.Opt
	}

	opt := t.mapStringToNoneTerminal(s, "opt")
	t.strings.table.Put(s, &stringsEntry{
		Opt: opt,
	})

	return opt
}

// GetGroup generates a new non-terminal symbol for grouping a list of grammar strings.
// If a name was previously generated for the same strings and purpose, it will be reused.
func (t *SymbolTable) GetGroup(s Strings) grammar.NonTerminal {
	t.Lock()
	defer t.Unlock()

	e, ok := t.strings.table.Get(s)
	if ok {
		return e.Group
	}

	group := t.mapStringToNoneTerminal(s, "group")
	t.strings.table.Put(s, &stringsEntry{
		Group: group,
	})

	return group
}

// GetStar generates a new non-terminal symbol for zero or more occurrences of a list of grammar strings.
// If a name was previously generated for the same strings and purpose, it will be reused.
func (t *SymbolTable) GetStar(s Strings) grammar.NonTerminal {
	t.Lock()
	defer t.Unlock()

	e, ok := t.strings.table.Get(s)
	if ok {
		return e.Star
	}

	star := t.mapStringToNoneTerminal(s, "star")
	t.strings.table.Put(s, &stringsEntry{
		Star: star,
	})

	return star
}

// GetPlus generates a new non-terminal symbol for one or more occurrences of a list of grammar strings.
// If a name was previously generated for the same strings and purpose, it will be reused.
func (t *SymbolTable) GetPlus(s Strings) grammar.NonTerminal {
	t.Lock()
	defer t.Unlock()

	e, ok := t.strings.table.Get(s)
	if ok {
		return e.Plus
	}

	plus := t.mapStringToNoneTerminal(s, "plus")
	t.strings.table.Put(s, &stringsEntry{
		Plus: plus,
	})

	return plus
}

func (t *SymbolTable) mapStringToNoneTerminal(s Strings, suffix string) grammar.NonTerminal {
	var name string

	if len(s) == 1 && len(s[0]) == 1 {
		switch v := s[0][0].(type) {
		case grammar.NonTerminal:
			name = string(v)

		case grammar.Terminal:
			name = terminalNames[v]
		}
	}

	if name == "" {
		t.strings.counter++
		name = fmt.Sprintf("gen%d_%s", t.strings.counter, suffix)
	} else {
		name = fmt.Sprintf("gen_%s_%s", name, suffix)
	}

	return grammar.NonTerminal(name)
}

var terminalNames = map[grammar.Terminal]string{
	"\t": "tab",
	"\n": "newline",
	" ":  "space",
	"!":  "exclam",
	"\"": "dquot",
	"#":  "hash",
	"$":  "dollar",
	"%":  "percent",
	"&":  "ampersand",
	"'":  "squot",
	"(":  "lparen",
	")":  "rparen",
	"*":  "star",
	"+":  "plus",
	",":  "comma",
	"-":  "dash",
	".":  "dot",
	"/":  "slash",
	":":  "colon",
	";":  "semi",
	"<":  "lt",
	"=":  "equal",
	">":  "gt",
	"?":  "question",
	"@":  "atsign",
	"[":  "lbrack",
	"\\": "backslash",
	"]":  "rbrack",
	"^":  "caret",
	"_":  "underscore",
	"`":  "backtick",
	"{":  "rbrace",
	"|":  "bar",
	"}":  "lbrace",
	"~":  "tilde",
}
