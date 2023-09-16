package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/components"
)

type ExampleApp struct {
	morphadon.Page
}

func NewExampleApp() *ExampleApp {
	return &ExampleApp{
		Page: morphadon.NewPage(),
	}
}

func (c *ExampleApp) Render(data morphadon.SetupData) morphadon.Renderable {
	return c.Context().H(components.HTML(
		morphadon.RouterView(c.Context()),
	))
}

type HelloPage struct {
	morphadon.Page
}

func NewHelloPage() *HelloPage {
	return &HelloPage{
		Page: morphadon.NewPage(),
	}
}

func (c *HelloPage) Middlewares() []morphadon.Middleware {
	return []morphadon.Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-HELLO", "HELLO")
				next.ServeHTTP(w, r)
			})
		},
	}
}

func (c *HelloPage) Render(data morphadon.SetupData) morphadon.Renderable {
	return Div(
		g.Text("Hello"),
		c.Context().H(morphadon.RouterView(c.Context())),
	)
}

type WorldPage struct {
	morphadon.Page
}

func NewWorldPage() *WorldPage {
	return &WorldPage{
		Page: morphadon.NewPage(),
	}
}

func (c *WorldPage) Middlewares() []morphadon.Middleware {
	return []morphadon.Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-WORLD", "WORLD")
				next.ServeHTTP(w, r)
			})
		},
	}
}

func (c *WorldPage) Assets() []morphadon.Asset {
	return []morphadon.Asset{
		morphadon.NewCSSAsset("example.css", morphadon.ScopeGlobal),
	}
}

func (c *WorldPage) Render(data morphadon.SetupData) morphadon.Renderable {
	return Div(
		g.Text("World"),
	)
}

func main() {
	app := morphadon.CreateApp()

	app.RegisterServerEndpoint(morphadon.NewServerEndpoint(http.MethodGet, "/api/foo.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"foo": "bar"}`))
	}))

	app.AssetManager().SetConfig(&morphadon.AssetManagerConfig{
		SrcDir:    "./web",
		OutputDir: "./web/public",
	})

	app.Use(morphadon.MorphadonRouter([]morphadon.Route{
		morphadon.NewRoute("hello", "/", NewHelloPage(), []morphadon.Route{
			morphadon.NewRoute("world", "/", NewWorldPage()),
		}...),
	}))

	app.RegisterMiddleware(middleware.Logger)
	if err := app.Mount(NewExampleApp()); err != nil {
		panic(err)
	}
}
