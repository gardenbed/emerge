// Package ast implements a minimal abstract syntax tree for regular expressions.
// The functions computed on each node are required by the construction alogorithm for finite state automata.
//
// It also provides a combinator parser for parsing regular expression into an abstract syntax tree.
package ast

import (
	"errors"
	"fmt"
	"sort"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/list"
	"github.com/moorara/algo/parser/combinator"

	"github.com/gardenbed/emerge/internal/char"
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
	Root      Node
	lastPos   Pos
	posToChar map[Pos]char.Range
	charToPos map[char.Range]Poses
	follows   map[Pos]Poses
}

func Parse(regex string) (*AST, error) {
	m := new(mappers)
	p := parser.New(m)

	out, ok := p.Parse(regex)
	if !ok {
		return nil, fmt.Errorf("invalid regular expression: %s", regex)
	}

	if m.errors != nil {
		return nil, m.errors
	}

	// Concat a unique right end-marker to the regular expression root node.
	// This is required for the construction of a DFA directly from a regular expression.
	root := &Concat{
		Exprs: []Node{
			out.Result.Val.(Node),
			&Char{Lo: endMarker, Hi: endMarker},
		},
	}

	a := &AST{
		Root:      root,
		lastPos:   0,
		posToChar: map[Pos]char.Range{},
		charToPos: map[char.Range]Poses{},
		follows:   map[Pos]Poses{},
	}

	// Assign one-based positions to Char nodes.
	a.backfillCharPos(a.Root)

	// Preprocessing: compute followpos function.
	a.computeFollows(a.Root)
	for _, poses := range a.follows {
		sort.Sort(poses)
	}

	return a, nil
}

// backfillCharPos backfills Pos for all Char nodes in the abstract syntaxt tree from left to right.
// These positions are one-based and used for directly converting a regular expression to a DFA.
// They are semantically different from the zero-based positions set by parsers (in the mappers).
//
// It also creates a map of positions to characters and vice versa.
// These mappings are used when constructing the DFA from the AST.
func (a *AST) backfillCharPos(n Node) {
	switch v := n.(type) {
	case *Concat:
		for _, e := range v.Exprs {
			a.backfillCharPos(e)
		}

	case *Alt:
		for _, e := range v.Exprs {
			a.backfillCharPos(e)
		}

	case *Star:
		a.backfillCharPos(v.Expr)

	case *Char:
		// Assign the next position to the current Char.
		a.lastPos++
		v.Pos = a.lastPos

		// Update the pos-to-char and char-to-pos mappings.
		c := char.Range{v.Lo, v.Hi}
		a.posToChar[v.Pos] = c
		a.charToPos[c] = append(a.charToPos[c], v.Pos)
	}
}

// There are only two ways that a position of a regular expression can be made to follow another.
//
//  1. If n is a concat node with left child n1 and right child n2,
//     then for every position i in lastPos(n1), all positions in firstPos(n2) are in followPos(i).
//
//  2. If n is a star node, and i is a position in lastPos(n1),
//     then all positions in firstPos(n1) are in followPos(i).
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

// followPos for a position p, is the set of positions q in the entire syntax tree such that there is some string
// x = a₁a₂...aₙ in L((r)µ) such that for some i, there is a way to explain the membership of x in L((r)µ) by
// matching aᵢ to position p of the syntax tree and aᵢ₊₁ to position q.
func (a *AST) followPos(p Pos) Poses {
	return a.follows[p]
}

