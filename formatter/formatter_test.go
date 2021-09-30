/*
	TODO: convert this symfony test case:
    public function testFormatToStringObject()
    {
        $formatter = new OutputFormatter(false);
        $this->assertEquals(
            'some info', $formatter->format(new TableCell())
        );
    }
*/

package formatter

import (
	"fmt"
	qt "github.com/frankban/quicktest"
	"os"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	checker := qt.New(t)
	f := NewFormatter()

	checker.Assert(f.IsDecorated(), qt.IsTrue)
	checker.Assert(f.HasStyle("error"), qt.IsTrue)
	checker.Assert(f.HasStyle("info"), qt.IsTrue)
	checker.Assert(f.HasStyle("comment"), qt.IsTrue)
	checker.Assert(f.HasStyle("question"), qt.IsTrue)
}

func TestFormatter_Decorated(t *testing.T) {
	checker := qt.New(t)
	f := NewFormatter()

	checker.Assert(f.IsDecorated(), qt.IsTrue)
	f.SetDecorated(false)
	checker.Assert(f.IsDecorated(), qt.IsFalse)
}

func TestFormatter_Style(t *testing.T) {
	checker := qt.New(t)
	f := NewFormatter()
	s1 := NewFormatterStyle("blue", "white")
	s2 := NewFormatterStyle("blue", "white")

	f.SetStyle("test", s1)
	checker.Assert(f.HasStyle("test"), qt.IsTrue)

	f.SetStyle("b", s2)
	expected := "\033[34;47msome \033[39;49m\033[34;47mcustom\033[39;49m\033[34;47m msg\033[39;49m"
	in := "<test>some <b>custom</b> msg</test>"
	checker.Assert(f.Format(in), qt.Equals, expected)
}

func TestFormatter_RedefineStyle(t *testing.T) {
	checker := qt.New(t)
	formatter := NewFormatter()
	s := NewFormatterStyle("blue", "white")

	formatter.SetStyle("info", s)
	checker.Assert(
		formatter.Format("<info>some custom msg</info>"),
		qt.Equals,
		"\033[34;47msome custom msg\033[39;49m",
	)
}

func TestFormatter_InlineStyle(t *testing.T) {
	checker := qt.New(t)
	f := NewFormatter()

	checker.Assert(
		f.Format("<fg=blue;bg=red>some text</>"),
		qt.Equals,
		"\033[34;41msome text\033[39;49m",
	)
	checker.Assert(
		f.Format("<fg=blue;bg=red>some text</fg=blue;bg=red>"),
		qt.Equals,
		"\033[34;41msome text\033[39;49m",
	)
}

type TestFormatter struct {
	formatter *Formatter
}

func (tf TestFormatter) TestCreateStyleFromString(string string) *Style {
	return tf.formatter.createStyleFromString(string)
}

func TestFormatter_InlineStyleWithOptions(t *testing.T) {
	colorTerm := os.Getenv("COLORTERM")
	os.Setenv("COLORTERM", "truecolor")
	type cs struct {
		Tag       string
		Expected  string
		Input     string
		TrueColor bool
	}
	cases := []cs{
		{
			Tag: "<unknown=_unknown_>",
		},
		{
			Tag: "<unknown=_unknown_;a=1;b>",
		},
		{
			Tag:      "<fg=green;>",
			Expected: "\033[32m[test]\033[39m",
			Input:    "[test]",
		},
		{
			Tag:      "<fg=green;bg=blue>",
			Expected: "\033[32;44ma\033[39;49m",
			Input:    "a",
		},
		{
			Tag:      "<fg=green;options=bold>",
			Expected: "\033[32;1mb\033[39;22m",
			Input:    "b",
		},
		{
			Tag:      "<fg=green;options=reverse;>",
			Expected: "\033[32;7m<a>\033[39;27m",
			Input:    "<a>",
		},
		{
			Tag:      "<fg=green;options=bold,underscore>",
			Expected: "\033[32;1;4mz\033[39;22;24m",
			Input:    "z",
		},
		{
			Tag:      "<fg=green;options=bold,underscore,reverse;>",
			Expected: "\033[32;1;4;7md\033[39;22;24;27m",
			Input:    "d",
		},
		{
			Tag:       "<fg=#00ff00;bg=#00f>",
			Expected:  "\033[38;2;0;255;0;48;2;0;0;255m[test]\033[39;49m",
			Input:     "[test]",
			TrueColor: true,
		},
	}

	for i, v := range cases {
		name := fmt.Sprintf("test_%d", i)
		t.Run(name, func(t *testing.T) {
			checker := qt.New(t)
			styleString := v.Tag[1 : len(v.Tag)-1]
			formatter := NewFormatter()
			tf := &TestFormatter{formatter: formatter}
			result := tf.TestCreateStyleFromString(styleString)
			if "" == v.Expected {
				checker.Assert(result, qt.IsNil)
				expected := v.Tag + v.Input + "</" + styleString + ">"
				checker.Assert(formatter.Format(expected), qt.Equals, expected)
			} else {
				checker.Assert(result, qt.IsNotNil)
				checker.Assert(
					formatter.Format(v.Tag+v.Input+"</>"),
					qt.Equals,
					v.Expected,
				)
				checker.Assert(
					formatter.Format(v.Tag+v.Input+"</"+styleString+">"),
					qt.Equals,
					v.Expected,
				)
			}
		})
	}

	os.Setenv("COLORTERM", colorTerm)
}

