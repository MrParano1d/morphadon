package web

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
)

type WebPlugin struct {
}

var _ morphadon.Plugin[*Context] = (*WebPlugin)(nil)

func NewWebPlugin() *WebPlugin {
	return &WebPlugin{}
}

func (p *WebPlugin) Init(app morphadon.App[*Context]) error {
	app.SetPresenter(NewHttpPresenter())
	app.Presenter().SetRenderer(NewWebRenderer())
	app.SetAssetManager(NewAssetManager())
	app.Presenter().RegisterAction(morphadon.NewActionFunc[*Context](OpHttpGet, "/", func(ctx *Context) any {
		return html.H1(g.Text("Hello, world!"))
	}))

	return nil
}
