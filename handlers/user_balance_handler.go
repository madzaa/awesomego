package handlers

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/user/")
	id = strings.TrimSuffix(id, "/balance")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	ctx := r.Context()
	balance, err := h.service.GetBalance(ctx, userId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := models.UserBalance{
		ID:      userId,
		Balance: utils.FormatBalance(balance),
	}
	h.respondWithJSON(w, http.StatusOK, response)
}
