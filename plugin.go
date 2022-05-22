package engine

import g "github.com/maragudk/gomponents"

type Plugin interface {
	Install(app *Marla, options interface{}) error
}

type SourcePlugin interface {
	Plugin
	RequestRender(filter interface{}, render func(data interface{}) g.Node) g.Node
	RequestRaw(filter interface{}, raw func(data interface{}) interface{}) interface{}
}

func EmptyRequest(render func(data interface{}) g.Node) g.Node {
	return render([]interface{}{})
}
