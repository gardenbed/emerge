# Theory

## Background and Context

A language is defined by four sets:

  - **Alphabet:** A finite set of symbols.
  - **Words:** A finite or infinite set of combinations of symbols from the language's alphabet.
  - **Grammar:** A finite set of rules determining how sentences in the language can be constructed.
  - **Semantic:** A finite set of rules determining which sentences in the language are valid (meaningful).

In the theory of formal languages, there are four types of grammars and respectively four types of languages.
These classes of formal grammars have been formulated by [Noam Chomsky](https://en.wikipedia.org/wiki/Noam_Chomsky) in 1956.

  - **Type-0:** The most general kind of grammars and languages (a.k.a. *unrestricted grammars*).
                These grammars can generate all languages that can be represented and decided by a [Turing machine](https://en.wikipedia.org/wiki/Turing_machine).
  - **Type-1:** *Context-sensitive* grammars and languages are a superset of Type-2 and Type-1 languages.
                These languages can be represented and decided by a [linear bounded automaton](https://en.wikipedia.org/wiki/Linear_bounded_automaton).
  - **Type-2:** *Context-free* grammars and languages are a superset of Type-3 languages.
                In context-free grammars, all production rules must have only one non-terminal symbol on the left-hand side.
                In other words, regardless of which context the non-terminal symbol appears, it should always be interpreted the same way.
                These languages can be represented and decided by a [pushdown automaton](https://en.wikipedia.org/wiki/Pushdown_automaton).
  - **Type-3:** The most restrictive kind of grammars and languages (a.k.a. *regular grammars/languages*).
                These languages can be represented and decided by a [finite state machine](https://en.wikipedia.org/wiki/Finite-state_machine).

### Type-2 Languages

Most programming languages are defined using Type-2 grammars.

A context-free grammar `G` is defined by the 4-tuple `{G=(V,Σ,R,S)}`, where:

  1. `V` is a finite set of non-terminals.
  2. `Σ` is a finite set of terminals, disjoint from `V`.
  3. `R` is a finite relation in `V × (V ∪ Σ)*`, where the asterisk represents the *Kleene star* operation.
     The members of `R` are called the production rules of the grammar (also commonly symbolized by a `P`).
  4. `S ∈ V` is the start symbol, used to represent the whole program.

There are several slightly different notations for context-free garmmars.
We use a variation of [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form).

### Type-3 Languages

Type-3 languages can be represented and decided by a *finite automata*.
Finite automata are recognizers; they simply accept or reject an input string.
In a finite automata, the next state only depends on the current state the next input symbol.

Finite automata come in two flavors:

  - **Deterministic finite automata (DFA):** For each state and each symbol from the language alphabet, there is exactly one edge with symbol to the next state.
  - **Nondeterministic finite automata (NFA):** A symbol from the language alphabet can label multiple edges from the same state.
                                                The empty string `ε` can also be used for labelling edges.

Both DFA and NFA can recognize the same languages.
In fact, every nondeterministic finite automata (NFA) can be converted to a deterministic finite automaton (DFA).

A nondeterministic finite automaton (NFA) is defined as follows:

  1. A finite set of states `S`.
  2. A set of input symbols `Σ`, the input alphabet (we assume the empty string `ε` is never a member of `Σ`).
  3. A transition function that determines a set of next states for each state and for each symbol in `Σ ∪ {ε}`.
  4. A state <code>s<sub>0</sub> ∈ S</code> that is distinguished as the *start state* or *initial state*.
  5. A set of states `F`, a subset of `S`, that is distinguished as the *accepting states* or *final states*.

A deterministic finite automaton (DFA) is a special case of an NFA where:

  - There are no moves on input `ε`, and
  - For each state `s` and input symbol `a`, there is exactly one edge out of `s` labeled `a`.

## Lexical Analysis

  - A **token** is a tuple consisting of a token name alongside some optional attribute values (*lexeme, position, etc.*).
    The token name is an abstract symbol representing a kind of lexical unit.
  - A **lexeme** is a sequence of characters (from the language's alphabet)
    that matches the pattern for a token and is identified by the lexical analyzer as an instance of that token.

## Syntax Analysis

  - [Top-Down Parsers](https://en.wikipedia.org/wiki/Top-down_parsing)
    - [LL Parser](https://en.wikipedia.org/wiki/LL_parser)
    - [Recursive Descent Parser](https://en.wikipedia.org/wiki/Recursive_descent_parser)
      - [Parser Combinator](https://en.wikipedia.org/wiki/Parser_combinator)
  - [Bottom-Up Parsering](https://en.wikipedia.org/wiki/Bottom-up_parsing)
    - [LR Parser](https://en.wikipedia.org/wiki/LR_parser)
      - [Simple LR Parser](https://en.wikipedia.org/wiki/Simple_LR_parser)
      - [LALR Parser](https://en.wikipedia.org/wiki/LALR_parser)
      - [Canonical LR Parser](https://en.wikipedia.org/wiki/Canonical_LR_parser)
      - [GLR Parser](https://en.wikipedia.org/wiki/GLR_parser)

## Semantic Analysis

## References

  - [Compilers: Principles, Techniques, and Tools, 2nd Edition](https://www.pearson.com/us/higher-education/program/Aho-Compilers-Principles-Techniques-and-Tools-2nd-Edition/PGM167067.html)
  - [Parsing Techniques: A Practical Guide, Second Edition](https://link.springer.com/book/10.1007/978-0-387-68954-8)
