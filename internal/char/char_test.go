package char

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeList_Dedup(t *testing.T) {
	tests := []struct {
		name           string
		l              RangeList
		x              RangeList
		expectedResult RangeList
	}{
		{
			name: "OK",
			l: RangeList{
				{'_', '_'},
				{'a', 'z'},
				{'0', '9'},
				{'A', 'Z'},
				{'0', '9'},
			},
			expectedResult: RangeList{
				{'0', '9'},
				{'A', 'Z'},
				{'_', '_'},
				{'a', 'z'},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.l.Dedup()

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestRangeList_Exclude(t *testing.T) {
	tests := []struct {
		name           string
		l              RangeList
		x              RangeList
		expectedResult RangeList
	}{
		{
			name: "OK",
			l:    Classes["ASCII"],
			x: RangeList{
				{'0', '9'},
				{'A', 'Z'},
				{'_', '_'},
				{'a', 'z'},
			},
			expectedResult: RangeList{
				{0x00, 0x2F},
				{0x3A, 0x40},
				{0x5B, 0x5E},
				{0x60, 0x60},
				{0x7B, 0x7F},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.l.Exclude(tc.x)

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
