package web

import (
	"fmt"
	"log"
	"reflect"

	g "github.com/maragudk/gomponents"
	"github.com/marlaone/engine/core"
)

func h(c any, ctx *Context) g.Node {
	switch v := c.(type) {
	case core.Component[*Context]:
		component := v
		component.SetContext(ctx)
		return component.Render(component.Setup()).(g.Node)
	case core.System[*Context]:
		page := v
		page.SetContext(ctx)
		return page.Render(page.Setup()).(g.Node)
	case g.Node:
		return c.(g.Node)
	case []g.Node:
		return g.Group(c.([]g.Node))
	}
	return nil
}

func MustRenderSlot(slotName string, c core.Component[*Context]) g.Node {
	node, err := RenderSlot(slotName, c)
	if err != nil {
		log.Printf("[warn] %v\n", err)
		return nil
	}
	return node
}

func RenderSlot(slotName string, c core.Component[*Context]) (g.Node, error) {
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

func resolvePageComponents(components []core.Component[*Context], pages ...core.System[*Context]) []core.Component[*Context] {
	for _, page := range pages {
		components = append(components, resolveComponents(page.Components()...)...)
	}
	return components
}

func resolveComponents(components ...core.Component[*Context]) []core.Component[*Context] {
	for _, c := range components {
		components = append(components, resolveComponents(c.Components()...)...)
	}
	return components
}
