package formatter

import "errors"

//StyleStack manages style stack during formatting text
type StyleStack struct {
	styles     []*Style
	emptyStyle *Style
}

//NewStyleStack creates new StyleStack object
func NewStyleStack() *StyleStack {
	return &StyleStack{
		emptyStyle: NewFormatterStyle("", ""),
		styles:     []*Style{},
	}
}

//GetEmptyStyle returns empty style for this stack
func (s *StyleStack) GetEmptyStyle() *Style {
	return s.emptyStyle
}

//SetEmptyStyle set empty Style for this StyleStack
func (s *StyleStack) SetEmptyStyle(style *Style) {
	s.emptyStyle = style
}

//Reset resets stack (empty internal styles array)
func (s *StyleStack) Reset() {
	s.styles = []*Style{}
}

//Push pushes a Style into stack
func (s *StyleStack) Push(style *Style) {
	s.styles = append(s.styles, style)
}

//Current computes current style with stacks top codes.
func (s *StyleStack) Current() *Style {
	if 0 == len(s.styles) {
		return s.emptyStyle
	}

	return s.styles[len(s.styles)-1]
}

//Pop remove last Style from the StyleStack.
func (s *StyleStack) Pop() (ps *Style, err error) {
	if 0 == len(s.styles) {
		return s.emptyStyle, nil
	}

	last := s.styles[len(s.styles)-1]

	var styles []*Style

	for i := 0; i < (len(s.styles) - 1); i++ {
		styles = append(styles, s.styles[i])
	}
	s.styles = styles

	return last, nil
}

//PopS remove given Style from StyleStack
func (s *StyleStack) PopS(style *Style) (ps *Style, err error) {
	if 0 == len(s.styles) {
		return s.emptyStyle, nil
	}

	for i := len(s.styles) - 1; i >= 0; i-- {
		stack := s.styles[i]
		if stack.Apply("") == style.Apply("") {
			s.styles = s.styles[0:i]
			return stack, nil
		}
	}

	return nil, errors.New("incorrectly nested style tag found")
}
