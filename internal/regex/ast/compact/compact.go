// package compact implements a more complex abstract syntax tree for regular expressions.
package compact

import (
	"github.com/gardenbed/emerge/internal/regex/token"
)

// Node is the type for all nodes in an abstract syntax tree.
// It represents both the terminal and non-terminal symbols in a grammar.
type Node interface {
	// position of the first character belonging to the node.
	Pos() token.Pos
}

// Regex represents a regular expression node (root node).
//
// Production Rule: regex --> "^"? expr
type Regex struct {
	SOS  bool
	Expr Expr
}

func (n *Regex) Pos() token.Pos {
	if n.SOS {
		return n.Expr.Pos() - 1
	}
	return n.Expr.Pos()
}

// Expr represents an expression node.
//
// Production Rule: expr --> subexpr ("|" expr)?
type Expr struct {
	Sub  Subexpr
	Expr *Expr
}

func (n *Expr) Pos() token.Pos {
	return n.Sub.Pos()
}

// Subexpr represents a subexpression node.
//
// Production Rule: subexpr --> subexpr_item+
//                          --> subexpr_item subexpr_item*
type Subexpr struct {
	Items []SubexprItem
}

func (n *Subexpr) Pos() token.Pos {
	return n.Items[0].Pos()
}

// SubexprItem represents a subexpression item node.
//
// Production Rule: subexpr_item --> group | anchor | backref | match
type SubexprItem interface {
	Node
	implSubexprItem()
}

// Anchor represents an anchor node.
//
// Production Rule: anchor --> "$"
type Anchor struct {
	TokPos token.Pos
}

func (n *Anchor) Tag() token.Tag {
	return token.END_OF_STRING
}

func (n *Anchor) Pos() token.Pos {
	return n.TokPos
}

func (n *Anchor) implSubexprItem() {}

// Backref represents a backreference node.
//
// Production Rule: backref --> "\" num
type Backref struct {
	SlashPos token.Pos
	Ref      Num
	Group    *Group
}

func (n *Backref) Pos() token.Pos {
	return n.SlashPos
}

func (n *Backref) implSubexprItem() {}

// Group represents a group node.
//
// Production Rule: group --> "(" expr ")" quantifier?
type Group struct {
	OpenPos  token.Pos
	ClosePos token.Pos
	Expr     Expr
	Quant    *Quantifier
}

func (n *Group) Pos() token.Pos {
	return n.OpenPos
}

func (n *Group) implSubexprItem() {}

// Match represents a match node.
//
// Production Rule: match --> match_item quantifier?
type Match struct {
	Item  MatchItem
	Quant *Quantifier
}

func (n *Match) Pos() token.Pos {
	return n.Item.Pos()
}

func (n *Match) implSubexprItem() {}

// MatchItem represents a match item node.
//
// Production Rule: match_item --> any_char | char_class | ascii_char_class | char_group | char /* excluding | ) */
type MatchItem interface {
	Node
	implMatchItem()
}

// AnyChar represents an any character node.
//
// Production Rule: any_char --> "."
type AnyChar struct {
	TokPos token.Pos
}

func (n *AnyChar) Tag() token.Tag {
	return token.ANY_CHAR
}

func (n *AnyChar) Pos() token.Pos {
	return n.TokPos
}

func (n *AnyChar) implMatchItem() {}

// CharClass represents a character class node.
//
// Production Rule: char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
type CharClass struct {
	StartPos token.Pos
	EndPos   token.Pos
	TokTag   token.Tag
}

func (n *CharClass) Tag() token.Tag {
	return n.TokTag
}

func (n *CharClass) Pos() token.Pos {
	return n.StartPos
}

func (n *CharClass) implMatchItem()     {}
func (n *CharClass) implCharGroupItem() {}

// ASCIICharClass represents an ASCII character class node.
//
// Production Rule: ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
type ASCIICharClass struct {
	StartPos token.Pos
	EndPos   token.Pos
	TokTag   token.Tag
}

func (n *ASCIICharClass) Tag() token.Tag {
	return n.TokTag
}

