package morphadon

import (
	"sync"
)

func NewDefaultApp() *App {
	return &App{
		mutex:           &sync.RWMutex{},
		am:              NewAssetManager(),
		services:        make(map[string]any),
		middleware:      make([]Middleware, 0),
		serverEndpoints: make([]ServerEndpoint, 0),
	}
}
