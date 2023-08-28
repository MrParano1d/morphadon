package main

import (
	"github.com/mrparano1d/morphadon/core"
	"github.com/mrparano1d/morphadon/web"
)

type ExampleSystem struct {
	*core.DefaultSystem[*web.Context]
}

var _ core.System[*web.Context] = (*ExampleSystem)(nil)

func NewExampleSystem() *ExampleSystem {
	return &ExampleSystem{
		DefaultSystem: core.NewDefaultSystem[*web.Context](),
	}
}

func (s *ExampleSystem) Actions() []core.Action[*web.Context] {
	return []core.Action[*web.Context]{
		web.NewPageAction(web.OpHttpGet, ExamplePageScope, NewExamplePage()),
	}
}
