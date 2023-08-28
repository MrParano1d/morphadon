package web

import (
	"fmt"
	"log"
	"reflect"

	g "github.com/maragudk/gomponents"
	"github.com/mrparano1d/morphadon"
)

func h(c any, ctx *Context) g.Node {
	switch v := c.(type) {
	case morphadon.Component[*Context]:
		component := v
		component.SetContext(ctx)
		return component.Render(component.Setup()).(g.Node)
	case g.Node:
		return c.(g.Node)
	case []g.Node:
		return g.Group(c.([]g.Node))
	}
	return nil
}

func MustRenderSlot(slotName string, c morphadon.Component[*Context]) g.Node {
	node, err := RenderSlot(slotName, c)
	if err != nil {
		log.Printf("[warn] %v\n", err)
		return nil
	}
	return node
}

func RenderSlot(slotName string, c morphadon.Component[*Context]) (g.Node, error) {
	slot, ok := c.Slots()[slotName]
	if !ok {
		return nil, nil
	}
	rendered := h(slot, c.Context())
	if rendered != nil {
		return rendered, nil
	}
	return nil, fmt.Errorf("invalid slot type: %v", reflect.ValueOf(slot).Type())
}
