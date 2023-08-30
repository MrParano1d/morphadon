package web

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
)

type WebPlugin struct {
	assetManagerconfig *AssetManagerConfig
}

var _ morphadon.Plugin[*Context] = (*WebPlugin)(nil)

type WebPluginOption func(*WebPlugin)

func PluginWithAssetManagerConfig(config *AssetManagerConfig) WebPluginOption {
	return func(p *WebPlugin) {
		p.assetManagerconfig = config
	}
}

func NewWebPlugin(opts ...WebPluginOption) *WebPlugin {
	plugin := &WebPlugin{}

	for _, opt := range opts {
		opt(plugin)
	}

	return plugin
}

func (p *WebPlugin) Init(app morphadon.App[*Context]) error {
	app.SetPresenter(NewHttpPresenter())
	app.Presenter().SetRenderer(NewWebRenderer())
	app.SetAssetManager(NewAssetManagerWithConfig(p.assetManagerconfig))
	app.Presenter().RegisterAction(morphadon.NewActionFunc[*Context](OpHttpGet, "/", func(ctx *Context) any {
		return html.H1(g.Text("Hello, world!"))
	}))

	return nil
}
