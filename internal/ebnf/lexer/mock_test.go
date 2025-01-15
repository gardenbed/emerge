package lexer

import "github.com/moorara/algo/lexer"

type (
	NextMock struct {
		OutRune  rune
		OutError error
	}

	PeekMock struct {
		OutRune  rune
		OutError error
	}

	LexemeMock struct {
		OutVal string
		OutPos lexer.Position
	}

	SkipMock struct {
		OutPos lexer.Position
	}

	mockInputBuffer struct {
		NextIndex int
		NextMocks []NextMock

		RetractIndex int

		LexemeIndex int
		LexemeMocks []LexemeMock

		SkipIndex int
		SkipMocks []SkipMock
	}
)

func (m *mockInputBuffer) Next() (rune, error) {
	i := m.NextIndex
	m.NextIndex++
	return m.NextMocks[i].OutRune, m.NextMocks[i].OutError
}

func (m *mockInputBuffer) Retract() {
	m.RetractIndex++
}

func (m *mockInputBuffer) Lexeme() (string, lexer.Position) {
	i := m.LexemeIndex
	m.LexemeIndex++
	return m.LexemeMocks[i].OutVal, m.LexemeMocks[i].OutPos
}

func (m *mockInputBuffer) Skip() lexer.Position {
	i := m.SkipIndex
	m.SkipIndex++
	return m.SkipMocks[i].OutPos
}
