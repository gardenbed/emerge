package nfa

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/algo/automata"
	"github.com/moorara/algo/parser/combinator"

	"github.com/gardenbed/emerge/internal/char"
)

var testNFA = map[string]*automata.NFA{
	"x": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 'x', 'x', []automata.State{1}).
		Build(),

	"\n": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '\n', '\n', []automata.State{1}).
		Build(),

	"\r": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '\r', '\r', []automata.State{1}).
		Build(),

	/* CHAR CLASSES */

	"ws": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, ' ', ' ', []automata.State{1}).
		AddTransition(0, '\t', '\r', []automata.State{1}).
		Build(),

	"not_ws": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x00, 0x08, []automata.State{1}).
		AddTransition(0, 0x0E, 0x1F, []automata.State{1}).
		AddTransition(0, 0x21, 0x10FFFF, []automata.State{1}).
		Build(),

	"digit": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '0', '9', []automata.State{1}).
		Build(),

	"not_digit": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x00, 0x2F, []automata.State{1}).
		AddTransition(0, 0x3A, 0x10FFFF, []automata.State{1}).
		Build(),

	"word": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '0', '9', []automata.State{1}).
		AddTransition(0, 'A', 'Z', []automata.State{1}).
		AddTransition(0, '_', '_', []automata.State{1}).
		AddTransition(0, 'a', 'z', []automata.State{1}).
		Build(),

	"not_word": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x00, 0x2F, []automata.State{1}).
		AddTransition(0, 0x3A, 0x40, []automata.State{1}).
		AddTransition(0, 0x5B, 0x5E, []automata.State{1}).
		AddTransition(0, 0x60, 0x60, []automata.State{1}).
		AddTransition(0, 0x7B, 0x10FFFF, []automata.State{1}).
		Build(),

	/* ASCII CLASSES */

	"blank": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, ' ', ' ', []automata.State{1}).
		AddTransition(0, '\t', '\t', []automata.State{1}).
		Build(),

	"space": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, ' ', ' ', []automata.State{1}).
		AddTransition(0, '\t', '\r', []automata.State{1}).
		Build(),

	// "digit" is already added.

	"xdigit": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '0', '9', []automata.State{1}).
		AddTransition(0, 'A', 'F', []automata.State{1}).
		AddTransition(0, 'a', 'f', []automata.State{1}).
		Build(),

	"upper": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 'A', 'Z', []automata.State{1}).
		Build(),

	"lower": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 'a', 'z', []automata.State{1}).
		Build(),

	"alpha": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 'A', 'Z', []automata.State{1}).
		AddTransition(0, 'a', 'z', []automata.State{1}).
		Build(),

	"alnum": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, '0', '9', []automata.State{1}).
		AddTransition(0, 'A', 'Z', []automata.State{1}).
		AddTransition(0, 'a', 'z', []automata.State{1}).
		Build(),

	// "word" is already added.

	"ascii": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x00, 0x7F, []automata.State{1}).
		Build(),

	/* UNICODE CLASSES */

	"Number": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x30, 0x39, []automata.State{1}).
		AddTransition(0, 0xB2, 0xB3, []automata.State{1}).
		AddTransition(0, 0xB9, 0xB9, []automata.State{1}).
		AddTransition(0, 0xBC, 0xBC, []automata.State{1}).
		AddTransition(0, 0xBD, 0xBE, []automata.State{1}).
		AddTransition(0, 0x0660, 0x0669, []automata.State{1}).
		AddTransition(0, 0x06F0, 0x06F9, []automata.State{1}).
		AddTransition(0, 0x07C0, 0x07C9, []automata.State{1}).
		AddTransition(0, 0x0966, 0x096F, []automata.State{1}).
		AddTransition(0, 0x09E6, 0x09EF, []automata.State{1}).
		AddTransition(0, 0x09F4, 0x09F9, []automata.State{1}).
		AddTransition(0, 0x0A66, 0x0A6F, []automata.State{1}).
		AddTransition(0, 0x0AE6, 0x0AEF, []automata.State{1}).
		AddTransition(0, 0x0B66, 0x0B6F, []automata.State{1}).
		AddTransition(0, 0x0B72, 0x0B77, []automata.State{1}).
		AddTransition(0, 0x0BE6, 0x0BF2, []automata.State{1}).
		AddTransition(0, 0x0C66, 0x0C6F, []automata.State{1}).
		AddTransition(0, 0x0C78, 0x0C7E, []automata.State{1}).
		AddTransition(0, 0x0CE6, 0x0CEF, []automata.State{1}).
		AddTransition(0, 0x0D58, 0x0D5E, []automata.State{1}).
		AddTransition(0, 0x0D66, 0x0D78, []automata.State{1}).
		AddTransition(0, 0x0DE6, 0x0DEF, []automata.State{1}).
		AddTransition(0, 0x0E50, 0x0E59, []automata.State{1}).
		AddTransition(0, 0x0ED0, 0x0ED9, []automata.State{1}).
		AddTransition(0, 0x0F20, 0x0F33, []automata.State{1}).
		AddTransition(0, 0x1040, 0x1049, []automata.State{1}).
		AddTransition(0, 0x1090, 0x1099, []automata.State{1}).
		AddTransition(0, 0x1369, 0x137C, []automata.State{1}).
		AddTransition(0, 0x16EE, 0x16F0, []automata.State{1}).
		AddTransition(0, 0x17E0, 0x17E9, []automata.State{1}).
		AddTransition(0, 0x17F0, 0x17F9, []automata.State{1}).
		AddTransition(0, 0x1810, 0x1819, []automata.State{1}).
		AddTransition(0, 0x1946, 0x194F, []automata.State{1}).
		AddTransition(0, 0x19D0, 0x19DA, []automata.State{1}).
		AddTransition(0, 0x1A80, 0x1A89, []automata.State{1}).
		AddTransition(0, 0x1A90, 0x1A99, []automata.State{1}).
		AddTransition(0, 0x1B50, 0x1B59, []automata.State{1}).
		AddTransition(0, 0x1BB0, 0x1BB9, []automata.State{1}).
		AddTransition(0, 0x1C40, 0x1C49, []automata.State{1}).
		AddTransition(0, 0x1C50, 0x1C59, []automata.State{1}).
		AddTransition(0, 0x2070, 0x2070, []automata.State{1}).
		AddTransition(0, 0x2074, 0x2074, []automata.State{1}).
		AddTransition(0, 0x2075, 0x2079, []automata.State{1}).
		AddTransition(0, 0x2080, 0x2089, []automata.State{1}).
		AddTransition(0, 0x2150, 0x2182, []automata.State{1}).
		AddTransition(0, 0x2185, 0x2189, []automata.State{1}).
		AddTransition(0, 0x2460, 0x249B, []automata.State{1}).
		AddTransition(0, 0x24EA, 0x24FF, []automata.State{1}).
		AddTransition(0, 0x2776, 0x2793, []automata.State{1}).
		AddTransition(0, 0x2CFD, 0x2CFD, []automata.State{1}).
		AddTransition(0, 0x3007, 0x3007, []automata.State{1}).
		AddTransition(0, 0x3021, 0x3029, []automata.State{1}).
		AddTransition(0, 0x3038, 0x303A, []automata.State{1}).
		AddTransition(0, 0x3192, 0x3195, []automata.State{1}).
		AddTransition(0, 0x3220, 0x3229, []automata.State{1}).
		AddTransition(0, 0x3248, 0x324F, []automata.State{1}).
		AddTransition(0, 0x3251, 0x325F, []automata.State{1}).
		AddTransition(0, 0x3280, 0x3289, []automata.State{1}).
		AddTransition(0, 0x32B1, 0x32BF, []automata.State{1}).
		AddTransition(0, 0xA620, 0xA629, []automata.State{1}).
		AddTransition(0, 0xA6E6, 0xA6EF, []automata.State{1}).
		AddTransition(0, 0xA830, 0xA835, []automata.State{1}).
		AddTransition(0, 0xA8D0, 0xA8D9, []automata.State{1}).
		AddTransition(0, 0xA900, 0xA909, []automata.State{1}).
		AddTransition(0, 0xA9D0, 0xA9D9, []automata.State{1}).
		AddTransition(0, 0xA9F0, 0xA9F9, []automata.State{1}).
		AddTransition(0, 0xAA50, 0xAA59, []automata.State{1}).
		AddTransition(0, 0xABF0, 0xABF9, []automata.State{1}).
		AddTransition(0, 0xFF10, 0xFF19, []automata.State{1}).
		AddTransition(0, 0x010107, 0x010133, []automata.State{1}).
		AddTransition(0, 0x010140, 0x010178, []automata.State{1}).
		AddTransition(0, 0x01018A, 0x01018B, []automata.State{1}).
		AddTransition(0, 0x0102E1, 0x0102FB, []automata.State{1}).
		AddTransition(0, 0x010320, 0x010323, []automata.State{1}).
		AddTransition(0, 0x010341, 0x010341, []automata.State{1}).
		AddTransition(0, 0x01034A, 0x01034A, []automata.State{1}).
		AddTransition(0, 0x0103D1, 0x0103D5, []automata.State{1}).
		AddTransition(0, 0x0104A0, 0x0104A9, []automata.State{1}).
		AddTransition(0, 0x010858, 0x01085F, []automata.State{1}).
		AddTransition(0, 0x010879, 0x01087F, []automata.State{1}).
		AddTransition(0, 0x0108A7, 0x0108AF, []automata.State{1}).
		AddTransition(0, 0x0108FB, 0x0108FF, []automata.State{1}).
		AddTransition(0, 0x010916, 0x01091B, []automata.State{1}).
		AddTransition(0, 0x0109BC, 0x0109BD, []automata.State{1}).
		AddTransition(0, 0x0109C0, 0x0109CF, []automata.State{1}).
		AddTransition(0, 0x0109D2, 0x0109FF, []automata.State{1}).
		AddTransition(0, 0x010A40, 0x010A48, []automata.State{1}).
		AddTransition(0, 0x010A7D, 0x010A7E, []automata.State{1}).
		AddTransition(0, 0x010A9D, 0x010A9F, []automata.State{1}).
		AddTransition(0, 0x010AEB, 0x010AEF, []automata.State{1}).
		AddTransition(0, 0x010B58, 0x010B5F, []automata.State{1}).
		AddTransition(0, 0x010B78, 0x010B7F, []automata.State{1}).
		AddTransition(0, 0x010BA9, 0x010BAF, []automata.State{1}).
		AddTransition(0, 0x010CFA, 0x010CFF, []automata.State{1}).
		AddTransition(0, 0x010D30, 0x010D39, []automata.State{1}).
		AddTransition(0, 0x010E60, 0x010E7E, []automata.State{1}).
		AddTransition(0, 0x010F1D, 0x010F26, []automata.State{1}).
		AddTransition(0, 0x010F51, 0x010F54, []automata.State{1}).
		AddTransition(0, 0x010FC5, 0x010FCB, []automata.State{1}).
		AddTransition(0, 0x011052, 0x01106F, []automata.State{1}).
		AddTransition(0, 0x0110F0, 0x0110F9, []automata.State{1}).
		AddTransition(0, 0x011136, 0x01113F, []automata.State{1}).
		AddTransition(0, 0x0111D0, 0x0111D9, []automata.State{1}).
		AddTransition(0, 0x0111E1, 0x0111F4, []automata.State{1}).
		AddTransition(0, 0x0112F0, 0x0112F9, []automata.State{1}).
		AddTransition(0, 0x011450, 0x011459, []automata.State{1}).
		AddTransition(0, 0x0114D0, 0x0114D9, []automata.State{1}).
		AddTransition(0, 0x011650, 0x011659, []automata.State{1}).
		AddTransition(0, 0x0116C0, 0x0116C9, []automata.State{1}).
		AddTransition(0, 0x011730, 0x01173B, []automata.State{1}).
		AddTransition(0, 0x0118E0, 0x0118F2, []automata.State{1}).
		AddTransition(0, 0x011950, 0x011959, []automata.State{1}).
		AddTransition(0, 0x011C50, 0x011C6C, []automata.State{1}).
		AddTransition(0, 0x011D50, 0x011D59, []automata.State{1}).
		AddTransition(0, 0x011DA0, 0x011DA9, []automata.State{1}).
		AddTransition(0, 0x011F50, 0x011F59, []automata.State{1}).
		AddTransition(0, 0x011FC0, 0x011FD4, []automata.State{1}).
		AddTransition(0, 0x012400, 0x01246E, []automata.State{1}).
		AddTransition(0, 0x016A60, 0x016A69, []automata.State{1}).
		AddTransition(0, 0x016AC0, 0x016AC9, []automata.State{1}).
		AddTransition(0, 0x016B50, 0x016B59, []automata.State{1}).
		AddTransition(0, 0x016B5B, 0x016B61, []automata.State{1}).
		AddTransition(0, 0x016E80, 0x016E96, []automata.State{1}).
		AddTransition(0, 0x01D2C0, 0x01D2D3, []automata.State{1}).
		AddTransition(0, 0x01D2E0, 0x01D2F3, []automata.State{1}).
		AddTransition(0, 0x01D360, 0x01D378, []automata.State{1}).
		AddTransition(0, 0x01D7CE, 0x01D7FF, []automata.State{1}).
		AddTransition(0, 0x01E140, 0x01E149, []automata.State{1}).
		AddTransition(0, 0x01E2F0, 0x01E2F9, []automata.State{1}).
		AddTransition(0, 0x01E4F0, 0x01E4F9, []automata.State{1}).
		AddTransition(0, 0x01E8C7, 0x01E8CF, []automata.State{1}).
		AddTransition(0, 0x01E950, 0x01E959, []automata.State{1}).
		AddTransition(0, 0x01EC71, 0x01ECAB, []automata.State{1}).
		AddTransition(0, 0x01ECAD, 0x01ECAF, []automata.State{1}).
		AddTransition(0, 0x01ECB1, 0x01ECB4, []automata.State{1}).
		AddTransition(0, 0x01ED01, 0x01ED2D, []automata.State{1}).
		AddTransition(0, 0x01ED2F, 0x01ED3D, []automata.State{1}).
		AddTransition(0, 0x01F100, 0x01F10C, []automata.State{1}).
		AddTransition(0, 0x01FBF0, 0x01FBF9, []automata.State{1}).
		Build(),

	"Unicode": automata.NewNFABuilder().
		SetStart(0).
		SetFinal([]automata.State{1}).
		AddTransition(0, 0x00, 0x10FFFF, []automata.State{1}).
		Build(),
}

