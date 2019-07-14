package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"rakuten.co.jp/chatroom/services"
)

type createUserReq struct {
	Name string `json:"name"`
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	req := createUserReq{}
	err := json.Unmarshal(requestBody, &req)
	if err != nil {
		log.WithError(err).Error("Decode json body failed")
		returnBadRequest(w, errorResponse{
			Message: "Invalid request body",
		})
		return
	}

	err = services.CreateUser(req.Name)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err := services.DeleteAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
