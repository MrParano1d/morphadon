package engine

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/spf13/viper"
)

const publicDir = "./public/_marla/"
const tmpDir = "./tmp/_marla/"

type AssetFlag int

const (
	AssetFlagAsset      AssetFlag = iota
	AssetFlagStylesheet AssetFlag = iota
	AssetFlagScript     AssetFlag = iota
)

type PageAssets struct {
	Stylesheets []string
	Scripts     []string
}

type PageBundle map[string]PageAssets

func fillBundles(bundles PageBundle, routes []*Route) PageBundle {
	for _, r := range routes {
		bundleName := r.Name

		if bundleName == "" {
			bundleName = "global"
		}

		if len(r.Children) > 0 {
			bundles = fillBundles(bundles, r.Children)
		}

		bundle, ok := bundles[r.Name]
		if !ok {
			bundles[r.Name] = PageAssets{
				Stylesheets: []string{},
				Scripts:     []string{},
			}
		}

		bundle.Stylesheets = append(bundle.Stylesheets, r.Page.Stylesheets()...)
		bundle.Scripts = append(bundle.Scripts, r.Page.Scripts()...)

		for _, c := range resolvePageComponents([]Component{}, r.Page) {
			bundle.Stylesheets = append(bundle.Stylesheets, c.Stylesheets()...)
			bundle.Scripts = append(bundle.Scripts, c.Scripts()...)
		}

		bundle.Scripts = removeDuplicateStr(bundle.Scripts)
		bundle.Stylesheets = removeDuplicateStr(bundle.Stylesheets)

		bundles[r.Name] = bundle
	}

	return bundles
}

func (app *Marla) BuildPageBundles() PageAssets {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	bundledAssets := PageAssets{
		Stylesheets: []string{},
		Scripts:     []string{},
	}

	bundles := fillBundles(PageBundle{}, app.router.routes)

	for bundleName, bundleAssets := range bundles {
		if bundleName == "" {
			bundleName = "global"
		}

		if len(bundleAssets.Scripts) > 0 {
			Logger.Debug("asset::bundling::js\n", zap.String("bundle", bundleName))
			jsOutFile := fmt.Sprintf("%s%s.chunk.js", publicDir, bundleName)
			hasModified := false
			scripts := []string{}
			for _, s := range bundleAssets.Scripts {
				if !hasModified {
					hasModified = compareModified(s, jsOutFile)
				}
				scripts = append(scripts, fmt.Sprintf("import '%s';", s))
			}
			if hasModified {
				js := api.Build(api.BuildOptions{
					Stdin: &api.StdinOptions{
						Contents:   strings.Join(scripts, "\n"),
						Loader:     api.LoaderTS,
						ResolveDir: cwd,
					},
					Loader: map[string]api.Loader{
						".ts":   api.LoaderTS,
						".js":   api.LoaderJS,
						".json": api.LoaderJSON,
					},
					MinifyWhitespace:  true,
					MinifyIdentifiers: true,
					MinifySyntax:      true,
					Bundle:            true,
					Write:             true,
					Outfile:           jsOutFile,
				})
				bundledAssets.Scripts = append(bundledAssets.Scripts, jsOutFile)
				Logger.Debug("asset::bundling::js finished with warnings and errors\n", zap.Int("warnings", len(js.Warnings)), zap.Int("errors", len(js.Errors)))

				for _, err := range js.Errors {
					Logger.Warn(err.Text)
				}
			}
		}

		if len(bundleAssets.Stylesheets) > 0 {
			Logger.Debug("asset::bundling::css\n", zap.String("bundle", bundleName))
			var styleSheets []string
			var tmpFiles []string
			cssOutFile := fmt.Sprintf("%s%s.chunk.css", publicDir, bundleName)
			hasModified := false

			for _, s := range bundleAssets.Stylesheets {
				if !hasModified {
					hasModified = compareModified(s, cssOutFile)
				}
				if hasModified {
					tmpFile := strings.TrimSuffix(s, filepath.Base(s)) + getHashedName(s) + ".css"
					if err := cssLoader(s, tmpFile); err != nil {
						panic(err)
					}
					styleSheets = append(styleSheets, fmt.Sprintf("@import '%s';", tmpFile))
					tmpFiles = append(tmpFiles, tmpFile)
				}
			}
			if hasModified {
				css := api.Build(api.BuildOptions{
					Stdin: &api.StdinOptions{
						Contents:   strings.Join(styleSheets, "\n"),
						Loader:     api.LoaderCSS,
						ResolveDir: cwd,
					},
					Loader: map[string]api.Loader{
						".css": api.LoaderCSS,
						".ttf": api.LoaderFile,
					},
					Bundle:            true,
					MinifyWhitespace:  true,
					MinifyIdentifiers: true,
					MinifySyntax:      true,
					Write:             true,
					Outfile:           cssOutFile,
				})
				bundledAssets.Stylesheets = append(bundledAssets.Stylesheets, cssOutFile)
				Logger.Debug("asset::bundling::css finished with warnings and errors", zap.Int("warnings", len(css.Warnings)), zap.Int("errors", len(css.Errors)))
				for _, err := range css.Errors {
					Logger.Warn(err.Text)
				}
			}
			for _, tmpFile := range tmpFiles {
				if err := os.Remove(tmpFile); err != nil {
					Logger.Warn("failed to remove tmp file", zap.Error(err))
				}
			}
		}

	}

	return bundledAssets
}

