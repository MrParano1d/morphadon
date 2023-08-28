package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/mrparano1d/morphadon/core"
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
	scopedAssets := make([]core.Asset[*Context], 0, len(assets))

	for _, asset := range assets {
		if asset.Scope() == core.ScopeGlobal || asset.Scope() == core.ScopeMultiple {
			globalAssets = append(globalAssets, asset)
			continue
		}

	}

	for _, asset := range assets {
		if !slices.Contains(globalAssets, asset) {
			scopedAssets = append(scopedAssets, asset)
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

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	entryPoints := make([]string, len(assets))

	// transform css with postcss to tmp files

	for i, asset := range assets {
		file := asset.Path()
		// transform css with postcss
		tmpFile := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + "." + getMD5Hash(file)[0:7] + ".css"

		// write tmp file
		command := exec.Command("npx", "tailwindcss", "-i", file, "-o", filepath.Join(a.config.outputDir, tmpFile))
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			return fmt.Errorf("failed to transform css %s: %w", file, err)
		}
		entryPoints[i] = filepath.Join(a.config.outputDir, tmpFile)
		asset.SetTargetPath(filepath.Join(a.config.outputDir, outputFile))
	}

	stylesheets := make([]string, len(entryPoints))
	for i, entryPoint := range entryPoints {
		stylesheets[i] = fmt.Sprintf("@import \"%s\";", entryPoint)
	}

	// build css with esbuild

	ctx, err := api.Context(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   strings.Join(stylesheets, "\n"),
			Loader:     api.LoaderCSS,
			ResolveDir: cwd,
		},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifySyntax:      true,
		MinifyIdentifiers: true,
		Write:             true,
		Outfile:           filepath.Join(a.config.outputDir, outputFile),
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
			".ttf": api.LoaderFile,
		},
	})

	css := ctx.Rebuild()

	if len(css.Warnings) > 0 {
		for _, warning := range css.Warnings {
			fmt.Printf("warning building css: %s\n", warning.Text)
		}
	}

	for _, file := range entryPoints {
		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("failed to remove tmp file: %w", err)
		}
	}

	if len(css.Errors) > 0 {
		return fmt.Errorf("failed to build css: %w", css.Errors)
	}

	return nil
}

func (a *AssetManager) transformJS(outputFile string, assets []core.Asset[*Context]) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	scripts := make([]string, len(assets))
	for i, asset := range assets {
		scripts[i] = fmt.Sprintf("import \"%s\";", asset.Path())
		asset.SetTargetPath(filepath.Join(a.config.outputDir, outputFile))
	}

	ctx, err := api.Context(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   strings.Join(scripts, "\n"),
			Loader:     api.LoaderTS,
			ResolveDir: cwd,
		},
		Loader: map[string]api.Loader{
			".ts":   api.LoaderTS,
			".js":   api.LoaderJS,
			".json": api.LoaderJSON,
			".png":  api.LoaderFile,
			".jpg":  api.LoaderFile,
			".gif":  api.LoaderFile,
			".svg":  api.LoaderFile,
			".ico":  api.LoaderFile,
		},
		MinifyWhitespace:  true,
		MinifySyntax:      true,
		MinifyIdentifiers: true,
		Bundle:            true,
		Write:             true,
		Outfile:           filepath.Join(a.config.outputDir, outputFile),
	})

	js := ctx.Rebuild()

	if len(js.Warnings) > 0 {
		for _, warning := range js.Warnings {
			fmt.Printf("warning building js: %s\n", warning.Text)
		}
	}

	if len(js.Errors) > 0 {
		for _, error := range js.Errors {
			fmt.Printf("error building js: %s\n", error.Text)
		}
		return fmt.Errorf("failed to build js see above for errors")
	}

	return nil
}

func (a *AssetManager) BuildCSS() error {

	globalStylesheets, scopedStylesheets := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypeCSS))

	if len(globalStylesheets) > 0 {
		if err := a.transformCSS("global.css", globalStylesheets); err != nil {
			return fmt.Errorf("failed to transform global stylesheets: %w", err)
		}
	}

	for scope, stylesheets := range scopedStylesheets {
		if err := a.transformCSS(scope.String()+".chunk.css", stylesheets); err != nil {
			return fmt.Errorf("failed to transform scoped stylesheet: %w", err)
		}
	}

	return nil
}

func (a *AssetManager) BuildJS() error {

	globalScripts, scopedScripts := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypeJS))

	if len(globalScripts) > 0 {
		// build global js
		if err := a.transformJS("global.js", globalScripts); err != nil {
			return fmt.Errorf("failed to transform global scripts: %w", err)
		}
	}

	for scope, scripts := range scopedScripts {
		// build scoped js
		if err := a.transformJS(scope.String()+".chunk.js", scripts); err != nil {
			return fmt.Errorf("failed to transform scoped scripts: %w", err)
		}
	}

	return nil
}

func (a *AssetManager) Build() error {

	// remove output dir
	err := os.RemoveAll(a.config.outputDir)
	if err != nil {
		return fmt.Errorf("failed to remove output dir: %w", err)
	}

	// create output dir if not exists
	if _, err := os.Stat(a.config.outputDir); os.IsNotExist(err) {
		err := os.Mkdir(a.config.outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output dir: %w", err)
		}
	}

	err = a.BuildCSS()
	if err != nil {
		return fmt.Errorf("failed to build css: %w", err)
	}

	err = a.BuildJS()
	if err != nil {
		return fmt.Errorf("failed to build js: %w", err)
	}

	// globalImages, scopedImages := a.filterGlobalAndScopedAssets(a.findAssetType(core.AssetTypePNG, core.AssetTypeJPG, core.AssetTypeGIF, core.AssetTypeSVG))

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
