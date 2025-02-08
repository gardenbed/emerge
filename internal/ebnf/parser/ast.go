package parser

import (
	"bytes"
	"fmt"

	"github.com/moorara/algo/dot"
	"github.com/moorara/algo/generic"
	"github.com/moorara/algo/lexer"
	"github.com/moorara/algo/parser/lr"
)

// Node represents a node in an EBNF parse tree, a.k.a. abstract syntax tree (AST).
type Node interface {
	fmt.Stringer
	generic.Equaler[Node]

	// Pos returns the leftmost position in the input string that a node represent.
	Pos() *lexer.Position
}

// InternalNode represents a non-leaf node in an EBNF abstract syntax tree (AST).
// It contains at least one child node.
type InternalNode interface {
	Node
	Children() []Node
}

// LeafNode represents a leaf node in an EBNF abstract syntax tree (AST).
type LeafNode interface {
	Node
	leaf()
}

// Traverse performs a depth-first traversal of an EBNF abstract syntax tree (AST),
// starting from the given root node.
// It visits each node according to the specified traversal order
// and passes each node to the provided visit function.
// If the visit function returns false, the traversal is stopped early.
//
// Valid traversal orders for an AST are VLR, VRL, LRV, and RLV.
func Traverse(n Node, order generic.TraverseOrder, visit generic.VisitFunc1[Node]) bool {
	if leaf, ok := n.(LeafNode); ok {
		return visit(leaf)
	}

	in, ok := n.(InternalNode)
	if !ok {
		return false
	}

	children := in.Children()

	switch order {
	case generic.VLR:
		res := visit(in)
		for i := range len(children) {
			res = res && Traverse(children[i], order, visit)
		}
		return res

	case generic.VRL:
		res := visit(in)
		for i := len(children) - 1; i >= 0; i-- {
			res = res && Traverse(children[i], order, visit)
		}
		return res

	case generic.LRV:
		res := true
		for i := range len(children) {
			res = res && Traverse(children[i], order, visit)
		}
		return res && visit(in)

	case generic.RLV:
		res := true
		for i := len(children) - 1; i >= 0; i-- {
			res = res && Traverse(children[i], order, visit)
		}
		return res && visit(in)

	default:
		return false
	}
}

// Grammar represents an EBNF grammar and serves as the root of an abstract syntax tree (AST).
// This node corresponds to the `grammar → name {decl}` production rule.
type Grammar struct {
	Name     string
	Decls    []Decl
	Position *lexer.Position
}

func (n *Grammar) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "Grammar::%s", n.Name)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *Grammar) Equal(rhs Node) bool {
	nn, ok := rhs.(*Grammar)
	if !ok {
		return false
	}

	if len(n.Decls) != len(nn.Decls) {
		return false
	}

	for i := range len(n.Decls) {
		if !n.Decls[i].Equal(nn.Decls[i]) {
			return false
		}
	}

	return equalPositions(n.Position, nn.Position)
}

func (n *Grammar) Pos() *lexer.Position {
	return n.Position
}

func (n *Grammar) Children() []Node {
	nodes := make([]Node, len(n.Decls))
	for i, decl := range n.Decls {
		nodes[i] = decl
	}

	return nodes
}

