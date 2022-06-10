package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gardenbed/emerge/internal/regex/ast"
)

var (
	digit = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: '0'},
			&ast.Char{Val: '1'},
			&ast.Char{Val: '2'},
			&ast.Char{Val: '3'},
			&ast.Char{Val: '4'},
			&ast.Char{Val: '5'},
			&ast.Char{Val: '6'},
			&ast.Char{Val: '7'},
			&ast.Char{Val: '8'},
			&ast.Char{Val: '9'},
		},
	}

	nonDigit = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 0},
			&ast.Char{Val: 1},
			&ast.Char{Val: 2},
			&ast.Char{Val: 3},
			&ast.Char{Val: 4},
			&ast.Char{Val: 5},
			&ast.Char{Val: 6},
			&ast.Char{Val: 7},
			&ast.Char{Val: 8},
			&ast.Char{Val: 9},
			&ast.Char{Val: 10},
			&ast.Char{Val: 11},
			&ast.Char{Val: 12},
			&ast.Char{Val: 13},
			&ast.Char{Val: 14},
			&ast.Char{Val: 15},
			&ast.Char{Val: 16},
			&ast.Char{Val: 17},
			&ast.Char{Val: 18},
			&ast.Char{Val: 19},
			&ast.Char{Val: 20},
			&ast.Char{Val: 21},
			&ast.Char{Val: 22},
			&ast.Char{Val: 23},
			&ast.Char{Val: 24},
			&ast.Char{Val: 25},
			&ast.Char{Val: 26},
			&ast.Char{Val: 27},
			&ast.Char{Val: 28},
			&ast.Char{Val: 29},
			&ast.Char{Val: 30},
			&ast.Char{Val: 31},
			&ast.Char{Val: 32},
			&ast.Char{Val: 33},
			&ast.Char{Val: 34},
			&ast.Char{Val: 35},
			&ast.Char{Val: 36},
			&ast.Char{Val: 37},
			&ast.Char{Val: 38},
			&ast.Char{Val: 39},
			&ast.Char{Val: 40},
			&ast.Char{Val: 41},
			&ast.Char{Val: 42},
			&ast.Char{Val: 43},
			&ast.Char{Val: 44},
			&ast.Char{Val: 45},
			&ast.Char{Val: 46},
			&ast.Char{Val: 47},
			&ast.Char{Val: 58},
			&ast.Char{Val: 59},
			&ast.Char{Val: 60},
			&ast.Char{Val: 61},
			&ast.Char{Val: 62},
			&ast.Char{Val: 63},
			&ast.Char{Val: 64},
			&ast.Char{Val: 65},
			&ast.Char{Val: 66},
			&ast.Char{Val: 67},
			&ast.Char{Val: 68},
			&ast.Char{Val: 69},
			&ast.Char{Val: 70},
			&ast.Char{Val: 71},
			&ast.Char{Val: 72},
			&ast.Char{Val: 73},
			&ast.Char{Val: 74},
			&ast.Char{Val: 75},
			&ast.Char{Val: 76},
			&ast.Char{Val: 77},
			&ast.Char{Val: 78},
			&ast.Char{Val: 79},
			&ast.Char{Val: 80},
			&ast.Char{Val: 81},
			&ast.Char{Val: 82},
			&ast.Char{Val: 83},
			&ast.Char{Val: 84},
			&ast.Char{Val: 85},
			&ast.Char{Val: 86},
			&ast.Char{Val: 87},
			&ast.Char{Val: 88},
			&ast.Char{Val: 89},
			&ast.Char{Val: 90},
			&ast.Char{Val: 91},
			&ast.Char{Val: 92},
			&ast.Char{Val: 93},
			&ast.Char{Val: 94},
			&ast.Char{Val: 95},
			&ast.Char{Val: 96},
			&ast.Char{Val: 97},
			&ast.Char{Val: 98},
			&ast.Char{Val: 99},
			&ast.Char{Val: 100},
			&ast.Char{Val: 101},
			&ast.Char{Val: 102},
			&ast.Char{Val: 103},
			&ast.Char{Val: 104},
			&ast.Char{Val: 105},
			&ast.Char{Val: 106},
			&ast.Char{Val: 107},
			&ast.Char{Val: 108},
			&ast.Char{Val: 109},
			&ast.Char{Val: 110},
			&ast.Char{Val: 111},
			&ast.Char{Val: 112},
			&ast.Char{Val: 113},
			&ast.Char{Val: 114},
			&ast.Char{Val: 115},
			&ast.Char{Val: 116},
			&ast.Char{Val: 117},
			&ast.Char{Val: 118},
			&ast.Char{Val: 119},
			&ast.Char{Val: 120},
			&ast.Char{Val: 121},
			&ast.Char{Val: 122},
			&ast.Char{Val: 123},
			&ast.Char{Val: 124},
			&ast.Char{Val: 125},
			&ast.Char{Val: 126},
			&ast.Char{Val: 127},
		},
	}

	whitespace = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: ' '},
			&ast.Char{Val: '\t'},
			&ast.Char{Val: '\n'},
			&ast.Char{Val: '\r'},
			&ast.Char{Val: '\f'},
		},
	}

	nonWhitespace = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 0},
			&ast.Char{Val: 1},
			&ast.Char{Val: 2},
			&ast.Char{Val: 3},
			&ast.Char{Val: 4},
			&ast.Char{Val: 5},
			&ast.Char{Val: 6},
			&ast.Char{Val: 7},
			&ast.Char{Val: 8},
			&ast.Char{Val: 11},
			&ast.Char{Val: 14},
			&ast.Char{Val: 15},
			&ast.Char{Val: 16},
			&ast.Char{Val: 17},
			&ast.Char{Val: 18},
			&ast.Char{Val: 19},
			&ast.Char{Val: 20},
			&ast.Char{Val: 21},
			&ast.Char{Val: 22},
			&ast.Char{Val: 23},
			&ast.Char{Val: 24},
			&ast.Char{Val: 25},
			&ast.Char{Val: 26},
			&ast.Char{Val: 27},
			&ast.Char{Val: 28},
			&ast.Char{Val: 29},
			&ast.Char{Val: 30},
			&ast.Char{Val: 31},
			&ast.Char{Val: 33},
			&ast.Char{Val: 34},
			&ast.Char{Val: 35},
			&ast.Char{Val: 36},
			&ast.Char{Val: 37},
			&ast.Char{Val: 38},
			&ast.Char{Val: 39},
			&ast.Char{Val: 40},
			&ast.Char{Val: 41},
			&ast.Char{Val: 42},
			&ast.Char{Val: 43},
			&ast.Char{Val: 44},
			&ast.Char{Val: 45},
			&ast.Char{Val: 46},
			&ast.Char{Val: 47},
			&ast.Char{Val: 48},
			&ast.Char{Val: 49},
			&ast.Char{Val: 50},
			&ast.Char{Val: 51},
			&ast.Char{Val: 52},
			&ast.Char{Val: 53},
			&ast.Char{Val: 54},
			&ast.Char{Val: 55},
			&ast.Char{Val: 56},
			&ast.Char{Val: 57},
			&ast.Char{Val: 58},
			&ast.Char{Val: 59},
			&ast.Char{Val: 60},
			&ast.Char{Val: 61},
			&ast.Char{Val: 62},
			&ast.Char{Val: 63},
			&ast.Char{Val: 64},
			&ast.Char{Val: 65},
			&ast.Char{Val: 66},
			&ast.Char{Val: 67},
			&ast.Char{Val: 68},
			&ast.Char{Val: 69},
			&ast.Char{Val: 70},
			&ast.Char{Val: 71},
			&ast.Char{Val: 72},
			&ast.Char{Val: 73},
			&ast.Char{Val: 74},
			&ast.Char{Val: 75},
			&ast.Char{Val: 76},
			&ast.Char{Val: 77},
			&ast.Char{Val: 78},
			&ast.Char{Val: 79},
			&ast.Char{Val: 80},
			&ast.Char{Val: 81},
			&ast.Char{Val: 82},
			&ast.Char{Val: 83},
			&ast.Char{Val: 84},
			&ast.Char{Val: 85},
			&ast.Char{Val: 86},
			&ast.Char{Val: 87},
			&ast.Char{Val: 88},
			&ast.Char{Val: 89},
			&ast.Char{Val: 90},
			&ast.Char{Val: 91},
			&ast.Char{Val: 92},
			&ast.Char{Val: 93},
			&ast.Char{Val: 94},
			&ast.Char{Val: 95},
			&ast.Char{Val: 96},
			&ast.Char{Val: 97},
			&ast.Char{Val: 98},
			&ast.Char{Val: 99},
			&ast.Char{Val: 100},
			&ast.Char{Val: 101},
			&ast.Char{Val: 102},
			&ast.Char{Val: 103},
			&ast.Char{Val: 104},
			&ast.Char{Val: 105},
			&ast.Char{Val: 106},
			&ast.Char{Val: 107},
			&ast.Char{Val: 108},
			&ast.Char{Val: 109},
			&ast.Char{Val: 110},
			&ast.Char{Val: 111},
			&ast.Char{Val: 112},
			&ast.Char{Val: 113},
			&ast.Char{Val: 114},
			&ast.Char{Val: 115},
			&ast.Char{Val: 116},
			&ast.Char{Val: 117},
			&ast.Char{Val: 118},
			&ast.Char{Val: 119},
			&ast.Char{Val: 120},
			&ast.Char{Val: 121},
			&ast.Char{Val: 122},
			&ast.Char{Val: 123},
			&ast.Char{Val: 124},
			&ast.Char{Val: 125},
			&ast.Char{Val: 126},
			&ast.Char{Val: 127},
		},
	}

	word = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: '0'},
			&ast.Char{Val: '1'},
			&ast.Char{Val: '2'},
			&ast.Char{Val: '3'},
			&ast.Char{Val: '4'},
			&ast.Char{Val: '5'},
			&ast.Char{Val: '6'},
			&ast.Char{Val: '7'},
			&ast.Char{Val: '8'},
			&ast.Char{Val: '9'},
			&ast.Char{Val: 'A'},
			&ast.Char{Val: 'B'},
			&ast.Char{Val: 'C'},
			&ast.Char{Val: 'D'},
			&ast.Char{Val: 'E'},
			&ast.Char{Val: 'F'},
			&ast.Char{Val: 'G'},
			&ast.Char{Val: 'H'},
			&ast.Char{Val: 'I'},
			&ast.Char{Val: 'J'},
			&ast.Char{Val: 'K'},
			&ast.Char{Val: 'L'},
			&ast.Char{Val: 'M'},
			&ast.Char{Val: 'N'},
			&ast.Char{Val: 'O'},
			&ast.Char{Val: 'P'},
			&ast.Char{Val: 'Q'},
			&ast.Char{Val: 'R'},
			&ast.Char{Val: 'S'},
			&ast.Char{Val: 'T'},
			&ast.Char{Val: 'U'},
			&ast.Char{Val: 'V'},
			&ast.Char{Val: 'W'},
			&ast.Char{Val: 'X'},
			&ast.Char{Val: 'Y'},
			&ast.Char{Val: 'Z'},
			&ast.Char{Val: '_'},
			&ast.Char{Val: 'a'},
			&ast.Char{Val: 'b'},
			&ast.Char{Val: 'c'},
			&ast.Char{Val: 'd'},
			&ast.Char{Val: 'e'},
			&ast.Char{Val: 'f'},
			&ast.Char{Val: 'g'},
			&ast.Char{Val: 'h'},
			&ast.Char{Val: 'i'},
			&ast.Char{Val: 'j'},
			&ast.Char{Val: 'k'},
			&ast.Char{Val: 'l'},
			&ast.Char{Val: 'm'},
			&ast.Char{Val: 'n'},
			&ast.Char{Val: 'o'},
			&ast.Char{Val: 'p'},
			&ast.Char{Val: 'q'},
			&ast.Char{Val: 'r'},
			&ast.Char{Val: 's'},
			&ast.Char{Val: 't'},
			&ast.Char{Val: 'u'},
			&ast.Char{Val: 'v'},
			&ast.Char{Val: 'w'},
			&ast.Char{Val: 'x'},
			&ast.Char{Val: 'y'},
			&ast.Char{Val: 'z'},
		},
	}

	nonWord = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 0},
			&ast.Char{Val: 1},
			&ast.Char{Val: 2},
			&ast.Char{Val: 3},
			&ast.Char{Val: 4},
			&ast.Char{Val: 5},
			&ast.Char{Val: 6},
			&ast.Char{Val: 7},
			&ast.Char{Val: 8},
			&ast.Char{Val: 9},
			&ast.Char{Val: 10},
			&ast.Char{Val: 11},
			&ast.Char{Val: 12},
			&ast.Char{Val: 13},
			&ast.Char{Val: 14},
			&ast.Char{Val: 15},
			&ast.Char{Val: 16},
			&ast.Char{Val: 17},
			&ast.Char{Val: 18},
			&ast.Char{Val: 19},
			&ast.Char{Val: 20},
			&ast.Char{Val: 21},
			&ast.Char{Val: 22},
			&ast.Char{Val: 23},
			&ast.Char{Val: 24},
			&ast.Char{Val: 25},
			&ast.Char{Val: 26},
			&ast.Char{Val: 27},
			&ast.Char{Val: 28},
			&ast.Char{Val: 29},
			&ast.Char{Val: 30},
			&ast.Char{Val: 31},
			&ast.Char{Val: 32},
			&ast.Char{Val: 33},
			&ast.Char{Val: 34},
			&ast.Char{Val: 35},
			&ast.Char{Val: 36},
			&ast.Char{Val: 37},
			&ast.Char{Val: 38},
			&ast.Char{Val: 39},
			&ast.Char{Val: 40},
			&ast.Char{Val: 41},
			&ast.Char{Val: 42},
			&ast.Char{Val: 43},
			&ast.Char{Val: 44},
			&ast.Char{Val: 45},
			&ast.Char{Val: 46},
			&ast.Char{Val: 47},
			&ast.Char{Val: 58},
			&ast.Char{Val: 59},
			&ast.Char{Val: 60},
			&ast.Char{Val: 61},
			&ast.Char{Val: 62},
			&ast.Char{Val: 63},
			&ast.Char{Val: 64},
			&ast.Char{Val: 91},
			&ast.Char{Val: 92},
			&ast.Char{Val: 93},
			&ast.Char{Val: 94},
			&ast.Char{Val: 96},
			&ast.Char{Val: 123},
			&ast.Char{Val: 124},
			&ast.Char{Val: 125},
			&ast.Char{Val: 126},
			&ast.Char{Val: 127},
		},
	}

	xdigit = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: '0'},
			&ast.Char{Val: '1'},
			&ast.Char{Val: '2'},
			&ast.Char{Val: '3'},
			&ast.Char{Val: '4'},
			&ast.Char{Val: '5'},
			&ast.Char{Val: '6'},
			&ast.Char{Val: '7'},
			&ast.Char{Val: '8'},
			&ast.Char{Val: '9'},
			&ast.Char{Val: 'A'},
			&ast.Char{Val: 'B'},
			&ast.Char{Val: 'C'},
			&ast.Char{Val: 'D'},
			&ast.Char{Val: 'E'},
			&ast.Char{Val: 'F'},
			&ast.Char{Val: 'a'},
			&ast.Char{Val: 'b'},
			&ast.Char{Val: 'c'},
			&ast.Char{Val: 'd'},
			&ast.Char{Val: 'e'},
			&ast.Char{Val: 'f'},
		},
	}

	upper = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 'A'},
			&ast.Char{Val: 'B'},
			&ast.Char{Val: 'C'},
			&ast.Char{Val: 'D'},
			&ast.Char{Val: 'E'},
			&ast.Char{Val: 'F'},
			&ast.Char{Val: 'G'},
			&ast.Char{Val: 'H'},
			&ast.Char{Val: 'I'},
			&ast.Char{Val: 'J'},
			&ast.Char{Val: 'K'},
			&ast.Char{Val: 'L'},
			&ast.Char{Val: 'M'},
			&ast.Char{Val: 'N'},
			&ast.Char{Val: 'O'},
			&ast.Char{Val: 'P'},
			&ast.Char{Val: 'Q'},
			&ast.Char{Val: 'R'},
			&ast.Char{Val: 'S'},
			&ast.Char{Val: 'T'},
			&ast.Char{Val: 'U'},
			&ast.Char{Val: 'V'},
			&ast.Char{Val: 'W'},
			&ast.Char{Val: 'X'},
			&ast.Char{Val: 'Y'},
			&ast.Char{Val: 'Z'},
		},
	}

	lower = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 'a'},
			&ast.Char{Val: 'b'},
			&ast.Char{Val: 'c'},
			&ast.Char{Val: 'd'},
			&ast.Char{Val: 'e'},
			&ast.Char{Val: 'f'},
			&ast.Char{Val: 'g'},
			&ast.Char{Val: 'h'},
			&ast.Char{Val: 'i'},
			&ast.Char{Val: 'j'},
			&ast.Char{Val: 'k'},
			&ast.Char{Val: 'l'},
			&ast.Char{Val: 'm'},
			&ast.Char{Val: 'n'},
			&ast.Char{Val: 'o'},
			&ast.Char{Val: 'p'},
			&ast.Char{Val: 'q'},
			&ast.Char{Val: 'r'},
			&ast.Char{Val: 's'},
			&ast.Char{Val: 't'},
			&ast.Char{Val: 'u'},
			&ast.Char{Val: 'v'},
			&ast.Char{Val: 'w'},
			&ast.Char{Val: 'x'},
			&ast.Char{Val: 'y'},
			&ast.Char{Val: 'z'},
		},
	}

	alpha = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 'A'},
			&ast.Char{Val: 'B'},
			&ast.Char{Val: 'C'},
			&ast.Char{Val: 'D'},
			&ast.Char{Val: 'E'},
			&ast.Char{Val: 'F'},
			&ast.Char{Val: 'G'},
			&ast.Char{Val: 'H'},
			&ast.Char{Val: 'I'},
			&ast.Char{Val: 'J'},
			&ast.Char{Val: 'K'},
			&ast.Char{Val: 'L'},
			&ast.Char{Val: 'M'},
			&ast.Char{Val: 'N'},
			&ast.Char{Val: 'O'},
			&ast.Char{Val: 'P'},
			&ast.Char{Val: 'Q'},
			&ast.Char{Val: 'R'},
			&ast.Char{Val: 'S'},
			&ast.Char{Val: 'T'},
			&ast.Char{Val: 'U'},
			&ast.Char{Val: 'V'},
			&ast.Char{Val: 'W'},
			&ast.Char{Val: 'X'},
			&ast.Char{Val: 'Y'},
			&ast.Char{Val: 'Z'},
			&ast.Char{Val: 'a'},
			&ast.Char{Val: 'b'},
			&ast.Char{Val: 'c'},
			&ast.Char{Val: 'd'},
			&ast.Char{Val: 'e'},
			&ast.Char{Val: 'f'},
			&ast.Char{Val: 'g'},
			&ast.Char{Val: 'h'},
			&ast.Char{Val: 'i'},
			&ast.Char{Val: 'j'},
			&ast.Char{Val: 'k'},
			&ast.Char{Val: 'l'},
			&ast.Char{Val: 'm'},
			&ast.Char{Val: 'n'},
			&ast.Char{Val: 'o'},
			&ast.Char{Val: 'p'},
			&ast.Char{Val: 'q'},
			&ast.Char{Val: 'r'},
			&ast.Char{Val: 's'},
			&ast.Char{Val: 't'},
			&ast.Char{Val: 'u'},
			&ast.Char{Val: 'v'},
			&ast.Char{Val: 'w'},
			&ast.Char{Val: 'x'},
			&ast.Char{Val: 'y'},
			&ast.Char{Val: 'z'},
		},
	}

	alnum = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: '0'},
			&ast.Char{Val: '1'},
			&ast.Char{Val: '2'},
			&ast.Char{Val: '3'},
			&ast.Char{Val: '4'},
			&ast.Char{Val: '5'},
			&ast.Char{Val: '6'},
			&ast.Char{Val: '7'},
			&ast.Char{Val: '8'},
			&ast.Char{Val: '9'},
			&ast.Char{Val: 'A'},
			&ast.Char{Val: 'B'},
			&ast.Char{Val: 'C'},
			&ast.Char{Val: 'D'},
			&ast.Char{Val: 'E'},
			&ast.Char{Val: 'F'},
			&ast.Char{Val: 'G'},
			&ast.Char{Val: 'H'},
			&ast.Char{Val: 'I'},
			&ast.Char{Val: 'J'},
			&ast.Char{Val: 'K'},
			&ast.Char{Val: 'L'},
			&ast.Char{Val: 'M'},
			&ast.Char{Val: 'N'},
			&ast.Char{Val: 'O'},
			&ast.Char{Val: 'P'},
			&ast.Char{Val: 'Q'},
			&ast.Char{Val: 'R'},
			&ast.Char{Val: 'S'},
			&ast.Char{Val: 'T'},
			&ast.Char{Val: 'U'},
			&ast.Char{Val: 'V'},
			&ast.Char{Val: 'W'},
			&ast.Char{Val: 'X'},
			&ast.Char{Val: 'Y'},
			&ast.Char{Val: 'Z'},
			&ast.Char{Val: 'a'},
			&ast.Char{Val: 'b'},
			&ast.Char{Val: 'c'},
			&ast.Char{Val: 'd'},
			&ast.Char{Val: 'e'},
			&ast.Char{Val: 'f'},
			&ast.Char{Val: 'g'},
			&ast.Char{Val: 'h'},
			&ast.Char{Val: 'i'},
			&ast.Char{Val: 'j'},
			&ast.Char{Val: 'k'},
			&ast.Char{Val: 'l'},
			&ast.Char{Val: 'm'},
			&ast.Char{Val: 'n'},
			&ast.Char{Val: 'o'},
			&ast.Char{Val: 'p'},
			&ast.Char{Val: 'q'},
			&ast.Char{Val: 'r'},
			&ast.Char{Val: 's'},
			&ast.Char{Val: 't'},
			&ast.Char{Val: 'u'},
			&ast.Char{Val: 'v'},
			&ast.Char{Val: 'w'},
			&ast.Char{Val: 'x'},
			&ast.Char{Val: 'y'},
			&ast.Char{Val: 'z'},
		},
	}

	ascii = &ast.Alt{
		Exprs: []ast.Node{
			&ast.Char{Val: 0},
			&ast.Char{Val: 1},
			&ast.Char{Val: 2},
			&ast.Char{Val: 3},
			&ast.Char{Val: 4},
			&ast.Char{Val: 5},
			&ast.Char{Val: 6},
			&ast.Char{Val: 7},
			&ast.Char{Val: 8},
			&ast.Char{Val: 9},
			&ast.Char{Val: 10},
			&ast.Char{Val: 11},
			&ast.Char{Val: 12},
			&ast.Char{Val: 13},
			&ast.Char{Val: 14},
			&ast.Char{Val: 15},
			&ast.Char{Val: 16},
			&ast.Char{Val: 17},
			&ast.Char{Val: 18},
			&ast.Char{Val: 19},
			&ast.Char{Val: 20},
			&ast.Char{Val: 21},
			&ast.Char{Val: 22},
			&ast.Char{Val: 23},
			&ast.Char{Val: 24},
			&ast.Char{Val: 25},
			&ast.Char{Val: 26},
			&ast.Char{Val: 27},
			&ast.Char{Val: 28},
			&ast.Char{Val: 29},
			&ast.Char{Val: 30},
			&ast.Char{Val: 31},
			&ast.Char{Val: 32},
			&ast.Char{Val: 33},
			&ast.Char{Val: 34},
			&ast.Char{Val: 35},
			&ast.Char{Val: 36},
			&ast.Char{Val: 37},
			&ast.Char{Val: 38},
			&ast.Char{Val: 39},
			&ast.Char{Val: 40},
			&ast.Char{Val: 41},
			&ast.Char{Val: 42},
			&ast.Char{Val: 43},
			&ast.Char{Val: 44},
			&ast.Char{Val: 45},
			&ast.Char{Val: 46},
			&ast.Char{Val: 47},
			&ast.Char{Val: 48},
			&ast.Char{Val: 49},
			&ast.Char{Val: 50},
			&ast.Char{Val: 51},
			&ast.Char{Val: 52},
			&ast.Char{Val: 53},
			&ast.Char{Val: 54},
			&ast.Char{Val: 55},
			&ast.Char{Val: 56},
			&ast.Char{Val: 57},
			&ast.Char{Val: 58},
			&ast.Char{Val: 59},
			&ast.Char{Val: 60},
			&ast.Char{Val: 61},
			&ast.Char{Val: 62},
			&ast.Char{Val: 63},
			&ast.Char{Val: 64},
			&ast.Char{Val: 65},
			&ast.Char{Val: 66},
			&ast.Char{Val: 67},
			&ast.Char{Val: 68},
			&ast.Char{Val: 69},
			&ast.Char{Val: 70},
			&ast.Char{Val: 71},
			&ast.Char{Val: 72},
			&ast.Char{Val: 73},
			&ast.Char{Val: 74},
			&ast.Char{Val: 75},
			&ast.Char{Val: 76},
			&ast.Char{Val: 77},
			&ast.Char{Val: 78},
			&ast.Char{Val: 79},
			&ast.Char{Val: 80},
			&ast.Char{Val: 81},
			&ast.Char{Val: 82},
			&ast.Char{Val: 83},
			&ast.Char{Val: 84},
			&ast.Char{Val: 85},
			&ast.Char{Val: 86},
			&ast.Char{Val: 87},
			&ast.Char{Val: 88},
			&ast.Char{Val: 89},
			&ast.Char{Val: 90},
			&ast.Char{Val: 91},
			&ast.Char{Val: 92},
			&ast.Char{Val: 93},
			&ast.Char{Val: 94},
			&ast.Char{Val: 95},
			&ast.Char{Val: 96},
			&ast.Char{Val: 97},
			&ast.Char{Val: 98},
			&ast.Char{Val: 99},
			&ast.Char{Val: 100},
			&ast.Char{Val: 101},
			&ast.Char{Val: 102},
			&ast.Char{Val: 103},
			&ast.Char{Val: 104},
			&ast.Char{Val: 105},
			&ast.Char{Val: 106},
			&ast.Char{Val: 107},
			&ast.Char{Val: 108},
			&ast.Char{Val: 109},
			&ast.Char{Val: 110},
			&ast.Char{Val: 111},
			&ast.Char{Val: 112},
			&ast.Char{Val: 113},
			&ast.Char{Val: 114},
			&ast.Char{Val: 115},
			&ast.Char{Val: 116},
			&ast.Char{Val: 117},
			&ast.Char{Val: 118},
			&ast.Char{Val: 119},
			&ast.Char{Val: 120},
			&ast.Char{Val: 121},
			&ast.Char{Val: 122},
			&ast.Char{Val: 123},
			&ast.Char{Val: 124},
			&ast.Char{Val: 125},
			&ast.Char{Val: 126},
			&ast.Char{Val: 127},
		},
	}
)

