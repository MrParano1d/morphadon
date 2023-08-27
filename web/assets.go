package web

import "github.com/marlaone/engine/core"

type WebAsset struct {
	path string
	assetType core.AssetType
	scope core.Scope
}

var _ core.Asset[*Context] = (*WebAsset)(nil)

func NewWebAsset(path string, assetType core.AssetType, scope core.Scope) *WebAsset {
	return &WebAsset{
		path: path,
		assetType: assetType,
		scope: scope,
	}
}

func (a *WebAsset) Init(app core.App[*Context]) error {
	return nil
}

func (a *WebAsset) Path() string {
	return a.path
}

func (a *WebAsset) SetPath(path string) {
	a.path = path
}

func (a *WebAsset) Type() core.AssetType {
	return a.assetType
}

func (a *WebAsset) Scope() core.Scope {
	return a.scope
}

func (a *WebAsset) SetScope(scope core.Scope) {
	a.scope = scope
}


