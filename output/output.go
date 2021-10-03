package output

import (
	"github.com/kilip/console/formatter"
	"io"
)

//Verbosity mode
const (
	VerbosityQuiet       = 16
	VerbosityNormal      = 32
	VerbosityVerbose     = 64
	VerbosityVeryVerbose = 128
	VerbosityDebug       = 256
)

//Format output
const (
	FormatNormal = 1
	FormatRaw    = 2
	FormatPlain  = 4
)

// Output is base class for output classes.
// There are five levels of verbosity:
// * normal: no option passed given the normal output
// * verbose: more output
// * very verbose: highly extended output
// * debug: all debug output
// * quite: no output at all
type Output struct {
	writer    io.Writer
	verbosity int
	formatter *formatter.Formatter
}

// NewOutput creates and returns new Output class
func NewOutput(writer io.Writer, formatter *formatter.Formatter) *Output {
	return &Output{
		writer:    writer,
		verbosity: VerbosityNormal,
		formatter: formatter,
	}
}

// SetVerbosity sets the verbosity of the output
func (o *Output) SetVerbosity(verbosity int) {
	o.verbosity = verbosity
}

// GetVerbosity gets current verbosity of the output
func (o *Output) GetVerbosity() int {
	return o.verbosity
}

// GetFormatter returns current output formatter.Formatter instance
func (o *Output) GetFormatter() *formatter.Formatter {
	return o.formatter
}

// IsDecorated returns whether this Output is decorated
func (o *Output) IsDecorated() bool {
	return o.formatter.IsDecorated()
}

// IsQuite returns whether verbosity is quite
func (o *Output) IsQuite() bool {
	return VerbosityQuiet == o.verbosity
}

// IsVerbose returns whether verbosity is verbose
func (o *Output) IsVerbose() bool {
	return VerbosityVerbose <= o.verbosity
}

// IsVeryVerbose returns whether verbosity is very verbose
func (o *Output) IsVeryVerbose() bool {
	return VerbosityVeryVerbose <= o.verbosity
}

// IsDebug returns whether verbosity is debug
func (o *Output) IsDebug() bool {
	return VerbosityDebug <= o.verbosity
}

// Write writes a message into the output
func (o *Output) Write(message string) error {
	return o.WriteO(message, FormatNormal)
}

// WriteO writes a message into the output with defined options
func (o *Output) WriteO(message string, options int) error {
	return o.doWrite(message, false, options)
}

// Writeln writes a message into the output and adds a new line at the end
func (o *Output) Writeln(message string) error {
	return o.WritelnO(message, FormatNormal)
}

// WritelnO writes a message into the output and adds a new line at the end
// with given options behavior
func (o *Output) WritelnO(message string, options int) error {
	return o.doWrite(message, true, options)
}

// doWrite perform an actual write to the io.Writer output
func (o *Output) doWrite(message string, newLine bool, options int) error {
	formatted := message
	types := FormatNormal | FormatRaw | FormatPlain
	formatType := types & options

	if 0 == formatType {
		formatType = FormatNormal
	}

	verbosities := VerbosityQuiet | VerbosityNormal | VerbosityVerbose | VerbosityVeryVerbose | VerbosityDebug
	verbosity := verbosities & options

	if 0 == verbosity {
		verbosity = VerbosityNormal
	}

	if verbosity > o.verbosity {
		return nil
	}

	switch formatType {
	case FormatNormal:
		formatted = o.formatter.Format(formatted)
		break
	}

	if newLine {
		formatted += "\n"
	}

	_, e := o.writer.Write([]byte(formatted))

	return e
}
