// Generates the parsing table used by the EBNF parser.
// Temporary bootstrap: once Emerge can generate this itself, this program can be removed.
package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/symboltable"

	"github.com/gardenbed/emerge/internal/char"
)

func main() {
	dfa := BuildDFA()
	specs := BuildTokenSpecs()

	code, err := GenerateLexer(dfa, specs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile("lexer.go", code, 0622)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("OK")
}

func BuildDFA() *automata.DFA {
	b := automata.NewDFABuilder().
		SetStart(0).
		SetFinal([]automata.State{
			1,                          // WS
			2,                          // EOL
			3,                          // DEF
			4,                          // SEMI
			5,                          // ALT
			6,                          // LPAREN
			7,                          // RPAREN
			8,                          // LBRACK
			9,                          // RBRACK
			10,                         // LBRACE
			11,                         // RBRACE
			12,                         // LLBRACE
			13,                         // RRBRACE
			14,                         // LANGLE
			15,                         // RANGLE
			17,                         // PREDEF
			22,                         // LASSOC
			27,                         // RASSOC
			31,                         // NOASSOC
			38,                         // GRAMMER
			32, 33, 34, 35, 36, 37, 39, // IDENT
			40,     // TOKEN
			61,     // STRING
			65,     // REGEX
			66, 69, // COMMENT
		})

	// WHITESPACES
	b.AddTransition(0, '\t', '\t', 1).AddTransition(0, ' ', ' ', 1).
		AddTransition(1, '\t', '\t', 1).AddTransition(1, ' ', ' ', 1)

	// NEWLINES
	b.AddTransition(0, '\n', '\n', 2).AddTransition(0, '\r', '\r', 2).
		AddTransition(2, '\n', '\n', 2).AddTransition(2, '\r', '\r', 2)

	// MISC TOKENS
	b.AddTransition(0, '=', '=', 3).
		AddTransition(0, ';', ';', 4).
		AddTransition(0, '|', '|', 5).
		AddTransition(0, '(', '(', 6).
		AddTransition(0, ')', ')', 7).
		AddTransition(0, '[', '[', 8).
		AddTransition(0, ']', ']', 9).
		AddTransition(0, '{', '{', 10).AddTransition(10, '{', '{', 12).
		AddTransition(0, '}', '}', 11).AddTransition(11, '}', '}', 13).
		AddTransition(0, '<', '<', 14).
		AddTransition(0, '>', '>', 15)

	// PREDEFINED TOKENS
	b.AddTransition(0, '$', '$', 16).
		AddTransition(16, 'A', 'Z', 17).
		AddTransition(17, '0', '9', 17).AddTransition(17, 'A', 'Z', 17).AddTransition(17, '_', '_', 17)

	// ASSOCIATIVITY TOKENS
	b.AddTransition(0, '@', '@', 18).
		AddTransition(18, 'l', 'l', 19).AddTransition(19, 'e', 'e', 20).AddTransition(20, 'f', 'f', 21).AddTransition(21, 't', 't', 22).
		AddTransition(18, 'r', 'r', 23).AddTransition(23, 'i', 'i', 24).AddTransition(24, 'g', 'g', 25).AddTransition(25, 'h', 'h', 26).AddTransition(26, 't', 't', 27).
		AddTransition(18, 'n', 'n', 28).AddTransition(28, 'o', 'o', 29).AddTransition(29, 'n', 'n', 30).AddTransition(30, 'e', 'e', 31)

	// GRAMMER & IDENT
	b.AddTransition(0, 'g', 'g', 32).AddTransition(32, 'r', 'r', 33).AddTransition(33, 'a', 'a', 34).AddTransition(34, 'm', 'm', 35).AddTransition(35, 'm', 'm', 36).AddTransition(36, 'a', 'a', 37).AddTransition(37, 'r', 'r', 38).
		AddTransition(0, 'a', 'f', 39).AddTransition(0, 'h', 'z', 39).
		AddTransition(32, '0', '9', 39).AddTransition(32, '_', '_', 39).AddTransition(32, 'a', 'q', 39).AddTransition(32, 's', 'z', 39).
		AddTransition(33, '0', '9', 39).AddTransition(33, '_', '_', 39).AddTransition(33, 'b', 'z', 39).
		AddTransition(34, '0', '9', 39).AddTransition(34, '_', '_', 39).AddTransition(34, 'a', 'l', 39).AddTransition(34, 'n', 'z', 39).
		AddTransition(35, '0', '9', 39).AddTransition(35, '_', '_', 39).AddTransition(35, 'a', 'l', 39).AddTransition(35, 'n', 'z', 39).
		AddTransition(36, '0', '9', 39).AddTransition(36, '_', '_', 39).AddTransition(36, 'b', 'z', 39).
		AddTransition(37, '0', '9', 39).AddTransition(37, '_', '_', 39).AddTransition(37, 'a', 'q', 39).AddTransition(37, 's', 'z', 39).
		AddTransition(38, '0', '9', 39).AddTransition(38, '_', '_', 39).AddTransition(38, 'a', 'z', 39).
		AddTransition(39, '0', '9', 39).AddTransition(39, '_', '_', 39).AddTransition(39, 'a', 'z', 39)

	// TOKEN
	b.AddTransition(0, 'A', 'Z', 40).
		AddTransition(40, '0', '9', 40).AddTransition(40, 'A', 'Z', 40).AddTransition(40, '_', '_', 40)

	// STRING
	b.AddTransition(0, '"', '"', 41).
		AddTransition(41, '\\', '\\', 42).
		AddTransition(41, '"', '"', 61)

	// Escapes: \\ \' \" \t \n \r
	b.AddTransition(42, '\\', '\\', 43).
		AddTransition(42, '\'', '\'', 43).
		AddTransition(42, '"', '"', 43).
		AddTransition(42, 't', 't', 43).
		AddTransition(42, 'n', 'n', 43).
		AddTransition(42, 'r', 'r', 43).
		AddTransition(43, '\\', '\\', 42).
		AddTransition(43, '"', '"', 61)

	// ASCII Escapes: \xhh
	b.AddTransition(42, 'x', 'x', 44).
		AddTransition(44, '0', '9', 45).AddTransition(44, 'A', 'F', 45).AddTransition(44, 'a', 'f', 45).
		AddTransition(45, '0', '9', 46).AddTransition(45, 'A', 'F', 46).AddTransition(45, 'a', 'f', 46).
		AddTransition(46, '\\', '\\', 42).
		AddTransition(46, '"', '"', 61)

	// Unicode Escapes: \uhhhh
	b.AddTransition(42, 'u', 'u', 47).
		AddTransition(47, '0', '9', 48).AddTransition(47, 'A', 'F', 48).AddTransition(47, 'a', 'f', 48).
		AddTransition(48, '0', '9', 49).AddTransition(48, 'A', 'F', 49).AddTransition(48, 'a', 'f', 49).
		AddTransition(49, '0', '9', 50).AddTransition(49, 'A', 'F', 50).AddTransition(49, 'a', 'f', 50).
		AddTransition(50, '0', '9', 51).AddTransition(50, 'A', 'F', 51).AddTransition(50, 'a', 'f', 51).
		AddTransition(51, '\\', '\\', 42).
		AddTransition(51, '"', '"', 61)

	// Unicode Escapes: \Uhhhhhhhh
	b.AddTransition(42, 'U', 'U', 52).
		AddTransition(52, '0', '9', 53).AddTransition(52, 'A', 'F', 53).AddTransition(52, 'a', 'f', 53).
		AddTransition(53, '0', '9', 54).AddTransition(53, 'A', 'F', 54).AddTransition(53, 'a', 'f', 54).
		AddTransition(54, '0', '9', 55).AddTransition(54, 'A', 'F', 55).AddTransition(54, 'a', 'f', 55).
		AddTransition(55, '0', '9', 56).AddTransition(55, 'A', 'F', 56).AddTransition(55, 'a', 'f', 56).
		AddTransition(56, '0', '9', 57).AddTransition(56, 'A', 'F', 57).AddTransition(56, 'a', 'f', 57).
		AddTransition(57, '0', '9', 58).AddTransition(57, 'A', 'F', 58).AddTransition(57, 'a', 'f', 58).
		AddTransition(58, '0', '9', 59).AddTransition(58, 'A', 'F', 59).AddTransition(58, 'a', 'f', 59).
		AddTransition(59, '0', '9', 60).AddTransition(59, 'A', 'F', 60).AddTransition(59, 'a', 'f', 60).
		AddTransition(60, '\\', '\\', 42).
		AddTransition(60, '"', '"', 61)

	// All Unicode characters except \ "
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'\\', '\\'}, {'"', '"'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])

		b.AddTransition(41, lo, hi, 41)
		b.AddTransition(43, lo, hi, 41)
		b.AddTransition(46, lo, hi, 41)
		b.AddTransition(51, lo, hi, 41)
		b.AddTransition(60, lo, hi, 41)
	}

	// REGEX
	b.AddTransition(0, '/', '/', 62).
		AddTransition(62, '\\', '\\', 63).
		AddTransition(64, '\\', '\\', 63).
		AddTransition(64, '/', '/', 65)

	// All Unicode characters except / \ *
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'/', '/'}, {'\\', '\\'}, {'*', '*'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(62, lo, hi, 64)
	}

	// All Unicode characters
	for _, r := range char.Classes["UNICODE"] {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(63, lo, hi, 64)
	}

	// All Unicode characters except / \
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'/', '/'}, {'\\', '\\'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(64, lo, hi, 64)
	}

	// SINGLE-LINE COMMENT
	b.AddTransition(62, '/', '/', 66)

	// All Unicode characters except \n \v \f \r
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'\n', '\r'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(66, lo, hi, 66)
	}

	// MULTI-LINE COMMENT
	b.AddTransition(62, '*', '*', 67).
		AddTransition(67, '*', '*', 68).
		AddTransition(68, '*', '*', 68).
		AddTransition(68, '/', '/', 69)

	// All Unicode characters except *
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'*', '*'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(67, lo, hi, 67)
	}

	// All Unicode characters except * /
	for _, r := range char.Classes["UNICODE"].Exclude(char.RangeList{{'*', '*'}, {'/', '/'}}) {
		lo, hi := automata.Symbol(r[0]), automata.Symbol(r[1])
		b.AddTransition(68, lo, hi, 67)
	}

	return b.Build()
}

