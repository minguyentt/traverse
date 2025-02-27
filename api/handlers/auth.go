package handlers

import (
	"net/http"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/api/models"
	"traverse/internal/services"
	"traverse/internal/storage"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	RegisterUser(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type authHandler struct {
	service  *services.UserService
	validate *validator.Validate
}

func (u *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Use RegistrationPayload json struct
	var userPayload models.RegistrationPayload

	// read the HTTP request
	err := json.Read(w, r, &userPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}
	// validate the json struct
	err = u.validate.Struct(userPayload)
	if err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	// call the service for user creation
	user, err := u.service.RegisterUser(r.Context(), &userPayload)
	if err != nil {
		switch err {
		case storage.ErrDuplicateUsername:
			errors.BadRequestResponse(w, r, err)
			return
		default:
			errors.InternalServerErr(w, r, err)
		}
	}

	if err := json.Response(w, http.StatusCreated, user); err != nil {
		errors.InternalServerErr(w, r, err)
	}
}

func (u *authHandler) Login(w http.ResponseWriter, r *http.Request) {
}
