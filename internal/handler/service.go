package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	servStatus, err := h.storage.GetServiceStatus(tx)
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

	servStatusBytes, _ := json.Marshal(servStatus)
	w.WriteHeader(http.StatusOK)
	w.Write(servStatusBytes)
}

func (h *Handler) ServiceClear(w http.ResponseWriter, r *http.Request) {
	tx, err := h.storage.DB.Begin()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.storage.ClearData(tx); err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