var testRanges = map[string]char.RangeList{
	/* CHAR CLASSES */

	"ws":        {{' ', ' '}, {'\t', '\r'}},
	"not_ws":    {{0x00, 0x08}, {0x0E, 0x1F}, {0x21, 0x10FFFF}},
	"digit":     {{'0', '9'}},
	"not_digit": {{0x00, 0x2F}, {0x3A, 0x10FFFF}},
	"word":      {{'0', '9'}, {'A', 'Z'}, {'_', '_'}, {'a', 'z'}},
	"not_word":  {{0x00, 0x2F}, {0x3A, 0x40}, {0x5B, 0x5E}, {0x60, 0x60}, {0x7B, 0x10FFFF}},

	/* ASCII CLASSES */

	"blank": {{' ', ' '}, {'\t', '\t'}},
	"space": {{' ', ' '}, {'\t', '\r'}},
	// "digit" is already added.
	"xdigit": {{'0', '9'}, {'A', 'F'}, {'a', 'f'}},
	"upper":  {{'A', 'Z'}},
	"lower":  {{'a', 'z'}},
	"alpha":  {{'A', 'Z'}, {'a', 'z'}},
	"alnum":  {{'0', '9'}, {'A', 'Z'}, {'a', 'z'}},
	// "word" is already added.
	"ascii": {{0x00, 0x7F}},

	/* UNICODE CLASSES */

	"Number": {
		{0x30, 0x39},
		{0xB2, 0xB3},
		{0xB9, 0xB9},
		{0xBC, 0xBC},
		{0xBD, 0xBE},
		{0x0660, 0x0669},
		{0x06F0, 0x06F9},
		{0x07C0, 0x07C9},
		{0x0966, 0x096F},
		{0x09E6, 0x09EF},
		{0x09F4, 0x09F9},
		{0x0A66, 0x0A6F},
		{0x0AE6, 0x0AEF},
		{0x0B66, 0x0B6F},
		{0x0B72, 0x0B77},
		{0x0BE6, 0x0BF2},
		{0x0C66, 0x0C6F},
		{0x0C78, 0x0C7E},
		{0x0CE6, 0x0CEF},
		{0x0D58, 0x0D5E},
		{0x0D66, 0x0D78},
		{0x0DE6, 0x0DEF},
		{0x0E50, 0x0E59},
		{0x0ED0, 0x0ED9},
		{0x0F20, 0x0F33},
		{0x1040, 0x1049},
		{0x1090, 0x1099},
		{0x1369, 0x137C},
		{0x16EE, 0x16F0},
		{0x17E0, 0x17E9},
		{0x17F0, 0x17F9},
		{0x1810, 0x1819},
		{0x1946, 0x194F},
		{0x19D0, 0x19DA},
		{0x1A80, 0x1A89},
		{0x1A90, 0x1A99},
		{0x1B50, 0x1B59},
		{0x1BB0, 0x1BB9},
		{0x1C40, 0x1C49},
		{0x1C50, 0x1C59},
		{0x2070, 0x2070},
		{0x2074, 0x2074},
		{0x2075, 0x2079},
		{0x2080, 0x2089},
		{0x2150, 0x2182},
		{0x2185, 0x2189},
		{0x2460, 0x249B},
		{0x24EA, 0x24FF},
		{0x2776, 0x2793},
		{0x2CFD, 0x2CFD},
		{0x3007, 0x3007},
		{0x3021, 0x3029},
		{0x3038, 0x303A},
		{0x3192, 0x3195},
		{0x3220, 0x3229},
		{0x3248, 0x324F},
		{0x3251, 0x325F},
		{0x3280, 0x3289},
		{0x32B1, 0x32BF},
		{0xA620, 0xA629},
		{0xA6E6, 0xA6EF},
		{0xA830, 0xA835},
		{0xA8D0, 0xA8D9},
		{0xA900, 0xA909},
		{0xA9D0, 0xA9D9},
		{0xA9F0, 0xA9F9},
		{0xAA50, 0xAA59},
		{0xABF0, 0xABF9},
		{0xFF10, 0xFF19},
		{0x010107, 0x010133},
		{0x010140, 0x010178},
		{0x01018A, 0x01018B},
		{0x0102E1, 0x0102FB},
		{0x010320, 0x010323},
		{0x010341, 0x010341},
		{0x01034A, 0x01034A},
		{0x0103D1, 0x0103D5},
		{0x0104A0, 0x0104A9},
		{0x010858, 0x01085F},
		{0x010879, 0x01087F},
		{0x0108A7, 0x0108AF},
		{0x0108FB, 0x0108FF},
		{0x010916, 0x01091B},
		{0x0109BC, 0x0109BD},
		{0x0109C0, 0x0109CF},
		{0x0109D2, 0x0109FF},
		{0x010A40, 0x010A48},
		{0x010A7D, 0x010A7E},
		{0x010A9D, 0x010A9F},
		{0x010AEB, 0x010AEF},
		{0x010B58, 0x010B5F},
		{0x010B78, 0x010B7F},
		{0x010BA9, 0x010BAF},
		{0x010CFA, 0x010CFF},
		{0x010D30, 0x010D39},
		{0x010E60, 0x010E7E},
		{0x010F1D, 0x010F26},
		{0x010F51, 0x010F54},
		{0x010FC5, 0x010FCB},
		{0x011052, 0x01106F},
		{0x0110F0, 0x0110F9},
		{0x011136, 0x01113F},
		{0x0111D0, 0x0111D9},
		{0x0111E1, 0x0111F4},
		{0x0112F0, 0x0112F9},
		{0x011450, 0x011459},
		{0x0114D0, 0x0114D9},
		{0x011650, 0x011659},
		{0x0116C0, 0x0116C9},
		{0x011730, 0x01173B},
		{0x0118E0, 0x0118F2},
		{0x011950, 0x011959},
		{0x011C50, 0x011C6C},
		{0x011D50, 0x011D59},
		{0x011DA0, 0x011DA9},
		{0x011F50, 0x011F59},
		{0x011FC0, 0x011FD4},
		{0x012400, 0x01246E},
		{0x016A60, 0x016A69},
		{0x016AC0, 0x016AC9},
		{0x016B50, 0x016B59},
		{0x016B5B, 0x016B61},
		{0x016E80, 0x016E96},
		{0x01D2C0, 0x01D2D3},
		{0x01D2E0, 0x01D2F3},
		{0x01D360, 0x01D378},
		{0x01D7CE, 0x01D7FF},
		{0x01E140, 0x01E149},
		{0x01E2F0, 0x01E2F9},
		{0x01E4F0, 0x01E4F9},
		{0x01E8C7, 0x01E8CF},
		{0x01E950, 0x01E959},
		{0x01EC71, 0x01ECAB},
		{0x01ECAD, 0x01ECAF},
		{0x01ECB1, 0x01ECB4},
		{0x01ED01, 0x01ED2D},
		{0x01ED2F, 0x01ED3D},
		{0x01F100, 0x01F10C},
		{0x01FBF0, 0x01FBF9},
	},

	"Unicode": {
		{0x00, 0x10FFFF},
	},
}

