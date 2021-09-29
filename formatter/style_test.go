package formatter

import (
	qt "github.com/frankban/quicktest"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	Name       string
	Foreground string
	Background string
	Options    []string
	Expected   string
	Text       string
}

func TestStyle_Apply(t *testing.T) {
	cases := []testCase{
		{
			Name:       "constructor",
			Foreground: "green",
			Background: "black",
			Options:    []string{"bold", "underscore"},
			Expected:   "\033[32;40;1;4mfoo\033[39;49;22;24m",
		},
		{
			Name:       "constructor",
			Foreground: "red",
			Background: "",
			Options:    []string{"blink"},
			Expected:   "\033[31;5mfoo\033[39;25m",
		},
		{
			Name:       "constructor",
			Foreground: "",
			Background: "white",
			Expected:   "\033[47mfoo\033[49m",
		},
		{
			Name:       "foreground",
			Foreground: "black",
			Expected:   "\033[30mfoo\033[39m",
		},
		{
			Name:       "foreground",
			Foreground: "blue",
			Expected:   "\033[34mfoo\033[39m",
		},
		{
			Name:       "foreground",
			Foreground: "default",
			Expected:   "\033[39mfoo\033[39m",
		},
		{
			Name:       "background",
			Background: "black",
			Expected:   "\033[40mfoo\033[49m",
		},
		{
			Name:       "background",
			Background: "yellow",
			Expected:   "\033[43mfoo\033[49m",
		},
		{
			Name:       "background",
			Background: "default",
			Expected:   "\033[49mfoo\033[49m",
		},
		{
			Name:     "options_reverse_conceal",
			Options:  []string{"reverse", "conceal"},
			Expected: "\033[7;8mfoo\033[27;28m",
		},
	}

	for _, v := range cases {
		name := strings.Join([]string{v.Name, v.Foreground, v.Background}, "_")
		t.Run(name, func(t *testing.T) {
			checker := qt.New(t)
			style := NewFormatterStyleO(
				v.Foreground,
				v.Background,
				v.Options,
			)

			if "" == v.Text {
				v.Text = "foo"
			}
			output := style.Apply(v.Text)
			checker.Assert(output, qt.Equals, v.Expected)
		})
	}
}

func TestStyle_Options(t *testing.T) {
	checker := qt.New(t)
	style := NewFormatterStyleO("", "", []string{"reverse", "conceal"})

	checker.Assert(style.Apply("foo"), qt.Equals, "\033[7;8mfoo\033[27;28m")

	style.SetOption("bold")
	checker.Assert(style.Apply("foo"), qt.Equals, "\033[7;8;1mfoo\033[27;28;22m")

	style.UnsetOption("reverse")
	checker.Assert(style.Apply("foo"), qt.Equals, "\033[8;1mfoo\033[28;22m")

	style.SetOption("bold")
	checker.Assert(style.Apply("foo"), qt.Equals, "\033[8;1mfoo\033[28;22m")

	style.SetOptions([]string{"bold"})
	checker.Assert(style.Apply("foo"), qt.Equals, "\033[1mfoo\033[22m")
}

func TestStyle_Apply_Href(t *testing.T) {
	prevEmulator := os.Getenv("TERMINAL_EMULATOR")
	os.Setenv("TERMINAL_EMULATOR", "")

	checker := qt.New(t)
	style := NewFormatterStyleO("", "", []string{})
	style.SetHref("idea://open/?file=/path/SomeFile.php&line=12")

	expected := "\x1b]8;;idea://open/?file=/path/SomeFile.php&line=12\x1bsome URL\x1b]8;;\x1b\\"
	checker.Assert(style.Apply("some URL"), qt.Equals, expected)

	os.Setenv("TERMINAL_EMULATOR", prevEmulator)
}
