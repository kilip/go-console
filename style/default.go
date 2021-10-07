package style

import (
	"fmt"
	"github.com/kilip/go-console/formatter"
	"github.com/kilip/go-console/helper"
	"github.com/kilip/go-console/output"
	"github.com/kilip/go-wordwrap"
	"io"
	"runtime"
	"strings"
)

type DefaultStyle struct {
	lineLength     int
	bufferedOutput *output.TrimmedBufferOutput
	*OutputStyle
}

func NewDefaultStyle(reader io.Reader, o output.IOutput) *DefaultStyle {

	maxLength := 2
	if runtime.GOOS == "windows" {
		maxLength = 4
	}
	bufferedOutput := output.NewTrimmedBufferOutput(maxLength)
	bufferedOutput.SetDecorated(false)

	return &DefaultStyle{
		lineLength:     120,
		bufferedOutput: bufferedOutput,
		OutputStyle: &OutputStyle{
			input:   reader,
			IOutput: o,
		},
	}
}

func (ds *DefaultStyle) SetDecorated(decorated bool) {
	ds.bufferedOutput.SetDecorated(decorated)
	ds.IOutput.SetDecorated(decorated)
}

// NewLine Add newline.
func (ds *DefaultStyle) NewLine() {
	ds.NewLineC(1)
}

// NewLineC Add given count newline(s).
func (ds *DefaultStyle) NewLineC(count int) {
	ds.WriteO(strings.Repeat("\n", count), false, output.FormatNormal)
}

func (ds *DefaultStyle) Block(messages interface{}, blockType string) {
	ds.BlockO(messages, blockType, "info", " ", false, true)
}

func (ds *DefaultStyle) BlockO(messages interface{}, blockType string, style string, prefix string, padding bool, escape bool) {
	ds.autoPrependBlock()
	text := ds.createBlock(messages, blockType, style, prefix, padding, escape)
	ds.Writeln(text)
	ds.NewLine()
}

func (ds *DefaultStyle) Caution(message string) {
	ds.BlockO(message, "CAUTION", "fg=white;bg=red", " ! ", true, true)
}

func (ds *DefaultStyle) Title(message string) {
	ds.autoPrependBlock()
	ds.Writeln(fmt.Sprintf("<comment>%s</>", formatter.EscapeTrailingBackslash(message)))
	ds.Writeln(fmt.Sprintf(
		"<comment>%s</>",
		strings.Repeat("=", helper.Width(helper.RemoveDecoration(ds.GetFormatter(), message))),
	),
	)
	ds.NewLine()
}

func (ds *DefaultStyle) Warning(message string) {
	ds.BlockO(message, "WARNING", "fg=black;bg=yellow", " ", true, true)
}

func (ds *DefaultStyle) Error(message string) {
	ds.BlockO(message, "ERROR", "fg=white;bg=red", " ", true, true)
}

func (ds *DefaultStyle) Success(message string) {
	ds.BlockO(message, "OK", "fg=black;bg=green", " ", true, true)
}

func (ds *DefaultStyle) Note(message string) {
	ds.BlockO(message, "NOTE", "fg=yellow", " ! ", true, true)
}

func (ds *DefaultStyle) Info(message string) {
	ds.BlockO(message, "INFO", "fg=yellow", " ", false, true)
}

func (ds *DefaultStyle) Listing(elements interface{}) {
	ds.autoPrependText()

	messages := helper.TextToSlices(elements)

	for _, v := range messages {
		ds.Writeln(fmt.Sprintf(" * %s", v))
	}
	ds.NewLine()
}

func (ds *DefaultStyle) Text(message interface{}) {
	ds.autoPrependText()

	conv := helper.TextToSlices(message)

	for _, v := range conv {
		ds.Writeln(fmt.Sprintf(" %s", v))
	}
}

func (ds *DefaultStyle) Comment(message interface{}) {
	messages := helper.TextToSlices(message)
	ds.BlockO(messages, "", "", "<fg=default;bg=default> // </>", false, false)
}

