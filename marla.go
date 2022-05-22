package engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/marlaone/engine/middleware"
	"path"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	"github.com/maragudk/gomponents/html"
	"go.uber.org/zap"
)

type Marla struct {
	Components []Component
	router     *Router
	HTML       func(ctx *Context, node interface{}, head *Head, scripts []string, stylesheets []string) g.Node
	ctx        context.Context
	logger     *zap.Logger
}

var App *Marla

func CreateApp() *Marla {
	App = &Marla{
		ctx:        context.Background(),
		Components: []Component{},
		logger:     Logger,
		HTML: func(ctx *Context, node interface{}, head *Head, scripts []string, stylesheets []string) g.Node {
			title := "Marla.ONE Dev"
			if head.Title != "" {
				title = head.Title
			}
			language, ok := head.HtmlAttrs["lang"]
			if !ok {
				language = "en"
			}
			return c.HTML5(c.HTML5Props{
				Title:    title,
				Language: language,
				Head: []g.Node{
					g.Group(g.Map(len(head.Meta), func(i int) g.Node {
						attributes := []g.Node{}
						for k, v := range head.Meta[i] {
							attributes = append(attributes, g.Attr(k, v))
						}
						return html.Meta(g.Group(attributes))
					})),
					g.Group(g.Map(len(stylesheets), func(i int) g.Node {
						return html.Link(html.Rel("stylesheet"), html.Href(stylesheets[i]))
					})),
					g.Group(g.Map(len(head.Link), func(i int) g.Node {
						attributes := []g.Node{}
						for k, v := range head.Link[i] {
							attributes = append(attributes, g.Attr(k, v))
						}
						return html.Link(g.Group(attributes))
					})),
					g.If(len(head.Base) > 0, head.baseNode()),
				},
				Body: []g.Node{
					head.bodyAttributes(),
					h(node, ctx),
					g.Group(g.Map(len(scripts), func(i int) g.Node {
						return html.Script(html.Src(scripts[i]))
					})),
					g.Group(g.Map(len(head.Scripts), func(i int) g.Node {
						var attributes []g.Node
						for k, v := range head.Scripts[i] {
							attributes = append(attributes, g.Attr(k, v))
						}
						return html.Script(g.Group(attributes))
					})),
				},
			})
		},
	}
	return App
}

func (app *Marla) Context() context.Context {
	return app.ctx
}

func (app *Marla) Logger() *zap.Logger {
	return app.logger
}

func (app *Marla) Use(plugin Plugin, options interface{}) {
	err := plugin.Install(app, options)
	if err != nil {
		panic(fmt.Errorf("insalling plugin %T failed: %v", plugin, err))
	}
}

func (app *Marla) Mount() {

	fibr := fiber.New()

	fibr.Use(middleware.NewLogger(middleware.WithLogger(app.Logger())).Handler())

	fibr.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	if err := app.Build(); err != nil {
		panic(fmt.Errorf("assets build failed: %v", err))
	}

	for _, route := range app.router.routes {
		for routePath, routeMap := range resolveRoutes("", "", []Page{}, route) {
			routePath = path.Clean(routePath)
			for routeName, pages := range routeMap {
				app.createRoute(routePath, pages, fibr, routeName)
			}
		}
	}

	fibr.Static("/static", "./public")

	if err := fibr.Listen(":7292"); err != nil {
		app.Logger().Fatal("failed to mount engine", zap.Error(err))
	}
}

func (app *Marla) createRoute(routePath string, pages []Page, fibr *fiber.App, routeName string) {
	fibr.All(routePath, func(c *fiber.Ctx) error {
		ctx := &Context{
			Fiber:     c,
			RouteName: routeName,
		}
		head := UseHead()
		for _, p := range pages {
			head = mergeHeads(p.Head(), head)
		}
		buf := bytes.NewBufferString("")
		if err := app.HTML(ctx, renderPages(pages, ctx), head, app.Scripts(routeName), app.Stylesheets(routeName)).Render(buf); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendString(buf.String())
	})
}

func renderPages(pages []Page, ctx *Context) g.Node {
	var node g.Node
	for i := len(pages) - 1; i >= 0; i-- {
		pages[i].SetContext(ctx)
		node = pages[i].Render(pages[i].Setup(), node)
	}
	return node
}
