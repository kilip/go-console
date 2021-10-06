package style

import (
	qt "github.com/frankban/quicktest"
	"github.com/kilip/go-console/formatter"
	"github.com/kilip/go-console/output"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func getFileContents(path string) string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	fileName := basePath + "/testdata/" + path
	if bVal, err := os.ReadFile(fileName); err == nil {
		return string(bVal)
	} else {
		panic(err)
	}
}

func TestNewDefaultStyle(t *testing.T) {
	c := qt.New(t)
	buffMock := new(buffMock)
	o := output.NewStreamOutput(buffMock, formatter.NewFormatter())
	ds := NewDefaultStyle(buffMock, o)

	ds.Write("hello world")
	c.Assert(buffMock.Output, qt.Equals, "hello world")
}

type outputTestCase struct {
	Name,
	Expected,
	Input string
	Style func(testCase outputTestCase, ds *DefaultStyle, input string)
}

func getBlockTestCase() []outputTestCase {
	return []outputTestCase{
		{
			Name: "case00",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Caution(input)
			},
		},
		{
			Name: "case10",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.BlockO(
					input,
					"CUSTOM",
					"fg=white;bg=green",
					"X ",
					true,
					true,
				)
			},
		},
	}
}

func TestDefaultStyle_Output(t *testing.T) {
	cases := getBlockTestCase()

	if runtime.GOOS == "windows" {
		t.Skip("Skip test on windows")
	}

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			c := qt.New(t)
			buffMock := new(buffMock)
			out := output.NewStreamOutput(buffMock, formatter.NewFormatter())
			ds := NewDefaultStyle(buffMock, out)
			tCase.Style(tCase, ds, getFileContents(tCase.Name+".in.txt"))
			c.Assert(buffMock.Output, qt.Equals, getFileContents(tCase.Name+".out.txt"))
		})
	}
}
