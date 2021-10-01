package style

import (
	"github.com/kilip/console/formatter"
	"github.com/kilip/console/output"
	"io"
)

type DefaultStyle struct {
	*OutputStyle
}

func NewDefaultStyle(reader io.Reader, writer io.Writer) *DefaultStyle {
	out := output.NewOutput(writer, formatter.NewFormatter())
	return &DefaultStyle{
		&OutputStyle{
			input:  reader,
			Output: out,
		},
	}
}
