package api

import (
	"arby-user-api/pkg/models"
	"net/http"
	"arby-user-api/pkg/models/responses"
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
