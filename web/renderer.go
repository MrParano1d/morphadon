package web

import (
	"fmt"
	"io"

	g "github.com/maragudk/gomponents"
	"github.com/marlaone/engine/core"
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
	node, ok := data.(g.Node)
	if !ok {
		return fmt.Errorf("data is not a g.Node")
	}
	return node.Render(w)
}
