package morphadon

func If(condition bool, trueRenderable, falseRenderable Renderable) Renderable {
	if condition {
		return trueRenderable
	}
	return falseRenderable
}
