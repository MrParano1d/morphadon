package engine

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (app *Marla) Stylesheets(routeName string) []string {
	stylesheets := []string{
		publicDir + "global.chunk.css",
		publicDir + fmt.Sprintf("%s.chunk.css", routeName),
	}

	for i, s := range stylesheets {
		if _, err := os.Stat(s); os.IsNotExist(err) {
			stylesheets = append(stylesheets[:i], stylesheets[i+1:]...)
		}
	}

	for i, s := range stylesheets {
		stylesheets[i] = path.Clean("/" + strings.Replace(s, "public", "/static", 1))
	}

	return stylesheets
}

func (app *Marla) Scripts(routeName string) []string {
	scripts := []string{
		publicDir + "global.chunk.js",
		publicDir + fmt.Sprintf("%s.chunk.js", routeName),
	}

	for i, s := range scripts {
		if _, err := os.Stat(s); os.IsNotExist(err) {
			scripts = append(scripts[:i], scripts[i+1:]...)
		}
	}

	for i, s := range scripts {
		scripts[i] = path.Clean("/" + strings.Replace(s, "public", "/static", 1))
	}

	return scripts
}
