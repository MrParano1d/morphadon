package web

import (
	"context"

	g "github.com/maragudk/gomponents"
	"github.com/marlaone/engine"
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

func (c *Context) H(n any) g.Node {
	return h(n, c)
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Assets() []core.Asset[*Context] {
	return engine.GetInstance[*Context]().AssetManager().Assets()
}

func (c *Context) ScopeAssets(scope core.Scope) []core.Asset[*Context] {
	return engine.GetInstance[*Context]().AssetManager().ScopeAssets(scope)
}
