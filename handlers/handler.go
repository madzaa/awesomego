package handlers

import (
	"awesomeProject/service"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/user/") {
		switch r.Method {
		case http.MethodGet:
			h.GetBalance(w, r)
			return
		case http.MethodPost:
			h.CreateTransaction(w, r)
			return
		default:
			h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}
	h.respondWithError(w, http.StatusNotFound, "not found")
}
