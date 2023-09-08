package morphadon

type useRoute struct {
	CurrentRoute *plainRoute
}

func UseRoute(ctx *Context) *useRoute {
	router := GetInstance().GetService(routerServiceKey).(*Router)

	return &useRoute{
		CurrentRoute: router.CurrentRoute(ctx),
	}
}
