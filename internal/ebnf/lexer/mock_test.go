package lexer

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
		OutPos int
	}

	SkipMock struct {
		OutPos int
	}

	mockInputBuffer struct {
		NextIndex int
		NextMocks []NextMock

		RetractIndex int

		PeekIndex int
		PeekMocks []PeekMock

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

func (m *mockInputBuffer) Peek() (rune, error) {
	i := m.PeekIndex
	m.PeekIndex++
	return m.PeekMocks[i].OutRune, m.PeekMocks[i].OutError
}

func (m *mockInputBuffer) Lexeme() (string, int) {
	i := m.LexemeIndex
	m.LexemeIndex++
	return m.LexemeMocks[i].OutVal, m.LexemeMocks[i].OutPos
}

func (m *mockInputBuffer) Skip() int {
	i := m.SkipIndex
	m.SkipIndex++
	return m.SkipMocks[i].OutPos
}
