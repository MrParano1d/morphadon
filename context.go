package morphadon

import (
	"context"
	"net/http"

	g "github.com/maragudk/gomponents"
)

type ContextWithOption func(*Context)

func ContextWithRequest(r *http.Request) ContextWithOption {
	return func(c *Context) {
		c.Req = r
	}
}

func ContextWithContext(ctx context.Context) ContextWithOption {
	return func(c *Context) {
		c.context = ctx
	}
}

type Context struct {
	context context.Context

	Req *http.Request

	Language string
	Title    string
	BaseURL  string

	BodyAttrs g.Node

	Meta []map[string]string
}

func NewContext(opts ...ContextWithOption) *Context {
	ctx := &Context{
		context: context.Background(),
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx
}

func (c *Context) Context() context.Context {
	return c.context
}

func (c *Context) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *Context) H(n any) g.Node {
	return h(n, c)
}

