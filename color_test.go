package console

import (
	qt "github.com/frankban/quicktest"
	"os"
	"testing"
)

func TestColor_Apply(t *testing.T) {
	type cs struct {
		Label string
		Expected string
		Foreground string
		Background string
		Options []string
		TestDegraded bool
	}

	os.Setenv("COLORTERM", "truecolor")
	cases := []cs{
		{
			Label:      "ansi Color",
			Expected:   "\033[31;43m \033[39;49m",
			Foreground: "red",
			Background: "yellow",
		},
		{
			Label:      "ansi Color with bright-red and bright-yellow",
			Expected:   "\033[91;103m \033[39;49m",
			Foreground: "bright-red",
			Background: "bright-yellow",
		},
		{
			Label:      "ansi Color with underline options",
			Expected:   "\033[31;43;4m \033[39;49;24m",
			Foreground: "red",
			Background: "yellow",
			Options:    []string{"underscore"},
		},
		{
			Label: "true Color with #fff and #000",
			Foreground: "#fff",
			Background: "#000",
			Expected: "\033[38;2;255;255;255;48;2;0;0;0m \033[39;49m",
		},
		{
			Label: "true Color with #ffffff and #000000",
			Foreground: "#ffffff",
			Background: "#000000",
			Expected: "\033[38;2;255;255;255;48;2;0;0;0m \033[39;49m",
		},
		{
			Label: "degraded Color with #f00 and #ff0",
			Foreground: "#f00",
			Background: "#ff0",
			Expected: "\033[31;43m \033[39;49m",
			TestDegraded: true,
		},
		{
			Label: "degraded Color with #c0392b and #f1c40f",
			Foreground: "#c0392b",
			Background: "#f1c40f",
			Expected: "\033[31;43m \033[39;49m",
			TestDegraded: true,
		},
	}

	for _, v := range cases {
		t.Run(v.Label, func(t *testing.T){
			if true == v.TestDegraded {
				os.Setenv("COLORTERM", "nocolor")
			}
			check := qt.New(t)
			color := NewColor(v.Foreground, v.Background)
			if len(v.Options) > 0 {
				color = NewColorWithOptions(v.Foreground, v.Background, v.Options)
			}


			output := color.Apply(" ")
			check.Assert(output, qt.Equals, v.Expected)
		})
	}
}