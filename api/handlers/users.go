package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/internal/services"
)

type UsersHandler interface {
	UserByID(w http.ResponseWriter, r *http.Request)
}

type usersHandler struct {
	service *services.UserService
}

func (u *usersHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		// handle server error when failed to parse url
		errors.BadRequestResponse(w, r, err)
		return
	}

	// call the service layer
	user, err := u.service.ByID(r.Context(), userID)
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
