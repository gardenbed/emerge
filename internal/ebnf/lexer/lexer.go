// Package lexer implements a lexical analyzer for an extension of EBNF language.
package lexer

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/lexer/input"
)

const (
	errorState = -1
	bufferSize = 4096
)

const (
	ERR     = grammar.Terminal("ERR")     // ERR is the terminal for an error case.
	WS      = grammar.Terminal("WS")      // WS is the terminal for whitespace characters.
	DEF     = grammar.Terminal("DEF")     // DEF is the terminal for "=".
	ALT     = grammar.Terminal("ALT")     // ALT is the terminal for "|".
	LPAREN  = grammar.Terminal("LPAREN")  // LPAREN is the terminal for "(".
	RPAREN  = grammar.Terminal("RPAREN")  // RPAREN is the terminal for ")".
	LBRACK  = grammar.Terminal("LBRACK")  // LBRACK is the terminal for "[".
	RBRACK  = grammar.Terminal("RBRACK")  // RBRACK is the terminal for "]".
	LBRACE  = grammar.Terminal("LBRACE")  // LBRACE is the terminal for "{".
	RBRACE  = grammar.Terminal("RBRACE")  // RBRACE is the terminal for "}".
	LLBRACE = grammar.Terminal("LLBRACE") // LLBRACE is the terminal for "{{".
	RRBRACE = grammar.Terminal("RRBRACE") // RRBRACE is the terminal for "}}".
	GRAMMER = grammar.Terminal("GRAMMER") // GRAMMER is the terminal for "grammar".
	LASSOC  = grammar.Terminal("LASSOC")  // LASSOC  is the terminal for "@left".
	RASSOC  = grammar.Terminal("RASSOC")  // RASSOC  is the terminal for "@right".
	NOASSOC = grammar.Terminal("NOASSOC") // NOASSOC is the terminal for "@none".
	IDENT   = grammar.Terminal("IDENT")   // IDENT is the terminal for /[a-z][0-9a-z_]*/.
	TOKEN   = grammar.Terminal("TOKEN")   // TOKEN is the terminal for /[A-Z][0-9A-Z_]*/.
	STRING  = grammar.Terminal("STRING")  // STRING is the terminal for /"([\x21\x23-\x5B\x5D-\x7E]|\\[\x21-\x7E])+"/.
	REGEX   = grammar.Terminal("REGEX")   // REGEX is the terminal for /\/([\x20-\x2E\x30-\x5B\x5D-\x7E]|\\[\x20-\x7E])*\//.
	PREDEF  = grammar.Terminal("PREDEF")  // PREDEF is the terminal for /\$[A-Z][0-9A-Z_]*/.
	COMMENT = grammar.Terminal("COMMENT") // COMMENT is the terminal for single-line and multi-line comments.
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
			case WS, COMMENT:
				// Skip whitespace and comments.
				return l.NextToken()
			default:
				return token, nil
			}
		}
	}
}

// evalDFA examines the final state of a deterministic finite automaton (DFA) after it has stopped processing input.
// Based on the last encountered state, it returns the corresponding token and advances the input buffer reader.
// If the final state is invalid, it returns an ERR token with the Lexeme set to an error message.
func (l *Lexer) evalDFA(state int) lexer.Token {
	switch state {
	// Whitespace
	case 1:
		pos := l.in.Skip()
		return lexer.Token{Terminal: WS, Lexeme: "", Pos: pos}

	// DEF
	case 2:
		pos := l.in.Skip()
		return lexer.Token{Terminal: DEF, Lexeme: "=", Pos: pos}

	// ALT
	case 3:
		pos := l.in.Skip()
		return lexer.Token{Terminal: ALT, Lexeme: "|", Pos: pos}

	// LPAREN
	case 4:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LPAREN, Lexeme: "(", Pos: pos}

	// RPAREN
	case 5:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RPAREN, Lexeme: ")", Pos: pos}

	// LBRACK
	case 6:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LBRACK, Lexeme: "[", Pos: pos}

	// RBRACK
	case 7:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RBRACK, Lexeme: "]", Pos: pos}

	// LBRACE
	case 8:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LBRACE, Lexeme: "{", Pos: pos}

	// RBRACE
	case 9:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RBRACE, Lexeme: "}", Pos: pos}

	// LLBRACE
	case 10:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LLBRACE, Lexeme: "{{", Pos: pos}

	// RRBRACE
	case 11:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RRBRACE, Lexeme: "}}", Pos: pos}

	// LASSOC
	case 16:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LASSOC, Lexeme: "@left", Pos: pos}

	// RASSOC
	case 21:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RASSOC, Lexeme: "@right", Pos: pos}

	// NOASSOC
	case 25:
		pos := l.in.Skip()
		return lexer.Token{Terminal: NOASSOC, Lexeme: "@none", Pos: pos}

	// PREDEF
	case 27:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: PREDEF, Lexeme: lexeme, Pos: pos}

	// GRAMMER
	case 34:
		pos := l.in.Skip()
		return lexer.Token{Terminal: GRAMMER, Lexeme: "grammar", Pos: pos}

	// IDENT
	case 36:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: IDENT, Lexeme: lexeme, Pos: pos}

	// TOKEN
	case 38:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: TOKEN, Lexeme: lexeme, Pos: pos}

	// STRING
	case 42:
		lexeme, pos := l.in.Lexeme()
		lexeme = lexeme[1 : len(lexeme)-1]
		return lexer.Token{Terminal: STRING, Lexeme: lexeme, Pos: pos}

	// REGEX
	case 46:
		lexeme, pos := l.in.Lexeme()
		lexeme = strings.Trim(lexeme, "/")
		return lexer.Token{Terminal: REGEX, Lexeme: lexeme, Pos: pos}

	// Single-Line COMMENT
	case 47:
		pos := l.in.Skip()
		return lexer.Token{Terminal: COMMENT, Lexeme: "", Pos: pos}

	// Multi-Line COMMENT
	case 50:
		pos := l.in.Skip()
		return lexer.Token{Terminal: COMMENT, Lexeme: "", Pos: pos}
	}

	// ERR
	val, pos := l.in.Lexeme()
	return lexer.Token{
		Terminal: ERR,
		Lexeme:   fmt.Sprintf("lexical error at %s:%s", pos, val),
		Pos:      pos,
	}
}

