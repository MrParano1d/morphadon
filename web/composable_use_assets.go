package web

import (
	"github.com/mrparano1d/morphadon"
)

type useAssetsInstance struct {
}

func useAssets() *useAssetsInstance {
	return &useAssetsInstance{}
}

func (i *useAssetsInstance) All() []morphadon.Asset[*Context] {
	return morphadon.GetInstance[*Context]().AssetManager().Assets()
}

func (i *useAssetsInstance) Scoped(scope morphadon.Scope) []morphadon.Asset[*Context] {
	return morphadon.GetInstance[*Context]().AssetManager().ScopeAssets(scope)
}
