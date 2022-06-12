package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"techpark_db/internal/domain/entity"
	"time"
)

const (
	DEFAULT_ORDER    = "ASC"
	DEFAULT_LIMIT    = 100
	DEFAULT_SINCE_ID = 0
)

var DEFAULT_SINCE_DATA_MIN = time.Date(1000, 00, 0, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
var DEFAULT_SINCE_DATA_MAX = time.Date(4000, 00, 0, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

func (h *Handler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var forumRequest entity.CreateForum
	if err := json.NewDecoder(r.Body).Decode(&forumRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.storage.GetUser(forumRequest.User)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoUser + forumRequest.User,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}
	forumRequest.User = user.Nickname

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
		resp := &entity.Error{
			Message: ErrNoForum + slug,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
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

	if threadRequest.Created == "" {
		threadRequest.Created = time.Now().Format(time.RFC3339)
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

func (h *Handler) ForumUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := DEFAULT_ORDER
	limit := DEFAULT_LIMIT
	since := DEFAULT_SINCE_ID
	if r.FormValue("limit") != "" {
		limit, _ = strconv.Atoi(r.FormValue("limit"))
	}
	if r.FormValue("desc") == "true" {
		order = "DESC"
	}
	if r.FormValue("since") != "" {
		since, _ = strconv.Atoi(r.FormValue("since"))
	}

	forum, err := h.storage.GetForumUsers(slug, order, limit, since)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	forumBytes, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(forumBytes)
}

func (h *Handler) ForumThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := DEFAULT_ORDER
	limit := DEFAULT_LIMIT
	since := DEFAULT_SINCE_DATA_MIN
	if r.FormValue("limit") != "" {
		limit, _ = strconv.Atoi(r.FormValue("limit"))
	}
	if r.FormValue("desc") == "true" {
		order = "DESC"
		since = DEFAULT_SINCE_DATA_MAX
	}
	if r.FormValue("since") != "" {
		since = r.FormValue("since")
	}

	if _, err := h.storage.GetForum(slug); err != nil {
		resp := &entity.Error{
			Message: ErrNoForum + slug,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	forum, err := h.storage.GetForumThreads(slug, order, limit, since)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forumBytes, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(forumBytes)
}
