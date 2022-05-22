package engine

import (
	g "github.com/maragudk/gomponents"
)

type SetupData map[string]interface{}

type Component interface {
	Setup() SetupData
	Assets() []string
	Stylesheets() []string
	Scripts() []string
	Slots() Slots
	Props() Properties
	Context() *Context
	SetContext(ctx *Context)
	Components() []Component
	Render(data SetupData) g.Node
	H(n interface{}) g.Node
}

func DefineComponent(c Component) Component {
	App.Components = append(App.Components, c)
	return c
}

type DefaultComponent struct {
	ctx   *Context
	props Properties
	slots Slots
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

func (c *DefaultComponent) Setup() SetupData {
	return SetupData{}
}

func (c *DefaultComponent) Context() *Context {
	return c.ctx
}

func (c *DefaultComponent) SetContext(ctx *Context) {
	c.ctx = ctx
}

func (c *DefaultComponent) Slots() Slots {
	return c.slots
}

func (c *DefaultComponent) Props() Properties {
	return c.props
}

func (c *DefaultComponent) Assets() []string {
	return []string{}
}

func (c *DefaultComponent) Scripts() []string {
	return []string{}
}

func (c *DefaultComponent) Stylesheets() []string {
	return []string{}
}

func (c *DefaultComponent) Components() []Component {
	return []Component{}
}

func (c *DefaultComponent) Render(data SetupData) g.Node {
	return nil
}

func (c *DefaultComponent) H(n interface{}) g.Node {
	return c.Context().H(n)
}
