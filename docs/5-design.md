# Design

In this document, we go over some of the design decisions and rationals behind **Emerge**.

## Regular Expression

### Language Design

The following features are NOT included in the Emerge's Regular Expression language.
They seem unnecessary for the purpose of designing and defining tokens of a language.

  - Backreference
  - Unicode character class
  - Non-capturing group modifier `?:`
  - Anchors:
      - Word boundary `\b`
      - Non-word boundary`\B`
      - Start of string only `\A`
      - End of string only `\Z`
      - End of string only (not newline) `\z`
      - Previous match end `\G`

### Parser Design

For building a *Lexer* from an EBNF input, we need to parse regular expression patterns in the input,
so we can construct a [DFA](https://en.wikipedia.org/wiki/Deterministic_finite_automaton) for each pattern.
To this end, we need to first build a parser for regular expressions.

Building a regex parser is fairly simple and straightforward.
Implementing a separate *lexer* and *parser* for regular expressions is an inessential complexity
(i.e., whitespace characters do not need to be stripped out).

We have built a simple parser for Emerge's regular expressions that takes care of terminal symbols as well.
This parser is implemented as a [Top-Down Parser](https://en.wikipedia.org/wiki/Top-down_parsing) using [Parser Combinators](https://en.wikipedia.org/wiki/Parser_combinator).

A parser combinator is a *higher-order function* that accepts a *stream of input characters* and returns a *parsing result*.
Using a *functional programming* style,
we can implement a *Type-2* grammar as a single function that receives an input stream and returns an *abstract syntax tree*.

We will later use regular expression ASTs to construct DFAs needed for generating a lexer for an EBNF grammar.

## Extended Backus-Naur Form

### Language Design

The following terminal symbols are removed from the Emerge's EBNF language for simplicity and brevity.

  - Concatenation (`,`)
  - Termination (`;`)
  - Single quotation (`'`)

The Solidus (Slash) character (`/`) is added to the Emerge's EBNF language for defining regex patterns.

### Lexer Design

![Lexer DFA](./lexer_dfa.png)

<details>
  <summary>Lexer DFA Code</summary>

  ```go
  dfa := auto.NewDFA(0, auto.States{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 18, 20, 22, 26, 30, 31, 34})

  dfa.Add(0, ' ', 1)
  dfa.Add(0, '\t', 1)
  dfa.Add(0, '\n', 1)
  dfa.Add(0, '\r', 1)
  dfa.Add(0, '\f', 1)
  dfa.Add(0, '\v', 1)
  dfa.Add(1, ' ', 1)
  dfa.Add(1, '\t', 1)
  dfa.Add(1, '\n', 1)
  dfa.Add(1, '\r', 1)
  dfa.Add(1, '\f', 1)
  dfa.Add(1, '\v', 1)

  dfa.Add(0, '=', 2)
  dfa.Add(0, '|', 3)
  dfa.Add(0, '(', 4)
  dfa.Add(0, ')', 5)
  dfa.Add(0, '[', 6)
  dfa.Add(0, ']', 7)
  dfa.Add(0, '{', 8)
  dfa.Add(8, '{', 9)
  dfa.Add(0, '}', 10)
  dfa.Add(10, '}', 11)

  dfa.Add(0, 'g', 12)
  dfa.Add(12, 'r', 13)
  dfa.Add(13, 'a', 14)
  dfa.Add(14, 'm', 15)
  dfa.Add(15, 'm', 16)
  dfa.Add(16, 'a', 17)
  dfa.Add(17, 'r', 18)

  //==================================================< IDENTIFIER >==================================================

  for r := 'a'; r <= 'z'; r++ {
    if r != 'g' {
      dfa.Add(0, auto.Symbol(r), 19)
    }

    if r != 'r' {
      dfa.Add(12, auto.Symbol(r), 20)
      dfa.Add(17, auto.Symbol(r), 20)
    }
    if r != 'a' {
      dfa.Add(13, auto.Symbol(r), 20)
      dfa.Add(16, auto.Symbol(r), 20)
    }
    if r != 'm' {
      dfa.Add(14, auto.Symbol(r), 20)
      dfa.Add(15, auto.Symbol(r), 20)
    }
    dfa.Add(18, auto.Symbol(r), 20)
    dfa.Add(19, auto.Symbol(r), 20)
    dfa.Add(20, auto.Symbol(r), 20)
  }

  for r := '0'; r <= '9'; r++ {
    dfa.Add(12, auto.Symbol(r), 20)
    dfa.Add(13, auto.Symbol(r), 20)
    dfa.Add(14, auto.Symbol(r), 20)
    dfa.Add(15, auto.Symbol(r), 20)
    dfa.Add(16, auto.Symbol(r), 20)
    dfa.Add(17, auto.Symbol(r), 20)
    dfa.Add(18, auto.Symbol(r), 20)
    dfa.Add(19, auto.Symbol(r), 20)
    dfa.Add(20, auto.Symbol(r), 20)
  }

  dfa.Add(12, '_', 20)
  dfa.Add(13, '_', 20)
  dfa.Add(14, '_', 20)
  dfa.Add(15, '_', 20)
  dfa.Add(16, '_', 20)
  dfa.Add(17, '_', 20)
  dfa.Add(18, '_', 20)
  dfa.Add(19, '_', 20)
  dfa.Add(20, '_', 20)

  //==================================================< TOKEN >==================================================

  for r := 'A'; r <= 'Z'; r++ {
    dfa.Add(0, auto.Symbol(r), 21)
    dfa.Add(21, auto.Symbol(r), 22)
    dfa.Add(22, auto.Symbol(r), 22)
  }

  for r := '0'; r <= '9'; r++ {
    dfa.Add(21, auto.Symbol(r), 22)
    dfa.Add(22, auto.Symbol(r), 22)
  }

  dfa.Add(21, '_', 22)
  dfa.Add(22, '_', 22)

  //==================================================< STRING >==================================================

  dfa.Add(0, '"', 23)
  dfa.Add(23, '\\', 24)
  dfa.Add(25, '\\', 24)
  dfa.Add(25, '"', 26)

  for r := 0x21; r <= 0x7E; r++ {
    dfa.Add(24, auto.Symbol(r), 25)
    if r != '"' && r != '\\' {
      dfa.Add(23, auto.Symbol(r), 25)
      dfa.Add(25, auto.Symbol(r), 25)
    }
  }

  //==================================================< REGEX >==================================================

  dfa.Add(0, '/', 27)
  dfa.Add(27, '\\', 28)
  dfa.Add(29, '\\', 28)
  dfa.Add(29, '/', 30)

  for r := 0x20; r <= 0x7E; r++ {
    if r != '*' && r != '/' && r != '\\' {
      dfa.Add(27, auto.Symbol(r), 29)
    }

    dfa.Add(28, auto.Symbol(r), 29)

    if r != '/' && r != '\\' {
      dfa.Add(29, auto.Symbol(r), 29)
    }
  }

  //==================================================< SINGLE-LINE COMMENT >==================================================

  dfa.Add(27, '/', 31)

  for r := 0x20; r <= 0x7E; r++ {
    dfa.Add(31, auto.Symbol(r), 31)
  }

  //==================================================< MULTI-LINE COMMENT >==================================================

  dfa.Add(27, '*', 32)
  dfa.Add(32, '*', 33)
  dfa.Add(33, '/', 34)

  for _, r := range []rune{'\t', '\n', '\r'} {
    dfa.Add(32, auto.Symbol(r), 32)
    dfa.Add(33, auto.Symbol(r), 32)
  }

  for r := 0x20; r <= 0x7E; r++ {
    if r != '*' {
      dfa.Add(32, auto.Symbol(r), 32)
    }

    if r != '/' {
      dfa.Add(33, auto.Symbol(r), 32)
    }
  }
  ```
</details>

#### Input Buffer

We employ the *two-buffer* scheme explained [here](./2-lexer_theory.md#input-buffering).
The two buffers are implemented as one buffer divided into two halves.

### Parser Design

## Resources

  - [Extended Backus–Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
  - [Parser Combinator](https://en.wikipedia.org/wiki/Parser_combinator)
  - [Let's Build a Regex Engine](https://kean.blog/post/lets-build-regex)
