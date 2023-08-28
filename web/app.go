package web

import (
	morphadon "github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/core"
)

func CreateWebApp() core.App[*Context] {
	return morphadon.CreateApp[*Context]()
}

func GetInstance() core.App[*Context] {
	return morphadon.GetInstance[*Context]()
}
