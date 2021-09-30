package style

import (
	qt "github.com/frankban/quicktest"
	"github.com/kilip/console/formatter"
	"github.com/kilip/console/output"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

type buffMock struct {
	Output string
	Input  string
	mock.Mock
}

func (bm *buffMock) Write(p []byte) (n int, err error) {
	bm.Called(p)
	bm.Output += string(p)
	return 0, nil
}

func (bm *buffMock) Read(p []byte) (n int, err error) {
	bm.Called(p)
	bm.Input += string(p)
	return 0, err
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
	out := output.NewOutput(buff, formatter.NewFormatter())
	os := &OutputStyle{input: buff, Output: out}

	buff.On("Write", []byte("\n"))
	err := os.NewLine()
	ch.Assert(err, qt.IsNil)
	ch.Assert(buff.Output, qt.Equals, "\n")

	buff.On("Write", []byte(strings.Repeat("\n", 4)))
	err = os.NewLineC(4)
	ch.Assert(err, qt.IsNil)
	ch.Assert(buff.Output, qt.Equals, strings.Repeat("\n", 5))

	err = os.Write("\n")
	ch.Assert(err, qt.IsNil)
}
