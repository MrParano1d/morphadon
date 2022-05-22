package engine

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return NewRouterWithRoutes([]*Route{})
}

func NewRouterWithRoutes(routes []*Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) GetRoutes() []*Route {
	return r.routes
}

func (r *Router) AddRoute(route *Route) {
	r.routes = append(r.routes, route)
}

func (app *Marla) Router(r *Router) {
	app.router = r
}

func resolveChildRoutes(routePath string, routeName string, pages []Page, route *Route) (string, string, []Page) {
	routePath += route.Route
	routeName += route.Name
	pages = append(pages, route.Page)

	for _, r := range route.Children {
		routePath, routeName, pages = resolveChildRoutes(routePath, routeName, pages, r)
	}

	return routePath, routeName, pages
}

func resolveRoutes(routePath string, routeName string, pages []Page, route *Route) map[string]map[string][]Page {
	routeMap := map[string]map[string][]Page{}

	pages = append(pages, route.Page)

	for _, r := range route.Children {
		routePath, routeName, pages := resolveChildRoutes(routePath, routeName, pages, r)
		routeMap[routePath] = map[string][]Page{
			routeName: pages,
		}
	}

	return routeMap
}

func flattenPages(routes map[string]map[string][]Page) []Page {
	var pages []Page
	for _, r := range routes {
		for _, p := range r {
			pages = append(pages, p...)
		}
	}
	return pages
}
