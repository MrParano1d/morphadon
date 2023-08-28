package morphadon

type SetupData map[string]interface{}

type Component[C Context] interface {
	Setup() SetupData
	Assets() []Asset[C]
	Slots() Slots
	Props() Properties
	Context() C
	SetContext(ctx C)
	Components() []Component[C]
	Render(data SetupData) any
}

type DefaultComponent[C Context] struct {
	ctx   C
	props Properties
	slots Slots
}

var _ Component[*TodoContext] = &DefaultComponent[*TodoContext]{}

func NewDefaultComponent[C Context]() *DefaultComponent[C] {
	return NewDefaultComponentWithPropsAndSlots[C](Properties{}, Slots{})
}

func NewDefaultComponentWithProps[C Context](props Properties) *DefaultComponent[C] {
	return NewDefaultComponentWithPropsAndSlots[C](props, Slots{})
}

func NewDefaultComponentWithSlots[C Context](slots Slots) *DefaultComponent[C] {
	return NewDefaultComponentWithPropsAndSlots[C](Properties{}, slots)
}

func NewDefaultComponentWithPropsAndSlots[C Context](props Properties, slots Slots) *DefaultComponent[C] {
	return &DefaultComponent[C]{
		props: props,
		slots: slots,
	}
}

func (c *DefaultComponent[C]) Setup() SetupData {
	return SetupData{}
}

func (c *DefaultComponent[C]) Context() C {
	return c.ctx
}

func (c *DefaultComponent[C]) SetContext(ctx C) {
	c.ctx = ctx
}

func (c *DefaultComponent[C]) Slots() Slots {
	return c.slots
}

func (c *DefaultComponent[C]) Props() Properties {
	return c.props
}

func (c *DefaultComponent[C]) Assets() []Asset[C] {
	return make([]Asset[C], 0)
}

func (c *DefaultComponent[C]) Components() []Component[C] {
	return make([]Component[C], 0)
}

func (c *DefaultComponent[C]) Render(data SetupData) any {
	return nil
}
