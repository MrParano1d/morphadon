package core

import "io"

type Scope string

func (s Scope) String() string {
	return string(s)
}

const (
	ScopeGlobal Scope = "global"
)

type System[C Context] interface {
	Init(App[C]) error

	Context() C
	SetContext(C)

	Actions() []Action[C]
	Components() []Component[C]
	Assets() []Asset[C]


	Setup() SetupData
	Render(data SetupData) any
}

type Plugin[C Context] interface {
	// Init is called once the plugin is registered.
	Init(App[C]) error
}

type Renderer[C Context] interface {
	Init(App[C]) error
	Render(any, io.Writer) error
}

type App[C Context] interface {
	// RegisterPlugin registers a plugin.
	Use(...Plugin[C]) App[C]

	Presenter() Presenter[C]

	SetPresenter(Presenter[C])

	RegisterSystem(System[C]) error

	AssetManager() AssetManager[C]
	SetAssetManager(AssetManager[C])

	// RegisterComponent registers a component.
	RegisterComponent(Component[C]) error

	Mount() error
}