func TestFormatter_UndecoratedFormatter(t *testing.T) {
	osTermEmulator := os.Getenv("TERMINAL_EMULATOR")

	type cs struct {
		Input               string
		ExpectedUndecorated string
		ExpectedDecorated   string
		TerminalEmulator    string
	}

	cases := []cs{
		{
			Input:               "<error>some error</error>",
			ExpectedUndecorated: "some error",
			ExpectedDecorated:   "\033[37;41msome error\033[39;49m",
		},
		{
			Input:               "<info>some info</info>",
			ExpectedUndecorated: "some info",
			ExpectedDecorated:   "\033[32msome info\033[39m",
		},
		{
			Input:               "<comment>some comment</comment>",
			ExpectedUndecorated: "some comment",
			ExpectedDecorated:   "\033[33msome comment\033[39m",
		},
		{
			Input:               "<question>some question</question>",
			ExpectedUndecorated: "some question",
			ExpectedDecorated:   "\033[30;46msome question\033[39;49m",
		},
		{
			Input:               "<fg=red>some text with inline style</>",
			ExpectedUndecorated: "some text with inline style",
			ExpectedDecorated:   "\033[31msome text with inline style\033[39m",
		},
		{
			Input:               "<href=idea://open/?file=/path/SomeFile.php&line=12>some URL</>",
			ExpectedUndecorated: "some URL",
			ExpectedDecorated:   "\033]8;;idea://open/?file=/path/SomeFile.php&line=12\033\\some URL\033]8;;\033\\",
		},
		{
			Input:               "<href=idea://open/?file=/path/SomeFile.php&line=12>some URL</>",
			ExpectedUndecorated: "some URL",
			ExpectedDecorated:   "some URL",
			TerminalEmulator:    "JetBrains-JediTerm",
		},
	}

	for i, v := range cases {
		if "" == v.TerminalEmulator {
			v.TerminalEmulator = "foo"
		}
		os.Setenv("TERMINAL_EMULATOR", v.TerminalEmulator)

		testName := fmt.Sprintf("case_%d", i)
		t.Run(testName, func(t *testing.T) {
			formatter := NewFormatter()
			checker := qt.New(t)
			checker.Assert(
				formatter.Format(v.Input),
				qt.Equals,
				v.ExpectedDecorated,
			)

			formatter.SetDecorated(false)
			checker.Assert(
				formatter.Format(v.Input),
				qt.Equals,
				v.ExpectedUndecorated,
			)
		})
	}

	os.Setenv("TERMINAL_EMULATOR", osTermEmulator)
}