func (n *Grammar) DOT() string {
	// Create a map of node --> id
	var id int
	nodeID := map[Node]int{}
	Traverse(n, generic.VLR, func(n Node) bool {
		id++
		nodeID[n] = id
		return true
	})

	graph := dot.NewGraph(true, true, false, "AST", "", "", "", "")

	Traverse(n, generic.VLR, func(n Node) bool {
		name := fmt.Sprintf("%d", nodeID[n])

		switch n := n.(type) {
		case *Grammar:
			label := fmt.Sprintf("Grammar::%s", n.Name)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorGold, dot.StyleFilled, dot.ShapeSquare, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *StringTokenDecl:
			label := fmt.Sprintf("TokenDecl::%s", n.Name)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorSkyBlue, dot.StyleFilled, dot.ShapeBox, "", ""))

		case *RegexTokenDecl:
			label := fmt.Sprintf("TokenDecl::%s", n.Name)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorSkyBlue, dot.StyleFilled, dot.ShapeBox, "", ""))

		case *PrecedenceDecl:
			label := fmt.Sprintf("PrecedenceDecl::%s", n.Associativity)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorBurlyWood, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *TerminalHandle:
			label := fmt.Sprintf("TerminalHandle::%s", n.Terminal)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorBurlyWood, dot.StyleFilled, dot.ShapeBox, "", ""))

		case *ProductionHandle:
			label := fmt.Sprintf("ProductionHandle::%s →", n.LHS)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorBurlyWood, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *RuleDecl:
			label := fmt.Sprintf("RuleDecl::%s →", n.LHS)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorLightPink, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *ConcatRHS:
			graph.AddNode(dot.NewNode(name, "", "CONCAT", dot.ColorLavender, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *AltRHS:
			graph.AddNode(dot.NewNode(name, "", "ALT", dot.ColorLavender, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *OptRHS:
			graph.AddNode(dot.NewNode(name, "", "ZERO OR ONE", dot.ColorLavender, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *StarRHS:
			graph.AddNode(dot.NewNode(name, "", "ZERO OR MORE", dot.ColorLavender, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *PlusRHS:
			graph.AddNode(dot.NewNode(name, "", "ONE OR MORE", dot.ColorLavender, dot.StyleFilled, dot.ShapeBox, "", ""))

			for _, m := range n.Children() {
				from := fmt.Sprintf("%s", name)
				to := fmt.Sprintf("%d", nodeID[m])
				graph.AddEdge(dot.NewEdge(from, to, dot.EdgeTypeDirected, "", "", "", "", "", ""))
			}

		case *NonTerminalRHS:
			label := fmt.Sprintf("NonTerminal::%s", n.NonTerminal)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorTurquoise, dot.StyleFilled, dot.ShapeBox, "", ""))

		case *TerminalRHS:
			label := fmt.Sprintf("Terminal::%s", n.Terminal)
			graph.AddNode(dot.NewNode(name, "", label, dot.ColorSpringGreen, dot.StyleFilled, dot.ShapeOval, "", ""))

		case *EmptyRHS:
			graph.AddNode(dot.NewNode(name, "", "ε", dot.ColorViolet, dot.StyleFilled, dot.ShapeCircle, "", ""))

		}

		return true
	})

	return graph.DOT()
}

// Decl represents a declaration in an EBNF grammar.
// This node corresponds to `decl → token | directive | rule` production rule.
type Decl interface {
	Node
	decl()
}

// StringTokenDecl represents a token declaration with a string value in an EBNF grammar.
// This node corresponds to the `token → TOKEN "=" STRING` production rule.
type StringTokenDecl struct {
	Name     string
	Value    string
	Position *lexer.Position
}

func (n *StringTokenDecl) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "TokenDecl::%s=%q", n.Name, n.Value)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *StringTokenDecl) Equal(rhs Node) bool {
	nn, ok := rhs.(*StringTokenDecl)
	return ok &&
		n.Name == nn.Name &&
		n.Value == nn.Value &&
		equalPositions(n.Position, nn.Position)
}

func (n *StringTokenDecl) Pos() *lexer.Position {
	return n.Position
}

func (n *StringTokenDecl) leaf() {}

func (n *StringTokenDecl) decl() {}

// RegexTokenDecl represents a token declaration with a regular expression in an EBNF grammar.
// This node corresponds to the `token → TOKEN "=" REGEX` production rule.
type RegexTokenDecl struct {
	Name     string
	Regex    string
	Position *lexer.Position
}

