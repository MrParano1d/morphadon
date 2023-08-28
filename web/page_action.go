package web

import "github.com/mrparano1d/morphadon"

type PageAction struct {
	operation morphadon.Operation
	scope     morphadon.Scope
	page      WebPage
}

var _ morphadon.Action[*Context] = (*PageAction)(nil)

func NewPageAction(operation morphadon.Operation, scope morphadon.Scope, page WebPage) *PageAction {
	return &PageAction{
		operation: operation,
		scope:     scope,
		page:      page,
	}
}

func (a *PageAction) Init(morphadon.App[*Context]) error {
	return nil
}

func (a *PageAction) Operation() morphadon.Operation {
	return a.operation
}

func (a *PageAction) Scope() morphadon.Scope {
	return a.scope
}

func (a *PageAction) Renderer() morphadon.Renderer[*Context] {
	return nil
}

func (a *PageAction) SetRenderer(morphadon.Renderer[*Context]) {
}

func (a *PageAction) Components() []morphadon.Component[*Context] {
	return []morphadon.Component[*Context]{a.page}
}

func (a *PageAction) Assets() []morphadon.Asset[*Context] {
	return make([]morphadon.Asset[*Context], 0)
}

func (a *PageAction) Execute(ctx *Context) any {
	provideScope(ctx, a.scope)
	return ctx.H(a.page)
}

type WebPage interface {
	morphadon.Component[*Context]
}

type DefaultWebPage struct {
	*morphadon.DefaultComponent[*Context]
}

var _ morphadon.Component[*Context] = (*DefaultWebPage)(nil)

func NewWebPage() WebPage {
	return &DefaultWebPage{
		DefaultComponent: morphadon.NewDefaultComponent[*Context](),
	}
}
