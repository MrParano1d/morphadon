package morphadon

type AssetManagerNoop struct {
}

var _ AssetManager = &AssetManagerNoop{}

func NewAssetManagerNoop() *AssetManagerNoop {
	return &AssetManagerNoop{}
}

func (a *AssetManagerNoop) SetConfig(config *AssetManagerConfig) {
}

func (a *AssetManagerNoop) Config() *AssetManagerConfig {
	return nil
}

func (a *AssetManagerNoop) Init(App) error {
	return nil
}

func (a *AssetManagerNoop) SetSrcDir(string) {
}

func (a *AssetManagerNoop) SrcDir() string {
	return ""
}

func (a *AssetManagerNoop) RegisterAsset(asset Asset) error {
	return nil
}

func (a *AssetManagerNoop) ScopeAssets(scope Scope) []Asset {
	return nil
}

func (a *AssetManagerNoop) Assets() []Asset {
	return nil
}

func (a *AssetManagerNoop) Build() error {
	return nil
}
