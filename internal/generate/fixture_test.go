package generate

import (
	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"

	"github.com/gardenbed/emerge/internal/ebnf/parser/spec"
)

var grammars = []*grammar.CFG{}

var precedences = []lr.PrecedenceLevels{}

func getDefinitions() []*spec.TerminalDef {
	return []*spec.TerminalDef{}
}
