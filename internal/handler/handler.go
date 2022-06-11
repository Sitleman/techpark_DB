package handler

import "techpark_db/internal/infra/psql"

type Handler struct {
	storage *psql.Storage
}

func NewNotesHandler(store *psql.Storage) *Handler {
	return &Handler{
		storage: store,
	}
}
