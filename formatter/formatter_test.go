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
	s := NewFormatterStyleO("black", "green", []string{})

	f.AddStyle("info", s)
	checker.Assert(f.HasStyle("info"), qt.IsTrue)
	ostyle, _ := f.GetStyle("info")
	checker.Assert(ostyle, qt.Equals, s)
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
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			checker := qt.New(t)
			f := NewFormatter()

			checker.Assert(f.Format(v.Message), qt.Equals, v.Expected)

		})
	}
}

func TestFoo(t *testing.T) {
	checker := qt.New(t)
	formatter := NewFormatter()

	output := formatter.Format("\\<info>some info\\</info>")
	checker.Assert(output, qt.Equals, "<info>some info</info>")
}
