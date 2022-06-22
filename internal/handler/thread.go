package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"techpark_db/internal/domain/entity"
	"time"
)

func (h *Handler) ThreadCreatePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug_or_id, ok := vars["slug_or_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var postReq []entity.CreatePost
	if err := json.NewDecoder(r.Body).Decode(&postReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := h.storage.GetThread(slug_or_id)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	for _, post := range postReq {
		if _, err := h.storage.GetUser(post.Author); err != nil {
			resp := &entity.Error{
				Message: ErrNoPostAuthor + post.Author,
			}
			respBytes, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusNotFound)
			w.Write(respBytes)
			return
		}

		if post.Parent != 0 {
			ok, err := h.storage.CheckParentPost(post.Parent, thread.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				resp := &entity.Error{
					Message: ErrNoThread + slug_or_id,
				}
				respBytes, _ := json.Marshal(resp)
				w.WriteHeader(http.StatusConflict)
				w.Write(respBytes)
				return
			}
		}
	}

	created := time.Now().Format(time.RFC3339Nano)

	if err := h.storage.SavePosts(postReq, thread.Forum, thread.Id, created); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := h.storage.GetPostsByCreated(created)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	postsBytes, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusCreated)
	w.Write(postsBytes)
}

func (h *Handler) ThreadVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug_or_id, ok := vars["slug_or_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var voteReq entity.Vote
	if err := json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := h.storage.GetThread(slug_or_id)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	user, err := h.storage.GetUser(voteReq.Nickname)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoUser + voteReq.Nickname,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}
	voteReq.Nickname = user.Nickname

	voteReq.IdThread = thread.Id
	voteReq.SlugThread = thread.Slug

	vote, err := h.storage.GetVote(voteReq.IdThread, voteReq.Nickname)
	if err != nil {
		h.storage.SaveVote(voteReq)
		thread.Votes += voteReq.Voice
	} else {
		h.storage.UpdateVote(voteReq)
		thread.Votes += voteReq.Voice - vote.Voice
	}

	if err := h.storage.UpdateThreadVote(*thread); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadBytes, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(threadBytes)
}

func (h *Handler) ThreadDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug_or_id, ok := vars["slug_or_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	thread, err := h.storage.GetThread(slug_or_id)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	threadBytes, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(threadBytes)
}

func (h *Handler) ThreadUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug_or_id, ok := vars["slug_or_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var threadReq entity.Thread
	if err := json.NewDecoder(r.Body).Decode(&threadReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := h.storage.GetThread(slug_or_id)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if threadReq.Title != "" {
		thread.Title = threadReq.Title
	}
	if threadReq.Message != "" {
		thread.Message = threadReq.Message
	}

	if err := h.storage.UpdateThread(*thread); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadBytes, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(threadBytes)
}

func (h *Handler) ThreadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	slug_or_id, ok := vars["slug_or_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := DEFAULT_LIMIT
	since := DEFAULT_SINCE_ID
	sort := DEFAUTL_SORT
	order := DEFAULT_ORDER
	if r.FormValue("limit") != "" {
		limit, _ = strconv.Atoi(r.FormValue("limit"))
	}
	if r.FormValue("since") != "" {
		since, _ = strconv.Atoi(r.FormValue("since"))
	}
	if r.FormValue("sort") != "" {
		sort = r.FormValue("sort")
	}
	if r.FormValue("desc") == "true" {
		order = "DESC"
	}

	thread, err := h.storage.GetThread(slug_or_id)
	if err != nil {
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	var posts *[]entity.Post
	switch sort {
	case "flat":
		posts, err = h.storage.GetPostsByThreadFlat(thread.Id, limit, since, sort, order)
	case "tree":
		posts, err = h.storage.GetPostsTree(thread.Id, limit, since, sort, order)
	case "parent_tree":
		posts, err = h.storage.GetPostsParentTree(thread.Id, limit, since, sort, order)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postsBytes, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(postsBytes)
}
