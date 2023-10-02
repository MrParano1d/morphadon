package morphadon

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chai2010/webp"
	"github.com/evanw/esbuild/pkg/api"
)

type AssetManagerDefault struct {
	assets []Asset

	srcDir string
}

var _ AssetManager = (*AssetManagerDefault)(nil)

func NewAssetManagerDefault() *AssetManagerDefault {
	return &AssetManagerDefault{
		assets: make([]Asset, 0),
		srcDir: ".",
	}
}

func (a *AssetManagerDefault) Init(app App) error {
	return nil
}

func (a *AssetManagerDefault) SetConfig(*AssetManagerConfig) {
}

func (a *AssetManagerDefault) Config() *AssetManagerConfig {
	return nil
}

func (a *AssetManagerDefault) SrcDir() string {
	return a.srcDir
}

func (a *AssetManagerDefault) SetSrcDir(srcDir string) {
	a.srcDir = srcDir
}

func (a *AssetManagerDefault) getAssetFilePath(asset Asset) string {
	return filepath.Join(a.srcDir, asset.Path())
}

func (a *AssetManagerDefault) RegisterAsset(asset Asset) error {
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
		abs1 := registeredAsset.Path()
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

func (a *AssetManagerDefault) ScopeAssets(scope Scope) []Asset {
	var assets []Asset
	for _, asset := range a.assets {
		if asset.Scope() == scope {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (a *AssetManagerDefault) Assets() []Asset {
	return a.assets
}

func (a *AssetManagerDefault) Build() error {
	return nil
}

type ScopeGroupedAssets map[Scope][]Asset

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

type WebAssetManager struct {
	*AssetManagerDefault

	config *AssetManagerConfig
}

var _ AssetManager = (*WebAssetManager)(nil)

func NewAssetManager() *WebAssetManager {
	return &WebAssetManager{
		AssetManagerDefault: NewAssetManagerDefault(),
		config:              NewDefaultAssetManagerConfig(),
	}
}

func NewAssetManagerWithConfig(config *AssetManagerConfig) *WebAssetManager {
	if config == nil {
		config = NewDefaultAssetManagerConfig()
	}
	assetManager := &WebAssetManager{
		AssetManagerDefault: NewAssetManagerDefault(),
		config:              config,
	}

	assetManager.SetSrcDir(config.SrcDir)

	return assetManager
}

func (a *WebAssetManager) SetConfig(config *AssetManagerConfig) {
	a.config = config
	a.SetSrcDir(config.SrcDir)
}

func (a *WebAssetManager) Config() *AssetManagerConfig {
	return a.config
}

func (a *WebAssetManager) findAssetType(assetTypes ...AssetType) []Asset {
	var assets []Asset
	for _, asset := range a.Assets() {
		if slices.Contains(assetTypes, asset.Type()) {
			assets = append(assets, asset)
		}
	}
	return assets
}

func (a *WebAssetManager) filterGlobalAndScopedAssets(assets []Asset) ([]Asset, ScopeGroupedAssets) {
	var globalAssets []Asset
	scopedAssets := make([]Asset, 0, len(assets))

	for _, asset := range assets {
		if asset.Scope() == ScopeGlobal || asset.Scope() == ScopeMultiple {
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
			scopeGroupedAssets[asset.Scope()] = make([]Asset, 0)
		}
		scopeGroupedAssets[asset.Scope()] = append(scopeGroupedAssets[asset.Scope()], asset)
	}

	globalAssets = slices.Compact(globalAssets)
	for k, v := range scopeGroupedAssets {
		scopeGroupedAssets[k] = slices.Compact(v)
	}

	return globalAssets, scopeGroupedAssets
}

func (a *WebAssetManager) transformCSS(outputFile string, assets []Asset) error {

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

		absTmpFilePath, err := filepath.Abs(filepath.Join(a.config.SrcDir, tmpFile))
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
		entryPoints[i] = filepath.Join(a.config.SrcDir, tmpFile)

		asset.SetTargetPath("/" + filepath.Join("public", outputFile))
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
		Loader: map[string]api.Loader{
			".css": api.LoaderCSS,
			".ttf": api.LoaderFile,
		},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifySyntax:      true,
		MinifyIdentifiers: true,
		Write:             true,
		Outfile:           filepath.Join(a.config.OutputDir, outputFile),
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
		for _, error := range css.Errors {
			fmt.Printf("error building css: %s\n", error.Text)
		}
		return fmt.Errorf("failed to build css see above for errors")
	}

	return nil
}

func (a *WebAssetManager) transformJS(outputFile string, assets []Asset) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	scripts := make([]string, len(assets))
	for i, asset := range assets {
		scripts[i] = fmt.Sprintf("import \"%s\";", asset.Path())
		asset.SetTargetPath("/" + filepath.Join("public", outputFile))
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

func (a *WebAssetManager) decodeImage(path string) (image image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image %s: %w", path, err)
	}
	defer file.Close()

	switch filepath.Ext(path) {
	case ".png":
		return png.Decode(file)
	case ".jpg":
		return jpeg.Decode(file)
	case ".gif":
		return gif.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image type: %s", filepath.Ext(path))
	}
}

func (a *WebAssetManager) transformImages(assets []Asset) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	outputDir := filepath.Join(cwd, a.config.OutputDir)
	imagesDir := filepath.Join(outputDir, "..")

	var buf bytes.Buffer
	for _, asset := range assets {
		rel, err := filepath.Rel(imagesDir, asset.Path())
		if err != nil {
			return fmt.Errorf("failed to get relative path for image %s: %w", asset.Path(), err)
		}

		imageDir := filepath.Join(outputDir, rel)
		err = os.MkdirAll(imageDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create image dir %s: %w", imageDir, err)
		}

		// check if image is webp or svg then just copy it
		switch filepath.Ext(asset.Path()) {
		case ".svg", ".webp":
			assetData, err := os.ReadFile(asset.Path())
			if err != nil {
				return fmt.Errorf("failed to read image %s: %w", asset.Path(), err)
			}

			err = os.WriteFile(imageDir, assetData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write image %s: %w", asset.Path(), err)
			}
			continue
		}

		decodedImg, err := a.decodeImage(asset.Path())
		if err != nil {
			return fmt.Errorf("failed to decode image %s: %w", asset.Path(), err)
		}
		if err = webp.Encode(&buf, decodedImg, &webp.Options{Lossless: false, Quality: 80}); err != nil {
			return fmt.Errorf("failed to encode image %s: %w", asset.Path(), err)
		}

		if err = os.WriteFile(filepath.Join(filepath.Dir(imageDir), filepath.Base(strings.TrimSuffix(imageDir, filepath.Ext(imageDir)))+".webp"), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write image %s: %w", asset.Path(), err)
		}
		buf.Reset()
	}

	return nil
}

func (a *WebAssetManager) AssetPathToBuiltPath(assetPath string) string {

	outputDir := a.config.OutputDir
	imageDir := filepath.Join(outputDir, assetPath)

	// check if image is webp or svg then use the same path
	switch filepath.Ext(assetPath) {
	case ".svg", ".webp":
		return filepath.Join(imageDir, assetPath)
	}

	// everything else is converted to webp
	return fmt.Sprintf("public/%s.webp", strings.TrimSuffix(assetPath, filepath.Ext(assetPath)))
}

func (a *WebAssetManager) BuildCSS() error {

	globalStylesheets, scopedStylesheets := a.filterGlobalAndScopedAssets(a.findAssetType(AssetTypeCSS))

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

func (a *WebAssetManager) BuildJS() error {

	globalScripts, scopedScripts := a.filterGlobalAndScopedAssets(a.findAssetType(AssetTypeJS))

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

func (a *WebAssetManager) BuildImages() error {

	globalImages, scopedImages := a.filterGlobalAndScopedAssets(a.findAssetType(AssetTypePNG, AssetTypeJPG, AssetTypeGIF, AssetTypeSVG))

	if len(globalImages) > 0 {
		if err := a.transformImages(globalImages); err != nil {
			return fmt.Errorf("failed to transform global images: %w", err)
		}
	}

	for scope, images := range scopedImages {
		if err := a.transformImages(images); err != nil {
			return fmt.Errorf("failed to transform images for %s scope: %w", scope.String(), err)
		}
	}

	return nil
}

func (a *WebAssetManager) Build() error {

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

	err = a.BuildImages()
	if err != nil {
		return fmt.Errorf("failed to build images: %w", err)
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
