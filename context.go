package engine

import (
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
)

type Context struct {
	Fiber     *fiber.Ctx
	RouteName string
}

func (ctx *Context) H(c interface{}) g.Node {
	return h(c, ctx)
}
