package style

import (
	qt "github.com/frankban/quicktest"
	"github.com/kilip/go-console/formatter"
	"github.com/kilip/go-console/output"
	"strings"
	"testing"
)

type buffMock struct {
	Output string
	Input  string
}

func (bm *buffMock) Write(p []byte) (n int, err error) {
	bm.Output += string(p)
	return 0, nil
}

func (bm *buffMock) Read(p []byte) (n int, err error) {
	bm.Input += string(p)
	return 0, err
}

func (bm *buffMock) Reset() {
	bm.Output = ""
	bm.Input = ""
}

type cs struct {
	Name     string
	Expected string
}

func NewReadWriterMock() *buffMock {
	return new(buffMock)
}

func TestOutputStyle_NewLineC(t *testing.T) {
	ch := qt.New(t)
	buff := NewReadWriterMock()
	out := output.NewStreamOutput(buff, formatter.NewFormatter())
	os := &OutputStyle{input: buff, IOutput: out}

	os.NewLine()
	ch.Assert(buff.Output, qt.Equals, "\n")

	os.NewLineC(4)
	ch.Assert(buff.Output, qt.Equals, strings.Repeat("\n", 5))
}
