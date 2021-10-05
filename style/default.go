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

	return &DefaultStyle{
		lineLength:     120,
		bufferedOutput: bufferedOutput,
		OutputStyle: &OutputStyle{
			input:   reader,
			IOutput: o,
		},
	}
}

func (ds *DefaultStyle) Block(messages string, blockType string) {
	ds.BlockO(messages, blockType, "info", " ", false, true)
}

func (ds *DefaultStyle) BlockO(messages, blockType string, style string, prefix string, padding bool, escape bool) {
	ds.autoPrependBlock()

	text := ds.createBlock(messages, blockType, style, prefix, padding, escape)
	ds.Writeln(text)
}

func (ds *DefaultStyle) Caution(message string) {
	ds.BlockO(message, "CAUTION", "fg=white;bg=red", " ! ", true, true)
}

func (ds *DefaultStyle) autoPrependBlock() {
	ds.NewLine()
}

func (ds *DefaultStyle) createBlock(messages string, blockType string, style string, prefix string, padding bool, escape bool) string {
	indentLength := 0
	prefixLength := helper.Width(helper.RemoveDecoration(ds.GetFormatter(), prefix))
	var lines []string
	var lineIndentation string

	if "" != blockType {
		blockType = fmt.Sprintf("[%s] ", blockType)
		indentLength = len(blockType)
		lineIndentation = strings.Repeat(" ", indentLength)
	}

	for key, message := range strings.Split(messages, "\n") {
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

		if len(messages) > 1 && key < len(messages)-1 {
			lines = append(lines, "")
		}
	}

	firstLineIndex := 0
	if padding && ds.IsDecorated() {
		firstLineIndex = 1
		lines = append([]string{""}, lines...)
		lines = append(lines, "")
	}

	for i, line := range lines {
		if "" != blockType {
			if i == firstLineIndex {
				line = blockType + line
			} else {
				line = lineIndentation + line
			}
		}
		trimmed := strings.Trim(strings.Trim(line, " "), "\n")
		if trimmed != "" {
			line = prefix + line
			max := ds.lineLength - helper.Width(helper.RemoveDecoration(ds.GetFormatter(), line))
			if max < 0 {
				max = 0
			}
			line += strings.Repeat(" ", max)
		} else {
			line = "\n"
		}

		if "" != style {
			line = fmt.Sprintf("<%s>%s</>", style, line)
		}
		lines[i] = line
	}

	ret := strings.Join(lines, "\n")

	return ret
}

func (ds *DefaultStyle) writeBuffer(message string, newLine bool, options int) {
	// We need to know if the last chars are newLine
	ds.WriteO(message, newLine, options)
	ds.bufferedOutput.WriteO(message, newLine, options)
}

func (ds *DefaultStyle) Write(message string) {
	ds.writeBuffer(message, false, output.FormatNormal)
}

func (ds *DefaultStyle) Writeln(messages string) {
	ds.writeBuffer(messages, true, output.FormatNormal)
}

func (ds *DefaultStyle) WritelnO(messages string, options int) {
	ds.writeBuffer(messages, true, options)
}