func (n *RegexTokenDecl) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "TokenDecl::%s=/%s/", n.Name, n.Regex)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *RegexTokenDecl) Equal(rhs Node) bool {
	nn, ok := rhs.(*RegexTokenDecl)
	return ok &&
		n.Name == nn.Name &&
		n.Regex == nn.Regex &&
		equalPositions(n.Position, nn.Position)
}

func (n *RegexTokenDecl) Pos() *lexer.Position {
	return n.Position
}

func (n *RegexTokenDecl) leaf() {}

func (n *RegexTokenDecl) decl() {}

// PrecedenceDecl represents a precedence declaration in an EBNF grammar.
// This node corresponds to the `token → ("@left" | "@right" | "@none") {{handle}}` production rule.
type PrecedenceDecl struct {
	Associativity lr.Associativity
	Handles       []PrecedenceHandle
	Position      *lexer.Position
}

func (n *PrecedenceDecl) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "PrecedenceDecl::%s=%s", n.Associativity, n.Handles)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *PrecedenceDecl) Equal(rhs Node) bool {
	nn, ok := rhs.(*PrecedenceDecl)
	if !ok {
		return false
	}

	if n.Associativity != nn.Associativity {
		return false
	}

	if len(n.Handles) != len(nn.Handles) {
		return false
	}

	for i := range len(n.Handles) {
		if !n.Handles[i].Equal(nn.Handles[i]) {
			return false
		}
	}

	return equalPositions(n.Position, nn.Position)
}

func (n *PrecedenceDecl) Pos() *lexer.Position {
	return n.Position
}

func (n *PrecedenceDecl) Children() []Node {
	nodes := make([]Node, len(n.Handles))
	for i, handle := range n.Handles {
		nodes[i] = handle
	}

	return nodes
}

func (n *PrecedenceDecl) decl() {}

// PrecedenceHandle represents a handle in a precedence level within an EBNF grammar.
// This node corresponds to the `handle → term | "<" rule ">"` production rule.
type PrecedenceHandle interface {
	Node
	precedenceHandle()
}

// TerminalHandle represents a terminal handle in a precedence level within an EBNF grammar.
// This node corresponds to the `handle → term` production rule.
type TerminalHandle struct {
	Terminal string
	Position *lexer.Position
}

func (n *TerminalHandle) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "TerminalHandle::%s", n.Terminal)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *TerminalHandle) Equal(rhs Node) bool {
	nn, ok := rhs.(*TerminalHandle)
	return ok &&
		n.Terminal == nn.Terminal &&
		equalPositions(n.Position, nn.Position)
}

func (n *TerminalHandle) Pos() *lexer.Position {
	return n.Position
}

func (n *TerminalHandle) leaf() {}

func (n *TerminalHandle) precedenceHandle() {}

// ProductionHandle represents a production handle in a precedence level within an EBNF grammar.
// This node corresponds to the `handle → "<" rule ">"` production rule.
type ProductionHandle struct {
	LHS      string
	RHS      RHS
	Position *lexer.Position
}

func (n *ProductionHandle) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "ProductionHandle::%s → %s", n.LHS, n.RHS)

	return b.String()
}

func (n *ProductionHandle) Equal(rhs Node) bool {
	nn, ok := rhs.(*ProductionHandle)
	return ok &&
		n.LHS == nn.LHS &&
		n.RHS.Equal(nn.RHS) &&
		equalPositions(n.Position, nn.Position)
}

func (n *ProductionHandle) Pos() *lexer.Position {
	return n.Position
}

func (n *ProductionHandle) Children() []Node {
	nodes := []Node{n.RHS}

	return nodes
}

func (n *ProductionHandle) precedenceHandle() {}

// RuleDecl represents a rule declaration in an EBNF grammar.
// This node corresponds to the `rule → lhs "=" [rhs]` production rule.
type RuleDecl struct {
	LHS      string
	RHS      RHS
	Position *lexer.Position
}

func (n *RuleDecl) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "RuleDecl::%s → %s", n.LHS, n.RHS)

	return b.String()
}

