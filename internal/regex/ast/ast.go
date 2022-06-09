// Package ast implements a minimal abstract syntax tree for regular expressions.
// The functions computed on each node are required by the construction alogorithm for finite state automata.
package ast

// Node is the interface for all nodes in an abstract syntax tree.
type Node interface {
	// Nullable(n) is true for a node n if and only if the subexpression represented by n has ε in its language.
	Nullable() bool

	// FirstPos(n) is the set of positions in the subtree rooted at n that
	// correspond to the first symbol of at least one string in the language of the subexpression rooted at n.
	FirstPos() []int

	// LastPos(n) is the set of positions in the subtree rooted at n that
	// correspond to the last symbol of at least one string in the language of the subexpression rooted at n.
	LastPos() []int
}

type computed struct {
	nullable bool
	firstPos []int
	lastPos  []int
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
		firstPos: []int{},
		lastPos:  []int{},
	}

	for _, expr := range n.Exprs {
		n.comp.nullable = n.comp.nullable && expr.Nullable()
	}

	for _, expr := range n.Exprs {
		n.comp.firstPos = append(n.comp.firstPos, expr.FirstPos()...)
		if !expr.Nullable() {
			break
		}
	}

	for i := len(n.Exprs) - 1; i >= 0; i-- {
		n.comp.lastPos = append(n.Exprs[i].LastPos(), n.comp.lastPos...)
		if !n.Exprs[i].Nullable() {
			break
		}
	}
}

func (n *Concat) Nullable() bool {
	n.compute()
	return n.comp.nullable
}

func (n *Concat) FirstPos() []int {
	n.compute()
	return n.comp.firstPos
}

func (n *Concat) LastPos() []int {
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
		firstPos: []int{},
		lastPos:  []int{},
	}

	for _, expr := range n.Exprs {
		n.comp.nullable = n.comp.nullable || expr.Nullable()
		n.comp.firstPos = append(n.comp.firstPos, expr.FirstPos()...)
		n.comp.lastPos = append(n.comp.lastPos, expr.LastPos()...)
	}
}

func (n *Alt) Nullable() bool {
	n.compute()
	return n.comp.nullable
}

func (n *Alt) FirstPos() []int {
	n.compute()
	return n.comp.firstPos
}

func (n *Alt) LastPos() []int {
	n.compute()
	return n.comp.lastPos
}

// Star represents a Kleene Star node.
type Star struct {
	Expr Node
}

func (n *Star) Nullable() bool {
	return true
}

func (n *Star) FirstPos() []int {
	return n.Expr.FirstPos()
}

func (n *Star) LastPos() []int {
	return n.Expr.LastPos()
}

// Empty represents an Empty string ε leaf node.
type Empty struct{}

func (n *Empty) Nullable() bool {
	return true
}

func (n *Empty) FirstPos() []int {
	return []int{}
}

func (n *Empty) LastPos() []int {
	return []int{}
}

// Char represents a Character leaf node.
type Char struct {
	Val rune
	Pos int
}

func (n *Char) Nullable() bool {
	return false
}

func (n *Char) FirstPos() []int {
	return []int{n.Pos}
}

func (n *Char) LastPos() []int {
	return []int{n.Pos}
}
