package web

import (
	"github.com/marlaone/engine/core"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	html "github.com/maragudk/gomponents/html"
)

type HTMLComponent struct {
	*core.DefaultComponent[*Context]
}

var _ core.Component[*Context] = (*HTMLComponent)(nil)

func HTML(children ...g.Node) *HTMLComponent {
	return &HTMLComponent{
		DefaultComponent: core.NewDefaultComponentWithSlots[*Context](core.Slots{
			"default": children,
		}),
	}
}

func (h *HTMLComponent) Render(data core.SetupData) any {

	var scripts []string
	for _, asset := range h.Context().Assets() {
		if asset.Type() != core.AssetTypeJS {
			continue
		}
		if asset.Scope() != core.ScopeGlobal {
			continue
		}
		scripts = append(scripts, asset.TargetPath())
	}

	var styles []string
	for _, asset := range h.Context().Assets() {
		if asset.Type() != core.AssetTypeCSS {
			continue
		}
		if asset.Scope() != core.ScopeGlobal {
			continue
		}
		styles = append(styles, asset.TargetPath())
	}

	var title string
	if h.Context().Title != "" {
		title = h.Context().Title
	} else {
		title = "Marla//Engine powered website"
	}

	return c.HTML5(
		c.HTML5Props{
			Title:    title,
			Language: h.Context().Language,
			Head: []g.Node{
				g.Group(g.Map(h.Context().Meta, func(meta map[string]string) g.Node {
					attributes := make([]g.Node, 0, len(meta))
					for k, v := range meta {
						attributes = append(attributes, g.Attr(k, v))
					}
					return html.Meta(
						g.Group(attributes),
					)
				})),
				g.Group(g.Map(scripts, func(script string) g.Node {
					return html.Script(
						html.Src(script),
						html.Defer(),
					)
				})),
				g.Group(g.Map(styles, func(style string) g.Node {
					return html.Link(
						html.Rel("stylesheet"),
						html.Href(style),
					)
				})),
			},
			Body: []g.Node{
				h.Context().BodyAttrs,
				h.Context().H(
					MustRenderSlot("default", h),
				),
			},
		},
	)
}