func (n *RuleDecl) Equal(rhs Node) bool {
	nn, ok := rhs.(*RuleDecl)
	return ok &&
		n.LHS == nn.LHS &&
		n.RHS.Equal(nn.RHS) &&
		equalPositions(n.Position, nn.Position)
}

func (n *RuleDecl) Pos() *lexer.Position {
	return n.Position
}

func (n *RuleDecl) Children() []Node {
	nodes := []Node{n.RHS}

	return nodes
}

func (n *RuleDecl) decl() {}

// RHS represents the right-hand side (rhs) non-terminal in an EBNF grammar.
// This node corresponds to the `rhs → rhs rhs | rhs "|" rhs | "(" rhs ")" | "[" rhs "]" | "{" rhs "}" | "{{" rhs "}}" | nonterm | term | ε` production rule.
type RHS interface {
	Node
	rhs()
}

// ConcatRHS represents a concatenation of RHS nodes in an EBNF grammar.
// This node corresponds to the `rhs → rhs rhs` production rule.
type ConcatRHS struct {
	Ops []RHS
}

func (n *ConcatRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "ConcatRHS::")
	for _, op := range n.Ops {
		fmt.Fprintf(&b, "%s ", op)
	}

	if len(n.Ops) > 0 {
		b.Truncate(b.Len() - 1)
	}

	return b.String()
}

func (n *ConcatRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*ConcatRHS)
	if !ok {
		return false
	}

	if len(n.Ops) != len(nn.Ops) {
		return false
	}

	for i := range len(n.Ops) {
		if !n.Ops[i].Equal(nn.Ops[i]) {
			return false
		}
	}

	return true
}

func (n *ConcatRHS) Pos() *lexer.Position {
	if len(n.Ops) > 0 {
		return n.Ops[0].Pos()
	}

	return nil
}

func (n *ConcatRHS) Children() []Node {
	nodes := make([]Node, len(n.Ops))
	for i, op := range n.Ops {
		nodes[i] = op
	}

	return nodes
}

func (n *ConcatRHS) rhs() {}

// AltRHS represents an alternative between RHS nodes in an EBNF grammar.
// This node corresponds to the `rhs → rhs "|" rhs` production rule.
type AltRHS struct {
	Ops []RHS
}

func (n *AltRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "AltRHS::")
	for _, op := range n.Ops {
		fmt.Fprintf(&b, "%s \"|\" ", op)
	}

	if len(n.Ops) > 0 {
		b.Truncate(b.Len() - 5)
	}

	return b.String()
}

func (n *AltRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*AltRHS)
	if !ok {
		return false
	}

	if len(n.Ops) != len(nn.Ops) {
		return false
	}

	for i := range len(n.Ops) {
		if !n.Ops[i].Equal(nn.Ops[i]) {
			return false
		}
	}

	return true
}

func (n *AltRHS) Pos() *lexer.Position {
	if len(n.Ops) > 0 {
		return n.Ops[0].Pos()
	}

	return nil
}

func (n *AltRHS) Children() []Node {
	nodes := make([]Node, len(n.Ops))
	for i, op := range n.Ops {
		nodes[i] = op
	}

	return nodes
}

func (n *AltRHS) rhs() {}

// OptRHS represents an optional RHS node in an EBNF grammar.
// This node corresponds to the `rhs → "[" rhs "]"` production rule.
type OptRHS struct {
	Op       RHS
	Position *lexer.Position
}

func (n *OptRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "OptRHS::%s", n.Op)

	return b.String()
}

func (n *OptRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*OptRHS)
	return ok &&
		n.Op.Equal(nn.Op) &&
		equalPositions(n.Position, nn.Position)
}

func (n *OptRHS) Pos() *lexer.Position {
	return n.Position
}

func (n *OptRHS) Children() []Node {
	nodes := []Node{n.Op}

	return nodes
}

func (n *OptRHS) rhs() {}

