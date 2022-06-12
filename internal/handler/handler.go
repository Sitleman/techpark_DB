package handler

import "techpark_db/internal/infra/psql"

type Handler struct {
	storage *psql.Storage
}

func NewHandler(store *psql.Storage) *Handler {
	return &Handler{
		storage: store,
	}
}

var ErrNoUser = "Can't find user by nickname: "
var ErrEmailAlreadyRegistered = "This email is already registered by user: "
var ErrNoForum = "Can't find forum with slug: "
