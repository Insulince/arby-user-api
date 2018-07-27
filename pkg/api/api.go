package api

import (
	"arby-user-api/pkg/models"
	"net/http"
	"arby-user-api/pkg/models/responses"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
	"arby-user-api/pkg/mongo"
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

	rawPostBody, err := ar.GetRequestBody()
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not process request."}, http.StatusInternalServerError)
		return
	}
	type PostBody struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}
	var postBody PostBody
	err = json.Unmarshal(rawPostBody, &postBody)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not read request body."}, http.StatusBadRequest)
		return
	}
	if postBody.Username == nil || postBody.Password == nil {
		log.Println("Malformed request body.")
		aw.Respond(responses.Error{Error: "Malformed request body."}, http.StatusBadRequest)
		return
	}

	user, err := mongo.FindUserByUsername(*postBody.Username)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not lookup provided username."}, http.StatusInternalServerError)
		return
	}
	if user != nil {
		log.Println("Username already taken.")
		aw.Respond(responses.Error{Error: "Username already taken."}, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*postBody.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not generate password hash."}, http.StatusInternalServerError)
		return
	}

	err = mongo.InsertUser(models.User{Username: *postBody.Username, PasswordHash: passwordHash, CreationTimestamp: time.Now().Unix()})
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not insert user."}, http.StatusInternalServerError)
		return
	}

	aw.Respond(responses.Message{Message: "User registered successfully."}, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) () {
	ar, aw := models.NewApiCommunication(r, w)

	rawPostBody, err := ar.GetRequestBody()
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not process request."}, http.StatusInternalServerError)
		return
	}
	type PostBody struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
	}
	var postBody PostBody
	err = json.Unmarshal(rawPostBody, &postBody)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Could not read request body."}, http.StatusBadRequest)
		return
	}
	if postBody.Username == nil || postBody.Password == nil {
		log.Println("Malformed request body.")
		aw.Respond(responses.Error{Error: "Malformed request body."}, http.StatusBadRequest)
		return
	}

	user, err := mongo.FindUserByUsername(*postBody.Username)
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Failed to validate credentials."}, http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Printf("No user found for provided username \"%v\"\n", *postBody.Username)
		aw.Respond(responses.Error{Error: "Invalid username or password."}, http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(*postBody.Password))
	if err != nil {
		log.Println(err)
		aw.Respond(responses.Error{Error: "Invalid username or password."}, http.StatusBadRequest)
		return
	}

	aw.Respond(responses.Message{Message: "Success."}, http.StatusOK)
}