type MapperTest struct {
	name           string
	r              combinator.Result
	expectedResult combinator.Result
	expectedOK     bool
	expectedError  string
}

func intPtr(v int) *int {
	return &v
}

func assertEqualResults(t *testing.T, expected, actual combinator.Result) {
	expectedNFA, ok := expected.Val.(*automata.NFA)
	if !ok {
		assert.Equal(t, expected, actual)
		return
	}

	nfa, ok := actual.Val.(*automata.NFA)
	if !ok {
		assert.Equal(t, expected, actual)
		return
	}

	assert.True(t, nfa.Equal(expectedNFA), "Expected NFA:\n%s\nGot:\n%s\n", expectedNFA, nfa)
	assert.Equal(t, expected.Pos, actual.Pos, "Expected Pos:\n%d\nGot:\n%d\n", expected.Pos, actual.Pos)
	assert.Equal(t, expected.Bag, actual.Bag, "Expected Bag:\n%v\nGot:\n%v\n", expected.Bag, actual.Bag)
}

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		regex         string
		expectedError string
		expectedNFA   *automata.NFA
	}{
		{
			name:          "InvalidRegex",
			regex:         "[",
			expectedError: "invalid regular expression: [",
		},
		{
			name:          "InvalidCharRange",
			regex:         "[9-0]",
			expectedError: "invalid character range 9-0",
		},
		{
			name:          "InvalidRepRange",
			regex:         "[0-9]{4,2}",
			expectedError: "invalid repetition range {4,2}",
		},
		{
			name:  "Success_EscapedChars",
			regex: `\n|\r|\r\n`,
			expectedNFA: testNFA["\n"].Union(
				testNFA["\r"].Union(
					testNFA["\r"].Concat(testNFA["\n"]),
				),
			),
		},
		{
			name:  "Success_CharRanges",
			regex: `^[A-Z]?[a-z][0-9A-Za-z]{1,}$`,
			expectedNFA: empty().Union(testNFA["upper"]).Concat(
				testNFA["lower"],
				testNFA["alnum"].Concat(testNFA["alnum"].Star()),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nfa, err := Parse(tc.regex)

			if tc.expectedError != "" {
				assert.Nil(t, nfa)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nfa)
				assert.True(t, nfa.Equal(tc.expectedNFA))
			}
		})
	}
}

