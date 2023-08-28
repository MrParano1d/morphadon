package engine

import "github.com/marlaone/engine/core"

var app any

func GetInstance[C core.Context]() core.App[C] {
	if app == nil {
		panic("App not initialized")
	}
	return app.(core.App[C])
}

func CreateApp[C core.Context]() core.App[C] {
	app = core.NewDefaultApp[C]()
	return app.(core.App[C])
}

func DefineComponent[C core.Context](c core.Component[C]) core.Component[C] {
	app := GetInstance[C]()
	app.RegisterComponent(c)
	return c
}
