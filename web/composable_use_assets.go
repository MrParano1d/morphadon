package web

import (
	"github.com/marlaone/engine"
	"github.com/marlaone/engine/core"
)

type useAssetsInstance struct {
}

func useAssets() *useAssetsInstance {
	return &useAssetsInstance{}
}

func (i *useAssetsInstance) All() []core.Asset[*Context] {
	return engine.GetInstance[*Context]().AssetManager().Assets()
}

func (i *useAssetsInstance) Scoped(scope core.Scope) []core.Asset[*Context] {
	return engine.GetInstance[*Context]().AssetManager().ScopeAssets(scope)
}
