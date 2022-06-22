package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"techpark_db/internal/domain/entity"
)

func (h *Handler) PostGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	idRaw, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	argsRaw := r.FormValue("related")
	args := strings.Split(argsRaw, ",")
	//log.Info(args, "///", argsRaw, "///")

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(idRaw)
	post, err := h.storage.GetPostById(tx, id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoPost + idRaw,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	var postDetails entity.PostDetails
	postDetails.DPost = post

	for _, arg := range args {
		switch arg {
		case "user":
			author, err := h.storage.GetUser(tx, post.Author)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DAuthor = author
		case "forum":
			forum, err := h.storage.GetForum(tx, post.Forum)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DForum = forum
		case "thread":
			thread, err := h.storage.GetThreadById(tx, post.Thread)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DThread = thread
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postDetailsBytes, _ := json.Marshal(postDetails)
	w.WriteHeader(http.StatusOK)
	w.Write(postDetailsBytes)
}

func (h *Handler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	idRaw, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//log.Info(r.FormValue("related"))
	id, _ := strconv.Atoi(idRaw)

	var postRequest entity.UpdatePost
	if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, err := h.storage.GetPostById(tx, id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoPost + idRaw,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if postRequest.Message == "" || postRequest.Message == post.Message {
		postWithoutEdited := entity.PostWithoutEdited{
			Id:      post.Id,
			Parent:  post.Parent,
			Author:  post.Author,
			Message: post.Message,
			Forum:   post.Forum,
			Thread:  post.Thread,
			Created: post.Created,
		}

		postWithoutEditedBytes, _ := json.Marshal(postWithoutEdited)
		w.WriteHeader(http.StatusOK)
		w.Write(postWithoutEditedBytes)
		return
	}

	if err := h.storage.UpdatePost(tx, id, postRequest.Message); err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post.Message = postRequest.Message
	post.IsEdited = true

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postBytes, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(postBytes)
}
