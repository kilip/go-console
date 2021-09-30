package question

import (
	"errors"
	"reflect"
)

// Question represents a Question
type Question struct {
	question              string
	attempts              int
	hidden                bool
	hiddenFallback        bool
	autoCompleterCallback func(input string) []string
	validator             func(input string) (isValid bool, err error)
	defaultValue          interface{}
	trimmable             bool
	multiline             bool
	normalizer            func(input string) interface{}
}

// NewQuestion creates new Question object
func NewQuestion(q string, defaultValue interface{}) *Question {
	return &Question{
		question:       q,
		defaultValue:   defaultValue,
		hidden:         false,
		hiddenFallback: true,
		trimmable:      false,
		multiline:      false,
	}
}

// GetQuestion returns the question
func (q *Question) GetQuestion() string {
	return q.question
}

// GetDefault returns the default answer value
func (q *Question) GetDefault() interface{} {
	return q.defaultValue
}

// IsMultiline returns whether the user response accepts newline characters
func (q *Question) IsMultiline() bool {
	return q.multiline
}

// SetMultiline set whether the user response should accept newline characters
func (q *Question) SetMultiline(multiline bool) {
	q.multiline = multiline
}

// IsHidden returns whether the user response must be hidden
func (q *Question) IsHidden() bool {
	return q.hidden
}

// SetHidden sets whether the user response must be hidden or not
func (q *Question) SetHidden(hidden bool) error {
	if nil != q.autoCompleterCallback {
		return errors.New("A hidden question cannot use the autocompleter")
	}
	q.hidden = hidden
	return nil
}

// IsHiddenFallback In case the response can not be hidden,
// whether to fallback on non-hidden question or not.
func (q *Question) IsHiddenFallback() bool {
	return q.hiddenFallback
}

// SetHiddenFallback Sets whether to fallback on non-hidden question
// if the response can not be hidden.
func (q *Question) SetHiddenFallback(fallback bool) {
	q.hiddenFallback = fallback
}

// GetAutoCompleterValues gets values for the autocompleter
func (q *Question) GetAutoCompleterValues() []string {
	callback := q.autoCompleterCallback

	return callback("")
}

// SetAutoCompleterValues sets values for the autocompleter
func (q *Question) SetAutoCompleterValues(v interface{}) error {
	var cb func(input string) []string
	var err error
	rv := reflect.ValueOf(v)

	if reflect.Slice == rv.Kind() {
		cb = func(input string) []string {
			return v.([]string)
		}

	} else if reflect.Func == rv.Kind() {
		cb = v.(func(input string) []string)
	}

	if nil != cb {
		q.autoCompleterCallback = cb
	} else {
		err = errors.New("can't set auto completer values with given values")
	}
	return err
}

// GetAutoCompleterCallback gets the callback function used for the autocompleter
func (q *Question) GetAutoCompleterCallback() func(input string) []string {
	return q.autoCompleterCallback
}

// SetAutoCompleterCallback sets the callback function used for the autocompleter.
// The callback is passed the user input as argument and should return an iterable of corresponding suggestions.
func (q *Question) SetAutoCompleterCallback(callback func(input string) []string) {
	q.autoCompleterCallback = callback
}

// SetValidator sets a validator for this question
func (q *Question) SetValidator(validator func(input string) (valid bool, error error)) {
	q.validator = validator
}

// GetValidator returns the validator for this question
func (q *Question) GetValidator() func(input string) (valid bool, error error) {
	return q.validator
}

// SetMaxAttempts sets the maximum number of attempts for answering this question
// zero means an unlimited number of attempts
func (q *Question) SetMaxAttempts(attempts int) {
	q.attempts = attempts
}

// GetMaxAttempts gets the maximum number of attempts for answering this question
// zero means an unlimited number of attempts
func (q *Question) GetMaxAttempts() int {
	return q.attempts
}

// SetNormalizer sets a normalizer for the response
func (q *Question) SetNormalizer(normalizer func(input string) interface{}) {
	q.normalizer = normalizer
}

// GetNormalizer gets the normalizer for the response
func (q *Question) GetNormalizer() func(input string) interface{} {
	return q.normalizer
}

// IsTrimmable check if the question can be trimmed
func (q *Question) IsTrimmable() bool {
	return q.trimmable
}

// SetTrimmable sets if the question can be trimmed
func (q *Question) SetTrimmable(trimmable bool) {
	q.trimmable = trimmable
}
