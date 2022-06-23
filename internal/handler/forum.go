package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"techpark_db/internal/domain/entity"
	"time"
)

func (h *Handler) ForumCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var forumRequest entity.CreateForum
	if err := json.NewDecoder(r.Body).Decode(&forumRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.storage.GetUser(tx, forumRequest.User)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoUser + forumRequest.User,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}
	forumRequest.User = user.Nickname

	forum, err := h.storage.GetForum(tx, forumRequest.Slug)
	if err == nil {
		tx.Rollback()
		forumBytes, _ := easyjson.Marshal(forum)
		w.WriteHeader(http.StatusConflict)
		w.Write(forumBytes)
		return
	}

	if err := h.storage.SaveForum(tx, forumRequest); err != nil {
		tx.Rollback()
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

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forumBytes, _ := easyjson.Marshal(forum)
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forum, err := h.storage.GetForum(tx, slug)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoForum + slug,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forumBytes, _ := easyjson.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(forumBytes)
}

func (h *Handler) ForumCreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slugForum, ok := vars["slug"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var threadRequest entity.CreateThread
	if err := json.NewDecoder(r.Body).Decode(&threadRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := h.storage.GetUser(tx, threadRequest.Author); err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThreadAuthor + threadRequest.Author,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	thread, err := h.storage.GetThread(tx, threadRequest.Slug)
	if err == nil {
		tx.Rollback()
		threadBytes, _ := easyjson.Marshal(thread)
		w.WriteHeader(http.StatusConflict)
		w.Write(threadBytes)
		return
	}

	forum, err := h.storage.GetForum(tx, slugForum)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThreadForum + slugForum,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if threadRequest.Created == "" {
		threadRequest.Created = time.Now().Format(time.RFC3339Nano)
	}

	insertId, err := h.storage.SaveThread(tx, threadRequest, forum.Slug)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err = h.storage.GetThreadById(tx, insertId)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadBytes, _ := easyjson.Marshal(thread)
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
	since := r.FormValue("since")
	if r.FormValue("limit") != "" {
		limit, _ = strconv.Atoi(r.FormValue("limit"))
	}
	if r.FormValue("desc") == "true" {
		order = "DESC"
	}

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts := r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	users, err := h.storage.GetForumUsers(tx, slug, order, limit, since)
	if err != nil {
		tx.Rollback()
		log.Warning("trouble GetForumUsers")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	if len(*users) == 0 {
		if _, err := h.storage.GetForum(tx, slug); err != nil {
			tx.Rollback()
			resp := &entity.Error{
				Message: ErrNoForum + slug,
			}
			respBytes, _ := easyjson.Marshal(resp)
			w.WriteHeader(http.StatusNotFound)
			w.Write(respBytes)
			return
		}
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usersBytes, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	w.Write(usersBytes)
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

	ts := r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	forum, err := h.storage.GetForumThreads(tx, slug, order, limit, since)
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	if len(*forum) == 0 {
		if _, err := h.storage.GetForum(tx, slug); err != nil {
			tx.Rollback()
			resp := &entity.Error{
				Message: ErrNoForum + slug,
			}
			respBytes, _ := easyjson.Marshal(resp)
			w.WriteHeader(http.StatusNotFound)
			w.Write(respBytes)
			return
		}
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts = r.Context().Value("timestamp").(*[]time.Time)
	*ts = append(*ts, time.Now())

	forumBytes, _ := json.Marshal(forum)
	w.WriteHeader(http.StatusOK)
	w.Write(forumBytes)
}
