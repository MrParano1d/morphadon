package main

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/marlaone/engine/core"
	"github.com/marlaone/engine/web"
)

type ExamplePage struct {
	web.WebPage
}

func NewExamplePage() *ExamplePage {
	return &ExamplePage{
		WebPage: web.NewWebPage(),
	}
}

func (p *ExamplePage) Assets() []core.Asset[*web.Context] {
	return []core.Asset[*web.Context]{
		web.NewCSSAsset("example.css", core.ScopeGlobal),
		web.NewJSAsset("example.ts", core.ScopeGlobal),
	}
}

func (p *ExamplePage) Setup() core.SetupData {
	return core.SetupData{
		"greeting": "example",
	}
}

func (p *ExamplePage) Render(data core.SetupData) any {
	return p.Context().H(
		ExampleLayout(
			H1(
				Class("text-blue-600"),
				g.Text(fmt.Sprintf("Hello %s", data["greeting"])),
			),
		),
	)
}
