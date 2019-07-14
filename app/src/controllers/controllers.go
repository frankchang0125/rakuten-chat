package controllers

import (
	"encoding/json"
	"net/http"

	"rakuten.co.jp/chatroom/routes"
)

var Routes []*routes.Route

func init() {
	Routes = []*routes.Route{
		// chatroom
		routes.NewRoute(http.MethodPost, "chatrooms", handleCreateChatroom),
		routes.NewRoute(http.MethodGet, "chatrooms", handleGetChatrooms),
		routes.NewRoute(http.MethodGet, "chatrooms/{chatroom-id}/users", handleGetChatroomUsers),
		routes.NewRoute(http.MethodGet, "chatrooms/{chatroom-id}/join", handleJoinChatroom),

		// users
		routes.NewRoute(http.MethodPost, "users", handleCreateUser),
		routes.NewRoute(http.MethodDelete, "users", handleDeleteAllUsers),
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func returnOK(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeResponseJSON(w, resp)
}

func returnBadRequest(w http.ResponseWriter, errResp errorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	writeResponseJSON(w, errResp)
}

func returnForbiddenRequest(w http.ResponseWriter, errResp errorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	writeResponseJSON(w, errResp)
}

func returnNotFoundRequest(w http.ResponseWriter, errResp errorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	writeResponseJSON(w, errResp)
}

func returnInternalServerError(w http.ResponseWriter, errResp errorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	writeResponseJSON(w, errResp)
}

func writeResponseJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
