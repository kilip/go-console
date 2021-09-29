package formatter

import (
	qt "github.com/frankban/quicktest"
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

func TestFormatter_Format(t *testing.T) {
	type cs struct {
		Name     string
		Expected string
		Message  string
	}

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
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			checker := qt.New(t)
			f := NewFormatter()

			checker.Assert(f.Format(v.Message), qt.Equals, v.Expected)
		})
	}
}
