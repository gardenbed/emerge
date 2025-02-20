package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneList(t *testing.T) {
	tests := []struct {
		name          string
		r             runeList
		expectedRunes []rune
	}{
		{
			name:          "OK",
			r:             runeList{' ', '\t', '\n', '\r'},
			expectedRunes: []rune{' ', '\t', '\n', '\r'},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRunes, tc.r.Runes())
		})
	}
}

func TestRuneRange(t *testing.T) {
	tests := []struct {
		name          string
		r             runeRange
		expectedRunes []rune
	}{
		{
			name:          "OK",
			r:             runeRange{'0', '9'},
			expectedRunes: []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRunes, tc.r.Runes())
		})
	}
}

func TestRuneClass(t *testing.T) {
	tests := []struct {
		name          string
		r             RuneClass
		expectedRunes []rune
	}{
		{
			name:          "OK",
			r:             RuneClass{runeList{'-', '.'}, runeRange{'0', '9'}},
			expectedRunes: []rune{'-', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRunes, tc.r.Runes())
		})
	}
}
