package rest

import (
	"fmt"
	"net/http"

	"mockLogin/internal/adapters/rest/handlers"
	"mockLogin/internal/config"

	"github.com/gorilla/mux"
)

type Server struct {
	Host      string
	Port      string
	SecretKey string
}

func NewRouter(configs *config.Config, handlers *handlers.Handlers) {
	s := Server{
		Host:      configs.Server.Host,
		Port:      configs.Server.Port,
		SecretKey: configs.Server.SecretKey,
	}

	router := mux.NewRouter()

	// asigna handlers dinámicamente según correspondan hallan sido creados
	for method, routes := range handlers.Handlers {
		for endpoint, handler := range routes {
			router.HandleFunc(endpoint, handler).Methods(method)
			fmt.Printf("HandleFunc agregado, método %v, endpoint %v\n", method, endpoint)
		}
	}

	http.Handle("/", router)

	fmt.Printf("Esto es config.Puestos: %+v \n", configs.Puestos)

	// arranca el listener del server
	fmt.Printf("\nIniciando Moca Mock en %s:%s ...\n", s.Host, s.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", s.Host, s.Port), router); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