func BuildTokenSpecs() tokenSpecs {
	specs := symboltable.NewRedBlack[automata.States, tokenSpec](automata.CmpStates, nil)

	specs.Put(automata.NewStates(1), tokenSpec{
		TerminalName: "WS",
		LexemeValue:  stringPtr(""),
	})

	specs.Put(automata.NewStates(2), tokenSpec{
		TerminalName: "EOL",
		LexemeValue:  stringPtr(""),
	})

	specs.Put(automata.NewStates(3), tokenSpec{
		TerminalName: "DEF",
		LexemeValue:  stringPtr("="),
	})

	specs.Put(automata.NewStates(4), tokenSpec{
		TerminalName: "SEMI",
		LexemeValue:  stringPtr(";"),
	})

	specs.Put(automata.NewStates(5), tokenSpec{
		TerminalName: "ALT",
		LexemeValue:  stringPtr("|"),
	})

	specs.Put(automata.NewStates(6), tokenSpec{
		TerminalName: "LPAREN",
		LexemeValue:  stringPtr("("),
	})

	specs.Put(automata.NewStates(7), tokenSpec{
		TerminalName: "RPAREN",
		LexemeValue:  stringPtr(")"),
	})

	specs.Put(automata.NewStates(8), tokenSpec{
		TerminalName: "LBRACK",
		LexemeValue:  stringPtr("["),
	})

	specs.Put(automata.NewStates(9), tokenSpec{
		TerminalName: "RBRACK",
		LexemeValue:  stringPtr("]"),
	})

	specs.Put(automata.NewStates(10), tokenSpec{
		TerminalName: "LBRACE",
		LexemeValue:  stringPtr("{"),
	})

	specs.Put(automata.NewStates(11), tokenSpec{
		TerminalName: "RBRACE",
		LexemeValue:  stringPtr("}"),
	})

	specs.Put(automata.NewStates(12), tokenSpec{
		TerminalName: "LLBRACE",
		LexemeValue:  stringPtr("{{"),
	})

	specs.Put(automata.NewStates(13), tokenSpec{
		TerminalName: "RRBRACE",
		LexemeValue:  stringPtr("}}"),
	})

	specs.Put(automata.NewStates(14), tokenSpec{
		TerminalName: "LANGLE",
		LexemeValue:  stringPtr("<"),
	})

	specs.Put(automata.NewStates(15), tokenSpec{
		TerminalName: "RANGLE",
		LexemeValue:  stringPtr(">"),
	})

	specs.Put(automata.NewStates(17), tokenSpec{
		TerminalName: "PREDEF",
	})

	specs.Put(automata.NewStates(22), tokenSpec{
		TerminalName: "LASSOC",
		LexemeValue:  stringPtr("@left"),
	})

	specs.Put(automata.NewStates(27), tokenSpec{
		TerminalName: "RASSOC",
		LexemeValue:  stringPtr("@right"),
	})

	specs.Put(automata.NewStates(31), tokenSpec{
		TerminalName: "NOASSOC",
		LexemeValue:  stringPtr("@none"),
	})

	specs.Put(automata.NewStates(38), tokenSpec{
		TerminalName: "GRAMMER",
		LexemeValue:  stringPtr("grammar"),
	})

	specs.Put(automata.NewStates(32, 33, 34, 35, 36, 37, 39), tokenSpec{
		TerminalName: "IDENT",
	})

	specs.Put(automata.NewStates(40), tokenSpec{
		TerminalName: "TOKEN",
	})

	specs.Put(automata.NewStates(61), tokenSpec{
		TerminalName: "STRING",
		TrimLexeme:   true,
	})

	specs.Put(automata.NewStates(65), tokenSpec{
		TerminalName: "REGEX",
		TrimLexeme:   true,
	})

	specs.Put(automata.NewStates(66), tokenSpec{
		TerminalName: "COMMENT",
		LexemeValue:  stringPtr(""),
	})

	specs.Put(automata.NewStates(69), tokenSpec{
		TerminalName: "COMMENT",
		LexemeValue:  stringPtr(""),
	})

	return specs
}

