package web

import (
	"fmt"
	"io"
	"log"

	g "github.com/maragudk/gomponents"
	"github.com/marlaone/morphadon/core"
)

type WebRenderer struct {
}

var _ core.Renderer[*Context] = &WebRenderer{}

func NewWebRenderer() *WebRenderer {
	return &WebRenderer{}
}

func (r *WebRenderer) Init(core.App[*Context]) error {
	return nil
}

func (r *WebRenderer) Render(data any, w io.Writer) error {
	if data == nil {
		log.Println("[warn] WebRenderer.Render called with nil data")
		return nil
	}
	node, ok := data.(g.Node)
	if !ok {
		return fmt.Errorf("data is not a g.Node, but a %T", data)
	}
	return node.Render(w)
}
