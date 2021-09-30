package formatter

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestNewStyleStack(t *testing.T) {
	checker := qt.New(t)
	stack := NewStyleStack()

	checker.Assert(stack.GetEmptyStyle, qt.IsNotNil)
}

func TestStyleStack_EmptyStyle(t *testing.T) {
	checker := qt.New(t)
	stack := NewStyleStack()
	style := NewFormatterStyle("", "")

	checker.Assert(stack.GetEmptyStyle(), qt.IsNotNil)
	stack.SetEmptyStyle(style)
	checker.Assert(stack.GetEmptyStyle(), qt.Equals, style)
}

func TestStyleStack_Push(t *testing.T) {
	checker := qt.New(t)
	stack := NewStyleStack()

	s1 := NewFormatterStyle("white", "black")
	s2 := NewFormatterStyle("yellow", "blue")
	stack.Push(s1)
	stack.Push(s2)

	checker.Assert(stack.Current(), qt.Equals, s2)

	s3 := NewFormatterStyle("green", "red")
	stack.Push(s3)
	checker.Assert(stack.Current(), qt.Equals, s3)
}

func TestStyleStack_Pop(t *testing.T) {
	c := qt.New(t)
	stack := NewStyleStack()
	s1 := NewFormatterStyle("white", "black")
	s2 := NewFormatterStyle("yellow", "blue")

	stack.Push(s1)
	stack.Push(s2)

	var o *Style
	o, _ = stack.Pop()
	c.Assert(o, qt.Equals, s2)
	o, _ = stack.Pop()
	c.Assert(o, qt.Equals, s1)
}

func TestStyleStack_PopEmpty(t *testing.T) {
	c := qt.New(t)
	stack := NewStyleStack()
	var o *Style
	o, _ = stack.Pop()
	c.Assert(o, qt.Equals, stack.GetEmptyStyle())
}

func TestStyleStack_PopNotLast(t *testing.T) {
	var o *Style
	c := qt.New(t)
	stack := NewStyleStack()
	s1 := NewFormatterStyle("white", "black")
	s2 := NewFormatterStyle("yellow", "blue")
	s3 := NewFormatterStyle("green", "red")

	stack.Push(s1)
	stack.Push(s2)
	stack.Push(s3)

	o, _ = stack.PopS(s2)
	c.Assert(o, qt.Equals, s2)
	o, _ = stack.Pop()
	c.Assert(o, qt.Equals, s1)
}

func TestStyleStack_InvalidPop(t *testing.T) {
	c := qt.New(t)
	stack := NewStyleStack()
	s1 := NewFormatterStyle("white", "black")
	s2 := NewFormatterStyle("yellow", "blue")

	stack.Push(s1)
	o, err := stack.PopS(s2)
	c.Assert(err, qt.IsNotNil)
	c.Assert(o, qt.IsNil)
	c.Assert(err.Error(), qt.Contains, "incorrectly nested style tag found")
}