type (
	tokenSpecs symboltable.SymbolTable[automata.States, tokenSpec]

	tokenSpec struct {
		TerminalName string
		LexemeValue  *string
		TrimLexeme   bool
	}
)

func stringPtr(s string) *string {
	return &s
}

func GenerateLexer(dfa *automata.DFA, specs tokenSpecs) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(`//go:generate go run ./generate

// Package lexer implements a lexical analyzer for the EBNF language.
package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/lexer/input"
)

const (
	errorState = -1
	bufferSize = 4096
)

const (
	ERR     = grammar.Terminal("ERR")     // ERR is the error token.
	WS      = grammar.Terminal("WS")      // WS is the token for whitespace characters.
	EOL     = grammar.Terminal("EOL")     // WS is the token for newline characters.
	DEF     = grammar.Terminal("=")       // DEF is the token for "=".
	SEMI    = grammar.Terminal(";")       // SEMI is the token for ";".
	ALT     = grammar.Terminal("|")       // ALT is the token for "|".
	LPAREN  = grammar.Terminal("(")       // LPAREN is the token for "(".
	RPAREN  = grammar.Terminal(")")       // RPAREN is the token for ")".
	LBRACK  = grammar.Terminal("[")       // LBRACK is the token for "[".
	RBRACK  = grammar.Terminal("]")       // RBRACK is the token for "]".
	LBRACE  = grammar.Terminal("{")       // LBRACE is the token for "{".
	RBRACE  = grammar.Terminal("}")       // RBRACE is the token for "}".
	LLBRACE = grammar.Terminal("{{")      // LLBRACE is the token for "{{".
	RRBRACE = grammar.Terminal("}}")      // RRBRACE is the token for "}}".
	LANGLE  = grammar.Terminal("<")       // LANGLE  is the token for "<".
	RANGLE  = grammar.Terminal(">")       // RANGLE  is the token for ">".
	PREDEF  = grammar.Terminal("PREDEF")  // PREDEF is the token for /\$[A-Z][0-9A-Z_]*/.
	LASSOC  = grammar.Terminal("@left")   // LASSOC  is the token for "@left".
	RASSOC  = grammar.Terminal("@right")  // RASSOC  is the token for "@right".
	NOASSOC = grammar.Terminal("@none")   // NOASSOC is the token for "@none".
	GRAMMER = grammar.Terminal("grammar") // GRAMMER is the token for "grammar".
	IDENT   = grammar.Terminal("IDENT")   // IDENT is the token for /[a-z][0-9a-z_]*/.
	TOKEN   = grammar.Terminal("TOKEN")   // TOKEN is the token for /[A-Z][0-9A-Z_]*/.
	STRING  = grammar.Terminal("STRING")  // STRING is the token for /"([^\\"]\|\\[\\"'tnr]\|\\x[0-9A-Fa-f]{2}\|\\u[0-9A-Fa-f]{4}\|\\U[0-9A-Fa-f]{8})*"/.
	REGEX   = grammar.Terminal("REGEX")   // REGEX is the token for /\/([^\/\\*]\|\\.)([^\/\\]\|\\.)*\//.
	COMMENT = grammar.Terminal("COMMENT") // COMMENT is the token for single-line and multi-line comments.
)

// inputBuffer is an interface for the input.Input struct.
type inputBuffer interface {
	Next() (rune, error)
	Retract()
	Lexeme() (string, lexer.Position)
	Skip() lexer.Position
}

// Lexer is a lexical analyzer for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
type Lexer struct {
	in inputBuffer
}

// New creates a new lexical analyzer for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
func New(filename string, src io.Reader) (*Lexer, error) {
	in, err := input.New(filename, src, bufferSize)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		in: in,
	}, nil
}

// NextToken scans the input stream until it recognizes a valid token, which it then returns.
// If the end of the input is reached, it returns an io.EOF error.
func (l *Lexer) NextToken() (lexer.Token, error) {
	for curr, next := 0, 0; ; curr = next {
		// Read the next character from the input stream.
		r, err := l.in.Next()
		if err != nil {
			return lexer.Token{}, err
		}

		// Keep running the DFA through the input symbols.
		next = advanceDFA(curr, r)

		if next == errorState {
			// Retract one character, as the last read character did not belong to the current token.
			l.in.Retract()

			// Evaluate the final state of the DFA.
			token := l.evalDFA(curr)

			switch token.Terminal {
			case ERR:
				return lexer.Token{}, errors.New(token.Lexeme)
			case WS, EOL, COMMENT:
				// Skip whitespaces, newlines, and comments.
				return l.NextToken()
			default:
				return token, nil
			}
		}
	}
}

`)

	b.Write(generateEvalDFA(specs))
	b.WriteString("\n")
	b.Write(generateAdvanceDFA(dfa))

	return b.Bytes(), nil
}

