package morphadon

import "io"

type lazy struct {
	render func() Renderable
}

func (l *lazy) Render(w io.Writer) error {
	return l.render().Render(w)
}

func Lazy(cb func() Renderable) *lazy {
	return &lazy{cb}
}
