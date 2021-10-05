package output

import (
	qt "github.com/frankban/quicktest"
	"github.com/kilip/go-console/formatter"
	"testing"
)

type TestOutput struct {
	Buffer string
	*Output
}

func (to *TestOutput) doWrite(message string, newLine bool) {
	to.Buffer += message

	if newLine {
		to.Buffer += "\n"
	}
}

func NewTestOutput() *TestOutput {
	output := NewOutput(formatter.NewFormatter())
	o := &TestOutput{
		Output: output,
	}

	output.doWrite = o.doWrite

	return o
}

func TestOutput_Verbosity(t *testing.T) {
	type cs struct {
		Label         string
		Verbosity     int
		IsQuite       bool
		IsVerbose     bool
		IsVeryVerbose bool
		IsDebug       bool
	}

	cases := []cs{
		{
			Label:         "quite",
			Verbosity:     VerbosityQuiet,
			IsQuite:       true,
			IsVerbose:     false,
			IsVeryVerbose: false,
			IsDebug:       false,
		},
		{
			Label:         "normal",
			Verbosity:     VerbosityNormal,
			IsQuite:       false,
			IsVerbose:     false,
			IsVeryVerbose: false,
			IsDebug:       false,
		},
		{
			Label:         "verbose",
			Verbosity:     VerbosityVerbose,
			IsQuite:       false,
			IsVerbose:     true,
			IsVeryVerbose: false,
			IsDebug:       false,
		},
		{
			Label:         "very verbose",
			Verbosity:     VerbosityVeryVerbose,
			IsQuite:       false,
			IsVerbose:     true,
			IsVeryVerbose: true,
			IsDebug:       false,
		},
		{
			Label:         "debug",
			Verbosity:     VerbosityDebug,
			IsQuite:       false,
			IsVerbose:     true,
			IsVeryVerbose: true,
			IsDebug:       true,
		},
	}

	for _, v := range cases {
		t.Run(v.Label, func(t *testing.T) {
			o := NewTestOutput()
			o.SetVerbosity(v.Verbosity)
			cc := qt.New(t)

			cc.Assert(o.IsQuite(), qt.Equals, v.IsQuite)
			cc.Assert(o.IsVerbose(), qt.Equals, v.IsVerbose)
			cc.Assert(o.IsVeryVerbose(), qt.Equals, v.IsVeryVerbose)
			cc.Assert(o.IsDebug(), qt.Equals, v.IsDebug)
		})
	}
}

func TestOutput_Write(t *testing.T) {
	type cs struct {
		Name      string
		Verbosity int
		Expected  string
	}

	cases := []cs{
		{
			Name:      "quite",
			Verbosity: VerbosityQuiet,
			Expected:  "2",
		},
		{
			Name:      "verbose",
			Verbosity: VerbosityNormal,
			Expected:  "123",
		},
		{
			Name:      "verbose",
			Verbosity: VerbosityVerbose,
			Expected:  "1234",
		},
		{
			Name:      "verbose",
			Verbosity: VerbosityVeryVerbose,
			Expected:  "12345",
		},
		{
			Name:      "verbose",
			Verbosity: VerbosityDebug,
			Expected:  "123456",
		},
	}

	for _, val := range cases {
		t.Run(val.Name, func(t *testing.T) {
			cc := qt.New(t)
			o := NewTestOutput()
			o.SetVerbosity(val.Verbosity)

			o.Write("1")
			o.WriteO("2", false, VerbosityQuiet)
			o.WriteO("3", false, VerbosityNormal)
			o.WriteO("4", false, VerbosityVerbose)
			o.WriteO("5", false, VerbosityVeryVerbose)
			o.WriteO("6", false, VerbosityDebug)

			cc.Assert(o.Buffer, qt.Equals, val.Expected)
		})
	}

}
