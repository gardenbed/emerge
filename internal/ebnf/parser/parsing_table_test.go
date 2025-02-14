package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/grammar"
	"github.com/moorara/algo/parser/lr"
	"github.com/moorara/algo/parser/lr/lookahead"
)

func productionIndex(p *grammar.Production) int {
	for i, prod := range productions {
		if prod.Equal(p) {
			return i
		}
	}

	return -1
}

func TestParsingTable(t *testing.T) {
	T, err := lookahead.BuildParsingTable(G, precedences)
	assert.NoError(t, err)

	t.Run("ACTION", func(t *testing.T) {
		for _, s := range T.States {
			for _, a := range T.Terminals {
				if action, err := T.ACTION(s, a); err == nil {
					actionType, param, err := ACTION(int(s), a)

					assert.NoError(t, err)
					assert.Equal(t, action.Type, actionType)

					switch actionType {
					case lr.SHIFT:
						assert.Equal(t, int(action.State), param)
					case lr.REDUCE:
						assert.Equal(t, productionIndex(action.Production), param)
					case lr.ACCEPT:
						assert.Zero(t, param)
					}
				}
			}
		}
	})

	t.Run("GOTO", func(t *testing.T) {
		for _, s := range T.States {
			for _, A := range T.NonTerminals {
				if next, err := T.GOTO(s, A); err == nil {
					assert.Equal(t, int(next), GOTO(int(s), A))
				}
			}
		}
	})
}
