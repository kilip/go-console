package output

// Null suppresses all output
type Null struct {
	*Output
}

// NewNullOutput creates and returns new Null output object
func NewNullOutput() *Null {
	output := &Output{}
	no := &Null{
		Output: output,
	}
	output.doWrite = no.doWrite
	return no
}

func (no Null) doWrite(message string, newLine bool) {

}
