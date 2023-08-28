package main

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon/core"
	"github.com/mrparano1d/morphadon/web"
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
	return []core.Asset[*web.Context]{}
}

func (c *ExampleLayoutComponent) Components() []core.Component[*web.Context] {
	return []core.Component[*web.Context]{
		web.HTML(),
	}
}

func (c *ExampleLayoutComponent) Render(data core.SetupData) any {
	c.Context().Title = "Example Page"
	return c.Context().H(web.HTML(
		Div(
			Class("example-layout"),
			web.MustRenderSlot("default", c),
		),
	))
}
