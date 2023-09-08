package morphadon

type Scope string

func (s Scope) String() string {
	return string(s)
}

const (
	ScopeGlobal    Scope = "global"
	ScopeMultiple  Scope = "multiple"
	ScopeComponent Scope = "component"
	ScopeNone      Scope = "none"
)

type Plugin interface {
	// Init is called once the plugin is registered.
	Init(*App) error
}