func generateEvalDFA(specs tokenSpecs) []byte {
	var b bytes.Buffer

	b.WriteString("// evalDFA examines the final state of a deterministic finite automaton (DFA) after it has stopped processing input.\n")
	b.WriteString("// Based on the last encountered state, it returns the corresponding token and advances the input buffer reader.\n")
	b.WriteString("// If the final state is invalid, it returns an ERR token with the Lexeme set to the error message.\n")
	b.WriteString("func (l *Lexer) evalDFA(state int) lexer.Token {\n")
	b.WriteString("	switch state {\n")

	for states, spec := range specs.All() {
		fmt.Fprintf(&b, "	case ")

		for s := range states.All() {
			fmt.Fprintf(&b, "%d, ", s)
		}

		// Remove last comma and space
		b.Truncate(b.Len() - 2)

		fmt.Fprintf(&b, ":\n")

		if spec.LexemeValue != nil {
			fmt.Fprintf(&b, "		pos := l.in.Skip()\n")
			fmt.Fprintf(&b, "		return lexer.Token{Terminal: %s, Lexeme: %q, Pos: pos}\n", spec.TerminalName, *spec.LexemeValue)
		} else if !spec.TrimLexeme {
			fmt.Fprintf(&b, "		lexeme, pos := l.in.Lexeme()\n")
			fmt.Fprintf(&b, "		return lexer.Token{Terminal: %s, Lexeme: lexeme, Pos: pos}\n", spec.TerminalName)
		} else {
			fmt.Fprintf(&b, "		lexeme, pos := l.in.Lexeme()\n")
			fmt.Fprintf(&b, "		lexeme = lexeme[1 : len(lexeme)-1]\n")
			fmt.Fprintf(&b, "		return lexer.Token{Terminal: %s, Lexeme: lexeme, Pos: pos}\n", spec.TerminalName)
		}

		fmt.Fprintln(&b)
	}

	b.WriteString("	}\n")
	b.WriteString("\n")
	b.WriteString("	// ERR\n")
	b.WriteString("	val, pos := l.in.Lexeme()\n")
	b.WriteString("	return lexer.Token{\n")
	b.WriteString("		Terminal: ERR,\n")
	b.WriteString("		Lexeme:   fmt.Sprintf(\"lexical error at %s:%s\", pos, val),\n")
	b.WriteString("		Pos:      pos,\n")
	b.WriteString("	}\n")
	b.WriteString("}\n")

	return b.Bytes()
}

