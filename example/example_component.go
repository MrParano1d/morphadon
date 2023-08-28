package main

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/marlaone/engine/core"
	"github.com/marlaone/engine/web"
)

type ExampleLayoutComponent struct {
	*core.DefaultComponent[*web.Context]
}

var _ core.Component[*web.Context] = (*ExampleLayoutComponent)(nil)

func ExampleLayout(children ...g.Node) *ExampleLayoutComponent {
	return &ExampleLayoutComponent{
		DefaultComponent: core.NewDefaultComponentWithSlots[*web.Context](
			core.Slots{
				"default": children,
			},
		),
	}
}

func (c *ExampleLayoutComponent) Assets() []core.Asset[*web.Context] {
	return []core.Asset[*web.Context]{
		web.NewCSSAsset("example.css", core.ScopeGlobal),
		web.NewJSAsset("example.ts", core.ScopeGlobal),
	}
}

func (c *ExampleLayoutComponent) Render(data core.SetupData) any {
	return c.Context().H(web.HTML(
	Div(
		Class("example-layout"),
		web.MustRenderSlot("default", c),
	),
	))
}
