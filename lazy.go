package morphadon

type lazy struct {
	render func() morphadon.Renderable
}

func (l *lazy) Render(w io.Writer) error {
	return l.render().Render(w)
}

func Lazy(cb func() morphadon.Renderable) *lazy {
	return &lazy{cb}
}
