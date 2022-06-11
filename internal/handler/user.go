package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"techpark_db/internal/domain/entity"
)

func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userReq entity.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := h.storage.FindUser(nickname, userReq.Email)
	if err == nil && len(*users) > 0 {
		usersBytes, _ := json.Marshal(users)
		w.WriteHeader(http.StatusConflict)
		w.Write(usersBytes)
		return
	}

	if err := h.storage.SaveUser(userReq, nickname); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		Nickname: nickname,
		Fullname: userReq.Fullname,
		About:    userReq.About,
		Email:    userReq.Email,
	}

	forumBytes, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(forumBytes)
}

func (h *Handler) UserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.storage.GetUser(nickname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userBytes, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userReq entity.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := h.storage.FindUser(nickname, userReq.Email)
	if err == nil && len(*users) > 1 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err := h.storage.UpdateUser(userReq, nickname); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user := &entity.User{
		Nickname: nickname,
		Fullname: userReq.Fullname,
		About:    userReq.About,
		Email:    userReq.Email,
	}

	userBytes, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return
}