// StarRHS represents a Kleene star (*) operation on an RHS node in an EBNF grammar.
// This node corresponds to the `rhs → "{" rhs "}"` production rule.
type StarRHS struct {
	Op       RHS
	Position *lexer.Position
}

func (n *StarRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "StarRHS::%s", n.Op)

	return b.String()
}

func (n *StarRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*StarRHS)
	return ok &&
		n.Op.Equal(nn.Op) &&
		equalPositions(n.Position, nn.Position)
}

func (n *StarRHS) Pos() *lexer.Position {
	return n.Position
}

func (n *StarRHS) Children() []Node {
	nodes := []Node{n.Op}

	return nodes
}

func (n *StarRHS) rhs() {}

// PlusRHS represents a Kleene plus (+) operation on an RHS node in an EBNF grammar.
// This node corresponds to the `rhs → "{{" rhs "}}"` production rule.
type PlusRHS struct {
	Op       RHS
	Position *lexer.Position
}

func (n *PlusRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "PlusRHS::%s", n.Op)

	return b.String()
}

func (n *PlusRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*PlusRHS)
	return ok &&
		n.Op.Equal(nn.Op) &&
		equalPositions(n.Position, nn.Position)
}

func (n *PlusRHS) Pos() *lexer.Position {
	return n.Position
}

func (n *PlusRHS) Children() []Node {
	nodes := []Node{n.Op}

	return nodes
}

func (n *PlusRHS) rhs() {}

// NonTerminalRHS represents a non-terminal symbol as the right-hand side of a rule in an EBNF grammar.
// This node corresponds to the `rhs → nonterm` production rule.
type NonTerminalRHS struct {
	NonTerminal string
	Position    *lexer.Position
}

func (n *NonTerminalRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "NonTerminalRHS::%s", n.NonTerminal)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *NonTerminalRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*NonTerminalRHS)
	return ok &&
		n.NonTerminal == nn.NonTerminal &&
		equalPositions(n.Position, nn.Position)
}

func (n *NonTerminalRHS) Pos() *lexer.Position {
	return n.Position
}

func (n *NonTerminalRHS) leaf() {}

func (n *NonTerminalRHS) rhs() {}

// TerminalRHS represents a terminal symbol as the right-hand side of a rule in an EBNF grammar.
// This node corresponds to the `rhs → term` production rule.
type TerminalRHS struct {
	Terminal string
	Position *lexer.Position
}

func (n *TerminalRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "TerminalRHS::%s", n.Terminal)
	if hasPositionValue(n.Position) {
		fmt.Fprintf(&b, " <%s>", n.Position)
	}

	return b.String()
}

func (n *TerminalRHS) Equal(rhs Node) bool {
	nn, ok := rhs.(*TerminalRHS)
	return ok &&
		n.Terminal == nn.Terminal &&
		equalPositions(n.Position, nn.Position)
}

func (n *TerminalRHS) Pos() *lexer.Position {
	return n.Position
}

func (n *TerminalRHS) leaf() {}

func (n *TerminalRHS) rhs() {}

// EmptyRHS represents the empty string ε as the right-hand side of a rule in an EBNF grammar.
// This node corresponds to the `rhs → ε` production rule.
type EmptyRHS struct{}

func (n *EmptyRHS) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "EmptyRHS::ε")

	return b.String()
}

func (n *EmptyRHS) Equal(rhs Node) bool {
	_, ok := rhs.(*EmptyRHS)
	return ok
}

func (n *EmptyRHS) Pos() *lexer.Position {
	return nil
}

func (n *EmptyRHS) leaf() {}

func (n *EmptyRHS) rhs() {}

// hasPositionValue checks if a position has a non-nil non-zero value.
func hasPositionValue(pos *lexer.Position) bool {
	return pos != nil && !pos.IsZero()
}

// equalPositions determines whether or not two positions are the same.
func equalPositions(lhs, rhs *lexer.Position) bool {
	if lhs == nil || rhs == nil {
		return lhs == rhs
	}
	return lhs.Equal(*rhs)
}
