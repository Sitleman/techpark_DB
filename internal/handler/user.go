package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := h.storage.FindUser(tx, nickname, userReq.Email)
	if err == nil && len(*users) > 0 {
		tx.Rollback()
		usersBytes, _ := json.Marshal(users)
		w.WriteHeader(http.StatusConflict)
		w.Write(usersBytes)
		return
	}

	if err := h.storage.SaveUser(tx, userReq, nickname); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		Nickname: nickname,
		Fullname: userReq.Fullname,
		About:    userReq.About,
		Email:    userReq.Email,
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.storage.GetUser(tx, nickname)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoUser + nickname,
		}
		errorBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(errorBytes)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.storage.GetUser(tx, nickname)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoUser + nickname,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if userReq.Fullname == "" {
		userReq.Fullname = user.Fullname
	}
	if userReq.About == "" {
		userReq.About = user.About
	}
	if userReq.Email == "" {
		userReq.Email = user.Email
	}

	users, err := h.storage.FindUser(tx, nickname, userReq.Email)
	if err == nil && len(*users) > 1 {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrEmailAlreadyRegistered + nickname,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusConflict)
		w.Write(respBytes)
		return
	}

	if err := h.storage.UpdateUser(tx, userReq, nickname); err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrEmailAlreadyRegistered + nickname,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusConflict)
		w.Write(respBytes)
		return
	}

	user = &entity.User{
		Nickname: nickname,
		Fullname: userReq.Fullname,
		About:    userReq.About,
		Email:    userReq.Email,
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userBytes, _ := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return
}
