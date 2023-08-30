package morphadon

import (
	"net/url"
	"path/filepath"
)

type AssetManagerDefault[C Context] struct {
	assets []Asset[C]

	srcDir string
}

var _ AssetManager[*TodoContext] = (*AssetManagerDefault[*TodoContext])(nil)

func NewAssetManagerDefault[C Context]() *AssetManagerDefault[C] {
	return &AssetManagerDefault[C]{
		assets: make([]Asset[C], 0),
		srcDir: ".",
	}
}

func (a *AssetManagerDefault[C]) Init(app App[C]) error {
	return nil
}

func (a *AssetManagerDefault[C]) SrcDir() string {
	return a.srcDir
}

func (a *AssetManagerDefault[C]) SetSrcDir(srcDir string) {
	a.srcDir = srcDir
}

func (a *AssetManagerDefault[C]) getAssetFilePath(asset Asset[C]) string {
	return filepath.Join(a.srcDir, asset.Path())
}

func (a *AssetManagerDefault[C]) RegisterAsset(asset Asset[C]) error {
	exists := false
	for _, registeredAsset := range a.assets {
		// check if asset is a url
		url1, _ := url.ParseRequestURI(registeredAsset.Path())
		url2, _ := url.ParseRequestURI(asset.Path())
		if url1 != nil && url2 != nil {
			if url1.String() == url2.String() {
				exists = true
				registeredAsset.SetScope(ScopeMultiple)
				break
			}
		}

		// check if asset is a file
		abs1, _ := filepath.Abs(a.getAssetFilePath(registeredAsset))
		abs2, _ := filepath.Abs(a.getAssetFilePath(asset))
		if abs1 == abs2 {
			exists = true
			registeredAsset.SetScope(ScopeMultiple)
			break
		}
	}

	if !exists {
		url1, _ := url.ParseRequestURI(asset.Path())
		if url1 != nil {
			asset.SetPath(url1.String())
		} else {
			abs, _ := filepath.Abs(a.getAssetFilePath(asset))
			asset.SetPath(abs)
		}
		a.assets = append(a.assets, asset)
	}

	return nil
}

func (a *AssetManagerDefault[C]) ScopeAssets(scope Scope) []Asset[C] {
	var assets []Asset[C]
	for _, asset := range a.assets {
		if asset.Scope() == scope {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (a *AssetManagerDefault[C]) Assets() []Asset[C] {
	return a.assets
}

func (a *AssetManagerDefault[C]) Build() error {
	return nil
}
