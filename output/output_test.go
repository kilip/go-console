package output

import (
	qt "github.com/frankban/quicktest"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockObject struct {
	mock.Mock
	Buffer string
}

func (w *mockObject) Write(p []byte) (nn int, err error) {
	w.Buffer += string(p)
	return 0, nil
}

func (w *mockObject) Format(message string) string {
	//w.Called(message)
	return message
}

func TestOutput_Verbosity(t *testing.T) {
	type cs struct {
		Label string
		Verbosity int
		IsQuite bool
		IsVerbose bool
		IsVeryVerbose bool
		IsDebug bool
	}

	cases := []cs{
		{
			Label: "quite",
			Verbosity: VerbosityQuiet,
			IsQuite: true,
			IsVerbose: false,
			IsVeryVerbose: false,
			IsDebug: false,

		},
		{
			Label: "normal",
			Verbosity: VerbosityNormal,
			IsQuite: false,
			IsVerbose: false,
			IsVeryVerbose: false,
			IsDebug: false,
		},
		{
			Label: "verbose",
			Verbosity: VerbosityVerbose,
			IsQuite: false,
			IsVerbose: true,
			IsVeryVerbose: false,
			IsDebug: false,
		},
		{
			Label: "very verbose",
			Verbosity: VerbosityVeryVerbose,
			IsQuite: false,
			IsVerbose: true,
			IsVeryVerbose: true,
			IsDebug: false,
		},
		{
			Label: "debug",
			Verbosity: VerbosityDebug,
			IsQuite: false,
			IsVerbose: true,
			IsVeryVerbose: true,
			IsDebug: true,
		},
	}

	for _, v := range cases {
		t.Run(v.Label, func(t *testing.T){
			w := new(mockObject)
			o := NewOutput(w, w)
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
		Name string
		Verbosity int
		Expected string
	}

	cases := []cs {
		{
			Name: "quite",
			Verbosity: VerbosityQuiet,
			Expected: "2",
		},
		{
			Name: "verbose",
			Verbosity: VerbosityNormal,
			Expected: "123",
		},
		{
			Name: "verbose",
			Verbosity: VerbosityVerbose,
			Expected: "1234",
		},
		{
			Name: "verbose",
			Verbosity: VerbosityVeryVerbose,
			Expected: "12345",
		},
		{
			Name: "verbose",
			Verbosity: VerbosityDebug,
			Expected: "123456",
		},
	}

	for _, val := range cases {
		t.Run(val.Name, func(t *testing.T){
			cc := qt.New(t)
			w := new(mockObject)
			o := NewOutput(w, w)

			o.SetVerbosity(val.Verbosity)

			w.On("Format")
			w.On("Write").Return(0, nil)

			o.Write("1")
			o.WriteO("2", VerbosityQuiet)
			o.WriteO("3", VerbosityNormal)
			o.WriteO("4", VerbosityVerbose)
			o.WriteO("5", VerbosityVeryVerbose)
			o.WriteO("6", VerbosityDebug)

			cc.Assert(w.Buffer, qt.Equals, val.Expected)
		})
	}

}