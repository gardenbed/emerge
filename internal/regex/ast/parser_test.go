package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	comb "github.com/gardenbed/emerge/internal/combinator"
)

var (
	digit = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
		},
	}

	nonDigit = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	whitespace = &Alt{
		Exprs: []Node{
			&Char{Val: ' '},
			&Char{Val: '\t'},
			&Char{Val: '\n'},
			&Char{Val: '\r'},
			&Char{Val: '\f'},
		},
	}

	nonWhitespace = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 11},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 48},
			&Char{Val: 49},
			&Char{Val: 50},
			&Char{Val: 51},
			&Char{Val: 52},
			&Char{Val: 53},
			&Char{Val: 54},
			&Char{Val: 55},
			&Char{Val: 56},
			&Char{Val: 57},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	word = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: '_'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	nonWord = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 96},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}

	xdigit = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
		},
	}

	upper = &Alt{
		Exprs: []Node{
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
		},
	}

	lower = &Alt{
		Exprs: []Node{
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	alpha = &Alt{
		Exprs: []Node{
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	alnum = &Alt{
		Exprs: []Node{
			&Char{Val: '0'},
			&Char{Val: '1'},
			&Char{Val: '2'},
			&Char{Val: '3'},
			&Char{Val: '4'},
			&Char{Val: '5'},
			&Char{Val: '6'},
			&Char{Val: '7'},
			&Char{Val: '8'},
			&Char{Val: '9'},
			&Char{Val: 'A'},
			&Char{Val: 'B'},
			&Char{Val: 'C'},
			&Char{Val: 'D'},
			&Char{Val: 'E'},
			&Char{Val: 'F'},
			&Char{Val: 'G'},
			&Char{Val: 'H'},
			&Char{Val: 'I'},
			&Char{Val: 'J'},
			&Char{Val: 'K'},
			&Char{Val: 'L'},
			&Char{Val: 'M'},
			&Char{Val: 'N'},
			&Char{Val: 'O'},
			&Char{Val: 'P'},
			&Char{Val: 'Q'},
			&Char{Val: 'R'},
			&Char{Val: 'S'},
			&Char{Val: 'T'},
			&Char{Val: 'U'},
			&Char{Val: 'V'},
			&Char{Val: 'W'},
			&Char{Val: 'X'},
			&Char{Val: 'Y'},
			&Char{Val: 'Z'},
			&Char{Val: 'a'},
			&Char{Val: 'b'},
			&Char{Val: 'c'},
			&Char{Val: 'd'},
			&Char{Val: 'e'},
			&Char{Val: 'f'},
			&Char{Val: 'g'},
			&Char{Val: 'h'},
			&Char{Val: 'i'},
			&Char{Val: 'j'},
			&Char{Val: 'k'},
			&Char{Val: 'l'},
			&Char{Val: 'm'},
			&Char{Val: 'n'},
			&Char{Val: 'o'},
			&Char{Val: 'p'},
			&Char{Val: 'q'},
			&Char{Val: 'r'},
			&Char{Val: 's'},
			&Char{Val: 't'},
			&Char{Val: 'u'},
			&Char{Val: 'v'},
			&Char{Val: 'w'},
			&Char{Val: 'x'},
			&Char{Val: 'y'},
			&Char{Val: 'z'},
		},
	}

	ascii = &Alt{
		Exprs: []Node{
			&Char{Val: 0},
			&Char{Val: 1},
			&Char{Val: 2},
			&Char{Val: 3},
			&Char{Val: 4},
			&Char{Val: 5},
			&Char{Val: 6},
			&Char{Val: 7},
			&Char{Val: 8},
			&Char{Val: 9},
			&Char{Val: 10},
			&Char{Val: 11},
			&Char{Val: 12},
			&Char{Val: 13},
			&Char{Val: 14},
			&Char{Val: 15},
			&Char{Val: 16},
			&Char{Val: 17},
			&Char{Val: 18},
			&Char{Val: 19},
			&Char{Val: 20},
			&Char{Val: 21},
			&Char{Val: 22},
			&Char{Val: 23},
			&Char{Val: 24},
			&Char{Val: 25},
			&Char{Val: 26},
			&Char{Val: 27},
			&Char{Val: 28},
			&Char{Val: 29},
			&Char{Val: 30},
			&Char{Val: 31},
			&Char{Val: 32},
			&Char{Val: 33},
			&Char{Val: 34},
			&Char{Val: 35},
			&Char{Val: 36},
			&Char{Val: 37},
			&Char{Val: 38},
			&Char{Val: 39},
			&Char{Val: 40},
			&Char{Val: 41},
			&Char{Val: 42},
			&Char{Val: 43},
			&Char{Val: 44},
			&Char{Val: 45},
			&Char{Val: 46},
			&Char{Val: 47},
			&Char{Val: 48},
			&Char{Val: 49},
			&Char{Val: 50},
			&Char{Val: 51},
			&Char{Val: 52},
			&Char{Val: 53},
			&Char{Val: 54},
			&Char{Val: 55},
			&Char{Val: 56},
			&Char{Val: 57},
			&Char{Val: 58},
			&Char{Val: 59},
			&Char{Val: 60},
			&Char{Val: 61},
			&Char{Val: 62},
			&Char{Val: 63},
			&Char{Val: 64},
			&Char{Val: 65},
			&Char{Val: 66},
			&Char{Val: 67},
			&Char{Val: 68},
			&Char{Val: 69},
			&Char{Val: 70},
			&Char{Val: 71},
			&Char{Val: 72},
			&Char{Val: 73},
			&Char{Val: 74},
			&Char{Val: 75},
			&Char{Val: 76},
			&Char{Val: 77},
			&Char{Val: 78},
			&Char{Val: 79},
			&Char{Val: 80},
			&Char{Val: 81},
			&Char{Val: 82},
			&Char{Val: 83},
			&Char{Val: 84},
			&Char{Val: 85},
			&Char{Val: 86},
			&Char{Val: 87},
			&Char{Val: 88},
			&Char{Val: 89},
			&Char{Val: 90},
			&Char{Val: 91},
			&Char{Val: 92},
			&Char{Val: 93},
			&Char{Val: 94},
			&Char{Val: 95},
			&Char{Val: 96},
			&Char{Val: 97},
			&Char{Val: 98},
			&Char{Val: 99},
			&Char{Val: 100},
			&Char{Val: 101},
			&Char{Val: 102},
			&Char{Val: 103},
			&Char{Val: 104},
			&Char{Val: 105},
			&Char{Val: 106},
			&Char{Val: 107},
			&Char{Val: 108},
			&Char{Val: 109},
			&Char{Val: 110},
			&Char{Val: 111},
			&Char{Val: 112},
			&Char{Val: 113},
			&Char{Val: 114},
			&Char{Val: 115},
			&Char{Val: 116},
			&Char{Val: 117},
			&Char{Val: 118},
			&Char{Val: 119},
			&Char{Val: 120},
			&Char{Val: 121},
			&Char{Val: 122},
			&Char{Val: 123},
			&Char{Val: 124},
			&Char{Val: 125},
			&Char{Val: 126},
			&Char{Val: 127},
		},
	}
)

func intPtr(v int) *int {
	return &v
}

// stringInput implements the input interface for strings.
type stringInput struct {
	pos   int
	runes []rune
}

func newStringInput(s string) comb.Input {
	return &stringInput{
		pos:   0,
		runes: []rune(s),
	}
}

func (s *stringInput) Current() (rune, int) {
	return s.runes[0], s.pos
}

func (s *stringInput) Remaining() comb.Input {
	if len(s.runes) == 1 {
		return nil
	}

	return &stringInput{
		pos:   s.pos + 1,
		runes: s.runes[1:],
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name             string
		in               comb.Input
		expectedError    string
		expectedAST      Node
		expectedNullable bool
		expectedFirstPos []int
		expectedLastPos  []int
	}{
		{
			name:          "InvalidRegex",
			in:            newStringInput("["),
			expectedError: "invalid regular expression",
		},
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
			expectedAST: &Concat{
				Exprs: []Node{
					&Alt{
						Exprs: []Node{
							&Empty{},
							&Alt{
								Exprs: []Node{
									&Char{Val: 'A', Pos: 1},
									&Char{Val: 'B', Pos: 2},
									&Char{Val: 'C', Pos: 3},
									&Char{Val: 'D', Pos: 4},
									&Char{Val: 'E', Pos: 5},
									&Char{Val: 'F', Pos: 6},
									&Char{Val: 'G', Pos: 7},
									&Char{Val: 'H', Pos: 8},
									&Char{Val: 'I', Pos: 9},
									&Char{Val: 'J', Pos: 10},
									&Char{Val: 'K', Pos: 11},
									&Char{Val: 'L', Pos: 12},
									&Char{Val: 'M', Pos: 13},
									&Char{Val: 'N', Pos: 14},
									&Char{Val: 'O', Pos: 15},
									&Char{Val: 'P', Pos: 16},
									&Char{Val: 'Q', Pos: 17},
									&Char{Val: 'R', Pos: 18},
									&Char{Val: 'S', Pos: 19},
									&Char{Val: 'T', Pos: 20},
									&Char{Val: 'U', Pos: 21},
									&Char{Val: 'V', Pos: 22},
									&Char{Val: 'W', Pos: 23},
									&Char{Val: 'X', Pos: 24},
									&Char{Val: 'Y', Pos: 25},
									&Char{Val: 'Z', Pos: 26},
								},
							},
						},
					},
					&Alt{
						Exprs: []Node{
							&Char{Val: 'a', Pos: 27},
							&Char{Val: 'b', Pos: 28},
							&Char{Val: 'c', Pos: 29},
							&Char{Val: 'd', Pos: 30},
							&Char{Val: 'e', Pos: 31},
							&Char{Val: 'f', Pos: 32},
							&Char{Val: 'g', Pos: 33},
							&Char{Val: 'h', Pos: 34},
							&Char{Val: 'i', Pos: 35},
							&Char{Val: 'j', Pos: 36},
							&Char{Val: 'k', Pos: 37},
							&Char{Val: 'l', Pos: 38},
							&Char{Val: 'm', Pos: 39},
							&Char{Val: 'n', Pos: 40},
							&Char{Val: 'o', Pos: 41},
							&Char{Val: 'p', Pos: 42},
							&Char{Val: 'q', Pos: 43},
							&Char{Val: 'r', Pos: 44},
							&Char{Val: 's', Pos: 45},
							&Char{Val: 't', Pos: 46},
							&Char{Val: 'u', Pos: 47},
							&Char{Val: 'v', Pos: 48},
							&Char{Val: 'w', Pos: 49},
							&Char{Val: 'x', Pos: 50},
							&Char{Val: 'y', Pos: 51},
							&Char{Val: 'z', Pos: 52},
						},
					},
					&Concat{
						Exprs: []Node{
							&Alt{
								Exprs: []Node{
									&Char{Val: '0', Pos: 53},
									&Char{Val: '1', Pos: 54},
									&Char{Val: '2', Pos: 55},
									&Char{Val: '3', Pos: 56},
									&Char{Val: '4', Pos: 57},
									&Char{Val: '5', Pos: 58},
									&Char{Val: '6', Pos: 59},
									&Char{Val: '7', Pos: 60},
									&Char{Val: '8', Pos: 61},
									&Char{Val: '9', Pos: 62},
									&Char{Val: 'a', Pos: 63},
									&Char{Val: 'b', Pos: 64},
									&Char{Val: 'c', Pos: 65},
									&Char{Val: 'd', Pos: 66},
									&Char{Val: 'e', Pos: 67},
									&Char{Val: 'f', Pos: 68},
									&Char{Val: 'g', Pos: 69},
									&Char{Val: 'h', Pos: 70},
									&Char{Val: 'i', Pos: 71},
									&Char{Val: 'j', Pos: 72},
									&Char{Val: 'k', Pos: 73},
									&Char{Val: 'l', Pos: 74},
									&Char{Val: 'm', Pos: 75},
									&Char{Val: 'n', Pos: 76},
									&Char{Val: 'o', Pos: 77},
									&Char{Val: 'p', Pos: 78},
									&Char{Val: 'q', Pos: 79},
									&Char{Val: 'r', Pos: 80},
									&Char{Val: 's', Pos: 81},
									&Char{Val: 't', Pos: 82},
									&Char{Val: 'u', Pos: 83},
									&Char{Val: 'v', Pos: 84},
									&Char{Val: 'w', Pos: 85},
									&Char{Val: 'x', Pos: 86},
									&Char{Val: 'y', Pos: 87},
									&Char{Val: 'z', Pos: 88},
								},
							},
							&Star{
								Expr: &Alt{
									Exprs: []Node{
										&Char{Val: '0', Pos: 89},
										&Char{Val: '1', Pos: 90},
										&Char{Val: '2', Pos: 91},
										&Char{Val: '3', Pos: 92},
										&Char{Val: '4', Pos: 93},
										&Char{Val: '5', Pos: 94},
										&Char{Val: '6', Pos: 95},
										&Char{Val: '7', Pos: 96},
										&Char{Val: '8', Pos: 97},
										&Char{Val: '9', Pos: 98},
										&Char{Val: 'a', Pos: 99},
										&Char{Val: 'b', Pos: 100},
										&Char{Val: 'c', Pos: 101},
										&Char{Val: 'd', Pos: 102},
										&Char{Val: 'e', Pos: 103},
										&Char{Val: 'f', Pos: 104},
										&Char{Val: 'g', Pos: 105},
										&Char{Val: 'h', Pos: 106},
										&Char{Val: 'i', Pos: 107},
										&Char{Val: 'j', Pos: 108},
										&Char{Val: 'k', Pos: 109},
										&Char{Val: 'l', Pos: 110},
										&Char{Val: 'm', Pos: 111},
										&Char{Val: 'n', Pos: 112},
										&Char{Val: 'o', Pos: 113},
										&Char{Val: 'p', Pos: 114},
										&Char{Val: 'q', Pos: 115},
										&Char{Val: 'r', Pos: 116},
										&Char{Val: 's', Pos: 117},
										&Char{Val: 't', Pos: 118},
										&Char{Val: 'u', Pos: 119},
										&Char{Val: 'v', Pos: 120},
										&Char{Val: 'w', Pos: 121},
										&Char{Val: 'x', Pos: 122},
										&Char{Val: 'y', Pos: 123},
										&Char{Val: 'z', Pos: 124},
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

func TestRegexMappers(t *testing.T) {
	r := newRegex()

	tests := []struct {
		name        string
		p           comb.Parser
		in          comb.Input
		expectedOK  bool
		expectedOut comb.Output
	}{
		{
			name:       "char_Successful",
			p:          r.char,
			in:         newStringInput(`!"#$%&'()*+,-./[\]^_{|}~`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: '!',
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune(`"#$%&'()*+,-./[\]^_{|}~`),
				},
			},
		},
		{
			name:       "unescapedChar_Successful",
			p:          r.unescapedChar,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '!'},
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "escapedChar_Successful",
			p:          r.escapedChar,
			in:         newStringInput(`\+tail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '+'},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "digit_Successful",
			p:          r.digit,
			in:         newStringInput("0123456789"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_ZeroOrOne_Successful",
			p:          r.quantifier,
			in:         newStringInput("??tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_ZeroOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("*?tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_OneOrMore_Successful",
			p:          r.quantifier,
			in:         newStringInput("+?tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_range_Fixed_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2}?tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_range_upper_Unbounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,}?tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			name:       "quantifier_Lazy_range_upper_Bounded_Successful",
			p:          r.quantifier,
			in:         newStringInput("{2,4}?tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_charClass_Successful",
			p:          r.charGroupItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   3,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_escapedChar_Successful",
			p:          r.charGroupItem,
			in:         newStringInput(`\+tail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '+'},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroupItem_unescapedChar_Successful",
			p:          r.charGroupItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '!'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_escapedChar_Successful",
			p:          r.charGroup,
			in:         newStringInput(`[\+]tail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Char{Val: '+'},
						},
					},
				},
				Remaining: &stringInput{
					pos:   4,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charGroup_unescapedChar_Successful",
			p:          r.charGroup,
			in:         newStringInput("[!]tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Char{Val: '!'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Char{Val: 0},
							&Char{Val: 1},
							&Char{Val: 2},
							&Char{Val: 3},
							&Char{Val: 4},
							&Char{Val: 5},
							&Char{Val: 6},
							&Char{Val: 7},
							&Char{Val: 8},
							&Char{Val: 9},
							&Char{Val: 10},
							&Char{Val: 11},
							&Char{Val: 12},
							&Char{Val: 13},
							&Char{Val: 14},
							&Char{Val: 15},
							&Char{Val: 16},
							&Char{Val: 17},
							&Char{Val: 18},
							&Char{Val: 19},
							&Char{Val: 20},
							&Char{Val: 21},
							&Char{Val: 22},
							&Char{Val: 23},
							&Char{Val: 24},
							&Char{Val: 25},
							&Char{Val: 26},
							&Char{Val: 27},
							&Char{Val: 28},
							&Char{Val: 29},
							&Char{Val: 30},
							&Char{Val: 31},
							&Char{Val: 32},
							&Char{Val: 33},
							&Char{Val: 34},
							&Char{Val: 37},
							&Char{Val: 38},
							&Char{Val: 39},
							&Char{Val: 40},
							&Char{Val: 41},
							&Char{Val: 42},
							&Char{Val: 43},
							&Char{Val: 44},
							&Char{Val: 45},
							&Char{Val: 46},
							&Char{Val: 47},
							&Char{Val: 48},
							&Char{Val: 49},
							&Char{Val: 50},
							&Char{Val: 51},
							&Char{Val: 52},
							&Char{Val: 53},
							&Char{Val: 54},
							&Char{Val: 55},
							&Char{Val: 56},
							&Char{Val: 57},
							&Char{Val: 58},
							&Char{Val: 59},
							&Char{Val: 60},
							&Char{Val: 61},
							&Char{Val: 62},
							&Char{Val: 63},
							&Char{Val: 64},
							&Char{Val: 65},
							&Char{Val: 66},
							&Char{Val: 67},
							&Char{Val: 68},
							&Char{Val: 69},
							&Char{Val: 70},
							&Char{Val: 71},
							&Char{Val: 72},
							&Char{Val: 73},
							&Char{Val: 74},
							&Char{Val: 75},
							&Char{Val: 76},
							&Char{Val: 77},
							&Char{Val: 78},
							&Char{Val: 79},
							&Char{Val: 80},
							&Char{Val: 81},
							&Char{Val: 82},
							&Char{Val: 83},
							&Char{Val: 84},
							&Char{Val: 85},
							&Char{Val: 86},
							&Char{Val: 87},
							&Char{Val: 88},
							&Char{Val: 89},
							&Char{Val: 90},
							&Char{Val: 91},
							&Char{Val: 92},
							&Char{Val: 93},
							&Char{Val: 94},
							&Char{Val: 95},
							&Char{Val: 96},
							&Char{Val: 97},
							&Char{Val: 98},
							&Char{Val: 99},
							&Char{Val: 100},
							&Char{Val: 101},
							&Char{Val: 102},
							&Char{Val: 103},
							&Char{Val: 104},
							&Char{Val: 105},
							&Char{Val: 106},
							&Char{Val: 107},
							&Char{Val: 108},
							&Char{Val: 109},
							&Char{Val: 110},
							&Char{Val: 111},
							&Char{Val: 112},
							&Char{Val: 113},
							&Char{Val: 114},
							&Char{Val: 115},
							&Char{Val: 116},
							&Char{Val: 117},
							&Char{Val: 118},
							&Char{Val: 119},
							&Char{Val: 120},
							&Char{Val: 121},
							&Char{Val: 122},
							&Char{Val: 123},
							&Char{Val: 124},
							&Char{Val: 125},
							&Char{Val: 126},
							&Char{Val: 127},
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
			name:       "asciiCharClass_Blank_Successful",
			p:          r.asciiCharClass,
			in:         newStringInput("[:blank:]tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Char{Val: ' '},
							&Char{Val: '\t'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Char{Val: ' '},
							&Char{Val: '\t'},
							&Char{Val: '\n'},
							&Char{Val: '\r'},
							&Char{Val: '\f'},
							&Char{Val: '\v'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: ascii,
				},
				Remaining: &stringInput{
					pos:   9,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "charClass_Digit_Successful",
			p:          r.charClass,
			in:         newStringInput(`\dtail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: nonWord,
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "anyChar_Successful",
			p:          r.anyChar,
			in:         newStringInput(".tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: ascii,
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_unescapedChar_Successful",
			p:          r.matchItem,
			in:         newStringInput("!tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '!'},
				},
				Remaining: &stringInput{
					pos:   1,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_escapedChar_Successful",
			p:          r.matchItem,
			in:         newStringInput(`\+tail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Char{Val: '+'},
				},
				Remaining: &stringInput{
					pos:   2,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "matchItem_charClass_Successful",
			p:          r.matchItem,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: digit,
				},
				Remaining: &stringInput{
					pos:   5,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_unescapedChar_quantifier_Successful",
			p:          r.match,
			in:         newStringInput("#{2,4}tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Char{Val: '#'},
							&Char{Val: '#'},
							&Alt{
								Exprs: []Node{
									&Empty{},
									&Char{Val: '#'},
								},
							},
							&Alt{
								Exprs: []Node{
									&Empty{},
									&Char{Val: '#'},
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
			name:       "match_escapedChar_quantifier_Successful",
			p:          r.match,
			in:         newStringInput(`\+{2,4}tail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Char{Val: '+'},
							&Char{Val: '+'},
							&Alt{
								Exprs: []Node{
									&Empty{},
									&Char{Val: '+'},
								},
							},
							&Alt{
								Exprs: []Node{
									&Empty{},
									&Char{Val: '+'},
								},
							},
						},
					},
				},
				Remaining: &stringInput{
					pos:   7,
					runes: []rune("tail"),
				},
			},
		},
		{
			name:       "match_charClass_Successful",
			p:          r.match,
			in:         newStringInput(`\wtail`),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							digit,
							digit,
							&Alt{
								Exprs: []Node{
									&Empty{},
									digit,
								},
							},
							&Alt{
								Exprs: []Node{
									&Empty{},
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
			name:       "group_Successful",
			p:          r.group,
			in:         newStringInput("(a|b)tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Alt{
						Exprs: []Node{
							&Concat{
								Exprs: []Node{
									&Char{Val: 'a'},
								},
							},
							&Concat{
								Exprs: []Node{
									&Char{Val: 'b'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Alt{
								Exprs: []Node{
									&Concat{
										Exprs: []Node{
											&Char{Val: 'a'},
										},
									},
									&Concat{
										Exprs: []Node{
											&Char{Val: 'b'},
										},
									},
								},
							},
							&Star{
								Expr: &Alt{
									Exprs: []Node{
										&Concat{
											Exprs: []Node{
												&Char{Val: 'a'},
											},
										},
										&Concat{
											Exprs: []Node{
												&Char{Val: 'b'},
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
			in:         newStringInput("$"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: '$',
				},
				Remaining: nil,
			},
		},
		{
			name:       "subexprItem_anchor_Successful",
			p:          r.subexprItem,
			in:         newStringInput("$"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: '$',
				},
				Remaining: nil,
			},
		},
		{
			name:       "subexprItem_group_Successful",
			p:          r.subexprItem,
			in:         newStringInput("(ab)+tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Concat{
								Exprs: []Node{
									&Char{Val: 'a'},
									&Char{Val: 'b'},
								},
							},
							&Star{
								Expr: &Concat{
									Exprs: []Node{
										&Char{Val: 'a'},
										&Char{Val: 'b'},
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
			name:       "subexprItem_match_charGroup_Successful",
			p:          r.subexprItem,
			in:         newStringInput("[0-9]+tail"),
			expectedOK: true,
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							digit,
							&Star{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Concat{
								Exprs: []Node{
									&Concat{
										Exprs: []Node{
											&Char{Val: 'a'},
											&Char{Val: 'b'},
										},
									},
									&Star{
										Expr: &Concat{
											Exprs: []Node{
												&Char{Val: 'a'},
												&Char{Val: 'b'},
											},
										},
									},
								},
							},
							&Star{
								Expr: digit,
							},
							&Char{Val: 't'},
							&Char{Val: 'a'},
							&Char{Val: 'i'},
							&Char{Val: 'l'},
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							word,
							&Star{
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
			expectedOut: comb.Output{
				Result: comb.Result{
					Val: &Concat{
						Exprs: []Node{
							&Char{Val: 'p'},
							&Char{Val: 'a'},
							&Char{Val: 'c'},
							&Char{Val: 'k'},
							&Char{Val: 'a'},
							&Char{Val: 'g'},
							&Char{Val: 'e'},
							&Concat{
								Exprs: []Node{
									whitespace,
									&Star{
										Expr: whitespace,
									},
								},
							},
							word,
							&Star{
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
