package web

import (
	"fmt"

	"github.com/mrparano1d/morphadon"
)

type WebAsset struct {
	path       string
	targetPath string
	assetType  morphadon.AssetType
	scope      morphadon.Scope
}

var _ morphadon.Asset[*Context] = (*WebAsset)(nil)

func NewWebAsset(path string, assetType morphadon.AssetType, scope morphadon.Scope) *WebAsset {
	return &WebAsset{
		path:       path,
		targetPath: "",
		assetType:  assetType,
		scope:      scope,
	}
}

func (a *WebAsset) Init(app morphadon.App[*Context]) error {
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

func (a *WebAsset) Type() morphadon.AssetType {
	return a.assetType
}

func (a *WebAsset) Scope() morphadon.Scope {
	return a.scope
}

func (a *WebAsset) SetScope(scope morphadon.Scope) {
	a.scope = scope
}

func NewJSAsset(path string, scope ...morphadon.Scope) *WebAsset {
	var s morphadon.Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = morphadon.ScopeComponent
	}
	return NewWebAsset(path, morphadon.AssetTypeJS, s)
}

func NewCSSAsset(path string, scope ...morphadon.Scope) *WebAsset {
	var s morphadon.Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = morphadon.ScopeComponent
	}
	return NewWebAsset(path, morphadon.AssetTypeCSS, s)
}

func NewImageAsset(path string, assetType morphadon.AssetType, scope ...morphadon.Scope) *WebAsset {
	if morphadon.IsImageAssetType(assetType) {
		var s morphadon.Scope
		if len(scope) > 0 {
			s = scope[0]
		} else {
			s = morphadon.ScopeGlobal
		}
		return NewWebAsset(path, assetType, s)
	}
	panic(fmt.Errorf("invalid image asset type: %s", assetType))
}
