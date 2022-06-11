package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"techpark_db/internal/domain/entity"
)

var ErrNoUser = entity.Error{
	Message: "Can't find user with id",
}

func (h *Handler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var forumRequest entity.CreateForum
	if err := json.NewDecoder(r.Body).Decode(&forumRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.storage.GetUser(forumRequest.User); err != nil {
		ErrNoUserBytes, _ := json.Marshal(ErrNoUser)
		w.WriteHeader(http.StatusNotFound)
		w.Write(ErrNoUserBytes)
		return
	}

	forum, err := h.storage.GetForum(forumRequest.Slug)
	if err == nil {
		forumBytes, _ := json.Marshal(forum)
		w.WriteHeader(http.StatusConflict)
		w.Write(forumBytes)
		return
	}

	if err := h.storage.SaveForum(forumRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	forum = &entity.Forum{
		Slug:    forumRequest.Slug,
		Title:   forumRequest.Title,
		User:    forumRequest.User,
		Posts:   0,
		Threads: 0,
	}

	forumBytes, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusCreated)
	w.Write(forumBytes)
}

func (h *Handler) ForumDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	forum, err := h.storage.GetForum(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	forumBytes, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(forumBytes)
}

func (h *Handler) ForumCreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var threadRequest entity.CreateThread
	if err := json.NewDecoder(r.Body).Decode(&threadRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.storage.GetUser(threadRequest.Author); err != nil {
		ErrNoUserBytes, _ := json.Marshal(ErrNoUser)
		w.WriteHeader(http.StatusNotFound)
		w.Write(ErrNoUserBytes)
		return
	}

	thread, err := h.storage.GetThreadByTitle(threadRequest.Title)
	if err == nil {
		threadBytes, _ := json.Marshal(thread)
		w.WriteHeader(http.StatusConflict)
		w.Write(threadBytes)
		return
	}

	if err := h.storage.SaveThread(threadRequest, slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err = h.storage.GetThreadByTitle(threadRequest.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	threadBytes, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusCreated)
	w.Write(threadBytes)
}
