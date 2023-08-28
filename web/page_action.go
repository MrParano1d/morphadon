package web

import (
	"github.com/mrparano1d/morphadon/core"
)

type PageAction struct {
	operation core.Operation
	scope     core.Scope
	page      WebPage
}

var _ core.Action[*Context] = (*PageAction)(nil)

func NewPageAction(operation core.Operation, scope core.Scope, page WebPage) *PageAction {
	return &PageAction{
		operation: operation,
		scope:     scope,
		page:      page,
	}
}

func (a *PageAction) Init(core.App[*Context]) error {
	return nil
}

func (a *PageAction) Operation() core.Operation {
	return a.operation
}

func (a *PageAction) Scope() core.Scope {
	return a.scope
}

func (a *PageAction) Renderer() core.Renderer[*Context] {
	return nil
}

func (a *PageAction) SetRenderer(core.Renderer[*Context]) {
}

func (a *PageAction) Components() []core.Component[*Context] {
	return []core.Component[*Context]{a.page}
}

func (a *PageAction) Assets() []core.Asset[*Context] {
	return make([]core.Asset[*Context], 0)
}

func (a *PageAction) Execute(ctx *Context) any {
	provideScope(ctx, a.scope)
	return ctx.H(a.page)
}

type WebPage interface {
	core.Component[*Context]
}

type DefaultWebPage struct {
	*core.DefaultComponent[*Context]
}

var _ core.Component[*Context] = (*DefaultWebPage)(nil)

func NewWebPage() WebPage {
	return &DefaultWebPage{
		DefaultComponent: core.NewDefaultComponent[*Context](),
	}
}
