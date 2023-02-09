# Theory of Languages and Compilers

## Languages

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

Colloquially, we say that *finnite automata* cannot count, meaning that
a finite automaton cannot accept a language like `{a`<sup>`n`</sup> `b`<sup>`n`</sup> `| n >= 1}`
that would require it to keep count of the number of a's before it sees the b's.
Likewise, a context-free grammar can count two items but not three.

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
  2. A set of input symbols `Σ`, the *input alphabet* (we assume the empty string `ε` is never a member of `Σ`).
  3. A *transition function* that determines a set of next states for each state and for each symbol in `Σ ∪ {ε}`.
  4. A state <code>s<sub>0</sub> ∈ S</code> that is distinguished as the *start state* or *initial state*.
  5. A set of states `F ⊆ S` that is distinguished as the *accepting states* or *final states*.

A deterministic finite automaton (DFA) is a special case of an NFA where:

  - There are no moves on input `ε`, and
  - For each state `s` and input symbol `a`, there is exactly one edge out of `s` labeled `a`.

## Compilers

A compiler is a program that can read a program in the source language and translate it into an equivalent program in the target language.
A compiler has two major parts:

  - **Analysis** (front-end)
  - **Synthesis** (back-end)

The analysis part breaks up the source program into constituent pieces and imposes a grammatical structure on them.
It then uses this structure to create an **intermediate representation** of the source program.
If the analysis part finds that the source program is either syntactically malformed or semantically unsound,
then it must provide informative messages for the user to take actions.
The analysis part also collects information about the source program and stores it in a data structure called **symbol table**.
The **symbol table** is then passed along with the **intermediate representation** to the synthesis part.

The synthesis part constructs the desired **target program** from the intermediate representation and the information in the symbol table.

```
                                          Character Stream
                                                 │
                                ┌────────────────┴───────────────────┐
           .....................│         Lexical Analyzer           │
           .                    └────────────────┬───────────────────┘
           .                                     │
           .                                Token Stream
           .                                     │
           .                    ┌────────────────┴───────────────────┐
           .....................│          Syntax Analyzer           │
           .                    └────────────────┬───────────────────┘
           .                                     │
           .                                Syntax Tree
           .                                     │
           .                    ┌────────────────┴───────────────────┐
           .....................│        Semantic Analyzer           │
           .                    └────────────────┬───────────────────┘
           .                                     │
           .                                Syntax Tree
 ┌───────────────────┐                           │
 │                   │          ┌────────────────┴───────────────────┐
 │   Symbol Table    │..........│    Intermediate Code Generator     │
 │                   │          └────────────────┬───────────────────┘
 └───────────────────┘                           │
           .                        Intermediate Representation
           .                                     │
           .                    ┌────────────────┴───────────────────┐
           .....................│ Machine-Independent Code Optimizer │
           .                    └────────────────┬───────────────────┘
           .                                     │
           .                        Intermediate Representation
           .                                     │
           .                    ┌────────────────┴───────────────────┐
           .....................│            Code Generator          │
           .                    └────────────────┬───────────────────┘
           .                                     │
           .                            Target-Machine Code
           .                                     │
           .                    ┌────────────────┴───────────────────┐
           .....................│  Machine-Dependent Code Optimizer  │
                                └────────────────┬───────────────────┘
                                                 │
                                        Target-Machine Code
                                                 │
                                                 ▼
```

In an implementation, activities from several phases may be grouped together into a **pass**
that reads an input file and writes an output file.

### Lexical Analysis

The first phase of a compiler is called *lexical analysis* or *scanning*.
There are a number of reasons for the separation of *lexical analysis* and *syntax analysis*:

  - Simplicity of design
  - Compiler efficiency
  - Compiler portability

The lexical analyzer performs the following tasks:

  - Reading the **stream of characters** of the source program, grouping them into **lexemes**, and producing a **stream of of tokens**.
  - Stripping out comments and whitespace (blank, newline, tab, ...).
  - Correlating error messages generated by the compiler with the source program.

It is common for the lexical analyzer to interact with the *symbol table* as well.
Often, information about an identifier (lexeme, type, location at which it is first found, etc.) is kept in the symbol table.
Hence, the appropriate attribute value for an identifier is a pointer to the symbol-table entry for that identifier.

It is hard for a lexical analyzer to tell, without the aid of other components, that there is a source-code error.

  - A **token** is a tuple consisting of a token name alongside some optional attributes (*lexeme*, *position*, etc.).
    The token name is an abstract symbol representing a kind of lexical unit.
  - A **lexeme** is a sequence of characters (from the language's alphabet)
    that matches the **pattern** for a token and is identified by the lexical analyzer as an instance of that token.

The parser calls the lexical analyzer using the `getNextToken` command, causes the lexical analyzer to read characters
from its input until it can identify the next lexeme and produce the next token, which it returns to the parser.

For a more in-depth theory of lexical analysis and lexers, please see [Theory of Lexers](./2-lexer_theory.md).

### Syntax Analysis

The second phase of the compiler is *syntax analysis* or *parsing*.
The parser creates a **syntax tree** (intermediate representation) that depicts the grammatical structure of the token stream.

There are no strong guidelines for what to put into the lexical rules vs. to the syntactic rules.
Regular expressions are favorable for formulating the structure of constructs such as keywords, identifiers, and whitespaces.
Grammars, on the other hand, are convenient for formulating nested structures such as matching parentheses, brackets, and so on.

  1. Separating the syntactic structure of a language into lexical and non-lexical parts provides
     a convenient way of modularizing the front end of a compiler into two manageable components.
  1. We do not need a notation as expressive as grammars for describing the lexical rules of a language that are quite simple.
  1. Regular expressions provide a more concise and easier-to-understand notation for tokens than grammars.
  1. More eficient lexical analyzers can be generated automatically from regular expressions than from arbitrary grammars.

The syntax of programming language constructs can be specified by context-free grammars or *BNF (Backus-Naur Form)* notation.

  - A grammar gives a precise and easy-to-understand syntactic specification of a language.
  - A grammar allows a language to be developed and evolved iteratively.
  - From certain classes of grammars, we can automatically construct an efficient parser.
    - During the parser construction, we can reveal syntactic ambiguities that might have slipped through the initial design of a language.

For a more in-depth theory of syntax analysis and parsers, please see [Theory of Parsers](./3-parser_theory.md).

### Semantic Analysis

The semantic analyzer uses the *syntax tree* and the information in the *symbol table*
to check the source program for semantic consistency with the language definition.
It also gathers **type information** and saves it in either the syntax tree or the symbol table.
An important part of semantic analysis is **type checking**, where the compiler checks that each operator has matching operands.

### Intermediate Code Generation

After syntax and semantic analysis of the source program, many compilers generate an explicit *low-level* or *machine-like* intermediate representation.
This intermediate representation can be thought of as a program for an *abstract machine*.

### Machine-Independent Code Optimization

The machine-independent code optimization phase attempts to improve the intermediate code so that better target code will result.
Usually better means faster, but other objectives may be desired, such as shorter code, or target code that consumes less power.

### Code Generation

The code generator takes as input an intermediate representation of the source program and maps it into the *target language*.

### Machine-Dependent Code Optimization

*Optimization* is a misnomer!
In fact, there is no way that the code produced by a compiler can be guaranteed to be as fast or faster than any other code that performs the same task.
