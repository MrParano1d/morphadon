package morphadon

type Presenter[C Context] interface {
	Init(App[C]) error
	RegisterAction(Action[C]) error
	Renderer() Renderer[C]
	SetRenderer(Renderer[C])
	Start() error
}
