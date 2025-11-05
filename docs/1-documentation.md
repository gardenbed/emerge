# Emerge Documentation

## Writing A Grammar

EBNF (Extended Backus-Naur Form) is a formal notation, similar to BNF,
used to define [context-free grammar](./4-parser_theory.md#context-free-grammars).
It introduces additional high-level constructs, such as the *Kleene Star* and *Kleene Plus* closures,
which make grammar definitions more concise and readable.

Context-free grammars belong to the Type-2 category in the *Chomsky hierarchy*
and are used to describe the syntax of nearly all programming languages.
Formally, a context-free grammar consists of:

  1. A set of **terminal** symbols.
  1. A set of **non-terminals** symbols.
  1. A set of **productions**, where each production consists of:
      - A **head** (a non-terminal symbol).
      - A **body** (a sequence of terminals and/or non-terminals).
  1. A designated **start symbol**, which defines the entry point of the grammar.

This section explores how each of these components is expressed using EBNF.

### Terminals

Terminals can be defined either implicitly or explicitly.

  - **Implicit terminals** are defined directly by using string values on the right-hand side of production rules.
  - **Explicit terminals** are defined by assigning a name to a value.
    In Emerge's terminology, this name is called a TOKEN, and the value can be either a string or a regular expression.

A *TOKEN* refers to a terminal with an explicit definition.
TOKENS must start with an uppercase letter (A–Z), followed by any number of uppercase letters or underscores.

Internally, implicit terminal definitions are converted into token definitions using string values.

The following predefined regular expressions can be referenced when explicitly defining a token.

  - `$DIGIT`
  - `$HEX`
  - `$LETTER`
  - `$INT`
  - `$FLOAT`
  - `$STRING`
  - `$COMMENT`

**Collisions:** If multiple tokens match the same string,
the lexer automaton may recognize more than one token in a final state.
Tokens defined by string values take precedence over those defined by regular expressions.
If multiple tokens share the same string value, lexer generation fails.

**Whitespaces** (space, tab, newline, carriage return, and all Unicode spacing and breaking characters)
are skipped by the lexer by default — you do not need to handle them in your grammar.
If you define a token that matches whitespace, only the characters matched by that token are emitted;
any remaining whitespace is still ignored.

### Non-Terminals

Non-terminal symbols are always defined and referenced in place, without needing prior declaration.

A non-terminal name must start with a lowercase letter (a–z),
followed by any number of lowercase letters or underscores.

### Productions

Production rules in a grammar are defined using the following syntax:

	non-terminal = a sequence of terminal and/or non-terminal symbols, ending with a semicolon.

#### Examples:

	stmt = "if" expr "then" stmt

##### Alternation

Multiple productions with the same non-terminal head can be combined using the alternation operator `|`:

	expr = expr "+" expr | expr "-" expr | expr "*" expr | expr "/" expr

##### Grouping

Grouping can be used with alternation to improve readability:

	expr = expr ("+" | "-" | "*" | "/") expr

##### Optional

Square brackets denote an optional sequence of symbols:

	decl = "var" ID ["=" VALUE]

##### Repetition (Kleene Star & Kleene Plus)

  - `{ }` (*Kleene Star*): zero or more occurrences
  - `{{ }}` (*Kleene Plus*): one or more occurrences

```
array = "[" {element} "]"
block = "{" {{stmt}}  "}"
```

### Start Symbol

By convention an EBNF grammar is expected to have a production rule with the special non-terminal `start`.
This is going to be considered as the start symbol of the grammar.

### Associativity and Precedence

Emerge can handle certain ambiguous grammars
if the ambiguities can be resolved using associativity and precedence rules.

In some cases, introducing controlled ambiguities can make a grammar more readable
while leading to a simpler and more efficient parser. To maintain clarity,
Emerge keeps the definition of a context-free grammar separate from associativity and precedence rules.
Unlike some systems that intertwine these concerns, Emerge ensures that
production rules retain their original semantics, just as in standard BNF.

Associativity and precedence can be assigned to either terminal symbols or production rules.
However, assigning them to terminals is preferred.

If a production rule contains a terminal symbol on its right-hand side,
it inherits the associativity and precedence of the leftmost terminal in the rule.
If a production rule does not contain any terminal symbols, associativity and precedence
can be explicitly assigned by enclosing the entire rule in angle brackets (`< >`).

Possible values for associativity are:

  - `@left` – Left-associative
  - `@right` – Right-associative
  - `@none` – No associativity

Rules listed earlier (higher in the list) have higher precedence than those listed later (lower in the list).

#### Examples

```
@left "*" "/"
@left "+" "-"
@left "!"
@left <expr = expr logop expr>
@none "=" "<>" "<" ">" "<=" ">="
@left "||" "&&"
```

## Generating A Parser

The generated parser offers three primary modes of operation, similar to the examples
[here](https://pkg.go.dev/github.com/moorara/algo/parser/lr/lookahead#pkg-examples):

  - **Tokenization and Production Extraction**: Outputs tokens and their corresponding productions.
  - **Abstract Syntax Tree (AST) Construction**: Builds an AST based on the grammar's production rules.
  - **Rule-based Evaluation and Direct Translation**: Evaluates production rules
    alongside previously computed values, enabling direct translation of the parsed input.
