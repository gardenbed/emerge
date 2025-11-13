[![Go Doc][godoc-image]][godoc-url]
[![CodeQL][codeql-image]][codeql-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][codecov-image]][codecov-url]

# Emerge

Emerge is a command-line tool for automatically generating
[lexers](https://en.wikipedia.org/wiki/Lexical_analysis) and [parsers](https://en.wikipedia.org/wiki/Parsing)
in [Go](https://go.dev) from an [EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form)
description of any arbitrary grammar (language).

For an in-depth description of the project, please see the full documentation [here](./docs/index.md).

## Quick Start

### Install

```
brew install gardenbed/brew/emerge
```

For other platforms, you can download the binary from the [latest release](https://github.com/gardenbed/emerge/releases/latest).

### Examples

Below is the context-free grammar for the JSON language in EBNF format. Save this as `json.ebnf`:

```
grammar json

NUMBER   = $FLOAT
STRING   = $STRING

start    = value;
value    = object | array | STRING | NUMBER | "true" | "false" | "null";
object   = "{" members "}" | "{" "}";
members  = members "," member | member;
member   = STRING ":" value;
array    = "[" elements "]" | "[" "]";
elements = elements "," value | value;
```

```bash
emerge json.ebnf
```


[godoc-url]: https://pkg.go.dev/github.com/gardenbed/emerge
[godoc-image]: https://pkg.go.dev/badge/github.com/gardenbed/emerge
[codeql-url]: https://github.com/gardenbed/emerge/actions/workflows/github-code-scanning/codeql
[codeql-image]: https://github.com/gardenbed/emerge/workflows/CodeQL/badge.svg
[workflow-url]: https://github.com/gardenbed/emerge/actions
[workflow-image]: https://github.com/gardenbed/emerge/workflows/Go/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/gardenbed/emerge
[goreport-image]: https://goreportcard.com/badge/github.com/gardenbed/emerge
[codecov-url]: https://codecov.io/gh/gardenbed/emerge
[codecov-image]: https://codecov.io/gh/gardenbed/emerge/branch/main/graph/badge.svg
