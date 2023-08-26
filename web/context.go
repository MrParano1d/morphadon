package web

import (
	"context"

	"github.com/marlaone/engine/core"
)

type Context struct {
	ctx context.Context
}

var _ core.Context = (*Context)(nil)

func NewContext() *Context {
	return &Context{
		ctx: context.Background(),
	}
}

func (c *Context) Context() context.Context {
	return c.ctx
}
