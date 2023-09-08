package morphadon

import (
	"context"
)

type ScopeKey string

const (
	ScopeSymbol ScopeKey = "scope"
)

type scopeInstance struct {
	scope Scope
}

func (s *scopeInstance) Scope() Scope {
	return s.scope
}

func ProvideScope(ctx *Context, scope Scope) {
	ctx.SetContext(
		context.WithValue(ctx.Context(), ScopeSymbol, &scopeInstance{
			scope: scope,
		}),
	)
}

func UseScope(ctx *Context) *scopeInstance {
	scope, ok := ctx.Context().Value(ScopeSymbol).(*scopeInstance)
	if !ok {
		panic("scope composable not provided")
	}
	return scope
}
