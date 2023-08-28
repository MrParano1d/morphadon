package web

import (
	"context"

	g "github.com/maragudk/gomponents"
	"github.com/marlaone/engine/core"
)

type Context struct {
	ctx context.Context

	Language string
	Title    string
	BaseURL  string

	BodyAttrs g.Node

	Meta []map[string]string
}

var _ core.Context = (*Context)(nil)

func NewContext() *Context {
	return &Context{
		ctx:       context.Background(),
		Language:  "en",
		Title:     "",
		BaseURL:   "/",
		BodyAttrs: nil,
		Meta:      make([]map[string]string, 0),
	}
}

func (c *Context) H(n any) g.Node {
	return h(n, c)
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) SetContext(ctx context.Context) {
	c.ctx = ctx
}
