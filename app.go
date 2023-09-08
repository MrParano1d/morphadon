package morphadon

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

var app *App

func GetInstance() *App {
	if app == nil {
		panic("App not initialized")
	}
	return app
}

func CreateApp() *App {
	app = NewDefaultApp()
	return app
}

type HttpMethod string

const (
	ANY    HttpMethod = "ANY"
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
	HEAD   HttpMethod = "HEAD"
)

type ServerEndpoint struct {
	Method  HttpMethod
	Path    string
	Handler http.HandlerFunc
}

func NewServerEndpoint(method HttpMethod, path string, handler http.HandlerFunc) ServerEndpoint {
	return ServerEndpoint{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

type App struct {
	mutex *sync.RWMutex

	am AssetManager

	middleware []func(http.Handler) http.Handler

	serverEndpoints []ServerEndpoint

	services map[string]any
	// RegisterPlugin registers a plugin.
}

func (a *App) Use(plugins ...Plugin) *App {
	for _, plugin := range plugins {
		plugin.Init(a)
	}
	return a
}

func (a *App) RegisterService(name string, service any) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.services[name] = service
}

func (a *App) RegisterServerEndpoint(endpoint ServerEndpoint) {
	a.serverEndpoints = append(a.serverEndpoints, endpoint)
}

func (a *App) GetService(name string) any {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.services[name]
}

func (a *App) AssetManager() AssetManager {
	return a.am
}
func (a *App) SetAssetManager(am AssetManager) {
	a.am = am
}

func (a *App) RegisterMiddleware(middleware ...func(http.Handler) http.Handler) {
	a.middleware = append(a.middleware, middleware...)
}

func (a *App) registerComponent(component Component) {

	for _, c := range component.Components() {
		if c == component {
			log.Printf("Component %T already registered", c)
			continue
		}
		a.registerComponent(c)
	}

	for _, asset := range component.Assets() {
		if err := a.am.RegisterAsset(asset); err != nil {
			log.Fatalf("Error registering asset %T in component %T: %w", asset, component, err)
		}
	}
}

func (a *App) Mount(component Component) error {

	r := chi.NewRouter()
	r.Use(a.middleware...)

	router := a.GetService(routerServiceKey).(*Router)

	for _, endpoint := range a.serverEndpoints {
		r.MethodFunc(string(endpoint.Method), endpoint.Path, endpoint.Handler)
	}

	for _, route := range router.Routes() {

		a.registerComponent(route.page)

		r.Group(func(r chi.Router) {
			r.Use(router.parentMiddlewares(&route)...)
			r.Use(route.page.Middlewares()...)
			r.Get(route.Path(), func(w http.ResponseWriter, r *http.Request) {

				ctx := NewContext(
					ContextWithRequest(r),
					ContextWithContext(r.Context()),
				)

				ProvideScope(ctx, Scope(r.URL.Path))

				w.Header().Set("Content-Type", "text/html")
				data := component.Setup()
				component.SetContext(ctx)
				renderer := component.Render(data)
				if err := renderer.Render(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})
		})
	}

	FileServer(r, "/public", http.Dir(a.am.Config().OutputDir))

	if err := a.am.Build(); err != nil {
		return fmt.Errorf("Error building assets: %w", err)
	}

	log.Printf("Listening on port 3000")
	return http.ListenAndServe(":3000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
