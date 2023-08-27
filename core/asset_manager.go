package core

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

type Asset[C Context] interface {
	// Init is called once the asset is registered.
	Init(App[C]) error

	Path() string
	SetPath(string)
	Type() AssetType
	Scope() Scope
	SetScope(Scope)
}

type AssetManager[C Context] interface {
	Init(app App[C]) error

	// RegisterAsset registers an asset to the asset manager.
	// The asset manager will use the asset's name to identify the asset.
	RegisterAsset(asset Asset[C]) error

	ScopeAssets(scope Scope) []Asset[C]

	Assets() []Asset[C]

	Build() error
}