func (ds *DefaultStyle) autoPrependText() {
	fetched := ds.bufferedOutput.Fetch()
	//Prepend new line if last char isn't EOL:
	lnFetch := len(fetched) - 1
	if fetched == "" {
		ds.NewLine()
	} else if "\n" != fetched[lnFetch:] {
		ds.NewLine()
	}
}

func (ds *DefaultStyle) autoPrependBlock() {
	buff := ds.bufferedOutput.Fetch()
	replaced := strings.Replace(buff, "\n", "\n", -1)

	if "" != replaced {
		chars := replaced[len(replaced)-2:]
		if "" == string(chars[0]) {
			//empty history, so we should start with a new line.
			ds.NewLine()
		} else {
			//Prepend new line for each non LF chars (This means no blank line was output before)
			count := strings.Count(chars, "\n")
			sum := 2 - count
			if sum > 0 {
				ds.NewLineC(2 - count)
			}
		}
	} else {
		ds.NewLine()
	}
}

func (ds *DefaultStyle) createBlock(messages interface{}, blockType string, style string, prefix string, padding bool, escape bool) []string {
	cMsgs := helper.TextToSlices(messages)
	indentLength := 0
	prefixLength := helper.Width(helper.RemoveDecoration(ds.GetFormatter(), prefix))
	var lines []string
	var lineIndentation string

	if "" != blockType {
		blockType = fmt.Sprintf("[%s] ", blockType)
		indentLength = len(blockType)
		lineIndentation = strings.Repeat(" ", indentLength)
	}

	for key, message := range cMsgs {
		if "" == message {
			continue
		}
		if escape {
			message = formatter.Escape(message)
		}

		decorationLength := helper.Width(message) - helper.Width(helper.RemoveDecoration(ds.GetFormatter(), message))
		messageLineLength := ds.lineLength - prefixLength - indentLength + decorationLength
		if ds.lineLength < messageLineLength {
			messageLineLength = ds.lineLength
		}
		wrapped := wordwrap.Wrap(message, uint(messageLineLength))
		messageLines := strings.Split(wrapped, "\n")
		for _, messageLine := range messageLines {
			lines = append(lines, messageLine)
		}

		if len(cMsgs) > 1 && key < len(cMsgs)-1 {
			lines = append(lines, "")
		}
	}

	firstLineIndex := 0
	if padding && ds.IsDecorated() {
		firstLineIndex = 1
		lines = append([]string{""}, lines...)
		lines = append(lines, "")
	}

	//ret := ""
	for i, line := range lines {
		if "" != blockType {
			if i == firstLineIndex {
				line = blockType + line
			} else {
				line = lineIndentation + line
			}
		}
		line = prefix + line
		max := ds.lineLength - helper.Width(helper.RemoveDecoration(ds.GetFormatter(), line))
		if max < 0 {
			max = 0
		}
		line += strings.Repeat(" ", max)
		if "" != style {
			line = fmt.Sprintf("<%s>%s</>", style, line)
		}
		lines[i] = line
	}

	return lines
}

func (ds *DefaultStyle) writeBuffer(message interface{}, newLine bool, options int) {
	// We need to know if the last chars are newLine
	ds.bufferedOutput.WriteO(message, newLine, options)
}

func (ds *DefaultStyle) Write(message interface{}) {
	ds.WriteO(message, false, output.FormatNormal)
}

func (ds *DefaultStyle) WriteO(message interface{}, newLine bool, options int) {
	msgs := helper.TextToSlices(message)

	for _, m := range msgs {
		ds.IOutput.WriteO(m, newLine, options)
		ds.writeBuffer(m, newLine, output.FormatNormal)
	}
}

func (ds *DefaultStyle) Writeln(messages interface{}) {
	ds.WritelnO(messages, output.FormatNormal)
}

func (ds *DefaultStyle) WritelnO(messages interface{}, options int) {
	msgs := helper.TextToSlices(messages)
	for _, m := range msgs {
		ds.IOutput.WritelnO(m, options)
		ds.writeBuffer(m, true, options)
	}
}
