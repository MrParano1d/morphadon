package core

import "context"

type Context interface {
	Context() context.Context
}

type TodoContext struct{}

var _ Context = (*TodoContext)(nil)

func (c *TodoContext) Context() context.Context {
	return context.TODO()
}
