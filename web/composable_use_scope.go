package web

import (
	"context"

	"github.com/mrparano1d/morphadon"
)

type ScopeKey string

const (
	ScopeSymbol ScopeKey = "scope"
)

type scopeInstance struct {
	scope morphadon.Scope
}

func (s *scopeInstance) Scope() morphadon.Scope {
	return s.scope
}

func provideScope(ctx *Context, scope morphadon.Scope) {
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
