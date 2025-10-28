package main

import (
	"mockLogin/internal/adapters/rest"
	"mockLogin/internal/adapters/rest/handlers"
	"mockLogin/internal/config"
)

func main() {
	// carga configuraciones
	config.Load()
	cfg := config.GetConfig()

	// setea handlers
	handlers := handlers.NewHandlers(cfg)

	// inicia listener
	rest.NewRouter(cfg, handlers)
}
