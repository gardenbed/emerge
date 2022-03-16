package regex

// Token represents the kind of a terminal symbol.
type Token int

const (
	END_OF_STRING     Token = iota // "$"
	WORD_BOUNDARY                  // "\b"
	NOT_WORD_BOUNDARY              // "\B"

	DIGIT_CHARS     // "\d"
	NOT_DIGIT_CHARS // "\D"
	WHITESPACE      // "\s"
	NOT_WHITESPACE  // "\S"
	WORD_CHARS      // "\w"
	NOT_WORD_CHARS  // "\W"

	BLANK  // "[:blank:]"
	SPACE  // "[:space:]"
	DIGIT  // "[:digit:]"
	XDIGIT // "[:xdigit:]"
	UPPER  // "[:upper:]"
	LOWER  // "[:lower:]"
	ALPHA  // "[:alpha:]"
	ALNUM  // "[:alnum:]"
	WORD   // "[:word:]"
	ASCII  // "[:ascii:]"
)

// Node is the type for all nodes in abstract syntax tree.
type Node interface {
	Pos() int // position of the first character belonging to the node
}

// Regex represents a regular expression node (root).
//
// Production Rule: regex --> "^"? expr
type Regex struct {
	Begin bool
	Expr  Expr
}

func (n *Regex) Pos() int {
	if n.Begin {
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

func (n *Expr) Pos() int { return n.Sub.Pos() }

// Subexpr represents a subexpression node.
//
// Production Rule: subexpr --> subexpr_item+
type Subexpr struct {
	Items []SubexprItem
}

func (n *Subexpr) Pos() int { return n.Items[0].Pos() }

// SubexprItem represents a subexpression item node.
//
// Production Rule: subexpr_item --> group | anchor | backref | match
type SubexprItem interface {
	Node
	implSubexprItem()
}

// Anchor represents an anchor node.
//
// Production Rule: anchor --> "$" | "\b" | "\B"
type Anchor struct {
	TokPos int
	Tok    Token
}

func (n *Anchor) Pos() int         { return n.TokPos }
func (n *Anchor) implSubexprItem() {}

// Backref represents a backreference node.
//
// Production Rule: backref --> "\" num
type Backref struct {
	SlashPos int
	Ref      Num
}

func (n *Backref) Pos() int         { return n.SlashPos }
func (n *Backref) implSubexprItem() {}

// Group represents a group node.
//
// Production Rule: group --> "(" "?:"? expr ")" quantifier?
type Group struct {
	OpenPos int
	NonCap  bool
	Expr    Expr
	Quant   *Quantifier
}

func (n *Group) Pos() int         { return n.OpenPos }
func (n *Group) implSubexprItem() {}

// Match represents a match node.
//
// Production Rule: match --> match_item quantifier?
type Match struct {
	Item  MatchItem
	Quant *Quantifier
}

func (n *Match) Pos() int         { return n.Item.Pos() }
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
	TokPos int
}

func (n *AnyChar) Pos() int       { return n.TokPos }
func (n *AnyChar) implMatchItem() {}

// CharClass represents a character class node.
//
// Production Rule: char_class --> "\d" | "\D" | "\s" | "\S" | "\w" | "\W"
type CharClass struct {
	TokPos int
	Tok    Token
}

func (n *CharClass) Pos() int           { return n.TokPos }
func (n *CharClass) implMatchItem()     {}
func (n *CharClass) implCharGroupItem() {}

// ASCIICharClass represents an ASCII character class node.
//
// Production Rule: ascii_char_class --> "[:blank:]" | "[:space:]" | "[:digit:]" | "[:xdigit:]" | "[:upper:]" | "[:lower:]" | "[:alpha:]" | "[:alnum:]" | "[:word:]" | "[:ascii:]"
type ASCIICharClass struct {
	TokPos int
	Tok    Token
}

func (n *ASCIICharClass) Pos() int           { return n.TokPos }
func (n *ASCIICharClass) implMatchItem()     {}
func (n *ASCIICharClass) implCharGroupItem() {}

// CharGroup represents a character group node.
//
// Production Rule: char_group --> "[" "^"? char_group_item+ "]"
type CharGroup struct {
	OpenPos int
	Neg     bool
	Items   []CharGroupItem
}

func (n *CharGroup) Pos() int       { return n.OpenPos }
func (n *CharGroup) implMatchItem() {}

// CharGroupItem represents a character group item node.
//
// Production Rule: char_group_item -->  char_class | ascii_char_class | char_range | char /* excluding ] */
type CharGroupItem interface {
	Node
	implCharGroupItem()
}

// CharRange represents a character range node.
//
// Production Rule: char_range --> char ("-" char)?
type CharRange struct {
	Low Char
	Up  Char
}

func (n *CharRange) Pos() int           { return n.Low.Pos() }
func (n *CharRange) implCharGroupItem() {}

// Quantifier represents a quantifier node.
//
// Production Rule: quantifier --> cardinality "?"?
type Quantifier struct {
	Card Cardinality
	Lazy bool
}

func (n *Quantifier) Pos() int { return n.Card.Pos() }

// Cardinality represents a cardinality node.
//
// Production Rule: cardinality --> zero_or_one | zero_or_more | one_or_more | range
type Cardinality interface {
	Node
	implCardinality()
}

// Range represents a zero-or-one (?) node.
//
// Production Rule: zero_or_one --> "?"
type ZeroOrOne struct {
	TokPos int
}

func (n *ZeroOrOne) Pos() int         { return n.TokPos }
func (n *ZeroOrOne) implCardinality() {}

// Range represents a zero-or-more (*) node.
//
// Production Rule: zero_or_more --> "*"
type ZeroOrMore struct {
	TokPos int
}

func (n *ZeroOrMore) Pos() int         { return n.TokPos }
func (n *ZeroOrMore) implCardinality() {}

// OneOrMore represents a one-or-more (+) node.
//
// Production Rule: one_or_more --> "+"
type OneOrMore struct {
	TokPos int
}

func (n *OneOrMore) Pos() int         { return n.TokPos }
func (n *OneOrMore) implCardinality() {}

// Range represents a range node.
//
// Production Rule: range --> "{" num upper_bound? "}"
type Range struct {
	OpenPos int
	Low     Num
	Up      *UpperBound
}

func (n *Range) Pos() int         { return n.OpenPos }
func (n *Range) implCardinality() {}

// UpperBound represents an upper bound node.
//
// Production Rule: upper_bound --> "," num?
type UpperBound struct {
	CommaPos int
	Val      *Num
}

func (n *UpperBound) Pos() int { return n.CommaPos }

// Num represents a number node.
//
// Production Rule: num --> INT
//
// The production rule here is implemented a bit different from the definition in the documentation for obvious reasons.
// We do not want to create a node in AST for every digit.
type Num struct {
	TokPos int
	Val    int
}

func (n *Num) Pos() int { return n.TokPos }

// Letters represents a letters node.
//
// Production Rule: letters --> STRING
//
// The production rule here is implemented a bit different from the definition in the documentation for obvious reasons.
// We do not want to create a node in AST for every letter.
type Letters struct {
	TokPos int
	Val    string
}

func (n *Letters) Pos() int { return n.TokPos }

// Char represents a character node.
//
// Production Rule: char --> CHAR
//
// The production rule here is implemented a bit different from the definition in the documentation for obvious reasons.
// We do not want to create a node in AST for every character.
type Char struct {
	TokPos int
	Val    rune
}

func (n *Char) Pos() int           { return n.TokPos }
func (n *Char) implMatchItem()     {}
func (n *Char) implCharGroupItem() {}