func generateAdvanceDFA(dfa *automata.DFA) []byte {
	var b bytes.Buffer

	b.WriteString("// advanceDFA determines the next state of a deterministic finite automaton (DFA)\n")
	b.WriteString("// given the current state and an input symbol.\n")
	b.WriteString("// It functions as a coded lookup table.\n")
	b.WriteString("func advanceDFA(state int, r rune) int {\n")
	b.WriteString("	switch state {\n")

	for s, seq := range dfa.Transitions() {
		fmt.Fprintf(&b, "	case %d:\n", s)
		fmt.Fprintf(&b, "		switch {\n")

		for ranges, next := range seq {
			fmt.Fprintf(&b, "		case ")

			for _, r := range ranges {
				if r.Lo == r.Hi {
					fmt.Fprintf(&b, "r == %q, ", r.Lo)
				} else {
					fmt.Fprintf(&b, "%q <= r && r <= %q, ", r.Lo, r.Hi)
				}
			}

			// Remove last comma and space
			b.Truncate(b.Len() - 2)

			fmt.Fprintf(&b, ":\n")
			fmt.Fprintf(&b, "			return %d\n", next)
		}

		fmt.Fprintf(&b, "		}\n\n")
	}

	// Remove last newline
	b.Truncate(b.Len() - 1)

	b.WriteString("	}\n")
	b.WriteString("\n")
	b.WriteString("	return errorState\n")
	b.WriteString("}\n")

	return b.Bytes()
}
