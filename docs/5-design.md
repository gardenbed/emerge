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

#### Input Buffer

We employ the *two-buffer* scheme explained [here](./2-lexer_theory.md#input-buffering).
The two buffers are implemented as one buffer divided into two halves.

### Parser Design

## Resources

  - [Extended Backus–Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
  - [Parser Combinator](https://en.wikipedia.org/wiki/Parser_combinator)
  - [Let's Build a Regex Engine](https://kean.blog/post/lets-build-regex)
