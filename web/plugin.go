package web

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
)

type WebPlugin struct {
	assetManagerconfig *AssetManagerConfig

	presenterHttpDir http.Dir
}

var _ morphadon.Plugin[*Context] = (*WebPlugin)(nil)

type WebPluginOption func(*WebPlugin)

func PluginWithAssetManagerConfig(config *AssetManagerConfig) WebPluginOption {
	return func(p *WebPlugin) {
		p.assetManagerconfig = config
	}
}

func PluginWithPresenterHttpDir(dir http.Dir) WebPluginOption {
	return func(p *WebPlugin) {
		p.presenterHttpDir = dir
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
	app.SetPresenter(NewHttpPresenter(HttpPresenterWithFilesDir(p.presenterHttpDir)))
	app.Presenter().SetRenderer(NewWebRenderer())
	app.SetAssetManager(NewAssetManagerWithConfig(p.assetManagerconfig))
	app.Presenter().RegisterAction(morphadon.NewActionFunc[*Context](OpHttpGet, "/", func(ctx *Context) any {
		return html.H1(g.Text("Hello, world!"))
	}))

	return nil
}
