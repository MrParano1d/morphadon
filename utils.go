package engine

import (
	"fmt"
	"reflect"

	g "github.com/maragudk/gomponents"
)

func MustRenderSlot(slotName string, c Component) g.Node {
	node, err := RenderSlot(slotName, c)
	if err != nil {
		Logger.Warn(err.Error())
		return nil
	}
	return node
}

func RenderSlot(slotName string, c Component) (g.Node, error) {
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

func h(c interface{}, ctx *Context) g.Node {
	switch c.(type) {
	case Page:
		page := c.(Page)
		page.SetContext(ctx)
		return page.Render(page.Setup(), nil)
	case Component:
		component := c.(Component)
		component.SetContext(ctx)
		return c.(Component).Render(component.Setup())
	case g.Node:
		return c.(g.Node)
	case []g.Node:
		return g.Group(c.([]g.Node))
	}
	return nil
}

func resolvePageComponents(components []Component, pages ...Page) []Component {
	for _, page := range pages {
		components = append(components, resolveComponents(page.Components()...)...)
	}
	return components
}

func resolveComponents(components ...Component) []Component {
	for _, c := range components {
		components = append(components, resolveComponents(c.Components()...)...)
	}
	return components
}
