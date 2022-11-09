package input

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f, err := os.Open("./fixture/lorem_ipsum")
	assert.NoError(t, err)
	defer f.Close()

	tests := []struct {
		name          string
		n             int
		src           io.Reader
		expectedError string
	}{
		{
			name:          "Success",
			n:             4096,
			src:           f,
			expectedError: "",
		},
		{
			name:          "Failure",
			n:             4096,
			src:           iotest.ErrReader(errors.New("io error")),
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in, err := New(tc.n, tc.src)

			if tc.expectedError == "" {
				assert.NotNil(t, in)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, in)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestInput_loadFirst(t *testing.T) {
	tests := []struct {
		name          string
		i             *Input
		expectedError string
	}{
		{
			name: "Success",
			i: &Input{
				src:  strings.NewReader("Lorem ipsum"),
				buff: make([]byte, 2048),
			},
			expectedError: "",
		},
		{
			name: "Failure",
			i: &Input{
				src:  iotest.ErrReader(errors.New("io error")),
				buff: make([]byte, 2048),
			},
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.i.loadFirst()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestInput_loadSecond(t *testing.T) {
	tests := []struct {
		name          string
		i             *Input
		expectedError string
	}{
		{
			name: "Success",
			i: &Input{
				src:  strings.NewReader("Lorem ipsum"),
				buff: make([]byte, 2048),
			},
			expectedError: "",
		},
		{
			name: "Failure",
			i: &Input{
				src:  iotest.ErrReader(errors.New("io error")),
				buff: make([]byte, 2048),
			},
			expectedError: "io error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.i.loadSecond()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestInput_Next(t *testing.T) {
	tests := []struct {
		name          string
		n             int
		file          string
		expectedCount int
	}{
		{
			name:          "Success",
			n:             1024,
			file:          "./fixture/lorem_ipsum",
			expectedCount: 3422,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := New(tc.n, f)
			assert.NoError(t, err)

			var r rune
			var count int

			for r, err = in.Next(); err == nil; r, err = in.Next() {
				count++
				assert.NotEmpty(t, r)
			}

			assert.Equal(t, io.EOF, err)
			assert.Equal(t, tc.expectedCount, count)
		})
	}
}

func TestInput_Retract(t *testing.T) {
	tests := []struct {
		name         string
		n            int
		file         string
		lexemeBegin  int
		forward      int
		expectedPeek rune
	}{
		{
			name:         "Success",
			n:            1024,
			file:         "./fixture/lorem_ipsum",
			lexemeBegin:  0,
			forward:      10,
			expectedPeek: 'u',
		},
		{
			name:         "Success_SecondHalfToFirstHalf",
			n:            1024,
			file:         "./fixture/lorem_ipsum",
			lexemeBegin:  1020,
			forward:      1024,
			expectedPeek: 's',
		},
		{
			name:         "Success_FirstHalfToSecondHalf",
			n:            1024,
			file:         "./fixture/lorem_ipsum",
			lexemeBegin:  2040,
			forward:      0,
			expectedPeek: 'p',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := New(tc.n, f)
			assert.NoError(t, err)

			err = in.loadSecond()
			assert.NoError(t, err)

			in.lexemeBegin = tc.lexemeBegin
			in.forward = tc.forward

			in.Retract()

			assert.Equal(t, tc.expectedPeek, in.Peek())
		})
	}
}

func TestInput_Peek(t *testing.T) {
	tests := []struct {
		name         string
		n            int
		file         string
		lexemeBegin  int
		forward      int
		expectedRune rune
	}{
		{
			name:         "Success",
			n:            1024,
			file:         "./fixture/lorem_ipsum",
			lexemeBegin:  0,
			forward:      10,
			expectedRune: 'm',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := New(tc.n, f)
			assert.NoError(t, err)

			err = in.loadSecond()
			assert.NoError(t, err)

			in.lexemeBegin = tc.lexemeBegin
			in.forward = tc.forward

			r := in.Peek()
			assert.Equal(t, tc.expectedRune, r)
		})
	}
}

func TestInput_Lexeme(t *testing.T) {
	tests := []struct {
		name           string
		n              int
		file           string
		lexemePos      int
		lexemeBegin    int
		forward        int
		expectedLexeme string
		expectedPos    int
	}{
		{
			name:           "Success",
			n:              1024,
			file:           "./fixture/lorem_ipsum",
			lexemePos:      0,
			lexemeBegin:    0,
			forward:        5,
			expectedLexeme: "Lorem",
			expectedPos:    0,
		},
		{
			name:           "Success_FirstHalfToSecondHalf",
			n:              1024,
			file:           "./fixture/lorem_ipsum",
			lexemePos:      1020,
			lexemeBegin:    1020,
			forward:        1030,
			expectedLexeme: "us sceleri",
			expectedPos:    1020,
		},
		{
			name:           "Success_SecondHalfToFirstHalf",
			n:              1024,
			file:           "./fixture/lorem_ipsum",
			lexemePos:      4040,
			lexemeBegin:    2044,
			forward:        5,
			expectedLexeme: "corpLorem",
			expectedPos:    4040,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := New(tc.n, f)
			assert.NoError(t, err)

			err = in.loadSecond()
			assert.NoError(t, err)

			in.lexemePos = tc.lexemePos
			in.lexemeBegin = tc.lexemeBegin
			in.forward = tc.forward

			lexeme, pos := in.Lexeme()
			assert.Equal(t, tc.expectedLexeme, lexeme)
			assert.Equal(t, tc.expectedPos, pos)
		})
	}
}

func TestInput_Skip(t *testing.T) {
	tests := []struct {
		name                string
		n                   int
		file                string
		lexemePos           int
		lexemeBegin         int
		forward             int
		expectedLexemePos   int
		expectedLexemeBegin int
	}{
		{
			name:                "Success",
			n:                   1024,
			file:                "./fixture/lorem_ipsum",
			lexemePos:           0,
			lexemeBegin:         0,
			forward:             10,
			expectedLexemePos:   10,
			expectedLexemeBegin: 10,
		},
		{
			name:                "Success_FirstHalfToSecondHalf",
			n:                   1024,
			file:                "./fixture/lorem_ipsum",
			lexemePos:           1020,
			lexemeBegin:         1020,
			forward:             1030,
			expectedLexemePos:   1030,
			expectedLexemeBegin: 1030,
		},
		{
			name:                "Success_SecondHalfToFirstHalf",
			n:                   1024,
			file:                "./fixture/lorem_ipsum",
			lexemePos:           4040,
			lexemeBegin:         2044,
			forward:             5,
			expectedLexemePos:   4049,
			expectedLexemeBegin: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.file)
			assert.NoError(t, err)
			defer f.Close()

			in, err := New(tc.n, f)
			assert.NoError(t, err)

			err = in.loadSecond()
			assert.NoError(t, err)

			in.lexemePos = tc.lexemePos
			in.lexemeBegin = tc.lexemeBegin
			in.forward = tc.forward

			in.Skip()

			assert.Equal(t, tc.expectedLexemePos, in.lexemePos)
			assert.Equal(t, tc.expectedLexemeBegin, in.lexemeBegin)
		})
	}
}
