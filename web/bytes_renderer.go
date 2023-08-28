package web

import (
	"fmt"
	"io"

	"github.com/marlaone/engine/core"
)

type BytesRenderer struct {
}

var _ core.Renderer[*Context] = &BytesRenderer{}

func NewBytesRenderer() *BytesRenderer {
	return &BytesRenderer{}
}

func (r *BytesRenderer) Init(app core.App[*Context]) error {
	return nil
}

func (r *BytesRenderer) Render(data any, w io.Writer) error {
	bs, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("Invalid action result type: %T", data)
	}
	_, err := w.Write(bs)
	return err
}
