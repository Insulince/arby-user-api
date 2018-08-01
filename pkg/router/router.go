package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"arby-user-api/pkg/api"
)

type Router struct {
	*mux.Router
}

func CreateRouter() (router *Router) {
	router = &Router{
		mux.NewRouter().StrictSlash(true),
	}

	router.HandleFunc("/", api.Home).Methods("GET")

	router.HandleFunc("/health", api.HealthCheck).Methods("GET")

	router.HandleFunc("/register", api.Register).Methods("GET")
	router.HandleFunc("/login", api.Login).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(api.NotFound)

	return router
}