// ToDFA converts the abstract syntax tree for a regular expression to a DFA.
//
// For more details, see Compilers: Principles, Techniques, and Tools (2nd Edition).
func (a *AST) ToDFA() *automata.DFA {
	end := char.Range{endMarker, endMarker}

	b := automata.NewDFABuilder().SetStart(0)

	Dstates := list.NewSoftQueue(func(p, q Poses) bool {
		return p.Equal(q)
	})

	// Initialize Dstates to contain only the firstpos(n0), where n0 is the root of syntax tree for (r)µ
	Dstates.Enqueue(a.Root.firstPos())

	for S, i := Dstates.Dequeue(); i >= 0; S, i = Dstates.Dequeue() {
		for c := range a.charToPos { // for each input symbol c
			if c != end {
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

				lo, hi := automata.Symbol(c[0]), automata.Symbol(c[1])
				b.AddTransition(automata.State(i), lo, hi, automata.State(j))
			}
		}
	}

	final := []automata.State{}

	for i, S := range Dstates.Values() {
		for _, f := range a.charToPos[end] {
			if S.Contains(f) {
				final = append(final, automata.State(i))
				break // The accepting states of D are all those sets of positions that include the position of the end-marker
			}
		}
	}

	return b.SetFinal(final).Build().Minimize()
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

	// equal checks if two nodes are equal.
	equal(Node) bool
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
		nullable: true,
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

func (n *Concat) equal(rhs Node) bool {
	nn, ok := rhs.(*Concat)
	if !ok {
		return false
	}

	if len(n.Exprs) != len(nn.Exprs) {
		return false
	}

	for i := 0; i < len(n.Exprs); i++ {
		if !n.Exprs[i].equal(nn.Exprs[i]) {
			return false
		}
	}

	return true
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

func (n *Alt) equal(rhs Node) bool {
	nn, ok := rhs.(*Alt)
	if !ok {
		return false
	}

	if len(n.Exprs) != len(nn.Exprs) {
		return false
	}

	for i := 0; i < len(n.Exprs); i++ {
		if !n.Exprs[i].equal(nn.Exprs[i]) {
			return false
		}
	}

	return true
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

func (n *Star) equal(rhs Node) bool {
	nn, ok := rhs.(*Star)
	if !ok {
		return false
	}

	return n.Expr.equal(nn.Expr)
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

func (n *Empty) equal(rhs Node) bool {
	_, ok := rhs.(*Empty)

	return ok
}

// Char represents a Character leaf node.
// It represents an inclusive range of characters.
type Char struct {
	Lo  rune
	Hi  rune
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

func (n *Char) equal(rhs Node) bool {
	nn, ok := rhs.(*Char)
	if !ok {
		return false
	}

	return n.Lo == nn.Lo && n.Hi == nn.Hi && n.Pos == nn.Pos
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

func (p Poses) Equal(q Poses) bool {
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

//==================================================< MAPPERS >==================================================

// Anchor is the type for a regular expression anchor.
type Anchor int

const (
	StartOfString Anchor = iota + 1
	EndOfString

	bagKeyCharRanges     combinator.BagKey = "char_ranges"
	bagKeyLazyQuantifier combinator.BagKey = "lazy_quantifier"
	BagKeyStartOfString  combinator.BagKey = "start_of_string"
)

// mappers implements the parser.Mappers interface.
type mappers struct {
	errors error
}

func (m *mappers) ToAnyChar(r combinator.Result) (combinator.Result, bool) {
	node, _ := charRangesToNode(false, char.Classes["UNICODE"])

	return combinator.Result{
		Val: node,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToSingleChar(r combinator.Result) (combinator.Result, bool) {
	c := r.Val.(rune)
	node, ranges := charRangesToNode(false, char.RangeList{{c, c}})

	return combinator.Result{
		Val: node,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, true
}

func (m *mappers) ToCharClass(r combinator.Result) (combinator.Result, bool) {
	class := r.Val.(string)

	var node Node
	var ranges char.RangeList

	switch class {
	case `\s`:
		node, ranges = charRangesToNode(false, char.Classes[`\s`])
	case `\S`:
		node, ranges = charRangesToNode(true, char.Classes[`\s`])
	case `\d`:
		node, ranges = charRangesToNode(false, char.Classes[`\d`])
	case `\D`:
		node, ranges = charRangesToNode(true, char.Classes[`\d`])
	case `\w`:
		node, ranges = charRangesToNode(false, char.Classes[`\w`])
	case `\W`:
		node, ranges = charRangesToNode(true, char.Classes[`\w`])
	default:
		return combinator.Result{}, false
	}

	return combinator.Result{
		Val: node,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, true
}

func (m *mappers) ToASCIICharClass(r combinator.Result) (combinator.Result, bool) {
	class := r.Val.(string)

	ranges, ok := char.Classes[class]
	if !ok {
		return combinator.Result{}, false
	}

	node, ranges := charRangesToNode(false, ranges)

	return combinator.Result{
		Val: node,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, true
}

func (m *mappers) ToUnicodeCategory(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToUnicodeCharClass(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r2, _ := r.Get(2)

	prop := r0.Val.(string)
	class := r2.Val.(string)

	ranges, ok := char.Classes[class]
	if !ok {
		return combinator.Result{}, false
	}

	node, ranges := charRangesToNode(prop == `\P`, ranges)

	return combinator.Result{
		Val: node,
		Pos: r.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, true
}

func (m *mappers) ToRepOp(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToUpperBound(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	var num *int
	if v, ok := r1.Val.(int); ok {
		num = &v
	}

	return combinator.Result{
		Val: num,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRange(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	// The upper bound is same as the lower bound if no upper bound is specified (default)
	low := r1.Val.(int)
	up := &low

	// If an upper bound is specified, it can be either bounded or unbounded
	if v, ok := r2.Val.(*int); ok {
		up = v
	}

	if up != nil && low > *up {
		// The input syntax is correct while its semantic is incorrect
		// We continue parsing the rest of input to find more errors
		m.errors = errors.Join(m.errors, fmt.Errorf("invalid repetition range {%d,%d}", low, *up))
	}

	return combinator.Result{
		Val: tuple[int, *int]{
			p: low,
			q: up,
		},
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRepetition(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToQuantifier(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	// Check whether or not the lazy modifier is present
	_, lazy := r1.Val.(rune)

	return combinator.Result{
		Val: tuple[any, bool]{
			p: r0.Val,
			q: lazy,
		},
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToCharInRange(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToCharRange(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r2, _ := r.Get(2)

	lo, hi := r0.Val.(rune), r2.Val.(rune)

	if lo > hi {
		m.errors = errors.Join(m.errors, fmt.Errorf("invalid character range %s-%s", string(lo), string(hi)))

		// The input syntax is correct while its semantic is incorrect.
		// We continue parsing the rest of input to find more errors.
		return combinator.Result{Pos: r0.Pos}, true
	}

	node, ranges := charRangesToNode(false, char.RangeList{{lo, hi}})

	return combinator.Result{
		Val: node,
		Pos: r0.Pos,
		Bag: combinator.Bag{
			bagKeyCharRanges: ranges,
		},
	}, true
}

func (m *mappers) ToCharGroupItem(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToCharGroup(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r2, _ := r.Get(2)

	// Check whether or not the negation modifier is present
	_, neg := r1.Val.(rune)

	// Collect all character ranges from the character group items
	var all char.RangeList
	for _, r := range r2.Val.(combinator.List) {
		if ranges, ok := r.Bag[bagKeyCharRanges].(char.RangeList); ok {
			all = append(all, ranges...)
		}
	}

	node, _ := charRangesToNode(neg, all.Dedup())

	return combinator.Result{
		Val: node,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToMatchItem(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToMatch(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	node := r0.Val.(Node)
	var bag combinator.Bag

	if t, ok := r1.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, t.p)
		if lazy := t.q; lazy {
			bag = combinator.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return combinator.Result{
		Val: node,
		Pos: r0.Pos,
		Bag: bag,
	}, true
}

func (m *mappers) ToGroup(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)
	r3, _ := r.Get(3)

	node := r1.Val.(Node)
	var bag combinator.Bag

	if t, ok := r3.Val.(tuple[any, bool]); ok {
		node = quantifyNode(node, t.p)
		if lazy := t.q; lazy {
			bag = combinator.Bag{
				bagKeyLazyQuantifier: true,
			}
		}
	}

	return combinator.Result{
		Val: node,
		Pos: r0.Pos,
		Bag: bag,
	}, true
}

func (m *mappers) ToAnchor(r combinator.Result) (combinator.Result, bool) {
	c := r.Val.(rune)

	var anchor Anchor
	switch c {
	case '$': // end-of-string
		anchor = EndOfString
	}

	return combinator.Result{
		Val: anchor,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToSubexprItem(r combinator.Result) (combinator.Result, bool) {
	// Passing the result up the parsing chain
	return r, true
}

func (m *mappers) ToSubexpr(r combinator.Result) (combinator.Result, bool) {
	items := r.Val.(combinator.List)

	concat := new(Concat)
	for _, r := range items {
		// TODO: Anchor result value is not a node
		if n, ok := r.Val.(Node); ok {
			concat.Exprs = append(concat.Exprs, n)
		}
	}

	return combinator.Result{
		Val: concat,
		Pos: r.Pos,
	}, true
}

func (m *mappers) ToExpr(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	node := r0.Val.(Node)

	if _, ok := r1.Val.(combinator.List); ok {
		r11, _ := r1.Get(1)
		expr := r11.Val.(Node)
		node = &Alt{
			Exprs: []Node{node, expr},
		}
	}

	return combinator.Result{
		Val: node,
		Pos: r0.Pos,
	}, true
}

func (m *mappers) ToRegex(r combinator.Result) (combinator.Result, bool) {
	r0, _ := r.Get(0)
	r1, _ := r.Get(1)

	pos := r1.Pos
	var bag combinator.Bag

	// Check whether or not the start-of-string is present
	if _, sos := r0.Val.(rune); sos {
		pos = r0.Pos
		bag = combinator.Bag{
			BagKeyStartOfString: true,
		}
	}

	expr := r1.Val.(Node)

	return combinator.Result{
		Val: expr,
		Pos: pos,
		Bag: bag,
	}, true
}

//==================================================< HELPERS >==================================================

type tuple[P, Q any] struct {
	p P
	q Q
}

// charRangesToNode converts a list of character ranges into a Node.
// The resulting Node is either a Char node or an Alt node.
//
// If neg is true, the resulting Alt node will represent the negation of the given ranges.
// It also returns the list of rune ranges represented by the resulting Alt node.
func charRangesToNode(neg bool, ranges char.RangeList) (Node, char.RangeList) {
	if neg {
		ranges = char.Classes["UNICODE"].Exclude(ranges)
	}

	if len(ranges) == 1 {
		lo, hi := ranges[0][0], ranges[0][1]
		return &Char{Lo: lo, Hi: hi /* Pos will be set later */}, ranges
	}

	alt := new(Alt)
	for _, r := range ranges {
		lo, hi := r[0], r[1]
		alt.Exprs = append(alt.Exprs, &Char{Lo: lo, Hi: hi /* Pos will be set later */})
	}

	return alt, ranges
}

// quantifyNFA applies the quantifier q to the given node n and returns the resulting node.
func quantifyNode(n Node, q any) Node {
	var node Node

	switch rep := q.(type) {
	// Simple repetition
	case rune:
		switch rep {
		case '?':
			node = &Alt{
				Exprs: []Node{
					&Empty{},
					cloneNode(n),
				},
			}

		case '*':
			node = &Star{
				Expr: cloneNode(n),
			}

		case '+':
			node = &Concat{
				Exprs: []Node{
					cloneNode(n),
					&Star{
						Expr: cloneNode(n),
					},
				},
			}
		}

	// Range repetition
	case tuple[int, *int]:
		low, up := rep.p, rep.q
		concat := new(Concat)

		for i := 0; i < low; i++ {
			concat.Exprs = append(concat.Exprs, cloneNode(n))
		}

		if up == nil {
			concat.Exprs = append(concat.Exprs, &Star{
				Expr: cloneNode(n),
			})
		} else {
			for i := 0; i < *up-low; i++ {
				concat.Exprs = append(concat.Exprs, &Alt{
					Exprs: []Node{
						&Empty{},
						cloneNode(n),
					},
				})
			}
		}

		node = concat
	}

	return node
}

// Cloning is required since each instance of Char will have a distinct value for Pos field.
func cloneNode(n Node) Node {
	switch v := n.(type) {
	case *Concat:
		concat := new(Concat)
		for _, e := range v.Exprs {
			concat.Exprs = append(concat.Exprs, cloneNode(e))
		}
		return concat

	case *Alt:
		alt := new(Alt)
		for _, e := range v.Exprs {
			alt.Exprs = append(alt.Exprs, cloneNode(e))
		}
		return alt

	case *Star:
		return &Star{
			Expr: cloneNode(v.Expr),
		}

	case *Empty:
		return new(Empty)

	case *Char:
		return &Char{
			Lo:  v.Lo,
			Hi:  v.Hi,
			Pos: v.Pos,
		}

	default:
		return nil
	}
}
