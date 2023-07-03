package data

// TODO: Add tests for meta methods

// Meta - stores Meta for analytics and error positioning
type Meta struct {
	line int
	pos  int
	// TODO: add an input buffer for better error printing
}

// NewMeta - constructor for meta
func NewMeta() *Meta {
	return &Meta{
		line: 1,
	}
}

// Current - copies the values for the current meta
func (m *Meta) Current() *Meta {
	return &Meta{
		line: m.line,
		pos:  m.pos,
	}
}

// Inc - increments position
func (m *Meta) Inc() {
	m.pos++
}

// NewLine - increments line and resets position
func (m *Meta) NewLine() {
	m.line++
	m.pos = 0
}

// IncNL - increments line and resets position,
// if s == '\n', otherwise increments position
func (m *Meta) IncNL(s rune) {
	if s == '\n' {
		m.line++
		m.pos = 0
		return
	}

	m.pos++
}
