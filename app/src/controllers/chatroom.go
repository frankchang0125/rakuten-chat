package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"rakuten.co.jp/chatroom/services"
)

type createChatroomResp struct {
	ID int `json:"id"`
}

type getChatroomsResp struct {
	IDs []int `json:"ids"`
}

type getChatroomUsersResp struct {
	Users []string `json:"users"`
}

func handleCreateChatroom(w http.ResponseWriter, r *http.Request) {
	id, err := services.CreateChatroom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := createChatroomResp{
		ID: id,
	}
	returnOK(w, resp)
}

func handleGetChatrooms(w http.ResponseWriter, r *http.Request) {
	ids, err := services.GetChatroomsList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := getChatroomsResp{
		IDs: ids,
	}
	returnOK(w, resp)
}

func handleGetChatroomUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatroomIDStr := vars["chatroom-id"]
	chatroomID, err := strconv.Atoi(chatroomIDStr)
	if err != nil {
		returnBadRequest(w, errorResponse{
			Message: fmt.Sprintf("Invaid chat room ID: %s", chatroomIDStr),
		})
		return
	}

	users, err := services.GetChatroomUsers(chatroomID)
	if err != nil {
		switch err {
		case services.ErrChatroomNotExists:
			returnBadRequest(w, errorResponse{
				Message: err.Error(),
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := getChatroomUsersResp{
		Users: users,
	}
	returnOK(w, resp)
}

func handleJoinChatroom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatroomIDStr := vars["chatroom-id"]
	chatroomID, err := strconv.Atoi(chatroomIDStr)
	if err != nil {
		returnBadRequest(w, errorResponse{
			Message: fmt.Sprintf("Invaid chat room ID: %s", chatroomIDStr),
		})
		return
	}

	userName := r.URL.Query().Get("user")
	if userName == "" {
		returnBadRequest(w, errorResponse{
			Message: "Missing 'user' parameter",
		})
		return
	}

	services.JoinChatroom(w, r, chatroomID, userName)
}
