package style

import (
	"github.com/kilip/go-console/output"
	"io"
	"strings"
)

//OutputStyle Decorates Output to add console style guide helpers.
type OutputStyle struct {
	input io.Reader
	*output.Output
}

//NewLine Add newline.
func (os *OutputStyle) NewLine() error {
	return os.NewLineC(1)
}

//NewLineC Add given count newline(s).
func (os *OutputStyle) NewLineC(count int) error {
	return os.Write(strings.Repeat("\n", count))
}
