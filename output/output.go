package output

import "io"

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

type Formatter interface {
	Format(message string) string
}

type Output struct {
	writer    io.Writer
	verbosity int
	formatter Formatter
}

func NewOutput(writer io.Writer, formatter Formatter) *Output {
	return &Output{
		writer:    writer,
		verbosity: VerbosityNormal,
		formatter: formatter,
	}
}

func (o *Output) SetVerbosity(verbosity int) {
	o.verbosity = verbosity
}

func (o *Output) GetVerbosity() int {
	return o.verbosity
}

func (o *Output) GetFormatter() Formatter {
	return o.formatter
}

func (o *Output) IsQuite() bool {
	return VerbosityQuiet == o.verbosity
}

func (o *Output) IsVerbose() bool {
	return VerbosityVerbose <= o.verbosity
}

func (o *Output) IsVeryVerbose() bool {
	return VerbosityVeryVerbose <= o.verbosity
}

func (o *Output) IsDebug() bool {
	return VerbosityDebug <= o.verbosity
}

func (o *Output) Write(message string) error {
	return o.WriteO(message, FormatNormal)
}

func (o *Output) WriteO(message string, options int) error {
	return o.doWrite(message, false, options)
}

func (o *Output) Writeln(message string) error {
	return o.WritelnO(message, FormatNormal)
}

func (o *Output) WritelnO(message string, options int) error {
	return o.doWrite(message, true, options)
}

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
