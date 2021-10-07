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
		return ""
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
			Name: "case01",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Title("Title")
				ds.Warning(input)
				ds.Title("Title")
			},
		},
		{
			Name: "case02",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Warning("Warning")
				ds.Caution("Caution")
				ds.Error("Error")
				ds.Success("Success")
				ds.Note("Note")
				ds.Info("Info")
				ds.BlockO("Custom block", "CUSTOM", "fg=white;bg=green", "X ", true, true)
			},
		},
		{
			Name: "case03",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Title("First title")
				ds.Title("Second title")
			},
		},
		{
			Name: "case04",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Write("Lorem ipsum dolor sit amet")
				ds.Title("First title")

				ds.Writeln("Lorem ipsum dolor sit amet")
				ds.Title("Second title")

				ds.Writeln("Lorem ipsum dolor sit amet")
				ds.Write("")
				ds.Title("Third title")

				//Ensure edge case by appending empty strings to history:
				ds.Write("Lorem ipsum dolor sit amet")
				ds.Write("")
				ds.Title("Fourth title")

				//Ensure have manual control over number of blank lines:
				ds.Writeln("Lorem ipsum dolor sit amet")
				ds.Writeln("")
				ds.Title("Fifth title")

				//Should append an extra blank line
				ds.Writeln("Lorem ipsum dolor sit amet")
				ds.NewLineC(2)
				ds.Title("Fifth title")
			},
		},
		{
			Name: "case05",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)

				ds.Writeln("Lorem ipsum dolor sit amet")
				ds.Listing([]string{
					"Lorem ipsum dolor sit amet",
					"consectetur adipiscing elit",
				})

				// even using write
				ds.Write("Lorem ipsum dolor sit amet")
				ds.Listing([]string{
					"Lorem ipsum dolor sit amet",
					"consectetur adipiscing elit",
				})

				ds.Write("Lorem ipsum dolor sit amet")
				ds.Text([]string{
					"Lorem ipsum dolor sit amet",
					"consectetur adipiscing elit",
				})

				ds.NewLine()

				ds.Write("Lorem ipsum dolor sit amet")
				ds.Comment([]string{
					"Lorem ipsum dolor sit amet",
					"consectetur adipiscing elit",
				})
			},
		},
		{
			Name: "case06",
			Style: func(testCase outputTestCase, ds *DefaultStyle, input string) {
				ds.SetDecorated(false)
				ds.Listing([]string{
					"Lorem ipsum dolor sit amet",
					"consectetur adipiscing elit",
				})
				ds.Success("Lorem ipsum dolor sit amet")
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
			expected := getFileContents(tCase.Name + ".out.txt")
			tCase.Style(tCase, ds, getFileContents(tCase.Name+".in.txt"))
			c.Assert(buffMock.Output+"\n", qt.Equals, expected)
		})
	}
}
