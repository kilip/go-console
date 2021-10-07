package output

import (
	"github.com/kilip/go-console/formatter"
	"github.com/kilip/go-console/helper"
)

// Output is base class for output classes.
// There are five levels of verbosity:
// * normal: no option passed given the normal output
// * verbose: more output
// * very verbose: highly extended output
// * debug: all debug output
// * quite: no output at all
type Output struct {
	verbosity int
	formatter *formatter.Formatter
	doWrite   func(message string, newLine bool)
}

// NewOutput creates and returns new Output class
func NewOutput(formatter *formatter.Formatter) *Output {
	o := &Output{
		verbosity: VerbosityNormal,
		formatter: formatter,
	}
	return o
}

// SetVerbosity sets the verbosity of the output
func (o *Output) SetVerbosity(verbosity int) {
	o.verbosity = verbosity
}

// GetVerbosity gets current verbosity of the output
func (o *Output) GetVerbosity() int {
	return o.verbosity
}

func (o *Output) SetFormatter(formatter *formatter.Formatter) {
	o.formatter = formatter
}

// GetFormatter returns current output formatter.Formatter instance
func (o *Output) GetFormatter() *formatter.Formatter {
	return o.formatter
}

// SetDecorated sets the Output decorated flag
func (o *Output) SetDecorated(decorated bool) {
	o.formatter.SetDecorated(decorated)
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
func (o *Output) Write(message interface{}) {
	o.WriteO(message, false, FormatNormal)
}

// WriteO writes a message into the output with defined options
func (o *Output) WriteO(message interface{}, newLine bool, options int) {
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
		return
	}

	messages := helper.TextToSlices(message)
	for _, m := range messages {
		formatted := m
		switch formatType {
		case FormatNormal:
			formatted = o.formatter.Format(formatted)
			break
		}
		o.doWrite(formatted, newLine)
	}
}

// Writeln writes a message into the output and adds a new line at the end
func (o *Output) Writeln(message string) {
	o.WritelnO(message, FormatNormal)
}

// WritelnO writes a message into the output and adds a new line at the end
// with given options behavior
func (o *Output) WritelnO(message string, options int) {
	o.WriteO(message, true, options)
}
