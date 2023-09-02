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
	"github.com/mrparano1d/morphadon"
)

type ScopeGroupedAssets map[morphadon.Scope][]morphadon.Asset[*Context]

type AssetManagerConfig struct {
	SrcDir    string
	OutputDir string
}

func NewDefaultAssetManagerConfig() *AssetManagerConfig {
	return &AssetManagerConfig{
		SrcDir:    ".",
		OutputDir: "public",
	}
}

type AssetManager struct {
	*morphadon.AssetManagerDefault[*Context]

	config *AssetManagerConfig
}

var _ morphadon.AssetManager[*Context] = (*AssetManager)(nil)

func NewAssetManager() *AssetManager {
	return &AssetManager{
		AssetManagerDefault: morphadon.NewAssetManagerDefault[*Context](),
		config:              NewDefaultAssetManagerConfig(),
	}
}

func NewAssetManagerWithConfig(config *AssetManagerConfig) *AssetManager {
	if config == nil {
		config = NewDefaultAssetManagerConfig()
	}
	assetManager := &AssetManager{
		AssetManagerDefault: morphadon.NewAssetManagerDefault[*Context](),
		config:              config,
	}

	assetManager.SetSrcDir(config.SrcDir)

	return assetManager
}

func (a *AssetManager) findAssetType(assetTypes ...morphadon.AssetType) []morphadon.Asset[*Context] {
	var assets []morphadon.Asset[*Context]
	for _, asset := range a.Assets() {
		if slices.Contains(assetTypes, asset.Type()) {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (a *AssetManager) filterGlobalAndScopedAssets(assets []morphadon.Asset[*Context]) ([]morphadon.Asset[*Context], ScopeGroupedAssets) {
	var globalAssets []morphadon.Asset[*Context]
	scopedAssets := make([]morphadon.Asset[*Context], 0, len(assets))

	for _, asset := range assets {
		if asset.Scope() == morphadon.ScopeGlobal || asset.Scope() == morphadon.ScopeMultiple {
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
			scopeGroupedAssets[asset.Scope()] = make([]morphadon.Asset[*Context], 0)
		}
		scopeGroupedAssets[asset.Scope()] = append(scopeGroupedAssets[asset.Scope()], asset)
	}

	return globalAssets, scopeGroupedAssets
}

func (a *AssetManager) transformCSS(outputFile string, assets []morphadon.Asset[*Context]) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	srcCwd := filepath.Join(cwd, a.config.SrcDir)

	entryPoints := make([]string, len(assets))

	// transform css with postcss to tmp files

	for i, asset := range assets {
		file := asset.Path()
		// transform css with postcss
		tmpFile := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + "." + getMD5Hash(file)[0:7] + ".css"

		absTmpFilePath, err := filepath.Abs(filepath.Join(a.config.OutputDir, tmpFile))
		if err != nil {
			return fmt.Errorf("failed to get absolute path for tmp file: %w", err)
		}

		// write tmp file
		command := exec.Command("npx", "tailwindcss", "-i", file, "-o", absTmpFilePath)
		command.Dir = srcCwd
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err = command.Run()
		if err != nil {
			return fmt.Errorf("failed to transform css %s: %w", file, err)
		}
		entryPoints[i] = filepath.Join(a.config.OutputDir, tmpFile)

		asset.SetTargetPath(filepath.Join("public", outputFile))
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
		Outfile:           filepath.Join(a.config.OutputDir, outputFile),
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

func (a *AssetManager) transformJS(outputFile string, assets []morphadon.Asset[*Context]) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	scripts := make([]string, len(assets))
	for i, asset := range assets {
		scripts[i] = fmt.Sprintf("import \"%s\";", asset.Path())
		asset.SetTargetPath(filepath.Join("public", outputFile))
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
		Outfile:           filepath.Join(a.config.OutputDir, outputFile),
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

	globalStylesheets, scopedStylesheets := a.filterGlobalAndScopedAssets(a.findAssetType(morphadon.AssetTypeCSS))

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

	globalScripts, scopedScripts := a.filterGlobalAndScopedAssets(a.findAssetType(morphadon.AssetTypeJS))

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
	err := os.RemoveAll(a.config.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to remove output dir: %w", err)
	}

	// create output dir if not exists
	if _, err := os.Stat(a.config.OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(a.config.OutputDir, 0755)
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

	// globalImages, scopedImages := a.filterGlobalAndScopedAssets(a.findAssetType(morphadon.AssetTypePNG, morphadon.AssetTypeJPG, morphadon.AssetTypeGIF, morphadon.AssetTypeSVG))

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
