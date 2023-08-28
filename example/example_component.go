package main

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/web"
)

type ExampleLayoutComponent struct {
	*morphadon.DefaultComponent[*web.Context]
}

var _ morphadon.Component[*web.Context] = (*ExampleLayoutComponent)(nil)

func ExampleLayout(children ...g.Node) *ExampleLayoutComponent {
	return &ExampleLayoutComponent{
		DefaultComponent: morphadon.NewDefaultComponentWithSlots[*web.Context](
			morphadon.Slots{
				"default": children,
			},
		),
	}
}

func (c *ExampleLayoutComponent) Assets() []morphadon.Asset[*web.Context] {
	return []morphadon.Asset[*web.Context]{}
}

func (c *ExampleLayoutComponent) Components() []morphadon.Component[*web.Context] {
	return []morphadon.Component[*web.Context]{
		web.HTML(),
	}
}

func (c *ExampleLayoutComponent) Render(data morphadon.SetupData) any {
	c.Context().Title = "Example Page"
	return c.Context().H(web.HTML(
		Div(
			Class("example-layout"),
			web.MustRenderSlot("default", c),
		),
	))
}
