package morphadon

import "net/http"

type Page interface {
	Component

	Middlewares() []func(http.Handler) http.Handler
}

type defaultPage struct {
	Component
}

func NewPage() Page {
	return &defaultPage{
		Component: NewDefaultComponent(),
	}
}

func (p *defaultPage) Middlewares() []func(http.Handler) http.Handler {
	return make([]func(http.Handler) http.Handler, 0)
}
