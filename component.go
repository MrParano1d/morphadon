package morphadon

import (
	"io"

	g "github.com/maragudk/gomponents"
)

type Renderable interface {
	Render(w io.Writer) error
}

type SetupData map[string]interface{}

type Component interface {
	Setup() SetupData
	Context() *Context
	SetContext(*Context)
	Assets() []Asset
	Slots() Slots
	Props() Properties
	Components() []Component
	Render(data SetupData) Renderable
}

type DefaultComponent struct {
	props Properties
	slots Slots

	ctx *Context
}

var _ Component = &DefaultComponent{}

func NewDefaultComponent() *DefaultComponent {
	return NewDefaultComponentWithPropsAndSlots(Properties{}, Slots{})
}

func NewDefaultComponentWithProps(props Properties) *DefaultComponent {
	return NewDefaultComponentWithPropsAndSlots(props, Slots{})
}

func NewDefaultComponentWithSlots(slots Slots) *DefaultComponent {
	return NewDefaultComponentWithPropsAndSlots(Properties{}, slots)
}

func NewDefaultComponentWithPropsAndSlots(props Properties, slots Slots) *DefaultComponent {
	return &DefaultComponent{
		props: props,
		slots: slots,
	}
}

func (c *DefaultComponent) Context() *Context {
	return c.ctx
}

func (c *DefaultComponent) SetContext(ctx *Context) {
	c.ctx = ctx
}

func (c *DefaultComponent) Setup() SetupData {
	return SetupData{}
}

func (c *DefaultComponent) Slots() Slots {
	return c.slots
}

func (c *DefaultComponent) Props() Properties {
	return c.props
}

func (c *DefaultComponent) Assets() []Asset {
	return make([]Asset, 0)
}

func (c *DefaultComponent) Components() []Component {
	return make([]Component, 0)
}

func (c *DefaultComponent) Render(data SetupData) Renderable {
	return g.Text("")
}
