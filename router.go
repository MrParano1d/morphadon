package morphadon

import (
	"regexp"
	"slices"

	"net/url"
)

type Router struct {
	routes []plainRoute
}

func NewRouter() *Router {
	return &Router{
		routes: make([]plainRoute, 0),
	}
}

func (r *Router) AddRoute(route Route) {
	r.routes = append(r.routes, flattenRoute(route, nil)...)
}

func (r *Router) Routes() []plainRoute {
	return r.routes
}

func (r *Router) findRoute(parentPath string, path string) *plainRoute {
	path, err := url.JoinPath(parentPath, path)
	if err != nil {
		panic(err)
	}

	routes := make([]plainRoute, len(r.routes))
	copy(routes, r.routes)

	slices.Reverse(routes)

	for _, route := range routes {
		if route.pathRegexp.MatchString(path) {
			return &route
		}
	}
	return nil
}

func (r *Router) CurrentRoute(ctx *Context) *plainRoute {
	return r.findRoute("/", ctx.Req.URL.Path)
}

func (r *Router) parentMiddlewares(route *plainRoute) []Middleware {
	if route.parent == nil {
		return make([]Middleware, 0)
	}
	return append(route.parent.page.Middlewares(), r.parentMiddlewares(route.parent)...)
}

type plainRoute struct {
	parent     *plainRoute
	name       string
	path       string
	pathRegexp *regexp.Regexp
	page       Page
}

func flattenRoutes(routes []Route) []plainRoute {
	result := make([]plainRoute, 0)
	for _, route := range routes {
		result = append(result, flattenRoute(route, nil)...)
	}
	return result
}

func flattenRoute(route Route, parent *plainRoute) []plainRoute {
	var err error
	result := make([]plainRoute, 0, 1+len(route.children))
	path := route.path
	if parent != nil {
		path, err = url.JoinPath(parent.path, route.path)
		if err != nil {
			panic(err)
		}
	}

	path, err = url.QueryUnescape(path)
	if err != nil {
		panic(err)
	}

	result = append(result, plainRoute{
		parent:     parent,
		name:       route.name,
		pathRegexp: regexp.MustCompile(urlParamRegexp.ReplaceAllString(path, "([^/]+)") + "$"),
		path:       path,
		page:       route.page,
	})
	resultLen := len(result)
	for _, child := range route.children {
		result = append(result, flattenRoute(child, &result[resultLen-1])...)
	}
	return result
}

type Route struct {
	name string
	path string
	page Page

	children []Route
}

var urlParamRegexp = regexp.MustCompile(`{([^}]+)}`)

func NewRoute(name string, path string, page Page, children ...Route) Route {
	return Route{
		name:     name,
		path:     path,
		page:     page,
		children: children,
	}
}

func (r *plainRoute) Parent() *plainRoute {
	return r.parent
}

func (r *plainRoute) Name() string {
	return r.name
}

func (r *plainRoute) Path() string {
	return r.path
}

func (r *plainRoute) Page() Page {
	return r.page
}
