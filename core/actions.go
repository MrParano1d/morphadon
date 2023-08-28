package core

type Operation uint32

func (f Operation) HasFlag(flag Operation) bool { return f&flag != 0 }
func (f *Operation) AddFlag(flag Operation)     { *f |= flag }
func (f *Operation) ClearFlag(flag Operation)   { *f &= ^flag }
func (f *Operation) ToggleFlag(flag Operation)  { *f ^= flag }

type Action[C Context] interface {
	// Init is called once the route is registered.
	Init(App[C]) error

	Operation() Operation
	Scope() Scope
	Execute(C) any

	Assets() []Asset[C]
	Components() []Component[C]

	Renderer() Renderer[C]
	SetRenderer(Renderer[C])
}

type ActionFunc[C Context] struct {
	operation Operation
	scope     Scope
	fn        func(ctx C) any
	renderer  Renderer[C]
}

var _ Action[*TodoContext] = &ActionFunc[*TodoContext]{}

func NewActionFunc[C Context](operation Operation, scope Scope, fn func(ctx C) any) *ActionFunc[C] {
	return &ActionFunc[C]{
		operation: operation,
		scope:     scope,
		fn:        fn,
		renderer:  nil,
	}
}

func (a *ActionFunc[C]) Init(app App[C]) error {
	return nil
}

func (a *ActionFunc[C]) Renderer() Renderer[C] {
	return a.renderer
}

func (a *ActionFunc[C]) SetRenderer(r Renderer[C]) {
	a.renderer = r
}

func (a *ActionFunc[C]) Operation() Operation {
	return a.operation
}

func (a *ActionFunc[C]) Scope() Scope {
	return a.scope
}

func (a *ActionFunc[C]) Execute(ctx C) any {
	return a.fn(ctx)
}

func (a *ActionFunc[C]) Assets() []Asset[C] {
	return make([]Asset[C], 0)
}

func (a *ActionFunc[C]) Components() []Component[C] {
	return make([]Component[C], 0)
}
