// Package ast implements a minimal abstract syntax tree for regular expressions.
// The functions computed on each node are required by the construction alogorithm for finite state automata.
//
// It also provides a combinator parser for parsing regular expression into an abstract syntax tree.
package ast

import (
	"errors"
	"sort"

	auto "github.com/moorara/algo/automata"
	"github.com/moorara/algo/list"

	comb "github.com/gardenbed/emerge/internal/combinator"
	"github.com/gardenbed/emerge/internal/regex/parser"
)

// End marker is a special character not used anywhere in the alphabet for regular expressions.
// This special unique character is taken from a Private Use Area (PUA) in Unicode.
//
// By concatenating a unique right end-marker µ to a regular expression r,
// we give the accepting state for r a transition on µ, making it an important state of the NFA for (r)µ.
// This is useful for directly constructing a DFA for a regular expression.
//
// For more details, see Compilers: Principles, Techniques, and Tools (2nd Edition).
const endMarker rune = 0xEEEE

// AST is the abstract syntax tree for a regular expression.
type AST struct {
	Root Node

	lastPos   Pos
	posToChar map[Pos]rune
	charToPos map[rune]Poses
	follows   map[Pos]Poses
}

func Parse(in comb.Input) (*AST, error) {
	m := new(mappers)
	p := parser.New(m)

	out, ok := p.Parse(in)
	if !ok {
		return nil, errors.New("invalid regular expression")
	}

	if m.errors != nil {
		return nil, m.errors
	}

	// Concat a unique right end-marker to the regular expression root node.
	// This is required for the construction of a DFA directly from a regular expression.
	root := &Concat{
		Exprs: []Node{
			out.Result.Val.(Node),
			&Char{Val: endMarker},
		},
	}

	a := &AST{
		Root:      root,
		lastPos:   0,
		posToChar: map[Pos]rune{},
		charToPos: map[rune]Poses{},
		follows:   map[Pos]Poses{},
	}

	// Preprocessing: assign positions to Char nodes and index them
	a.indexChars(a.Root)

	// Preprocessing: compute followpos function
	a.computeFollows(a.Root)
	for _, list := range a.follows {
		sort.Sort(list)
	}

	return a, nil
}

// indexChars backfills Pos for all Char nodes in the abstract syntaxt tree from left to right.
// These positions are one-based and used for directly converting a regular expression to a DFA.
// They are semantically different from the zero-based positions set by parsers (in the mappers).
//
// It also creates a map of positions to characters and vice versa.
func (a *AST) indexChars(n Node) {
	switch v := n.(type) {
	case *Concat:
		for _, e := range v.Exprs {
			a.indexChars(e)
		}

	case *Alt:
		for _, e := range v.Exprs {
			a.indexChars(e)
		}

	case *Star:
		a.indexChars(v.Expr)

	case *Char:
		a.lastPos++
		p, c := a.lastPos, v.Val
		v.Pos = p
		a.posToChar[p] = c
		a.charToPos[c] = append(a.charToPos[c], p)
	}
}

// There are only two ways that a position of a regular expression can be made to follow another.
//
// 1. If n is a concat node with left child n1 and right child n2,
// then for every position i in lastPos(n1), all positions in firstPos(n2) are in followPos(i).
//
// 2. If n is a star node, and i is a position in lastPos(n1), then all positions in firstPos(n1) are in followPos(i).
func (a *AST) computeFollows(n Node) {
	switch v := n.(type) {
	case *Concat:
		for i := 0; i < len(v.Exprs)-1; i++ {
			for _, p := range v.Exprs[i].lastPos() {
				a.follows[p] = append(a.follows[p], v.Exprs[i+1].firstPos()...)
			}
		}

		for _, e := range v.Exprs {
			a.computeFollows(e)
		}

	case *Alt:
		for _, e := range v.Exprs {
			a.computeFollows(e)
		}

	case *Star:
		for _, p := range v.Expr.lastPos() {
			a.follows[p] = append(a.follows[p], v.Expr.firstPos()...)
		}

		a.computeFollows(v.Expr)
	}
}

// for a position p, is the set of positions q in the entire syntax tree
// such that there is some string x = a1a2...an in L((r)µ)
// such that for some i, there is a way to explain the membership of x in L((r)µ)
// by matching ai to position p of the syntax tree and ai+1 to position q.
func (a *AST) followPos(p Pos) Poses {
	return a.follows[p]
}

