package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	servStatus, err := h.storage.GetServiceStatus()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	servStatusBytes, _ := json.Marshal(servStatus)
	w.WriteHeader(http.StatusOK)
	w.Write(servStatusBytes)
}

func (h *Handler) ServiceClear(w http.ResponseWriter, r *http.Request) {

	if err := h.storage.ClearData(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