func intPtr(v int) *int {
	return &v
}

func TestParse(t *testing.T) {
	tests := []struct {
		name             string
		in               input
		expectedError    string
		expectedAST      ast.Node
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
	}{
		{
			name:          "InvalidCharRange",
			in:            newStringInput("[9-0]"),
			expectedError: "1 error occurred:\n\t* invalid character range 9-0\n\n",
		},
		{
			name:          "InvalidRepRange",
			in:            newStringInput("[0-9]{4,2}"),
			expectedError: "1 error occurred:\n\t* invalid repetition range {4,2}\n\n",
		},
		{
			name: "Successful",
			in:   newStringInput(`[A-Z]?[a-z][0-9a-z]{1,}`),
			expectedAST: &ast.Concat{
				Exprs: []ast.Node{
					&ast.Alt{
						Exprs: []ast.Node{
							&ast.Empty{},
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Char{Val: 'A', Pos: 1},
									&ast.Char{Val: 'B', Pos: 2},
									&ast.Char{Val: 'C', Pos: 3},
									&ast.Char{Val: 'D', Pos: 4},
									&ast.Char{Val: 'E', Pos: 5},
									&ast.Char{Val: 'F', Pos: 6},
									&ast.Char{Val: 'G', Pos: 7},
									&ast.Char{Val: 'H', Pos: 8},
									&ast.Char{Val: 'I', Pos: 9},
									&ast.Char{Val: 'J', Pos: 10},
									&ast.Char{Val: 'K', Pos: 11},
									&ast.Char{Val: 'L', Pos: 12},
									&ast.Char{Val: 'M', Pos: 13},
									&ast.Char{Val: 'N', Pos: 14},
									&ast.Char{Val: 'O', Pos: 15},
									&ast.Char{Val: 'P', Pos: 16},
									&ast.Char{Val: 'Q', Pos: 17},
									&ast.Char{Val: 'R', Pos: 18},
									&ast.Char{Val: 'S', Pos: 19},
									&ast.Char{Val: 'T', Pos: 20},
									&ast.Char{Val: 'U', Pos: 21},
									&ast.Char{Val: 'V', Pos: 22},
									&ast.Char{Val: 'W', Pos: 23},
									&ast.Char{Val: 'X', Pos: 24},
									&ast.Char{Val: 'Y', Pos: 25},
									&ast.Char{Val: 'Z', Pos: 26},
								},
							},
						},
					},
					&ast.Alt{
						Exprs: []ast.Node{
							&ast.Char{Val: 'a', Pos: 27},
							&ast.Char{Val: 'b', Pos: 28},
							&ast.Char{Val: 'c', Pos: 29},
							&ast.Char{Val: 'd', Pos: 30},
							&ast.Char{Val: 'e', Pos: 31},
							&ast.Char{Val: 'f', Pos: 32},
							&ast.Char{Val: 'g', Pos: 33},
							&ast.Char{Val: 'h', Pos: 34},
							&ast.Char{Val: 'i', Pos: 35},
							&ast.Char{Val: 'j', Pos: 36},
							&ast.Char{Val: 'k', Pos: 37},
							&ast.Char{Val: 'l', Pos: 38},
							&ast.Char{Val: 'm', Pos: 39},
							&ast.Char{Val: 'n', Pos: 40},
							&ast.Char{Val: 'o', Pos: 41},
							&ast.Char{Val: 'p', Pos: 42},
							&ast.Char{Val: 'q', Pos: 43},
							&ast.Char{Val: 'r', Pos: 44},
							&ast.Char{Val: 's', Pos: 45},
							&ast.Char{Val: 't', Pos: 46},
							&ast.Char{Val: 'u', Pos: 47},
							&ast.Char{Val: 'v', Pos: 48},
							&ast.Char{Val: 'w', Pos: 49},
							&ast.Char{Val: 'x', Pos: 50},
							&ast.Char{Val: 'y', Pos: 51},
							&ast.Char{Val: 'z', Pos: 52},
						},
					},
					&ast.Concat{
						Exprs: []ast.Node{
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Char{Val: '0', Pos: 53},
									&ast.Char{Val: '1', Pos: 54},
									&ast.Char{Val: '2', Pos: 55},
									&ast.Char{Val: '3', Pos: 56},
									&ast.Char{Val: '4', Pos: 57},
									&ast.Char{Val: '5', Pos: 58},
									&ast.Char{Val: '6', Pos: 59},
									&ast.Char{Val: '7', Pos: 60},
									&ast.Char{Val: '8', Pos: 61},
									&ast.Char{Val: '9', Pos: 62},
									&ast.Char{Val: 'a', Pos: 63},
									&ast.Char{Val: 'b', Pos: 64},
									&ast.Char{Val: 'c', Pos: 65},
									&ast.Char{Val: 'd', Pos: 66},
									&ast.Char{Val: 'e', Pos: 67},
									&ast.Char{Val: 'f', Pos: 68},
									&ast.Char{Val: 'g', Pos: 69},
									&ast.Char{Val: 'h', Pos: 70},
									&ast.Char{Val: 'i', Pos: 71},
									&ast.Char{Val: 'j', Pos: 72},
									&ast.Char{Val: 'k', Pos: 73},
									&ast.Char{Val: 'l', Pos: 74},
									&ast.Char{Val: 'm', Pos: 75},
									&ast.Char{Val: 'n', Pos: 76},
									&ast.Char{Val: 'o', Pos: 77},
									&ast.Char{Val: 'p', Pos: 78},
									&ast.Char{Val: 'q', Pos: 79},
									&ast.Char{Val: 'r', Pos: 80},
									&ast.Char{Val: 's', Pos: 81},
									&ast.Char{Val: 't', Pos: 82},
									&ast.Char{Val: 'u', Pos: 83},
									&ast.Char{Val: 'v', Pos: 84},
									&ast.Char{Val: 'w', Pos: 85},
									&ast.Char{Val: 'x', Pos: 86},
									&ast.Char{Val: 'y', Pos: 87},
									&ast.Char{Val: 'z', Pos: 88},
								},
							},
							&ast.Star{
								Expr: &ast.Alt{
									Exprs: []ast.Node{
										&ast.Char{Val: '0', Pos: 89},
										&ast.Char{Val: '1', Pos: 90},
										&ast.Char{Val: '2', Pos: 91},
										&ast.Char{Val: '3', Pos: 92},
										&ast.Char{Val: '4', Pos: 93},
										&ast.Char{Val: '5', Pos: 94},
										&ast.Char{Val: '6', Pos: 95},
										&ast.Char{Val: '7', Pos: 96},
										&ast.Char{Val: '8', Pos: 97},
										&ast.Char{Val: '9', Pos: 98},
										&ast.Char{Val: 'a', Pos: 99},
										&ast.Char{Val: 'b', Pos: 100},
										&ast.Char{Val: 'c', Pos: 101},
										&ast.Char{Val: 'd', Pos: 102},
										&ast.Char{Val: 'e', Pos: 103},
										&ast.Char{Val: 'f', Pos: 104},
										&ast.Char{Val: 'g', Pos: 105},
										&ast.Char{Val: 'h', Pos: 106},
										&ast.Char{Val: 'i', Pos: 107},
										&ast.Char{Val: 'j', Pos: 108},
										&ast.Char{Val: 'k', Pos: 109},
										&ast.Char{Val: 'l', Pos: 110},
										&ast.Char{Val: 'm', Pos: 111},
										&ast.Char{Val: 'n', Pos: 112},
										&ast.Char{Val: 'o', Pos: 113},
										&ast.Char{Val: 'p', Pos: 114},
										&ast.Char{Val: 'q', Pos: 115},
										&ast.Char{Val: 'r', Pos: 116},
										&ast.Char{Val: 's', Pos: 117},
										&ast.Char{Val: 't', Pos: 118},
										&ast.Char{Val: 'u', Pos: 119},
										&ast.Char{Val: 'v', Pos: 120},
										&ast.Char{Val: 'w', Pos: 121},
										&ast.Char{Val: 'x', Pos: 122},
										&ast.Char{Val: 'y', Pos: 123},
										&ast.Char{Val: 'z', Pos: 124},
									},
								},
							},
						},
					},
				},
			},
			expectedNullable: false,
			expectedFirstPos: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
				27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
			},
			expectedLastPos: []int{
				53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
				89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ast, err := Parse(tc.in)

			if tc.expectedError != "" {
				assert.Nil(t, ast)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAST, ast)
				assert.Equal(t, tc.expectedNullable, ast.Nullable())
				assert.Equal(t, tc.expectedFirstPos, ast.FirstPos())
				assert.Equal(t, tc.expectedLastPos, ast.LastPos())
			}
		})
	}
}

