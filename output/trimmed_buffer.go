package output

import "github.com/kilip/go-console/formatter"

type TrimmedBufferOutput struct {
	maxLength int
	buffer    string
	*Output
}

func NewTrimmedBufferOutput(maxLength int) *TrimmedBufferOutput {
	output := NewOutput(formatter.NewFormatter())
	trimmedBuffer := &TrimmedBufferOutput{
		maxLength: maxLength,
		Output:    output,
	}

	output.doWrite = trimmedBuffer.doWrite

	return trimmedBuffer
}

// doWrite performs an actual writes for Output.doWrite function
func (tb *TrimmedBufferOutput) doWrite(message string, newLine bool) {
	tb.buffer += message

	if newLine {
		tb.buffer += "\n"
	}

	tb.buffer = tb.buffer[0:tb.maxLength]
}

// Fetch empties buffer and returns its content
func (tb *TrimmedBufferOutput) Fetch() string {
	content := tb.buffer
	tb.buffer = ""

	return content
}
