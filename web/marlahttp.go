package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marlaone/morphadon/core"
)

type MarlaHttpPresenter struct {
	app core.App[*Context]

	router   *chi.Mux
	renderer core.Renderer[*Context]
}

var _ core.Presenter[*Context] = &MarlaHttpPresenter{}

func NewHttpPresenter() *MarlaHttpPresenter {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/public", filesDir)

	return &MarlaHttpPresenter{
		router:   r,
		renderer: NewBytesRenderer(),
	}
}

func (p *MarlaHttpPresenter) Init(app core.App[*Context]) error {
	p.app = app
	return nil
}

func (p *MarlaHttpPresenter) Renderer() core.Renderer[*Context] {
	return p.renderer
}

func (p *MarlaHttpPresenter) SetRenderer(r core.Renderer[*Context]) {
	p.renderer = r
}

func (p *MarlaHttpPresenter) actionHandler(action core.Action[*Context]) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var renderer core.Renderer[*Context]
		if action.Renderer() != nil {
			renderer = action.Renderer()
		} else {
			renderer = p.renderer
		}
		ctx := NewContext()

		if err := renderer.Render(action.Execute(ctx), rw); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
	}
}

func (p *MarlaHttpPresenter) RegisterAction(action core.Action[*Context]) error {

	switch action.Operation() {
	case OpHttpGet:
		p.router.Get(action.Scope().String(), p.actionHandler(action))
	case OpHttpPost:
		p.router.Post(action.Scope().String(), p.actionHandler(action))
	case OpHttpPut:
		p.router.Put(action.Scope().String(), p.actionHandler(action))
	case OpHttpDelete:
		p.router.Delete(action.Scope().String(), p.actionHandler(action))
	case OpHttpPatch:
		p.router.Patch(action.Scope().String(), p.actionHandler(action))
	case OpHttpHead:
		p.router.Head(action.Scope().String(), p.actionHandler(action))
	case OpHttpOptions:
		p.router.Options(action.Scope().String(), p.actionHandler(action))
	default:
		return fmt.Errorf("Unknown operation: %d", action.Operation())

	}
	return nil
}

func (p *MarlaHttpPresenter) Start() error {
	log.Println("Starting HTTP server on port 8080")
	return http.ListenAndServe(":8080", p.router)
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
