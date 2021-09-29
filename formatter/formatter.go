package formatter

import (
	"errors"
	"fmt"

	"regexp"
	"strings"
)

//Formatter class for output
type Formatter struct {
	styles    map[string]*Style
	decorated bool
}

//NewFormatter create and returns new Formatter
func NewFormatter() *Formatter {
	formatter := &Formatter{
		decorated: true,
		styles:    make(map[string]*Style),
	}

	formatter.AddStyle("error", NewFormatterStyle("white", "red"))
	formatter.AddStyle("info", NewFormatterStyle("green", ""))
	formatter.AddStyle("comment", NewFormatterStyle("yellow", ""))
	formatter.AddStyle("question", NewFormatterStyle("black", "cyan"))
	return formatter
}

//SetDecorated sets the decorated flag
func (f *Formatter) SetDecorated(decorated bool) {
	f.decorated = decorated
}

//IsDecorated returns true if the output will decorate messages
func (f *Formatter) IsDecorated() bool {
	return f.decorated
}

//AddStyle add a new Style into this Formatter
func (f *Formatter) AddStyle(name string, style *Style) {
	f.styles[name] = style
}

//HasStyle returns true if this formatter has given style name
func (f *Formatter) HasStyle(name string) bool {
	if _, ok := f.styles[name]; ok {
		return true
	}
	return false
}

//GetStyle returns Style with given name
func (f *Formatter) GetStyle(name string) (s *Style, err error) {
	if s, ok := f.styles[name]; ok {
		return s, nil
	}
	errMsg := fmt.Sprintf(`This formatter doesn't have style with name: "%s"`, name)
	return nil, errors.New(errMsg)
}

func (f *Formatter) Format(message string) string {
	return f.FormatAndWrap(message, 0)
}

func (f *Formatter) FormatAndWrap(message string, width int) string {
	tagRegex := `[a-z][a-z0-9,_=;-]*`
	sRegex := fmt.Sprintf("<((%s)|/(%s)?)>", tagRegex, tagRegex)
	regex := regexp.MustCompile(sRegex)
	tags := regex.FindAllString(message, -1)
	indexes := regex.FindAllStringIndex(message, -1)

	offset := 0
	currentLineLength := 0
	output := ""

	for i := 0; i < len(tags); i++ {
		pos := indexes[i][0]
		text := tags[i]

		if 0 != pos && '\\' == message[pos-1] {
			continue
		}

		fmt.Println(text)
		//substrMessage := message[offset:(pos - offset)]
		//output += applyCurrentStyle(substrMessage, output, width, currentLineLength)
	}
	substrMessage := message[offset:]
	output += f.applyCurrentStyle(substrMessage, output, width, currentLineLength)
	if strings.Contains(output, "\000") {
		output = strings.ReplaceAll(output, "\000", "\\")
		output = strings.ReplaceAll(output, "\\<", "<")
	}
	return strings.ReplaceAll(output, "\\<", "<")
}

//Escape will escapes "<" special char in given text
func Escape(text string) string {
	regex := regexp.MustCompile("/([^\\\\\\\\]?)</")
	replace := "$1\\<"
	replaced := regex.ReplaceAllString(text, replace)

	return EscapeTrailingBackslash(replaced)
}

//EscapeTrailingBackslash escapes trailing "\" in given text
func EscapeTrailingBackslash(text string) string {
	last := text[len(text)-1:]

	if "\\" == last {
		strlen := len(text)
		text = strings.TrimRight(text, "\\")
		text = strings.ReplaceAll(text, "\000", "")
		text += strings.Repeat("\000", strlen-len(text))
	}

	return text
}

func (f *Formatter) applyCurrentStyle(text string, current string, width int, currentLineLength int) string {
	if "" == text {
		return ""
	}

	if 0 == width {
		if f.decorated {
			//TODO: apply current style stack
			f.styles["info"].Apply(text)
		}
		return text
	}

	if 0 == currentLineLength && "" != current {
		text = strings.TrimLeft(text, "")
	}

	prefix := ""
	if currentLineLength > 0 {
		i := width - currentLineLength
		prefix = text[0:i] + "\n"
		text = text[i:]
	}

	// preg_match('~(\\n)$~', $text, $matches);
	regex := regexp.MustCompile(`(\\n)$`)
	matches := regex.FindAllString(text, -1)

	return "<info>some info</info>" + strings.Join(matches, " ") + prefix
}
