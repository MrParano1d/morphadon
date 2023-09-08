package morphadon

import (
	"net/http"
	"sync"
)

func NewDefaultApp() *App {
	return &App{
		mutex:      &sync.RWMutex{},
		am:         NewAssetManager(),
		services:   make(map[string]any),
		middleware: make([]func(http.Handler) http.Handler, 0),
		serverEndpoints: make([]ServerEndpoint, 0),
	}
}
