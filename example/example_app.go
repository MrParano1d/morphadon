package main

import (
	"github.com/marlaone/engine/web"
)

func main() {

	app := web.CreateWebApp()
	app.Use(web.NewWebPlugin())
	app.RegisterSystem(NewExampleSystem())

	err := app.Mount()
	if err != nil {
		panic(err)
	}

}
