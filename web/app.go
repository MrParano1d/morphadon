package web

import (
	engine "github.com/marlaone/morphadon"
	"github.com/marlaone/morphadon/core"
)

func CreateWebApp() core.App[*Context] {
	return engine.CreateApp[*Context]()
}

func GetInstance() core.App[*Context] {
	return engine.GetInstance[*Context]()
}
