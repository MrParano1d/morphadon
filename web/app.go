package web

import (
	"github.com/marlaone/engine"
	"github.com/marlaone/engine/core"
)

func CreateWebApp() core.App[*Context] {
	return engine.CreateApp[*Context]()
}

func GetInstance() core.App[*Context] {
	return engine.GetInstance[*Context]()
}
