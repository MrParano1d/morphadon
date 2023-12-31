package morphadon

type Properties map[string]any

func (props Properties) Has(propName string) bool {
	_, ok := props[propName]
	return ok
}

func (props Properties) Get(propName string) any {
	return props[propName]
}

func (props Properties) GetStr(propName string) string {
	return PropStr(propName, props)
}

func (props Properties) GetBool(propName string) bool {
	return PropBool(propName, props)
}

func (props Properties) GetInt(propName string) int {
	prop, ok := props[propName]
	if !ok {
		return 0
	}
	return prop.(int)
}

func (props Properties) GetFloat(propName string) float64 {
	prop, ok := props[propName]
	if !ok {
		return 0
	}
	return prop.(float64)
}

func PropStr(propName string, props Properties) string {
	return PropStrWithDefault(propName, props, "")
}

func PropStrWithDefault(propName string, props Properties, defaultValue string) string {
	str, ok := props[propName].(string)
	if !ok {
		return defaultValue
	}
	return str
}

func PropBool(propName string, props Properties) bool {
	return PropBoolWithDefault(propName, props, false)
}

func PropBoolWithDefault(propName string, props Properties, defaultValue bool) bool {
	b, ok := props[propName].(bool)
	if !ok {
		return defaultValue
	}
	return b
}
