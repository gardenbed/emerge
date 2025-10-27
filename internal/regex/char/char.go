// Package char is a common package for regex parsers.
// It provides ranges of characters and helper functions for working with character classes, Unicode categories, etc.
package char

import (
	"unicode"

	"github.com/moorara/algo/range/disc"
)

// Range represents an inclusive range of characters.
type Range [2]rune

// RangeList represents a list of inclusive ranges of characters.
type RangeList []Range

// Dedup returns a new RangeList with overlapping and adjacent ranges merged.
func (l RangeList) Dedup() RangeList {
	list := disc.NewRangeList[rune](nil)
	for _, r := range l {
		list.Add(disc.Range[rune]{Lo: r[0], Hi: r[1]})
	}

	var res RangeList
	for r := range list.All() {
		res = append(res, Range{r.Lo, r.Hi})
	}

	return res
}

// Exclude returns a new RangeList that is the result of excluding the ranges in x from l.
func (l RangeList) Exclude(x RangeList) RangeList {
	list := disc.NewRangeList[rune](nil)
	for _, r := range l {
		list.Add(disc.Range[rune]{Lo: r[0], Hi: r[1]})
	}

	for _, r := range x {
		list.Remove(disc.Range[rune]{Lo: r[0], Hi: r[1]})
	}

	var res RangeList
	for r := range list.All() {
		res = append(res, Range{r.Lo, r.Hi})
	}

	return res
}

// unicodeCategoryToRanges converts a Go unicode.RangeTable into a flat list of Range.
func unicodeCategoryToRanges(c *unicode.RangeTable) RangeList {
	var ranges RangeList

	for _, r := range c.R16 {
		if r.Stride == 1 {
			ranges = append(ranges, Range{rune(r.Lo), rune(r.Hi)})
		} else {
			for v := r.Lo; v <= r.Hi; v += r.Stride {
				ranges = append(ranges, Range{rune(v), rune(v)})
			}
		}
	}

	for _, r := range c.R32 {
		if r.Stride == 1 {
			ranges = append(ranges, Range{rune(r.Lo), rune(r.Hi)})
		} else {
			for v := r.Lo; v <= r.Hi; v += r.Stride {
				ranges = append(ranges, Range{rune(v), rune(v)})
			}
		}
	}

	return ranges
}

