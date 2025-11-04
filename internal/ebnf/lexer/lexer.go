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

// evalDFA examines the final state of a deterministic finite automaton (DFA) after it has stopped processing input.
// Based on the last encountered state, it returns the corresponding token and advances the input buffer reader.
// If the final state is invalid, it returns an ERR token with the Lexeme set to the error message.
func (l *Lexer) evalDFA(state int) lexer.Token {
	switch state {
	// Whitespace
	case 1:
		pos := l.in.Skip()
		return lexer.Token{Terminal: WS, Lexeme: "", Pos: pos}

	// Newline
	case 2:
		pos := l.in.Skip()
		return lexer.Token{Terminal: EOL, Lexeme: "", Pos: pos}

	// DEF
	case 3:
		pos := l.in.Skip()
		return lexer.Token{Terminal: DEF, Lexeme: "=", Pos: pos}

	// SEMI
	case 4:
		pos := l.in.Skip()
		return lexer.Token{Terminal: SEMI, Lexeme: ";", Pos: pos}

	// ALT
	case 5:
		pos := l.in.Skip()
		return lexer.Token{Terminal: ALT, Lexeme: "|", Pos: pos}

	// LPAREN
	case 6:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LPAREN, Lexeme: "(", Pos: pos}

	// RPAREN
	case 7:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RPAREN, Lexeme: ")", Pos: pos}

	// LBRACK
	case 8:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LBRACK, Lexeme: "[", Pos: pos}

	// RBRACK
	case 9:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RBRACK, Lexeme: "]", Pos: pos}

	// LBRACE
	case 10:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LBRACE, Lexeme: "{", Pos: pos}

	// RBRACE
	case 11:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RBRACE, Lexeme: "}", Pos: pos}

	// LLBRACE
	case 12:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LLBRACE, Lexeme: "{{", Pos: pos}

	// RRBRACE
	case 13:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RRBRACE, Lexeme: "}}", Pos: pos}

	// LANGLE
	case 14:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LANGLE, Lexeme: "<", Pos: pos}

	// RANGLE
	case 15:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RANGLE, Lexeme: ">", Pos: pos}

	// PREDEF
	case 17:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: PREDEF, Lexeme: lexeme, Pos: pos}

	// LASSOC
	case 22:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LASSOC, Lexeme: "@left", Pos: pos}

	// RASSOC
	case 27:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RASSOC, Lexeme: "@right", Pos: pos}

	// NOASSOC
	case 31:
		pos := l.in.Skip()
		return lexer.Token{Terminal: NOASSOC, Lexeme: "@none", Pos: pos}

	// GRAMMER
	case 38:
		pos := l.in.Skip()
		return lexer.Token{Terminal: GRAMMER, Lexeme: "grammar", Pos: pos}

	// IDENT
	case 32, 33, 34, 35, 36, 37, 39:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: IDENT, Lexeme: lexeme, Pos: pos}

	// TOKEN
	case 40:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: TOKEN, Lexeme: lexeme, Pos: pos}

	// STRING
	case 61:
		lexeme, pos := l.in.Lexeme()
		lexeme = lexeme[1 : len(lexeme)-1]
		return lexer.Token{Terminal: STRING, Lexeme: lexeme, Pos: pos}

	// REGEX
	case 65:
		lexeme, pos := l.in.Lexeme()
		lexeme = lexeme[1 : len(lexeme)-1]
		return lexer.Token{Terminal: REGEX, Lexeme: lexeme, Pos: pos}

	// Single-Line COMMENT
	case 66:
		pos := l.in.Skip()
		return lexer.Token{Terminal: COMMENT, Lexeme: "", Pos: pos}

	// Multi-Line COMMENT
	case 69:
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
		switch {
		case r == '\t', r == ' ':
			return 1
		case r == '\n', r == '\r':
			return 2
		case r == '=':
			return 3
		case r == ';':
			return 4
		case r == '|':
			return 5
		case r == '(':
			return 6
		case r == ')':
			return 7
		case r == '[':
			return 8
		case r == ']':
			return 9
		case r == '{':
			return 10
		case r == '}':
			return 11
		case r == '<':
			return 14
		case r == '>':
			return 15
		case r == '$':
			return 16
		case r == '@':
			return 18
		case r == 'g':
			return 32
		case 'a' <= r && r <= 'f',
			'h' <= r && r <= 'z':
			return 39
		case 'A' <= r && r <= 'Z':
			return 40
		case r == '"':
			return 41
		case r == '/':
			return 62
		}

	case 1:
		switch {
		case r == '\t', r == ' ':
			return 1
		}

	case 2:
		switch {
		case r == '\n', r == '\r':
			return 2
		}

	case 10:
		switch {
		case r == '{':
			return 12
		}

	case 11:
		switch {
		case r == '}':
			return 13
		}

	case 16:
		switch {
		case 'A' <= r && r <= 'Z':
			return 17
		}

	case 17:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'Z',
			r == '_':
			return 17
		}

	case 18:
		switch {
		case r == 'l':
			return 19
		case r == 'n':
			return 28
		case r == 'r':
			return 23
		}

	case 19:
		switch {
		case r == 'e':
			return 20
		}

	case 20:
		switch {
		case r == 'f':
			return 21
		}

	case 21:
		switch {
		case r == 't':
			return 22
		}

	case 23:
		switch {
		case r == 'i':
			return 24
		}

	case 24:
		switch {
		case r == 'g':
			return 25
		}

	case 25:
		switch {
		case r == 'h':
			return 26
		}

	case 26:
		switch {
		case r == 't':
			return 27
		}

	case 28:
		switch {
		case r == 'o':
			return 29
		}

	case 29:
		switch {
		case r == 'n':
			return 30
		}

	case 30:
		switch {
		case r == 'e':
			return 31
		}

	case 32:
		switch {
		case r == 'r':
			return 33
		case '0' <= r && r <= '9',
			r == '_',
			'a' <= r && r <= 'q',
			's' <= r && r <= 'z':
			return 39
		}

	case 33:
		switch {
		case r == 'a':
			return 34
		case '0' <= r && r <= '9',
			r == '_',
			'b' <= r && r <= 'z':
			return 39
		}

	case 34:
		switch {
		case r == 'm':
			return 35
		case '0' <= r && r <= '9',
			r == '_',
			'a' <= r && r <= 'l',
			'n' <= r && r <= 'z':
			return 39
		}

	case 35:
		switch {
		case r == 'm':
			return 36
		case '0' <= r && r <= '9',
			r == '_',
			'a' <= r && r <= 'l',
			'n' <= r && r <= 'z':
			return 39
		}

	case 36:
		switch {
		case r == 'a':
			return 37
		case '0' <= r && r <= '9',
			r == '_',
			'b' <= r && r <= 'z':
			return 39
		}

	case 37:
		switch {
		case r == 'r':
			return 38
		case '0' <= r && r <= '9',
			r == '_',
			'a' <= r && r <= 'q',
			's' <= r && r <= 'z':
			return 39
		}

	case 38:
		switch {
		case '0' <= r && r <= '9':
			return 39
		case r == '_':
			return 39
		case 'a' <= r && r <= 'z':
			return 39
		}

	case 39:
		switch {
		case '0' <= r && r <= '9',
			r == '_',
			'a' <= r && r <= 'z':
			return 39
		}

	case 40:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'Z',
			r == '_':
			return 40
		}

	case 41:
		switch {
		case r == '"':
			return 61
		case r == '\\':
			return 42
		case 0x00 <= r && r <= 0x21,
			0x23 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 41
		}

	case 42:
		switch {
		case r == '"',
			r == '\'',
			r == '\\',
			r == 'n',
			r == 'r',
			r == 't':
			return 43
		case r == 'x':
			return 44
		case r == 'u':
			return 47
		case r == 'U':
			return 52
		}

	case 43:
		switch {
		case r == '"':
			return 61
		case r == '\\':
			return 42
		case 0x00 <= r && r <= 0x21,
			0x23 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 41
		}

	case 44:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 45
		}

	case 45:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 46
		}

	case 46:
		switch {
		case r == '"':
			return 61
		case r == '\\':
			return 42
		case 0x00 <= r && r <= 0x21,
			0x23 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 41
		}

	case 47:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 48
		}

	case 48:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 49
		}

	case 49:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 50
		}

	case 50:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 51
		}

	case 51:
		switch {
		case r == '"':
			return 61
		case r == '\\':
			return 42
		case 0x00 <= r && r <= 0x21,
			0x23 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 41
		}

	case 52:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 53
		}

	case 53:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 54
		}

	case 54:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 55
		}

	case 55:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 56
		}

	case 56:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 57
		}

	case 57:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 58
		}

	case 58:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 59
		}

	case 59:
		switch {
		case '0' <= r && r <= '9',
			'A' <= r && r <= 'F',
			'a' <= r && r <= 'f':
			return 60
		}

	case 60:
		switch {
		case r == '"':
			return 61
		case r == '\\':
			return 42
		case 0x00 <= r && r <= 0x21,
			0x23 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 41
		}

	case 62:
		switch {
		case r == '*':
			return 67
		case r == '/':
			return 66
		case r == '\\':
			return 63
		case 0x00 <= r && r <= 0x29,
			0x2B <= r && r <= 0x2E,
			0x30 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 64
		}

	case 63:
		switch {
		case 0x00 <= r && r <= 0x10FFFF:
			return 64
		}

	case 64:
		switch {
		case r == '/':
			return 65
		case r == '\\':
			return 63
		case 0x00 <= r && r <= 0x2E,
			0x30 <= r && r <= 0x5B,
			0x5D <= r && r <= 0x10FFFF:
			return 64
		}

	case 66:
		switch {
		case 0x00 <= r && r <= 0x09,
			0x0E <= r && r <= 0x10FFFF:
			return 66
		}

	case 67:
		switch {
		case r == '*':
			return 68
		case 0x00 <= r && r <= 0x29,
			0x2B <= r && r <= 0x10FFFF:
			return 67
		}

	case 68:
		switch {
		case r == '*':
			return 68
		case r == '/':
			return 69
		case 0x00 <= r && r <= 0x29,
			0x2B <= r && r <= 0x2E,
			0x30 <= r && r <= 0x10FFFF:
			return 67
		}
	}

	return errorState
}
