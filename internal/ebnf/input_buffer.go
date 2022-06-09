package ebnf

import "io"

const sentinel byte = 4

// inputBuffer implements the two-buffer scheme for reading the input characters.
type inputBuffer struct {
	src io.Reader

	// The first and second halves of the buff are alternatively reloaded.
	// Each half is of the same size N plus an additional space for the sentinel character.
	// Usually, N should be the size of a disk block (4096 bytes).
	buff []byte

	lexemeBegin int // Pointer lexemeBegin marks the beginning of the current lexeme.
	forward     int // Pointer forward scans ahead until a pattern match is found.

	char byte  // Last character read
	err  error // Last error encountered
}

// newInputBuffer creates a new input buffer of size N.
// N usually should be the size of a disk block (4096 bytes).
func newInputBuffer(size int, src io.Reader) (*inputBuffer, error) {
	// buff is divided into two sub-buffers (first half and second half).
	// Each sub-buffer has an additional space for the sentinel character.
	l := (size + 1) * 2
	buff := make([]byte, l)

	in := &inputBuffer{
		src:         src,
		buff:        buff,
		lexemeBegin: 0,
		forward:     0,
	}

	if err := in.loadFirst(); err != nil {
		return nil, err
	}

	return in, nil
}

// GetNextChar reads the next character from the input source.
// Next advances the input buffer to the next character, which will then be available through the Char method.
// It returns false when either an error occurs or the end of the input is reached.
// After Next returns false, the Err method will return any error that occurred.
func (i *inputBuffer) Next() bool {
	if i.err != nil {
		return false
	}

	i.char = i.buff[i.forward]
	i.forward++

	if i.buff[i.forward] == sentinel {
		if i.isForwardAtEndOfFirst() {
			if i.err = i.loadSecond(); i.err == nil {
				i.forward++
			}
		} else if i.isForwardAtEndOfSecond() {
			if i.err = i.loadFirst(); i.err == nil {
				i.forward = 0
			}
		} else {
			// Sentinel within a sub-buffer signifies the end of input
			i.err = io.EOF
		}
	}

	// The current read is valid, but the next read is not possible
	return true
}

// Char returns the most recent character read by a call to Next method.
func (i *inputBuffer) Char() byte {
	return i.char
}

// Err returns the first non-EOF error encountered by a call to Next method.
func (i *inputBuffer) Err() error {
	if i.err == io.EOF {
		return nil
	}
	return i.err
}

// isForwardAtEndOfFirst determines whether or not forward is at the end of the first half.
func (i *inputBuffer) isForwardAtEndOfFirst() bool {
	high := (len(i.buff) / 2) - 1
	return i.forward == high
}

// isForwardAtEndOfSecond determines whether or not forward is at the end of the second half.
func (i *inputBuffer) isForwardAtEndOfSecond() bool {
	high := len(i.buff) - 1
	return i.forward == high
}

// loadFirst reads the input and loads the first sub-buffer.
func (i *inputBuffer) loadFirst() error {
	high := (len(i.buff) / 2) - 1
	n, err := i.src.Read(i.buff[:high])
	if err != nil {
		return err
	}

	i.buff[n] = sentinel

	return nil
}

// loadSecond reads the input and loads the second sub-buffer.
func (i *inputBuffer) loadSecond() error {
	low, high := len(i.buff)/2, len(i.buff)-1
	n, err := i.src.Read(i.buff[low:high])
	if err != nil {
		return err
	}

	i.buff[low+n] = sentinel

	return nil
}