func (n *ASCIICharClass) Pos() token.Pos {
	return n.StartPos
}

func (n *ASCIICharClass) implMatchItem()     {}
func (n *ASCIICharClass) implCharGroupItem() {}

// CharGroup represents a character group node.
//
// Production Rule: char_group --> "[" "^"? char_group_item+ "]"
type CharGroup struct {
	OpenPos  token.Pos
	ClosePos token.Pos
	Negated  bool
	Items    []CharGroupItem
}

func (n *CharGroup) Pos() token.Pos {
	return n.OpenPos
}

func (n *CharGroup) implMatchItem() {}

// CharGroupItem represents a character group item node.
//
// Production Rule: char_group_item --> char_class | ascii_char_class | char_range | char /* excluding ] */
type CharGroupItem interface {
	Node
	implCharGroupItem()
}

// CharRange represents a character range node.
//
// Production Rule: char_range --> char "-" char
type CharRange struct {
	Low Char
	Up  Char
}

func (n *CharRange) Pos() token.Pos {
	return n.Low.Pos()
}

func (n *CharRange) implCharGroupItem() {}

// Quantifier represents a quantifier node.
//
// Production Rule: quantifier --> repetition "?"?
type Quantifier struct {
	Rep  Repetition
	Lazy bool
}

func (n *Quantifier) Pos() token.Pos {
	return n.Rep.Pos()
}

// Repetition represents a repetition node.
//
// Production Rule: repetition --> rep_op | range
type Repetition interface {
	Node
	implRepetition()
}

// RepOp represents a repetition operator node.
//
// Production Rule: rep_op --> "?" | "*" | "+"
type RepOp struct {
	TokPos token.Pos
	TokTag token.Tag
}

func (n *RepOp) Tag() token.Tag {
	return n.TokTag
}

func (n *RepOp) Pos() token.Pos {
	return n.TokPos
}

func (n *RepOp) implRepetition() {}

// Range represents a range node.
//
// Production Rule: range --> "{" num upper_bound? "}"
type Range struct {
	OpenPos  token.Pos
	ClosePos token.Pos
	Low      Num
	Up       *UpperBound
}

func (n *Range) Pos() token.Pos {
	return n.OpenPos
}

func (n *Range) implRepetition() {}

// UpperBound represents an upper bound node.
//
// Production Rule: upper_bound --> "," num?
type UpperBound struct {
	CommaPos token.Pos
	Val      *Num
}

func (n *UpperBound) Pos() token.Pos {
	return n.CommaPos
}

// Num represents a number node.
//
// Production Rule: num --> NUM
//
// NUM is a token defined by regex [0-9]+
// The production rule here is implemented a bit different from the definition in the documentation for practical reasons.
// We do not want to create a node in AST for every terminal symbol.
type Num struct {
	StartPos token.Pos
	EndPos   token.Pos
	Val      int
}

func (n *Num) Tag() token.Tag {
	return token.NUM
}

func (n *Num) Pos() token.Pos {
	return n.StartPos
}

// Letters represents a letters node.
//
// Production Rule: letters --> LETTERS
//
// is a token defined by regex [A-Za-z]+
// The production rule here is implemented a bit different from the definition in the documentation for practical reasons.
// We do not want to create a node in AST for every terminal symbol.
type Letters struct {
	StartPos token.Pos
	EndPos   token.Pos
	Val      string
}

func (n *Letters) Tag() token.Tag {
	return token.LETTERS
}

func (n *Letters) Pos() token.Pos {
	return n.StartPos
}

// Char represents a character node.
//
// Production Rule: char --> CHAR
//
// is a token defined by regex [\x20-\x7E]
// The production rule here is implemented a bit different from the definition in the documentation for practical reasons.
// We do not want to create a node in AST for every terminal symbol.
type Char struct {
	TokPos token.Pos
	Val    rune
}

func (n *Char) Tag() token.Tag {
	return token.CHAR
}

func (n *Char) Pos() token.Pos {
	return n.TokPos
}

func (n *Char) implMatchItem()     {}
func (n *Char) implCharGroupItem() {}
