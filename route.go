package engine

type Route struct {
	Name     string
	Route    string
	Page     Page
	Children []*Route
}

func NewRoute(name string, route string, p Page) *Route {
	return NewRouteWithChildrenAndName(name, route, p, []*Route{})
}

func NewRouteWithChildrenAndName(name string, route string, p Page, children []*Route) *Route {
	return &Route{
		Name:     name,
		Route:    route,
		Page:     p,
		Children: children,
	}
}

func NewRouteWithChildren(route string, p Page, children []*Route) *Route {
	return NewRouteWithChildrenAndName("", route, p, children)
}
