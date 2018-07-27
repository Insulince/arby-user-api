package configuration

import (
	"os"
	"strings"
	"strconv"
	"log"
)

type Config struct {
	Port  int       `json:"port"`
	Cors  CorsConfig  `json:"cors"`
	Mongo MongoConfig `json:"mongo"`
}

type MongoConfig struct {
	ConnectionString string `json:"connectionString"`
	DatabaseName     string `json:"databaseName"`
	CollectionNames   []string `json:"collectionNames"`
}

type CorsConfig struct {
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedHeaders   []string `json:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
}

func LoadConfig() (config *Config) {
	log.Printf("Loading config from environment variables...\n")
	config = &Config{}

	if value, present := os.LookupEnv("PORT"); present {
		port, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Fatalf("Environment variable \"PORT\" is invalid! Must be an integer, but got \"%v\".\n", value)
		}
		config.Port = int(port)
	} else {
		log.Fatalf("Environment variable \"PORT\" not provided!\n")
	}

	if value, present := os.LookupEnv("MONGO_CONNECTION_STRING"); present {
		config.Mongo.ConnectionString = value
	} else {
		log.Fatalf("Environment variable \"MONGO_CONNECTION_STRING\" not provided!\n")
	}
	if value, present := os.LookupEnv("MONGO_DATABASE_NAME"); present {
		config.Mongo.DatabaseName = value
	} else {
		log.Fatalf("Environment variable \"MONGO_DATABASE_NAME\" not provided!\n")
	}

	if value, present := os.LookupEnv("CORS_ALLOWED_ORIGINS"); present {
		config.Cors.AllowedOrigins = strings.Split(value, ",")
	} else {
		log.Fatalf("Environment variable \"CORS_ALLOWED_ORIGINS\" not provided!\n")
	}
	if value, present := os.LookupEnv("CORS_ALLOWED_METHODS"); present {
		config.Cors.AllowedMethods = strings.Split(value, ",")
	} else {
		log.Fatalf("Environment variable \"CORS_ALLOWED_METHODS\" not provided!\n")
	}
	if value, present := os.LookupEnv("CORS_ALLOWED_HEADERS"); present {
		config.Cors.AllowedHeaders = strings.Split(value, ",")
	} else {
		log.Fatalf("Environment variable \"CORS_ALLOWED_HEADERS\" not provided!\n")
	}
	if value, present := os.LookupEnv("CORS_ALLOW_CREDENTIALS"); present {
		allowCredentials, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Environment variable \"CORS_ALLOW_CREDENTIALS\" is invalid! Must be one of [\"true\", \"false\"], but got \"%v\".\n", value)
		}
		config.Cors.AllowCredentials = allowCredentials
	} else {
		log.Fatalf("Environment variable \"CORS_ALLOW_CREDENTIALS\" not provided!\n")
	}

	log.Printf("Successfully loaded config.\n")
	return config
}
