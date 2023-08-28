package web

import (
	"github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/core"
)

type useAssetsInstance struct {
}

func useAssets() *useAssetsInstance {
	return &useAssetsInstance{}
}

func (i *useAssetsInstance) All() []core.Asset[*Context] {
	return morphadon.GetInstance[*Context]().AssetManager().Assets()
}

func (i *useAssetsInstance) Scoped(scope core.Scope) []core.Asset[*Context] {
	return morphadon.GetInstance[*Context]().AssetManager().ScopeAssets(scope)
}
