package engine

import (
	g "github.com/maragudk/gomponents"
)

type Page interface {
	Setup() SetupData
	Context() *Context
	SetContext(ctx *Context)
	Head() *Head
	Assets() []string
	Stylesheets() []string
	Scripts() []string
	Components() []Component
	Render(SetupData, g.Node) g.Node
	H(n interface{}) g.Node
}

type DefaultPage struct {
	ctx *Context
}

var _ Page = &DefaultPage{}

func NewDefaultPage() *DefaultPage {
	return &DefaultPage{}
}

func (l *DefaultPage) Setup() SetupData {
	return SetupData{}
}

func (l *DefaultPage) Context() *Context {
	return l.ctx
}

func (l *DefaultPage) SetContext(ctx *Context) {
	l.ctx = ctx
}

func (l *DefaultPage) Assets() []string {
	return []string{}
}

func (l *DefaultPage) Scripts() []string {
	return []string{}
}

func (l *DefaultPage) Stylesheets() []string {
	return []string{}
}

func (l *DefaultPage) Components() []Component {
	return []Component{}
}

func (l *DefaultPage) Head() *Head {
	return UseHead()
}

func (l *DefaultPage) Render(_ SetupData, _ g.Node) g.Node {
	return nil
}

func (l *DefaultPage) H(n interface{}) g.Node {
	return l.Context().H(n)
}
