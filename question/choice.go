package question

import (
	"errors"
	"expvar"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ChoiceQuestion struct {
	choices      interface{}
	multiselect  bool
	prompt       string
	errorMessage string
	*Question
}

func NewChoiceQuestion(question string, choices interface{}) *ChoiceQuestion {
	q := &ChoiceQuestion{
		Question: NewQuestion(question),
	}

	q.SetChoices(choices)
	q.SetMultiSelect(false)
	q.SetPrompt(" > ")
	q.SetErrorMessage(`Value "%s" is invalid`)
	q.SetValidator(q.getDefaultValidator())
	q.SetTrimmable(true)

	return q
}

func (cq *ChoiceQuestion) SetChoices(choices interface{}) {
	cq.choices = choices
}

func (cq *ChoiceQuestion) GetChoices() interface{} {
	return cq.choices
}

func (cq *ChoiceQuestion) SetMultiSelect(multiselect bool) {
	cq.multiselect = multiselect
}

func (cq *ChoiceQuestion) GetMultiSelect() bool {
	return cq.multiselect
}

func (cq *ChoiceQuestion) SetPrompt(prompt string) {
	cq.prompt = prompt
}

func (cq *ChoiceQuestion) GetPrompt() string {
	return cq.prompt
}

func (cq *ChoiceQuestion) SetErrorMessage(message string) {
	cq.errorMessage = message
}

func (cq *ChoiceQuestion) GetErrorMessage() string {
	return cq.errorMessage
}

func (cq *ChoiceQuestion) getDefaultValidator() func(input string) (v interface{}, err error) {
	return func(selected string) (v interface{}, err error) {
		var sChoices []string

		if cq.multiselect {
			// check for separated comma values
			regex := regexp.MustCompile(`^[^,]+(?:,[^,]+)*$`)
			if false == regex.MatchString(selected) {
				return nil, errors.New(fmt.Sprintf(cq.errorMessage, selected))
			}
			sChoices = strings.Split(selected, ",")
		} else {
			sChoices = []string{selected}
		}

		if cq.IsTrimmable() {
			var trimmed []string

			for _, v := range sChoices {
				trimmed = append(trimmed, strings.Trim(v, " "))
			}
			sChoices = trimmed
		}

		var multiSelectChoices []string
		var currentChoice string
		for _, choiceValue := range sChoices {
			var results []string
			var result string
			cType := reflect.TypeOf(cq.choices)
			if reflect.Map == cType.Kind() {
				for key, val := range cq.choices.(map[string]interface{}) {
					converted := val

					if _, ok := val.(expvar.Var); ok {
						converted = val.(expvar.Var).String()
					}

					if key == choiceValue {
						results = append(results, key)
						result = choiceValue
					} else if choiceValue == converted {
						results = append(results, key)
						result = key
					}
				}
			} else if reflect.Slice == cType.Kind() {
				choicesType := reflect.TypeOf(cq.choices).Elem().String()
				if "string" == choicesType {
					for key, choice := range cq.choices.([]string) {
						strChoice := fmt.Sprintf("%v", choice)

						if strChoice == choiceValue {
							results = append(results, strChoice)
							result = strChoice
						} else if kToStr := strconv.Itoa(key); kToStr == choiceValue {
							results = append(results, strChoice)
							result = strChoice
						}
					}
				} else if "interface {}" == choicesType {
					for key, choice := range cq.choices.([]interface{}) {
						strChoice := fmt.Sprintf("%v", choice)

						if strChoice == choiceValue {
							results = append(results, strChoice)
							result = strChoice
						} else if kToStr := strconv.Itoa(key); kToStr == choiceValue {
							results = append(results, strChoice)
							result = strChoice
						}
					}
				}
			}

			if len(results) > 1 {
				errMsg := fmt.Sprintf(
					`The provided answer is ambigous. Value should be one of "%s"`,
					strings.Join(results, `" or "`),
				)
				return nil, errors.New(errMsg)
			}

			if "" == result {
				return nil, errors.New(fmt.Sprintf(cq.errorMessage, choiceValue))
			}

			multiSelectChoices = append(multiSelectChoices, result)
			currentChoice = result
		}

		if cq.multiselect {
			return multiSelectChoices, nil
		}
		return currentChoice, nil
	}
}
