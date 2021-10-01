package question

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func createDefQuestion() *Question {
	q := NewQuestion("Test question")
	q.SetDefault("default")

	return q
}

func TestQuestion_GetQuestion(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.GetQuestion(), qt.Equals, "Test question")
}

func TestQuestion_GetDefault(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.GetDefault(), qt.Equals, "default")
}

func TestQuestion_Multiline(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.IsMultiline(), qt.IsFalse)
	q.SetMultiline(true)
	c.Assert(q.IsMultiline(), qt.IsTrue)
}

func TestQuestion_Hidden(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.IsHidden(), qt.IsFalse)
	q.SetHidden(true)
	c.Assert(q.IsHidden(), qt.IsTrue)

	q.SetAutoCompleterCallback(func(input string) []string {
		return nil
	})
	err := q.SetHidden(true)
	c.Assert(err, qt.IsNotNil)
}

func TestQuestion_HiddenFallback(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.IsHiddenFallback(), qt.IsTrue)
	q.SetHiddenFallback(false)
	c.Assert(q.IsHiddenFallback(), qt.IsFalse)
}

func TestQuestion_AutoCompleterValues(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	vals := []string{"hello", "world"}
	err := q.SetAutoCompleterValues(vals)
	out := q.GetAutoCompleterValues()
	c.Assert(err, qt.IsNil)
	c.Assert(len(out), qt.Equals, len(vals))
	c.Assert(q.GetAutoCompleterCallback(), qt.IsNotNil)

	q.SetAutoCompleterCallback(func(input string) []string {
		return vals
	})
	c.Assert(q.GetAutoCompleterValues(), qt.IsNotNil)
}

func TestQuestion_Validator(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	validator := func(input string) (valid interface{}, err error) {
		return true, nil
	}

	q.SetValidator(validator)
	valid, err := q.GetValidator()("")
	c.Assert(err, qt.IsNil)
	c.Assert(valid, qt.IsTrue)
}

func TestQuestion_SetMaxAttempts(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.GetMaxAttempts(), qt.Equals, 0)
	q.SetMaxAttempts(1)
	c.Assert(q.GetMaxAttempts(), qt.Equals, 1)
}

func TestQuestion_Normalizer(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	normalizer := func(input string) interface{} {
		return input
	}
	c.Assert(q.GetNormalizer(), qt.IsNil)
	q.SetNormalizer(normalizer)
	c.Assert(q.GetNormalizer()("input"), qt.Equals, "input")
}

func TestQuestion_IsTrimmable(t *testing.T) {
	c := qt.New(t)
	q := createDefQuestion()

	c.Assert(q.IsTrimmable(), qt.IsFalse)
	q.SetTrimmable(true)
	c.Assert(q.IsTrimmable(), qt.IsTrue)
}
