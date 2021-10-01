package question

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestConfirmationQuestion_RegexUseCases(t *testing.T) {
	type cs struct {
		Name         string
		DefaultValue bool
		Answers      []string
		Expected     bool
	}

	cases := []cs{
		{
			DefaultValue: true,
			Answers:      []string{"y", "Y", "yes", "YES", "yEs", ""},
			Expected:     true,
			Name:         "when default is true, the normalizer must return true for %s",
		},
		{
			DefaultValue: true,
			Answers:      []string{"n", "N", "no", "NO", "nO", "foo", "0"},
			Expected:     false,
			Name:         "When default is true, the normalizer must return false for %s",
		},
		{
			DefaultValue: false,
			Answers:      []string{"y", "Y", "yes", "YES", "yEs"},
			Expected:     true,
			Name:         "When default is false, the normalizer must return true for %s",
		},
		{
			DefaultValue: false,
			Answers:      []string{"n", "N", "no", "NO", "nO", "foo", "0", ""},
			Expected:     false,
			Name:         "When default is false, the normalizer must return false for %s",
		},
	}

	for _, testCase := range cases {

		t.Run(testCase.Name, func(t *testing.T) {
			c := qt.New(t)
			q := NewConfirmationQuestion("A question", testCase.DefaultValue)
			normalizer := q.GetNormalizer()
			for _, answer := range testCase.Answers {
				c.Assert(normalizer(answer), qt.Equals, testCase.Expected)
			}
		})
	}
}
