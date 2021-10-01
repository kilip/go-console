package style

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestNewDefaultStyle(t *testing.T) {
	c := qt.New(t)
	buffMock := new(buffMock)
	ds := NewDefaultStyle(buffMock, buffMock)

	ds.Write("hello world")
	c.Assert(buffMock.Output, qt.Equals, "hello world")
}