func (a *AST) ToDFA() *auto.DFA {
	dfa := auto.NewDFA(0, nil)
	Dstates := list.NewSoftQueue[Poses](func(p, q Poses) bool {
		return p.Equals(q)
	})

	// Initialize Dstates to contain only the firstpos(n0), where n0 is the root of syntax tree for (r)µ
	Dstates.Enqueue(a.Root.firstPos())

	for S, i := Dstates.Dequeue(); i >= 0; S, i = Dstates.Dequeue() {
		for c := range a.charToPos { // for each input symbol c
			if c != endMarker {
				// Let U be the union of followpos(p) for all p in S that correspond to c
				U := Poses{}
				for _, p := range S {
					if a.posToChar[p] == c {
						U = U.Union(a.followPos(p))
					}
				}

				// If U is not in Dstates, add U to Dstates
				j := Dstates.Contains(U)
				if j == -1 {
					j = Dstates.Enqueue(U)
				}

				dfa.Add(auto.State(i), auto.Symbol(c), auto.State(j))
			}
		}
	}

	dfa.Start = auto.State(0)
	dfa.Final = auto.States{}

	for i, S := range Dstates.Values() {
		for _, f := range a.charToPos[endMarker] {
			if S.Contains(f) {
				dfa.Final = append(dfa.Final, auto.State(i))
				break // The accepting states of D are all those sets of positions that include the position of the end-marker
			}
		}
	}

	// Minimize the no. of states
	dfa = dfa.Minimize()

	return dfa
}

// Node is the interface for all nodes in an abstract syntax tree.
type Node interface {
	// nullable(n) is true for a node n if and only if the subexpression represented by n has ε in its language.
	nullable() bool

	// firstPos(n) is the set of positions in the subtree rooted at n that
	// correspond to the first symbol of at least one string in the language of the subexpression rooted at n.
	firstPos() Poses

	// lastPos(n) is the set of positions in the subtree rooted at n that
	// correspond to the last symbol of at least one string in the language of the subexpression rooted at n.
	lastPos() Poses
}

type computed struct {
	nullable bool
	firstPos Poses
	lastPos  Poses
}

// Concat represents a concatenation node.
type Concat struct {
	Exprs []Node
	comp  *computed
}

func (n *Concat) compute() {
	if n.comp != nil {
		return
	}

	n.comp = &computed{
		nullable: false,
		firstPos: Poses{},
		lastPos:  Poses{},
	}

	for _, expr := range n.Exprs {
		n.comp.nullable = n.comp.nullable && expr.nullable()
	}

	for _, expr := range n.Exprs {
		n.comp.firstPos = append(n.comp.firstPos, expr.firstPos()...)
		if !expr.nullable() {
			break
		}
	}

	for i := len(n.Exprs) - 1; i >= 0; i-- {
		n.comp.lastPos = append(n.Exprs[i].lastPos(), n.comp.lastPos...)
		if !n.Exprs[i].nullable() {
			break
		}
	}
}

func (n *Concat) nullable() bool {
	n.compute()
	return n.comp.nullable
}

func (n *Concat) firstPos() Poses {
	n.compute()
	return n.comp.firstPos
}

func (n *Concat) lastPos() Poses {
	n.compute()
	return n.comp.lastPos
}

// Alt represents an Alternation node.
type Alt struct {
	Exprs []Node
	comp  *computed
}

func (n *Alt) compute() {
	if n.comp != nil {
		return
	}

	n.comp = &computed{
		nullable: false,
		firstPos: Poses{},
		lastPos:  Poses{},
	}

	for _, expr := range n.Exprs {
		n.comp.nullable = n.comp.nullable || expr.nullable()
		n.comp.firstPos = append(n.comp.firstPos, expr.firstPos()...)
		n.comp.lastPos = append(n.comp.lastPos, expr.lastPos()...)
	}
}

func (n *Alt) nullable() bool {
	n.compute()
	return n.comp.nullable
}

func (n *Alt) firstPos() Poses {
	n.compute()
	return n.comp.firstPos
}

func (n *Alt) lastPos() Poses {
	n.compute()
	return n.comp.lastPos
}

// Star represents a Kleene Star node.
type Star struct {
	Expr Node
}

func (n *Star) nullable() bool {
	return true
}

func (n *Star) firstPos() Poses {
	return n.Expr.firstPos()
}

func (n *Star) lastPos() Poses {
	return n.Expr.lastPos()
}

// Empty represents an Empty string ε leaf node.
type Empty struct{}

func (n *Empty) nullable() bool {
	return true
}

func (n *Empty) firstPos() Poses {
	return Poses{}
}

func (n *Empty) lastPos() Poses {
	return Poses{}
}

// Char represents a Character leaf node.
type Char struct {
	Val rune
	Pos Pos // one-based
}

func (n *Char) nullable() bool {
	return false
}

func (n *Char) firstPos() Poses {
	return Poses{n.Pos}
}

func (n *Char) lastPos() Poses {
	return Poses{n.Pos}
}

// Pos is the type for positions.
type Pos int

// Poses is the type for a set of positions.
type Poses []Pos

func (p Poses) Len() int {
	return len(p)
}

func (p Poses) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p Poses) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Poses) Contains(q Pos) bool {
	for _, pos := range p {
		if pos == q {
			return true
		}
	}

	return false
}

func (p Poses) Equals(q Poses) bool {
	for _, pos := range p {
		if !q.Contains(pos) {
			return false
		}
	}

	for _, pos := range q {
		if !p.Contains(pos) {
			return false
		}
	}

	return true
}

func (p Poses) Union(q Poses) Poses {
	u := make(Poses, len(p))
	copy(u, p)

	for _, pos := range q {
		if !u.Contains(pos) {
			u = append(u, pos)
		}
	}

	return u
}