// Classes defines the collection of characters for different classes and categories of characters.
var Classes = map[string]RangeList{
	// All Characters
	`ASCII`:   {{0x00, 0x7F}},
	`UNICODE`: {{0x000000, 0x10FFFF}},

	// Character Classes
	`\s`: {{' ', ' '}, {'\t', '\n'}, {'\f', '\r'}},
	`\d`: {{'0', '9'}},
	`\w`: {{'0', '9'}, {'A', 'Z'}, {'_', '_'}, {'a', 'z'}},

	// ASCII Classes
	`[:blank:]`:  {{' ', ' '}, {'\t', '\t'}},
	`[:space:]`:  {{' ', ' '}, {'\t', '\r'}},
	`[:digit:]`:  {{'0', '9'}},
	`[:xdigit:]`: {{'0', '9'}, {'A', 'F'}, {'a', 'f'}},
	`[:upper:]`:  {{'A', 'Z'}},
	`[:lower:]`:  {{'a', 'z'}},
	`[:alpha:]`:  {{'A', 'Z'}, {'a', 'z'}},
	`[:alnum:]`:  {{'0', '9'}, {'A', 'Z'}, {'a', 'z'}},
	`[:word:]`:   {{'0', '9'}, {'A', 'Z'}, {'_', '_'}, {'a', 'z'}},
	`[:ascii:]`:  {{0x00, 0x7F}},

	/* Unicode Classes */

	// General - Letters
	`Letter`: unicodeCategoryToRanges(unicode.L),
	`L`:      unicodeCategoryToRanges(unicode.L),
	`Lu`:     unicodeCategoryToRanges(unicode.Lu),
	`Ll`:     unicodeCategoryToRanges(unicode.Ll),
	`Lt`:     unicodeCategoryToRanges(unicode.Lt),
	`Lm`:     unicodeCategoryToRanges(unicode.Lm),
	`Lo`:     unicodeCategoryToRanges(unicode.Lo),

	// General - Marks
	`Mark`: unicodeCategoryToRanges(unicode.M),
	`M`:    unicodeCategoryToRanges(unicode.M),
	`Mn`:   unicodeCategoryToRanges(unicode.Mn),
	`Mc`:   unicodeCategoryToRanges(unicode.Mc),
	`Me`:   unicodeCategoryToRanges(unicode.Me),

	// Numbers
	`Number`: unicodeCategoryToRanges(unicode.N),
	`N`:      unicodeCategoryToRanges(unicode.N),
	`Nd`:     unicodeCategoryToRanges(unicode.Nd),
	`Nl`:     unicodeCategoryToRanges(unicode.Nl),
	`No`:     unicodeCategoryToRanges(unicode.No),

	// General - Punctuations
	`Punctuation`: unicodeCategoryToRanges(unicode.P),
	`P`:           unicodeCategoryToRanges(unicode.P),
	`Pc`:          unicodeCategoryToRanges(unicode.Pc),
	`Pd`:          unicodeCategoryToRanges(unicode.Pd),
	`Ps`:          unicodeCategoryToRanges(unicode.Ps),
	`Pe`:          unicodeCategoryToRanges(unicode.Pe),
	`Pi`:          unicodeCategoryToRanges(unicode.Pi),
	`Pf`:          unicodeCategoryToRanges(unicode.Pf),
	`Po`:          unicodeCategoryToRanges(unicode.Po),

	// General - Symbols
	`Symbol`: unicodeCategoryToRanges(unicode.S),
	`S`:      unicodeCategoryToRanges(unicode.S),
	`Sm`:     unicodeCategoryToRanges(unicode.Sm),
	`Sc`:     unicodeCategoryToRanges(unicode.Sc),
	`Sk`:     unicodeCategoryToRanges(unicode.Sk),
	`So`:     unicodeCategoryToRanges(unicode.So),

	// General - Separator
	`Separator`: unicodeCategoryToRanges(unicode.Z),
	`Z`:         unicodeCategoryToRanges(unicode.Z),
	`Zs`:        unicodeCategoryToRanges(unicode.Zs),
	`Zl`:        unicodeCategoryToRanges(unicode.Zl),
	`Zp`:        unicodeCategoryToRanges(unicode.Zp),

	/* Scripts */

	`Latin`: {
		{0x0000, 0x007F}, // C0 Controls and Basic Latin
		{0x0080, 0x00FF}, // C1 Controls and Latin-1 Supplement
		{0x0100, 0x017F}, // Latin Extended-A
		{0x0180, 0x024F}, // Latin Extended-B
		{0x1E00, 0x1EFF}, // Latin Extended Additional
	},

	`Greek`: {
		{0x0370, 0x03FF}, // Greek and Coptic
		{0x1F00, 0x1FFF}, // Greek Extended
	},

	`Cyrillic`: {
		{0x0400, 0x04FF}, // Cyrillic
		{0x0500, 0x052F}, // Cyrillic Supplement
		{0x2DE0, 0x2DFF}, // Cyrillic Extended-A
		{0xA640, 0xA69F}, // Cyrillic Extended-B
		{0x1C80, 0x1C8F}, // Cyrillic Extended-C
	},

	`Han`: {
		{0x4E00, 0x9FFF},     // CJK Unified Ideographs
		{0x3400, 0x4DBF},     // CJK Unified Ideographs Extension A
		{0x020000, 0x02A6DF}, // CJK Unified Ideographs Extension B
		{0x02A700, 0x02B738}, // CJK Unified Ideographs Extension C
		{0x02B740, 0x02B81D}, // CJK Unified Ideographs Extension D
		{0x02B820, 0x02CEA1}, // CJK Unified Ideographs Extension E
		{0x02CEB0, 0x02EBE0}, // CJK Unified Ideographs Extension F
		{0x030000, 0x03134A}, // CJK Unified Ideographs Extension G
	},

	`Persian`: {
		{0x0600, 0x06FF},     // Arabic
		{0x0750, 0x077F},     // Arabic Supplement
		{0x08A0, 0x08FF},     // Arabic Extended-A
		{0x0870, 0x089F},     // Arabic Extended-B
		{0xFB50, 0xFDFF},     // Arabic Presentation Forms-A
		{0xFE70, 0xFEFF},     // Arabic Presentation Forms-B
		{0x0103A0, 0x0103DF}, // Old Persian
	},

	/* Derived */

	`Math`: {
		{0x2200, 0x22FF},     // Mathematical Operators
		{0x27C0, 0x27EF},     // Miscellaneous Mathematical Symbols-A
		{0x2980, 0x29FF},     // Miscellaneous Mathematical Symbols-B
		{0x2A00, 0x2AFF},     // Supplemental Mathematical Operators
		{0x01D400, 0x01D7FF}, // Mathematical Alphanumeric Symbols
	},

	`Emoji`: {
		{0x01F300, 0x01F5FF}, // Miscellaneous Symbols and Pictographs
		{0x01F600, 0x01F64F}, // Emoticons
		{0x01F680, 0x01F6FF}, // Transport and Map Symbols
		{0x01F900, 0x01F9FF}, // Supplemental Symbols and Pictographs
		{0x01FA70, 0x01FAFF}, // Symbols and Pictographs Extended-A
	},
}
