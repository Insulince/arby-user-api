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

var config *configuration.Config

func init() () {
	var err error
	config, err = configuration.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = mongo.InitializeDatabase(config)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() () {
	r := router.CreateRouter(config)

	c := cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowCredentials: config.Cors.AllowCredentials,
	}

	log.Printf("Server listening on port %v.\n", config.Port)
	log.Fatalln(
		http.ListenAndServe(
			":"+strconv.Itoa(config.Port),
			cors.New(c).Handler(r),
		),
	)
}
