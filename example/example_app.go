package main

import (
	"github.com/marlaone/engine/web"
)

func main() {

	app := web.CreateWebApp()
	app.Use(web.NewWebPlugin())

	err := app.Mount()
	if err != nil {
		panic(err)
	}

}
