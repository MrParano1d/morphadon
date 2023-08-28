package core

type AssetManagerNoop[C Context] struct {
}

var _ AssetManager[*TodoContext] = &AssetManagerNoop[*TodoContext]{}

func NewAssetManagerNoop[C Context]() *AssetManagerNoop[C] {
	return &AssetManagerNoop[C]{}
}

func (a *AssetManagerNoop[C]) Init(App[C]) error {
	return nil
}

func (a *AssetManagerNoop[C]) RegisterAsset(asset Asset[C]) error {
	return nil
}

func (a *AssetManagerNoop[C]) ScopeAssets(scope Scope) []Asset[C] {
	return nil
}

func (a *AssetManagerNoop[C]) Assets() []Asset[C] {
	return nil
}

func (a *AssetManagerNoop[C]) Build() error {
	return nil
}
