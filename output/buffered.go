package output

import "github.com/kilip/go-console/formatter"

// Buffered represents a buffered output struct
type Buffered struct {
	buffer string
	*Output
}

// NewBuffered creates new buffered output
func NewBuffered(formatter *formatter.Formatter) *Buffered {
	output := NewOutput(formatter)
	buffered := &Buffered{
		Output: output,
	}
	output.doWrite = buffered.doWrite

	return buffered
}

// doWrite performs an actual writes for Output.doWrite function
func (b *Buffered) doWrite(message string, newLine bool) {
	b.buffer += message

	if newLine {
		b.buffer += message
	}
}

// Fetch empties Buffer and returns its content
func (b *Buffered) Fetch() string {
	content := b.buffer
	b.buffer = ""

	return content
}
