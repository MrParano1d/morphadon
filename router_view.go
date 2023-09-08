package morphadon

import (
	"context"

	g "github.com/maragudk/gomponents"
)

type RouterViewContextKey string

const lastRouteKey RouterViewContextKey = "lastRoute"

type routerView struct {
	*DefaultComponent
}

var _ Component = (*routerView)(nil)

func RouterView(ctx *Context) Renderable {
	rv := NewRouterView()
	rv.SetContext(ctx)
	return rv.Render(SetupData{})
}

func NewRouterView() *routerView {
	return &routerView{
		DefaultComponent: NewDefaultComponent(),
	}
}

func (v *routerView) Render(data SetupData) Renderable {

	currentRoute := UseRoute(v.Context()).CurrentRoute

	if currentRoute == nil {
		return g.Text("")
	}

	lastRoute, ok := v.Context().Context().Value(lastRouteKey).(*plainRoute)
	if !ok {
		lastRoute = nil
	}
	nextRoute := v.searchNext(currentRoute, lastRoute)

	if lastRoute == nextRoute {
		return g.Text("")
	}

	v.Context().SetContext(context.WithValue(v.Context().Context(), lastRouteKey, nextRoute))

	if nextRoute != nil {
		data := nextRoute.page.Setup()
		nextRoute.page.SetContext(v.Context())
		return nextRoute.page.Render(data)
	}

	return g.Text("")
}

func (v *routerView) searchNext(route *plainRoute, parent *plainRoute) *plainRoute {
	if route.parent == parent {
		return route
	}
	if route.parent == nil {
		return nil
	}
	return v.searchNext(route.parent, parent)
}
