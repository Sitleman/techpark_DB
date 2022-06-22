package handler

import (
	"techpark_db/internal/infra/psql"
	"time"
)

const (
	DEFAULT_ORDER      = "ASC"
	DEFAULT_LIMIT      = 100
	DEFAULT_SINCE_ID   = 0
	DEFAULT_SINCE_ASC  = ""
	DEFAULT_SINCE_DESC = "ZZZZZZZZZZZZZZ"
	DEFAUTL_SORT       = "flat"
)

var DEFAULT_SINCE_DATA_MIN = time.Date(1000, 00, 0, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
var DEFAULT_SINCE_DATA_MAX = time.Date(4000, 00, 0, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

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
var ErrNoThread = "Can't find thread with slug: "
var ErrNoThreadAuthor = "Can't find thread author by nickname: "
var ErrNoThreadForum = "Can't find thread forum by slug: "
var ErrNoPost = "Can't find post by id: "
var ErrNoPostAuthor = "Can't find post author by nickname: "
