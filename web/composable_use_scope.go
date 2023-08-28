package web

import (
	"context"

	"github.com/marlaone/morphadon/core"
)

type ScopeKey string

const (
	ScopeSymbol ScopeKey = "scope"
)

type scopeInstance struct {
	scope core.Scope
}

func (s *scopeInstance) Scope() core.Scope {
	return s.scope
}

func provideScope(ctx *Context, scope core.Scope) {
	ctx.SetContext(
		context.WithValue(ctx.Context(), ScopeSymbol, &scopeInstance{
			scope: scope,
		}),
	)
}

func useScope(ctx *Context) *scopeInstance {
	scope, ok := ctx.Context().Value(ScopeSymbol).(*scopeInstance)
	if !ok {
		panic("scope composable not provided")
	}
	return scope
}
