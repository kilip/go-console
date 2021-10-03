package formatter

import (
	"github.com/kilip/go-console"
	"os"
	"strconv"
	"strings"
)

type Style struct {
	color                 *console.Color
	foreground            string
	background            string
	options               []string
	href                  string
	handlesHrefGracefully bool
}

//NewFormatterStyle creates new Style object
func NewFormatterStyle(foreground string, background string) *Style {
	return NewFormatterStyleO(foreground, background, []string{})
}

//NewFormatterStyleO creates new Style object and set options
func NewFormatterStyleO(foreground string, background string, options []string) *Style {
	clr := console.NewColorWithO(foreground, background, options)

	style := &Style{
		color:                 clr,
		options:               options,
		foreground:            foreground,
		background:            background,
		handlesHrefGracefully: defaultHandlesHrefGracefully(),
	}
	return style
}

func (s *Style) SetForeground(color string) {
	s.foreground = color
	s.color = console.NewColorWithO(s.foreground, s.background, s.options)
}

func (s *Style) SetBackground(color string) {
	s.background = color
	s.color = console.NewColorWithO(s.foreground, s.background, s.options)
}

func (s *Style) Apply(text string) string {
	if "" != s.href && s.handlesHrefGracefully {
		text = "\033]8;;" + s.href + "\033\\" + text + "\033]8;;\033\\"
	}

	return s.color.Apply(text)
}

func (s *Style) SetHref(href string) {
	s.href = href
}

func (s *Style) SetOption(option string) {
	if false == strings.Contains(strings.Join(s.options, " "), option) {
		s.options = append(s.options, option)
		s.color = console.NewColorWithO(s.foreground, s.background, s.options)
	}
}

func (s *Style) UnsetOption(option string) {
	var options []string

	for _, val := range s.options {
		if val != option {
			options = append(options, val)
		}
	}

	s.options = options
	s.color = console.NewColorWithO(s.foreground, s.background, s.options)
}

func (s *Style) SetOptions(options []string) {
	s.options = options
	s.color = console.NewColorWithO(s.foreground, s.background, s.options)
}

func defaultHandlesHrefGracefully() bool {
	emulator := os.Getenv("TERMINAL_EMULATOR")
	konsole := os.Getenv("KONSOLE_VERSION")
	inKonsole, _ := strconv.Atoi(konsole)

	if "JetBrains-JediTerm" != emulator && "" == konsole || inKonsole > 201100 {
		return true
	}
	return false
}
