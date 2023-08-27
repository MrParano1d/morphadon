package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/marlaone/engine/core"
)

type ScopeGroupedAssets map[core.Scope][]core.Asset[*Context]

type AssetManagerConfig struct {
	outputDir string
}

func NewDefaultAssetManagerConfig() *AssetManagerConfig {
	return &AssetManagerConfig{
		outputDir: "public",
	}
}

type AssetManager struct {
	*core.AssetManagerDefault[*Context]

	config *AssetManagerConfig
}

var _ core.AssetManager[*Context] = (*AssetManager)(nil)

func NewAssetManager() *AssetManager {
	return &AssetManager{
		AssetManagerDefault: core.NewAssetManagerDefault[*Context](),
		config:              NewDefaultAssetManagerConfig(),
	}
}

func NewAssetManagerWithConfig(config *AssetManagerConfig) *AssetManager {
	return &AssetManager{
		AssetManagerDefault: core.NewAssetManagerDefault[*Context](),
		config:              config,
	}
}

func (a *AssetManager) findAssetType(assetTypes ...core.AssetType) []core.Asset[*Context] {
	var assets []core.Asset[*Context]
	for _, asset := range a.Assets() {
		if slices.Contains(assetTypes, asset.Type()) {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (a *AssetManager) filterGlobalAndScopedAssets(assets []core.Asset[*Context]) ([]core.Asset[*Context], ScopeGroupedAssets) {
	var globalAssets []core.Asset[*Context]
	scopedAssets := make([]core.Asset[*Context], len(assets))
	copy(scopedAssets, assets)

	for i, asset := range assets {
		if asset.Scope() == core.ScopeGlobal {
			globalAssets = append(globalAssets, asset)
			scopedAssets = append(scopedAssets[:i], scopedAssets[i+1:]...)
			continue
		}

		if asset.Scope() == core.ScopeMultiple {
			globalAssets = append(globalAssets, asset)
			scopedAssets = append(scopedAssets[:i], scopedAssets[i+1:]...)
			continue
		}
	}

	scopeGroupedAssets := make(ScopeGroupedAssets)

	for _, asset := range scopedAssets {
		if _, ok := scopeGroupedAssets[asset.Scope()]; !ok {
			scopeGroupedAssets[asset.Scope()] = make([]core.Asset[*Context], 0)
		}
		scopeGroupedAssets[asset.Scope()] = append(scopeGroupedAssets[asset.Scope()], asset)
	}

	return globalAssets, scopeGroupedAssets
}

func (a *AssetManager) transformCSS(outputFile string, assets []core.Asset[*Context]) error {

	entryPoints := make([]string, len(assets))

	// transform css with postcss to tmp files

	for _, asset := range assets {
		file := asset.Path()
		// transform css with postcss
		tmpFile := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + "." + getMD5Hash(file)[0:7] + ".css"

		// write tmp file
		command := exec.Command("npx", "postcss", file, "-o", filepath.Join(a.config.outputDir, tmpFile))
		err := command.Run()
		if err != nil {
			return fmt.Errorf("failed to transform css %s: %w", file, err)
		}
		entryPoints = append(entryPoints, filepath.Join(a.config.outputDir, tmpFile))
	}

	// build css with esbuild


	ctx, err := api.Context(api.BuildOptions{
		EntryPoints: entryPoints,
		Bundle:      true,
		Outfile:     filepath.Join(a.config.outputDir, outputFile),
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
			".ttf": api.LoaderFile,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create esbuild context: %w", err)
	}

	result := ctx.Rebuild()

	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to build css: %w", result.Errors)
	}

	return nil
}

func (a *AssetManager) BuildCSS() error {
	globalStylesheets, scopedStylesheets := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypeCSS))

	if err := a.transformCSS("global.css", globalStylesheets); err != nil {
		return fmt.Errorf("failed to transform global stylesheets: %w", err)
	}

	for scope, stylesheets := range scopedStylesheets {
		if err := a.transformCSS(scope.String()+".chunk.css", stylesheets); err != nil {
			return fmt.Errorf("failed to transform scoped stylesheet: %w", err)
		}
	}

	return nil
}

func (a *AssetManager) Build() error {

	err := a.BuildCSS()
	if err != nil {
		return fmt.Errorf("failed to build css: %w", err)
	}

	// globalScripts, scopedScripts := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypeJS))

	// globalImages, scopedImages := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypePNG, core.AssetTypeJPG, core.AssetTypeGIF, core.AssetTypeSVG))

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
