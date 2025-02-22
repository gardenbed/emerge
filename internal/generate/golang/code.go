package golang

import (
	"regexp"

	"github.com/moorara/algo/generic"
)

var idRegex = regexp.MustCompile(`^[\p{L}_][\p{L}\p{Nd}_]*$`)

var builtin = []string{
	// Keywords
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
	// Types
	"any", "bool", "byte", "comparable",
	"complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64", "rune", "string",
	"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
	// Constants
	"true", "false", "iota",
	// Zero value
	"nil",
	// Functions
	"append", "cap", "clear", "close", "complex", "copy", "delete", "imag", "len",
	"make", "max", "min", "new", "panic", "print", "println", "real", "recover",
}

// isIDValid checks if a name is a valid identifier in Go.
func isIDValid(name string) bool {
	return idRegex.MatchString(name) && !generic.AnyMatch(builtin, func(s string) bool {
		return s == name
	})
}
