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
	// ERR is the terminal for an error case.
	ERR = grammar.Terminal("ERR")
	// WS is the terminal for whitespace characters.
	WS = grammar.Terminal("WS")
	// DEF is the terminal for "=".
	DEF = grammar.Terminal("DEF")
	// ALT is the terminal for "|".
	ALT = grammar.Terminal("ALT")
	// LPAREN is the terminal for "(".
	LPAREN = grammar.Terminal("LPAREN")
	// RPAREN is the terminal for ")".
	RPAREN = grammar.Terminal("RPAREN")
	// LBRACK is the terminal for "[".
	LBRACK = grammar.Terminal("LBRACK")
	// RBRACK is the terminal for "]".
	RBRACK = grammar.Terminal("RBRACK")
	// LBRACE is the terminal for "{".
	LBRACE = grammar.Terminal("LBRACE")
	// RBRACE is the terminal for "}".
	RBRACE = grammar.Terminal("RBRACE")
	// LLBRACE is the terminal for "{{".
	LLBRACE = grammar.Terminal("LLBRACE")
	// RRBRACE is the terminal for "}}".
	RRBRACE = grammar.Terminal("RRBRACE")
	// GRAMMER is the terminal for "grammar".
	GRAMMER = grammar.Terminal("GRAMMER")
	// IDENT is the terminal for /[a-z][0-9a-z_]*/.
	IDENT = grammar.Terminal("IDENT")
	// TOKEN is the terminal for /[A-Z][0-9A-Z_]*/.
	TOKEN = grammar.Terminal("TOKEN")
	// STRING is the terminal for /"([\x21\x23-\x5B\x5D-\x7E]|\\[\x21-\x7E])+"/.
	STRING = grammar.Terminal("STRING")
	// REGEX is the terminal for /\/([\x20-\x2E\x30-\x5B\x5D-\x7E]|\\[\x20-\x7E])*\//.
	REGEX = grammar.Terminal("REGEX")
	// COMMENT is the terminal for single-line and multi-line comments.
	COMMENT = grammar.Terminal("COMMENT")
)

// inputBuffer is an interface for the input.Input struct.
type inputBuffer interface {
	Next() (rune, error)
	Retract()
	Lexeme() (string, lexer.Position)
	Skip() lexer.Position
}

// Lexer is a lexical analyzer for an extension of EBNF language.
type Lexer struct {
	in inputBuffer
}

// New creates a new lexical analyzer for an extension of EBNF language.
func New(filename string, src io.Reader) (*Lexer, error) {
	in, err := input.New(filename, src, bufferSize)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		in: in,
	}, nil
}

// NextToken returns the next token from the input.
// An io.EOF error will be returned when the end of input is reached.
func (l *Lexer) NextToken() (lexer.Token, error) {
	for curr, next := 0, 0; ; curr = next {
		// Read the next character from the input
		r, err := l.in.Next()
		if err != nil {
			return lexer.Token{}, err
		}

		// Simulate the DFA
		next = advanceDFA(curr, r)
		if next == errorState {
			// Move one character backward in the input after DFA being stuck with the last character
			l.in.Retract()

			token := l.evalDFA(curr)
			switch token.Terminal {
			case ERR:
				return lexer.Token{}, errors.New(token.Lexeme)
			case WS, COMMENT: // Skip whitespaces and comments
				return l.NextToken()
			default:
				return token, nil
			}
		}
	}
}

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

	// LLBRACE
	case 9:
		pos := l.in.Skip()
		return lexer.Token{Terminal: LLBRACE, Lexeme: "{{", Pos: pos}

	// RBRACE
	case 10:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RBRACE, Lexeme: "}", Pos: pos}

	// RRBRACE
	case 11:
		pos := l.in.Skip()
		return lexer.Token{Terminal: RRBRACE, Lexeme: "}}", Pos: pos}

	// GRAMMER
	case 18:
		pos := l.in.Skip()
		return lexer.Token{Terminal: GRAMMER, Lexeme: "grammar", Pos: pos}

	// IDENT
	case 20:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: IDENT, Lexeme: lexeme, Pos: pos}

	// TOKEN
	case 22:
		lexeme, pos := l.in.Lexeme()
		return lexer.Token{Terminal: TOKEN, Lexeme: lexeme, Pos: pos}

	// STRING
	case 26:
		lexeme, pos := l.in.Lexeme()
		lexeme = lexeme[1 : len(lexeme)-1]
		return lexer.Token{Terminal: STRING, Lexeme: lexeme, Pos: pos}

	// REGEX
	case 30:
		lexeme, pos := l.in.Lexeme()
		lexeme = strings.Trim(lexeme, "/")
		return lexer.Token{Terminal: REGEX, Lexeme: lexeme, Pos: pos}

	// Single-Line COMMENT
	case 31:
		pos := l.in.Skip()
		return lexer.Token{Terminal: COMMENT, Lexeme: "", Pos: pos}

	// Multi-Line COMMENT
	case 34:
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

func advanceDFA(state int, r rune) int {
	switch state {
	case 0:
		switch r {
		case ' ', '\t', '\n', '\r', '\f', '\v':
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
			return 10

		case 'g':
			return 12

		case 'a', 'b', 'c', 'd', 'e', 'f', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 19

		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 21

		case '"':
			return 23

		case '/':
			return 27
		}

	case 1:
		switch r {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			return 1
		}

	case 8:
		switch r {
		case '{':
			return 9
		}

	case 10:
		switch r {
		case '}':
			return 11
		}

	case 12:
		switch r {
		case 'r':
			return 13

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 13:
		switch r {
		case 'a':
			return 14

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 14:
		switch r {
		case 'm':
			return 15

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 15:
		switch r {
		case 'm':
			return 16

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 16:
		switch r {
		case 'a':
			return 17

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 17:
		switch r {
		case 'r':
			return 18

		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 18:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 19:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 20:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			return 20
		}

	case 21:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 22
		}

	case 22:
		switch r {
		case '_',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			return 22
		}

	case 23:
		switch r {
		case '\\':
			return 24

		case '!', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 25
		}

	case 24:
		switch r {
		case '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 25
		}

	case 25:
		switch r {
		case '"':
			return 26

		case '\\':
			return 24

		case '!', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 25
		}

	case 27:
		switch r {
		case '*':
			return 32

		case '/':
			return 31

		case '\\':
			return 28

		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '+', ',', '-', '.',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 29
		}

	case 28:
		switch r {
		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 29
		}

	case 29:
		switch r {
		case '/':
			return 30

		case '\\':
			return 28

		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 29
		}

	case 31:
		switch r {
		case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 31
		}

	case 32:
		switch r {
		case '*':
			return 33

		case '\t', '\n', '\r',
			' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '+', ',', '-', '.', '/',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 32
		}

	case 33:
		switch r {
		case '/':
			return 34

		case '\t', '\n', '\r',
			' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			':', ';', '<', '=', '>', '?', '@',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
			'[', '\\', ']', '^', '_', '`',
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'{', '|', '}', '~':
			return 32
		}
	}

	return errorState
}
