package main

import (
	"github.com/mrparano1d/morphadon/web"
)

func main() {

	app := web.CreateWebApp()
	app.Use(web.NewWebPlugin(
		web.PluginWithAssetManagerConfig(&web.AssetManagerConfig{
			SrcDir:    "./web",
			OutputDir: "./web/public",
		}),
		web.PluginWithPresenterHttpDir("./web/public"),
	))
	app.RegisterSystem(NewExampleSystem())

	err := app.Mount()
	if err != nil {
		panic(err)
	}

}
