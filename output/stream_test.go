package output

import (
	qt "github.com/frankban/quicktest"
	"github.com/kilip/go-console/formatter"
	"testing"
)

type writerMock struct {
	Buffer     string
	BufferRead string
}

func (wm *writerMock) Write(p []byte) (nn int, err error) {
	wm.Buffer += string(p)
	return 0, nil
}

func (wm *writerMock) Read(p []byte) (n int, err error) {
	wm.BufferRead += string(p)
	return 0, nil
}

func TestNewStreamOutput(t *testing.T) {
	c := qt.New(t)
	wm := &writerMock{}
	o := NewStreamOutput(wm, formatter.NewFormatter())

	c.Assert(o.GetVerbosity(), qt.Equals, VerbosityNormal)
	c.Assert(o.IsDecorated(), qt.IsFalse)
	c.Assert(o.GetWriter(), qt.Equals, wm)
}

func TestStream_Write(t *testing.T) {
	c := qt.New(t)
	wm := &writerMock{}
	o := NewStreamOutput(wm, formatter.NewFormatter())

	o.Writeln("foo")
	c.Assert(wm.Buffer, qt.Equals, "foo\n")
}
