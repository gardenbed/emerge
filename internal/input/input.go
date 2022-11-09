package input

import (
	"io"
	"strings"
)

const eof byte = 0x00

// Input implements the two-buffer scheme for reading the input characters.
//
// For more details, see Compilers: Principles, Techniques, and Tools (2nd Edition).
type Input struct {
	src io.Reader

	// The first and second halves of the buff are alternatively reloaded.
	// Each half is of the same size N. Usually, N should be the size of a disk block.
	buff []byte

	lexemePos   int // Counter lexemePos tracks the position of the current lexeme in the input file.
	lexemeBegin int // Pointer lexemeBegin marks the beginning of the current lexeme.
	forward     int // Pointer forward scans ahead until a pattern match is found.

	err error // Last error encountered
}

// New creates a new input buffer of size N.
// N usually should be the size of a disk block.
func New(n int, src io.Reader) (*Input, error) {
	// buff is divided into two sub-buffers (first half and second half).
	l := 2 * n
	buff := make([]byte, l)

	in := &Input{
		src:         src,
		buff:        buff,
		lexemePos:   0,
		lexemeBegin: 0,
		forward:     0,
	}

	if err := in.loadFirst(); err != nil {
		return nil, err
	}

	return in, nil
}

// loadFirst reads the input and loads the first sub-buffer.
func (i *Input) loadFirst() error {
	high := len(i.buff) / 2
	n, err := i.src.Read(i.buff[:high])
	if err != nil {
		return err
	}

	if n < high {
		i.buff[n] = eof
	}

	return nil
}

// loadSecond reads the input and loads the second sub-buffer.
func (i *Input) loadSecond() error {
	low, high := len(i.buff)/2, len(i.buff)
	n, err := i.src.Read(i.buff[low:high])
	if err != nil {
		return err
	}

	if n < high-low {
		i.buff[low+n] = eof
	}

	return nil
}

// Next advances to the next rune in the input and returns it.
func (i *Input) Next() (rune, error) {
	if i.err != nil {
		return 0, i.err
	}

	r := i.buff[i.forward]
	i.forward++

	// Determine whether or not the forward pointer has reached the end of any halves.
	// If so, it loads the other half and set the forward pointer to the beginning of it.
	// If the forward pointer has reached to the end of input, an io.EOF error will be returned.
	if i.forward == len(i.buff)/2 { // Is forward at the end of first half?
		i.err = i.loadSecond()
	} else if i.forward == len(i.buff) { // Is forward at the end of second half?
		if i.err = i.loadFirst(); i.err == nil {
			i.forward = 0 // beginning of the first half
		}
	} else if i.buff[i.forward] == eof {
		i.err = io.EOF
	}

	// The current read is fine, but the next read may return an error
	return rune(r), nil
}

// Retract recedes to the previous rune in the input.
// It can only be called once per each call of Next.
func (i *Input) Retract() {
	if i.forward == 0 { // Is forward at the beginning of first half?
		i.forward = len(i.buff) - 1 // end of the second half
	} else {
		i.forward--
	}
}

// Peek returns the next rune in the input without consuming it.
func (i *Input) Peek() rune { // Is forward at the end of second half?
	r := i.buff[i.forward]

	return rune(r)
}

// Lexeme returns the current lexeme alongside its position.
func (i *Input) Lexeme() (string, int) {
	var lexeme strings.Builder
	pos := i.lexemePos

	for i.lexemeBegin != i.forward {
		lexeme.WriteByte(i.buff[i.lexemeBegin])

		i.lexemePos++
		i.lexemeBegin++

		if i.lexemeBegin == len(i.buff) { // Is lexemeBegin at the end of second half?
			i.lexemeBegin = 0 // beginning of the first half
		}
	}

	return lexeme.String(), pos
}

// Skip skips over the pending lexeme in the input.
func (i *Input) Skip() {
	if d := i.forward - i.lexemeBegin; d >= 0 {
		i.lexemePos += d
	} else {
		i.lexemePos += len(i.buff) + d
	}

	i.lexemeBegin = i.forward
}
