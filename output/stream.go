package output

import (
	"github.com/kilip/go-console/formatter"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"strconv"
)

// Stream writes output by using given io.Writer
type Stream struct {
	writer io.Writer
	*Output
}

// NewStreamOutput creates new Stream output object
func NewStreamOutput(writer io.Writer, formatter *formatter.Formatter) *Stream {
	decorated := hasColorSupport()
	formatter.SetDecorated(decorated)

	output := NewOutput(formatter)
	stream := &Stream{
		writer: writer,
		Output: output,
	}
	output.doWrite = stream.doWrite

	return stream
}

// doWrite perform an actual write for Output.doWrite
func (so *Stream) doWrite(message string, newLine bool) {
	_, _ = so.writer.Write([]byte(message))
	if newLine {
		_, _ = so.writer.Write([]byte("\n"))
	}
}

// GetWriter returns io.Writer attached to this Stream instance
func (so *Stream) GetWriter() io.Writer {
	return so.writer
}

// hasColorSupport determines if current environment has color supports
func hasColorSupport() bool {
	/*
		TODO: remove isatty dependencies
	*/
	noColorEnv, err := strconv.ParseBool(os.Getenv("NO_COLOR"))
	if err == nil {
		return noColorEnv
	}

	if "Hyper" == os.Getenv("TERM_PROGRAM") {
		return true
	}

	return isatty.IsTerminal(os.Stdout.Fd())
}