func TestRegexConverters(t *testing.T) {
	r := newRegex()

	tests := []struct {
		name        string
		p           parser
		in          input
		expectedOK  bool
		expectedOut output
	}{
		{
			name:       "char_Successful",
			p:          r.char,
			in:         newStringInput(`!"#$%&'()*+,-./[\]^_{|}~`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Char{Val: '!'},
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune(`"#$%&'()*+,-./[\]^_{|}~`),
				},
			},
		},
		{
			name:       "digit_Successful",
			p:          r.digit,
			in:         newStringInput("0123456789"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '0',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("123456789"),
				},
			},
		},
		{
			name:       "letter_Successful",
			p:          r.letter,
			in:         newStringInput("abcdefghijklmnopqrstuvwxyz"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: 'a',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("bcdefghijklmnopqrstuvwxyz"),
				},
			},
		},
		{
			name:       "num_Successful",
			p:          r.num,
			in:         newStringInput("2022tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: 2022,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "letters_Successful",
			p:          r.letters,
			in:         newStringInput("head2022"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: "head",
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("2022"),
				},
			},
		},
		{
			name:       "repOp_ZeroOrOne_Successful",
			p:          r.repOp,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '?',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repOp_ZeroOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '*',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repOp_OneOrMore_Successful",
			p:          r.repOp,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '+',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "upperBound_Unbounded_Successful",
			p:          r.upperBound,
			in:         newStringInput(",}"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: (*int)(nil),
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("}"),
				},
			},
		},
		{
			name:       "upperBound_Bounded_Successful",
			p:          r.upperBound,
			in:         newStringInput(",4}"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: intPtr(4),
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("}"),
				},
			},
		},
		{
			name:       "range_Fixed_Successful",
			p:          r.range_,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: intPtr(2),
					},
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "range_upper_Unbounded_Successful",
			p:          r.range_,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: nil,
					},
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "range_upper_Bounded_Successful",
			p:          r.range_,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: intPtr(4),
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_ZeroOrOne_Successful",
			p:          r.repetition,
			in:         newStringInput("?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '?',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_ZeroOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '*',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_OneOrMore_Successful",
			p:          r.repetition,
			in:         newStringInput("+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '+',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_range_Fixed_Successful",
			p:          r.repetition,
			in:         newStringInput("{2}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: intPtr(2),
					},
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_range_upper_Unbounded_Successful",
			p:          r.repetition,
			in:         newStringInput("{2,}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: nil,
					},
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "repetition_range_upper_Bounded_Successful",
			p:          r.repetition,
			in:         newStringInput("{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[int, *int]{
						p: 2,
						q: intPtr(4),
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_ZeroOrOne_Successful",
			p:          r.quantifier,
			in:         newStringInput("??tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: '?',
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_ZeroOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("*?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: '*',
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_OneOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("+?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: '+',
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_Fixed_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: tuple[int, *int]{
							p: 2,
							q: intPtr(2),
						},
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_upper_Unbounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: tuple[int, *int]{
							p: 2,
							q: nil,
						},
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "quantifier_range_upper_Bounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,4}?tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: tuple[any, bool]{
						p: tuple[int, *int]{
							p: 2,
							q: intPtr(4),
						},
						q: true,
					},
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charRange_Range_Successful",
			p:          r.charRange,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Digit_Successful",
			p:          r.charClass,
			in:         newStringInput(`\dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotDigit_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Dtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: nonDigit,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Whitespace_Successful",
			p:          r.charClass,
			in:         newStringInput(`\stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: whitespace,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotWhitespace_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Stail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: nonWhitespace,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Word_Successful",
			p:          r.charClass,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_NotWord_Successful",
			p:          r.charClass,
			in:         newStringInput(`\Wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: nonWord,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Blank_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:blank:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Alt{
						Exprs: []ast.Node{
							&ast.Char{Val: ' '},
							&ast.Char{Val: '\t'},
						},
					},
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Space_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:space:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Alt{
						Exprs: []ast.Node{
							&ast.Char{Val: ' '},
							&ast.Char{Val: '\t'},
							&ast.Char{Val: '\n'},
							&ast.Char{Val: '\r'},
							&ast.Char{Val: '\f'},
							&ast.Char{Val: '\v'},
						},
					},
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Digit_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:digit:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_XDigit_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:xdigit:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: xdigit,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Upper_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:upper:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: upper,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Lower_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:lower:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: lower,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Alpha_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:alpha:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: alpha,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Alnum_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:alnum:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: alnum,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_Word_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "asciiCharClass_ASCII_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:ascii:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ascii,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_charClass_Successful",
			p:          r.charGroupItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_asciiCharClass_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_charRange_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("0-9tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_char_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Char{Val: '!'},
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_charClass_Successful",
			p:          r.charGroup,
			in:         newStringInput(`[\w]tail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_asciiCharClass_Successful",
			p:          r.charGroup,
			in:         newStringInput("[[:word:]]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_charRange_Successful",
			p:          r.charGroup,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_char_Successful",
			p:          r.charGroup,
			in:         newStringInput("[!]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Alt{
						Exprs: []ast.Node{
							&ast.Char{Val: '!'},
						},
					},
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_Negated_Successful",
			p:          r.charGroup,
			in:         newStringInput("[^#$]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Alt{
						Exprs: []ast.Node{
							&ast.Char{Val: 0},
							&ast.Char{Val: 1},
							&ast.Char{Val: 2},
							&ast.Char{Val: 3},
							&ast.Char{Val: 4},
							&ast.Char{Val: 5},
							&ast.Char{Val: 6},
							&ast.Char{Val: 7},
							&ast.Char{Val: 8},
							&ast.Char{Val: 9},
							&ast.Char{Val: 10},
							&ast.Char{Val: 11},
							&ast.Char{Val: 12},
							&ast.Char{Val: 13},
							&ast.Char{Val: 14},
							&ast.Char{Val: 15},
							&ast.Char{Val: 16},
							&ast.Char{Val: 17},
							&ast.Char{Val: 18},
							&ast.Char{Val: 19},
							&ast.Char{Val: 20},
							&ast.Char{Val: 21},
							&ast.Char{Val: 22},
							&ast.Char{Val: 23},
							&ast.Char{Val: 24},
							&ast.Char{Val: 25},
							&ast.Char{Val: 26},
							&ast.Char{Val: 27},
							&ast.Char{Val: 28},
							&ast.Char{Val: 29},
							&ast.Char{Val: 30},
							&ast.Char{Val: 31},
							&ast.Char{Val: 32},
							&ast.Char{Val: 33},
							&ast.Char{Val: 34},
							&ast.Char{Val: 37},
							&ast.Char{Val: 38},
							&ast.Char{Val: 39},
							&ast.Char{Val: 40},
							&ast.Char{Val: 41},
							&ast.Char{Val: 42},
							&ast.Char{Val: 43},
							&ast.Char{Val: 44},
							&ast.Char{Val: 45},
							&ast.Char{Val: 46},
							&ast.Char{Val: 47},
							&ast.Char{Val: 48},
							&ast.Char{Val: 49},
							&ast.Char{Val: 50},
							&ast.Char{Val: 51},
							&ast.Char{Val: 52},
							&ast.Char{Val: 53},
							&ast.Char{Val: 54},
							&ast.Char{Val: 55},
							&ast.Char{Val: 56},
							&ast.Char{Val: 57},
							&ast.Char{Val: 58},
							&ast.Char{Val: 59},
							&ast.Char{Val: 60},
							&ast.Char{Val: 61},
							&ast.Char{Val: 62},
							&ast.Char{Val: 63},
							&ast.Char{Val: 64},
							&ast.Char{Val: 65},
							&ast.Char{Val: 66},
							&ast.Char{Val: 67},
							&ast.Char{Val: 68},
							&ast.Char{Val: 69},
							&ast.Char{Val: 70},
							&ast.Char{Val: 71},
							&ast.Char{Val: 72},
							&ast.Char{Val: 73},
							&ast.Char{Val: 74},
							&ast.Char{Val: 75},
							&ast.Char{Val: 76},
							&ast.Char{Val: 77},
							&ast.Char{Val: 78},
							&ast.Char{Val: 79},
							&ast.Char{Val: 80},
							&ast.Char{Val: 81},
							&ast.Char{Val: 82},
							&ast.Char{Val: 83},
							&ast.Char{Val: 84},
							&ast.Char{Val: 85},
							&ast.Char{Val: 86},
							&ast.Char{Val: 87},
							&ast.Char{Val: 88},
							&ast.Char{Val: 89},
							&ast.Char{Val: 90},
							&ast.Char{Val: 91},
							&ast.Char{Val: 92},
							&ast.Char{Val: 93},
							&ast.Char{Val: 94},
							&ast.Char{Val: 95},
							&ast.Char{Val: 96},
							&ast.Char{Val: 97},
							&ast.Char{Val: 98},
							&ast.Char{Val: 99},
							&ast.Char{Val: 100},
							&ast.Char{Val: 101},
							&ast.Char{Val: 102},
							&ast.Char{Val: 103},
							&ast.Char{Val: 104},
							&ast.Char{Val: 105},
							&ast.Char{Val: 106},
							&ast.Char{Val: 107},
							&ast.Char{Val: 108},
							&ast.Char{Val: 109},
							&ast.Char{Val: 110},
							&ast.Char{Val: 111},
							&ast.Char{Val: 112},
							&ast.Char{Val: 113},
							&ast.Char{Val: 114},
							&ast.Char{Val: 115},
							&ast.Char{Val: 116},
							&ast.Char{Val: 117},
							&ast.Char{Val: 118},
							&ast.Char{Val: 119},
							&ast.Char{Val: 120},
							&ast.Char{Val: 121},
							&ast.Char{Val: 122},
							&ast.Char{Val: 123},
							&ast.Char{Val: 124},
							&ast.Char{Val: 125},
							&ast.Char{Val: 126},
							&ast.Char{Val: 127},
						},
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anyChar_Successful",
			p:          r.anyChar,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ascii,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_anyChar_Successful",
			p:          r.matchItem,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: ascii,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charClass_Successful",
			p:          r.matchItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_asciiCharClass_Successful",
			p:          r.matchItem,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charGroup_charRange_Successful",
			p:          r.matchItem,
			in:         newStringInput("[0-9]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_char_Successful",
			p:          r.matchItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Char{Val: '!'},
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charClass_Successful",
			p:          r.match,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_asciiCharClass_Successful",
			p:          r.match,
			in:         newStringInput("[:word:]tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: word,
				},
				Remaining: &stringInput{
					pos:   8,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charGroup_charRange_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("[0-9]{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							digit,
							digit,
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Empty{},
									digit,
								},
							},
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Empty{},
									digit,
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   10,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_char_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("#{2,4}tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							&ast.Char{Val: '#'},
							&ast.Char{Val: '#'},
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Empty{},
									&ast.Char{Val: '#'},
								},
							},
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Empty{},
									&ast.Char{Val: '#'},
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "group_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Alt{
						Exprs: []ast.Node{
							&ast.Concat{
								Exprs: []ast.Node{
									&ast.Char{Val: 'a'},
								},
							},
							&ast.Concat{
								Exprs: []ast.Node{
									&ast.Char{Val: 'b'},
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "group_quantifier_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							&ast.Alt{
								Exprs: []ast.Node{
									&ast.Concat{
										Exprs: []ast.Node{
											&ast.Char{Val: 'a'},
										},
									},
									&ast.Concat{
										Exprs: []ast.Node{
											&ast.Char{Val: 'b'},
										},
									},
								},
							},
							&ast.Star{
								Expr: &ast.Alt{
									Exprs: []ast.Node{
										&ast.Concat{
											Exprs: []ast.Node{
												&ast.Char{Val: 'a'},
											},
										},
										&ast.Concat{
											Exprs: []ast.Node{
												&ast.Char{Val: 'b'},
											},
										},
									},
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anchor_Successful",
			p:          r.anchor,
			in:         newStringInput("$tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '$',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_group_Successful",
			p:          r.subexprItem,
			in:         newStringInput("(ab)+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							&ast.Concat{
								Exprs: []ast.Node{
									&ast.Char{Val: 'a'},
									&ast.Char{Val: 'b'},
								},
							},
							&ast.Star{
								Expr: &ast.Concat{
									Exprs: []ast.Node{
										&ast.Char{Val: 'a'},
										&ast.Char{Val: 'b'},
									},
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexprItem_anchor_Successful",
			p:          r.subexprItem,
			in:         newStringInput("$"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: '$',
				},
				Remaining: nil,
			},
		},
		{
			name:       "subexprItem_match_charGroup_Successful",
			p:          r.subexprItem,
			in:         newStringInput("[0-9]+tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							digit,
							&ast.Star{
								Expr: digit,
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   6,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "subexpr_group_matches_Successful",
			p:          r.subexpr,
			in:         newStringInput("(ab)+[0-9]*tail"),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							&ast.Concat{
								Exprs: []ast.Node{
									&ast.Concat{
										Exprs: []ast.Node{
											&ast.Char{Val: 'a'},
											&ast.Char{Val: 'b'},
										},
									},
									&ast.Star{
										Expr: &ast.Concat{
											Exprs: []ast.Node{
												&ast.Char{Val: 'a'},
												&ast.Char{Val: 'b'},
											},
										},
									},
								},
							},
							&ast.Star{
								Expr: digit,
							},
							&ast.Char{Val: 't'},
							&ast.Char{Val: 'a'},
							&ast.Char{Val: 'i'},
							&ast.Char{Val: 'l'},
						},
					},
				},
				Remaining: nil,
			},
		},
		{
			name:       "expr_Successful",
			p:          r.expr,
			in:         newStringInput(`[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							word,
							&ast.Star{
								Expr: word,
							},
						},
					},
				},
				Remaining: nil,
			},
		},
		{
			name:       "regex_Successful",
			p:          r.regex,
			in:         newStringInput(`^package\s+[0-9A-Za-z_][\d\w]*`),
			expectedOK: true,
			expectedOut: output{
				Result: result{
					Val: &ast.Concat{
						Exprs: []ast.Node{
							&ast.Char{Val: 'p'},
							&ast.Char{Val: 'a'},
							&ast.Char{Val: 'c'},
							&ast.Char{Val: 'k'},
							&ast.Char{Val: 'a'},
							&ast.Char{Val: 'g'},
							&ast.Char{Val: 'e'},
							&ast.Concat{
								Exprs: []ast.Node{
									whitespace,
									&ast.Star{
										Expr: whitespace,
									},
								},
							},
							word,
							&ast.Star{
								Expr: word,
							},
						},
					},
				},
				Remaining: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the state
			r.errors = nil

			out, ok := tc.p(tc.in)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedOut, out)
		})
	}
}
