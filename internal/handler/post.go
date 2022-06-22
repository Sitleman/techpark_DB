package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	id, _ := strconv.Atoi(idRaw)
	post, err := h.storage.GetPostById(id)
	if err != nil {
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
			author, err := h.storage.GetUser(post.Author)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DAuthor = author
		case "forum":
			forum, err := h.storage.GetForum(post.Forum)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DForum = forum
		case "thread":
			thread, err := h.storage.GetThreadById(post.Thread)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			postDetails.DThread = thread
		}
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

	post, err := h.storage.GetPostById(id)
	if err != nil {
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

	if err := h.storage.UpdatePost(id, postRequest.Message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post.Message = postRequest.Message
	post.IsEdited = true

	postBytes, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(postBytes)
}