func TestFormatter_Format(t *testing.T) {
	type cs struct {
		Name     string
		Expected string
		Message  string
	}

	repeatedStrings := strings.Repeat("\\", 14000)
	cases := []cs{
		{
			Name:     "escaping",
			Expected: "foo<bar",
			Message:  "foo\\<bar",
		},
		{
			Name:     "escaping",
			Expected: "foo << bar",
			Message:  "foo << bar",
		},
		{
			Name:     "escaping",
			Expected: "foo << bar \\",
			Message:  "foo << bar \\",
		},
		{
			Name:     "escaping",
			Expected: "foo << \033[32mbar \\ baz\033[39m \\",
			Message:  "foo << <info>bar \\ baz</info> \\",
		},
		{
			Name:     "escaping",
			Expected: "<info>some info</info>",
			Message:  "\\<info>some info\\</info>",
		},
		{
			Name:     "escaping",
			Expected: "\033[33mThis\\Console\\Component does work very well!\033[39m",
			Message:  "<comment>This\\Console\\Component does work very well!</comment>",
		},
		{
			Name:     "bundled.error",
			Expected: "\033[37;41msome error\033[39;49m",
			Message:  "<error>some error</error>",
		},
		{
			Name:     "bundled.info",
			Expected: "\033[32msome info\033[39m",
			Message:  "<info>some info</info>",
		},
		{
			Name:     "bundled.comment",
			Expected: "\033[33msome comment\033[39m",
			Message:  "<comment>some comment</comment>",
		},
		{
			Name:     "bundled.question",
			Expected: "\033[30;46msome question\033[39;49m",
			Message:  "<question>some question</question>",
		},
		{
			Name:     "nested_styles",
			Expected: "\033[37;41msome \033[39;49m\033[32msome info\033[39m\033[37;41m error\033[39;49m",
			Message:  "<error>some <info>some info</info> error</error>",
		},
		{
			Name:     "deep_nested_styles",
			Expected: "\033[37;41merror\033[39;49m\033[32minfo\033[39m\033[33mcomment\033[39m\033[37;41merror\033[39;49m",
			Message:  "<error>error<info>info<comment>comment</info>error</error>",
		},
		{
			Name:     "adjacent_styles",
			Expected: "\033[37;41msome error\033[39;49m\033[32msome info\033[39m",
			Message:  "<error>some error</error><info>some info</info>",
		},
		{
			Name:     "ungreedy_matching",
			Expected: "(\033[32m>=2.0,<2.3\033[39m)",
			Message:  "(<info>>=2.0,<2.3</info>)",
		},
		{
			Name:     "style_escaping",
			Expected: "(\033[32mz>=2.0,<<<a2.3\\\033[39m)",
			Message:  "(<info>" + Escape("z>=2.0,<\\<<a2.3\\") + "</info>)",
		},
		{
			Name:     "style_escaping",
			Expected: "\033[32m<error>some error</error>\033[39m",
			Message:  "<info>" + Escape("<error>some error</error>") + "</info>",
		},
		{
			Name:     "non_style_tag",
			Expected: "\033[32msome \033[39m\033[32m<tag>\033[39m\033[32m \033[39m\033[32m<setting=value>\033[39m\033[32m styled \033[39m\033[32m<p>\033[39m\033[32msingle-char tag\033[39m\033[32m</p>\033[39m",
			Message:  "<info>some <tag> <setting=value> styled <p>single-char tag</p></info>",
		},
		{
			Name:     "format_long_string",
			Expected: "\033[37;41msome error\033[39;49m" + repeatedStrings,
			Message:  "<error>some error</error>" + repeatedStrings,
		},
		{
			Name: "content_with_line_breaks",
			Message: `
<info>
some text</info>
`,
			Expected: "\n\033[32m\nsome text\033[39m\n",
		},
		{
			Name: "content_with_line_breaks",
			Message: `
<info>some text
</info>
`,
			Expected: "\n\033[32msome text\n\033[39m\n",
		},
		{
			Name: "content_with_line_breaks",
			Message: `
<info>
some text
</info>
`,
			Expected: "\n\033[32m\nsome text\n\033[39m\n",
		},
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			checker := qt.New(t)
			f := NewFormatter()

			checker.Assert(f.Format(v.Message), qt.Equals, v.Expected)
		})
	}
}

