package generator

import (
	"fmt"
	"strings"

	"github.com/moorara/algo/errors"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/symboltable"
)

type (
	// SymbolTable is used by an EBNF parser during parsing.
	// It keeps track of grammar symbols encountered, their occurrences, and other relevant information.
	// It helps validate aspects of a grammar definition beyond the syntactic structure, such as identifier definitions.
	SymbolTable struct {
		precedences struct {
			list lr.PrecedenceLevels
		}

		tokenDefs struct {
			strings []*tokenDefEntry
			regexes []*tokenDefEntry
		}

		terminals struct {
			counter int
			strings symboltable.SymbolTable[grammar.Terminal, *terminalEntry]
			tokens  symboltable.SymbolTable[grammar.Terminal, *terminalEntry]
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

	// tokenDefEntry is the table entry for a token definition.
	// token → TOKEN "=" (STRING | REGEX | PREDEF)
	tokenDefEntry struct {
		token      grammar.Terminal
		value      string
		occurrence *lexer.Position
	}

	// terminalEntry is the table entry for a terminal.
	// term → TOKEN | STRING
	terminalEntry struct {
		index       int
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

	st.tokenDefs.strings = make([]*tokenDefEntry, 0)
	st.tokenDefs.regexes = make([]*tokenDefEntry, 0)

	st.terminals.strings = symboltable.NewQuadraticHashTable[grammar.Terminal, *terminalEntry](
		grammar.HashTerminal,
		grammar.EqTerminal,
		nil,
		opts,
	)

	st.terminals.tokens = symboltable.NewQuadraticHashTable[grammar.Terminal, *terminalEntry](
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

// Verify is called after parsing is complete and the symbol table is populated.
// It checks for errors beyond the syntactic structure, such as missing or duplicate identifier definitions.
// If any issues are found, it returns an error with a descriptive message.
func (t *SymbolTable) Verify() error {
	err := &errors.MultiError{
		Format: errors.BulletErrorFormat,
	}

	// Check if there is a definition for every terminal referenced by a token name.
	for token := range t.terminals.tokens.All() {
		if !generic.AnyMatch(t.tokenDefs.strings, func(e *tokenDefEntry) bool {
			return e.token.Equal(token)
		}) && !generic.AnyMatch(t.tokenDefs.regexes, func(e *tokenDefEntry) bool {
			return e.token.Equal(token)
		}) {
			err = errors.Append(err, fmt.Errorf("no definition for terminal %s", token))
		}
	}

	// Aggregate token definitions by their names.
	agg := map[grammar.Terminal][]string{}
	for _, e := range append(t.tokenDefs.strings, t.tokenDefs.regexes...) {
		agg[e.token] = append(agg[e.token], fmt.Sprintf("  %s", e.occurrence))
	}

	// Check if there is more than one definition for any tokens.
	for token, occurrs := range agg {
		if len(occurrs) > 1 {
			err = errors.Append(err,
				fmt.Errorf("multiple definitions for terminal %s:\n%s", token, strings.Join(occurrs, "\n")),
			)
		}
	}

	return err.ErrorOrNil()
}

// Precedences returns the set of precedence levels added to the symbol table.
func (t *SymbolTable) Precedences() lr.PrecedenceLevels {
	return t.precedences.list
}

// Terminals returns the set of terminal symbols added to the symbol table.
func (t *SymbolTable) Terminals() []grammar.Terminal {
	var all []grammar.Terminal

	for a := range t.terminals.strings.All() {
		all = append(all, a)
	}

	for a := range t.terminals.tokens.All() {
		all = append(all, a)
	}

	return all
}

// NonTerminals returns the set of non-terminal symbols added to the symbol table.
func (t *SymbolTable) NonTerminals() []grammar.NonTerminal {
	var all []grammar.NonTerminal

	for A := range t.nonTerminals.table.All() {
		all = append(all, A)
	}

	return all
}

// Productions returns the set of production rules added to the symbol table.
func (t *SymbolTable) Productions() []*grammar.Production {
	var all []*grammar.Production

	for p := range t.productions.table.All() {
		all = append(all, p)
	}

	return all
}

// AddPrecedence
func (t *SymbolTable) AddPrecedence(p *lr.PrecedenceLevel) {
	t.precedences.list = append(t.precedences.list, p)
}

// AddStringTokenDef adds a new token definition based on a string value to the symbol table.
func (t *SymbolTable) AddStringTokenDef(token grammar.Terminal, value string, pos *lexer.Position) {
	t.tokenDefs.strings = append(t.tokenDefs.strings, &tokenDefEntry{
		token:      token,
		value:      value,
		occurrence: pos,
	})
}

// AddRegexTokenDef adds a new token definition based on a regex value to the symbol table.
func (t *SymbolTable) AddRegexTokenDef(token grammar.Terminal, value string, pos *lexer.Position) {
	t.tokenDefs.regexes = append(t.tokenDefs.regexes, &tokenDefEntry{
		token:      token,
		value:      value,
		occurrence: pos,
	})
}

// AddStringTerminal adds a terminal symbol, defined by its string value, to the symbol table.
// If the terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddStringTerminal(a grammar.Terminal, pos *lexer.Position) {
	if e, ok := t.terminals.strings.Get(a); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.terminals.counter++
	t.terminals.strings.Put(a, &terminalEntry{
		index:       t.terminals.counter,
		occurrences: []*lexer.Position{pos},
	})
}

// AddTokenTerminal adds a terminal symbol, referenced by its token name, to the symbol table.
// If the terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddTokenTerminal(a grammar.Terminal, pos *lexer.Position) {
	if e, ok := t.terminals.tokens.Get(a); ok {
		e.occurrences = append(e.occurrences, pos)
		return
	}

	t.terminals.counter++
	t.terminals.tokens.Put(a, &terminalEntry{
		index:       t.terminals.counter,
		occurrences: []*lexer.Position{pos},
	})
}

// AddNonTerminal adds a non-terminal symbol to the symbol table.
// If the non-terminal symbol already exists, the position is added to its occurrences.
func (t *SymbolTable) AddNonTerminal(A grammar.NonTerminal, pos *lexer.Position) {
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
