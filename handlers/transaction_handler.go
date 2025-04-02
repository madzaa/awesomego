package handlers

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var transaction models.TransactionRequest
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/user/")
	id = strings.TrimSuffix(id, "/transaction")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid body")
		return
	}
	source := r.Header.Get("Source-Type")
	if source == "" {
		h.respondWithError(w, http.StatusBadRequest, "header not set")
		return
	}
	amount, err := utils.ParseAmount(transaction.Amount)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.ProcessTransaction(ctx, transaction.TransactionId, userId, amount, transaction.State, source)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
