package parser

// Runes represents a collection of runes.
type Runes interface {
	Runes() []rune
}

// runeList is a list of runes that implements the Runes interface.
type runeList []rune

// Runes returns the slice of runes stored in runeList.
func (r runeList) Runes() []rune {
	return r
}

// runeRange represents a contiguous range of runes, defined by a start and end rune.
// It implements the Runes interface.
type runeRange [2]rune

// Runes returns a slice of runes representing all runes in the range, inclusive.
func (r runeRange) Runes() []rune {
	runes := make([]rune, 0, r[1]-r[0]+1)
	for i := r[0]; i <= r[1]; i++ {
		runes = append(runes, i)
	}

	return runes
}

// RuneClass represents a class of runes, consisting of one or more instances of Runes.
// It also implements the Runes interface to aggregate all runes in the class.
type RuneClass []Runes

// Runes aggregates and returns all runes in the class.
func (r RuneClass) Runes() []rune {
	runes := make([]rune, 0)
	for _, rs := range r {
		runes = append(runes, rs.Runes()...)
	}

	return runes
}

// RuneClasses defines the collection of runes for different classes and categories of characters.
var RuneClasses = map[string]RuneClass{
	// All Characters
	`ASCII`: {runeRange{0x00, 0x7F}},
	`UTF-8`: {runeRange{0x000000, 0x10FFFF}},

	// Character Classes
	`\s`: {runeList{' ', '\t', '\n', '\r', '\f'}},
	`\d`: {runeRange{'0', '9'}},
	`\w`: {runeRange{'0', '9'}, runeRange{'A', 'Z'}, runeList{'_'}, runeRange{'a', 'z'}},

	// ASCII Classes
	`[:blank:]`:  {runeList{' ', '\t'}},
	`[:space:]`:  {runeList{' ', '\t', '\n', '\r', '\f', '\v'}},
	`[:digit:]`:  {runeRange{'0', '9'}},
	`[:xdigit:]`: {runeRange{'0', '9'}, runeRange{'A', 'F'}, runeRange{'a', 'f'}},
	`[:upper:]`:  {runeRange{'A', 'Z'}},
	`[:lower:]`:  {runeRange{'a', 'z'}},
	`[:alpha:]`:  {runeRange{'A', 'Z'}, runeRange{'a', 'z'}},
	`[:alnum:]`:  {runeRange{'0', '9'}, runeRange{'A', 'Z'}, runeRange{'a', 'z'}},
	`[:word:]`:   {runeRange{'0', '9'}, runeRange{'A', 'Z'}, runeList{'_'}, runeRange{'a', 'z'}},
	`[:ascii:]`:  {runeRange{0x00, 0x7F}},

	/* TODO: Unicode Classes */

	// General - Letters
	`Letter`: {runeRange{'A', 'Z'}, runeRange{'a', 'z'}},
	`L`:      {runeRange{'A', 'Z'}, runeRange{'a', 'z'}},
	`Lu`:     {runeRange{'A', 'Z'}},
	`Ll`:     {runeRange{'a', 'z'}},
	`Lt`:     {},
	`Lm`:     {},
	`Lo`:     {},

	// General - Marks
	`Mark`: {},
	`M`:    {},
	`Mn`:   {},
	`Mc`:   {},
	`Me`:   {},

	// Numbers - Marks
	`Number`: {},
	`N`:      {},
	`Nd`:     {},
	`Nl`:     {},
	`No`:     {},

	// General - Punctuations
	`Punctuation`: {},
	`P`:           {},
	`Pc`:          {},
	`Pd`:          {},
	`Ps`:          {},
	`Pe`:          {},
	`Pi`:          {},
	`Pf`:          {},
	`Po`:          {},

	// General - Symbols
	`Symbol`: {},
	`S`:      {},
	`Sm`:     {},
	`Sc`:     {},
	`Sk`:     {},
	`So`:     {},

	// General - Separator
	`Separator`: {},
	`Z`:         {},
	`Zs`:        {},
	`Zl`:        {},
	`Zp`:        {},

	// Scripts
	`Latin`: {
		runeRange{0x0000, 0x007F}, // C0 Controls and Basic Latin
		runeRange{0x0080, 0x00FF}, // C1 Controls and Latin-1 Supplement
		runeRange{0x0100, 0x017F}, // Latin Extended-A
		runeRange{0x0180, 0x024F}, // Latin Extended-B
		runeRange{0x1E00, 0x1EFF}, // Latin Extended Additional
	},
	`Greek`: {
		runeRange{0x0370, 0x03FF}, // Greek and Coptic
		runeRange{0x1F00, 0x1FFF}, // Greek Extended
	},
	`Cyrillic`: {
		runeRange{0x0400, 0x04FF}, // Cyrillic
		runeRange{0x0500, 0x052F}, // Cyrillic Supplement
		runeRange{0x2DE0, 0x2DFF}, // Cyrillic Extended-A
		runeRange{0xA640, 0xA69F}, // Cyrillic Extended-B
		runeRange{0x1C80, 0x1C8F}, // Cyrillic Extended-C
	},
	`Han`: {
		runeRange{0x4E00, 0x9FFF},     // CJK Unified Ideographs
		runeRange{0x3400, 0x4DBF},     // CJK Unified Ideographs Extension A
		runeRange{0x020000, 0x02A6DF}, // CJK Unified Ideographs Extension B
		runeRange{0x02A700, 0x02B738}, // CJK Unified Ideographs Extension C
		runeRange{0x02B740, 0x02B81D}, // CJK Unified Ideographs Extension D
		runeRange{0x02B820, 0x02CEA1}, // CJK Unified Ideographs Extension E
		runeRange{0x02CEB0, 0x02EBE0}, // CJK Unified Ideographs Extension F
		runeRange{0x030000, 0x03134A}, // CJK Unified Ideographs Extension G
	},
	`Persian`: {
		runeRange{0x0600, 0x06FF},     // Arabic
		runeRange{0x0750, 0x077F},     // Arabic Supplement
		runeRange{0x08A0, 0x08FF},     // Arabic Extended-A
		runeRange{0x0870, 0x089F},     // Arabic Extended-B
		runeRange{0xFB50, 0xFDFF},     // Arabic Presentation Forms-A
		runeRange{0xFE70, 0xFEFF},     // Arabic Presentation Forms-B
		runeRange{0x0103A0, 0x0103DF}, // Old Persian
	},

	// Derived
	`Math`: {
		runeRange{0x2200, 0x22FF},     // Mathematical Operators
		runeRange{0x27C0, 0x27EF},     // Miscellaneous Mathematical Symbols-A
		runeRange{0x2980, 0x29FF},     // Miscellaneous Mathematical Symbols-B
		runeRange{0x2A00, 0x2AFF},     // Supplemental Mathematical Operators
		runeRange{0x01D400, 0x01D7FF}, // Mathematical Alphanumeric Symbols
	},
	`Emoji`: {
		runeRange{0x01F300, 0x01F5FF}, // Miscellaneous Symbols and Pictographs
		runeRange{0x01F600, 0x01F64F}, // Emoticons
		runeRange{0x01F680, 0x01F6FF}, // Transport and Map Symbols
		runeRange{0x01F900, 0x01F9FF}, // Supplemental Symbols and Pictographs
		runeRange{0x01FA70, 0x01FAFF}, // Symbols and Pictographs Extended-A
	},
}