func TestFormatter_FormatAndWrap(t *testing.T) {
	type cs struct {
		Expected    string
		Message     string
		Width       int
		Undecorated bool
	}

	cases := []cs{
		{
			Expected: "fo\no\x1b[37;41mb\x1b[39;49m\n\x1b[37;41mar\x1b[39;49m\nba\nz",
			Message:  "foo<error>bar</error> baz",
			Width:    2,
		},
		{
			Expected: "pr\ne \x1b[37;41m\x1b[39;49m\n\x1b[37;41mfo\x1b[39;49m\n\x1b[37;41mo \x1b[39;49m\n\x1b[37;41mba\x1b[39;49m\n\x1b[37;41mr \x1b[39;49m\n\x1b[37;41mba\x1b[39;49m\n\x1b[37;41mz\x1b[39;49m \npo\nst",
			Message:  "pre <error>foo bar baz</error> post",
			Width:    2,
		},
		{
			Expected: "pre\x1b[37;41m\x1b[39;49m\n\x1b[37;41mfoo\x1b[39;49m\n\x1b[37;41mbar\x1b[39;49m\n\x1b[37;41mbaz\x1b[39;49m\npos\nt",
			Message:  "pre <error>foo bar baz</error> post",
			Width:    3,
		},
		{
			Expected: "pre \x1b[37;41m\x1b[39;49m\n\x1b[37;41mfoo \x1b[39;49m\n\x1b[37;41mbar \x1b[39;49m\n\x1b[37;41mbaz\x1b[39;49m \npost",
			Message:  "pre <error>foo bar baz</error> post",
			Width:    4,
		},
		{
			Expected: "pre \x1b[37;41mf\x1b[39;49m\n\x1b[37;41moo ba\x1b[39;49m\n\x1b[37;41mr baz\x1b[39;49m\npost",
			Message:  "pre <error>foo bar baz</error> post",
			Width:    5,
		},
		{
			Expected: "Lore\nm \x1b[37;41mip\x1b[39;49m\n\x1b[37;41msum\x1b[39;49m \ndolo\nr \x1b[32msi\x1b[39m\n\x1b[32mt\x1b[39m am\net",
			Message:  "Lorem <error>ipsum</error> dolor <info>sit</info> amet",
			Width:    4,
		},
		{
			Expected: "Lorem \x1b[37;41mip\x1b[39;49m\n\x1b[37;41msum\x1b[39;49m dolo\nr \x1b[32msit\x1b[39m am\net",
			Message:  "Lorem <error>ipsum</error> dolor <info>sit</info> amet",
			Width:    8,
		},
		{
			Expected: "Lorem \x1b[37;41mipsum\x1b[39;49m dolor \x1b[32m\x1b[39m\n\x1b[32msit\x1b[39m, \x1b[37;41mamet\x1b[39;49m et \x1b[32mlauda\x1b[39m\n\x1b[32mntium\x1b[39m architecto",
			Message:  "Lorem <error>ipsum</error> dolor <info>sit</info>, <error>amet</error> et <info>laudantium</info> architecto",
			Width:    18,
		},
		{
			Expected:    "fo\nob\nar\nba\nz",
			Message:     "foo<error>bar</error> baz",
			Width:       2,
			Undecorated: true,
		},
		{
			Expected:    "pr\ne \nfo\no \nba\nr \nba\nz \npo\nst",
			Message:     "pre <error>foo bar baz</error> post",
			Width:       2,
			Undecorated: true,
		},
		{
			Expected:    "pre\nfoo\nbar\nbaz\npos\nt",
			Message:     "pre <error>foo bar baz</error> post",
			Width:       3,
			Undecorated: true,
		},
		{
			Expected:    "pre \nfoo \nbar \nbaz \npost",
			Message:     "pre <error>foo bar baz</error> post",
			Width:       4,
			Undecorated: true,
		},
		{
			Expected:    "pre f\noo ba\nr baz\npost",
			Message:     "pre <error>foo bar baz</error> post",
			Width:       5,
			Undecorated: true,
		},
	}

	for _, v := range cases {
		t.Run("format_and_wrap", func(t *testing.T) {
			formatter := NewFormatter()
			checker := qt.New(t)

			if v.Undecorated {
				formatter.SetDecorated(false)
			}
			checker.Assert(
				formatter.FormatAndWrap(v.Message, v.Width),
				qt.Equals,
				v.Expected,
			)
		})
	}
}
