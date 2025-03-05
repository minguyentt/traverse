package handlers

import (
	"log/slog"
	"net/http"
	"traverse/api/errors"
	"traverse/api/json"
	"traverse/internal/services"
	"traverse/internal/storage"
	"traverse/models"

	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	RegisterUser(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type authHandler struct {
	service  *services.Service
	validate *validator.Validate
}

func NewAuthHandler(s *services.Service, v *validator.Validate) *authHandler {
	return &authHandler{
		service:  s,
		validate: v,
	}
}

func (u *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.RegistrationPayload
	// Use RegistrationPayload json struct

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
	user, err := u.service.Users.Register(r.Context(), &userPayload)
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
	var payload models.UserLoginPayload
	if err := json.Read(w, r, &payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	if err := u.validate.Struct(payload); err != nil {
		errors.BadRequestResponse(w, r, err)
		return
	}

	slog.Info("login handler", "payload", &payload)

	userToken, err := u.service.Users.Login(r.Context(), &payload)
	if err != nil {
		errors.UnauthorizedErr(w, r, err)
		return
	}

	if err := json.Response(w, http.StatusOK, userToken); err != nil {
		errors.InternalServerErr(w, r, err)
	}
}