func (app *Marla) Assets(flags ...AssetFlag) []string {

	if len(flags) == 0 {
		flags = []AssetFlag{AssetFlagAsset, AssetFlagScript, AssetFlagStylesheet}
	}

	var assets []string
	var components []Component
	for _, route := range app.router.routes {
		pagesMap := resolveRoutes("", "", []Page{}, route)
		pages := flattenPages(pagesMap)
		for _, p := range pages {
			if hasFlag(flags, AssetFlagAsset) {
				assets = append(assets, p.Assets()...)
			}
			if hasFlag(flags, AssetFlagScript) {
				assets = append(assets, p.Scripts()...)
			}
			if hasFlag(flags, AssetFlagStylesheet) {
				assets = append(assets, p.Stylesheets()...)
			}
		}

		components = append(components, resolvePageComponents([]Component{}, pages...)...)
	}

	for _, c := range components {
		if hasFlag(flags, AssetFlagAsset) {
			assets = append(assets, c.Assets()...)
		}
		if hasFlag(flags, AssetFlagScript) {
			assets = append(assets, c.Scripts()...)
		}
		if hasFlag(flags, AssetFlagStylesheet) {
			assets = append(assets, c.Stylesheets()...)
		}
	}

	return removeDuplicateStr(assets)
}

func hasFlag(flags []AssetFlag, flagToCheck AssetFlag) bool {
	for _, f := range flags {
		if f == flagToCheck {
			return true
		}
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
	var list []string
	allKeys := make(map[string]bool)
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := source.Close(); err != nil {
			Logger.Warn("failed to close source file", zap.Error(err))
		}
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := destination.Close(); err != nil {
			Logger.Warn("failed to close destination file", zap.Error(err))
		}
	}()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func glob(dir string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ext) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func cssLoader(filePath string, outPath string) error {
	Logger.Debug("asset::loading::css", zap.String("file", filePath))
	env := viper.GetString("MARLA_ENV")
	command := "npx cross-env NODE_ENV=" + env + " postcss " + filePath + " -o " + outPath
	parts := strings.Fields(command)
	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return fmt.Errorf("build css failed: %v - %s", err, string(data))
	}
	return nil
}

func jsLoader(filePath string) error {
	return nil
}

func decodeImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("image loader open file failed: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			Logger.Warn("failed to close image to decode", zap.Error(err))
		}
	}()

	switch filepath.Ext(filePath) {
	case ".png":
		return png.Decode(bufio.NewReader(f))
	case ".jpg", ".jpeg":
		return jpeg.Decode(bufio.NewReader(f))
	case ".gif":
		return gif.Decode(bufio.NewReader(f))
	}
	return nil, fmt.Errorf("DecodeImage - unknown extension: %s", filepath.Ext(filePath))
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getHashedName(filePath string) string {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) + "." + getMD5Hash(filePath)[0:7]
}

func GetAssetUrl(filePath string) string {
	url := getPublicFilePath(filePath)
	url = strings.Replace(url, "./public", "/static", 1)
	return url
}

func imageLoader(filePath string) error {
	Logger.Debug("assets::loading::image", zap.String("file", filePath))
	var buf bytes.Buffer
	img, err := decodeImage(filePath)
	if err != nil {
		return err
	}
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 80}); err != nil {
		return fmt.Errorf("image loader webp encode failed: %v", err)
	}

	if err = ioutil.WriteFile(publicDir+getHashedName(filePath)+".webp", buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("image loader write file failed: %v", err)
	}
	return nil
}

func tsLoader(filePath string) error {
	Logger.Debug("asset::loading::ts", zap.String("file", filePath))
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{filePath},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Loader: map[string]api.Loader{
			".png":  api.LoaderFile,
			".jpg":  api.LoaderFile,
			".jpeg": api.LoaderFile,
		},
		Outfile: publicDir + getHashedName(filePath) + ".js",
		Write:   true,
	})
	if len(result.Errors) > 0 {
		var errorMessages []string
		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Text)
		}
		return fmt.Errorf("ts loader errors: %s", strings.Join(errorMessages, ", "))
	}
	return nil
}

func copyLoader(filePath string) error {
	from := filePath
	to := publicDir + getHashedName(filePath) + filepath.Ext(filePath)
	Logger.Debug("assets::loading::copy copying files", zap.String("from", from), zap.String("to", to))
	_, err := copyFile(from, to)
	return err
}

func ImportGlob(directory string, extension string) {
	files, err := glob(directory, extension)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := Import(f); err != nil {
			panic(err)
		}
	}
}

func Import(filePath string) error {

	if !isModified(filePath) {
		return nil
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".css":
		return cssLoader(filePath, publicDir+getHashedName(filePath)+filepath.Ext(filePath))
	case ".js":
		return jsLoader(filePath)
	case ".ts":
		return tsLoader(filePath)
	case ".gif", ".png", ".jpg", ".jpeg":
		return imageLoader(filePath)
	case ".svg":
		return copyLoader(filePath)
	}

	return nil
}

func getPublicFilePath(filePath string) string {
	hashName := getHashedName(filePath)
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png", ".jpeg", ".jpg", ".gif":
		ext = ".webp"
	case ".ts":
		ext = ".js"
	}
	return publicDir + hashName + ext
}

func compareModified(newerFile string, olderFile string) bool {
	file, err := os.Stat(newerFile)

	if err != nil {
		return true
	}

	modifiedtime := file.ModTime()

	publicFile, err := os.Stat(olderFile)

	if err != nil {
		return true
	}

	publicModifiedTime := publicFile.ModTime()

	return modifiedtime.After(publicModifiedTime)
}

func isModified(filePath string) bool {
	return compareModified(filePath, getPublicFilePath(filePath))
}

func (app *Marla) Build() error {
	assets := app.Assets(AssetFlagAsset)

	app.BuildPageBundles()

	for _, a := range assets {
		if err := Import(a); err != nil {
			return err
		}
	}

	return nil
}
