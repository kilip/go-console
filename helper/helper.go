package helper

import (
	"github.com/kilip/go-console/formatter"
	"regexp"
)

func RemoveDecoration(formatter *formatter.Formatter, text string) string {
	isDecorated := formatter.IsDecorated()
	formatter.SetDecorated(false)

	// remove <...> formatting
	formatted := formatter.Format(text)
	// remove already formatted characters
	regex := regexp.MustCompilePOSIX(`\033[[^m]*m`)
	formatted = regex.ReplaceAllString(formatted, "")

	formatter.SetDecorated(isDecorated)

	return formatted
}

// Width calculate strings width
// TODO: handle unicode string
// TODO: handle mb_strwidth encoding
func Width(text string) int {
	return len(text)
}
