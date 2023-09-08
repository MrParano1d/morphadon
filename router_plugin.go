package morphadon

const routerServiceKey = "router"

type routerPlugin struct {
	routes []Route
}

func MorphadonRouter(routes []Route) Plugin {
	return &routerPlugin{
		routes: routes,
	}
}

func (p *routerPlugin) Init(app *App) error {
	app.RegisterService(routerServiceKey, NewRouter())
	for _, route := range p.routes {
		app.GetService(routerServiceKey).(*Router).AddRoute(route)
	}
	return nil
}
