package morphadon

import "net/http"

type Middleware func(http.Handler) http.Handler

func MiddlewaresToHandlerSlice(middlewares []Middleware) []func(http.Handler) http.Handler {
	handlers := make([]func(http.Handler) http.Handler, len(middlewares))
	for i, middleware := range middlewares {
		handlers[i] = middleware
	}
	return handlers
}

type Page interface {
	Component

	Middlewares() []Middleware
}

type defaultPage struct {
	Component
}

func NewPage() Page {
	return &defaultPage{
		Component: NewDefaultComponent(),
	}
}

func (p *defaultPage) Middlewares() []Middleware {
	return make([]Middleware, 0)
}
