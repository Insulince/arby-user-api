package api

import (
	"arby-user-api/pkg/models"
	"net/http"
	"arby-user-api/pkg/models/responses"
	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
	"arby-user-api/pkg/mongo"
	"encoding/base64"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) () {
	_, aw := models.NewApiCommunication(r, w)
	aw.Respond(&responses.Message{Message: "Welcome!"}, http.StatusOK)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) () {
	_, aw := models.NewApiCommunication(r, w)
	aw.Respond(&responses.Message{Message: "OK"}, http.StatusOK)
}

func NotFound(w http.ResponseWriter, r *http.Request) () {
	_, aw := models.NewApiCommunication(r, w)
	aw.Respond(&responses.Error{Error: "Unsupported URL provided."}, http.StatusNotFound)
}

func Register(w http.ResponseWriter, r *http.Request) () {
	ar, aw := models.NewApiCommunication(r, w)

	authorizationHeader := ar.GetHeader("Authorization")
	if authorizationHeader == "" {
		log.Println("Empty/absent \"Authorization\" header value.")
		aw.Respond(responses.Error{Error: "Empty/absent \"Authorization\" header value."}, http.StatusBadRequest)
		return
	}
	if strings.Index(authorizationHeader, "Basic ") != 0 {
		log.Println("Missing \"Basic \" in Authorization header.")
		aw.Respond(responses.Error{Error: "Malformed Authorization header."}, http.StatusBadRequest)
		return
	}
	authorizationHeader = authorizationHeader[6:] // Remove "Basic " from header.

	type Authorization struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var authorization Authorization
	decodedAuthorizationBytes, err := base64.StdEncoding.DecodeString(authorizationHeader)
	decodedAuthorizationString := string(decodedAuthorizationBytes)
	authorization.Email = decodedAuthorizationString[:strings.Index(decodedAuthorizationString, ":")]
	authorization.Password = decodedAuthorizationString[strings.Index(decodedAuthorizationString, ":")+1:]

	if authorization.Email == "" || authorization.Password == "" {
		log.Println("Malformed authorization header.")
		aw.Respond(responses.Error{Error: "Malformed Authorization header."}, http.StatusBadRequest)
		return
	}

	user, err := mongo.FindUserByEmail(authorization.Email)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not lookup provided email."}, http.StatusInternalServerError)
		return
	}
	if user != nil {
		log.Println("Email already taken.")
		aw.Respond(responses.Error{Error: "Email already taken."}, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authorization.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not generate password hash."}, http.StatusInternalServerError)
		return
	}

	err = mongo.InsertUser(models.User{Email: authorization.Email, PasswordHash: passwordHash, CreationTimestamp: time.Now().Unix()})
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not insert user."}, http.StatusInternalServerError)
		return
	}

	aw.Respond(responses.Message{Message: "User registered successfully."}, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) () {
	ar, aw := models.NewApiCommunication(r, w)

	authorizationHeader := ar.GetHeader("Authorization")
	if authorizationHeader == "" {
		log.Println("Empty/absent \"Authorization\" header value.")
		aw.Respond(responses.Error{Error: "Empty/absent \"Authorization\" header value."}, http.StatusBadRequest)
		return
	}
	if strings.Index(authorizationHeader, "Basic ") != 0 {
		log.Println("Missing \"Basic \" in Authorization header.")
		aw.Respond(responses.Error{Error: "Malformed Authorization header."}, http.StatusBadRequest)
		return
	}
	authorizationHeader = authorizationHeader[6:] // Remove "Basic " from header.

	type Authorization struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var authorization Authorization
	decodedAuthorizationBytes, err := base64.StdEncoding.DecodeString(authorizationHeader)
	decodedAuthorizationString := string(decodedAuthorizationBytes)
	authorization.Email = decodedAuthorizationString[:strings.Index(decodedAuthorizationString, ":")]
	authorization.Password = decodedAuthorizationString[strings.Index(decodedAuthorizationString, ":")+1:]

	if authorization.Email == "" || authorization.Password == "" {
		log.Println("Malformed authorization header.")
		aw.Respond(responses.Error{Error: "Malformed Authorization header."}, http.StatusBadRequest)
		return
	}

	user, err := mongo.FindUserByEmail(authorization.Email)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Failed to validate credentials."}, http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Printf("No user found for provided email \"%v\"\n", authorization.Email)
		aw.Respond(responses.Error{Error: "Invalid email or password."}, http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(authorization.Password))
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Invalid email or password."}, http.StatusBadRequest)
		return
	}

	aw.Respond(responses.Message{Message: "Success."}, http.StatusOK)
}
