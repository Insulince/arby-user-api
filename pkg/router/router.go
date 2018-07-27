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

	router.HandleFunc("/user/register", api.Register).Methods("POST")
	router.HandleFunc("/user/login", api.Login).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(api.NotFound)

	return router
}
