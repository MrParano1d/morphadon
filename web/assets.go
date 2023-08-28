package web

import (
	"fmt"

	"github.com/marlaone/engine/core"
)

type WebAsset struct {
	path       string
	targetPath string
	assetType  core.AssetType
	scope      core.Scope
}

var _ core.Asset[*Context] = (*WebAsset)(nil)

func NewWebAsset(path string, assetType core.AssetType, scope core.Scope) *WebAsset {
	return &WebAsset{
		path:       path,
		targetPath: "",
		assetType:  assetType,
		scope:      scope,
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

func (a *WebAsset) TargetPath() string {
	return a.targetPath
}

func (a *WebAsset) SetTargetPath(targetPath string) {
	a.targetPath = targetPath
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

func NewJSAsset(path string, scope ...core.Scope) *WebAsset {
	var s core.Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = core.ScopeNone
	}
	return NewWebAsset(path, core.AssetTypeJS, s)
}

func NewCSSAsset(path string, scope ...core.Scope) *WebAsset {
	var s core.Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = core.ScopeGlobal
	}
	return NewWebAsset(path, core.AssetTypeCSS, s)
}

func NewImageAsset(path string, assetType core.AssetType, scope ...core.Scope) *WebAsset {
	if core.IsImageAssetType(assetType) {
		var s core.Scope
		if len(scope) > 0 {
			s = scope[0]
		} else {
			s = core.ScopeGlobal
		}
		return NewWebAsset(path, assetType, s)
	}
	panic(fmt.Errorf("invalid image asset type: %s", assetType))
}
