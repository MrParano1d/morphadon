package web

import (
	"github.com/mrparano1d/morphadon"
)

func CreateWebApp() morphadon.App[*Context] {
	return morphadon.CreateApp[*Context]()
}

func GetInstance() morphadon.App[*Context] {
	return morphadon.GetInstance[*Context]()
}
