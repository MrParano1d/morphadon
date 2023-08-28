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

	return c.HTML5(
		c.HTML5Props{
			Title:    "Marla One",
			Language: "en",
			Head: []g.Node{
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
				h.Context().H(
					MustRenderSlot("default", h),
				),
			},
		},
	)
}
