package morphadon

import (
	"fmt"
	"path/filepath"
)

type WebAsset struct {
	path       string
	targetPath string
	assetType  AssetType
	scope      Scope
}

var _ Asset = (*WebAsset)(nil)

func NewWebAsset(path string, assetType AssetType, scope Scope) *WebAsset {
	return &WebAsset{
		path:       path,
		targetPath: "",
		assetType:  assetType,
		scope:      scope,
	}
}

func (a *WebAsset) Init(app App) error {
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

func (a *WebAsset) Type() AssetType {
	return a.assetType
}

func (a *WebAsset) Scope() Scope {
	return a.scope
}

func (a *WebAsset) SetScope(scope Scope) {
	a.scope = scope
}

func NewJSAsset(path string, scope ...Scope) *WebAsset {
	var s Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = ScopeComponent
	}
	return NewWebAsset(path, AssetTypeJS, s)
}

func NewCSSAsset(path string, scope ...Scope) *WebAsset {
	var s Scope
	if len(scope) > 0 {
		s = scope[0]
	} else {
		s = ScopeComponent
	}
	return NewWebAsset(path, AssetTypeCSS, s)
}

func NewImageAsset(path string, assetType ...AssetType) *WebAsset {
	var imageType AssetType
	if len(assetType) > 0 {
		imageType = assetType[0]
	} else {
		switch filepath.Ext(path) {
		case ".jpg", ".jpeg":
			imageType = AssetTypeJPG
		case ".png":
			imageType = AssetTypePNG
		case ".svg":
			imageType = AssetTypeSVG
		case ".gif":
			imageType = AssetTypeGIF
		default:
			imageType = AssetTypeAny
		}
	}
	if IsImageAssetType(imageType) {
		return NewWebAsset(path, imageType, ScopeGlobal)
	}
	panic(fmt.Errorf("invalid image asset type: %s", assetType))
}
