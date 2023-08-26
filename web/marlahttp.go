package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marlaone/engine/core"
)

type MarlaHttpPresenter struct {
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

	return &MarlaHttpPresenter{
		router:   r,
		renderer: NewBytesRenderer(),
	}
}

func (p *MarlaHttpPresenter) Init(app core.App[*Context]) error {
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
