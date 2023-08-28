package morphadon

var app any

func GetInstance[C Context]() App[C] {
	if app == nil {
		panic("App not initialized")
	}
	return app.(App[C])
}

func CreateApp[C Context]() App[C] {
	app = NewDefaultApp[C]()
	return app.(App[C])
}

func DefineComponent[C Context](c Component[C]) Component[C] {
	app := GetInstance[C]()
	app.RegisterComponent(c)
	return c
}
