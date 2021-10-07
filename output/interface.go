package output

import "github.com/kilip/go-console/formatter"

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

type IOutput interface {
	SetVerbosity(level int)
	GetVerbosity() int
	SetFormatter(formatter *formatter.Formatter)
	GetFormatter() *formatter.Formatter
	SetDecorated(decorated bool)
	IsDecorated() bool
	IsQuite() bool
	IsVerbose() bool
	IsVeryVerbose() bool
	Write(message interface{})
	WriteO(message interface{}, newLine bool, options int)
	Writeln(message string)
	WritelnO(message string, options int)
}
