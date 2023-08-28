package main

import (
	"github.com/mrparano1d/morphadon"
	"github.com/mrparano1d/morphadon/web"
)

type ExampleSystem struct {
	*morphadon.DefaultSystem[*web.Context]
}

var _ morphadon.System[*web.Context] = (*ExampleSystem)(nil)

func NewExampleSystem() *ExampleSystem {
	return &ExampleSystem{
		DefaultSystem: morphadon.NewDefaultSystem[*web.Context](),
	}
}

func (s *ExampleSystem) Actions() []morphadon.Action[*web.Context] {
	return []morphadon.Action[*web.Context]{
		web.NewPageAction(web.OpHttpGet, ExamplePageScope, NewExamplePage()),
	}
}