// advanceDFA determines the next state of a deterministic finite automaton (DFA)
// given the current state and an input symbol.
// It functions as a coded lookup table.
func advanceDFA(state int, r rune) int {
	switch state {
	case 0:
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return 1

		case '=':
			return 2

		case '|':
			return 3

		case '(':
			return 4

		case ')':
			return 5

		case '[':
			return 6

		case ']':
			return 7

		case '{':
			return 8

		case '}':
			return 9

		case '@':
			return 12

		case '$':
			return 26

		case 'g':
			return 28

		case 'a', 'b', 'c', 'd', 'e', 'f' /*g*/, 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 35

		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 37

		case '"':
			return 39

		case '/':
			return 43
		}

	case 1:
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return 1
		}

	case 8:
		switch r {
		case '{':
			return 10
		}

	case 9:
		switch r {
		case '}':
			return 11
		}

	case 12:
		switch r {
		case 'l':
			return 13

		case 'r':
			return 17

		case 'n':
			return 22
		}

	case 13:
		switch r {
		case 'e':
			return 14
		}

	case 14:
		switch r {
		case 'f':
			return 15
		}

	case 15:
		switch r {
		case 't':
			return 16
		}

	case 17:
		switch r {
		case 'i':
			return 18
		}

	case 18:
		switch r {
		case 'g':
			return 19
		}

	case 19:
		switch r {
		case 'h':
			return 20
		}

	case 20:
		switch r {
		case 't':
			return 21
		}

	case 22:
		switch r {
		case 'o':
			return 23
		}

	case 23:
		switch r {
		case 'n':
			return 24
		}

	case 24:
		switch r {
		case 'e':
			return 25
		}

	case 26:
		switch r {
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 27
		}

	case 27:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 27
		}

	case 28:
		switch r {
		case 'r':
			return 29

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q' /*r*/, 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 29:
		switch r {
		case 'a':
			return 30

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			/*a*/ 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 30:
		switch r {
		case 'm':
			return 31

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l' /*m*/, 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 31:
		switch r {
		case 'm':
			return 32

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l' /*m*/, 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 32:
		switch r {
		case 'a':
			return 33

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			/*a*/ 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 33:
		switch r {
		case 'r':
			return 34

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q' /*r*/, 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 34:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 35:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 36:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 36
		}

	case 37:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 38
		}

	case 38:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 38
		}

	case 39:
		switch r {
		case '\\':
			return 40

		case '!' /*"*/, '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 41
		}

	case 40:
		switch r {
		case '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 41
		}

	case 41:
		switch r {
		case '\\':
			return 40

		case '!' /*"*/, '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[' /*\*/, ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 41

		case '"':
			return 42
		}

	case 43:
		switch r {
		case '\\':
			return 44

		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')' /***/, '+', ',', '-', '.', /*/*/
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[' /*\*/, ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 45

		case '/':
			return 47

		case '*':
			return 48
		}

	case 44:
		switch r {
		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 45
		}

	case 45:
		switch r {
		case '\\':
			return 44

		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', /*/*/
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[' /*\*/, ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 45

		case '/':
			return 46
		}

	case 47:
		switch r {
		case '\t',
			' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 47
		}

	case 48:
		switch r {
		case '\t', '\n', '\r',
			' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')' /***/, '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 48

		case '*':
			return 49
		}

	case 49:
		switch r {
		case '\t', '\n', '\r',
			' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', /*/*/
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 48

		case '/':
			return 50
		}
	}

	return errorState
}
