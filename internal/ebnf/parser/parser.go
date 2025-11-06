// Package parser implements a parser for the EBNF language.
package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/list"
	"github.com/moorara/algo/parser"
	"github.com/moorara/algo/parser/lr"

	ebnflexer "github.com/gardenbed/emerge/internal/ebnf/lexer"
)

// Predefs defines the acceptable values for the PREDEF token in the EBNF specification.
// Each value is a predefined regular expression for defining a token.
var Predefs = map[string]string{
	"$DIGIT":   `[0-9]`,
	"$HEX":     `[0-9A-Fa-f]`,
	"$LETTER":  `[A-Za-z]`,
	"$INT":     `-?[0-9]+`,
	"$FLOAT":   `-?[0-9]+(\.[0-9]+)?`,
	"$STRING":  `"([^\\"]|\\[\\"'tnr])*"`,
	"$COMMENT": `(#|//)[^\n\r]*|/\*.*?\*/`,
}

// ProductionFunc is a function that is invoked each time a production rule
// is matched or applied during the parsing process of an input string.
// It passes the index of a production rule instead of the production itself.
//
// It executes the actions associated with the matched production rule,
// such as semantic processing, constructing abstract syntax trees (AST),
// or performing other custom logic required for the parsing process.
//
// The function may return an error, indicating an issue during production rule processing.
// The parser may stop immediately or continue parsing and accumulate more errors.
type ProductionFunc func(int) error

// EvaluateFunc is a function invoked every time a production rule
// is matched or applied during the parsing of an input string.
// It passes the index of a production rule instead of the production itself.
//
// It receives a list of values corresponding to the right-hand side of the matched production
// and expects a value to be returned representing the left-hand side of the production.
//
// The returned value will be subsequently used as an input in the evaluation of other production rules.
// Both the input and output values are of the generic type any.
//
// The caller is responsible for ensuring that each value is converted to the appropriate type based on
// the production rule and the position of the symbol corresponding to the value in the production's right-hand side.
// The input values must retain the same type they were originally evaluated as when returned.
//
// The function may return an error if there are issues with the input values,
// such as mismatched types or unexpected inputs.
type EvaluateFunc func(int, []*lr.Value) (any, error)

// Parser is a parser (a.k.a. syntax analyzer) for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
type Parser struct {
	L lexer.Lexer
}

// New creates a new parser (a.k.a. syntax analyzer) for the EBNF language.
// EBNF (Extended Backus-Naur Form) is used to define context-free grammars and their corresponding languages.
func New(filename string, src io.Reader) (*Parser, error) {
	L, err := ebnflexer.New(filename, src)
	if err != nil {
		return nil, err
	}

	return &Parser{
		L: L,
	}, nil
}

// nextToken wraps the Lexer.NextToken method and ensures an Endmarker token is returned when the end of input is reached.
func (p *Parser) nextToken() (lexer.Token, error) {
	token, err := p.L.NextToken()
	if err != nil && errors.Is(err, io.EOF) {
		token.Terminal, token.Lexeme = grammar.Endmarker, ""
		return token, nil
	}

	return token, err
}

