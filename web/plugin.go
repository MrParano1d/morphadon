package web

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/marlaone/engine/core"
)

type WebPlugin struct {
}

var _ core.Plugin[*Context] = (*WebPlugin)(nil)

func NewWebPlugin() *WebPlugin {
	return &WebPlugin{}
}

func (p *WebPlugin) Init(app core.App[*Context]) error {
	app.SetPresenter(NewHttpPresenter())
	app.Presenter().SetRenderer(NewWebRenderer())
	app.Presenter().RegisterAction(core.NewActionFunc[*Context](OpHttpGet, "/", func(ctx *Context) any {
		return html.H1(g.Text("Hello, world!"))
	}))

	return nil
}
