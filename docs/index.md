# Emerge

## What

Emerge is a command-line tool for automatically generating
[lexers](https://en.wikipedia.org/wiki/Lexical_analysis) and [parsers](https://en.wikipedia.org/wiki/Parsing)
in [Go](https://go.dev) from an [EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form)
description of any arbitrary grammar (language).

Since Emerge is agnostic about a given input grammar (language),
building other's modules of the compiler's front-end
(semantic analyzer, intermediate code generator, and machine-independent code optimizer)
are beyond the scope of this project.

For building the back-end for a compiler's front-end you can tap into
out-of-the-shell projects such as [LLVM](https://www.llvm.org).

## Why

Currently, Go does not provide a first-class support for building compilers written completely in Go.
When it comes to lexical analysis, you can either use the built-in [`text/scanner`](https://pkg.go.dev/text/scanner)
package for simple use-cases or build the scanner by hand for more sophisticated use-cases.
For the syntactic analysis part, you can use [`goyacc`](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc)
for auto-generating a parser from a [Yacc](https://en.wikipedia.org/wiki/Yacc) file.
Since, `goyacc` is just a port of [yacc](https://www.tuhs.org/cgi-bin/utree.pl?file=V6/usr/source/yacc)
For Go, it does not play very well with the language. In particular,

  - The Yacc domain-specific language (DSL) is more than 50 years old,
    and as a result it is very hard-to-understand and hard-to-maintain.
  - You need to remember many conventions to integrate it with Go in a working manner.
  - When things go wrong, debugging the Yacc language and the generated parser takes intuition and ad-hoc solutions.
  - For larger projecst and the projects that need collaboration, Yacc does not offer an easy modular approach.

## How

Emerge is a response to all these challenges by reimplementing the underlying concepts and algorithms,
that are well-established for a long time in the field of compiler design,
in Go with having the requirements and best practices of today's world in mind.

## Similar Projects

Here is a summary of some of the other similar projects.

  - [Lex](https://minnie.tuhs.org/cgi-bin/utree.pl?file=4BSD/usr/src/cmd/lex)
    is a program written in C that generates lexical analyzers (a.k.a. scanners).
  - [Flex](https://github.com/westes/flex) is the fast lexical analyzer generator.
    It is a free and open-source alternative to lex also written in C.
  - [Yacc](https://www.tuhs.org/cgi-bin/utree.pl?file=V6/usr/source/yacc)
    is program written in C that generates a parser from a notation similar to BNF.
  - [Bison](https://www.gnu.org/software/bison/) is a [GNU](https://en.wikipedia.org/wiki/GNU_Project)
    distribution of Yacc that takes an annotated context-free grammar and generates a deterministic parser.
  - [Goyacc](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc) is a version of Yacc for Go.
    It is written in Go and generates parsers in Go.
  - [Participle](https://github.com/alecthomas/participle) is simple Go library for building parsers from grammars.
    A grammar is an annotated Go structure used for defining both the grammar and the abstract syntax tree.

## Table of Contents

  - [Theory of Languages and Compilers](./1-theory.md)
      - [Theory of Lexers](./2-lexer_theory.md)
      - [Theory of Parser](./3-parser_theory.md)
  - [Definitions](./4-definitions.md)
  - [Design](./5-design.md)

## References

  - [Compilers: Principles, Techniques, and Tools, 2nd Edition](https://www.pearson.com/us/higher-education/program/Aho-Compilers-Principles-Techniques-and-Tools-2nd-Edition/PGM167067.html)
  - [Parsing Techniques: A Practical Guide, Second Edition](https://link.springer.com/book/10.1007/978-0-387-68954-8)
