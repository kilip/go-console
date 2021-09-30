package formatter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//Formatter class for output
type Formatter struct {
	styles     map[string]*Style
	decorated  bool
	styleStack *StyleStack
}

//NewFormatter create and returns new Formatter
func NewFormatter() *Formatter {
	formatter := &Formatter{
		decorated:  true,
		styles:     make(map[string]*Style),
		styleStack: NewStyleStack(),
	}
	formatter.SetStyle("error", NewFormatterStyle("white", "red"))
	formatter.SetStyle("info", NewFormatterStyle("green", ""))
	formatter.SetStyle("comment", NewFormatterStyle("yellow", ""))
	formatter.SetStyle("question", NewFormatterStyle("black", "cyan"))
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

//SetStyle add a new Style into this Formatter
func (f *Formatter) SetStyle(name string, style *Style) {
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

//Format returns formatted message according to the given styles.
func (f *Formatter) Format(message string) string {
	return f.FormatAndWrap(message, 0)
}

//FormatAndWrap returns formatted message according to the given styles.
//This method will also wrap the formatted message according to the given width.
//Passing width to 0 will disable wrap method
func (f *Formatter) FormatAndWrap(message string, width int) string {
	//tagRegex := `[a-z][a-z0-9,_=;-]*`
	//sRegex := fmt.Sprintf("<((%s)|/(%s)?)>", tagRegex, tagRegex)
	tagRegex := "[a-z][^<>]*+"
	tagRegex = fmt.Sprintf("<((%s)|/(%s)?)>", tagRegex, tagRegex)
	regex := regexp.MustCompilePOSIX(tagRegex)
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

		mSubStr := message[offset:pos]
		cLength, oCstyle := f.applyCurrentStyle(mSubStr, output, width, currentLineLength)
		output += oCstyle
		currentLineLength = cLength
		offset = pos + len(text)

		// opening tag?
		tag := tags[i]
		tag = tag[1 : len(tag)-1]
		open := false
		ix := strings.Index(tag, "/")
		if 0 == ix {
			tag = tag[1:]
		} else {
			open = true
		}

		style := f.createStyleFromString(tag)
		if !open && tag == "" {
			f.styleStack.Pop()
		} else if nil == style {
			cLength, so := f.applyCurrentStyle(text, output, width, currentLineLength)
			currentLineLength = cLength
			output += so
		} else if open {
			f.styleStack.Push(style)
		} else {
			f.styleStack.PopS(style)
		}
	}
	substrMessage := message[offset:]
	_, cStyle := f.applyCurrentStyle(substrMessage, output, width, currentLineLength)
	output += cStyle
	if strings.Contains(output, "\000") {
		output = strings.ReplaceAll(output, "\000", "\\")
		output = strings.ReplaceAll(output, "\\<", "<")
	}
	return strings.ReplaceAll(output, "\\<", "<")
}

//Escape will escapes "<" special char in given text
func Escape(text string) string {
	regex := regexp.MustCompile("([^\\\\\\\\]?)<")
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

//applyCurrentStyle applies current style from stack to text, if must be applied.
func (f *Formatter) applyCurrentStyle(text string, current string, width int, currentLineLength int) (c int, output string) {
	if "" == text {
		return currentLineLength, ""
	}

	if 0 == width {
		if f.decorated {
			return currentLineLength, f.styleStack.Current().Apply(text)
		}
		return currentLineLength, text
	}

	if 0 == currentLineLength && "" != current {
		text = strings.TrimLeft(text, " ")
		text = strings.TrimLeft(text, "\n")
	}

	prefix := ""
	if currentLineLength > 0 {
		i := width - currentLineLength
		prefixLen := i
		if len(text) < prefixLen {
			prefixLen = len(text)
		}
		prefix = text[0:prefixLen] + "\n"

		if len(text) < i {
			text = ""
		} else {
			text = text[i:]
		}

	}

	// preg_match('~(\\n)$~', $text, $matches);
	regex := regexp.MustCompile(`(\\n)$`)
	matches := regex.FindAllString(text, -1)

	rReplace := fmt.Sprintf("([^\\n]{%d})\\ *", width)
	regReplace := regexp.MustCompilePOSIX(rReplace)
	text = prefix + regReplace.ReplaceAllString(text, "$1\n")
	text = strings.TrimRight(text, "\n")
	if len(matches) > 0 {
		text += matches[1]
	}

	if 0 == currentLineLength && "" != current && current != "\n" {
		text = "\n" + text
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		currentLineLength += len(line)
		if width <= currentLineLength {
			currentLineLength = 0
		}
	}

	if f.IsDecorated() {
		for i, line := range lines {
			lines[i] = f.styleStack.Current().Apply(line)
		}
	}
	return currentLineLength, strings.Join(lines, "\n")
}

func (f *Formatter) createStyleFromString(string string) *Style {
	if _, ok := f.styles[string]; ok {
		return f.styles[string]
	}

	regex := regexp.MustCompile(`([^=]+)=([^;]+)(;|$)`)
	matches := regex.FindAllStringSubmatch(string, -1)

	if 0 == len(matches) {
		return nil
	}

	style := NewFormatterStyle("", "")

	for _, match := range matches {
		match[1] = strings.ToLower(match[1])

		if "fg" == match[1] {
			style.SetForeground(match[2])
		} else if "bg" == match[1] {
			style.SetBackground(match[2])
		} else if "href" == match[1] {
			style.SetHref(match[2])
		} else if "options" == match[1] {
			regex := regexp.MustCompile("([^,;]+)")
			lower := strings.ToLower(match[2])
			options := regex.FindAllString(lower, -1)

			for _, option := range options {
				style.SetOption(option)
			}
		} else {
			return nil
		}
	}

	return style
}
