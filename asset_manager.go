package morphadon

import (
	"slices"
)

type AssetType string

const (
	AssetTypeCSS  AssetType = "css"
	AssetTypeJS   AssetType = "js"
	AssetTypeHTML AssetType = "html"
	AssetTypeJSON AssetType = "json"
	AssetTypePNG  AssetType = "png"
	AssetTypeJPG  AssetType = "jpg"
	AssetTypeGIF  AssetType = "gif"
	AssetTypeICO  AssetType = "ico"
	AssetTypeSVG  AssetType = "svg"
	AssetTypeAny  AssetType = "any"
)

var AssetImageTypes = []AssetType{
	AssetTypePNG,
	AssetTypeJPG,
	AssetTypeGIF,
	AssetTypeICO,
	AssetTypeSVG,
}

func IsImageAssetType(assetType AssetType) bool {
	return slices.Contains(AssetImageTypes, assetType)
}

type Asset interface {
	// Init is called once the asset is registered.
	Init(App) error

	Path() string
	SetPath(string)
	TargetPath() string
	SetTargetPath(string)
	Type() AssetType
	Scope() Scope
	SetScope(Scope)
}

type AssetManager interface {
	Init(app App) error

	SetConfig(*AssetManagerConfig)
	Config() *AssetManagerConfig

	// RegisterAsset registers an asset to the asset manager.
	// The asset manager will use the asset's name to identify the asset.
	RegisterAsset(asset Asset) error

	ScopeAssets(scope Scope) []Asset

	Assets() []Asset

	SrcDir() string
	SetSrcDir(string)

	Build() error
}
