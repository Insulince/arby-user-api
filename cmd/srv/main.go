package main

import (
	"arby-user-api/pkg/configuration"
	"log"
	"arby-user-api/pkg/mongo"
	"net/http"
	"github.com/rs/cors"
	"strconv"
	"arby-user-api/pkg/router"
)

func main() () {
	config := configuration.LoadConfig()

	err := mongo.InitializeDatabase(config)
	if err != nil {
		log.Fatalln(err)
	}

	r := router.CreateRouter()

	log.Printf("Server listening on port %v.\n", config.Port)
	log.Fatalln(
		http.ListenAndServe(
			":"+strconv.Itoa(config.Port),
			cors.New(
				cors.Options{
					AllowedOrigins:   config.Cors.AllowedOrigins,
					AllowedHeaders:   config.Cors.AllowedHeaders,
					AllowedMethods:   config.Cors.AllowedMethods,
					AllowCredentials: config.Cors.AllowCredentials,
				},
			).Handler(r),
		),
	)
}
