package style

import (
	"github.com/kilip/go-console/output"
	"io"
	"strings"
)

// OutputStyle Decorates Output to add console style guide helpers.
type OutputStyle struct {
	input io.Reader
	output.IOutput
}

// NewLine Add newline.
func (os *OutputStyle) NewLine() {
	os.NewLineC(1)
}

// NewLineC Add given count newline(s).
func (os *OutputStyle) NewLineC(count int) {
	os.Write(strings.Repeat("\n", count))
}
