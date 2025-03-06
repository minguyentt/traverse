package handlers

import (
	"net/http"
	"strconv"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/internal/services"

	"github.com/go-chi/chi/v5"
)

type UsersHandler interface {
	UserByIDHandler(w http.ResponseWriter, r *http.Request)
}

type usersHandler struct {
	service services.UserService
}

func NewUserHandler(services services.UserService) *usersHandler {
    return &usersHandler{service: services}
}

func (h *usersHandler) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		// handle server error when failed to parse url
		errors.BadRequestResponse(w, r, err)
		return
	}

	// call the service layer
	user, err := h.service.UserByID(r.Context(), userID)
	if err != nil {
		// handle either internal server error AND no user found in DB
		errors.InternalServerErr(w, r, err)
		return
	}

	// respond with JSON for any errors that happened
	if err := json.Response(w, http.StatusOK, user); err != nil {
		return
	}
}