func TestMappers_ToAnyChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '.',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["Unicode"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnyChar(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSingleChar(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: 'x',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: char.RangeList{
						{'x', 'x'},
					},
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSingleChar(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Whitespace",
			r: combinator.Result{
				Val: `\s`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["ws"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["ws"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWhitespace",
			r: combinator.Result{
				Val: `\S`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["not_ws"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_ws"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Digit",
			r: combinator.Result{
				Val: `\d`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotDigit",
			r: combinator.Result{
				Val: `\D`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["not_digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: combinator.Result{
				Val: `\w`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["word"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_NotWord",
			r: combinator.Result{
				Val: `\W`,
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["not_word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["not_word"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharClass(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToASCIICharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Blank",
			r: combinator.Result{
				Val: "[:blank:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["blank"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["blank"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Space",
			r: combinator.Result{
				Val: "[:space:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["space"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["space"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Digit",
			r: combinator.Result{
				Val: "[:digit:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_XDigit",
			r: combinator.Result{
				Val: "[:xdigit:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["xdigit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["xdigit"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Upper",
			r: combinator.Result{
				Val: "[:upper:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["upper"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["upper"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Lower",
			r: combinator.Result{
				Val: "[:lower:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["lower"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["lower"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alpha",
			r: combinator.Result{
				Val: "[:alpha:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["alpha"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["alpha"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Alnum",
			r: combinator.Result{
				Val: "[:alnum:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["alnum"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["alnum"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_Word",
			r: combinator.Result{
				Val: "[:word:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["word"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["word"],
				},
			},
			expectedOK: true,
		},
		{
			name: "Success_ASCII",
			r: combinator.Result{
				Val: "[:ascii:]",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["ascii"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["ascii"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToASCIICharClass(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUnicodeCategory(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: "Letter",
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: "Letter",
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUnicodeCategory(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUnicodeCharClass(t *testing.T) {
	tests := []MapperTest{
		{
			name: "InvalidClass",
			r: combinator.Result{
				Val: combinator.List{
					{Val: `\p`, Pos: 2},
					{Val: '{', Pos: 4},
					{Val: "Runic", Pos: 5},
					{Val: '}', Pos: 11},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{},
			expectedOK:     false,
		},
		{
			name: "Success_Number",
			r: combinator.Result{
				Val: combinator.List{
					{Val: `\p`, Pos: 2},
					{Val: '{', Pos: 4},
					{Val: "Number", Pos: 5},
					{Val: '}', Pos: 11},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["Number"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["Number"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUnicodeCharClass(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepOp(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepOp(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToUpperBound(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Unbounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: ',', Pos: 2},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: (*int)(nil),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: ',', Pos: 2},
					{Val: 4, Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: intPtr(4),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToUpperBound(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_Fixed",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{Val: combinator.Empty{}},
					{Val: '}', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Unbounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: (*int)(nil),
						Pos: 4,
					},
					{Val: '}', Pos: 5},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: (*int)(nil),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Bounded",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 2, Pos: 3},
					{
						Val: intPtr(6),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 2,
					q: intPtr(6),
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '{', Pos: 2},
					{Val: 6, Pos: 3},
					{
						Val: intPtr(2),
						Pos: 4,
					},
					{Val: '}', Pos: 6},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[int, *int]{
					p: 6,
					q: intPtr(2),
				},
				Pos: 2,
			},
			expectedOK:    true,
			expectedError: "invalid repetition range {6,2}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRange(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRepetition(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: '*',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRepetition(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToQuantifier(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success_NonLazy",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '*', Pos: 2},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: false,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '*', Pos: 2},
					{Val: '?', Pos: 3},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: tuple[any, bool]{
					p: '*',
					q: true,
				},
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToQuantifier(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharInRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: 'a',
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharInRange(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharRange(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: 'a', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'f', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: automata.NewNFABuilder().
					SetStart(0).
					SetFinal([]automata.State{1}).
					AddTransition(0, 'a', 'f', []automata.State{1}).
					Build(),
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: char.RangeList{
						{'a', 'f'},
					},
				},
			},
			expectedOK: true,
		},
		{
			name: "InvalidRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: 'f', Pos: 2},
					{Val: '-', Pos: 3},
					{Val: 'a', Pos: 4},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: nil,
				Pos: 2,
				Bag: nil,
			},
			expectedOK:    true,
			expectedError: "invalid character range f-a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharRange(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroupItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroupItem(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToCharGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '[', Pos: 2},
					{Val: combinator.Empty{}},
					{
						Val: combinator.List{
							{
								Val: testNFA["xdigit"],
								Pos: 3,
								Bag: combinator.Bag{
									bagKeyCharRanges: testRanges["xdigit"],
								},
							},
							{
								Val: automata.NewNFABuilder().
									SetStart(0).
									SetFinal([]automata.State{1}).
									AddTransition(0, '-', '-', []automata.State{1}).
									Build(),
								Pos: 12,
								Bag: combinator.Bag{
									bagKeyCharRanges: char.RangeList{
										{'-', '-'},
									},
								},
							},
						},
						Pos: 3,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: automata.NewNFABuilder().
					SetStart(0).
					SetFinal([]automata.State{1}).
					AddTransition(0, '-', '-', []automata.State{1}).
					AddTransition(0, '0', '9', []automata.State{1}).
					AddTransition(0, 'A', 'F', []automata.State{1}).
					AddTransition(0, 'a', 'f', []automata.State{1}).
					Build(),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Negated",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '[', Pos: 2},
					{Val: '^', Pos: 3},
					{
						Val: combinator.List{
							{
								Val: testNFA["alnum"],
								Pos: 4,
								Bag: combinator.Bag{
									bagKeyCharRanges: testRanges["alnum"],
								},
							},
						},
						Pos: 4,
					},
					{Val: ']', Pos: 13},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: automata.NewNFABuilder().
					SetStart(0).
					SetFinal([]automata.State{1}).
					AddTransition(0, 0x00, 0x2F, []automata.State{1}).
					AddTransition(0, 0x3A, 0x40, []automata.State{1}).
					AddTransition(0, 0x5B, 0x60, []automata.State{1}).
					AddTransition(0, 0x7B, 0x10FFFF, []automata.State{1}).
					Build(),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToCharGroup(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatchItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyCharRanges: testRanges["digit"],
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatchItem(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToMatch(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"],
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: empty().Union(testNFA["x"]),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Star(),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"].Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"]),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"], testNFA["x"].Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(
					testNFA["x"],
					empty().Union(testNFA["x"]),
					empty().Union(testNFA["x"]),
				),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["x"],
						Pos: 2,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"].Star()),
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToMatch(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToGroup(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"],
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrOne",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '?',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: empty().Union(testNFA["x"]),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_ZeroOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '*',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Star(),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"].Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_FixedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(2),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"]),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_UnboundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: (*int)(nil),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"], testNFA["x"].Star()),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_BoundedRange",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: tuple[int, *int]{
								p: 2,
								q: intPtr(4),
							},
							q: false,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(
					testNFA["x"],
					empty().Union(testNFA["x"]),
					empty().Union(testNFA["x"]),
				),
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_Lazy_OneOrMany",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '(', Pos: 2},
					{
						Val: testNFA["x"],
						Pos: 3,
						Bag: combinator.Bag{
							bagKeyCharRanges: []rune{'x'},
						},
					},
					{Val: ')', Pos: 4},
					{
						Val: tuple[any, bool]{
							p: '+',
							q: true,
						},
						Pos: 5,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["x"].Concat(testNFA["x"].Star()),
				Pos: 2,
				Bag: combinator.Bag{
					bagKeyLazyQuantifier: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToGroup(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToAnchor(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: '$',
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: EndOfString,
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToAnchor(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexprItem(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexprItem(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToSubexpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["digit"],
						Pos: 2,
					},
					{
						Val: EndOfString,
						Pos: 3,
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToSubexpr(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToExpr(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["upper"],
						Pos: 2,
					},
					{Val: combinator.Empty{}},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["upper"],
				Pos: 2,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithExpr",
			r: combinator.Result{
				Val: combinator.List{
					{
						Val: testNFA["upper"],
						Pos: 2,
					},
					{
						Val: combinator.List{
							{Val: '|', Pos: 3},
							{
								Val: testNFA["lower"],
								Pos: 4,
							},
						},
					},
				},
				Pos: 2,
			},
			expectedResult: combinator.Result{
				Val: testNFA["upper"].Union(testNFA["lower"]),
				Pos: 2,
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToExpr(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}

func TestMappers_ToRegex(t *testing.T) {
	tests := []MapperTest{
		{
			name: "Success",
			r: combinator.Result{
				Val: combinator.List{
					{Val: combinator.Empty{}},
					{
						Val: testNFA["digit"],
						Pos: 0,
					},
				},
				Pos: 0,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 0,
			},
			expectedOK: true,
		},
		{
			name: "Success_WithStartOfString",
			r: combinator.Result{
				Val: combinator.List{
					{Val: '^', Pos: 0},
					{
						Val: testNFA["digit"],
						Pos: 1,
					},
				},
				Pos: 0,
			},
			expectedResult: combinator.Result{
				Val: testNFA["digit"],
				Pos: 0,
				Bag: combinator.Bag{
					BagKeyStartOfString: true,
				},
			},
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := new(mappers)
			res, ok := m.ToRegex(tc.r)

			assertEqualResults(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectedError != "" {
				assert.EqualError(t, m.errors, tc.expectedError)
			}
		})
	}
}
