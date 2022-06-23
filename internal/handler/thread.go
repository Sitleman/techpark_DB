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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread, err := h.storage.GetThread(tx, slug_or_id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if len(postReq) == 0 {
		tx.Rollback()
		postsBytes, _ := json.Marshal(postReq)
		w.WriteHeader(http.StatusCreated)
		w.Write(postsBytes)
		return
	}

	if _, err := h.storage.GetUser(tx, postReq[0].Author); err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoPostAuthor + postReq[0].Author,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	if postReq[0].Parent != 0 {
		ok, err := h.storage.CheckParentPost(tx, postReq[0].Parent, thread.Id)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			tx.Rollback()
			resp := &entity.Error{
				Message: ErrNoThread + slug_or_id,
			}
			respBytes, _ := easyjson.Marshal(resp)
			w.WriteHeader(http.StatusConflict)
			w.Write(respBytes)
			return
		}
	}

	created := time.Now().Format(time.RFC3339Nano)
	created = created[:len(created)-4]

	ids, err := h.storage.SavePosts(tx, postReq, thread.Forum, thread.Id, created)
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

	posts := make([]entity.Post, len(*ids))
	for i := 0; i < len(*ids); i++ {
		posts[i] = entity.Post{
			Id:       (*ids)[i],
			Parent:   postReq[i].Parent,
			Author:   postReq[i].Author,
			Message:  postReq[i].Message,
			Created:  created,
			IsEdited: false,
			Forum:    thread.Forum,
			Thread:   thread.Id,
		}
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread, err := h.storage.GetThread(tx, slug_or_id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}
	voteReq.IdThread = thread.Id

	user, err := h.storage.GetUser(tx, voteReq.Nickname)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoUser + voteReq.Nickname,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}
	voteReq.Nickname = user.Nickname

	err = h.storage.SetVote(tx, voteReq)
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusConflict)
		return
	}

	voteCount, err := h.storage.CountVote(tx, thread.Id)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	thread.Votes = *voteCount

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadBytes, _ := easyjson.Marshal(thread)
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread, err := h.storage.GetThread(tx, slug_or_id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
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

	threadBytes, _ := easyjson.Marshal(thread)
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread, err := h.storage.GetThread(tx, slug_or_id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := easyjson.Marshal(resp)
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

	if err := h.storage.UpdateThread(tx, *thread); err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	threadBytes, _ := easyjson.Marshal(thread)
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

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread, err := h.storage.GetThread(tx, slug_or_id)
	if err != nil {
		tx.Rollback()
		resp := &entity.Error{
			Message: ErrNoThread + slug_or_id,
		}
		respBytes, _ := easyjson.Marshal(resp)
		w.WriteHeader(http.StatusNotFound)
		w.Write(respBytes)
		return
	}

	var posts *[]entity.Post
	switch sort {
	case "flat":
		posts, err = h.storage.GetPostsByThreadFlat(tx, thread.Id, limit, since, sort, order)
	case "tree":
		posts, err = h.storage.GetPostsTree(tx, thread.Id, limit, since, sort, order)
	case "parent_tree":
		posts, err = h.storage.GetPostsParentTree(tx, thread.Id, limit, since, sort, order)
	}

	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postsBytes, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(postsBytes)
}

//
//user, err := h.storage.GetUser(voteReq.Nickname)
//if err != nil {
//	resp := &entity.Error{
//		Message: ErrNoUser + voteReq.Nickname,
//	}
//	respBytes, _ := json.Marshal(resp)
//	w.WriteHeader(http.StatusNotFound)
//	w.Write(respBytes)
//	return
//}
//voteReq.Nickname = user.Nickname
//
//voteReq.IdThread = thread.Id
//voteReq.SlugThread = thread.Slug
//
//vote, err := h.storage.GetVote(voteReq.IdThread, voteReq.Nickname)
//if err != nil {
//	h.storage.SaveVote(voteReq)
//	thread.Votes += voteReq.Voice
//} else {
//	h.storage.UpdateVote(voteReq)
//	thread.Votes += voteReq.Voice - vote.Voice
//}
//
//if err := h.storage.UpdateThreadVote(*thread); err != nil {
//	w.WriteHeader(http.StatusInternalServerError)
//	return
//}
