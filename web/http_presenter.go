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
	"github.com/mrparano1d/morphadon"
)

type MarlaHttpPresenterOption func(*MarlaHttpPresenter)

func HttpPresenterWithFilesDir(filesDir http.Dir) MarlaHttpPresenterOption {
	return func(p *MarlaHttpPresenter) {
		p.filesDir = filesDir
	}
}

type MarlaHttpPresenter struct {
	app morphadon.App[*Context]

	filesDir http.Dir

	router   *chi.Mux
	renderer morphadon.Renderer[*Context]
}

var _ morphadon.Presenter[*Context] = &MarlaHttpPresenter{}

func NewHttpPresenter(opts ...MarlaHttpPresenterOption) *MarlaHttpPresenter {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))

	presenter := &MarlaHttpPresenter{
		router:   r,
		renderer: NewBytesRenderer(),
		filesDir: filesDir,
	}

	for _, opt := range opts {
		opt(presenter)
	}

	FileServer(r, "/public", presenter.filesDir)

	return presenter
}

func (p *MarlaHttpPresenter) Init(app morphadon.App[*Context]) error {
	p.app = app
	return nil
}

func (p *MarlaHttpPresenter) Renderer() morphadon.Renderer[*Context] {
	return p.renderer
}

func (p *MarlaHttpPresenter) SetRenderer(r morphadon.Renderer[*Context]) {
	p.renderer = r
}

func (p *MarlaHttpPresenter) actionHandler(action morphadon.Action[*Context]) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var renderer morphadon.Renderer[*Context]
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

func (p *MarlaHttpPresenter) RegisterAction(action morphadon.Action[*Context]) error {

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
