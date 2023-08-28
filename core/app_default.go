package core

import (
	"fmt"
	"log"
)

type DefaultApp[C Context] struct {
	presenter    Presenter[C]
	assetManager AssetManager[C]
}

var _ App[*TodoContext] = &DefaultApp[*TodoContext]{}

func NewDefaultApp[C Context]() *DefaultApp[C] {
	return &DefaultApp[C]{
		assetManager: NewAssetManagerNoop[C](),
	}
}

func (app *DefaultApp[C]) Use(plugins ...Plugin[C]) App[C] {
	for _, plugin := range plugins {
		if err := plugin.Init(app); err != nil {
			panic(err)
		}
	}
	return app
}

func (app *DefaultApp[C]) SetPresenter(p Presenter[C]) {
	if err := p.Init(app); err != nil {
		log.Fatalf("failed to init presenter %T: %v", p, err)
	}
	app.presenter = p
}

func (app *DefaultApp[C]) Presenter() Presenter[C] {
	return app.presenter
}

func (app *DefaultApp[C]) RegisterSystem(s System[C]) error {

	for _, action := range s.Actions() {

		for _, component := range action.Components() {
			if err := app.RegisterComponent(component); err != nil {
				return fmt.Errorf("failed to register component %T in action %T: %w", component, action, err)
			}
		}

		for _, asset := range action.Assets() {
			if err := app.assetManager.RegisterAsset(asset); err != nil {
				return fmt.Errorf("failed to register asset %T in action %T: %w", asset, action, err)
			}
		}

		if err := app.presenter.RegisterAction(action); err != nil {
			return fmt.Errorf("failed to register action %T in system %T: %w", action, s, err)
		}
	}

	return nil
}

func (app *DefaultApp[C]) AssetManager() AssetManager[C] {
	return app.assetManager
}

func (app *DefaultApp[C]) SetAssetManager(am AssetManager[C]) {
	app.assetManager = am
}

func (app *DefaultApp[C]) RegisterComponent(c Component[C]) error {

	for _, component := range c.Components() {
		if component == c {
			log.Printf("Component %T is trying to register itself, skipping\n", component)
			continue
		}
		if err := app.RegisterComponent(component); err != nil {
			return fmt.Errorf("failed to register component %T in component %T: %w", component, c, err)
		}
	}

	for _, asset := range c.Assets() {
		if err := app.assetManager.RegisterAsset(asset); err != nil {
			return fmt.Errorf("failed to register asset %T in component %T: %w", asset, c, err)
		}
	}

	return nil
}

func (app *DefaultApp[C]) Init() error {
	return nil
}

func (app *DefaultApp[C]) Mount() error {
	log.Println("Mounting app")

	if err := app.Init(); err != nil {
		return err
	}

	if app.assetManager != nil {
		log.Println("Building assets")
		if err := app.assetManager.Build(); err != nil {
			return err
		}
	}

	log.Println("Starting presenter")

	return app.presenter.Start()
}
