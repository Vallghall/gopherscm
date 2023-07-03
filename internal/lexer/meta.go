package lexer

// meta - stores meta for analytics and error positioning
type meta struct {
	line int
	pos  int
}

func newMeta() *meta {
	return &meta{
		line: 1,
	}
}

// inc - increments position
func (m *meta) inc() {
	m.pos++
}

// newLine - increments line and resets position
func (m *meta) newLine() {
	m.line++
	m.pos = 0
}

// incNL - increments line and resets position,
// if s == '\n', otherwise increments position
func (m *meta) incNL(s rune) {
	if s == '\n' {
		m.line++
		m.pos = 0
		return
	}

	m.pos++
}
