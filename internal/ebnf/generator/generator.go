package generator

import auto "github.com/moorara/algo/automata"

func StringToDFA(s string) *auto.DFA {
	start := auto.State(0)
	dfa := auto.NewDFA(start, nil)

	curr, next := start, start+1
	for _, r := range s {
		dfa.Add(curr, auto.Symbol(r), next)
		curr, next = next, next+1
	}

	dfa.Final = auto.States{curr}

	return dfa
}
