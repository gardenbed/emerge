package spec

import (
	"fmt"
	"strings"
	"sync"

	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/sort"
	"github.com/moorara/algo/symboltable"
)

// start is always the start symbol of an EBNF grammar by convention.
const start = grammar.NonTerminal("start")

// TerminalDefKind indicates whether a terminal definition is based on a string or a regex.
type TerminalDefKind int

const (
	StringDef TerminalDefKind = iota
	RegexDef
)

// TerminalDef represents a terminal symbol along with a deterministic finite automaton (DFA) for recognizing it.
type TerminalDef struct {
	grammar.Terminal
	Kind  TerminalDefKind
	Value string
	Pos   *lexer.Position
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
		definitions []*TerminalDef
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

	errs := &errors.MultiError{
		Format: errors.BulletErrorFormat,
	}

	if err := t.ensureSingleDefs(); err != nil {
		errs = errors.Append(errs, err)
	}

	if err := t.ensureDistinctDefs(); err != nil {
		errs = errors.Append(errs, err)
	}

	if err := t.ensureStartSymbol(); err != nil {
		errs = errors.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

// ensureSingleDefs verifies that each terminal has exactly one definition.
// It reports an error if a terminal is missing a definition or has multiple definitions.
func (t *SymbolTable) ensureSingleDefs() error {
	var errs error

	for a, e := range t.terminals.table.All() {
		if count := len(e.definitions); count == 0 {
			errs = errors.Append(errs, fmt.Errorf("no definition for terminal %s", a))
		} else if count > 1 {
			poses := generic.Transform(e.definitions, func(def *TerminalDef) string {
				return fmt.Sprintf("  %s", def.Pos)
			})

			errs = errors.Append(errs,
				fmt.Errorf("multiple definitions for terminal %s:\n%s", a, strings.Join(poses, "\n")),
			)
		}
	}

	return errs
}

// ensureDistinctDefs verifies that terminal definitions use distinct values.
// It reports an error when more than one terminal is defined with the same value.
func (t *SymbolTable) ensureDistinctDefs() error {
	var errs error

	reverse := make(map[string][]*TerminalDef)
	for _, e := range t.terminals.table.All() {
		if len(e.definitions) == 1 {
			def := e.definitions[0]
			reverse[def.Value] = append(reverse[def.Value], def)
		}
	}

	for val, defs := range reverse {
		if len(defs) > 1 {
			poses := generic.Transform(defs, func(def *TerminalDef) string {
				return fmt.Sprintf("  %s: %s", def.Pos, def.Terminal)
			})

			errs = errors.Append(errs,
				fmt.Errorf("multiple definitions with the same value: %q\n%s", val, strings.Join(poses, "\n")),
			)
		}
	}

	return errs
}

// ensureStartSymbol ensures a production rule exists with the start symbol as the head non-terminal.
func (t *SymbolTable) ensureStartSymbol() error {
	hasStart := t.productions.table.AnyMatch(func(p *grammar.Production, _ *productionEntry) bool {
		return p.Head.Equal(start)
	})

	if !hasStart {
		return fmt.Errorf("missing production rule with the start symbol: %s", start)
	}

	return nil
}

// Precedences returns the set of precedence levels added to the symbol table.
func (t *SymbolTable) Precedences() lr.PrecedenceLevels {
	t.Lock()
	defer t.Unlock()

	return t.precedences.list
}

// Definitions constructs and returns an ordered list of definitions,
// representing deterministic finite automata (DFAs) for all terminal symbols in the symbol table.
func (t *SymbolTable) Definitions() []*TerminalDef {
	t.Lock()
	defer t.Unlock()

	defs := make([]*TerminalDef, 0, t.terminals.table.Size())
	for _, e := range t.terminals.table.All() {
		if len(e.definitions) == 1 {
			defs = append(defs, e.definitions[0])
		}
	}

	// Sort terminals, placing string-based terminals before regex-based ones and shorter terminals before longer ones.
	sort.Quick(defs, func(lhs, rhs *TerminalDef) int {
		if lhs.Kind == StringDef && rhs.Kind == RegexDef {
			return -1
		} else if lhs.Kind == RegexDef && rhs.Kind == StringDef {
			return 1
		}

		if len(lhs.Terminal) < len(rhs.Terminal) {
			return -1
		} else if len(lhs.Terminal) > len(rhs.Terminal) {
			return 1
		}

		return grammar.CmpTerminal(lhs.Terminal, rhs.Terminal)
	})

	return defs
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

// AddPrecedence adds a new precedence level to the symbol table.
func (t *SymbolTable) AddPrecedence(p *lr.PrecedenceLevel) {
	t.Lock()
	defer t.Unlock()

	t.precedences.list = append(t.precedences.list, p)
}

// AddStringTokenDef adds a token definition with a string value to the symbol table.
func (t *SymbolTable) AddStringTokenDef(token grammar.Terminal, value string, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	e, ok := t.terminals.table.Get(token)

	if !ok {
		t.terminals.counter++
		e = &terminalEntry{
			index:       t.terminals.counter,
			definitions: []*TerminalDef{},
			occurrences: []*lexer.Position{},
		}

		t.terminals.table.Put(token, e)
	}

	e.definitions = append(e.definitions, &TerminalDef{
		Terminal: token,
		Kind:     StringDef,
		Value:    value,
		Pos:      pos,
	})
}

// AddRegexTokenDef adds a token definition with a regex value to the symbol table.
func (t *SymbolTable) AddRegexTokenDef(token grammar.Terminal, regex string, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	e, ok := t.terminals.table.Get(token)

	if !ok {
		t.terminals.counter++
		e = &terminalEntry{
			index:       t.terminals.counter,
			definitions: []*TerminalDef{},
			occurrences: []*lexer.Position{},
		}

		t.terminals.table.Put(token, e)
	}

	e.definitions = append(e.definitions, &TerminalDef{
		Terminal: token,
		Kind:     RegexDef,
		Value:    regex,
		Pos:      pos,
	})
}

// AddStringTerminal records an occurrence of a terminal symbol referenced by its string value in the symbol table.
//
// Because the terminal is referenced by its string value,
// if the terminal does not yet have a definition in the table, a new terminal definition is created.
// If the terminal already have a definition in the table,
// the position of the new occurrence is appended to its occurrences.
func (t *SymbolTable) AddStringTerminal(a grammar.Terminal, pos *lexer.Position) {
	t.Lock()
	defer t.Unlock()

	if e, ok := t.terminals.table.Get(a); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.terminals.counter++
	def := &TerminalDef{
		Terminal: a,
		Kind:     StringDef,
		Value:    string(a),
		// Pos will not be set because the same terminal can be referenced by its string in multiple places,
		// so there is no single position to associate with the definition.
	}

	t.terminals.table.Put(a, &terminalEntry{
		index:       t.terminals.counter,
		definitions: []*TerminalDef{def},
		occurrences: []*lexer.Position{pos},
	})
}

// AddTokenTerminal records an occurrence of a terminal symbol referenced by its token name in the symbol table.
//
// If the terminal already have a definition in the table,
// the position of the new occurrence is appended to its occurrences.
// If the terminal does not yet have a definition in the table,
// a new entry is created with an empty list of definitions and the occurrence position;
// the terminal's actual definition will be added later when it is encountered during parsing.
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
		definitions: []*TerminalDef{}, // The actual definition will be added later, when it is encountered during parsing.
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