// Parse implements the LR parsing algorithm.
// It analyzes a sequence of input tokens (terminal symbols) provided by the lexical analyzer.
// It attempts to parse the input according to the production rules of the EBNF grammar.
//
// The Parse method invokes the provided functions each time a token or a production rule is matched.
// This allows the caller to process or react to each step of the parsing process.
//
// An error is returned if the input fails to conform to the grammar rules, indicating a syntax issue,
// or if any of the provided functions return an error, indicating a semantic issue.
func (p *Parser) Parse(tokenF parser.TokenFunc, prodF ProductionFunc) error {
	stack := list.NewStack(1024, generic.NewEqualFunc[int]())
	stack.Push(0)

	// Read the first input token.
	token, err := p.nextToken()
	if err != nil {
		return &parser.ParseError{Cause: err}
	}

	for {
		s, _ := stack.Peek()
		a := token.Terminal

		action, param, err := ACTION(s, a)
		if err != nil {
			return &parser.ParseError{
				Description: fmt.Sprintf("unexpected string %q", token.Lexeme),
				Cause:       err,
				Pos:         token.Pos,
			}
		}

		switch action {
		case lr.SHIFT:
			stack.Push(param)

			// Yield the token.
			if tokenF != nil {
				if err := tokenF(&token); err != nil {
					return &parser.ParseError{
						Cause: err,
						Pos:   token.Pos,
					}
				}
			}

			// Read the next input token.
			token, err = p.nextToken()
			if err != nil {
				return &parser.ParseError{Cause: err}
			}

		case lr.REDUCE:
			A, β := productions[param].Head, productions[param].Body

			for range len(β) {
				stack.Pop()
			}

			// An LR parser detects an error when it consults the ACTION table.
			// Errors are never identified by consulting the GOTO table.
			// If ACTION(s, a) is not an error entry, GOTO(t, A) will also not be an error entry.

			t, _ := stack.Peek()
			next := GOTO(t, A)
			stack.Push(next)

			// Yield the production.
			if prodF != nil {
				if err := prodF(param); err != nil {
					return &parser.ParseError{Cause: err}
				}
			}

		case lr.ACCEPT:
			// Accept the input string.
			return nil

		case lr.ERROR:
			// TODO: This is unreachable currently, since ACTION handles the error.
		}
	}
}

// ParseAndBuildAST implements the LR parsing algorithm.
// It analyzes a sequence of input tokens (terminal symbols) provided by a lexical analyzer.
// It attempts to parse the input according to the production rules of the EBNF grammar.
//
// If the input string is valid, the root node of the BNF AST is returned,
// representing the syntactic structure of the input string.
//
// An error is returned if the input fails to conform to the grammar rules, indicating a syntax issue.
func (p *Parser) ParseAndBuildAST() (parser.Node, error) {
	// Stack for constructing the abstract syntax tree.
	nodes := list.NewStack(1024, parser.EqNode)

	err := p.Parse(
		func(token *lexer.Token) error {
			nodes.Push(&parser.LeafNode{
				Terminal: token.Terminal,
				Lexeme:   token.Lexeme,
				Position: token.Pos,
			})

			return nil
		},
		func(i int) error {
			prod := productions[i]

			in := &parser.InternalNode{
				NonTerminal: prod.Head,
				Production:  prod,
			}

			for range len(prod.Body) {
				child, _ := nodes.Pop()
				in.Children = append([]parser.Node{child}, in.Children...) // Maintain correct production body order
			}

			nodes.Push(in)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	// The nodes stack only contains the root of AST at this point.
	root, _ := nodes.Pop()

	return root, nil
}

// ParseAndEvaluate implements the LR parsing algorithm.
// It analyzes a sequence of input tokens (terminal symbols) provided by a lexical analyzer.
// It attempts to parse the input according to the production rules of the EBNF grammar.
//
// During the parsing process, the provided EvaluateFunc is invoked each time a production rule is matched.
// The function is called with values corresponding to the symbols in the body of the production,
// enabling the caller to process and evaluate the input incrementally.
//
// An error is returned if the input fails to conform to the grammar rules, indicating a syntax issue,
// or if the evaluation function returns an error, indicating a semantic issue.
func (p *Parser) ParseAndEvaluate(eval EvaluateFunc) (*lr.Value, error) {
	// Stack for constructing the evaluation hierarchy.
	values := list.NewStack[*lr.Value](1024, nil)

	err := p.Parse(
		func(token *lexer.Token) error {
			copy := token.Pos
			values.Push(&lr.Value{
				Val: token.Lexeme,
				Pos: &copy,
			})

			return nil
		},
		func(i int) error {
			l := len(productions[i].Body)
			rhs := make([]*lr.Value, l)

			// Maintain correct production body order
			for i := l - 1; i >= 0; i-- {
				v, _ := values.Pop()
				rhs[i] = v
			}

			lhs, err := eval(i, rhs)
			if err != nil {
				return err
			}

			v := &lr.Value{Val: lhs}
			if l > 0 {
				v.Pos = rhs[0].Pos
			}

			values.Push(v)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	// The values stack only contains the root of AST at this point.
	root, _ := values.Pop()

	return root, nil
}
