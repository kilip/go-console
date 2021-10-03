package question

import (
	qt "github.com/frankban/quicktest"
	"reflect"
	"strings"
	"testing"
)

func getChoiceQuestion() *ChoiceQuestion {
	cq := NewChoiceQuestion("Test choice question", []string{"foo", "bar"})
	return cq
}

func TestNewChoiceQuestion(t *testing.T) {
	c := qt.New(t)
	q := getChoiceQuestion()

	c.Assert(q.GetQuestion(), qt.Equals, "Test choice question")
	c.Assert(q.GetChoices(), qt.HasLen, 2)
	c.Assert(q.GetMultiSelect(), qt.IsFalse)
	c.Assert(q.GetPrompt(), qt.Equals, " > ")
	c.Assert(q.GetErrorMessage(), qt.Equals, `Value "%s" is invalid`)
}

func TestChoiceQuestion_SelectUseCases(t *testing.T) {
	type cs struct {
		MultiSelect bool
		Answers     []string
		Expected    interface{}
		Message     string
	}

	cases := []cs{
		{
			Answers:  []string{"First response", "First response", " First response", " First response "},
			Expected: "First response",
			Message:  "When passed single answer on singleSelect, the defaultValidator must return this answer as a string",
		},
		{
			MultiSelect: true,
			Answers:     []string{"First response", "First response ", " First response", " First response "},
			Expected:    []string{"First response"},
			Message:     "When passed single answer on MultiSelect, the defaultValidator must return this answer as an array",
		},
		{
			MultiSelect: true,
			Answers:     []string{"First response,Second response", "First response , Second response "},
			Expected:    []string{"First response", "Second response"},
			Message:     "When passed multiple answers on MultiSelect, the defaultValidator must return these answers as an array",
		},
		{
			Answers:  []string{"0"},
			Expected: "First response",
			Message:  "When passed single answer using choice's key, the defaultValidator must return the choice value",
		},
		{
			MultiSelect: true,
			Answers:     []string{"0, 2"},
			Expected:    []string{"First response", "Third response"},
			Message:     "When passed multiple answers using choices' key, the defaultValidator must return the choice values in an array",
		},
	}

	for _, testCase := range cases {
		t.Run("selectUseCase", func(t *testing.T) {
			c := qt.New(t)
			q := NewChoiceQuestion("A question", []string{
				"First response",
				"Second response",
				"Third response",
				"Fourth response",
			})
			q.SetMultiSelect(testCase.MultiSelect)

			for _, answer := range testCase.Answers {
				val := q.GetValidator()
				actual, err := val(answer)
				c.Assert(err, qt.IsNil)

				r := reflect.TypeOf(actual)
				if reflect.String == r.Kind() {
					c.Assert(actual, qt.Equals, testCase.Expected)
				} else {
					actJoin := strings.Join(actual.([]string), ",")
					expJoin := strings.Join(testCase.Expected.([]string), ",")
					c.Assert(actJoin, qt.Equals, expJoin)
				}
			}
		})
	}
}

func TestChoiceQuestion_NonTrimmable(t *testing.T) {
	c := qt.New(t)
	q := NewChoiceQuestion("A question", []string{
		"First response ",
		" Second response",
		"  Third response  ",
	})
	val := q.GetValidator()

	q.SetTrimmable(false)

	// let's begin the tests
	out, err := val("  Third response  ")
	c.Assert(err, qt.IsNil)
	c.Assert(out, qt.Equals, "  Third response  ")

	q.SetMultiSelect(true)
	out, err = val("First response , Second response")
	c.Assert(err, qt.IsNil)
	c.Assert(strings.Join(out.([]string), ","), qt.Equals, "First response , Second response")
}

type StringChoice struct {
	string string
}

func (sc *StringChoice) String() string {
	return sc.string
}

func TestChoiceQuestion_SelectAssociativeChoices(t *testing.T) {
	type cs struct {
		Name     string
		Answer   string
		Expected string
	}
	cases := []cs{
		{
			Name:     "select '0' choice by key",
			Answer:   "0",
			Expected: "0",
		},
		{
			Name:     "select '0' choice by value",
			Answer:   "First choice",
			Expected: "0",
		},
		{
			Name:     "select by key",
			Answer:   "foo",
			Expected: "foo",
		},
		{
			Name:     "select by value",
			Answer:   "Foo",
			Expected: "foo",
		},
		{
			Name:     "select by key, with numeric key",
			Answer:   "99",
			Expected: "99",
		},
		{
			Name:     "select by value, with numeric key",
			Answer:   "N°99",
			Expected: "99",
		},
		{
			Name:     "select by key, with string object value",
			Answer:   "string object",
			Expected: "string object",
		},
		{
			Name:     "select by value, with string object value",
			Answer:   "String Object",
			Expected: "string object",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			c := qt.New(t)
			stringChoice := &StringChoice{string: "String Object"}
			q := NewChoiceQuestion("A question", map[string]interface{}{
				"0":             "First choice",
				"foo":           "Foo",
				"99":            "N°99",
				"string object": stringChoice,
			})

			validator := q.GetValidator()
			out, err := validator(testCase.Answer)
			c.Assert(err, qt.IsNil)

			c.Assert(out, qt.Equals, testCase.Expected)
		})
	}
}

func TestChoiceQuestion_SelectWithNonStringChoices(t *testing.T) {
	c := qt.New(t)
	q := NewChoiceQuestion("A question", []interface{}{
		&StringChoice{string: "foo"},
		&StringChoice{string: "bar"},
		&StringChoice{string: "baz"},
	})
	validator := q.GetValidator()

	// begin testing

	// answer can be selected by its string value
	out, err := validator("foo")
	c.Assert(err, qt.IsNil)
	c.Assert(out, qt.Equals, "foo")

	// answer can be selected by index
	out, err = validator("0")
	c.Assert(err, qt.IsNil)
	c.Assert(out, qt.Equals, "foo")

	// test multi select
	q.SetMultiSelect(true)
	out, err = validator("baz, bar")
	c.Assert(strings.Join(out.([]string), ","), qt.Equals, "baz,bar")
}
