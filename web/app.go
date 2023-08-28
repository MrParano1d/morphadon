package web

import (
	morphadon "github.com/marlaone/morphadon"
	"github.com/marlaone/morphadon/core"
)

func CreateWebApp() core.App[*Context] {
	return morphadon.CreateApp[*Context]()
}

func GetInstance() core.App[*Context] {
	return morphadon.GetInstance[*Context]()
}
