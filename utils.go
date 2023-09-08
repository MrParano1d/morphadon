package morphadon

import (
	"fmt"
	"log"
	"reflect"

	g "github.com/maragudk/gomponents"
)

func h(c any, ctx *Context) g.Node {
	switch v := c.(type) {
	case Renderable:
		return v.(g.Node)
	case Component:
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

func MustRenderSlot(slotName string, c Component) g.Node {
	node, err := RenderSlot(slotName, c)
	if err != nil {
		log.Printf("[warn] %v\n", err)
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

func ImageSrc(src string) string {
	app := GetInstance()

	webAssetManager, ok := app.AssetManager().(*WebAssetManager)
	if !ok {
		panic("asset manager is not a web.AssetManager")
	}

	return webAssetManager.AssetPathToBuiltPath(src)
}
