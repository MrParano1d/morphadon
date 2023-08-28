package main

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/web"
)

const ExamplePageScope morphadon.Scope = "/example"

type ExamplePage struct {
	web.WebPage
}

func NewExamplePage() *ExamplePage {
	return &ExamplePage{
		WebPage: web.NewWebPage(),
	}
}

func (p *ExamplePage) Assets() []morphadon.Asset[*web.Context] {
	return []morphadon.Asset[*web.Context]{
		web.NewCSSAsset("example.css", morphadon.ScopeGlobal),
		web.NewJSAsset("example.ts", morphadon.ScopeGlobal),
		web.NewCSSAsset("example_button.css", ExamplePageScope),
		web.NewJSAsset("example_button.ts", ExamplePageScope),
	}
}

func (p *ExamplePage) Components() []morphadon.Component[*web.Context] {
	return []morphadon.Component[*web.Context]{
		ExampleLayout(),
	}
}

func (p *ExamplePage) Setup() morphadon.SetupData {
	return morphadon.SetupData{
		"greeting": "example",
	}
}

func (p *ExamplePage) Render(data morphadon.SetupData) any {
	return p.Context().H(
		ExampleLayout(
			H1(
				Class("text-blue-600"),
				g.Text(fmt.Sprintf("Hello %s", data["greeting"])),
			),
		),
	)
}
