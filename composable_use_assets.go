package morphadon

type useAssetsInstance struct {
}

func UseAssets() *useAssetsInstance {
	return &useAssetsInstance{}
}

func (i *useAssetsInstance) All() []Asset {
	return GetInstance().AssetManager().Assets()
}

func (i *useAssetsInstance) Scoped(scope Scope) []Asset {
	return GetInstance().AssetManager().ScopeAssets(scope)
}
