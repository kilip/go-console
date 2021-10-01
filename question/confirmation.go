package question

import (
	"regexp"
	"strconv"
)

// ConfirmationQuestion Represents a yes/no question.
// * true: when answer with yes or 1
// * false: when answer is no or 0
// Answering this question with other value like foo, bar will result false
type ConfirmationQuestion struct {
	trueAnswerRegex string
	*Question
}

// NewConfirmationQuestion creates new ConfirmationQuestion object
func NewConfirmationQuestion(question string, defaultValue bool) *ConfirmationQuestion {
	q := &ConfirmationQuestion{
		trueAnswerRegex: "^y|^Y",
		Question:        NewQuestion(question),
	}
	q.SetDefault(defaultValue)
	q.SetNormalizer(q.getDefaultNormalizer())

	return q
}

// getDefaultNormalizer will sets default normalization for ConfirmationQuestion
func (cq *ConfirmationQuestion) getDefaultNormalizer() func(input string) interface{} {
	return func(answer string) interface{} {
		defaultValue := cq.defaultValue
		if val, err := strconv.ParseBool(answer); err == nil {
			return val
		}
		regex := regexp.MustCompile(cq.trueAnswerRegex)
		answerIsTrue := regex.MatchString(answer)
		if false == defaultValue {
			return "" != answer && answerIsTrue
		}

		return "" == answer || answerIsTrue
	}
}
